package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*vrackDedicatedCloudResource)(nil)

func NewVrackDedicatedCloudResource() resource.Resource {
	return &vrackDedicatedCloudResource{}
}

type vrackDedicatedCloudResource struct {
	config *Config
}

func (r *vrackDedicatedCloudResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrack_dedicated_cloud"
}

func (d *vrackDedicatedCloudResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackDedicatedCloudResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackDedicatedCloudResourceSchema(ctx)
}

func (r *vrackDedicatedCloudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData VrackDedicatedCloudModel
	var task VrackTask

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloud"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to attach dedicatedCloud %v: %s", task.ServiceName, data.DedicatedCloud.String(), err),
			err.Error(),
		)
		return
	}

	endpoint = "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloud/" + url.PathEscape(data.DedicatedCloud.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)
	data.ID = ovhtypes.NewTfStringValue(fmt.Sprintf("%s/%s", data.ServiceName.ValueString(), data.DedicatedCloud.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackDedicatedCloudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackDedicatedCloudModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloud/" + url.PathEscape(data.DedicatedCloud.ValueString())

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

func (r *vrackDedicatedCloudResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "update func should never be called")
}

func (r *vrackDedicatedCloudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VrackDedicatedCloudModel
	var task VrackTask

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloud/" + url.PathEscape(data.DedicatedCloud.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to detach dedicatedCloud %v: %s", task.ServiceName, data.DedicatedCloud.String(), err),
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
