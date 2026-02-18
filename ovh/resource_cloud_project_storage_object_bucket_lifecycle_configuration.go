package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectStorageLifecycleConfigurationResource)(nil)
var _ resource.ResourceWithImportState = (*cloudProjectStorageLifecycleConfigurationResource)(nil)

func NewCloudProjectStorageLifecycleConfigurationResource() resource.Resource {
	return &cloudProjectStorageLifecycleConfigurationResource{}
}

type cloudProjectStorageLifecycleConfigurationResource struct {
	config *Config
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage_object_bucket_lifecycle_configuration"
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectStorageLifecycleConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectStorageLifecycleConfigurationResourceSchema(ctx)
}

func (r *cloudProjectStorageLifecycleConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError(
			"Given ID is malformed",
			"ID must be formatted like: <service_name>/<region_name>/<container_name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("container_name"), splits[2])...)
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudProjectStorageLifecycleConfigurationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() {
		data.ServiceName.StringValue = basetypes.NewStringValue(os.Getenv("OVH_CLOUD_PROJECT_SERVICE"))
	}

	endpoint := r.endpoint(data)

	var responseData CloudProjectStorageLifecycleConfigurationModel
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling PUT %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.ContainerName = data.ContainerName
	responseData.ID = ovhtypes.NewTfStringValue(r.compositeID(data))

	reorderRulesByID(&responseData, &data)
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudProjectStorageLifecycleConfigurationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data)

	var responseData CloudProjectStorageLifecycleConfigurationModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Lifecycle configuration not found, removing from state: %s", endpoint))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.ContainerName = data.ContainerName
	responseData.ID = ovhtypes.NewTfStringValue(r.compositeID(data))

	reorderRulesByID(&responseData, &data)
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudProjectStorageLifecycleConfigurationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data)

	var responseData CloudProjectStorageLifecycleConfigurationModel
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, planData.ToUpdate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling PUT %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.ContainerName = data.ContainerName
	responseData.ID = ovhtypes.NewTfStringValue(r.compositeID(data))

	reorderRulesByID(&responseData, &planData)
	responseData.MergeWith(&planData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageLifecycleConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectStorageLifecycleConfigurationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data)

	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling DELETE %s", endpoint),
			err.Error(),
		)
	}
}

func (r *cloudProjectStorageLifecycleConfigurationResource) endpoint(data CloudProjectStorageLifecycleConfigurationModel) string {
	return "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.ContainerName.ValueString()) +
		"/lifecycle"
}

func (r *cloudProjectStorageLifecycleConfigurationResource) compositeID(data CloudProjectStorageLifecycleConfigurationModel) string {
	return fmt.Sprintf("%s/%s/%s",
		data.ServiceName.ValueString(),
		data.RegionName.ValueString(),
		data.ContainerName.ValueString(),
	)
}

// reorderRulesByID reorders the rules in response to match the position order of rules in target by id.
// This is needed because the API may return rules in alphabetical order while the config has them in
// a different order. Without reordering, MergeWith (which operates by index) would pair wrong rules.
func reorderRulesByID(response, target *CloudProjectStorageLifecycleConfigurationModel) {
	if response.Rules.IsNull() || response.Rules.IsUnknown() ||
		target.Rules.IsNull() || target.Rules.IsUnknown() {
		return
	}
	responseElems := response.Rules.Elements()
	targetElems := target.Rules.Elements()
	if len(responseElems) == 0 || len(targetElems) == 0 || len(responseElems) != len(targetElems) {
		return
	}

	// Index response rules by id.
	responseByID := make(map[string]attr.Value, len(responseElems))
	for _, e := range responseElems {
		rule := e.(LifecycleRuleValue)
		if !rule.Id.IsNull() && !rule.Id.IsUnknown() {
			responseByID[rule.Id.ValueString()] = e
		}
	}

	// Build slice in the same order as target.
	reordered := make([]attr.Value, 0, len(responseElems))
	for _, e := range targetElems {
		rule := e.(LifecycleRuleValue)
		if !rule.Id.IsNull() && !rule.Id.IsUnknown() {
			if v, ok := responseByID[rule.Id.ValueString()]; ok {
				reordered = append(reordered, v)
			}
		}
	}
	if len(reordered) != len(responseElems) {
		return // ids don't fully match; leave response order unchanged
	}

	response.Rules = ovhtypes.TfListNestedValue[LifecycleRuleValue]{
		ListValue: basetypes.NewListValueMust(LifecycleRuleValue{}.Type(context.Background()), reordered),
	}
}
