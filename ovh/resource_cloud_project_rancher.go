package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

var _ resource.ResourceWithConfigure = (*cloudProjectRancherResource)(nil)
var _ resource.ResourceWithImportState = (*cloudProjectRancherResource)(nil)

func NewCloudProjectRancherResource() resource.Resource {
	return &cloudProjectRancherResource{}
}

type cloudProjectRancherResource struct {
	config *Config
}

func (r *cloudProjectRancherResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_rancher"
}

func (r *cloudProjectRancherResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudProjectRancherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectRancherResourceSchema(ctx)
}

func (r *cloudProjectRancherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<rancher ID>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudProjectRancherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectRancherModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create Rancher
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// The password is only returned in response of the initial POST, so
	// we save it to set it in the state in the end
	savedBootstrapPassword := responseData.CurrentState.BootstrapPassword

	// Wait for service to be ready
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(responseData.Id.ValueString())
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.RawOVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Fetch up-to-date service info
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	responseData.MergeWith(&data)
	responseData.CurrentState.BootstrapPassword = savedBootstrapPassword

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectRancherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectRancherModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectRancherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData CloudProjectRancherModel

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

	// Version is mandatory when updating resource, so use the current one if not defined in config
	updateData := planData.ToUpdate()
	if updateData.TargetSpec.Version == nil {
		updateData.TargetSpec.Version = &data.CurrentState.Version
	}

	// Update resource
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	// Wait for service to be ready
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.RawOVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	responseData.MergeWith(&planData)
	responseData.MergeWith(&data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectRancherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectRancherModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
	}
}
