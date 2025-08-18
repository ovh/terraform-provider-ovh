package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.ResourceWithConfigure = (*cloudProjectStorageResource)(nil)

func NewCloudProjectStorageResource() resource.Resource {
	return &cloudProjectStorageResource{}
}

type cloudProjectStorageResource struct {
	config *Config
}

func (r *cloudProjectStorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage"
}

func (d *cloudProjectStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.config = config
}

func (d *cloudProjectStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectRegionStorageResourceSchema(ctx)
}

func (r *cloudProjectStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<storage_name>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), splits[2])...)
}

func (r *cloudProjectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectRegionStorageModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectRegionStorageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add filters if defined
	queryParams := url.Values{}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(data.Limit.ValueInt64(), 10))
	}
	if !data.Marker.IsNull() && !data.Marker.IsUnknown() {
		queryParams.Add("marker", data.Marker.ValueString())
	}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		queryParams.Add("prefix", data.Prefix.ValueString())
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString()) + "?" + queryParams.Encode()
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData CloudProjectRegionStorageModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Add filters if defined
	queryParams := url.Values{}
	if !planData.Limit.IsNull() && !planData.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(planData.Limit.ValueInt64(), 10))
	}
	if !planData.Marker.IsNull() && !planData.Marker.IsUnknown() {
		queryParams.Add("marker", planData.Marker.ValueString())
	}
	if !planData.Prefix.IsNull() && !planData.Prefix.IsUnknown() {
		queryParams.Add("prefix", planData.Prefix.ValueString())
	}

	// Read updated resource
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString()) + "?" + queryParams.Encode()
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectRegionStorageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString()) +
		"/object?withVersions=true"
	bulkDeleteEndpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString()) +
		"/bulkDeleteObjects"

	// Try to empty bucket to be able to destroy it completely.
	// This operation can fail if objects are locked. In this case, they must be unlocked manually
	// before retrying the operation.
	// A maximum of 1000 objects are returned at each GET call, so we have to iterate to empty bucket.
	for {
		var (
			objects     []CloudProjectStorageObjectsValue
			idsToDelete []map[string]string
		)

		if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &objects); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		if len(objects) == 0 {
			break
		}

		for _, obj := range objects {
			idsToDelete = append(idsToDelete, map[string]string{
				"key":       obj.Key.ValueString(),
				"versionId": obj.VersionId.ValueString(),
			})
		}

		tflog.Info(ctx, fmt.Sprintf("removing objects %s", idsToDelete))
		if err := r.config.OVHClient.PostWithContext(ctx, bulkDeleteEndpoint, map[string]any{
			"objects": idsToDelete,
		}, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", bulkDeleteEndpoint),
				err.Error(),
			)
			return
		}
	}

	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString())

	// Delete bucket itself
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
