package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketReplicationResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketReplicationResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketReplicationResource)(nil)
)

func NewCloudStorageObjectBucketReplicationResource() resource.Resource {
	return &cloudStorageObjectBucketReplicationResource{}
}

type cloudStorageObjectBucketReplicationResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketReplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication"
}

func (r *cloudStorageObjectBucketReplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var storageObjectBucketReplicationStatusValues = []string{"DISABLED", "ENABLED"}

var storageObjectBucketReplicationStorageClassValues = []string{
	"DEEP_ARCHIVE",
	"GLACIER",
	"GLACIER_IR",
	"HIGH_PERF",
	"INTELLIGENT_TIERING",
	"ONEZONE_IA",
	"STANDARD",
	"STANDARD_IA",
}

func (r *cloudStorageObjectBucketReplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the replication configuration of an S3 compatible object storage bucket in a public cloud project using the publicCloud API. " +
			"The resource owns the whole configuration: destroying it removes every replication rule of the bucket.",
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
			"bucket_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Source bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Source bucket identifier, in the `<REGION>_<name>` form",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Replication configuration identifier, the bucket identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "Replication rules of the bucket",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Replication rule unique identifier",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Whether the rule is enabled or disabled (DISABLED, ENABLED)",
							Validators: []validator.String{
								stringvalidator.OneOf(storageObjectBucketReplicationStatusValues...),
							},
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Priority of the replication rule",
						},
						"delete_marker_replication": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Whether delete markers are replicated (DISABLED, ENABLED). Must be DISABLED when filter.tags is set",
							Validators: []validator.String{
								stringvalidator.OneOf(storageObjectBucketReplicationStatusValues...),
							},
						},
						"destination": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Destination bucket and settings for replicated objects",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Required:    true,
									Description: "Name of the destination bucket",
								},
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Required:    true,
									Description: "Region of the destination bucket, matched case-insensitively",
								},
								"storage_class": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Computed:    true,
									Description: "Storage class for replicated objects in the destination bucket (" + strings.Join(storageObjectBucketReplicationStorageClassValues, ", ") + ")",
									Validators: []validator.String{
										stringvalidator.OneOf(storageObjectBucketReplicationStorageClassValues...),
									},
								},
							},
						},
						"filter": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Filter criteria for the rule. If omitted, the rule applies to every object of the bucket",
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Key prefix to filter objects",
								},
								"tags": schema.MapAttribute{
									ElementType: ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Key-value tag pairs to filter objects",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *cloudStorageObjectBucketReplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func storageObjectBucketReplicationEndpoint(serviceName, bucketId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/object/bucket/" + url.PathEscape(bucketId) + "/replication"
}

// apply PUTs the rules then reads them back: the PUT response echoes the request, not what S3 stored.
// The PUT is synchronous despite its 202, so there is nothing to poll.
func (r *cloudStorageObjectBucketReplicationResource) apply(ctx context.Context, data *CloudStorageObjectBucketReplicationModel) diag.Diagnostics {
	payload, diags := data.ToPut(ctx)
	if diags.HasError() {
		return diags
	}

	endpoint := storageObjectBucketReplicationEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, payload, nil); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return diags
	}

	var responseData CloudStorageObjectBucketReplicationAPIRules
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return diags
	}

	diags.Append(data.MergeWith(ctx, &responseData)...)
	return diags
}

func (r *cloudStorageObjectBucketReplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketReplicationEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	var responseData CloudStorageObjectBucketReplicationAPIRules
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		if storageObjectBucketReadUnavailable(err, endpoint, &resp.Diagnostics) {
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// The API answers a bucket without replication configuration with an empty rule list.
	if len(responseData.Rules) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.MergeWith(ctx, &responseData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketReplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketReplicationEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	// An empty rule list deletes the whole replication configuration.
	payload := CloudStorageObjectBucketReplicationAPIRules{Rules: []CloudStorageObjectBucketReplicationAPIRule{}}
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, payload, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
	}
}
