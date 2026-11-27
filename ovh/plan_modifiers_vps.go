package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const vpsRequiresReplaceDescription = "Changing it terminates the VPS and orders a new one. " +
	"This includes the first plan after an import or the migration of the state, when your configuration differs from the value read from the VPS. " +
	"Only when the value could not be read, and the state has none, setting it is just recorded in the state."

// vpsRequiresReplaceIfPriorKnown replaces the VPS only when the state holds a
// value and the plan another one. A null state value, when the value could not
// be read from the VPS, never destroys the VPS. After an import or the
// migration of the state, Read fills the value: a configuration that differs
// replaces the VPS.
func vpsRequiresReplaceIfPriorKnown() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = !req.StateValue.IsNull() && !req.StateValue.IsUnknown()
		},
		vpsRequiresReplaceDescription,
		vpsRequiresReplaceDescription,
	)
}

// vpsChangeWithoutEffect tells whether changing the attribute at the dotted
// path, whose state value is state, leaves the VPS as it is: the change is
// only recorded in the state. Update and the plan of current_state share it.
// installationKnown tells that Terraform installed the VPS (the private state
// holds installation_known): a null installation value in the state then
// means none, not one the state does not know.
func vpsChangeWithoutEffect(path string, state tftypes.Value, installationKnown bool) bool {
	switch path {
	case "id", "current_state", "order", "timeouts", "installation.send_password_by_email":
		return true
	case "installation.ssh_key", "installation.post_installation_script",
		"installation.from.image_id", "installation.from.image_name":
		// Null after an import or the migration of the state of the former
		// ovh_vps: the key and the script are never read back, and the image
		// is not when GET /vps/{serviceName}/images/current does not find it.
		return state.IsNull() && !installationKnown
	case "plan_code", "location.region", "power_state":
		// Null after an import or the migration of the state of the former
		// ovh_vps, or for power_state when read during a transitional state:
		// the VPS already has a value, which the state did not know. The next
		// Read reads the power state.
		return state.IsNull()
	}
	return false
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
	flags, diags := vpsReadPrivate(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if vpsUpdatePlanned(req.Plan.Raw, req.State.Raw, req.Config.Raw, flags.installationKnown) {
		resp.PlanValue = types.ObjectUnknown(req.PlanValue.AttributeTypes(ctx))
		return
	}
	resp.PlanValue = req.StateValue
}

// vpsUpdatePlanned tells whether the plan changes the VPS, which takes an API
// call, rather than only the state (vpsChangeWithoutEffect).
func vpsUpdatePlanned(plan, state, config tftypes.Value, installationKnown bool) bool {
	return vpsValueChanged("", plan, state, config, installationKnown)
}

// vpsValueChanged compares a plan value with the state value, object attribute
// by object attribute. Once any attribute differs from the state, the framework
// marks every computed attribute the configuration leaves out as unknown: such
// an unknown is not a change the user asked for, so it does not count.
func vpsValueChanged(path string, plan, state, config tftypes.Value, installationKnown bool) bool {
	if path != "" && vpsChangeWithoutEffect(path, state, installationKnown) {
		return false
	}
	if !plan.IsKnown() {
		return !config.IsNull()
	}
	if _, ok := plan.Type().(tftypes.Object); !ok || plan.IsNull() {
		return !plan.Equal(state)
	}

	var pm, sm, cm map[string]tftypes.Value
	if plan.As(&pm) != nil {
		return true
	}
	if state.IsKnown() && !state.IsNull() && state.As(&sm) != nil {
		return true
	}
	if config.IsKnown() && !config.IsNull() && config.As(&cm) != nil {
		return true
	}
	for name, p := range pm {
		s, ok := sm[name]
		if !ok {
			s = tftypes.NewValue(p.Type(), nil)
		}
		if vpsValueChanged(vpsJoinPath(path, name), p, s, cm[name], installationKnown) {
			return true
		}
	}
	return false
}

func vpsJoinPath(path, name string) string {
	if path == "" {
		return name
	}
	return path + "." + name
}
