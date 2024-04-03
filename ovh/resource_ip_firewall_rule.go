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

var _ resource.ResourceWithConfigure = (*ipFirewallRuleResource)(nil)

func NewIpFirewallRuleResource() resource.Resource {
	return &ipFirewallRuleResource{}
}

type ipFirewallRuleResource struct {
	config *Config
}

func (r *ipFirewallRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_firewall_rule"
}

func (d *ipFirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *ipFirewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IpFirewallRuleResourceSchema(ctx)
}

func (r *ipFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData IpFirewallRuleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
	)
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	endpoint = fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)

	// Wait for state to be ok
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

func (r *ipFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData IpFirewallRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)
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

func (r *ipFirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update on API side
}

func (r *ipFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IpFirewallRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for rule to be removed
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
