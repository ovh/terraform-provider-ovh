package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// vrackServicesBGPPlanCodePrefix identifies the BGP option plan codes of a vRack Services
// (e.g. bgp-beta-vrack-services, bgp-beta-vrack-services-eu).
const vrackServicesBGPPlanCodePrefix = "bgp-"

var _ resource.ResourceWithConfigure = (*vrackServicesResource)(nil)
var _ resource.ResourceWithImportState = (*vrackServicesResource)(nil)

func NewVrackServicesResource() resource.Resource {
	return &vrackServicesResource{}
}

type vrackServicesResource struct {
	config *Config
}

func (r *vrackServicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vrackservices"
}

func (d *vrackServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *vrackServicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VrackServicesResourceSchema(ctx)
}

func (d *vrackServicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *vrackServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VrackServicesModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create order and wait for service to be delivered
	order := data.ToOrder()
	if err := orderCreate(order, r.config, "vrackServices", true, defaultOrderTimeout); err != nil {
		resp.Diagnostics.AddError("failed to create order", err.Error())
		return
	}

	// Find service name from order
	orderID := order.Order.OrderId.ValueInt64()
	plans := []PlanValue{}
	resp.Diagnostics.Append(data.Plan.ElementsAs(ctx, &plans, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve service name
	serviceName, err := serviceNameFromOrder(r.config.OVHClient, orderID, plans[0].PlanCode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve service name", err.Error())
		return
	}
	data.Id = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(serviceName),
	}

	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())

	// The BGP option ordered alongside the vRack Services may be delivered after it
	if bgpPlanOption(data.PlanOption) != nil {
		if err := waitForVrackServicesBGP(ctx, r.config, endpoint); err != nil {
			resp.Diagnostics.AddError("failed to wait for BGP option delivery", err.Error())
			return
		}
	}

	// Read created resource to get the Checksum for update
	var createdData VrackServicesModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &createdData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}
	data.Checksum = ovhtypes.TfStringValue{
		StringValue: basetypes.NewStringValue(createdData.Checksum.ValueString()),
	}

	// Update resource
	updateData := data.ToUpdate()
	updateData.TargetSpec.Bgp = bgpTargetSpecToUpdate(data.TargetSpec.Bgp, &createdData)
	if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for update
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Fetch up-to-date service info
	var responseData VrackServicesModel
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&responseData, true)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData VrackServicesModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData, true)

	// targetSpec.bgp only changes outside Terraform through the BGP option lifecycle (delivery,
	// termination), so the remote one is always right.
	if !data.TargetSpec.IsNull() && !data.TargetSpec.IsUnknown() && !responseData.TargetSpec.Bgp.IsUnknown() {
		data.TargetSpec.Bgp = responseData.TargetSpec.Bgp
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vrackServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var responseData, data, planData VrackServicesModel

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

	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(data.Id.ValueString())

	// Order or terminate the BGP option before touching the target specification:
	// targetSpec.bgp can only be set while the option is subscribed.
	checksum := data.Checksum
	stateBGPOption := bgpPlanOption(data.PlanOption)
	planBGPOption := bgpPlanOption(planData.PlanOption)
	switch {
	case planBGPOption != nil && stateBGPOption == nil:
		if err := r.orderBGPOption(ctx, data.Id.ValueString(), planData.OvhSubsidiary.ValueString(), *planBGPOption); err != nil {
			resp.Diagnostics.AddError("failed to order BGP option", err.Error())
			return
		}
	case planBGPOption == nil && stateBGPOption != nil:
		if err := r.terminateBGPOption(ctx, data.Id.ValueString()); err != nil {
			resp.Diagnostics.AddError("failed to terminate BGP option", err.Error())
			return
		}
	}

	if !planData.TargetSpec.IsNull() && !planData.TargetSpec.IsUnknown() {
		var currentData VrackServicesModel
		if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &currentData); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
			return
		}

		// Ordering or terminating the BGP option changes the target specification, hence its checksum
		if (planBGPOption == nil) != (stateBGPOption == nil) {
			checksum = currentData.Checksum
		}

		updateData := planData.ToUpdate()
		updateData.Checksum = &checksum
		updateData.TargetSpec.Bgp = bgpTargetSpecToUpdate(planData.TargetSpec.Bgp, &currentData)
		if err := r.config.OVHClient.Put(endpoint, updateData, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Put %s", endpoint),
				err.Error(),
			)
			return
		}
	}

	// Wait for update
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for resource to be ready", err.Error())
		return
	}

	// Fetch up-to-date service info
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	responseData.MergeWith(&planData, false)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *vrackServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VrackServicesModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check terminate business rules
	services := data.TargetSpec.Services
	if !services.IsNull() && !services.IsUnknown() && !services.ServiceRange.IsNull() && !services.ServiceRange.IsUnknown() {
		for _, elem := range services.ServiceRange.Elements() {
			serviceRange := elem.(TargetSpecServicesServiceRangeValue)
			if serviceRange.ServiceEndpoints.IsNull() || serviceRange.ServiceEndpoints.IsUnknown() {
				continue
			}

			if len(serviceRange.ServiceEndpoints.Elements()) > 0 {
				resp.Diagnostics.AddError("failed to delete resource",
					"every existing ServiceEndpoints must be deleted before terminate the resource")
				return
			}
		}
	}

	// A subscribed BGP option is terminated along with the vRack Services
	if err := orderTerminateServiceByResourceName(r.config, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete resource", err.Error())
		return
	}
}

// bgpPlanOption returns the BGP option among the given plan options, nil when absent.
func bgpPlanOption(planOptions ovhtypes.TfListNestedValue[PlanOptionValue]) *PlanOptionValue {
	if planOptions.IsNull() || planOptions.IsUnknown() {
		return nil
	}

	for _, elem := range planOptions.Elements() {
		option := elem.(PlanOptionValue)
		if strings.HasPrefix(option.PlanCode.ValueString(), vrackServicesBGPPlanCodePrefix) {
			return &option
		}
	}

	return nil
}

// orderBGPOption orders the BGP option on a delivered vRack Services and waits for its delivery.
// Nothing is ordered when the option is already subscribed (e.g. ordered outside of Terraform).
func (r *vrackServicesResource) orderBGPOption(ctx context.Context, serviceName, ovhSubsidiary string, option PlanOptionValue) error {
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(serviceName)

	var vrackServices VrackServicesModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &vrackServices); err != nil {
		return fmt.Errorf("calling Get %s: %w", endpoint, err)
	}
	if !vrackServices.CurrentState.Bgp.IsNull() {
		log.Printf("[DEBUG] BGP option already subscribed on %s, nothing to order", serviceName)
		return nil
	}

	if err := orderServiceOption(ctx, r.config, "vrackServices", serviceName, ovhSubsidiary, option); err != nil {
		return err
	}

	return waitForVrackServicesBGP(ctx, r.config, endpoint)
}

// terminateBGPOption requests the termination of the BGP option of a vRack Services.
//   - Like a vRack Services deletion, it returns once the termination is confirmed: the BGP is
//     suspended then deleted later, when Agora processes the termination.
//   - Nothing is terminated when the option is not subscribed anymore.
func (r *vrackServicesResource) terminateBGPOption(ctx context.Context, serviceName string) error {
	endpoint := "/v2/vrackServices/resource/" + url.PathEscape(serviceName)

	var vrackServices VrackServicesModel
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &vrackServices); err != nil {
		return fmt.Errorf("calling Get %s: %w", endpoint, err)
	}
	if vrackServices.CurrentState.Bgp.IsNull() {
		log.Printf("[DEBUG] BGP option not subscribed on %s, nothing to terminate", serviceName)
		return nil
	}

	// The BGP option service is named after the BGP identifier
	return orderTerminateServiceByResourceName(r.config, vrackServices.CurrentState.Bgp.Id.ValueString())
}

// waitForVrackServicesBGP waits until the vRack Services is READY with its BGP option delivered.
func waitForVrackServicesBGP(ctx context.Context, config *Config, endpoint string) error {
	return retry.RetryContext(ctx, defaultOrderTimeout, func() *retry.RetryError {
		var vrackServices VrackServicesModel
		if err := config.OVHClient.GetWithContext(ctx, endpoint, &vrackServices); err != nil {
			return retry.RetryableError(fmt.Errorf("calling Get %s: %w", endpoint, err))
		}

		if vrackServices.CurrentState.Bgp.IsNull() {
			return retry.RetryableError(fmt.Errorf("waiting for BGP option delivery"))
		}

		if status := vrackServices.ResourceStatus.ValueString(); status != "READY" {
			return retry.RetryableError(fmt.Errorf("resource not ready (current status: %s)", status))
		}

		return nil
	})
}

// bgpTargetSpecToUpdate returns the targetSpec.bgp to send, given the remote vRack Services:
//   - the API rejects removing targetSpec.bgp while the BGP option is subscribed, so an unset
//     target_spec.bgp keeps the remote one;
//   - unset (computed) attributes of a configured target_spec.bgp keep their remote value.
func bgpTargetSpecToUpdate(planned TargetSpecBgpValue, remote *VrackServicesModel) *TargetSpecBgpWritableValue {
	remoteBGP := remote.TargetSpec.Bgp

	if planned.IsNull() || planned.IsUnknown() {
		if remoteBGP.IsNull() || remoteBGP.IsUnknown() {
			return nil
		}
		return remoteBGP.ToUpdate()
	}

	if !remoteBGP.IsNull() && !remoteBGP.IsUnknown() {
		if planned.AllowedPrefixes.IsUnknown() {
			planned.AllowedPrefixes = remoteBGP.AllowedPrefixes
		} else if !planned.AllowedPrefixes.IsNull() {
			if planned.AllowedPrefixes.Ipv4.IsUnknown() {
				planned.AllowedPrefixes.Ipv4 = remoteBGP.AllowedPrefixes.Ipv4
			}
			if planned.AllowedPrefixes.Ipv6.IsUnknown() {
				planned.AllowedPrefixes.Ipv6 = remoteBGP.AllowedPrefixes.Ipv6
			}
		}
		if planned.BgpRange.IsUnknown() {
			planned.BgpRange = remoteBGP.BgpRange
		}
		if planned.ClientAsn.IsUnknown() {
			planned.ClientAsn = remoteBGP.ClientAsn
		}
		if planned.InfrastructureAsn.IsUnknown() {
			planned.InfrastructureAsn = remoteBGP.InfrastructureAsn
		}
	}

	return planned.ToUpdate()
}
