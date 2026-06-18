package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudSshKeyResource)(nil)

func NewCloudSshKeyResource() resource.Resource {
	return &cloudSshKeyResource{}
}

type cloudSshKeyResource struct {
	config *Config
}

func (r *cloudSshKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_ssh_key"
}

func (r *cloudSshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudSshKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudSshKeyResourceSchema(ctx)
}

func cloudSshKeyBaseEndpoint(serviceName string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/sshKey"
}

func cloudSshKeyEndpoint(serviceName, name string) string {
	return cloudSshKeyBaseEndpoint(serviceName) + "/" + url.PathEscape(name)
}

func (r *cloudSshKeyResource) resolveServiceName(data *CloudSshKeyModel, diags interface{ AddError(string, string) }) string {
	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() && data.ServiceName.ValueString() != "" {
		return data.ServiceName.ValueString()
	}
	envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
	if envServiceName == "" {
		diags.AddError(
			"Missing service_name",
			"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
		)
		return ""
	}
	data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	return envServiceName
}

func (r *cloudSshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudSshKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := r.resolveServiceName(&data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := cloudSshKeyCreateRequest{
		Name:      data.Name.ValueString(),
		PublicKey: data.PublicKey.ValueString(),
	}

	var sshKey cloudSshKeyResponse
	endpoint := cloudSshKeyBaseEndpoint(serviceName)
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, createReq, &sshKey); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling POST %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData := sshKey.toModel()
	responseData.ServiceName = data.ServiceName

	resp.Diagnostics.Append(resp.State.Set(ctx, responseData)...)
}

func (r *cloudSshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudSshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := cloudSshKeyEndpoint(
		data.ServiceName.ValueString(),
		data.Name.ValueString(),
	)

	var sshKey cloudSshKeyResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &sshKey); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData := sshKey.toModel()
	responseData.ServiceName = data.ServiceName

	resp.Diagnostics.Append(resp.State.Set(ctx, responseData)...)
}

func (r *cloudSshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "SSH keys are immutable — this func should never be called")
}

func (r *cloudSshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudSshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := cloudSshKeyEndpoint(
		data.ServiceName.ValueString(),
		data.Name.ValueString(),
	)

	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling DELETE %s", endpoint),
			err.Error(),
		)
	}
}
