package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"golang.org/x/exp/slices"
)

var _ resource.ResourceWithConfigure = (*vrackDedicatedCloudDatacenterResource)(nil)

func NewVrackDedicatedCloudDatacenterResource() resource.Resource {
	return &vrackDedicatedCloudDatacenterResource{}
}

type vrackDedicatedCloudDatacenterResource struct {
	config *Config
}

func (r *vrackDedicatedCloudDatacenterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrack_dedicated_cloud_datacenter"
}

func (d *vrackDedicatedCloudDatacenterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackDedicatedCloudDatacenterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackDedicatedCloudDatacenterResourceSchema(ctx)
}

func (r *vrackDedicatedCloudDatacenterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData VrackDedicatedCloudDatacenterModel
	var allowedVrack []string
	var task VrackTask

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(data.Datacenter.ValueString()) + "/allowedVrack"
	if err := r.config.OVHClient.Get(endpoint, &allowedVrack); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s, 404 means either resource does not exist or the datacenter is not in the vrack", endpoint),
			err.Error(),
		)
		return
	}

	if !slices.Contains(allowedVrack, data.TargetServiceName.ValueString()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error target_service_name '%s' is not an allowed vrack", data.TargetServiceName.ValueString()),
			fmt.Sprintf("Allowed vrack %v", allowedVrack),
		)
		return
	}

	endpoint = "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(data.Datacenter.ValueString()) + "/move"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to attach dedicatedCloudDatacenter %v: %s", task.ServiceName, data.Datacenter.String(), err),
			err.Error(),
		)
		return
	}

	endpoint = "/vrack/" + url.PathEscape(data.TargetServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(data.Datacenter.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// there's no serviceName in the Get response above.
	responseData.ServiceName = responseData.Vrack

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackDedicatedCloudDatacenterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackDedicatedCloudDatacenterModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(data.Datacenter.ValueString())

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

func (r *vrackDedicatedCloudDatacenterResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "update func should never be called")
}

func (r *vrackDedicatedCloudDatacenterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}
