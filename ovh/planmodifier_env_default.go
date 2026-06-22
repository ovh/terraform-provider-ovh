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
// If the env var is unset/empty the plan value is left unchanged so any
// downstream guard (e.g. an apply-time required-attribute check) still applies.
func EnvDefaultString(envVar string) planmodifier.String {
	return envDefaultString{envVar: envVar}
}

type envDefaultString struct {
	envVar string
}

func (m envDefaultString) Description(_ context.Context) string {
	return fmt.Sprintf("If unset in configuration, defaults to the %s environment variable", m.envVar)
}

func (m envDefaultString) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If unset in configuration, defaults to the `%s` environment variable", m.envVar)
}

func (m envDefaultString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// User set the attribute explicitly in config: keep their value untouched.
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// Config is null/unknown: fall back to the environment variable at plan time.
	envValue := os.Getenv(m.envVar)
	if envValue == "" {
		// No env value to inject. Leave the plan value unchanged so any
		// apply-time guard remains the final arbiter.
		return
	}

	resp.PlanValue = types.StringValue(envValue)
}
