package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithConfigure = (*cloudProjectGatewayInterfaceResource)(nil)
var _ resource.ResourceWithImportState = (*cloudProjectGatewayInterfaceResource)(nil)

func NewCloudProjectGatewayInterfaceResource() resource.Resource {
	return &cloudProjectGatewayInterfaceResource{}
}

type cloudProjectGatewayInterfaceResource struct {
	config *Config
}

func (r *cloudProjectGatewayInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_gateway_interface"
}

func (d *cloudProjectGatewayInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *cloudProjectGatewayInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectGatewayInterfaceResourceSchema(ctx)
}

func (r *cloudProjectGatewayInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) < 4 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <serviceName>/<region>/<gatewayId>/<interfaceId>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("interface_id"), splits[3])...)
}

func (r *cloudProjectGatewayInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data         CloudProjectGatewayInterfaceModel
		responseData CloudProjectGatewayInterfaceReadModel
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.Region.ValueString()) +
		"/gateway/" + url.PathEscape(data.Id.ValueString()) + "/interface"

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

func (r *cloudProjectGatewayInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data         CloudProjectGatewayInterfaceModel
		responseData CloudProjectGatewayInterfaceReadModel
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.Region.ValueString()) +
		"/gateway/" + url.PathEscape(data.Id.ValueString()) +
		"/interface/" + url.PathEscape(data.InterfaceId.ValueString())

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

func (r *cloudProjectGatewayInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("error unreachable code", "update function should never be reached, resource has no updatable field")
}

func (r *cloudProjectGatewayInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectGatewayInterfaceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.Region.ValueString()) +
		"/gateway/" + url.PathEscape(data.Id.ValueString()) +
		"/interface/" + url.PathEscape(data.InterfaceId.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
