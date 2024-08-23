package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ resource.ResourceWithConfigure = (*okmsServiceKeyResource)(nil)

func NewOkmsServiceKeyResource() resource.Resource {
	return &okmsServiceKeyResource{}
}

type okmsServiceKeyResource struct {
	config *Config
}

func (r *okmsServiceKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_service_key"
}

func (d *okmsServiceKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *okmsServiceKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OkmsServiceKeyResourceSchema(ctx)
}

func (r *okmsServiceKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData OkmsServiceKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *okmsServiceKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData OkmsServiceKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/okms/resource/" + url.PathEscape(data.OkmsId.ValueString()) + "/serviceKey/" + url.PathEscape(data.Id.ValueString())

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

func (r *okmsServiceKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData OkmsServiceKeyResourceModel

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

func (r *okmsServiceKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsServiceKeyResourceModel

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
