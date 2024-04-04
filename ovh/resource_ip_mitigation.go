package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
)

var _ resource.ResourceWithConfigure = (*ipMitigationResource)(nil)

func NewIpMitigationResource() resource.Resource {
	return &ipMitigationResource{}
}

type ipMitigationResource struct {
	config *Config
}

func (r *ipMitigationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_mitigation"
}

func (d *ipMitigationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *ipMitigationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IpMitigationResourceSchema(ctx)
}

func (r *ipMitigationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData IpMitigationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ip/" + url.PathEscape(data.Ip.ValueString()) + "/mitigation"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for state to be ok
	endpoint = "/ip/" + url.PathEscape(data.Ip.ValueString()) + "/mitigation/" + url.PathEscape(data.IpOnMitigation.ValueString())
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		readErr := r.config.OVHClient.Get(endpoint, &responseData)
		if readErr != nil {
			return retry.NonRetryableError(readErr)
		}

		if responseData.State.ValueString() == "ok" {
			return nil
		}

		return retry.RetryableError(errors.New("waiting for resource state to be ok"))
	})

	if err != nil {
		resp.Diagnostics.AddError("error waiting status to be ok", err.Error())
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *ipMitigationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData IpMitigationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ip/" + url.PathEscape(data.Ip.ValueString()) + "/mitigation/" + url.PathEscape(data.IpOnMitigation.ValueString())
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

func (r *ipMitigationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData IpMitigationModel

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
	endpoint := "/ip/" + url.PathEscape(data.Ip.ValueString()) + "/mitigation/" + url.PathEscape(data.IpOnMitigation.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		readErr := r.config.OVHClient.Get(endpoint, &responseData)
		if readErr != nil {
			return retry.NonRetryableError(readErr)
		}

		if responseData.State.ValueString() == "ok" {
			return nil
		}

		return retry.RetryableError(errors.New("waiting for resource state to be ok"))
	})

	if err != nil {
		resp.Diagnostics.AddError("error waiting status to be ok", err.Error())
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *ipMitigationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IpMitigationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/ip/" + url.PathEscape(data.Ip.ValueString()) + "/mitigation/" + url.PathEscape(data.IpOnMitigation.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}

	// Wait for resource to be removed
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		readErr := r.config.OVHClient.Get(endpoint, nil)
		if readErr != nil {
			if errOvh, ok := readErr.(*ovh.APIError); ok && errOvh.Code == 404 {
				return nil
			} else {
				return retry.NonRetryableError(readErr)
			}
		}

		return retry.RetryableError(errors.New("waiting for resource to be removed"))
	})

	if err != nil {
		resp.Diagnostics.AddError("error verifying that resource was deleted", err.Error())
	}
}
