package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Attributes removed from a policy must be cleared from the state, not kept with
// their previous value.
func TestIamPolicyToMap_clearsRemovedAttributes(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceIamPolicy().Schema, map[string]interface{}{
		"name":               "policy",
		"description":        "previous description",
		"identities":         []interface{}{"urn:v1:eu:identity:account:xx1111-ovh"},
		"resources":          []interface{}{"urn:v1:eu:resource:vps:*"},
		"allow":              []interface{}{"vps:apiovh:reboot"},
		"except":             []interface{}{"vps:apiovh:reinstall"},
		"permissions_groups": []interface{}{"urn:v1:eu:permissionsGroup:ovh:defaultMaintenance"},
		"expired_at":         "2030-01-01T00:00:00Z",
		"conditions": []interface{}{
			map[string]interface{}{"operator": "AND"},
		},
	})

	// The policy as returned by the API once those attributes are removed
	pol := IamPolicy{
		Name:        "policy",
		Identities:  []string{"urn:v1:eu:identity:account:xx1111-ovh"},
		Resources:   []IamResource{{URN: "urn:v1:eu:resource:vps:*"}},
		Permissions: IamPermissions{Allow: []IamAction{{Action: "vps:apiovh:reboot"}}},
	}
	for k, v := range pol.ToMap() {
		if err := d.Set(k, v); err != nil {
			t.Fatalf("failed to set %s: %s", k, err)
		}
	}

	if got := d.Get("conditions").([]interface{}); len(got) != 0 {
		t.Errorf("conditions: expected empty, got %v", got)
	}
	if got := d.Get("except").(*schema.Set).Len(); got != 0 {
		t.Errorf("except: expected empty, got %d element(s)", got)
	}
	if got := d.Get("permissions_groups").(*schema.Set).Len(); got != 0 {
		t.Errorf("permissions_groups: expected empty, got %d element(s)", got)
	}
	if got := d.Get("description").(string); got != "" {
		t.Errorf("description: expected empty, got %q", got)
	}
	if got := d.Get("expired_at").(string); got != "" {
		t.Errorf("expired_at: expected empty, got %q", got)
	}
	if got := d.Get("allow").(*schema.Set).Len(); got != 1 {
		t.Errorf("allow: expected 1 element, got %d", got)
	}
}
