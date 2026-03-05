package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
