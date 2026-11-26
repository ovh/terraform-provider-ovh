package ovh

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var (
	_ resource.ResourceWithConfigure   = (*vpsResource)(nil)
	_ resource.ResourceWithImportState = (*vpsResource)(nil)
)

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
		Description: "Orders and manages a VPS.",
		Attributes:  vpsResourceAttributes(),
	}
}

// ImportState imports a VPS by its service name. The Read that follows fills
// the state.
func (r *vpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError("Not implemented", "Creating a VPS is not implemented yet")
}

func (r *vpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VPSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	srv, diags, err := vpsV1BackendOf(r.config).read(ctx, id)
	resp.Diagnostics.Append(diags...)
	switch {
	case errors.Is(err, errVPSNotFound):
		resp.State.RemoveResource(ctx)
		return
	case errors.Is(err, errVPSExpired):
		resp.Diagnostics.AddError(
			fmt.Sprintf("VPS %s is expired", id),
			fmt.Sprintf("Renew the VPS to manage it again, or remove it from the state with terraform state rm. It stays in the state, so that Terraform does not order a new VPS: %s", err),
		)
		return
	case err != nil:
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading VPS %s", id), err.Error())
		return
	}

	var billing *ServiceBilling
	if data.billingNeeded() {
		billing, err = vpsBillingOf(ctx, r.config.OVHClient, id)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error reading the billing of VPS %s", id), err.Error())
			return
		}
	}

	resp.Diagnostics.Append(data.mergeRead(ctx, srv, billing)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Not implemented", "Updating a VPS is not implemented yet")
}

func (r *vpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Not implemented", "Deleting a VPS is not implemented yet")
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
