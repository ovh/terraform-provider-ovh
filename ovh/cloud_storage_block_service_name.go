package ovh

import (
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// requireResolvedServiceName is an apply-time guard for block storage resources
// whose `service_name` attribute is Optional+Computed and resolved from the
// OVH_CLOUD_PROJECT_SERVICE environment variable at plan time via the
// EnvDefaultString plan modifier.
//
// By the time Create runs, the plan modifier has already injected the env value
// (if any) into the planned value, so `serviceName` holds the resolved value.
// This guard only catches the pathological case where neither the configuration
// nor the environment variable provided a value, mirroring the ssh_key
// "Missing service_name" diagnostic.
func requireResolvedServiceName(serviceName ovhtypes.TfStringValue, diags interface{ AddError(string, string) }) bool {
	if serviceName.IsNull() || serviceName.IsUnknown() || serviceName.ValueString() == "" {
		diags.AddError(
			"Missing service_name",
			"The service_name attribute is required. Please provide it in the resource configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
		)
		return false
	}
	return true
}
