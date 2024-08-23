package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var _ resource.ResourceWithConfigure = (*okmsCredentialResource)(nil)
var _ resource.ResourceWithImportState = (*okmsCredentialResource)(nil)

func NewOkmsCredentialResource() resource.Resource {
	return &okmsCredentialResource{}
}

type okmsCredentialResource struct {
	config *Config
}

func (r *okmsCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_credential"
}

func (d *okmsCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *okmsCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OkmsCredentialResourceSchema(ctx)
}

func (r *okmsCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	ids := strings.Split(req.ID, "/")
	if len(ids) != 2 {
		resp.Diagnostics.AddError("Error importing okms_credential resource.", "ID should be of the format '{okmsId}/{credentialId}'")
		return
	}

	// Set ID attributes in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("okms_id"), ids[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ids[1])...)

	var data OkmsCredentialResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + ids[0] + "/credential/" + ids[1]
	if err := r.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error importing credential in GET %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData OkmsCredentialResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	okmsId := url.PathEscape(data.OkmsId.ValueString())
	endpoint := "/v2/okms/resource/" + okmsId + "/credential"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	var updateData OkmsCredentialResourceModel
	// Read updated credential to get the certificate
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		// Read updated resource
		credId := url.PathEscape(responseData.Id.ValueString())
		endpoint := "/v2/okms/resource/" + okmsId + "/credential/" + credId

		if err := r.config.OVHClient.Get(endpoint, &updateData); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error reading credential %s", credId))
		}

		status := updateData.Status.ValueString()
		if status == "READY" {
			// KMS was created successfully, return
			return nil
		} else if status != "CREATING" {
			return retry.NonRetryableError(
				fmt.Errorf("Unexpected credential status : %s",
					status,
				))
		}

		return retry.RetryableError(errors.New("Waiting for KMS credential readiness"))
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create KMS, creation timeout", err.Error())
	}

	responseData.MergeWith(&updateData)
	responseData.Status = updateData.Status

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *okmsCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData OkmsCredentialResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/credential/" + url.PathEscape(data.Id.ValueString()) + ""

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading credential with GET %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData OkmsCredentialResourceModel

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

	// We should only get an update request if the CSR can be replaced without recreation
	data.Csr = planData.Csr
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsCredentialResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/credential/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
