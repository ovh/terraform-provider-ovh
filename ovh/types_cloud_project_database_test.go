package ovh

import (
	"encoding/json"
	"strings"
	"testing"
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
