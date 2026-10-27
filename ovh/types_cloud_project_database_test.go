package ovh

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCloudProjectDatabaseUpdateOpts_networkSerialization(t *testing.T) {
	t.Run("no network change - fields omitted", func(t *testing.T) {
		opts := CloudProjectDatabaseUpdateOpts{
			Description: "test",
			Plan:        "essential",
		}
		data, err := json.Marshal(opts)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		if strings.Contains(string(data), "networkId") {
			t.Errorf("networkId should be omitted when no network change, got: %s", data)
		}
		if strings.Contains(string(data), "subnetId") {
			t.Errorf("subnetId should be omitted when no network change, got: %s", data)
		}
	})

	t.Run("switch to public - sends null", func(t *testing.T) {
		opts := CloudProjectDatabaseUpdateOpts{
			Description: "test",
			Plan:        "essential",
			NetworkID:   json.RawMessage("null"),
			SubnetID:    json.RawMessage("null"),
		}
		data, err := json.Marshal(opts)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		if !strings.Contains(string(data), `"networkId":null`) {
			t.Errorf("expected null networkId, got: %s", data)
		}
		if !strings.Contains(string(data), `"subnetId":null`) {
			t.Errorf("expected null subnetId, got: %s", data)
		}
	})

	t.Run("switch to private - sends UUID", func(t *testing.T) {
		networkUUID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
		subnetUUID := "11111111-2222-3333-4444-555555555555"

		opts := CloudProjectDatabaseUpdateOpts{
			Description: "test",
			Plan:        "essential",
		}
		opts.NetworkID, _ = json.Marshal(networkUUID)
		opts.SubnetID, _ = json.Marshal(subnetUUID)

		data, err := json.Marshal(opts)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		if !strings.Contains(string(data), `"networkId":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`) {
			t.Errorf("expected UUID networkId, got: %s", data)
		}
		if !strings.Contains(string(data), `"subnetId":"11111111-2222-3333-4444-555555555555"`) {
			t.Errorf("expected UUID subnetId, got: %s", data)
		}
	})
}

func TestCloudProjectDatabaseNodesDiff(t *testing.T) {
	r := resourceCloudProjectDatabase()
	state := &terraform.InstanceState{
		ID: "42",
		Attributes: map[string]string{
			"id":             "42",
			"service_name":   "p",
			"engine":         "mongodb",
			"version":        "8.0",
			"plan":           "production",
			"flavor":         "db2-2",
			"description":    "d",
			"nodes.#":        "3",
			"nodes.0.region": "GRA",
			"nodes.1.region": "GRA",
			"nodes.2.region": "GRA",
		},
	}

	tests := []struct {
		name        string
		regions     []string
		existing    bool
		expectErr   bool
		requiresNew bool
	}{
		{"same nodes", []string{"GRA", "GRA", "GRA"}, true, false, false},
		{"add a node", []string{"GRA", "GRA", "GRA", "GRA"}, true, false, false},
		{"remove a node", []string{"GRA", "GRA"}, true, false, false},
		{"move every node", []string{"RBX", "RBX", "RBX"}, true, false, true},
		{"nodes in different regions", []string{"GRA", "RBX", "SBG"}, true, true, false},
		{"new cluster in different regions", []string{"GRA", "RBX"}, false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := []interface{}{}
			for _, region := range tt.regions {
				nodes = append(nodes, map[string]interface{}{"region": region})
			}
			cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
				"service_name": "p",
				"engine":       "mongodb",
				"version":      "8.0",
				"plan":         "production",
				"flavor":       "db2-2",
				"description":  "d",
				"nodes":        nodes,
			})
			s := state
			if !tt.existing {
				s = nil
			}

			diff, err := r.Diff(context.Background(), s, cfg, nil)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := diff != nil && diff.RequiresNew(); got != tt.requiresNew {
				t.Fatalf("expected RequiresNew = %v, got %v", tt.requiresNew, got)
			}
		})
	}
}
