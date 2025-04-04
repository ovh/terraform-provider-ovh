package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*vrackIpv6RoutedSubrangeResource)(nil)
var _ resource.ResourceWithImportState = (*vrackIpv6RoutedSubrangeResource)(nil)

func NewVrackIpv6RoutedSubrangeResource() resource.Resource {
	return &vrackIpv6RoutedSubrangeResource{}
}

type vrackIpv6RoutedSubrangeResource struct {
	config *Config
}

func (r *vrackIpv6RoutedSubrangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrack_ipv6_routed_subrange"
}

func (d *vrackIpv6RoutedSubrangeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackIpv6RoutedSubrangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackIpv6RoutedSubrangeResourceSchema(ctx)
}

func (r *vrackIpv6RoutedSubrangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var responseData VrackIpv6RoutedSubrangeModel

	splits := strings.Split(req.ID, ",")
	if len(splits) != 3 {
		resp.Diagnostics.AddError(
			"Given ID is malformed",
			"import ID must be SERVICE_NAME,IPv6-block,RoutedSubrange formatted",
		)
		return
	}

	serviceName := splits[0]
	block := splits[1]
	routedSubrange := splits[2]

	endpoint := "/vrack/" + url.PathEscape(serviceName) + "/ipv6/" + url.PathEscape(block) + "/routedSubrange/" + url.PathEscape(routedSubrange)

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.ServiceName = ovhtypes.NewTfStringValue(serviceName)
	responseData.Block = ovhtypes.NewTfStringValue(block)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
	// with an id
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"),
		types.StringValue(fmt.Sprintf("vrack_%s-block_%s-routed_subrange_%s", serviceName, block, routedSubrange)))...)
}

func (r *vrackIpv6RoutedSubrangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData VrackIpv6RoutedSubrangeModel
	var task VrackTask

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/ipv6/" + url.PathEscape(data.Block.ValueString()) + "/routedSubrange"
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for ipv6 (%s) to route subrange (%s) from vrack (%s): %s", data.Block.ValueString(), data.RoutedSubrange.ValueString(), task.ServiceName, err),
			err.Error(),
		)
		return
	}

	endpoint = "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/ipv6/" + url.PathEscape(data.Block.ValueString()) + "/routedSubrange/" + url.PathEscape(data.RoutedSubrange.ValueString())

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	// with an id
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"),
		types.StringValue(fmt.Sprintf("vrack_%s-block_%s-routed_subrange_%s", task.ServiceName, data.Block.ValueString(), data.RoutedSubrange.ValueString())))...)
}

func (r *vrackIpv6RoutedSubrangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackIpv6RoutedSubrangeModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/ipv6/" + url.PathEscape(data.Block.ValueString()) + "/routedSubrange/" + url.PathEscape(data.RoutedSubrange.ValueString())

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

func (r *vrackIpv6RoutedSubrangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "update func should never be called")
}

func (r *vrackIpv6RoutedSubrangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VrackIpv6RoutedSubrangeModel
	var task VrackTask

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/vrack/" + url.PathEscape(data.ServiceName.ValueString()) + "/ipv6/" + url.PathEscape(data.Block.ValueString()) + "/routedSubrange/" + url.PathEscape(data.RoutedSubrange.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}

	if err := waitForVrackTask(&task, r.config.OVHClient); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("error waiting for ipv6 (%s) to unroute subrange (%s) from vrack (%s): %s", data.Block.ValueString(), data.RoutedSubrange.ValueString(), task.ServiceName, err),
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
