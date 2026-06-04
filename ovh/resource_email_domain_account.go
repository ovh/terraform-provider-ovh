package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*emailDomainAccountResource)(nil)

func NewEmailDomainAccountResource() resource.Resource {
	return &emailDomainAccountResource{}
}

type emailDomainAccountResource struct {
	config *Config
}

func (r *emailDomainAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain_account"
}

func (r *emailDomainAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *emailDomainAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = EmailDomainAccountResourceSchema(ctx)
}

func (r *emailDomainAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData EmailDomainAccountModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the account
	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read the created account
	endpoint = "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	responseData.ID = ovhtypes.NewTfStringValue(
		data.Domain.ValueString() + "/" + data.AccountName.ValueString(),
	)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *emailDomainAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData EmailDomainAccountModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString())
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

func (r *emailDomainAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData EmailDomainAccountModel

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

	// Update description and/or size via PUT
	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// If the password has changed, call the changePassword endpoint
	if !planData.Password.Equal(data.Password) {
		changePasswordEndpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString()) + "/changePassword"
		changePasswordPayload := struct {
			Password string `json:"password"`
		}{
			Password: planData.Password.ValueString(),
		}
		if err := r.config.OVHClient.Post(changePasswordEndpoint, changePasswordPayload, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", changePasswordEndpoint),
				err.Error(),
			)
			return
		}
	}

	// Read updated resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *emailDomainAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EmailDomainAccountModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
