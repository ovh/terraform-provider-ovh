package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.ResourceWithConfigure   = (*iploadbalancingUdpFarmServerResource)(nil)
	_ resource.ResourceWithImportState = (*iploadbalancingUdpFarmServerResource)(nil)
)

func NewIploadbalancingUdpFarmServerResource() resource.Resource {
	return &iploadbalancingUdpFarmServerResource{}
}

type iploadbalancingUdpFarmServerResource struct {
	config *Config
}

func (r *iploadbalancingUdpFarmServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iploadbalancing_udp_farm_server"
}

func (d *iploadbalancingUdpFarmServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *iploadbalancingUdpFarmServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IploadbalancingUdpFarmServerResourceSchema(ctx)
}

func (r *iploadbalancingUdpFarmServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData IploadbalancingUdpFarmServerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ipLoadbalancing/" + url.PathEscape(data.ServiceName.ValueString()) + "/udp/farm/" + strconv.FormatInt(data.FarmId.ValueInt64(), 10) + "/server"
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

func (r *iploadbalancingUdpFarmServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData IploadbalancingUdpFarmServerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/ipLoadbalancing/" + url.PathEscape(data.ServiceName.ValueString()) + "/udp/farm/" + strconv.FormatInt(data.FarmId.ValueInt64(), 10) + "/server/" + strconv.FormatInt(data.ServerId.ValueInt64(), 10)

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

func (r *iploadbalancingUdpFarmServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData, responseData IploadbalancingUdpFarmServerModel

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
	endpoint := "/ipLoadbalancing/" + url.PathEscape(data.ServiceName.ValueString()) + "/udp/farm/" + strconv.FormatInt(data.FarmId.ValueInt64(), 10) + "/server/" + strconv.FormatInt(data.ServerId.ValueInt64(), 10)
	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Read updated resource
	endpoint = "/ipLoadbalancing/" + url.PathEscape(data.ServiceName.ValueString()) + "/udp/farm/" + strconv.FormatInt(data.FarmId.ValueInt64(), 10) + "/server/" + strconv.FormatInt(data.ServerId.ValueInt64(), 10)
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

func (r *iploadbalancingUdpFarmServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IploadbalancingUdpFarmServerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/ipLoadbalancing/" + url.PathEscape(data.ServiceName.ValueString()) + "/udp/farm/" + strconv.FormatInt(data.FarmId.ValueInt64(), 10) + "/server/" + strconv.FormatInt(data.ServerId.ValueInt64(), 10)
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}

func (r *iploadbalancingUdpFarmServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<farmId>/<serverId>")
		return
	}

	serviceName := splits[0]
	farmId, err := strconv.Atoi(splits[1])
	if err != nil {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<farmId>/<serverId> where farmId is a number")
		return
	}
	serverId, err := strconv.Atoi(splits[2])
	if err != nil {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<farmId>/<serverId> where serverId is a number")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("farm_id"), farmId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_id"), serverId)...)
}
