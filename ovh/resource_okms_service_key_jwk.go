package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ resource.ResourceWithConfigure = (*okmsServiceKeyJwkResource)(nil)

func NewOkmsServiceKeyJwkResource() resource.Resource {
	return &okmsServiceKeyJwkResource{}
}

type okmsServiceKeyJwkResource struct {
	config *Config
}

func (r *okmsServiceKeyJwkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_service_key_jwk"
}

func (d *okmsServiceKeyJwkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *okmsServiceKeyJwkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OkmsServiceKeyJwkResourceSchema(ctx)
}

func (r *okmsServiceKeyJwkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData OkmsServiceKeyJwkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey/"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(&resp.Diagnostics, ctx), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Get the jwk details
	endpoint = endpoint + url.PathEscape(data.Id.ValueString()) + "?format=JWK"
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

func (r *okmsServiceKeyJwkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData OkmsServiceKeyJwkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey/" + url.PathEscape(data.Id.ValueString()) + "?format=JWK"

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

func (r *okmsServiceKeyJwkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData OkmsServiceKeyJwkResourceModel

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
	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
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

func (r *okmsServiceKeyJwkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsServiceKeyJwkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.State = ovhtypes.NewTfStringValue("DEACTIVATED")
	data.DeactivationReason = ovhtypes.NewTfStringValue("UNSPECIFIED")
	// Deactivate key first
	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deactivating key %s", endpoint),
			err.Error(),
		)
		return
	}

	// Delete API call logic
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
