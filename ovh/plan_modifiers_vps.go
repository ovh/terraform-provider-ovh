package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const vpsRequiresReplaceDescription = "Changing it from one value to another terminates the VPS and orders a new one. Setting it when the state has none, as after an import, only changes the state."

// vpsRequiresReplaceIfPriorKnown replaces the VPS only when the state holds a
// value and the plan another one. A null state value, as after an import or the
// migration of the state of the former ovh_vps, never destroys the VPS.
func vpsRequiresReplaceIfPriorKnown() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = !req.StateValue.IsNull() && !req.StateValue.IsUnknown()
		},
		vpsRequiresReplaceDescription,
		vpsRequiresReplaceDescription,
	)
}

// vpsUpdateIgnoredAttributes are the root attributes whose change does not
// update the VPS.
var vpsUpdateIgnoredAttributes = map[string]bool{
	"id":            true,
	"order":         true,
	"timeouts":      true,
	"current_state": true,
}

// vpsCurrentStatePlanModifier keeps current_state from the state, unless the
// plan updates the VPS: it is then unknown until the update is applied.
func vpsCurrentStatePlanModifier() planmodifier.Object { return vpsCurrentStateModifier{} }

type vpsCurrentStateModifier struct{}

func (vpsCurrentStateModifier) Description(context.Context) string {
	return "Keeps the current state unless the plan updates the VPS"
}

func (m vpsCurrentStateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (vpsCurrentStateModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if vpsUpdatePlanned(req.Plan.Raw, req.State.Raw, req.Config.Raw) {
		resp.PlanValue = types.ObjectUnknown(req.PlanValue.AttributeTypes(ctx))
		return
	}
	resp.PlanValue = req.StateValue
}

// vpsUpdatePlanned tells whether the plan changes an attribute of the VPS.
func vpsUpdatePlanned(plan, state, config tftypes.Value) bool {
	var pm, sm, cm map[string]tftypes.Value
	if plan.As(&pm) != nil || state.As(&sm) != nil || config.As(&cm) != nil {
		return true
	}
	for name, p := range pm {
		if !vpsUpdateIgnoredAttributes[name] && vpsValueChanged(p, sm[name], cm[name]) {
			return true
		}
	}
	return false
}

// vpsValueChanged compares a plan value with the state value, object attribute
// by object attribute. Once any attribute differs from the state, the framework
// marks every computed attribute the configuration leaves out as unknown: such
// an unknown is not a change the user asked for, so it does not count.
func vpsValueChanged(plan, state, config tftypes.Value) bool {
	if !plan.IsKnown() {
		return !config.IsNull()
	}
	if _, ok := plan.Type().(tftypes.Object); !ok || plan.IsNull() || state.IsNull() {
		return !plan.Equal(state)
	}

	var pm, sm, cm map[string]tftypes.Value
	if plan.As(&pm) != nil || state.As(&sm) != nil {
		return true
	}
	if config.IsKnown() && !config.IsNull() && config.As(&cm) != nil {
		return true
	}
	for name, p := range pm {
		if vpsValueChanged(p, sm[name], cm[name]) {
			return true
		}
	}
	return false
}
