package ovh

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// cloudProjectServiceNameOrEnv returns the service_name of an ovh_cloud_* data
// source: the configured value, or the OVH_CLOUD_PROJECT_SERVICE environment
// variable when the attribute is omitted. An error diagnostic is added when
// neither is set.
func cloudProjectServiceNameOrEnv(value ovhtypes.TfStringValue, diags *diag.Diagnostics) ovhtypes.TfStringValue {
	if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
		return value
	}
	if env := os.Getenv("OVH_CLOUD_PROJECT_SERVICE"); env != "" {
		return ovhtypes.NewTfStringValue(env)
	}
	diags.AddError(
		"Missing service_name",
		"The service_name attribute is required. Provide it in the data source "+
			"configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
	)
	return value
}
