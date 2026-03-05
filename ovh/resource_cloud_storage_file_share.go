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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageFileShareResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageFileShareResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageFileShareResource)(nil)
)

func NewCloudStorageFileShareResource() resource.Resource {
	return &cloudStorageFileShareResource{}
}

type cloudStorageFileShareResource struct {
	config *Config
}

func (r *cloudStorageFileShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share"
}

func (r *cloudStorageFileShareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var fileShareMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description"},
	Int64s:  []string{"size"},
	Lists:   []string{"access_rules"},
}

func (r *cloudStorageFileShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a file storage share (NFS) in a public cloud project using the publicCloud API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "File share name",
				MarkdownDescription: "File share name",
			},
			"size": schema.Int64Attribute{
				Required:            true,
				Description:         "Size of the file share in GB",
				MarkdownDescription: "Size of the file share in GB",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the file share will be created",
				MarkdownDescription: "Region where the file share will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "File share protocol (NFS)",
				MarkdownDescription: "File share protocol (`NFS`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"share_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "File share type (e.g. STANDARD_1AZ)",
				MarkdownDescription: "File share type (e.g. `STANDARD_1AZ`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Network ID to attach the file share to",
				MarkdownDescription: "Network ID to attach the file share to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Subnet ID to attach the file share to",
				MarkdownDescription: "Subnet ID to attach the file share to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"availability_zone": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "The availability zone where the file share will be created",
				MarkdownDescription: "The availability zone where the file share will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "File share description",
				MarkdownDescription: "File share description",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_rules": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Access rules for the file share",
				MarkdownDescription: "Access rules for the file share",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_to": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "IP address or CIDR to grant access to",
							MarkdownDescription: "IP address or CIDR to grant access to",
						},
						"access_level": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Access level (READ_WRITE, READ_ONLY)",
							MarkdownDescription: "Access level (`READ_WRITE`, `READ_ONLY`)",
						},
					},
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "File share ID",
				MarkdownDescription: "File share ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileShareMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the file share",
				MarkdownDescription: "Creation date of the file share",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the file share",
				MarkdownDescription: "Last update date of the file share",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileShareMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "File share readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "File share readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the file storage share",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(fileShareMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Current location",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Region",
								MarkdownDescription: "Region",
							},
							"availability_zone": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Availability zone",
								MarkdownDescription: "Availability zone",
							},
						},
					},
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "File share name",
						MarkdownDescription: "File share name",
					},
					"description": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "File share description",
						MarkdownDescription: "File share description",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						Description:         "Size of the file share in GB",
						MarkdownDescription: "Size of the file share in GB",
					},
					"protocol": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "File share protocol",
						MarkdownDescription: "File share protocol",
					},
					"share_type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "File share type",
						MarkdownDescription: "File share type",
					},
					"network_id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Network ID",
						MarkdownDescription: "Network ID",
					},
					"subnet_id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Subnet ID",
						MarkdownDescription: "Subnet ID",
					},
					"export_locations": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Export locations for the file share",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"path": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Export path",
									MarkdownDescription: "Export path",
								},
								"preferred": schema.BoolAttribute{
									Computed:            true,
									Description:         "Whether this is the preferred export location",
									MarkdownDescription: "Whether this is the preferred export location",
								},
							},
						},
					},
					"access_rules": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Current access rules for the file share",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Access rule ID",
									MarkdownDescription: "Access rule ID",
								},
								"access_to": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "IP address or CIDR",
									MarkdownDescription: "IP address or CIDR",
								},
								"access_level": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Access level",
									MarkdownDescription: "Access level",
								},
								"state": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Access rule state",
									MarkdownDescription: "Access rule state",
								},
								"created_at": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Access rule creation date",
									MarkdownDescription: "Access rule creation date",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *cloudStorageFileShareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<file_share_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudStorageFileShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageFileShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate(ctx)

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file"

	var responseData CloudStorageFileShareAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for file share to be READY
	_, err := r.waitForFileShareReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for file storage share to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/" + url.PathEscape(responseData.Id)
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

func (r *cloudStorageFileShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageFileShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageFileShareAPIResponse
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

func (r *cloudStorageFileShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudStorageFileShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(ctx, data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageFileShareAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for file share to be READY
	_, err := r.waitForFileShareReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for file storage share to be ready after update",
			err.Error(),
		)
		return
	}

	// Read final state
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

func (r *cloudStorageFileShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageFileShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/" + url.PathEscape(data.Id.ValueString())

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
		Refresh: func() (any, string, error) {
			res := &CloudStorageFileShareAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for file storage share to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudStorageFileShareResource) waitForFileShareReady(ctx context.Context, serviceName, fileShareId string) (any, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudStorageFileShareAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/file/" + url.PathEscape(fileShareId)
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
