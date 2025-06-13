package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var _ resource.ResourceWithConfigure = (*vpsReinstallResource)(nil)

func NewVpsReinstallResource() resource.Resource {
	return &vpsReinstallResource{}
}

type vpsReinstallResource struct {
	config *Config
}

func (r *vpsReinstallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_reinstall"
}

func (d *vpsReinstallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vpsReinstallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VpsReinstallResourceSchema(ctx)
}

func (r *vpsReinstallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData VpsReinstallModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vps/" + url.PathEscape(data.ServiceName.ValueString()) + "/rebuild"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}
	responseData.MergeWith(&data)

	// Wait for reinstallation to complete
	err := r.waitForVPSReinstall(ctx, data.ServiceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching updated resource",
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (d *vpsReinstallResource) waitForVPSReinstall(ctx context.Context, serviceName string) error {
	var responseData VpsModel

	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		// Read updated resource
		endpoint := "/vps/" + url.PathEscape(serviceName)
		if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error calling GET %s", endpoint))
		}

		// Update is succesfull, return
		if responseData.State.ValueString() == "installing" {
			return retry.RetryableError(errors.New("waiting for vps to finish reinstallation"))
		}

		return nil

	})

	return err
}

func (r *vpsReinstallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VpsReinstallModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsReinstallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData VpsReinstallModel

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
	endpoint := "/vps/" + url.PathEscape(planData.ServiceName.ValueString()) + "/rebuild"
	if err := r.config.OVHClient.Post(endpoint, planData.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}
	responseData.MergeWith(&planData)

	// Wait for reinstallation to complete
	err := r.waitForVPSReinstall(ctx, data.ServiceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching updated resource",
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *vpsReinstallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VpsReinstallModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}
