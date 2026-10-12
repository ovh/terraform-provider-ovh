package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// genericOrderSchema is shared, so this guards every order-based resource:
// changing an order-cart configuration value must force replacement
func TestDomainZoneConfigurationForcesNew(t *testing.T) {
	r := resourceDomainZone()

	state := &terraform.InstanceState{
		ID: "myzone.example.com",
		Attributes: map[string]string{
			"id":                           "myzone.example.com",
			"name":                         "myzone.example.com",
			"plan.#":                       "1",
			"plan.0.plan_code":             "zone",
			"plan.0.duration":              "P1Y",
			"plan.0.pricing_mode":          "create-default",
			"plan.0.configuration.#":       "2",
			"plan.0.configuration.0.label": "zone",
			"plan.0.configuration.0.value": "myzone.example.com",
			"plan.0.configuration.1.label": "template",
			"plan.0.configuration.1.value": "minimized",
		},
	}

	newConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"plan": []interface{}{
			map[string]interface{}{
				"plan_code":    "zone",
				"duration":     "P1Y",
				"pricing_mode": "create-default",
				"configuration": []interface{}{
					map[string]interface{}{"label": "zone", "value": "myzone.example.com"},
					map[string]interface{}{"label": "template", "value": "stub"},
				},
			},
		},
	})

	diff, err := r.Diff(context.Background(), state, newConfig, nil)
	if err != nil {
		t.Fatalf("Diff error: %s", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for the template change, got nil")
	}
	if !diff.RequiresNew() {
		t.Fatalf("changing template must force replacement, but diff does not require new: %#v", diff.Attributes)
	}
}
