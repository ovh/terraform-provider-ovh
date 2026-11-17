package ovh

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// derivedVolumeList plans the volume attribute the configuration leaves out.
//
// volumes and the deprecated volume_ids are two views of targetSpec.volumes, and
// both are Optional+Computed. When the configuration does not set this one, it
// keeps its state value, unless the other one changes: it is then unknown, as the
// API will recompute it.
func derivedVolumeList(other string) planmodifier.List {
	return derivedVolumeListModifier{other: other}
}

type derivedVolumeListModifier struct{ other string }

func (m derivedVolumeListModifier) Description(context.Context) string {
	return "Keeps the state value unless the other volume attribute changes"
}

func (m derivedVolumeListModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m derivedVolumeListModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Set in the configuration: the plan is that value.
	if !req.ConfigValue.IsNull() {
		return
	}
	// Creation: nothing to keep, the value is known after apply.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if m.otherChanged(ctx, req) {
		resp.PlanValue = types.ListUnknown(req.StateValue.ElementType(ctx))
		return
	}
	resp.PlanValue = req.StateValue
}

func (m derivedVolumeListModifier) otherChanged(ctx context.Context, req planmodifier.ListRequest) bool {
	var plan, state attr.Value
	if m.other == "volumes" {
		var p, s types.List
		req.Plan.GetAttribute(ctx, path.Root(m.other), &p)
		req.State.GetAttribute(ctx, path.Root(m.other), &s)
		plan, state = p, s
	} else {
		var p, s ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
		req.Plan.GetAttribute(ctx, path.Root(m.other), &p)
		req.State.GetAttribute(ctx, path.Root(m.other), &s)
		plan, state = p, s
	}
	return !plan.Equal(state)
}

// roleDefaultsToData gives a volume without role the role DATA, but only when
// volumes is set in the configuration. A schema Default cannot be used: the
// framework applies it to the elements of volumes taken from the state when the
// configuration uses volume_ids instead, turning a BOOT volume into DATA in the
// plan. The plan then differs from the state, and every computed attribute is
// marked unknown.
func roleDefaultsToData() planmodifier.String { return roleDefaultsToDataModifier{} }

type roleDefaultsToDataModifier struct{}

func (roleDefaultsToDataModifier) Description(context.Context) string {
	return "Defaults the role of a configured volume to DATA"
}

func (m roleDefaultsToDataModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (roleDefaultsToDataModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}
	var volumes types.List
	if d := req.Config.GetAttribute(ctx, path.Root("volumes"), &volumes); d.HasError() || volumes.IsNull() {
		return
	}
	resp.PlanValue = types.StringValue(instanceVolumeRoleData)
}

// keepStateUnlessConfigured keeps the state value of an attribute the
// configuration leaves out, null included. UseStateForUnknown skips a null state,
// so a null image_id stays unknown once the framework has marked it.
func keepStateUnlessConfigured() planmodifier.String { return keepStateUnlessConfiguredModifier{} }

type keepStateUnlessConfiguredModifier struct{}

func (keepStateUnlessConfiguredModifier) Description(context.Context) string {
	return "Keeps the state value, null included, when the configuration leaves the attribute out"
}

func (m keepStateUnlessConfiguredModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (keepStateUnlessConfiguredModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}

// instanceUpdatePlanned tells whether the plan changes a mutable attribute of the
// instance, or the instance is OUT_OF_SYNC.
//
// It differs from isUpdatePlanned on one point. Once any attribute differs from
// the state (volumes, when the configuration uses volume_ids), the framework marks
// every computed attribute the configuration leaves out as unknown, and the plan
// modifiers then restore them one by one, in no fixed order. An unknown that
// comes from there is not a change the user asked for, so it does not count.
func instanceUpdatePlanned(state tfsdk.State, plan tfsdk.Plan, config tfsdk.Config) bool {
	var sm, pm, cm map[string]tftypes.Value
	if state.Raw.As(&sm) != nil || plan.Raw.As(&pm) != nil || config.Raw.As(&cm) != nil {
		return true
	}
	if s, ok := sm["resource_status"]; ok {
		var v string
		if s.As(&v) == nil && v == "OUT_OF_SYNC" {
			return true
		}
	}
	attrs := append(append(append([]string{}, instanceMutableAttrs.Strings...), instanceMutableAttrs.Lists...), instanceMutableAttrs.CustomStringLists...)
	for _, name := range attrs {
		p := pm[name]
		if !p.IsKnown() && cm[name].IsNull() {
			continue
		}
		if !p.Equal(sm[name]) {
			return true
		}
	}
	return false
}

// instanceKeepStateUnlessUpdate marks a computed instance attribute unknown when
// an update is planned, and keeps its state value otherwise.
func instanceKeepStateUnlessUpdate() planmodifier.String {
	return instanceKeepStateString{}
}

type instanceKeepStateString struct{}

func (instanceKeepStateString) Description(context.Context) string {
	return "Unknown when an update is planned, the state value otherwise"
}

func (m instanceKeepStateString) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (instanceKeepStateString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if instanceUpdatePlanned(req.State, req.Plan, req.Config) {
		resp.PlanValue = types.StringUnknown()
		return
	}
	resp.PlanValue = req.StateValue
}

// instanceKeepStateObject is instanceKeepStateUnlessUpdate for an object.
func instanceKeepStateObject() planmodifier.Object { return instanceKeepStateObjectModifier{} }

type instanceKeepStateObjectModifier struct{}

func (instanceKeepStateObjectModifier) Description(context.Context) string {
	return "Unknown when an update is planned, the state value otherwise"
}

func (m instanceKeepStateObjectModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (instanceKeepStateObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if instanceUpdatePlanned(req.State, req.Plan, req.Config) {
		resp.PlanValue = types.ObjectUnknown(req.StateValue.AttributeTypes(ctx))
		return
	}
	resp.PlanValue = req.StateValue
}
