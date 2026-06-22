package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EnvDefaultString returns a plan modifier that replicates the SDKv2
// schema.EnvDefaultFunc behavior for terraform-plugin-framework string
// attributes.
//
// When the attribute is absent from the configuration (null or unknown), the
// modifier resolves its value from the given environment variable at plan time
// and uses that as the planned value. This keeps the planned value equal to the
// value that the CRUD layer (which also falls back to the env var) will persist
// in state, avoiding a phantom diff/replace on every plan when a user relies on
// the env var instead of setting the attribute explicitly.
//
// The attribute MUST be declared as both Optional and Computed for the
// framework to accept a planned value that differs from the (null) config.
//
// When required is true, the modifier raises a "missing value" diagnostic at
// plan time if the value cannot be resolved from the configuration, the
// environment variable, or pre-existing state. When required is false, an
// unresolvable value is left untouched.
func EnvDefaultString(envVar string, required bool) planmodifier.String {
	return envDefaultString{envVar: envVar, required: required}
}

type envDefaultString struct {
	envVar   string
	required bool
}

func (m envDefaultString) Description(_ context.Context) string {
	required := ""
	if m.required {
		required = " (required)"
	}
	return fmt.Sprintf("If unset in configuration, defaults to the %s environment variable%s", m.envVar, required)
}

func (m envDefaultString) MarkdownDescription(_ context.Context) string {
	required := ""
	if m.required {
		required = " (required)"
	}
	return fmt.Sprintf("If unset in configuration, defaults to the `%s` environment variable%s", m.envVar, required)
}

func (m envDefaultString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// 1. Destroy: req.Plan is null — nothing to default, never error.
	if req.Plan.Raw.IsNull() {
		return
	}
	// 2. Explicit config wins.
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}
	// 3. Fall back to the env var at plan time.
	if env := os.Getenv(m.envVar); env != "" {
		resp.PlanValue = types.StringValue(env)
		return
	}
	// 4. No config, no env: if the value is already in state, keep it (UseStateForUnknown supplies it). Never error here.
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		return
	}
	// 5. Truly unresolvable (fresh create, nothing anywhere). Error only if required.
	if m.required {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Missing %s", req.Path),
			fmt.Sprintf("Set it in the configuration or via the %s environment variable.", m.envVar),
		)
	}
}
