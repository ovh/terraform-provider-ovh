package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithConfigure = (*ovhcloudConnectPopDatacenterConfigResource)(nil)

func NewOvhcloudConnectPopDatacenterConfigResource() resource.Resource {
	return &ovhcloudConnectPopDatacenterConfigResource{}
}

type ovhcloudConnectPopDatacenterConfigResource struct {
	config *Config
}

func (r *ovhcloudConnectPopDatacenterConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connect_pop_datacenter_config"
}

func (d *ovhcloudConnectPopDatacenterConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *ovhcloudConnectPopDatacenterConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OvhcloudConnectPopDatacenterConfigResourceSchema(ctx)
}

func (r *ovhcloudConnectPopDatacenterConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData OvhcloudConnectPopDatacenterConfigModel
	task := OccTask{}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	log.Printf("[DEBUG] Waiting for occ config datacenter %d to be done", task.ResourceID)

	if err := waitForOccTask(ctx, r.config.OVHClient, data.ServiceName.ValueString(), task.TaskID); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for %s config datacenter to be created", data.ServiceName.ValueString()),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = fmt.Sprintf("/ovhCloudConnect/%s/config/pop/%s/datacenter/%d", url.PathEscape(data.ServiceName.ValueString()), url.PathEscape(data.ConfigPopId.String()), task.ResourceID)

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s task %v", endpoint, task),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *ovhcloudConnectPopDatacenterConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData OvhcloudConnectPopDatacenterConfigModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/" + url.PathEscape(data.Id.String())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ovhcloudConnectPopDatacenterConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "update func should never be called")
}

func (r *ovhcloudConnectPopDatacenterConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OvhcloudConnectPopDatacenterConfigModel
	task := OccTask{}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/" + url.PathEscape(data.Id.String())
	if err := r.config.OVHClient.Delete(endpoint, &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}

	log.Printf("[DEBUG] Waiting for occ config datacenter %s deletion to be done", data.Id.String())

	if err := waitForOccTask(ctx, r.config.OVHClient, data.ServiceName.ValueString(), task.TaskID); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for %s config datacenter to be deleted", data.ServiceName.ValueString()),
			err.Error(),
		)
		return
	}
}
