package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ resource.ResourceWithConfigure   = (*cloudProjectFileStorageShareNetworkResource)(nil)
	_ resource.ResourceWithImportState = (*cloudProjectFileStorageShareNetworkResource)(nil)
)

func NewCloudProjectFileStorageShareNetworkResource() resource.Resource {
	return &cloudProjectFileStorageShareNetworkResource{}
}

type cloudProjectFileStorageShareNetworkResource struct {
	config *Config
}

func (r *cloudProjectFileStorageShareNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_file_storage_share_network"
}

func (r *cloudProjectFileStorageShareNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectFileStorageShareNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectFileStorageShareNetworkResourceSchema(ctx)
}

func (r *cloudProjectFileStorageShareNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<share_network_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudProjectFileStorageShareNetworkResource) endpoint(serviceName, regionName string) string {
	return "/cloud/project/" + url.PathEscape(serviceName) +
		"/region/" + url.PathEscape(regionName) + "/sharenetwork"
}

func (r *cloudProjectFileStorageShareNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudProjectFileStorageShareNetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data.ServiceName.ValueString(), data.RegionName.ValueString())

	var responseData CloudProjectFileStorageShareNetworkModel
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectFileStorageShareNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectFileStorageShareNetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data.ServiceName.ValueString(), data.RegionName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Resource not found, removing from state: %s", endpoint))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = data.ServiceName
	responseData.RegionName = data.RegionName
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectFileStorageShareNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectFileStorageShareNetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(data.ServiceName.ValueString(), data.RegionName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}
}

// Update is a no-op: all input attributes are ForceNew. Terraform will not
// invoke Update, but the method must exist to satisfy the resource.Resource
// interface.
func (r *cloudProjectFileStorageShareNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
