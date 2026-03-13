package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectSshKeyResource)(nil)

func NewCloudProjectSshKeyResource() resource.Resource {
	return &cloudProjectSshKeyResource{}
}

type cloudProjectSshKeyResource struct {
	config *Config
}

func (r *cloudProjectSshKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_ssh_key"
}

func (d *cloudProjectSshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *cloudProjectSshKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectSshKeyResourceSchema(ctx)
}

// sshKeyBaseEndpoint returns the base API v2 endpoint for SSH keys.
func sshKeyBaseEndpoint(serviceName string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/sshKey"
}

// sshKeyEndpoint returns the API v2 endpoint for a specific SSH key.
func sshKeyEndpoint(serviceName, region, name string) string {
	return sshKeyBaseEndpoint(serviceName) + "/" + url.PathEscape(region) + "/" + url.PathEscape(name)
}

// resolveServiceName returns the service_name from the model or environment variable.
func (r *cloudProjectSshKeyResource) resolveServiceName(data *CloudProjectSshKeyModel, resp *resource.CreateResponse) string {
	if !data.ServiceName.IsNull() && !data.ServiceName.IsUnknown() && data.ServiceName.ValueString() != "" {
		return data.ServiceName.ValueString()
	}
	envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
	if envServiceName == "" {
		resp.Diagnostics.AddError(
			"Missing service_name",
			"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
		)
		return ""
	}
	data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	return envServiceName
}

func (r *cloudProjectSshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudProjectSshKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := r.resolveServiceName(&data, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build APIv2 create request
	createReq := sshKeyCreateRequest{
		TargetSpec: sshKeyTargetSpec{
			Name:      data.Name.ValueString(),
			PublicKey: data.PublicKey.ValueString(),
			Location:  sshKeyLocation{Region: data.Region.ValueString()},
		},
	}

	var envelope sshKeyEnvelope
	endpoint := sshKeyBaseEndpoint(serviceName)
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, createReq, &envelope); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling POST %s", endpoint),
			err.Error(),
		)
		return
	}

	// Poll until READY (SSH key creation is fast but async)
	region := data.Region.ValueString()
	name := data.Name.ValueString()
	getEndpoint := sshKeyEndpoint(serviceName, region, name)

	for i := 0; i < 60; i++ {
		if envelope.ResourceStatus == "READY" {
			break
		}
		time.Sleep(2 * time.Second)
		tflog.Debug(ctx, "Waiting for SSH key to be READY", map[string]interface{}{
			"status":  envelope.ResourceStatus,
			"attempt": i + 1,
		})
		if err := r.config.OVHClient.GetWithContext(ctx, getEndpoint, &envelope); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error polling SSH key status at %s", getEndpoint),
				err.Error(),
			)
			return
		}
	}

	if envelope.ResourceStatus != "READY" {
		resp.Diagnostics.AddError(
			"SSH key creation timed out",
			fmt.Sprintf("SSH key did not reach READY status (current: %s)", envelope.ResourceStatus),
		)
		return
	}

	responseData := envelope.toModel()
	responseData.ServiceName = data.ServiceName
	responseData.MergeWith(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, responseData)...)
}

func (r *cloudProjectSshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudProjectSshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := sshKeyEndpoint(
		data.ServiceName.ValueString(),
		data.Region.ValueString(),
		data.Name.ValueString(),
	)

	var envelope sshKeyEnvelope
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &envelope); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData := envelope.toModel()
	responseData.ServiceName = data.ServiceName
	data.MergeWith(responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectSshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "SSH keys are immutable — this func should never be called")
}

func (r *cloudProjectSshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectSshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := sshKeyEndpoint(
		data.ServiceName.ValueString(),
		data.Region.ValueString(),
		data.Name.ValueString(),
	)

	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling DELETE %s", endpoint),
			err.Error(),
		)
		return
	}

	// Poll until the key is gone (404)
	for i := 0; i < 30; i++ {
		time.Sleep(2 * time.Second)
		var envelope sshKeyEnvelope
		err := r.config.OVHClient.GetWithContext(ctx, endpoint, &envelope)
		if err != nil {
			// 404 means deletion is complete
			break
		}
		if envelope.ResourceStatus != "DELETING" {
			break
		}
		tflog.Debug(ctx, "Waiting for SSH key deletion", map[string]interface{}{
			"status":  envelope.ResourceStatus,
			"attempt": i + 1,
		})
	}
}
