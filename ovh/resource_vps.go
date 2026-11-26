package ovh

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*vpsResource)(nil)
	_ resource.ResourceWithImportState = (*vpsResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*vpsResource)(nil)
)

// vpsPollDelay and vpsPollInterval pace the waits for the order and for the
// VPS. The tests shorten them.
var (
	vpsPollDelay    = 10 * time.Second
	vpsPollInterval = 30 * time.Second
)

// vpsDeliveryPendingHint ends the error of a delivery that did not end in
// time.
const vpsDeliveryPendingHint = "The order is paid and kept in the state (order.order_id). " +
	"Terraform marks the resource as tainted: run terraform untaint on it, so that the next apply does not order a new VPS. " +
	"The next plan or apply checks the delivery again."

// NewVPSResource returns the ovh_vps resource.
func NewVPSResource() resource.Resource {
	return &vpsResource{}
}

type vpsResource struct {
	config *Config
}

func (r *vpsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps"
}

func (r *vpsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Orders and manages a VPS.",
		Attributes:  vpsResourceAttributes(),
	}
}

// ImportState imports a VPS by its service name. The Read that follows fills
// the state.
func (r *vpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ModifyPlan refuses, when the VPS is to be created, the settings that only an
// update of the ordered VPS could apply: an error after the payment would
// leave a VPS that differs from its configuration.
func (r *vpsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var config VPSResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(config.createRefusals()...)
}

func (r *vpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VPSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// The plan holds the values of the variables unknown at plan time.
	resp.Diagnostics.Append(data.createRefusals()...)
	if resp.Diagnostics.HasError() {
		return
	}

	order, diags := r.orderFor(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := orderCreate(order, r.config, "vps", false, 0); err != nil {
		resp.Diagnostics.AddError("Error ordering the VPS", err.Error())
		return
	}

	// Saved before the delivery, so that a timeout keeps the order.
	duration, pricingMode := data.orderTerms()
	data.Order = types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
		"duration":     tfStr(duration),
		"pricing_mode": tfStr(pricingMode),
		"order_id":     types.Int64Value(order.Order.OrderId.ValueInt64()),
	})
	resp.Diagnostics.Append(vpsSetState(ctx, &resp.State, &data, tftypes.NewValue(req.Plan.Raw.Type(), nil))...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.deliver(ctx, &data, &resp.State)...)
}

func (r *vpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VPSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		if vpsAttributeNull(data.Order, "order_id") {
			// Nothing was ordered.
			resp.State.RemoveResource(ctx)
			return
		}
		// The order was not delivered when Create gave up waiting.
		resp.Diagnostics.Append(r.checkDelivery(ctx, &data, &resp.State)...)
		return
	}

	found, diags := r.refresh(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update records the changes that leave the VPS as it is
// (vpsChangeWithoutEffect), and refuses any other.
func (r *vpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if vpsUpdatePlanned(req.Plan.Raw, req.State.Raw, req.Config.Raw) {
		resp.Diagnostics.AddError(
			"Updating a VPS is not supported yet",
			"This version of the provider only records, without calling the API, the changes of timeouts, order, installation.send_password_by_email, "+
				"and the values the state has none of, as after an import. Set the other attributes back to their value in the state.",
		)
		return
	}

	raw, err := vpsResolveUnknowns(req.Plan.Raw, req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Error updating the state of the VPS", err.Error())
		return
	}
	resp.State.Raw = raw
}

func (r *vpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VPSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"The VPS is not delivered yet",
			fmt.Sprintf("The VPS of order %s is not delivered yet, so it cannot be terminated. "+
				"Run terraform untaint on the resource and wait for its delivery with terraform apply, "+
				"or, if the order was cancelled, remove the resource from the state with terraform state rm.",
				vpsAttributeAt(data.Order, "order_id")),
		)
		return
	}

	id := data.Id.ValueString()
	base := "/vps/" + url.PathEscape(id)
	terminate := func() (string, error) {
		endpoint := base + "/terminate"
		if err := r.config.OVHClient.PostWithContext(ctx, endpoint, nil, nil); err != nil {
			if status := vpsAPIStatus(err); status == http.StatusNotFound || status == vpsV1StatusExpired {
				return "", nil
			}
			return "", fmt.Errorf("calling POST %s: %w", endpoint, err)
		}
		return id, nil
	}
	confirm := func(token string) error {
		endpoint := base + "/confirmTermination"
		if err := r.config.OVHClient.PostWithContext(ctx, endpoint, &ConfirmTerminationOpts{Token: token}, nil); err != nil {
			return fmt.Errorf("calling POST %s: %w", endpoint, err)
		}
		return nil
	}

	if err := orderDeleteWithTimeout(r.config, data.timeout("delete", vpsDeleteTimeout), terminate, confirm); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error terminating VPS %s", id), err.Error())
	}
}

// createRefusals refuses the settings that ordering a VPS cannot apply. Only
// known values count: a value unknown at plan time is checked again by Create.
func (m *VPSResourceModel) createRefusals() diag.Diagnostics {
	var diags diag.Diagnostics
	refuse := func(p path.Path, name string) {
		diags.AddAttributeError(p,
			"Not supported when creating a VPS",
			fmt.Sprintf("Setting %s when creating a VPS is not supported yet", name),
		)
	}
	set := func(value attr.Value) bool {
		return value != nil && !value.IsNull() && !value.IsUnknown()
	}

	installation := path.Root("installation")
	if set(vpsAttributeAt(m.Installation, "from", "image_id")) {
		refuse(installation.AtName("from").AtName("image_id"), "installation.from.image_id")
	}
	if set(vpsAttributeAt(m.Installation, "ssh_key")) {
		refuse(installation.AtName("ssh_key"), "installation.ssh_key")
	}
	if set(vpsAttributeAt(m.Installation, "post_installation_script")) {
		refuse(installation.AtName("post_installation_script"), "installation.post_installation_script")
	}
	if send, ok := vpsAttributeAt(m.Installation, "send_password_by_email").(types.Bool); ok && set(send) && !send.ValueBool() {
		refuse(installation.AtName("send_password_by_email"), "installation.send_password_by_email to false")
	}
	if set(m.PowerState) && m.PowerState.ValueString() != vpsDefaultPowerState {
		refuse(path.Root("power_state"), "power_state to "+m.PowerState.ValueString())
	}
	if set(m.DisplayName) {
		refuse(path.Root("display_name"), "display_name")
	}
	if set(vpsAttributeAt(m.Options, "auto_backup", "cron")) {
		refuse(path.Root("options").AtName("auto_backup").AtName("cron"), "options.auto_backup.cron")
	}
	return diags
}

// orderTerms returns the duration and the pricing mode of the order, or their
// defaults.
func (m *VPSResourceModel) orderTerms() (duration, pricingMode string) {
	duration, pricingMode = vpsDefaultDuration, vpsDefaultPricingMode
	if value, ok := vpsAttributeAt(m.Order, "duration").(ovhtypes.TfStringValue); ok && !value.IsNull() && !value.IsUnknown() {
		duration = value.ValueString()
	}
	if value, ok := vpsAttributeAt(m.Order, "pricing_mode").(ovhtypes.TfStringValue); ok && !value.IsNull() && !value.IsUnknown() {
		pricingMode = value.ValueString()
	}
	return duration, pricingMode
}

// orderFor checks the plan code, the region and the image against the order
// catalog of the subsidiary of the account, and returns the order of the VPS.
// Nothing is ordered yet.
func (r *vpsResource) orderFor(ctx context.Context, m *VPSResourceModel) (*OrderModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.config.OVHClient

	subsidiary, err := getOVHSubsidiary(ctx, client)
	if err != nil {
		diags.AddError("Error reading the subsidiary of your account", err.Error())
		return nil, diags
	}
	planCode := m.PlanCode.ValueString()
	offer, err := vpsOrderOfferFor(ctx, client, subsidiary, planCode)
	if err != nil {
		diags.AddAttributeError(path.Root("plan_code"), "Error reading the VPS order catalog", err.Error())
		return nil, diags
	}
	locations, err := vpsV1BackendOf(r.config).locations.get(ctx, client)
	if err != nil {
		diags.AddError("Error reading the regions",
			fmt.Sprintf("Ordering a VPS calls GET /v2/location to find the datacenter of its region: %s", err))
		return nil, diags
	}

	region := vpsStringAt(m.Location, "region")
	datacenter, ok := vpsDatacenterCodeForRegion(region, locations, offer.Datacenters)
	if !ok {
		diags.AddAttributeError(path.Root("location").AtName("region"),
			"Region not available for this plan",
			fmt.Sprintf("Plan %s cannot be ordered in region %s. Possible regions: %s",
				planCode, region, strings.Join(vpsOrderableRegions(offer.Datacenters, locations), ", ")),
		)
	}
	imageName := vpsStringAt(m.Installation, "from", "image_name")
	if !slices.Contains(offer.Images, imageName) {
		diags.AddAttributeError(path.Root("installation").AtName("from").AtName("image_name"),
			"Image not available for this plan",
			fmt.Sprintf("Plan %s cannot be ordered with image %q. Possible images: %s",
				planCode, imageName, strings.Join(offer.Images, ", ")),
		)
	}
	if diags.HasError() {
		return nil, diags
	}

	duration, pricingMode := m.orderTerms()
	order, orderDiags := vpsOrderModel(ctx, planCode, datacenter, imageName, duration, pricingMode, subsidiary)
	diags.Append(orderDiags...)
	return order, diags
}

// vpsOrderableRegions returns the regions of the datacenter codes, sorted; the
// code itself for a code that matches no single region.
func vpsOrderableRegions(codes []string, locations []vpsLocation) []string {
	regions := make([]string, 0, len(codes))
	for _, code := range codes {
		if region, ok := vpsRegionFromDatacenterCode(code, locations); ok {
			regions = append(regions, region)
			continue
		}
		regions = append(regions, code)
	}
	slices.Sort(regions)
	return slices.Compact(regions)
}

// deliver waits for the order of the VPS, then for the VPS to leave its
// transitional state, and reads it, all within the create timeout. Only Create
// waits: Read checks the delivery once (checkDelivery). It saves
// the state at each step: a timeout keeps the order and the id of the VPS.
func (r *vpsResource) deliver(ctx context.Context, data *VPSResourceModel, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics
	client := r.config.OVHClient
	deadline := time.Now().Add(data.timeout("create", vpsCreateTimeout))

	if data.Id.IsNull() || data.Id.IsUnknown() {
		orderIdValue, _ := vpsAttributeAt(data.Order, "order_id").(types.Int64)
		orderID := orderIdValue.ValueInt64()
		err := vpsWaitOrderDelivered(ctx, client, orderID, time.Until(deadline))
		if vpsTimedOut(err) {
			diags.AddError(fmt.Sprintf("VPS order %d is not delivered yet", orderID), vpsDeliveryPendingHint)
			return diags
		}
		if err != nil {
			diags.AddError(fmt.Sprintf("Error waiting for the delivery of VPS order %d", orderID), err.Error())
			return diags
		}

		id, err := serviceNameFromOrder(client, orderID, data.PlanCode.ValueString())
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reading the VPS of order %d", orderID), err.Error())
			return diags
		}
		data.Id = tfStr(id)
		diags.Append(state.Set(ctx, data)...)
		if diags.HasError() {
			return diags
		}
	}

	id := data.Id.ValueString()
	err := vpsWaitReady(ctx, client, id, time.Until(deadline))
	if vpsTimedOut(err) {
		diags.AddError(fmt.Sprintf("VPS %s is not ready yet", id), vpsDeliveryPendingHint)
		return diags
	}
	if err != nil {
		diags.AddError(fmt.Sprintf("Error waiting for VPS %s", id), err.Error())
		return diags
	}

	found, refreshDiags := r.refresh(ctx, data)
	diags.Append(refreshDiags...)
	if diags.HasError() {
		return diags
	}
	if !found {
		diags.AddError(fmt.Sprintf("VPS %s not found", id), "The VPS of the order is not found after its delivery.")
		return diags
	}
	diags.Append(state.Set(ctx, data)...)
	return diags
}

// checkDelivery checks once, without waiting, whether the order of the VPS is
// delivered: a plan never blocks on it. Delivered, it sets the id and reads the
// VPS as it is; otherwise the state stays as it is, with a warning.
func (r *vpsResource) checkDelivery(ctx context.Context, data *VPSResourceModel, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics
	client := r.config.OVHClient
	orderIdValue, _ := vpsAttributeAt(data.Order, "order_id").(types.Int64)
	orderID := orderIdValue.ValueInt64()

	var status string
	endpoint := fmt.Sprintf("/me/order/%d/status", orderID)
	if err := client.GetWithContext(ctx, endpoint, &status); err != nil {
		diags.AddError(fmt.Sprintf("Error reading VPS order %d", orderID), fmt.Sprintf("calling GET %s: %s", endpoint, err))
		return diags
	}
	if status != "delivered" {
		diags.AddWarning(
			fmt.Sprintf("VPS order %d is being delivered", orderID),
			fmt.Sprintf("The order is %s: the VPS is not in the state yet, plan or apply again later. "+
				"If Terraform marked the resource as tainted, run terraform untaint on it, so that the next apply does not order a new VPS.", status),
		)
		return diags
	}

	id, err := serviceNameFromOrder(client, orderID, data.PlanCode.ValueString())
	if err != nil {
		diags.AddError(fmt.Sprintf("Error reading the VPS of order %d", orderID), err.Error())
		return diags
	}
	data.Id = tfStr(id)

	found, refreshDiags := r.refresh(ctx, data)
	diags.Append(refreshDiags...)
	if diags.HasError() {
		return diags
	}
	if !found {
		// Kept, so that the next apply does not order a new VPS.
		diags.AddWarning(
			fmt.Sprintf("VPS %s cannot be read yet", id),
			fmt.Sprintf("Order %d is delivered, but GET /vps/%s does not find the VPS yet: its id is kept in the state, the next plan reads it.", orderID, id),
		)
	}
	diags.Append(state.Set(ctx, data)...)
	return diags
}

// refresh reads the VPS of data.Id and merges it into data. found is false
// when the VPS no longer exists.
func (r *vpsResource) refresh(ctx context.Context, data *VPSResourceModel) (found bool, diags diag.Diagnostics) {
	id := data.Id.ValueString()
	srv, readDiags, err := vpsV1BackendOf(r.config).read(ctx, id)
	diags.Append(readDiags...)
	switch {
	case errors.Is(err, errVPSNotFound):
		return false, diags
	case errors.Is(err, errVPSExpired):
		diags.AddError(
			fmt.Sprintf("VPS %s is expired", id),
			fmt.Sprintf("Renew the VPS to manage it again, or remove it from the state with terraform state rm. It stays in the state, so that Terraform does not order a new VPS: %s", err),
		)
		return true, diags
	case err != nil:
		diags.AddError(fmt.Sprintf("Error reading VPS %s", id), err.Error())
		return true, diags
	}

	var billing *ServiceBilling
	if data.billingNeeded() {
		billing, err = vpsBillingOf(ctx, r.config.OVHClient, id)
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reading the billing of VPS %s", id), err.Error())
			return true, diags
		}
	}

	diags.Append(data.mergeRead(ctx, srv, billing)...)
	return true, diags
}

// vpsWaitOrderDelivered waits at most timeout for the order to be delivered.
func vpsWaitOrderDelivered(ctx context.Context, client *ovhwrap.Client, orderID int64, timeout time.Duration) error {
	endpoint := fmt.Sprintf("/me/order/%d/status", orderID)
	conf := &retry.StateChangeConf{
		Pending: []string{"checking", "delivering"},
		Target:  []string{"delivered"},
		Refresh: func() (any, string, error) {
			var status string
			if err := client.GetWithContext(ctx, endpoint, &status); err != nil {
				return nil, "", fmt.Errorf("calling GET %s: %w", endpoint, err)
			}
			return status, status, nil
		},
		Timeout:      timeout,
		Delay:        vpsPollDelay,
		PollInterval: vpsPollInterval,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

// vpsWaitReady waits at most timeout for the VPS to leave its transitional
// states (installing, rebooting...). A VPS API v1 does not know yet is waited
// for too.
func vpsWaitReady(ctx context.Context, client *ovhwrap.Client, id string, timeout time.Duration) error {
	endpoint := "/vps/" + url.PathEscape(id)
	conf := &retry.StateChangeConf{
		Pending: []string{"unknown", "transitional"},
		Target:  []string{"ready"},
		Refresh: func() (any, string, error) {
			var vps vpsV1
			if err := client.GetWithContext(ctx, endpoint, &vps); err != nil {
				if vpsAPIStatus(err) == http.StatusNotFound {
					return vps, "unknown", nil
				}
				return nil, "", fmt.Errorf("calling GET %s: %w", endpoint, err)
			}
			if target, _ := vpsPowerStatesFromV1(vps.State); target == nil {
				return vps, "transitional", nil
			}
			return vps, "ready", nil
		},
		Timeout:      timeout,
		Delay:        vpsPollDelay,
		PollInterval: vpsPollInterval,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

// vpsTimedOut tells whether a wait ended on its timeout.
func vpsTimedOut(err error) bool {
	var timeout *retry.TimeoutError
	return errors.As(err, &timeout) || errors.Is(err, context.DeadlineExceeded)
}

// vpsSetState sets the state from data, the unknown values replaced by the
// ones of prior at the same path, or null; then reads data back from it.
func vpsSetState(ctx context.Context, state *tfsdk.State, data *VPSResourceModel, prior tftypes.Value) diag.Diagnostics {
	diags := state.Set(ctx, data)
	if diags.HasError() {
		return diags
	}
	raw, err := vpsResolveUnknowns(state.Raw, prior)
	if err != nil {
		diags.AddError("Error setting the state of the VPS", err.Error())
		return diags
	}
	state.Raw = raw
	diags.Append(state.Get(ctx, data)...)
	return diags
}

// vpsResolveUnknowns replaces each unknown value of raw by the known value of
// prior at the same path, or by null: a state holds no unknown value.
func vpsResolveUnknowns(raw, prior tftypes.Value) (tftypes.Value, error) {
	return tftypes.Transform(raw, func(p *tftypes.AttributePath, value tftypes.Value) (tftypes.Value, error) {
		if value.IsKnown() {
			return value, nil
		}
		if found, _, err := tftypes.WalkAttributePath(prior, p); err == nil {
			if priorValue, ok := found.(tftypes.Value); ok && priorValue.IsFullyKnown() {
				return priorValue, nil
			}
		}
		return tftypes.NewValue(value.Type(), nil), nil
	})
}

// vpsV1BackendOf returns the API v1 backend: the resource tells a VPS not
// found from an expired one, which the vpsBackend interface does not.
func vpsV1BackendOf(config *Config) *vpsV1Backend {
	return newVPSBackend(config).(*vpsV1Backend)
}

// vpsBillingOf returns the billing of the VPS: its plan code, and the duration
// and pricing mode of its renewal.
func vpsBillingOf(ctx context.Context, client *ovhwrap.Client, id string) (*ServiceBilling, error) {
	serviceId, err := serviceIdFromResourceName(client, id)
	if err != nil {
		return nil, err
	}

	var service Service
	endpoint := "/services/" + url.PathEscape(strconv.Itoa(serviceId))
	if err := client.GetWithContext(ctx, endpoint, &service); err != nil {
		return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}
	return &service.Billing, nil
}

// billingNeeded tells whether the state lacks a value that only the billing
// of the VPS holds, as after an import.
func (m *VPSResourceModel) billingNeeded() bool {
	return m.PlanCode.IsNull() ||
		vpsAttributeNull(m.Order, "duration") ||
		vpsAttributeNull(m.Order, "pricing_mode")
}

// mergeRead sets the model from a VPS read through API v1, and from its
// billing when not nil. The values API v1 cannot tell apart from the ones the
// configuration sets (plan code, region, image, installation settings, order)
// are only filled when null, as after an import.
func (m *VPSResourceModel) mergeRead(ctx context.Context, srv *VPSAPIServer, billing *ServiceBilling) diag.Diagnostics {
	var diags diag.Diagnostics

	var read VPSDataModel
	read.MergeWith(ctx, srv)

	m.Id = read.Id
	m.DisplayName = read.DisplayName
	// Null in a transitional state (installing, rebooting...): the requested
	// power state is kept.
	if !read.PowerState.IsNull() {
		m.PowerState = read.PowerState
	}
	m.Options = read.Options
	m.CurrentState = read.CurrentState

	if vpsAttributeNull(m.Location, "region") {
		m.Location = read.Location
	}

	installation := map[string]attr.Value{
		"from":                     types.ObjectNull(vpsInstallationFromAttrTypes()),
		"ssh_key":                  vpsTfStringPtr(nil),
		"post_installation_script": vpsTfStringPtr(nil),
		"send_password_by_email":   types.BoolNull(),
	}
	if !m.Installation.IsNull() && !m.Installation.IsUnknown() {
		maps.Copy(installation, m.Installation.Attributes())
	}
	if from, _ := installation["from"].(types.Object); vpsAttributeNull(from, "image_id") && vpsAttributeNull(from, "image_name") {
		installation["from"] = read.Installation.Attributes()["from"]
	}
	if installation["send_password_by_email"].IsNull() {
		installation["send_password_by_email"] = types.BoolValue(true)
	}
	var objectDiags diag.Diagnostics
	m.Installation, objectDiags = types.ObjectValue(vpsResourceInstallationAttrTypes(), installation)
	diags.Append(objectDiags...)

	if billing == nil {
		return diags
	}
	if m.PlanCode.IsNull() {
		m.PlanCode = vpsTfStringPtr(vpsNonEmpty(billing.Plan.Code))
	}
	order := map[string]attr.Value{
		"duration":     vpsTfStringPtr(nil),
		"pricing_mode": vpsTfStringPtr(nil),
		"order_id":     types.Int64Null(),
	}
	if !m.Order.IsNull() && !m.Order.IsUnknown() {
		maps.Copy(order, m.Order.Attributes())
	}
	if order["duration"].IsNull() {
		order["duration"] = vpsTfStringPtr(vpsNonEmpty(billing.Pricing.Duration))
	}
	if order["pricing_mode"].IsNull() {
		order["pricing_mode"] = vpsTfStringPtr(vpsNonEmpty(billing.Pricing.PricingMode))
	}
	m.Order, objectDiags = types.ObjectValue(vpsOrderAttrTypes(), order)
	diags.Append(objectDiags...)

	return diags
}

// vpsAttributeNull tells whether the attribute of an object is null, the
// object being null or unknown included.
func vpsAttributeNull(object types.Object, name string) bool {
	if object.IsNull() || object.IsUnknown() {
		return true
	}
	value, ok := object.Attributes()[name]
	return !ok || value.IsNull()
}

// vpsAttributeAt returns the value at the path of names under the object, nil
// when the object or one of the objects on the way is null or unknown.
func vpsAttributeAt(object types.Object, names ...string) attr.Value {
	var value attr.Value = object
	for _, name := range names {
		parent, ok := value.(types.Object)
		if !ok || parent.IsNull() || parent.IsUnknown() {
			return nil
		}
		value = parent.Attributes()[name]
	}
	return value
}

// vpsStringAt returns the string at the path of names under the object, empty
// when it is not set.
func vpsStringAt(object types.Object, names ...string) string {
	value, _ := vpsAttributeAt(object, names...).(ovhtypes.TfStringValue)
	return value.ValueString()
}
