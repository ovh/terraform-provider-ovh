package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudS3BucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudS3BucketResource)(nil)
	_ resource.ResourceWithImportState = (*cloudS3BucketResource)(nil)
)

func NewCloudS3BucketResource() resource.Resource {
	return &cloudS3BucketResource{}
}

type cloudS3BucketResource struct {
	config *Config
}

func (r *cloudS3BucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_s3_bucket"
}

func (r *cloudS3BucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var s3BucketMutableAttrs = MutableAttrs{
	Strings: []string{"owner_user_id"},
	Maps:    []string{"tags"},
	Objects: []string{"encryption", "versioning"},
}

type s3BucketObjectLockRequiresVersioningEnabled struct{}

func (v s3BucketObjectLockRequiresVersioningEnabled) Description(_ context.Context) string {
	return "object_lock requires versioning.status to be ENABLED"
}

func (v s3BucketObjectLockRequiresVersioningEnabled) MarkdownDescription(ctx context.Context) string {
	return "`object_lock` requires `versioning.status` to be `ENABLED`"
}

func (v s3BucketObjectLockRequiresVersioningEnabled) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Also covers objectvalidator.AlsoRequires, which is not vendored.
	var versioning types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("versioning"), &versioning)...)
	if resp.Diagnostics.HasError() || versioning.IsUnknown() {
		return
	}
	if versioning.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid object_lock configuration",
			"object_lock requires versioning to be set with status \"ENABLED\".",
		)
		return
	}

	var status ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("versioning").AtName("status"), &status)...)
	if resp.Diagnostics.HasError() || status.IsNull() || status.IsUnknown() {
		return
	}

	if status.ValueString() != "ENABLED" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid object_lock configuration",
			fmt.Sprintf("object_lock requires versioning.status to be \"ENABLED\", got %q.", status.ValueString()),
		)
	}
}

func (r *cloudS3BucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an S3 compatible object storage bucket in a public cloud project using the publicCloud API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers: []planmodifier.String{
					EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Bucket name (must be globally unique and DNS-compatible)",
				MarkdownDescription: "Bucket name (must be globally unique and DNS-compatible)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`),
						"must be lowercase alphanumeric, dots and hyphens, not starting or ending with a hyphen",
					),
				},
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Upper-case region identifier where the bucket will be created (e.g. GRA, SBG, BHS)",
				MarkdownDescription: "Upper-case region identifier where the bucket will be created (e.g. `GRA`, `SBG`, `BHS`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					// The API upper-cases the region; a lower-case value would never match state.
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Z0-9-]+$`),
						"must be upper-case (e.g. \"GRA\")",
					),
				},
			},
			"owner_user_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Owner user identifier",
				MarkdownDescription: "Owner user identifier",
			},
			"tags": schema.MapAttribute{
				ElementType:         ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Metadata tags for the bucket",
				MarkdownDescription: "Metadata tags for the bucket",
			},
			"encryption": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Server-side encryption configuration",
				MarkdownDescription: "Server-side encryption configuration",
				Attributes: map[string]schema.Attribute{
					"algorithm": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Encryption algorithm (AES256, PLAINTEXT)",
						MarkdownDescription: "Encryption algorithm (`AES256`, `PLAINTEXT`)",
						Validators: []validator.String{
							stringvalidator.OneOf("AES256", "PLAINTEXT"),
						},
					},
				},
			},
			"versioning": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Versioning configuration",
				MarkdownDescription: "Versioning configuration",
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Versioning status (DISABLED, ENABLED, SUSPENDED)",
						MarkdownDescription: "Versioning status (`DISABLED`, `ENABLED`, `SUSPENDED`)",
						Validators: []validator.String{
							stringvalidator.OneOf("DISABLED", "ENABLED", "SUSPENDED"),
						},
					},
				},
			},
			"object_lock": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Object lock (WORM) configuration; requires versioning.status to be ENABLED. Can only be set at bucket creation: changing it recreates the bucket.",
				MarkdownDescription: "Object lock (WORM) configuration; requires `versioning.status` to be `ENABLED`. Can only be set at bucket creation: changing it recreates the bucket.",
				// S3 arms object lock only at CreateBucket.
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Object{
					s3BucketObjectLockRequiresVersioningEnabled{},
				},
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Object lock retention mode (COMPLIANCE, GOVERNANCE)",
						MarkdownDescription: "Object lock retention mode (`COMPLIANCE`, `GOVERNANCE`)",
						Validators: []validator.String{
							stringvalidator.OneOf("COMPLIANCE", "GOVERNANCE"),
						},
					},
					"retention_days": schema.Int64Attribute{
						Required:            true,
						Description:         "Number of days to retain objects",
						MarkdownDescription: "Number of days to retain objects",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Bucket identifier",
				MarkdownDescription: "Bucket identifier",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(s3BucketMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the bucket",
				MarkdownDescription: "Creation date of the bucket",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the bucket",
				MarkdownDescription: "Last update date of the bucket",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(s3BucketMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Bucket readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN, UPDATING)",
				MarkdownDescription: "Bucket readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `SUSPENDED`, `UNKNOWN`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current observed state of the bucket",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(s3BucketMutableAttrs),
				},
				Attributes: s3BucketCurrentStateResourceAttributes(),
			},
		},
	}
}

// s3BucketCurrentStateResourceAttributes returns the current_state attributes of
// the resource schema. The data sources declare the same shape with their own
// schema package.
func s3BucketCurrentStateResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Bucket name",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Geographic region where the bucket is located",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region identifier",
				},
			},
		},
		"encryption": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current encryption configuration",
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
			Description: "Current versioning configuration",
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
			Description: "Current object lock configuration",
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Object lock retention mode",
				},
				"retention_days": schema.Int64Attribute{
					Computed:    true,
					Description: "Number of days to retain objects",
				},
				"retention_years": schema.Int64Attribute{
					Computed:    true,
					Description: "Number of years to retain objects",
				},
			},
		},
		"tags": schema.MapAttribute{
			ElementType: ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Current metadata tags",
		},
		"virtual_host": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Bucket virtual host, as a hostname without scheme. Only returned on a single bucket read.",
		},
		"objects_count": schema.Int64Attribute{
			Computed:    true,
			Description: "Bucket total objects count. Only returned on a single bucket read.",
		},
		"objects_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Bucket total objects size in bytes. Only returned on a single bucket read.",
		},
	}
}

func (r *cloudS3BucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudS3BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudS3BucketModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket"

	var responseData CloudS3BucketAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the bucket id is tracked even if the workflow fails
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if _, err := r.waitForS3BucketReady(ctx, data.ServiceName.ValueString(), responseData.Id); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for bucket to be ready",
			err.Error(),
		)
		return
	}

	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(responseData.Id)
	// json.Unmarshal merges into a non-zero struct: stale map keys/pointers would leak.
	responseData = CloudS3BucketAPIResponse{}
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

func (r *cloudS3BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudS3BucketModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudS3BucketAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudS3BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudS3BucketModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudS3BucketAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := r.waitForS3BucketReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for bucket to be ready after update",
			err.Error(),
		)
		return
	}

	responseData = CloudS3BucketAPIResponse{}
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

func (r *cloudS3BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudS3BucketModel

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

	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (any, string, error) {
			res := &CloudS3BucketAPIResponse{}
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("bucket", data.Id.ValueString(), res.CurrentTasks)
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

func (r *cloudS3BucketResource) waitForS3BucketReady(ctx context.Context, serviceName, bucketId string) (any, error) {
	stateConf := &retry.StateChangeConf{
		// UNKNOWN: transient S3 5xx on HeadBucket.
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC", "UNKNOWN"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudS3BucketAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/object/bucket/" + url.PathEscape(bucketId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			// ERROR is terminal: surface it with the task reason instead of a generic unexpected-state.
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("bucket", bucketId, res.CurrentTasks)
			}
			if res.ResourceStatus == "SUSPENDED" {
				return res, res.ResourceStatus, fmt.Errorf("bucket %s is SUSPENDED: its region is in maintenance, retry once the maintenance is over", bucketId)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
