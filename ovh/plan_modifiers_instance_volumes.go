package ovh

import (
	"context"

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
