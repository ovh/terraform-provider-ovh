package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithConfigure = (*dbaasLogsTokenResource)(nil)

func NewDbaasLogsTokenResource() resource.Resource {
	return &dbaasLogsTokenResource{}
}

type dbaasLogsTokenResource struct {
	config *Config
}

func (r *dbaasLogsTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_logs_token"
}

func (d *dbaasLogsTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *dbaasLogsTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DbaasLogsTokenResourceSchema(ctx)
}

func (r *dbaasLogsTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'serviceName/tokenId', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("token_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *dbaasLogsTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data, responseData DbaasLogsTokenModel
		operationData      DbaasLogsTokenReadModel
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create resource
	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/token"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &operationData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for operation to be done
	op, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operationData.OperationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error waiting for operation to be done", err.Error())
		return
	}

	if op.TokenID == nil {
		resp.Diagnostics.AddError("invalid operation state", "a tokenId should be returned in operation but is missing")
		return
	}

	// Read created resource
	endpoint = "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/token/" + url.PathEscape(*op.TokenID)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	responseData.ID = responseData.TokenId

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dbaasLogsTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data         DbaasLogsTokenModel
		responseData DbaasLogsTokenReadModel
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/token/" + url.PathEscape(data.TokenId.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(responseData.toModel())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasLogsTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "update func should never be called")
}

func (r *dbaasLogsTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		data          DbaasLogsTokenModel
		operationData DbaasLogsTokenReadModel
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/token/" + url.PathEscape(data.TokenId.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, &operationData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for deletion to be done
	if _, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operationData.OperationId.ValueString()); err != nil {
		resp.Diagnostics.AddError("error waiting for delete operation to be done", err.Error())
	}
}
