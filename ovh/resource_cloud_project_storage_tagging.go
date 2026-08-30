package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectStorageTaggingResource)(nil)

func NewCloudProjectStorageTaggingResource() resource.Resource {
	return &cloudProjectStorageTaggingResource{}
}

type cloudProjectStorageTaggingResource struct {
	config *Config
}

type cloudProjectStorageTaggingModel struct {
	ID          ovhtypes.TfStringValue                            `tfsdk:"id"`
	ServiceName ovhtypes.TfStringValue                            `tfsdk:"service_name"`
	RegionName  ovhtypes.TfStringValue                            `tfsdk:"region_name"`
	Name        ovhtypes.TfStringValue                            `tfsdk:"name"`
	Tags        ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] `tfsdk:"tags"`
}

type cloudProjectStorageTaggingPayload struct {
	Tags *ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] `json:"tags"`
}

func (r *cloudProjectStorageTaggingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage_tagging"
}

func (r *cloudProjectStorageTaggingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectStorageTaggingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages tags on an existing S3-compatible storage container without owning the bucket lifecycle. " +
			"Use this resource when the bucket is managed by another resource or team.",
		MarkdownDescription: "Manages tags on an existing S3-compatible storage container without owning the bucket lifecycle. " +
			"Use this resource when the bucket is managed by another resource or team.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Unique identifier for the resource (service_name/region_name/name)",
			},
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				MarkdownDescription: "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region name",
				MarkdownDescription: "Region name",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the storage container",
				MarkdownDescription: "Name of the storage container",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				CustomType:          ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](context.Background()),
				Required:            true,
				Description:         "Tags to apply to the storage container. Set to an empty map to remove all tags.",
				MarkdownDescription: "Tags to apply to the storage container. Set to an empty map to remove all tags.",
			},
		},
	}
}


func (r *cloudProjectStorageTaggingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cloudProjectStorageTaggingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() || data.ServiceName.ValueString() == "" {
		envVal := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envVal == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Please provide it or set OVH_CLOUD_PROJECT_SERVICE.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envVal)
	}

	if err := r.putTags(ctx, data); err != nil {
		resp.Diagnostics.AddError("Error setting tags", err.Error())
		return
	}

	data.ID = ovhtypes.NewTfStringValue(compositeID(data))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectStorageTaggingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cloudProjectStorageTaggingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString())

	var response CloudProjectRegionStorageModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &response); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling GET %s", endpoint), err.Error())
		return
	}

	data.Tags = response.Tags
	data.ID = ovhtypes.NewTfStringValue(compositeID(data))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectStorageTaggingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data cloudProjectStorageTaggingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putTags(ctx, data); err != nil {
		resp.Diagnostics.AddError("Error updating tags", err.Error())
		return
	}

	data.ID = ovhtypes.NewTfStringValue(compositeID(data))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectStorageTaggingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cloudProjectStorageTaggingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Clear all tags by sending an empty map.
	emptyTags, diags := ovhtypes.NewTfMapNestedValue[ovhtypes.TfStringValue](ctx, map[string]attr.Value{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Tags = emptyTags

	if err := r.putTags(ctx, data); err != nil {
		resp.Diagnostics.AddError("Error clearing tags", err.Error())
	}
}

func (r *cloudProjectStorageTaggingResource) putTags(ctx context.Context, data cloudProjectStorageTaggingModel) error {
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString())

	payload := cloudProjectStorageTaggingPayload{Tags: &data.Tags}
	return r.config.OVHClient.PutWithContext(ctx, endpoint, payload, nil)
}

func compositeID(data cloudProjectStorageTaggingModel) string {
	return fmt.Sprintf("%s/%s/%s",
		data.ServiceName.ValueString(),
		data.RegionName.ValueString(),
		data.Name.ValueString())
}
