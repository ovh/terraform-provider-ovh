package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// OutOfSyncPlanModifier returns a plan modifier that forces an update when
// the resource_status is OUT_OF_SYNC. This state indicates the resource was
// modified outside of Terraform and needs to be re-synced via a PUT request.
func OutOfSyncPlanModifier() planmodifier.String {
	return outOfSyncPlanModifier{}
}

type outOfSyncPlanModifier struct{}

func (m outOfSyncPlanModifier) Description(_ context.Context) string {
	return "Forces update when resource is OUT_OF_SYNC"
}

func (m outOfSyncPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Forces update when resource is OUT_OF_SYNC"
}

func (m outOfSyncPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	if req.StateValue.ValueString() == "OUT_OF_SYNC" {
		resp.PlanValue = types.StringValue("READY")
	}
}

// MutableAttrs describes the mutable config attributes for a resource,
// grouped by Terraform type so isUpdatePlanned can compare them correctly.
type MutableAttrs struct {
	Strings           []string // checked with ovhtypes.TfStringValue
	Bools             []string // checked with types.Bool
	Int64s            []string // checked with types.Int64
	Lists             []string // checked with types.List
	Maps              []string // checked with types.Map
	Objects           []string // checked with types.Object
	CustomStringLists []string // checked with ovhtypes.TfListNestedValue[TfStringValue]
}

// isUpdatePlanned checks if resource_status is OUT_OF_SYNC or any mutable
// config attribute changed between state and plan.
func isUpdatePlanned(ctx context.Context, state tfsdk.State, plan tfsdk.Plan, cfg MutableAttrs) bool {
	var status ovhtypes.TfStringValue
	state.GetAttribute(ctx, path.Root("resource_status"), &status)
	if status.ValueString() == "OUT_OF_SYNC" {
		return true
	}
	for _, attr := range cfg.Strings {
		var sv, pv ovhtypes.TfStringValue
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.Bools {
		var sv, pv types.Bool
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.Int64s {
		var sv, pv types.Int64
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.Lists {
		var sv, pv types.List
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.Maps {
		var sv, pv types.Map
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.Objects {
		var sv, pv types.Object
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	for _, attr := range cfg.CustomStringLists {
		var sv, pv ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
		state.GetAttribute(ctx, path.Root(attr), &sv)
		plan.GetAttribute(ctx, path.Root(attr), &pv)
		if !sv.Equal(pv) {
			return true
		}
	}
	return false
}

// UnknownDuringUpdateStringModifier marks a computed string attribute as
// unknown when an update is planned, so Terraform accepts the new value.
func UnknownDuringUpdateStringModifier(cfg MutableAttrs) planmodifier.String {
	return unknownDuringUpdateString{cfg: cfg}
}

type unknownDuringUpdateString struct {
	cfg MutableAttrs
}

func (m unknownDuringUpdateString) Description(_ context.Context) string {
	return "Sets value to unknown during updates"
}

func (m unknownDuringUpdateString) MarkdownDescription(_ context.Context) string {
	return "Sets value to unknown during updates"
}

func (m unknownDuringUpdateString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if isUpdatePlanned(ctx, req.State, req.Plan, m.cfg) {
		resp.PlanValue = types.StringUnknown()
	}
}

// UnknownDuringUpdateObjectModifier marks a computed object attribute as
// unknown when an update is planned, so Terraform accepts the new value.
func UnknownDuringUpdateObjectModifier(cfg MutableAttrs) planmodifier.Object {
	return unknownDuringUpdateObject{cfg: cfg}
}

type unknownDuringUpdateObject struct {
	cfg MutableAttrs
}

func (m unknownDuringUpdateObject) Description(_ context.Context) string {
	return "Sets value to unknown during updates"
}

func (m unknownDuringUpdateObject) MarkdownDescription(_ context.Context) string {
	return "Sets value to unknown during updates"
}

func (m unknownDuringUpdateObject) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if isUpdatePlanned(ctx, req.State, req.Plan, m.cfg) {
		resp.PlanValue = types.ObjectUnknown(req.StateValue.AttributeTypes(ctx))
	}
}
