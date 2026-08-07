package ovh

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpperCaseString returns a plan modifier that normalizes a string attribute to
// upper case at plan time. It mirrors the SDKv2 StateFunc `strings.ToUpper`
// pattern used elsewhere in this provider for enum-backed attributes whose API
// canonicalizes the value.
//
// Use it when the API always echoes the value upper-cased (so state will hold
// the upper-case form) while the user may write it in any case. Normalizing the
// planned value keeps plan == state, avoiding an "inconsistent result after
// apply" error and, on a RequiresReplace attribute, a perpetual forced replace.
// Place it before RequiresReplace() so the comparison runs on the normalized
// value.
func UpperCaseString() planmodifier.String {
	return upperCaseString{}
}

type upperCaseString struct{}

func (m upperCaseString) Description(_ context.Context) string {
	return "Normalizes the value to upper case"
}

func (m upperCaseString) MarkdownDescription(_ context.Context) string {
	return "Normalizes the value to upper case"
}

func (m upperCaseString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = types.StringValue(strings.ToUpper(req.PlanValue.ValueString()))
}
