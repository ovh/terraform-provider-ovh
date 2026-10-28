package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func testPrivateEndpointModel() CloudPrivateEndpointModel {
	return CloudPrivateEndpointModel{
		Name:        ovhtypes.NewTfStringValue("my-endpoint"),
		Description: ovhtypes.NewTfStringValue("an endpoint"),
		Region:      ovhtypes.NewTfStringValue("GRA11"),
		Service: types.ObjectValueMust(privateEndpointServiceAttrTypes(), map[string]attr.Value{
			"name": ovhtypes.NewTfStringValue("s3"),
			"tier": ovhtypes.NewTfStringValue("S"),
		}),
		SubnetId:         ovhtypes.NewTfStringValue("subnet-123"),
		VipAddress:       ovhtypes.NewTfStringValue("10.0.0.42"),
		PublicDnsEnabled: types.BoolValue(true),
	}
}

// TestCloudPrivateEndpointToCreate_SinglePlacement verifies the create
// payload carries the whole target spec and only the placement reference that
// was configured.
func TestCloudPrivateEndpointToCreate_SinglePlacement(t *testing.T) {
	model := testPrivateEndpointModel()

	payload := model.ToCreate()
	if payload == nil || payload.TargetSpec == nil {
		t.Fatal("expected non-nil create payload and target spec")
	}
	ts := payload.TargetSpec

	if ts.Name != "my-endpoint" {
		t.Fatalf("unexpected name: got %q", ts.Name)
	}
	if ts.Description != "an endpoint" {
		t.Fatalf("unexpected description: got %q", ts.Description)
	}
	if ts.Location == nil || ts.Location.Region != "GRA11" {
		t.Fatalf("unexpected location: got %+v", ts.Location)
	}
	if ts.Service == nil || ts.Service.Name != "s3" || ts.Service.Tier != "S" {
		t.Fatalf("unexpected service: got %+v", ts.Service)
	}
	if !ts.PublicDnsEnabled {
		t.Fatal("expected publicDnsEnabled to be true")
	}
	if ts.Subnet == nil || ts.Subnet.Id != "subnet-123" {
		t.Fatalf("unexpected subnet: got %+v", ts.Subnet)
	}
	if ts.VipAddress != "10.0.0.42" {
		t.Fatalf("unexpected vipAddress: got %q", ts.VipAddress)
	}
	if ts.Network != nil {
		t.Fatalf("expected network to be nil, got %+v", ts.Network)
	}
}

// TestCloudPrivateEndpointToCreate_OmittedVipAddress verifies an unset
// vip_address is left out of the payload so the server allocates one.
func TestCloudPrivateEndpointToCreate_OmittedVipAddress(t *testing.T) {
	model := testPrivateEndpointModel()
	model.VipAddress = ovhtypes.TfStringValue{StringValue: types.StringUnknown()}

	ts := model.ToCreate().TargetSpec
	if ts.VipAddress != "" {
		t.Fatalf("expected vipAddress to be empty, got %q", ts.VipAddress)
	}
}

// TestCloudPrivateEndpointToUpdate_OnlyMutableFields verifies the update
// payload carries the checksum plus only the mutable fields (name and
// description), excluding placement, location and service.
func TestCloudPrivateEndpointToUpdate_OnlyMutableFields(t *testing.T) {
	model := testPrivateEndpointModel()

	payload := model.ToUpdate("checksum-abc")
	if payload == nil || payload.TargetSpec == nil {
		t.Fatal("expected non-nil update payload and target spec")
	}
	if payload.Checksum != "checksum-abc" {
		t.Fatalf("unexpected checksum: got %q", payload.Checksum)
	}
	if payload.TargetSpec.Name != "my-endpoint" {
		t.Fatalf("unexpected name: got %q", payload.TargetSpec.Name)
	}
	if payload.TargetSpec.Description != "an endpoint" {
		t.Fatalf("unexpected description: got %q", payload.TargetSpec.Description)
	}
}

// TestCloudPrivateEndpointMergeWith_PreservesNullDescription verifies a
// null description in config stays null after a response with an empty
// description (avoids a null-versus-empty-string perpetual diff).
func TestCloudPrivateEndpointMergeWith_PreservesNullDescription(t *testing.T) {
	model := CloudPrivateEndpointModel{
		Description: ovhtypes.TfStringValue{StringValue: types.StringNull()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id: "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{
			Name:     "my-endpoint",
			Location: &CloudPrivateEndpointLocation{Region: "GRA11"},
		},
	}

	model.MergeWith(context.Background(), response)

	if !model.Description.IsNull() {
		t.Fatalf("expected description to remain null, got %q", model.Description.ValueString())
	}
}

// TestCloudPrivateEndpointMergeWith_ClearsDroppedDescription verifies a
// description cleared outside Terraform is reflected in the model, so the next
// plan shows the drift instead of silently keeping the stale value.
func TestCloudPrivateEndpointMergeWith_ClearsDroppedDescription(t *testing.T) {
	model := CloudPrivateEndpointModel{
		Description: ovhtypes.NewTfStringValue("an endpoint"),
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id:         "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{Name: "my-endpoint"},
	}

	model.MergeWith(context.Background(), response)

	if model.Description.ValueString() != "" {
		t.Fatalf("expected description to be cleared, got %q", model.Description.ValueString())
	}
}

// TestCloudPrivateEndpointMergeWith_VipAddressFromTargetSpec verifies a
// caller-pinned address echoed in targetSpec lands on the root attribute.
func TestCloudPrivateEndpointMergeWith_VipAddressFromTargetSpec(t *testing.T) {
	model := CloudPrivateEndpointModel{
		VipAddress: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id: "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{
			Name:       "my-endpoint",
			Subnet:     &CloudPrivateEndpointRef{Id: "subnet-123"},
			VipAddress: "10.0.0.42",
		},
		CurrentState: &CloudPrivateEndpointCurrentState{
			VipAddress: "10.0.0.42",
			Subnet:     &CloudPrivateEndpointRef{Id: "subnet-123"},
		},
	}

	model.MergeWith(context.Background(), response)

	if model.VipAddress.ValueString() != "10.0.0.42" {
		t.Fatalf("unexpected vip_address: got %q", model.VipAddress.ValueString())
	}
	subnet := model.CurrentState.Attributes()["subnet"].(types.Object)
	if got := subnet.Attributes()["id"].(ovhtypes.TfStringValue).ValueString(); got != "subnet-123" {
		t.Fatalf("unexpected current_state.subnet.id: got %q", got)
	}
}

// TestCloudPrivateEndpointMergeWith_VipAddressFromCurrentState verifies
// the Computed vip_address is filled from currentState when the caller did not
// pin one (targetSpec then carries no address).
func TestCloudPrivateEndpointMergeWith_VipAddressFromCurrentState(t *testing.T) {
	model := CloudPrivateEndpointModel{
		VipAddress: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id: "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{
			Name:   "my-endpoint",
			Subnet: &CloudPrivateEndpointRef{Id: "subnet-123"},
		},
		CurrentState: &CloudPrivateEndpointCurrentState{
			VipAddress: "10.0.0.7",
		},
	}

	model.MergeWith(context.Background(), response)

	if model.VipAddress.ValueString() != "10.0.0.7" {
		t.Fatalf("unexpected vip_address: got %q", model.VipAddress.ValueString())
	}
}

// TestCloudPrivateEndpointMergeWith_VipAddressUnknownResolvesToNull
// verifies a Computed vip_address never stays unknown, which would break the
// apply with a missing-value error.
func TestCloudPrivateEndpointMergeWith_VipAddressUnknownResolvesToNull(t *testing.T) {
	model := CloudPrivateEndpointModel{
		VipAddress: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id: "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{
			Name:    "my-endpoint",
			Network: &CloudPrivateEndpointRef{Id: "network-123"},
		},
	}

	model.MergeWith(context.Background(), response)

	if model.VipAddress.IsUnknown() {
		t.Fatal("expected vip_address to be resolved, got unknown")
	}
	if !model.VipAddress.IsNull() {
		t.Fatalf("expected vip_address to be null, got %q", model.VipAddress.ValueString())
	}
}

// TestCloudPrivateEndpointMergeWith_Placement covers the three placement
// shapes the API accepts: the counterpart the platform resolved must land on the
// root attribute, whichever reference the caller supplied.
func TestCloudPrivateEndpointMergeWith_Placement(t *testing.T) {
	cases := []struct {
		name         string
		targetSpec   *CloudPrivateEndpointTargetSpec
		currentState *CloudPrivateEndpointCurrentState
		wantNetwork  string
		wantSubnet   string
	}{
		{
			name: "network only in targetSpec, subnet resolved",
			targetSpec: &CloudPrivateEndpointTargetSpec{
				Network: &CloudPrivateEndpointRef{Id: "network-123"},
			},
			currentState: &CloudPrivateEndpointCurrentState{
				Network: &CloudPrivateEndpointRef{Id: "network-123"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-456"},
			},
			wantNetwork: "network-123",
			wantSubnet:  "subnet-456",
		},
		{
			name: "subnet only in targetSpec, network resolved",
			targetSpec: &CloudPrivateEndpointTargetSpec{
				Subnet: &CloudPrivateEndpointRef{Id: "subnet-456"},
			},
			currentState: &CloudPrivateEndpointCurrentState{
				Network: &CloudPrivateEndpointRef{Id: "network-123"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-456"},
			},
			wantNetwork: "network-123",
			wantSubnet:  "subnet-456",
		},
		{
			name: "both in targetSpec",
			targetSpec: &CloudPrivateEndpointTargetSpec{
				Network: &CloudPrivateEndpointRef{Id: "network-123"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-456"},
			},
			currentState: &CloudPrivateEndpointCurrentState{
				Network: &CloudPrivateEndpointRef{Id: "network-123"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-456"},
			},
			wantNetwork: "network-123",
			wantSubnet:  "subnet-456",
		},
		{
			name: "targetSpec wins over currentState",
			targetSpec: &CloudPrivateEndpointTargetSpec{
				Network: &CloudPrivateEndpointRef{Id: "network-from-spec"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-from-spec"},
			},
			currentState: &CloudPrivateEndpointCurrentState{
				Network: &CloudPrivateEndpointRef{Id: "network-from-state"},
				Subnet:  &CloudPrivateEndpointRef{Id: "subnet-from-state"},
			},
			wantNetwork: "network-from-spec",
			wantSubnet:  "subnet-from-spec",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := CloudPrivateEndpointModel{
				NetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
				SubnetId:  ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
			}
			tc.targetSpec.Name = "my-endpoint"
			response := &CloudPrivateEndpointAPIResponse{
				Id:           "pe-1",
				TargetSpec:   tc.targetSpec,
				CurrentState: tc.currentState,
			}

			model.MergeWith(context.Background(), response)

			if got := model.NetworkId.ValueString(); got != tc.wantNetwork {
				t.Fatalf("unexpected network_id: got %q, want %q", got, tc.wantNetwork)
			}
			if got := model.SubnetId.ValueString(); got != tc.wantSubnet {
				t.Fatalf("unexpected subnet_id: got %q, want %q", got, tc.wantSubnet)
			}
		})
	}
}

// TestCloudPrivateEndpointMergeWith_PlacementUnknownResolvesToNull
// verifies neither Computed placement attribute stays unknown when the response
// carries no reference at all, which would fail the apply.
func TestCloudPrivateEndpointMergeWith_PlacementUnknownResolvesToNull(t *testing.T) {
	model := CloudPrivateEndpointModel{
		NetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		SubnetId:  ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id:         "pe-1",
		TargetSpec: &CloudPrivateEndpointTargetSpec{Name: "my-endpoint"},
	}

	model.MergeWith(context.Background(), response)

	if model.NetworkId.IsUnknown() || !model.NetworkId.IsNull() {
		t.Fatalf("expected network_id to be null, got unknown=%v value=%q", model.NetworkId.IsUnknown(), model.NetworkId.ValueString())
	}
	if model.SubnetId.IsUnknown() || !model.SubnetId.IsNull() {
		t.Fatalf("expected subnet_id to be null, got unknown=%v value=%q", model.SubnetId.IsUnknown(), model.SubnetId.ValueString())
	}
}

// TestCloudPrivateEndpointMergeWith_PlacementNoTargetSpec verifies the
// fallback still runs when the envelope carries no targetSpec at all.
func TestCloudPrivateEndpointMergeWith_PlacementNoTargetSpec(t *testing.T) {
	model := CloudPrivateEndpointModel{
		NetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		SubnetId:  ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
	}
	response := &CloudPrivateEndpointAPIResponse{
		Id: "pe-1",
		CurrentState: &CloudPrivateEndpointCurrentState{
			Network: &CloudPrivateEndpointRef{Id: "network-123"},
			Subnet:  &CloudPrivateEndpointRef{Id: "subnet-456"},
		},
	}

	model.MergeWith(context.Background(), response)

	if got := model.NetworkId.ValueString(); got != "network-123" {
		t.Fatalf("unexpected network_id: got %q", got)
	}
	if got := model.SubnetId.ValueString(); got != "subnet-456" {
		t.Fatalf("unexpected subnet_id: got %q", got)
	}
}

// TestPrivateEndpointAPIDecode pins the JSON field names to the
// publicCloud.privateEndpoint and publicCloud.reference.privateEndpoint models.
func TestPrivateEndpointAPIDecode(t *testing.T) {
	var pe CloudPrivateEndpointAPIResponse
	if err := json.Unmarshal([]byte(`{
		"targetSpec": {"service": {"name": "s3", "tier": "M"}},
		"currentState": {
			"endpointId": "pe-0123abcd",
			"publicDnsEnabled": true,
			"service": {"name": "s3", "tier": "M", "port": 443, "transport": "TCP", "bandwidth": 1000}
		}
	}`), &pe); err != nil {
		t.Fatal(err)
	}
	cs := pe.CurrentState
	if pe.TargetSpec.Service.Tier != "M" || cs.Service.Tier != "M" || cs.EndpointId != "pe-0123abcd" || !cs.PublicDnsEnabled {
		t.Errorf("unexpected decode: targetSpec.service=%+v currentState=%+v", *pe.TargetSpec.Service, *cs)
	}
	if cs.Service.Bandwidth == nil || *cs.Service.Bandwidth != 1000 {
		t.Errorf("unexpected currentState.service.bandwidth: %v", cs.Service.Bandwidth)
	}

	var sa CloudPrivateEndpointServiceAvailability
	if err := json.Unmarshal([]byte(`{"name": "s3", "tiers": [{"name": "S", "location": {"region": "GRA11"}, "port": 443, "transport": "TCP", "bandwidth": null}]}`), &sa); err != nil {
		t.Fatal(err)
	}
	if len(sa.Tiers) != 1 || sa.Tiers[0].Name != "S" || sa.Tiers[0].Location.Region != "GRA11" || sa.Tiers[0].Bandwidth != nil {
		t.Errorf("unexpected service availability decode: %+v", sa)
	}
}
