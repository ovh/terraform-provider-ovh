package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudProjectRegionNetworkResource)(nil)
var _ resource.ResourceWithImportState = (*cloudProjectRegionNetworkResource)(nil)

func NewCloudProjectRegionNetworkResource() resource.Resource {
	return &cloudProjectRegionNetworkResource{}
}

type cloudProjectRegionNetworkResource struct {
	config *Config
}

func (r *cloudProjectRegionNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_region_network"
}

func (d *cloudProjectRegionNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *cloudProjectRegionNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudProjectRegionNetworkResourceSchema(ctx)
}

func (r *cloudProjectRegionNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<region_name>/<network_id>")
		return
	}

	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "is_imported", []byte("true"))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudProjectRegionNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudProjectRegionNetworkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle service_name: use provided value or fall back to environment variable
	if data.ServiceName.IsNull() || data.ServiceName.IsUnknown() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	}

	// Create network and retrieve an operation
	var operation CloudProjectOperationResponse
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/network"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &operation); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for operation to complete
	networkID, err := waitForCloudProjectOperation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), operation.Id, "", defaultCloudOperationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error waiting for operation %s", operation.Id), err.Error())
		return
	}

	// Fetch created network
	endpoint = "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/network/" + url.PathEscape(networkID)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudProjectRegionNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudProjectRegionNetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/network/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// If resource was imported, add a placeholder for the `subnet` attribute
	priv, privDiags := req.Private.GetKey(ctx, "is_imported")
	if privDiags.HasError() {
		resp.Diagnostics.Append(privDiags...)
		return
	}
	if priv != nil && string(priv) == "true" {
		r.addImportedResourcePlaceholders(&data)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudProjectRegionNetworkResource) addImportedResourcePlaceholders(data *CloudProjectRegionNetworkModel) {
	// Set the placeholder for the `subnet` property only if it was not already
	// defined in the configuration
	if data.Subnet.IsNull() {
		data.Subnet = SubnetValue{
			state:           attr.ValueStateKnown,
			Cidr:            ovhtypes.NewTfStringValue("placeholder"),
			EnableDhcp:      ovhtypes.NewTfBoolValue(true),
			EnableGatewayIp: ovhtypes.NewTfBoolValue(true),
			IpVersion:       ovhtypes.NewTfInt64Value(4),
		}
	}
}

func (r *cloudProjectRegionNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("update not implemented", "no field is updatable, this function should never be called")
}

func (r *cloudProjectRegionNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudProjectRegionNetworkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/network/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
