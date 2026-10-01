package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

const (
	hostingPrivateDatabaseWebhostingNetworkEnabled   = "enabled"
	hostingPrivateDatabaseWebhostingNetworkEnabling  = "enabling"
	hostingPrivateDatabaseWebhostingNetworkDisabled  = "disabled"
	hostingPrivateDatabaseWebhostingNetworkDisabling = "disabling"

	hostingPrivateDatabaseWebhostingNetworkTimeout = 2 * time.Minute
)

var (
	_ resource.Resource                = (*hostingPrivateDatabaseWebhostingNetworkResource)(nil)
	_ resource.ResourceWithConfigure   = (*hostingPrivateDatabaseWebhostingNetworkResource)(nil)
	_ resource.ResourceWithImportState = (*hostingPrivateDatabaseWebhostingNetworkResource)(nil)
)

func NewHostingPrivateDatabaseWebhostingNetworkResource() resource.Resource {
	return &hostingPrivateDatabaseWebhostingNetworkResource{}
}

type hostingPrivateDatabaseWebhostingNetworkResource struct {
	config *Config
}

// HostingPrivateDatabaseWebhostingNetwork is the API representation of
// /hosting/privateDatabase/{serviceName}/webhostingNetwork. The GET returns the status of the
// access, while the POST and the DELETE return the task in charge of the change.
type HostingPrivateDatabaseWebhostingNetwork struct {
	Status string `json:"status,omitempty"`
	TaskId int    `json:"id,omitempty"`
}

type HostingPrivateDatabaseWebhostingNetworkModel struct {
	ID          ovhtypes.TfStringValue `tfsdk:"id"`
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Enabled     ovhtypes.TfBoolValue   `tfsdk:"enabled"`
	Status      ovhtypes.TfStringValue `tfsdk:"status"`
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosting_privatedatabase_webhosting_network"
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allow or deny the OVHcloud webhosting network to reach a private database.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "The internal name of your private database",
				MarkdownDescription: "The internal name of your private database",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				CustomType:          ovhtypes.TfBoolType{},
				Required:            true,
				Description:         "Allow the OVHcloud webhosting network to connect to the private database. A new private database has it enabled",
				MarkdownDescription: "Allow the OVHcloud webhosting network to connect to the private database. A new private database has it enabled",
			},
			"status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Webhosting network status",
				MarkdownDescription: "Webhosting network status",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Unique identifier for the resource",
				MarkdownDescription: "Unique identifier for the resource",
			},
		},
	}
}

func hostingPrivateDatabaseWebhostingNetworkEndpoint(serviceName string) string {
	return "/hosting/privateDatabase/" + url.PathEscape(serviceName) + "/webhostingNetwork"
}

// hostingPrivateDatabaseWebhostingNetworkIsEnabled maps an API status to the enabled flag, a
// transient status being reported as the state it is heading to.
func hostingPrivateDatabaseWebhostingNetworkIsEnabled(status string) bool {
	return status == hostingPrivateDatabaseWebhostingNetworkEnabled ||
		status == hostingPrivateDatabaseWebhostingNetworkEnabling
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) getStatus(serviceName string) (string, error) {
	var res HostingPrivateDatabaseWebhostingNetwork

	endpoint := hostingPrivateDatabaseWebhostingNetworkEndpoint(serviceName)
	if err := r.config.OVHClient.Get(endpoint, &res); err != nil {
		return "", fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	return res.Status, nil
}

// waitForStableStatus waits for the API to leave the enabling and disabling transient statuses.
func (r *hostingPrivateDatabaseWebhostingNetworkResource) waitForStableStatus(ctx context.Context, serviceName string) (string, error) {
	var status string

	err := retry.RetryContext(ctx, hostingPrivateDatabaseWebhostingNetworkTimeout, func() *retry.RetryError {
		var err error

		status, err = r.getStatus(serviceName)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		switch status {
		case hostingPrivateDatabaseWebhostingNetworkEnabled, hostingPrivateDatabaseWebhostingNetworkDisabled:
			return nil
		default:
			return retry.RetryableError(fmt.Errorf("webhosting network of %s is %s", serviceName, status))
		}
	})

	return status, err
}

// apply brings the webhosting network access to the wanted state and returns the reached status.
// Nothing is sent to the API when the wanted state is already the current one, so that a resource
// created on a brand new private database, where the access is enabled by default, is a no-op.
func (r *hostingPrivateDatabaseWebhostingNetworkResource) apply(ctx context.Context, serviceName string, enabled bool) (string, error) {
	status, err := r.waitForStableStatus(ctx, serviceName)
	if err != nil {
		return "", err
	}

	wanted := hostingPrivateDatabaseWebhostingNetworkDisabled
	if enabled {
		wanted = hostingPrivateDatabaseWebhostingNetworkEnabled
	}

	if status == wanted {
		return status, nil
	}

	var task HostingPrivateDatabaseWebhostingNetwork
	endpoint := hostingPrivateDatabaseWebhostingNetworkEndpoint(serviceName)

	if enabled {
		err = r.config.OVHClient.Post(endpoint, nil, &task)
	} else {
		err = r.config.OVHClient.Delete(endpoint, &task)
	}
	if err != nil {
		return "", fmt.Errorf("error calling %s: %w", endpoint, err)
	}

	taskEndpoint := fmt.Sprintf("/hosting/privateDatabase/%s/tasks/%d", url.PathEscape(serviceName), task.TaskId)
	if err := WaitArchivedHostingPrivateDabaseTask(r.config.OVHClient, taskEndpoint, hostingPrivateDatabaseWebhostingNetworkTimeout); err != nil {
		return "", fmt.Errorf("error waiting for task %s: %w", taskEndpoint, err)
	}

	return r.waitForStableStatus(ctx, serviceName)
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostingPrivateDatabaseWebhostingNetworkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	status, err := r.apply(ctx, serviceName, data.Enabled.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error setting the webhosting network access of %s", serviceName),
			err.Error(),
		)
		return
	}

	data.ID = data.ServiceName
	data.Status = ovhtypes.NewTfStringValue(status)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostingPrivateDatabaseWebhostingNetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	status, err := r.getStatus(serviceName)
	if err != nil {
		// getStatus and apply wrap the client error, so unwrap it to reach the API error.
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the webhosting network access of %s", serviceName),
			err.Error(),
		)
		return
	}

	data.ID = data.ServiceName
	data.Status = ovhtypes.NewTfStringValue(status)
	data.Enabled = ovhtypes.NewTfBoolValue(hostingPrivateDatabaseWebhostingNetworkIsEnabled(status))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HostingPrivateDatabaseWebhostingNetworkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	status, err := r.apply(ctx, serviceName, data.Enabled.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error setting the webhosting network access of %s", serviceName),
			err.Error(),
		)
		return
	}

	data.ID = data.ServiceName
	data.Status = ovhtypes.NewTfStringValue(status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostingPrivateDatabaseWebhostingNetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	// The webhosting network access is enabled by default on a private database, restore that
	// default instead of leaving behind an access closed by a resource that no longer exists.
	if _, err := r.apply(ctx, serviceName, true); err != nil {
		// getStatus and apply wrap the client error, so unwrap it to reach the API error.
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error restoring the webhosting network access of %s", serviceName),
			err.Error(),
		)
	}
}

func (r *hostingPrivateDatabaseWebhostingNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
