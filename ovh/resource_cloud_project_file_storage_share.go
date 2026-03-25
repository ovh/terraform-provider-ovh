package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var (
	_ resource.ResourceWithConfigure   = (*cloudProjectFileStorageShareResource)(nil)
	_ resource.ResourceWithImportState = (*cloudProjectFileStorageShareResource)(nil)
)

func NewCloudProjectFileStorageShareResource() resource.Resource {
	return &cloudProjectFileStorageShareResource{}
}

type cloudProjectFileStorageShareResource struct {
	config *Config
}

func (r *cloudProjectFileStorageShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_file_storage_share"
}

func (r *cloudProjectFileStorageShareResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectFileStorageShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectFileStorageShareResourceSchema(ctx)
}

func (r *cloudProjectFileStorageShareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<share_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudProjectFileStorageShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudProjectFileStorageShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) + "/share"

	var operation CloudProjectOperationResponse
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &operation); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for the operation to complete and get the resource ID
	shareID, err := waitForCloudProjectOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operation.Id, "", defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for operation %s", operation.Id),
			err.Error(),
		)
		return
	}

	// Read the created share
	var responseData CloudProjectFileStorageShareModel
	readEndpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/share/" + url.PathEscape(shareID)

	if err := r.config.OVHClient.GetWithContext(ctx, readEndpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", readEndpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectFileStorageShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectFileStorageShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/share/" + url.PathEscape(data.Id.ValueString())

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

	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectFileStorageShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudProjectFileStorageShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/share/" + url.PathEscape(data.Id.ValueString())

	var operation CloudProjectOperationResponse
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, planData.ToUpdate(), &operation); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for the operation to complete
	_, err := waitForCloudProjectOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operation.Id, "", defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for operation %s", operation.Id),
			err.Error(),
		)
		return
	}

	// Wait for the share to be available (size changes go through extending/shrinking states)
	_, err = waitForCloudProjectFileStorageShareReady(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.RegionName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for share to be ready after update",
			err.Error(),
		)
		return
	}

	// Read updated resource
	var responseData CloudProjectFileStorageShareModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectFileStorageShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectFileStorageShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/share/" + url.PathEscape(data.Id.ValueString())

	var operation CloudProjectOperationResponse
	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, &operation); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for the operation to complete
	_, err := waitForCloudProjectOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operation.Id, "", defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for operation %s", operation.Id),
			err.Error(),
		)
	}
}

func waitForCloudProjectFileStorageShareReady(ctx context.Context, client *ovhwrap.Client, serviceName, regionName, shareID string) (any, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"creating", "creating_from_snapshot", "extending", "shrinking"},
		Target:  []string{"available"},
		Refresh: func() (any, string, error) {
			res := &CloudProjectFileStorageShareModel{}
			endpoint := "/cloud/project/" + url.PathEscape(serviceName) +
				"/region/" + url.PathEscape(regionName) +
				"/share/" + url.PathEscape(shareID)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res.Id.ValueString(), res.Status.ValueString(), nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
