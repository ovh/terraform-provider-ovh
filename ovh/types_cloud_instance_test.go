package ovh

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func strVal(s string) ovhtypes.TfStringValue {
	return ovhtypes.TfStringValue{StringValue: types.StringValue(s)}
}

func strNull() ovhtypes.TfStringValue {
	return ovhtypes.TfStringValue{StringValue: types.StringNull()}
}

func customStringList(ids ...string) ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] {
	vals := make([]attr.Value, len(ids))
	for i, id := range ids {
		vals[i] = strVal(id)
	}
	return ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
		ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals),
	}
}

func TestUnitCloudInstanceModelToCreate(t *testing.T) {
	m := &CloudInstanceModel{
		ServiceName:      strVal("proj"),
		Region:           strVal("GRA11"),
		AvailabilityZone: strNull(),
		Name:             strVal("web-1"),
		FlavorId:         strVal("flavor-uuid"),
		ImageId:          strVal("image-uuid"),
		PowerState:       strNull(),
		SSHKeyName:       strVal("mykey"),
		GroupId:          strNull(),
		Networks:         types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
		VolumeIds:        ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})},
		SecurityGroupIds: customStringList("sg-1"),
		Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
	}

	payload := m.ToCreate()

	if payload.TargetSpec.Name != "web-1" {
		t.Fatalf("name = %q, want web-1", payload.TargetSpec.Name)
	}
	if payload.TargetSpec.Flavor == nil || payload.TargetSpec.Flavor.Id != "flavor-uuid" {
		t.Fatalf("flavor.id not set correctly: %+v", payload.TargetSpec.Flavor)
	}
	if payload.TargetSpec.Image == nil || payload.TargetSpec.Image.Id != "image-uuid" {
		t.Fatalf("image.id not set correctly: %+v", payload.TargetSpec.Image)
	}
	if payload.TargetSpec.Location == nil || payload.TargetSpec.Location.Region != "GRA11" {
		t.Fatalf("location.region not set correctly: %+v", payload.TargetSpec.Location)
	}
	if payload.TargetSpec.Location.AvailabilityZone != "" {
		t.Fatalf("availabilityZone should be empty when null, got %q", payload.TargetSpec.Location.AvailabilityZone)
	}
	if payload.TargetSpec.SSHKeyName != "mykey" {
		t.Fatalf("sshKeyName = %q, want mykey", payload.TargetSpec.SSHKeyName)
	}
	if len(payload.TargetSpec.SecurityGroups) != 1 || payload.TargetSpec.SecurityGroups[0].Id != "sg-1" {
		t.Fatalf("securityGroups not set correctly: %+v", payload.TargetSpec.SecurityGroups)
	}
	// powerState omitted (null) so server defaults it
	b, _ := json.Marshal(payload)
	if got := string(b); got == "" {
		t.Fatalf("payload did not marshal")
	}
}

func TestUnitCloudInstanceModelToUpdate(t *testing.T) {
	m := &CloudInstanceModel{
		Region:           strVal("GRA11"),   // immutable — must NOT appear
		AvailabilityZone: strVal("GRA11-a"), // immutable — must NOT appear
		SSHKeyName:       strVal("mykey"),   // immutable — must NOT appear
		GroupId:          strVal("grp-1"),   // immutable — must NOT appear
		Name:             strVal("web-2"),
		FlavorId:         strVal("flavor-uuid"),
		ImageId:          strVal("image-uuid"),
		PowerState:       strVal("SHUTOFF"),
		Networks:         types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
		VolumeIds:        ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})},
		SecurityGroupIds: ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})},
		Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
	}

	payload := m.ToUpdate("chk-123")

	if payload.Checksum != "chk-123" {
		t.Fatalf("checksum = %q, want chk-123", payload.Checksum)
	}
	if payload.TargetSpec.Name != "web-2" {
		t.Fatalf("name = %q, want web-2", payload.TargetSpec.Name)
	}
	if payload.TargetSpec.PowerState != "SHUTOFF" {
		t.Fatalf("powerState = %q, want SHUTOFF", payload.TargetSpec.PowerState)
	}
	// Immutable fields must be absent from the update target spec JSON.
	b, _ := json.Marshal(payload.TargetSpec)
	for _, forbidden := range []string{"location", "sshKeyName", "group"} {
		if strings.Contains(string(b), forbidden) {
			t.Fatalf("update target spec must not contain %q: %s", forbidden, string(b))
		}
	}
}

func TestUnitCloudInstanceModelMergeWith(t *testing.T) {
	ctx := context.Background()
	resp := &CloudInstanceAPIResponse{
		Id:             "inst-1",
		Checksum:       "chk-9",
		CreatedAt:      "2026-07-09T10:00:00Z",
		UpdatedAt:      "2026-07-09T10:05:00Z",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:       "web-1",
			Flavor:     &CloudInstanceRef{Id: "flavor-uuid"},
			Image:      &CloudInstanceRef{Id: "image-uuid"},
			Location:   &CloudInstanceAPILocation{Region: "GRA11", AvailabilityZone: "GRA11-a"},
			PowerState: "ACTIVE",
			SSHKeyName: "mykey",
		},
		CurrentState: &CloudInstanceAPICurrentState{
			Name:       "web-1",
			Flavor:     &CloudInstanceAPIFlavor{Id: "flavor-uuid", Name: "b2-7", Vcpus: 2, Ram: 7000, Disk: 50},
			Image:      &CloudInstanceAPIImage{Id: "image-uuid", Name: "Debian 12", Size: 2000000000, Status: "active"},
			Location:   &CloudInstanceAPILocation{Region: "GRA11", AvailabilityZone: "GRA11-a"},
			PowerState: "ACTIVE",
			Networks: []CloudInstanceAPINetworkState{
				{Id: "net-1", Public: false, SubnetId: "sub-1", Addresses: []CloudInstanceAPIAddress{{Ip: "10.0.0.5", Mac: "fa:16:3e", Type: "private", Version: 4}}},
			},
			Volumes:   []CloudInstanceAPIVolume{{Id: "vol-1", Name: "data", Size: 100}},
			Locked:    false,
			ProjectId: "proj",
		},
	}

	m := &CloudInstanceModel{ServiceName: strVal("proj")}
	m.MergeWith(ctx, resp)

	if m.Id.ValueString() != "inst-1" {
		t.Fatalf("id = %q", m.Id.ValueString())
	}
	if m.Checksum.ValueString() != "chk-9" {
		t.Fatalf("checksum = %q", m.Checksum.ValueString())
	}
	if m.ResourceStatus.ValueString() != "READY" {
		t.Fatalf("resource_status = %q", m.ResourceStatus.ValueString())
	}
	// region/flattened fields come from TargetSpec
	if m.Region.ValueString() != "GRA11" {
		t.Fatalf("region = %q, want GRA11", m.Region.ValueString())
	}
	if m.FlavorId.ValueString() != "flavor-uuid" {
		t.Fatalf("flavor_id = %q", m.FlavorId.ValueString())
	}
	if m.ImageId.ValueString() != "image-uuid" {
		t.Fatalf("image_id = %q", m.ImageId.ValueString())
	}
	if m.CurrentState.IsNull() {
		t.Fatalf("current_state should not be null")
	}
	csAttrs := m.CurrentState.Attributes()
	flavorObj, ok := csAttrs["flavor"].(types.Object)
	if !ok || flavorObj.IsNull() {
		t.Fatalf("current_state.flavor missing")
	}
	if v, _ := flavorObj.Attributes()["vcpus"].(types.Int64); v.ValueInt64() != 2 {
		t.Fatalf("current_state.flavor.vcpus = %d, want 2", v.ValueInt64())
	}
}

func TestUnitCloudInstanceMergeWithNilCurrentState(t *testing.T) {
	ctx := context.Background()
	resp := &CloudInstanceAPIResponse{
		Id:             "inst-2",
		ResourceStatus: "CREATING",
		TargetSpec:     &CloudInstanceAPITargetSpec{Name: "n", Location: &CloudInstanceAPILocation{Region: "GRA11"}},
		CurrentState:   nil,
	}
	m := &CloudInstanceModel{}
	m.MergeWith(ctx, resp)
	if !m.CurrentState.IsNull() {
		t.Fatalf("current_state should be null when API currentState is nil")
	}
}

// Boot-from-volume: image nil on both specs must yield a null image_id + null current_state.image.
func TestUnitCloudInstanceMergeWithNilImage(t *testing.T) {
	ctx := context.Background()
	resp := &CloudInstanceAPIResponse{
		Id:             "inst-3",
		ResourceStatus: "READY",
		TargetSpec:     &CloudInstanceAPITargetSpec{Name: "n", Location: &CloudInstanceAPILocation{Region: "GRA11"}, Image: nil},
		CurrentState:   &CloudInstanceAPICurrentState{Name: "n", Image: nil, Location: &CloudInstanceAPILocation{Region: "GRA11"}},
	}
	m := &CloudInstanceModel{}
	m.MergeWith(ctx, resp)
	if !m.ImageId.IsNull() {
		t.Fatalf("image_id should be null for boot-from-volume, got %q", m.ImageId.ValueString())
	}
	imgObj, _ := m.CurrentState.Attributes()["image"].(types.Object)
	if !imgObj.IsNull() {
		t.Fatalf("current_state.image should be null for boot-from-volume")
	}
}
