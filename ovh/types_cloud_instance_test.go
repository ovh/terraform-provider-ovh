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
				{Id: "net-1", SubnetId: "sub-1", Addresses: []CloudInstanceAPIAddress{{Ip: "10.0.0.5", Mac: "fa:16:3e", Type: "FIXED", Version: 4}}},
			},
			Volumes:   []CloudInstanceAPIVolume{{Id: "vol-1", Name: "data", Size: 100}},
			Locked:    false,
			ProjectId: "proj",
		},
	}

	m := &CloudInstanceModel{ServiceName: strVal("proj")}
	m.MergeWith(ctx, resp, m.priorSpec())

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
	m.MergeWith(ctx, resp, m.priorSpec())
	if !m.CurrentState.IsNull() {
		t.Fatalf("current_state should be null when API currentState is nil")
	}
}

func instanceNetworkRefObj(networkId, subnetId, ip attr.Value, autoAssign attr.Value) attr.Value {
	obj, _ := types.ObjectValue(instanceNetworkRefAttrTypes(), map[string]attr.Value{
		"network_id":            networkId,
		"subnet_id":             subnetId,
		"ip":                    ip,
		"auto_assign_public_ip": autoAssign,
	})
	return obj
}

func TestUnitCloudInstanceNetworksRoundTrip(t *testing.T) {
	list := types.ListValueMust(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}, []attr.Value{
		// private NIC with an address (pinned fixed IP or existing floating IP)
		instanceNetworkRefObj(strVal("net-1"), strVal("sub-1"), strVal("203.0.113.10"), types.BoolNull()),
		// platform-assigned public IP
		instanceNetworkRefObj(strNull(), strNull(), strNull(), types.BoolValue(true)),
		// public IP the project already owns
		instanceNetworkRefObj(strNull(), strNull(), strVal("203.0.113.20"), types.BoolNull()),
		// private NIC, IPAM address
		instanceNetworkRefObj(strVal("net-2"), strVal("sub-2"), strNull(), types.BoolNull()),
	})

	refs := networksToAPI(list)
	if len(refs) != 4 {
		t.Fatalf("networksToAPI len = %d, want 4", len(refs))
	}
	if refs[0].Id != "net-1" || refs[0].SubnetId != "sub-1" || refs[0].IP != "203.0.113.10" || refs[0].AutoAssignPublicIp {
		t.Fatalf("private ref with ip not converted: %+v", refs[0])
	}
	if refs[1].Id != "" || refs[1].SubnetId != "" || refs[1].IP != "" || !refs[1].AutoAssignPublicIp {
		t.Fatalf("auto-assign ref not converted: %+v", refs[1])
	}
	if refs[2].Id != "" || refs[2].IP != "203.0.113.20" || refs[2].AutoAssignPublicIp {
		t.Fatalf("owned public ip ref not converted: %+v", refs[2])
	}
	if refs[3].Id != "net-2" || refs[3].SubnetId != "sub-2" || refs[3].IP != "" || refs[3].AutoAssignPublicIp {
		t.Fatalf("private IPAM ref not converted: %+v", refs[3])
	}

	b, _ := json.Marshal(refs)
	for _, forbidden := range []string{"floatingIpId", `"public"`, `"publicIp"`} {
		if strings.Contains(string(b), forbidden) {
			t.Fatalf("networks JSON must not contain %q: %s", forbidden, string(b))
		}
	}

	rebuilt := buildInstanceNetworksRootList(types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}), &CloudInstanceAPITargetSpec{Networks: refs})
	elems := rebuilt.Elements()
	if len(elems) != 4 {
		t.Fatalf("rebuilt networks len = %d, want 4", len(elems))
	}
	attrs := elems[0].(types.Object).Attributes()
	if ip, _ := attrs["ip"].(ovhtypes.TfStringValue); ip.ValueString() != "203.0.113.10" {
		t.Fatalf("rebuilt networks.0.ip = %q", ip.ValueString())
	}
	attrs = elems[1].(types.Object).Attributes()
	if auto, _ := attrs["auto_assign_public_ip"].(types.Bool); !auto.ValueBool() {
		t.Fatalf("rebuilt networks.1.auto_assign_public_ip should be true")
	}
	if ip, _ := attrs["ip"].(ovhtypes.TfStringValue); !ip.IsNull() {
		t.Fatalf("rebuilt networks.1.ip should be null, got %q", ip.ValueString())
	}
	// A false auto_assign_public_ip must round-trip as null, otherwise a config
	// that omitted the Optional-only attribute gets an inconsistent-result error.
	for _, i := range []int{0, 2, 3} {
		auto, _ := elems[i].(types.Object).Attributes()["auto_assign_public_ip"].(types.Bool)
		if !auto.IsNull() {
			t.Fatalf("rebuilt networks.%d.auto_assign_public_ip should be null, got %v", i, auto.ValueBool())
		}
	}
	if ip, _ := elems[3].(types.Object).Attributes()["ip"].(ovhtypes.TfStringValue); !ip.IsNull() {
		t.Fatalf("rebuilt networks.3.ip should be null, got %q", ip.ValueString())
	}
}

func instanceNetworksList(elems ...attr.Value) types.List {
	return types.ListValueMust(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}, elems)
}

// networkRefTuple flattens one rebuilt `networks` element for comparison.
func networkRefTuple(v attr.Value) [4]string {
	attrs := v.(types.Object).Attributes()
	str := func(k string) string {
		s, _ := attrs[k].(ovhtypes.TfStringValue)
		if s.IsNull() {
			return ""
		}
		return s.ValueString()
	}
	auto := "false"
	if b, _ := attrs["auto_assign_public_ip"].(types.Bool); b.ValueBool() {
		auto = "true"
	}
	return [4]string{str("network_id"), str("subnet_id"), str("ip"), auto}
}

// apiv2 sorts the derived target spec networks by network UUID (id-less public
// entry first). `networks` is Optional-only, so a private-first config must come
// back private-first or the apply fails with an inconsistent-result error.
func TestUnitCloudInstanceNetworksKeepPriorOrder(t *testing.T) {
	ctx := context.Background()

	private := instanceNetworkRefObj(strVal("net-1"), strVal("sub-1"), strNull(), types.BoolNull())
	public := instanceNetworkRefObj(strNull(), strNull(), strNull(), types.BoolValue(true))

	// API answers public-first (sorted), config was written private-first.
	resp := &CloudInstanceAPIResponse{
		Id:             "inst-order",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:     "n",
			Location: &CloudInstanceAPILocation{Region: "GRA11"},
			Networks: []CloudInstanceAPINetworkRef{
				{AutoAssignPublicIp: true},
				{Id: "net-1", SubnetId: "sub-1"},
			},
		},
	}

	m := &CloudInstanceModel{Networks: instanceNetworksList(private, public)}
	m.MergeWith(ctx, resp, m.priorSpec())

	got := m.Networks.Elements()
	if len(got) != 2 {
		t.Fatalf("networks len = %d, want 2", len(got))
	}
	if want := [4]string{"net-1", "sub-1", "", "false"}; networkRefTuple(got[0]) != want {
		t.Fatalf("networks[0] = %v, want %v", networkRefTuple(got[0]), want)
	}
	if want := [4]string{"", "", "", "true"}; networkRefTuple(got[1]) != want {
		t.Fatalf("networks[1] = %v, want %v", networkRefTuple(got[1]), want)
	}

	// An entry that correlates to nothing in the prior list is real drift and is
	// appended after the correlated ones.
	resp.TargetSpec.Networks = append(resp.TargetSpec.Networks, CloudInstanceAPINetworkRef{Id: "net-9", SubnetId: "sub-9"})
	m = &CloudInstanceModel{Networks: instanceNetworksList(private, public)}
	m.MergeWith(ctx, resp, m.priorSpec())
	got = m.Networks.Elements()
	if len(got) != 3 {
		t.Fatalf("networks len = %d, want 3", len(got))
	}
	if want := [4]string{"net-9", "sub-9", "", "false"}; networkRefTuple(got[2]) != want {
		t.Fatalf("networks[2] = %v, want %v", networkRefTuple(got[2]), want)
	}
}

// terraform import has no prior list: the API order must be kept as-is.
func TestUnitCloudInstanceNetworksNullPriorKeepsAPIOrder(t *testing.T) {
	ctx := context.Background()

	resp := &CloudInstanceAPIResponse{
		Id:             "inst-import",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:     "n",
			Location: &CloudInstanceAPILocation{Region: "GRA11"},
			Networks: []CloudInstanceAPINetworkRef{
				{AutoAssignPublicIp: true},
				{Id: "net-1", SubnetId: "sub-1"},
			},
		},
	}

	m := &CloudInstanceModel{}
	m.MergeWith(ctx, resp, m.priorSpec())

	got := m.Networks.Elements()
	if len(got) != 2 {
		t.Fatalf("networks len = %d, want 2", len(got))
	}
	if want := [4]string{"", "", "", "true"}; networkRefTuple(got[0]) != want {
		t.Fatalf("networks[0] = %v, want %v", networkRefTuple(got[0]), want)
	}
	if want := [4]string{"net-1", "sub-1", "", "false"}; networkRefTuple(got[1]) != want {
		t.Fatalf("networks[1] = %v, want %v", networkRefTuple(got[1]), want)
	}
}

// shares[].access_level is Optional-only: apiv2 normalises an omitted level to
// READ_WRITE and echoes it back, so the merged value must stay null when the
// config omitted it and carry the API value otherwise.
func TestUnitCloudInstanceSharesAccessLevelRoundTrip(t *testing.T) {
	ctx := context.Background()

	shareObj := func(id string, level attr.Value) attr.Value {
		obj, _ := types.ObjectValue(instanceShareRefAttrTypes(), map[string]attr.Value{
			"id":           strVal(id),
			"access_level": level,
		})
		return obj
	}
	prior := types.ListValueMust(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}, []attr.Value{
		shareObj("sh-omitted", strNull()),
		shareObj("sh-explicit", strVal("READ_ONLY")),
		shareObj("sh-explicit-default", strVal("READ_WRITE")),
	})

	resp := &CloudInstanceAPIResponse{
		Id:             "inst-shares",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:     "n",
			Location: &CloudInstanceAPILocation{Region: "GRA11"},
			Shares: []CloudInstanceAPIShareRef{
				{Id: "sh-omitted", AccessLevel: "READ_WRITE"},
				{Id: "sh-explicit", AccessLevel: "READ_ONLY"},
				{Id: "sh-explicit-default", AccessLevel: "READ_WRITE"},
			},
		},
	}

	m := &CloudInstanceModel{Shares: prior}
	m.MergeWith(ctx, resp, m.priorSpec())

	want := []string{"", "READ_ONLY", "READ_WRITE"}
	elems := m.Shares.Elements()
	if len(elems) != len(want) {
		t.Fatalf("shares len = %d, want %d", len(elems), len(want))
	}
	for i, expected := range want {
		level, _ := elems[i].(types.Object).Attributes()["access_level"].(ovhtypes.TfStringValue)
		got := ""
		if !level.IsNull() {
			got = level.ValueString()
		}
		if got != expected {
			t.Fatalf("shares.%d.access_level = %q, want %q", i, got, expected)
		}
	}

	// On import there is no prior list, so the API value is kept verbatim.
	m = &CloudInstanceModel{}
	m.MergeWith(ctx, resp, m.priorSpec())
	level, _ := m.Shares.Elements()[0].(types.Object).Attributes()["access_level"].(ovhtypes.TfStringValue)
	if level.IsNull() || level.ValueString() != "READ_WRITE" {
		t.Fatalf("imported shares.0.access_level = %+v, want READ_WRITE", level)
	}
}

// A response without a target spec must still leave every list attribute with a
// typed null; a zero-value list has a nil element type and panics at State.Set.
func TestUnitCloudInstanceMergeWithNoTargetSpecTypedNulls(t *testing.T) {
	ctx := context.Background()

	m := &CloudInstanceModel{}
	m.MergeWith(ctx, &CloudInstanceAPIResponse{Id: "inst-6", ResourceStatus: "CREATING"}, m.priorSpec())

	if m.Networks.ElementType(ctx) == nil {
		t.Fatal("networks must carry an element type")
	}
	if m.Shares.ElementType(ctx) == nil {
		t.Fatal("shares must carry an element type")
	}
	if m.VolumeIds.ElementType(ctx) == nil {
		t.Fatal("volume_ids must carry an element type")
	}
	if m.SecurityGroupIds.ElementType(ctx) == nil {
		t.Fatal("security_group_ids must carry an element type")
	}
}

func emptyCustomStringList() ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] {
	return ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
		ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, []attr.Value{}),
	}
}

func nullCustomStringList() ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] {
	return ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
		ListValue: basetypes.NewListNull(ovhtypes.TfStringType{}),
	}
}

// The API distinguishes a null securityGroups (apply the project default on
// create, leave unchanged on update) from an explicit [] (no security group), so
// both must survive JSON marshalling instead of being dropped by omitempty.
func TestUnitCloudInstanceSecurityGroupsNullVersusEmpty(t *testing.T) {
	base := CloudInstanceModel{
		Region:    strVal("GRA11"),
		Name:      strVal("web-1"),
		FlavorId:  strVal("flavor-uuid"),
		Networks:  types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
		VolumeIds: nullCustomStringList(),
		Shares:    types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
	}

	cases := []struct {
		name string
		sgs  ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
		want string
	}{
		{name: "null", sgs: nullCustomStringList(), want: `"securityGroups":null`},
		{name: "empty", sgs: emptyCustomStringList(), want: `"securityGroups":[]`},
		{name: "one", sgs: customStringList("sg-1"), want: `"securityGroups":[{"id":"sg-1"}]`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := base
			m.SecurityGroupIds = tc.sgs

			b, err := json.Marshal(m.ToCreate().TargetSpec)
			if err != nil {
				t.Fatalf("marshal create target spec: %s", err)
			}
			if !strings.Contains(string(b), tc.want) {
				t.Fatalf("create target spec must contain %s: %s", tc.want, string(b))
			}

			b, err = json.Marshal(m.ToUpdate("chk-1").TargetSpec)
			if err != nil {
				t.Fatalf("marshal update target spec: %s", err)
			}
			if !strings.Contains(string(b), tc.want) {
				t.Fatalf("update target spec must contain %s: %s", tc.want, string(b))
			}
		})
	}
}

// security_group_ids is Optional+Computed: it must always come back known, and
// "no security group" round-trips as an empty list rather than null so that
// refresh and import stay stable whether the API echoes [] or omits the field.
func TestUnitCloudInstanceMergeWithSecurityGroups(t *testing.T) {
	ctx := context.Background()

	newResponse := func(sgs []CloudInstanceRef) *CloudInstanceAPIResponse {
		return &CloudInstanceAPIResponse{
			Id:             "inst-4",
			ResourceStatus: "READY",
			TargetSpec: &CloudInstanceAPITargetSpec{
				Name:           "n",
				Location:       &CloudInstanceAPILocation{Region: "GRA11"},
				SecurityGroups: sgs,
			},
		}
	}

	// Omitted in config (unknown in the plan) + API-assigned default group.
	m := &CloudInstanceModel{SecurityGroupIds: ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListUnknown(ovhtypes.TfStringType{})}}
	m.MergeWith(ctx, newResponse([]CloudInstanceRef{{Id: "sg-default"}}), m.priorSpec())
	if m.SecurityGroupIds.IsUnknown() || m.SecurityGroupIds.IsNull() {
		t.Fatalf("security_group_ids must be known after merge: %+v", m.SecurityGroupIds)
	}
	if got := m.SecurityGroupIds.Elements(); len(got) != 1 || got[0].(ovhtypes.TfStringValue).ValueString() != "sg-default" {
		t.Fatalf("security_group_ids = %+v, want [sg-default]", got)
	}

	// Explicit empty list echoed back as an empty array.
	m = &CloudInstanceModel{SecurityGroupIds: emptyCustomStringList()}
	m.MergeWith(ctx, newResponse([]CloudInstanceRef{}), m.priorSpec())
	if m.SecurityGroupIds.IsNull() || len(m.SecurityGroupIds.Elements()) != 0 {
		t.Fatalf("explicit empty security_group_ids must stay an empty list: %+v", m.SecurityGroupIds)
	}

	// Explicit empty list, but the API omitted the field entirely.
	m = &CloudInstanceModel{SecurityGroupIds: emptyCustomStringList()}
	m.MergeWith(ctx, newResponse(nil), m.priorSpec())
	if m.SecurityGroupIds.IsNull() || len(m.SecurityGroupIds.Elements()) != 0 {
		t.Fatalf("omitted securityGroups must map to an empty list: %+v", m.SecurityGroupIds)
	}

	// No target spec at all (early state save right after POST).
	m = &CloudInstanceModel{SecurityGroupIds: ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListUnknown(ovhtypes.TfStringType{})}}
	m.MergeWith(ctx, &CloudInstanceAPIResponse{Id: "inst-5", ResourceStatus: "CREATING"}, m.priorSpec())
	if m.SecurityGroupIds.IsUnknown() {
		t.Fatal("security_group_ids must never stay unknown after merge")
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
	m.MergeWith(ctx, resp, m.priorSpec())
	if !m.ImageId.IsNull() {
		t.Fatalf("image_id should be null for boot-from-volume, got %q", m.ImageId.ValueString())
	}
	imgObj, _ := m.CurrentState.Attributes()["image"].(types.Object)
	if !imgObj.IsNull() {
		t.Fatalf("current_state.image should be null for boot-from-volume")
	}
}
