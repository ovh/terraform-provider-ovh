package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithUpgradeState = (*vpsResource)(nil)

// UpgradeState migrates the state of the former ovh_vps resource (version 0).
func (r *vpsResource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {PriorSchema: vpsResourceSchemaV0(), StateUpgrader: upgradeVPSStateV0},
	}
}

// upgradeVPSStateV0 keeps what the former state knows and the new schema holds:
// the id, the plan code and the terms of the order, the order id, the display
// name and the installation settings. The other attributes are null: the Read
// that follows fills them, as after an import.
func upgradeVPSStateV0(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var (
		id, serviceName, displayName, sshKey types.String
		doNotSendPassword                    types.Bool
		orderId                              types.Int64
		plans                                types.List
	)
	for p, target := range map[string]any{
		"id":                   &id,
		"service_name":         &serviceName,
		"display_name":         &displayName,
		"public_ssh_key":       &sshKey,
		"do_not_send_password": &doNotSendPassword,
		"plan":                 &plans,
	} {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(p), target)...)
	}
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("order").AtName("order_id"), &orderId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if id.IsNull() || id.ValueString() == "" {
		id = serviceName
	}
	planCode, duration, pricingMode := types.StringNull(), types.StringNull(), types.StringNull()
	if elements := plans.Elements(); len(elements) > 0 {
		if plan, ok := elements[0].(types.Object); ok && !plan.IsNull() {
			planCode = vpsV0PlanString(plan, "plan_code")
			duration = vpsV0PlanString(plan, "duration")
			pricingMode = vpsV0PlanString(plan, "pricing_mode")
		}
	}

	data := VPSResourceModel{
		Id:           ovhtypes.TfStringValue{StringValue: id},
		PlanCode:     ovhtypes.TfStringValue{StringValue: planCode},
		DisplayName:  ovhtypes.TfStringValue{StringValue: displayName},
		PowerState:   vpsTfStringPtr(nil),
		Location:     types.ObjectNull(vpsLocationAttrTypes()),
		Installation: types.ObjectNull(vpsResourceInstallationAttrTypes()),
		Options:      types.ObjectNull(vpsOptionsAttrTypes()),
		Order:        types.ObjectNull(vpsOrderAttrTypes()),
		Timeouts:     types.ObjectNull(vpsTimeoutsAttrTypes()),
		CurrentState: types.ObjectNull(vpsCurrentStateAttrTypes()),
	}

	sendPasswordByEmail := types.BoolNull()
	if !doNotSendPassword.IsNull() {
		sendPasswordByEmail = types.BoolValue(!doNotSendPassword.ValueBool())
	}
	if !sshKey.IsNull() || !sendPasswordByEmail.IsNull() {
		data.Installation = vpsObject(&resp.Diagnostics, vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from":                     types.ObjectNull(vpsInstallationFromAttrTypes()),
			"ssh_key":                  ovhtypes.TfStringValue{StringValue: sshKey},
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   sendPasswordByEmail,
		})
	}
	if !duration.IsNull() || !pricingMode.IsNull() || !orderId.IsNull() {
		data.Order = vpsObject(&resp.Diagnostics, vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     ovhtypes.TfStringValue{StringValue: duration},
			"pricing_mode": ovhtypes.TfStringValue{StringValue: pricingMode},
			"order_id":     orderId,
		})
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// vpsV0PlanString is the string attribute name of a plan of the former state,
// null if the plan has none.
func vpsV0PlanString(object types.Object, name string) types.String {
	if s, ok := object.Attributes()[name].(types.String); ok {
		return s
	}
	return types.StringNull()
}

// vpsObject is types.ObjectValue, appending its diagnostics to diags.
func vpsObject(diags *diag.Diagnostics, attrTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	object, d := types.ObjectValue(attrTypes, values)
	diags.Append(d...)
	return object
}

// vpsResourceSchemaV0 is a frozen copy of the schema of the former ovh_vps
// resource (schema version 0), from ovh/resource_vps_gen.go and the order
// attributes of ovh/order_resource_gen.go, with the same names, types and
// nesting. It only decodes the states it wrote, so
// it has no custom types, validators or plan modifiers.
func vpsResourceSchemaV0() *schema.Schema {
	str := func(optional bool) schema.StringAttribute {
		return schema.StringAttribute{Optional: optional, Computed: !optional}
	}
	optionalComputedString := schema.StringAttribute{Optional: true, Computed: true}
	optionalComputedInt64 := schema.Int64Attribute{Optional: true, Computed: true}
	optionalComputedStrings := schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true}
	configuration := func(computed bool) schema.ListNestedAttribute {
		return schema.ListNestedAttribute{
			Optional: true,
			Computed: computed,
			NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"label": schema.StringAttribute{Required: true},
				"value": schema.StringAttribute{Required: true},
			}},
		}
	}

	return &schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":                   str(false),
			"cluster":              str(false),
			"display_name":         optionalComputedString,
			"do_not_send_password": schema.BoolAttribute{Optional: true, Computed: true},
			"iam": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"display_name": str(false),
					"id":           str(false),
					"tags":         schema.MapAttribute{ElementType: types.StringType, Computed: true},
					"urn":          str(false),
				},
			},
			"image_id":     str(true),
			"keymap":       optionalComputedString,
			"memory_limit": optionalComputedInt64,
			"model": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"available_options":      optionalComputedStrings,
					"datacenter":             optionalComputedStrings,
					"disk":                   optionalComputedInt64,
					"maximum_additionnal_ip": optionalComputedInt64,
					"memory":                 optionalComputedInt64,
					"name":                   optionalComputedString,
					"offer":                  optionalComputedString,
					"vcore":                  optionalComputedInt64,
					"version":                optionalComputedString,
				},
			},
			"monitoring_ip_blocks": optionalComputedStrings,
			"name":                 optionalComputedString,
			"netboot_mode":         optionalComputedString,
			"offer_type":           optionalComputedString,
			"public_ssh_key":       str(true),
			"service_name":         str(false),
			"sla_monitoring":       schema.BoolAttribute{Optional: true, Computed: true},
			"state":                optionalComputedString,
			"vcore":                optionalComputedInt64,
			"zone":                 optionalComputedString,

			// OrderResourceSchema
			"order": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"date": str(false),
					"details": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"description":     str(false),
							"detail_type":     str(false),
							"domain":          str(false),
							"order_detail_id": schema.Int64Attribute{Computed: true},
							"quantity":        str(false),
						}},
					},
					"expiration_date": str(false),
					"order_id":        schema.Int64Attribute{Computed: true},
				},
			},
			"ovh_subsidiary": str(true),
			"plan": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"configuration": configuration(true),
					"duration":      schema.StringAttribute{Required: true},
					"item_id":       schema.Int64Attribute{Optional: true},
					"plan_code":     schema.StringAttribute{Required: true},
					"pricing_mode":  schema.StringAttribute{Required: true},
					"quantity":      schema.Int64Attribute{Optional: true},
				}},
			},
			"plan_option": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"configuration": configuration(false),
					"duration":      schema.StringAttribute{Required: true},
					"plan_code":     schema.StringAttribute{Required: true},
					"pricing_mode":  schema.StringAttribute{Required: true},
					"quantity":      schema.Int64Attribute{Required: true},
				}},
			},
		},
	}
}
