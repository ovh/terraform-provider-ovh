package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                = (*cloudStorageObjectBucketLifecycleResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketLifecycleResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketLifecycleResource)(nil)
)

func NewCloudStorageObjectBucketLifecycleResource() resource.Resource {
	return &cloudStorageObjectBucketLifecycleResource{}
}

type cloudStorageObjectBucketLifecycleResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketLifecycleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_lifecycle"
}

func (r *cloudStorageObjectBucketLifecycleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var storageObjectBucketLifecycleDateValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
	"must be a date in the YYYY-MM-DD format",
)

func (r *cloudStorageObjectBucketLifecycleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the lifecycle configuration of an S3 compatible object storage bucket in a public cloud project using the publicCloud API. " +
			"The resource owns the whole configuration: destroying it removes every lifecycle rule of the bucket.",
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
				Description:         "Bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Lifecycle configuration identifier, the bucket identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "Lifecycle rules of the bucket",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Lifecycle rule unique identifier",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Whether the rule is enabled or disabled (ENABLED, DISABLED)",
							Validators: []validator.String{
								stringvalidator.OneOf("ENABLED", "DISABLED"),
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
								"object_size_greater_than": schema.Int64Attribute{
									Optional:    true,
									Description: "Minimum object size in bytes for the rule to apply",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
								"object_size_less_than": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum object size in bytes for the rule to apply",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
							},
						},
						"expiration": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Expiration settings for current version objects",
							Attributes: map[string]schema.Attribute{
								"date": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Date (YYYY-MM-DD) when objects expire",
									Validators: []validator.String{
										storageObjectBucketLifecycleDateValidator,
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("days")),
									},
								},
								"days": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of days after creation when objects expire",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									Optional:    true,
									Description: "Whether to remove expired object delete markers",
								},
							},
						},
						"transitions": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Transition rules for current version objects",
							Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"date": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Optional:    true,
										Description: "Date (YYYY-MM-DD) when objects transition to another storage class",
										Validators: []validator.String{
											storageObjectBucketLifecycleDateValidator,
											stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("days")),
										},
									},
									"days": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of days after creation when objects transition",
										Validators:  []validator.Int64{int64validator.AtLeast(1)},
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Required:    true,
										Description: "Target storage class for the transition (DEEP_ARCHIVE, GLACIER_IR, STANDARD, STANDARD_IA)",
										Validators: []validator.String{
											stringvalidator.OneOf("DEEP_ARCHIVE", "GLACIER_IR", "STANDARD", "STANDARD_IA"),
										},
									},
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Expiration settings for noncurrent version objects",
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									Required:    true,
									Description: "Number of days after an object version becomes noncurrent before it expires",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
								"newer_noncurrent_versions": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of newer noncurrent versions to retain",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
							},
						},
						"noncurrent_version_transitions": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Transition rules for noncurrent version objects",
							Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"noncurrent_days": schema.Int64Attribute{
										Required:    true,
										Description: "Number of days after an object version becomes noncurrent before it transitions",
										Validators:  []validator.Int64{int64validator.AtLeast(1)},
									},
									"newer_noncurrent_versions": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of newer noncurrent versions to retain before transitioning",
										Validators:  []validator.Int64{int64validator.AtLeast(1)},
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Required:    true,
										Description: "Target storage class for the transition (DEEP_ARCHIVE, GLACIER_IR, HIGH_PERF, STANDARD, STANDARD_IA)",
										Validators: []validator.String{
											stringvalidator.OneOf("DEEP_ARCHIVE", "GLACIER_IR", "HIGH_PERF", "STANDARD", "STANDARD_IA"),
										},
									},
								},
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Settings for aborting incomplete multipart uploads",
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									Required:    true,
									Description: "Number of days after initiation to abort incomplete multipart uploads",
									Validators:  []validator.Int64{int64validator.AtLeast(1)},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *cloudStorageObjectBucketLifecycleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func storageObjectBucketLifecycleEndpoint(serviceName, bucketId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/object/bucket/" + url.PathEscape(bucketId) + "/lifecycle"
}

// apply PUTs the rules then reads them back: the PUT response echoes the request, not what S3 stored.
// The PUT is synchronous despite its 202, so there is nothing to poll.
func (r *cloudStorageObjectBucketLifecycleResource) apply(ctx context.Context, data *CloudStorageObjectBucketLifecycleModel) diag.Diagnostics {
	payload, diags := data.ToPut(ctx)
	if diags.HasError() {
		return diags
	}

	endpoint := storageObjectBucketLifecycleEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, payload, nil); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return diags
	}

	var responseData CloudStorageObjectBucketLifecycleAPIRules
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return diags
	}

	diags.Append(data.MergeWith(ctx, &responseData)...)
	return diags
}

func (r *cloudStorageObjectBucketLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketLifecycleModel

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

func (r *cloudStorageObjectBucketLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketLifecycleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketLifecycleEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	var responseData CloudStorageObjectBucketLifecycleAPIRules
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
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

	// The API answers a bucket without lifecycle configuration with an empty rule list.
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

func (r *cloudStorageObjectBucketLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudStorageObjectBucketLifecycleModel

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

func (r *cloudStorageObjectBucketLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketLifecycleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketLifecycleEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	// An empty rule list deletes the whole lifecycle configuration.
	payload := CloudStorageObjectBucketLifecycleAPIRules{Rules: []CloudStorageObjectBucketLifecycleAPIRule{}}
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
