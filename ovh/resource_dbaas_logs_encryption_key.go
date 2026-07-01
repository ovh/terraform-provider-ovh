package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithConfigure = (*dbaasLogsEncryptionKeyResource)(nil)
var _ resource.ResourceWithImportState = (*dbaasLogsEncryptionKeyResource)(nil)

func NewDbaasLogsEncryptionKeyResource() resource.Resource {
	return &dbaasLogsEncryptionKeyResource{}
}

type dbaasLogsEncryptionKeyResource struct {
	config *Config
}

func (r *dbaasLogsEncryptionKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_logs_encryption_key"
}

func (r *dbaasLogsEncryptionKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dbaasLogsEncryptionKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DbaasLogsEncryptionKeyResourceSchema(ctx)
}

func (r *dbaasLogsEncryptionKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'service_name/encryption_key_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("encryption_key_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *dbaasLogsEncryptionKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data         DbaasLogsEncryptionKeyModel
		responseData DbaasLogsEncryptionKeyModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/encryptionKey"
	opRes := &DbaasLogsOperation{}
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), opRes); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	op, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), opRes.OperationId)
	if err != nil {
		resp.Diagnostics.AddError("error waiting for operation", err.Error())
		return
	}

	if op.EncryptionKeyId == nil {
		resp.Diagnostics.AddError("invalid operation state", "an encryptionKeyId should be returned in operation but is missing")
		return
	}

	endpoint = "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/encryptionKey/" + url.PathEscape(*op.EncryptionKeyId)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	responseData.ID = responseData.EncryptionKeyId

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dbaasLogsEncryptionKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data         DbaasLogsEncryptionKeyModel
		responseData DbaasLogsEncryptionKeyModel
	)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/encryptionKey/" + url.PathEscape(data.EncryptionKeyId.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasLogsEncryptionKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		data         DbaasLogsEncryptionKeyModel
		responseData DbaasLogsEncryptionKeyModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/encryptionKey/" + url.PathEscape(data.EncryptionKeyId.ValueString())
	opRes := &DbaasLogsOperation{}
	if err := r.config.OVHClient.Put(endpoint, data.ToUpdate(), opRes); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), opRes.OperationId); err != nil {
		resp.Diagnostics.AddError("error waiting for operation", err.Error())
		return
	}

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dbaasLogsEncryptionKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DbaasLogsEncryptionKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/dbaas/logs/" + url.PathEscape(data.ServiceName.ValueString()) + "/encryptionKey/" + url.PathEscape(data.EncryptionKeyId.ValueString())
	opRes := &DbaasLogsOperation{}
	if err := r.config.OVHClient.Delete(endpoint, opRes); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), opRes.OperationId); err != nil {
		resp.Diagnostics.AddError("error waiting for delete operation", err.Error())
	}
}
