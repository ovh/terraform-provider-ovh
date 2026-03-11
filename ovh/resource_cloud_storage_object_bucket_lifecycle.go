package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (r *cloudStorageObjectBucketLifecycleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_lifecycle"
}

func (r *cloudStorageObjectBucketLifecycleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudStorageObjectBucketLifecycleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the lifecycle configuration of an S3 object storage bucket using the v2 API. " +
			"Setting an empty rules list removes all lifecycle rules.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource identifier (service_name/bucket_name)",
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of lifecycle rules. An empty list removes all rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Unique identifier for the rule",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Rule status (ENABLED or DISABLED)",
						},
						"filter": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Filter that identifies objects to apply the rule to",
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Object key prefix",
								},
								"tags": schema.MapAttribute{
									ElementType: ovhtypes.TfStringType{},
									Optional:    true,
									Description: "Object tags to match",
								},
								"object_size_greater_than": schema.Int64Attribute{
									Optional:    true,
									Description: "Minimum object size in bytes",
								},
								"object_size_less_than": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum object size in bytes",
								},
							},
						},
						"expiration": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "When to expire (delete) objects",
							Attributes: map[string]schema.Attribute{
								"date": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Optional:    true,
									Description: "ISO 8601 date after which objects expire (e.g. 2025-01-01)",
								},
								"days": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of days after object creation before expiration",
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									Optional:    true,
									Description: "Remove expired object delete markers",
								},
							},
						},
						"transitions": schema.ListNestedAttribute{
							Optional:    true,
							Description: "When to transition objects to another storage class",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"date": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Optional:    true,
										Description: "ISO 8601 date of the transition",
									},
									"days": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of days after object creation before transition",
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Required:    true,
										Description: "Target storage class (DEEP_ARCHIVE, GLACIER_IR, HIGH_PERF, STANDARD, STANDARD_IA)",
									},
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "When to expire noncurrent object versions",
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of days after a version becomes noncurrent before expiration",
								},
								"newer_noncurrent_versions": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of newer noncurrent versions to retain",
								},
							},
						},
						"noncurrent_version_transitions": schema.ListNestedAttribute{
							Optional:    true,
							Description: "When to transition noncurrent object versions to another storage class",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"noncurrent_days": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of days after a version becomes noncurrent before transition",
									},
									"newer_noncurrent_versions": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of newer noncurrent versions to retain",
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Required:    true,
										Description: "Target storage class",
									},
								},
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "When to abort incomplete multipart uploads",
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									Required:    true,
									Description: "Number of days after initiation before aborting",
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
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like: <service_name>/<bucket_name>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *cloudStorageObjectBucketLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketLifecycleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putLifecycle(data); err != nil {
		resp.Diagnostics.AddError("Error setting lifecycle rules", err.Error())
		return
	}

	if err := r.readLifecycle(&data); err != nil {
		resp.Diagnostics.AddError("Error reading lifecycle rules", err.Error())
		return
	}

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(data.ServiceName.ValueString() + "/" + data.BucketName.ValueString())}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketLifecycleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readLifecycle(&data); err != nil {
		resp.Diagnostics.AddError("Error reading lifecycle rules", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state CloudStorageObjectBucketLifecycleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id

	if err := r.putLifecycle(data); err != nil {
		resp.Diagnostics.AddError("Error updating lifecycle rules", err.Error())
		return
	}

	if err := r.readLifecycle(&data); err != nil {
		resp.Diagnostics.AddError("Error reading lifecycle rules after update", err.Error())
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

	// Delete all rules by sending an empty list
	emptyRules, _ := types.ListValue(types.ObjectType{AttrTypes: lifecycleRuleAttrTypes()}, []attr.Value{})
	data.Rules = emptyRules

	if err := r.putLifecycle(data); err != nil {
		resp.Diagnostics.AddError("Error deleting lifecycle rules", err.Error())
	}
}

func (r *cloudStorageObjectBucketLifecycleResource) putLifecycle(data CloudStorageObjectBucketLifecycleModel) error {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) +
		"/lifecycle"

	payload := CloudBucketLifecycleRequest{
		Rules: tfListToAPILifecycleRules(data.Rules),
	}

	var response CloudBucketLifecycleResponse
	if err := r.config.OVHClient.Put(endpoint, payload, &response); err != nil {
		return fmt.Errorf("calling Put %s: %w", endpoint, err)
	}
	return nil
}

func (r *cloudStorageObjectBucketLifecycleResource) readLifecycle(data *CloudStorageObjectBucketLifecycleModel) error {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketName.ValueString()) +
		"/lifecycle"

	var response CloudBucketLifecycleResponse
	if err := r.config.OVHClient.Get(endpoint, &response); err != nil {
		return fmt.Errorf("calling Get %s: %w", endpoint, err)
	}

	data.Rules = apiLifecycleRulesToTFList(response.Rules, data.Rules)
	return nil
}
