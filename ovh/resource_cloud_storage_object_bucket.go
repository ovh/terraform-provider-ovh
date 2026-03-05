package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketResource)(nil)
)

func NewCloudStorageObjectBucketResource() resource.Resource {
	return &cloudStorageObjectBucketResource{}
}

type cloudStorageObjectBucketResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket"
}

func (r *cloudStorageObjectBucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

var bucketMutableAttrs = MutableAttrs{
	Strings: []string{"name", "region", "owner_user_id"},
	Objects: []string{"encryption", "versioning", "object_lock"},
	Maps:    []string{"tags"},
}

func (r *cloudStorageObjectBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an S3 object storage bucket in a public cloud project using the v2 API.",
		Attributes: map[string]schema.Attribute{
			// Required attributes
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Bucket name (must be globally unique and DNS-compatible)",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Region where the bucket will be created",
			},

			// Optional attributes
			"encryption": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Server-side encryption configuration",
				Attributes: map[string]schema.Attribute{
					"algorithm": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Encryption algorithm (e.g. AES256)",
					},
				},
			},
			"versioning": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Versioning configuration",
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Versioning status (DISABLED, ENABLED, SUSPENDED)",
					},
				},
			},
			"object_lock": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Object lock (WORM) configuration; requires versioning to be enabled",
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Object lock retention mode (COMPLIANCE, GOVERNANCE)",
					},
					"retention_days": schema.Int64Attribute{
						Required:    true,
						Description: "Number of days to retain objects",
					},
					"retention_years": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of years to retain objects (alternative to retention_days)",
					},
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Metadata tags for the bucket",
			},
			"owner_user_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Owner user identifier",
			},

			// Computed attributes
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Bucket ID (same as bucket name)",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(bucketMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date of the bucket",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the bucket",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(bucketMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Bucket readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},

			// Current state (computed, read-only)
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the bucket",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(bucketMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Bucket name",
					},
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Location details",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Region code",
							},
						},
					},
					"encryption": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Encryption configuration",
						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Encryption algorithm",
							},
						},
					},
					"versioning": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Versioning configuration",
						Attributes: map[string]schema.Attribute{
							"status": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Versioning status",
							},
						},
					},
					"object_lock": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Object lock configuration",
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Object lock retention mode",
							},
							"retention_days": schema.Int64Attribute{
								Computed:    true,
								Description: "Retention period in days",
							},
							"retention_years": schema.Int64Attribute{
								Computed:    true,
								Description: "Retention period in years",
							},
						},
					},
					"tags": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "Metadata tags",
					},
				},
			},
		},
	}
}

func (r *cloudStorageObjectBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_name>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudStorageObjectBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate(ctx)

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket"

	var responseData CloudBucketAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for bucket to be READY
	_, err := r.waitForReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for bucket to be ready",
			err.Error(),
		)
		return
	}

	// Read the final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(responseData.Id)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudBucketAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudStorageObjectBucketModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(ctx, data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudBucketAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for bucket to be READY
	_, err := r.waitForReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for bucket to be ready after update",
			err.Error(),
		)
		return
	}

	// Read the final state
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	planData.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudStorageObjectBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for deletion to complete
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudBucketAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for bucket to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudStorageObjectBucketResource) waitForReady(ctx context.Context, serviceName, bucketName string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudBucketAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/object/bucket/" + url.PathEscape(bucketName)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
