package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*emailDomainAccountResource)(nil)
	_ resource.ResourceWithImportState = (*emailDomainAccountResource)(nil)
)

// Changes to an account are processed asynchronously. Until they are done, the
// API answers 409 ("This pop is already being processed") to any other change,
// and a deleted account still counts against the domain quota.
const (
	emailAccountPollInterval = 5 * time.Second
	emailAccountPollTimeout  = 5 * time.Minute
)

// emailAccountRetryOnConflict calls f until the API stops answering 409, or the timeout.
func emailAccountRetryOnConflict(ctx context.Context, f func() error) error {
	deadline := time.Now().Add(emailAccountPollTimeout)
	for {
		err := f()
		var apiErr *ovh.APIError
		if err == nil || !errors.As(err, &apiErr) || apiErr.Code != http.StatusConflict || time.Now().After(deadline) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(emailAccountPollInterval):
		}
	}
}

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

	// The password is optional in the schema only so that an imported account,
	// whose password the API never returns, can be managed without setting it.
	// Creating an account still needs one.
	if data.Password.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing password",
			"A password is required to create an email account.",
		)
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
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// The id is not returned by the API. Updates used to plan it as unknown and
	// store it as null, so it is set again here for the states that lost it.
	data.ID = ovhtypes.NewTfStringValue(data.Domain.ValueString() + "/" + data.AccountName.ValueString())

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
	if err := emailAccountRetryOnConflict(ctx, func() error {
		return r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil)
	}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// If the password has changed, call the changePassword endpoint. A null
	// password in the plan leaves the current one untouched: the API never
	// returns it, so there is nothing to compare it with.
	if !planData.Password.IsNull() && !planData.Password.Equal(data.Password) {
		changePasswordEndpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) + "/account/" + url.PathEscape(data.AccountName.ValueString()) + "/changePassword"
		changePasswordPayload := struct {
			Password string `json:"password"`
		}{
			Password: planData.Password.ValueString(),
		}
		if err := emailAccountRetryOnConflict(ctx, func() error {
			return r.config.OVHClient.Post(changePasswordEndpoint, changePasswordPayload, nil)
		}); err != nil {
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
	if err := emailAccountRetryOnConflict(ctx, func() error {
		return r.config.OVHClient.Delete(endpoint, nil)
	}); err != nil {
		if helpers.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// The deletion is asynchronous: wait for the account to be gone, so that
	// it no longer counts against the quota of the domain.
	deadline := time.Now().Add(emailAccountPollTimeout)
	for {
		var account EmailDomainAccountModel
		err := r.config.OVHClient.Get(endpoint, &account)
		if helpers.IsNotFound(err) {
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
			return
		}
		if time.Now().After(deadline) {
			resp.Diagnostics.AddError(
				"Email account not deleted",
				fmt.Sprintf("%s@%s was still there %s after its deletion was accepted",
					data.AccountName.ValueString(), data.Domain.ValueString(), emailAccountPollTimeout),
			)
			return
		}

		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError("Email account not deleted", ctx.Err().Error())
			return
		case <-time.After(emailAccountPollInterval):
		}
	}
}

func (r *emailDomainAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	domain, accountName, found := strings.Cut(req.ID, "/")
	if !found || domain == "" || accountName == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected an import ID of the form <domain>/<account_name>, got %q", req.ID),
		)
		return
	}

	// The password is not set: the API never returns it.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), accountName)...)
}
