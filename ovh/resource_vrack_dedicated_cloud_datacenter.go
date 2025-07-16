package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"golang.org/x/exp/slices"
)

var (
	_ resource.ResourceWithConfigure   = (*vrackDedicatedCloudDatacenterResource)(nil)
	_ resource.ResourceWithImportState = (*vrackDedicatedCloudDatacenterResource)(nil)
)

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

func (r *vrackDedicatedCloudDatacenterResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError("not implemented", "create func should never be called")
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

func (r *vrackDedicatedCloudDatacenterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var planData VrackDedicatedCloudDatacenterModel
	var task VrackTask

	splits := strings.Split(req.ID, "/")
	if len(splits) < 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <serviceName>/<datacenter>/<targetServiceName>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("datacenter"), splits[1])...)

	serviceName := splits[0]
	datacenter := splits[1]
	targetServiceName := splits[2]

	var allowedVrack []string
	endpoint := "/vrack/" + url.PathEscape(serviceName) + "/dedicatedCloudDatacenter/" + url.PathEscape(datacenter) + "/allowedVrack"
	if err := r.config.OVHClient.Get(endpoint, &allowedVrack); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s, 404 means either resource does not exist or the datacenter is not in the vrack", endpoint),
			err.Error(),
		)
		return
	}

	if !slices.Contains(allowedVrack, targetServiceName) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error target_service_name '%s' is not an allowed vrack", targetServiceName),
			fmt.Sprintf("Allowed vrack %v", allowedVrack),
		)
		return
	}

	planData.TargetServiceName = types.NewTfStringValue(targetServiceName)
	endpoint = "/vrack/" + url.PathEscape(serviceName) + "/dedicatedCloudDatacenter/" + url.PathEscape(datacenter) + "/move"
	if err := r.config.OVHClient.Post(endpoint, planData.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to attach dedicatedCloudDatacenter %v: %s", task.ServiceName, datacenter, err),
			err.Error(),
		)
		return
	}

	// service_name is not updated in Read GET response.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), targetServiceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *vrackDedicatedCloudDatacenterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData VrackDedicatedCloudDatacenterModel
	var allowedVrack []string
	var task VrackTask

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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), planData.ServiceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dedicated_cloud"), data.DedicatedCloud)...)

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(planData.Datacenter.ValueString()) + "/allowedVrack"
	if err := r.config.OVHClient.Get(endpoint, &allowedVrack); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s, 404 means either resource does not exist or the datacenter is not in the vrack", endpoint),
			err.Error(),
		)
		return
	}

	if !slices.Contains(allowedVrack, planData.ServiceName.ValueString()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error service_name '%s' is not an allowed vrack", planData.ServiceName.ValueString()),
			fmt.Sprintf("Allowed vrack %v", allowedVrack),
		)
		return
	}

	// do not store target_service_name in state, should be null.
	toCreatePayload := &VrackDedicatedCloudDatacenterModel{TargetServiceName: planData.ServiceName}

	endpoint = "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(planData.Datacenter.ValueString()) + "/move"
	if err := r.config.OVHClient.Post(endpoint, toCreatePayload, &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for vrack (%s) to attach dedicatedCloudDatacenter %v: %s", task.ServiceName, planData.Datacenter.String(), err),
			err.Error(),
		)
		return
	}

	endpoint = "/vrack/" + url.PathEscape(planData.ServiceName.ValueString()) + "/dedicatedCloudDatacenter/" + url.PathEscape(planData.Datacenter.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// service_name will be saved with planData value.
	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *vrackDedicatedCloudDatacenterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("not implemented", "delete func should never be called")
}
