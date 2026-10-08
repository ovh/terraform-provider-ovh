package ovh

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func loadbalancerNetworkObject(id, subnetId string, ip attr.Value) types.Object {
	return types.ObjectValueMust(
		LoadbalancerNetworkAttrTypes(),
		map[string]attr.Value{
			"id":        ovhtypes.NewTfStringValue(id),
			"subnet_id": ovhtypes.NewTfStringValue(subnetId),
			"ip":        ip,
		},
	)
}

func TestCloudLoadbalancerToCreate_Network(t *testing.T) {
	model := CloudLoadbalancerModel{
		Name:       ovhtypes.NewTfStringValue("test-lb"),
		Region:     ovhtypes.NewTfStringValue("GRA11"),
		FlavorName: ovhtypes.NewTfStringValue("SMALL"),
		Network:    loadbalancerNetworkObject("net-id", "subnet-id", ovhtypes.NewTfStringValue("10.0.0.37")),
	}

	payload := model.ToCreate()
	if payload.TargetSpec == nil {
		t.Fatal("expected targetSpec to be non-nil")
	}

	if payload.TargetSpec.Network == nil {
		t.Fatal("expected targetSpec.network to be non-nil")
	}

	if payload.TargetSpec.Network.ID != "net-id" {
		t.Fatalf("unexpected targetSpec.network.id: got %q", payload.TargetSpec.Network.ID)
	}

	if payload.TargetSpec.Network.SubnetID != "subnet-id" {
		t.Fatalf("unexpected targetSpec.network.subnetId: got %q", payload.TargetSpec.Network.SubnetID)
	}

	if payload.TargetSpec.Network.IP != "10.0.0.37" {
		t.Fatalf("unexpected targetSpec.network.ip: got %q", payload.TargetSpec.Network.IP)
	}
}

func TestCloudLoadbalancerToCreate_NetworkWithoutIP(t *testing.T) {
	model := CloudLoadbalancerModel{
		Name:       ovhtypes.NewTfStringValue("test-lb"),
		Region:     ovhtypes.NewTfStringValue("GRA11"),
		FlavorName: ovhtypes.NewTfStringValue("SMALL"),
		Network:    loadbalancerNetworkObject("net-id", "subnet-id", ovhtypes.TfStringValue{StringValue: types.StringNull()}),
	}

	payload := model.ToCreate()
	if payload.TargetSpec.Network.IP != "" {
		t.Fatalf("expected targetSpec.network.ip to be empty, got %q", payload.TargetSpec.Network.IP)
	}
}

func TestCloudLoadbalancerToUpdate_DoesNotIncludeNetwork(t *testing.T) {
	model := CloudLoadbalancerModel{
		Name:        ovhtypes.NewTfStringValue("test-lb"),
		Description: ovhtypes.NewTfStringValue("desc"),
		Network:     loadbalancerNetworkObject("net-id", "subnet-id", ovhtypes.NewTfStringValue("10.0.0.37")),
	}

	payload := model.ToUpdate("checksum-123")
	if payload.Checksum != "checksum-123" {
		t.Fatalf("unexpected checksum: got %q", payload.Checksum)
	}

	if payload.TargetSpec.Name != "test-lb" {
		t.Fatalf("unexpected targetSpec.name: got %q", payload.TargetSpec.Name)
	}

	if payload.TargetSpec.Description != "desc" {
		t.Fatalf("unexpected targetSpec.description: got %q", payload.TargetSpec.Description)
	}
}

func TestCloudLoadbalancerToUpdate_IncludesFlavor(t *testing.T) {
	model := CloudLoadbalancerModel{
		Name:       ovhtypes.NewTfStringValue("test-lb"),
		FlavorName: ovhtypes.NewTfStringValue("MEDIUM"),
		Network:    loadbalancerNetworkObject("net-id", "subnet-id", ovhtypes.NewTfStringValue("10.0.0.37")),
	}

	payload := model.ToUpdate("checksum-123")
	if payload.TargetSpec.Flavor == nil {
		t.Fatal("expected targetSpec.flavor to be non-nil")
	}

	if payload.TargetSpec.Flavor.Name != "MEDIUM" {
		t.Fatalf("unexpected targetSpec.flavor.name: got %q", payload.TargetSpec.Flavor.Name)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("unexpected marshal error: %s", err)
	}

	var raw struct {
		TargetSpec map[string]any `json:"targetSpec"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("unexpected unmarshal error: %s", err)
	}

	if _, ok := raw.TargetSpec["network"]; ok {
		t.Fatal("expected targetSpec.network to be absent from the update payload")
	}

	if _, ok := raw.TargetSpec["location"]; ok {
		t.Fatal("expected targetSpec.location to be absent from the update payload")
	}
}

func TestCloudLoadbalancerToUpdate_OmitsUnknownFlavor(t *testing.T) {
	model := CloudLoadbalancerModel{
		Name:       ovhtypes.NewTfStringValue("test-lb"),
		FlavorName: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}

	payload := model.ToUpdate("checksum-123")
	if payload.TargetSpec.Flavor != nil {
		t.Fatalf("expected targetSpec.flavor to be nil, got %+v", payload.TargetSpec.Flavor)
	}
}

func TestCloudLoadbalancerMergeWith_FlavorFromTargetSpec(t *testing.T) {
	model := CloudLoadbalancerModel{
		FlavorName: ovhtypes.NewTfStringValue("SMALL"),
		Network:    types.ObjectNull(LoadbalancerNetworkAttrTypes()),
	}

	model.MergeWith(context.Background(), &CloudLoadbalancerAPIResponse{
		Id:             "lb-id",
		ResourceStatus: "READY",
		TargetSpec: &CloudLoadbalancerAPITargetSpec{
			Name:   "test-lb",
			Flavor: &CloudLoadbalancerAPIFlavorRef{Name: "LARGE"},
		},
	})

	if got := model.FlavorName.ValueString(); got != "LARGE" {
		t.Fatalf("unexpected flavor_name: got %q", got)
	}
}

func TestCloudLoadbalancerSchema_FlavorNameIsMutable(t *testing.T) {
	r := &cloudLoadbalancerResource{}
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	flavorName, ok := resp.Schema.Attributes["flavor_name"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected flavor_name to be a StringAttribute, got %T", resp.Schema.Attributes["flavor_name"])
	}

	if !flavorName.Required {
		t.Fatal("expected flavor_name to be required")
	}

	if len(flavorName.PlanModifiers) != 0 {
		t.Fatalf("expected flavor_name to have no plan modifier (no RequiresReplace), got %d", len(flavorName.PlanModifiers))
	}

	if !slices.Contains(loadbalancerMutableAttrs.Strings, "flavor_name") {
		t.Fatal("expected flavor_name to be listed in loadbalancerMutableAttrs")
	}
}

func TestCloudLoadbalancerMergeWith_CurrentStateAddresses(t *testing.T) {
	model := CloudLoadbalancerModel{
		Network: types.ObjectNull(LoadbalancerNetworkAttrTypes()),
	}

	model.MergeWith(context.Background(), &CloudLoadbalancerAPIResponse{
		Id:             "lb-id",
		ResourceStatus: "READY",
		TargetSpec: &CloudLoadbalancerAPITargetSpec{
			Name:     "test-lb",
			Location: &CloudLoadbalancerAPILocation{Region: "GRA11"},
			Network: &CloudLoadbalancerAPINetworkRef{
				ID:       "net-id",
				SubnetID: "subnet-id",
			},
			Flavor: &CloudLoadbalancerAPIFlavorRef{Name: "SMALL"},
		},
		CurrentState: &CloudLoadbalancerAPICurrentState{
			Name:     "test-lb",
			Location: &CloudLoadbalancerAPILocation{Region: "GRA11"},
			Network: &CloudLoadbalancerAPINetwork{
				ID:       "net-id",
				SubnetID: "subnet-id",
				Addresses: []CloudLoadbalancerAPIAddress{
					{IP: "10.0.0.37", Type: "FIXED"},
					{IP: "51.1.2.3", Type: "FLOATING"},
				},
			},
		},
	})

	networkAttrs := model.Network.Attributes()
	if got := loadbalancerObjectString(networkAttrs, "subnet_id"); got != "subnet-id" {
		t.Fatalf("unexpected network.subnet_id: got %q", got)
	}

	if !networkAttrs["ip"].IsNull() {
		t.Fatalf("expected network.ip to stay null, got %v", networkAttrs["ip"])
	}

	currentNetwork, ok := model.CurrentState.Attributes()["network"].(types.Object)
	if !ok {
		t.Fatalf("expected current_state.network to be an object, got %T", model.CurrentState.Attributes()["network"])
	}

	addresses, ok := currentNetwork.Attributes()["addresses"].(types.List)
	if !ok {
		t.Fatalf("expected current_state.network.addresses to be a list, got %T", currentNetwork.Attributes()["addresses"])
	}

	elems := addresses.Elements()
	if len(elems) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(elems))
	}

	first := elems[0].(types.Object).Attributes()
	if got := loadbalancerObjectString(first, "ip"); got != "10.0.0.37" {
		t.Fatalf("unexpected addresses[0].ip: got %q", got)
	}
	if got := loadbalancerObjectString(first, "type"); got != "FIXED" {
		t.Fatalf("unexpected addresses[0].type: got %q", got)
	}

	second := elems[1].(types.Object).Attributes()
	if got := loadbalancerObjectString(second, "type"); got != "FLOATING" {
		t.Fatalf("unexpected addresses[1].type: got %q", got)
	}
}

func TestCloudLoadbalancerMergeWith_KeepsConfiguredVipIP(t *testing.T) {
	model := CloudLoadbalancerModel{
		Network: loadbalancerNetworkObject("net-id", "subnet-id", ovhtypes.NewTfStringValue("51.1.2.3")),
	}

	model.MergeWith(context.Background(), &CloudLoadbalancerAPIResponse{
		Id: "lb-id",
		TargetSpec: &CloudLoadbalancerAPITargetSpec{
			Name:    "test-lb",
			Network: &CloudLoadbalancerAPINetworkRef{ID: "net-id", SubnetID: "subnet-id"},
		},
	})

	if got := loadbalancerObjectString(model.Network.Attributes(), "ip"); got != "51.1.2.3" {
		t.Fatalf("expected configured network.ip to be preserved, got %q", got)
	}
}

func TestCloudLoadbalancerSchema_NetworkRequiresReplace(t *testing.T) {
	r := &cloudLoadbalancerResource{}
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if _, ok := resp.Schema.Attributes["vip_subnet_id"]; ok {
		t.Fatal("expected vip_subnet_id attribute to be removed")
	}

	if _, ok := resp.Schema.Attributes["network_id"]; ok {
		t.Fatal("expected network_id attribute to be removed")
	}

	networkAttr, ok := resp.Schema.Attributes["network"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected network attribute to be SingleNestedAttribute, got %T", resp.Schema.Attributes["network"])
	}

	if !networkAttr.Required {
		t.Fatal("expected network to be required")
	}

	if len(networkAttr.PlanModifiers) == 0 {
		t.Fatal("expected network object attribute to have a RequiresReplace plan modifier")
	}

	id, ok := networkAttr.Attributes["id"].(schema.StringAttribute)
	if !ok || !id.Required {
		t.Fatal("expected network.id to be a required string attribute")
	}

	subnetId, ok := networkAttr.Attributes["subnet_id"].(schema.StringAttribute)
	if !ok || !subnetId.Required {
		t.Fatal("expected network.subnet_id to be a required string attribute")
	}

	ip, ok := networkAttr.Attributes["ip"].(schema.StringAttribute)
	if !ok || !ip.Optional || ip.Computed {
		t.Fatal("expected network.ip to be an optional, non-computed string attribute")
	}

	currentState, ok := resp.Schema.Attributes["current_state"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected current_state to be SingleNestedAttribute, got %T", resp.Schema.Attributes["current_state"])
	}

	if _, ok := currentState.Attributes["vip_address"]; ok {
		t.Fatal("expected current_state.vip_address to be removed")
	}

	if _, ok := currentState.Attributes["vip_subnet"]; ok {
		t.Fatal("expected current_state.vip_subnet to be removed")
	}

	currentNetwork, ok := currentState.Attributes["network"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected current_state.network to be SingleNestedAttribute, got %T", currentState.Attributes["network"])
	}

	if _, ok := currentNetwork.Attributes["subnet_id"]; !ok {
		t.Fatal("expected current_state.network.subnet_id")
	}

	if _, ok := currentNetwork.Attributes["addresses"].(schema.ListNestedAttribute); !ok {
		t.Fatalf("expected current_state.network.addresses to be ListNestedAttribute, got %T", currentNetwork.Attributes["addresses"])
	}
}
