package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func TestCloudProjectServiceNameOrEnv(t *testing.T) {
	t.Run("configured value wins over the environment", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		got := cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringValue("from-config"), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got.ValueString() != "from-config" {
			t.Fatalf("got %q, want %q", got.ValueString(), "from-config")
		}
	})

	t.Run("null value falls back to the environment", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		got := cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringNull(), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got.ValueString() != "from-env" {
			t.Fatalf("got %q, want %q", got.ValueString(), "from-env")
		}
	})

	t.Run("null value without environment is an error", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "")
		var diags diag.Diagnostics
		cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringNull(), &diags)
		if !diags.HasError() {
			t.Fatal("expected an error diagnostic")
		}
	})
}
