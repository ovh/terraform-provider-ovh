package ovh

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
	m := &CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
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
			Volumes:          types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
			SecurityGroupIds: customStringList("sg-1"),
			Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
		},
		UserData: strVal("I2Nsb3VkLWNvbmZpZw=="),
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
	if payload.TargetSpec.UserData != "I2Nsb3VkLWNvbmZpZw==" {
		t.Fatalf("userData = %q, want I2Nsb3VkLWNvbmZpZw==", payload.TargetSpec.UserData)
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
	m := &CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
			Region:           strVal("GRA11"),   // immutable — must NOT appear
			AvailabilityZone: strVal("GRA11-a"), // immutable — must NOT appear
			SSHKeyName:       strVal("mykey"),   // immutable — must NOT appear
			GroupId:          strVal("grp-1"),   // immutable — must NOT appear
			Name:             strVal("web-2"),
			FlavorId:         strVal("flavor-uuid"),
			ImageId:          strVal("image-uuid"),
			PowerState:       strVal("SHUTOFF"),
			Networks:         types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
			Volumes:          types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
			SecurityGroupIds: ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})},
			Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
		},
		UserData: strNull(),
	}

	payload := m.ToUpdate("chk-123", strNull())

	if payload.Checksum != "chk-123" {
		t.Fatalf("checksum = %q, want chk-123", payload.Checksum)
	}
	if payload.TargetSpec.Name != "web-2" {
		t.Fatalf("name = %q, want web-2", payload.TargetSpec.Name)
	}
	if payload.TargetSpec.PowerState == nil || *payload.TargetSpec.PowerState != "SHUTOFF" {
		t.Fatalf("powerState = %v, want SHUTOFF", payload.TargetSpec.PowerState)
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
	if m.Volumes.ElementType(ctx) == nil {
		t.Fatal("volumes must carry an element type")
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
	base := CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
			Region:   strVal("GRA11"),
			Name:     strVal("web-1"),
			FlavorId: strVal("flavor-uuid"),
			Networks: types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
			Volumes:  types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
			Shares:   types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
		},
		UserData: strNull(),
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

			b, err = json.Marshal(m.ToUpdate("chk-1", strNull()).TargetSpec)
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

// A response with no target spec says nothing about the security groups. Letting
// it null a known list would put that null back on the wire on the next PUT,
// which the API reads as "remove every group" — a running instance stranded on a
// deny-all port.
func TestUnitCloudInstanceMergeWithoutTargetSpecKeepsKnownSecurityGroups(t *testing.T) {
	ctx := context.Background()
	noTargetSpec := func() *CloudInstanceAPIResponse {
		return &CloudInstanceAPIResponse{Id: "inst-6", ResourceStatus: "UPDATING"}
	}

	m := &CloudInstanceModel{SecurityGroupIds: customStringList("sg-1", "sg-2")}
	m.MergeWith(ctx, noTargetSpec(), m.priorSpec())
	got := m.SecurityGroupIds.Elements()
	if len(got) != 2 || got[0].(ovhtypes.TfStringValue).ValueString() != "sg-1" || got[1].(ovhtypes.TfStringValue).ValueString() != "sg-2" {
		t.Fatalf("security_group_ids = %+v, want the prior [sg-1 sg-2] preserved", m.SecurityGroupIds)
	}

	// An explicit empty list is an intent of its own and survives too.
	m = &CloudInstanceModel{SecurityGroupIds: emptyCustomStringList()}
	m.MergeWith(ctx, noTargetSpec(), m.priorSpec())
	if m.SecurityGroupIds.IsNull() || len(m.SecurityGroupIds.Elements()) != 0 {
		t.Fatalf("explicit empty security_group_ids must survive: %+v", m.SecurityGroupIds)
	}

	// Nothing known to preserve: no value is manufactured, and the list keeps a
	// usable element type so resp.State.Set cannot panic on it.
	m = &CloudInstanceModel{SecurityGroupIds: nullCustomStringList()}
	m.MergeWith(ctx, noTargetSpec(), m.priorSpec())
	if !m.SecurityGroupIds.IsNull() {
		t.Fatalf("a null prior must stay null: %+v", m.SecurityGroupIds)
	}
	if m.SecurityGroupIds.ElementType(ctx) == nil {
		t.Fatal("security_group_ids must keep a typed element type")
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

const (
	testUserDataA = "I2Nsb3VkLWNvbmZpZwpwYWNrYWdlczoKICAtIG5naW54Cg=="
	testUserDataB = "I2Nsb3VkLWNvbmZpZwpwYWNrYWdlczoKICAtIGh0b3AK"
)

// userDataUpdateModel builds a minimal resource model whose only interesting
// field is user_data, for the tri-state assertions below.
func userDataUpdateModel(planned ovhtypes.TfStringValue) *CloudInstanceResourceModel {
	return &CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
			Name:             strVal("web-1"),
			FlavorId:         strVal("flavor-uuid"),
			Networks:         types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
			Volumes:          types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
			SecurityGroupIds: nullCustomStringList(),
			Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
		},
		UserData: planned,
	}
}

// userData is the documented exception to the real PUT: it is write-only — GET
// returns only its checksum — so an absent key still keeps the stored payload,
// "" clears it, and a non-empty value replaces it. Both a clear and a replace
// REINSTALL the instance (Nova rebuild, root disk wiped), so an attribute the
// user never set, or left untouched, MUST marshal to no key at all: emitting ""
// would wipe the root disk on every unrelated update.
func TestUnitCloudInstanceUserDataUpdateTriState(t *testing.T) {
	cases := []struct {
		name  string
		prior ovhtypes.TfStringValue
		plan  ovhtypes.TfStringValue
		want  string // "" means the key must be absent entirely
	}{
		{name: "never configured", prior: strNull(), plan: strNull(), want: ""},
		{name: "unchanged", prior: strVal(testUserDataA), plan: strVal(testUserDataA), want: ""},
		{name: "removed from config", prior: strVal(testUserDataA), plan: strNull(), want: ""},
		{
			name:  "unknown plan",
			prior: strVal(testUserDataA),
			plan:  ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
			want:  "",
		},
		{name: "first set", prior: strNull(), plan: strVal(testUserDataA), want: `"userData":"` + testUserDataA + `"`},
		{name: "changed", prior: strVal(testUserDataA), plan: strVal(testUserDataB), want: `"userData":"` + testUserDataB + `"`},
		{name: "cleared", prior: strVal(testUserDataA), plan: strVal(""), want: `"userData":""`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(userDataUpdateModel(tc.plan).ToUpdate("chk-1", tc.prior).TargetSpec)
			if err != nil {
				t.Fatalf("marshal update target spec: %s", err)
			}
			got := string(b)

			if tc.want == "" {
				if strings.Contains(got, "userData") {
					t.Fatalf("update target spec must carry NO userData key — an empty one clears the user data and wipes the root disk: %s", got)
				}
				return
			}
			if !strings.Contains(got, tc.want) {
				t.Fatalf("update target spec must contain %s: %s", tc.want, got)
			}
		})
	}
}

// A real PUT clears whatever the body omits, so every other field of the API's
// update model is always present — even unset — and the three fields it does not
// accept never are.
func TestUnitCloudInstanceUpdateBodyCarriesEveryMutableField(t *testing.T) {
	networks := types.ListValueMust(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}, []attr.Value{
		types.ObjectValueMust(instanceNetworkRefAttrTypes(), map[string]attr.Value{
			"network_id":            strVal("net-1"),
			"subnet_id":             strVal("sub-1"),
			"ip":                    strNull(),
			"auto_assign_public_ip": types.BoolNull(),
		}),
	})
	shares := types.ListValueMust(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}, []attr.Value{
		types.ObjectValueMust(instanceShareRefAttrTypes(), map[string]attr.Value{
			"id":           strVal("share-1"),
			"access_level": strVal("READ_ONLY"),
		}),
	})

	populated := CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
			Region:           strVal("GRA11"),
			AvailabilityZone: strVal("GRA11-a"),
			SSHKeyName:       strVal("mykey"),
			GroupId:          strVal("grp-1"),
			Name:             strVal("web-1"),
			FlavorId:         strVal("flavor-uuid"),
			ImageId:          strVal("image-uuid"),
			PowerState:       strVal("ACTIVE"),
			Networks:         networks,
			Volumes:          volumeRefList(map[string]string{"vol-1": ""}),
			SecurityGroupIds: customStringList("sg-1"),
			Shares:           shares,
		},
	}

	// Same instance with every optional attribute unset: the keys must still be
	// there, stating emptiness, instead of vanishing from the body.
	bare := CloudInstanceResourceModel{
		CloudInstanceModel: CloudInstanceModel{
			Region:           strVal("GRA11"),
			SSHKeyName:       strVal("mykey"),
			GroupId:          strVal("grp-1"),
			Name:             strVal("web-1"),
			FlavorId:         strVal("flavor-uuid"),
			ImageId:          strNull(),
			PowerState:       strNull(),
			Networks:         types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}),
			Volumes:          types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
			SecurityGroupIds: emptyCustomStringList(),
			Shares:           types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}),
		},
	}

	cases := []struct {
		name  string
		model CloudInstanceResourceModel
		want  []string
	}{
		{
			name:  "populated",
			model: populated,
			want: []string{
				`"name":"web-1"`,
				`"flavor":{"id":"flavor-uuid"}`,
				`"image":{"id":"image-uuid"}`,
				`"networks":[{"id":"net-1","subnetId":"sub-1"}]`,
				`"volumes":[{"id":"vol-1","role":"DATA"}]`,
				`"shares":[{"id":"share-1","accessLevel":"READ_ONLY"}]`,
				`"powerState":"ACTIVE"`,
				`"securityGroups":[{"id":"sg-1"}]`,
			},
		},
		{
			name:  "everything unset",
			model: bare,
			want: []string{
				`"name":"web-1"`,
				`"flavor":{"id":"flavor-uuid"}`,
				`"image":null`,
				`"networks":[]`,
				`"volumes":[]`,
				`"shares":[]`,
				`"powerState":null`,
				`"securityGroups":[]`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.model
			b, err := json.Marshal(m.ToUpdate("chk-1", strNull()).TargetSpec)
			if err != nil {
				t.Fatalf("marshal update target spec: %s", err)
			}
			got := string(b)
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Fatalf("update target spec must contain %s: %s", want, got)
				}
			}
			// The API's update model rejects these as unknown fields; the server
			// carries the stored ones over instead.
			for _, forbidden := range []string{"location", "region", "availabilityZone", "sshKeyName", "group"} {
				if strings.Contains(got, forbidden) {
					t.Fatalf("update target spec must not contain %q: %s", forbidden, got)
				}
			}
		})
	}
}

// Same guarantee on create: an unset user_data sends no userData at all.
func TestUnitCloudInstanceUserDataOnCreate(t *testing.T) {
	m := userDataUpdateModel(strNull())
	b, err := json.Marshal(m.ToCreate().TargetSpec)
	if err != nil {
		t.Fatalf("marshal create target spec: %s", err)
	}
	if strings.Contains(string(b), "userData") {
		t.Fatalf("create target spec must carry no userData key when unset: %s", string(b))
	}

	m = userDataUpdateModel(strVal(testUserDataA))
	if b, err = json.Marshal(m.ToCreate().TargetSpec); err != nil {
		t.Fatalf("marshal create target spec: %s", err)
	}
	if want := `"userData":"` + testUserDataA + `"`; !strings.Contains(string(b), want) {
		t.Fatalf("create target spec must contain %s: %s", want, string(b))
	}
}

// user_data is write-only: GET and LIST never return it, so a merge must leave the
// configured value alone. Terraform state holds the only copy.
func TestUnitCloudInstanceUserDataSurvivesMerge(t *testing.T) {
	ctx := context.Background()

	m := userDataUpdateModel(strVal(testUserDataA))
	m.MergeWith(ctx, &CloudInstanceAPIResponse{
		Id:             "inst-ud",
		ResourceStatus: "READY",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:     "web-1",
			Location: &CloudInstanceAPILocation{Region: "GRA11"},
			// Should never happen (write-only), but even then the merge must not
			// let a response value displace the configured one.
			UserData: testUserDataB,
		},
	}, m.priorSpec())

	if m.UserData.IsNull() || m.UserData.IsUnknown() {
		t.Fatalf("user_data must survive a merge, got %+v", m.UserData)
	}
	if got := m.UserData.ValueString(); got != testUserDataA {
		t.Fatalf("user_data = %q, want the configured %q preserved", got, testUserDataA)
	}
}

// tfsdkTags collects a model's tfsdk tag names, following embedded structs the way
// terraform-plugin-framework's reflection does.
func tfsdkTags(t *testing.T, typ reflect.Type) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Anonymous {
			if _, ok := f.Tag.Lookup("tfsdk"); ok {
				t.Fatalf("embedded field %s must not carry a tfsdk tag, the framework rejects it", f.Name)
			}
			for k := range tfsdkTags(t, f.Type) {
				out[k] = true
			}
			continue
		}
		tag, ok := f.Tag.Lookup("tfsdk")
		if !ok {
			t.Fatalf("field %s has no tfsdk tag", f.Name)
		}
		if tag == "-" {
			continue
		}
		out[tag] = true
	}
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// The framework demands exact parity between a model's tfsdk tags and the schema
// it is read into and written from, in BOTH directions: an extra struct field or
// an extra schema attribute fails at Get/Set time, at runtime, with no compile
// error to warn you. The resource and the ovh_cloud_instance data source share
// CloudInstanceModel, so an attribute added to the shared model silently breaks
// the data source — which is exactly why write-only user_data lives on
// CloudInstanceResourceModel instead.
func TestUnitCloudInstanceModelsMatchSchemas(t *testing.T) {
	ctx := context.Background()

	schemaResp := &fwresource.SchemaResponse{}
	(&cloudInstanceResource{}).Schema(ctx, fwresource.SchemaRequest{}, schemaResp)

	resourceAttrs := map[string]bool{}
	for name := range schemaResp.Schema.Attributes {
		resourceAttrs[name] = true
	}
	resourceFields := tfsdkTags(t, reflect.TypeOf(CloudInstanceResourceModel{}))
	if !reflect.DeepEqual(sortedKeys(resourceAttrs), sortedKeys(resourceFields)) {
		t.Fatalf("resource schema and CloudInstanceResourceModel disagree:\n schema = %v\n  model = %v",
			sortedKeys(resourceAttrs), sortedKeys(resourceFields))
	}

	dsAttrs := map[string]bool{}
	for name := range instanceDataSourceAttributes(ctx) {
		dsAttrs[name] = true
	}
	dsFields := tfsdkTags(t, reflect.TypeOf(CloudInstanceModel{}))
	if !reflect.DeepEqual(sortedKeys(dsAttrs), sortedKeys(dsFields)) {
		t.Fatalf("data source schema and CloudInstanceModel disagree:\n schema = %v\n  model = %v",
			sortedKeys(dsAttrs), sortedKeys(dsFields))
	}

	if !resourceFields["user_data"] {
		t.Fatal("the resource must expose user_data")
	}
	if dsFields["user_data"] || dsAttrs["user_data"] {
		t.Fatal("the data source must NOT expose user_data: the API never returns it, so there is nothing to read back")
	}
}

// volumeRefList builds a root volumes list from id -> role, in id order. An
// empty role stands for an omitted one.
func volumeRefList(roles map[string]string) types.List {
	ids := make([]string, 0, len(roles))
	for id := range roles {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	items := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		role := strNull()
		if roles[id] != "" {
			role = strVal(roles[id])
		}
		obj, _ := types.ObjectValue(instanceVolumeRefAttrTypes(), map[string]attr.Value{"id": strVal(id), "role": role})
		items = append(items, obj)
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}, items)
}

func TestUnitCloudInstanceVolumesToAPI(t *testing.T) {
	got := volumesToAPI(volumeRefList(map[string]string{"vol-a": "BOOT", "vol-b": ""}))
	want := []CloudInstanceAPIVolumeRef{{Id: "vol-a", Role: "BOOT"}, {Id: "vol-b", Role: "DATA"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("volumesToAPI = %+v, want %+v", got, want)
	}
	if ids := bootVolumeIds(volumeRefList(map[string]string{"vol-a": "BOOT", "vol-b": "DATA"})); !reflect.DeepEqual(ids, []string{"vol-a"}) {
		t.Fatalf("bootVolumeIds = %v, want [vol-a]", ids)
	}
	if volumesToAPI(types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()})) != nil {
		t.Fatal("a null list must give no volumes")
	}
}

// A boot-from-volume instance: the API echoes the BOOT volume, image is null,
// and the role lands in both targetSpec-shaped volumes and current_state.volumes.
func TestUnitCloudInstanceMergeWithVolumeRole(t *testing.T) {
	ctx := context.Background()
	m := &CloudInstanceModel{}
	resp := &CloudInstanceAPIResponse{
		Id: "inst-boot",
		TargetSpec: &CloudInstanceAPITargetSpec{
			Name:    "boot",
			Flavor:  &CloudInstanceRef{Id: "f"},
			Volumes: []CloudInstanceAPIVolumeRef{{Id: "vol-boot", Role: "BOOT"}, {Id: "vol-data"}},
		},
		CurrentState: &CloudInstanceAPICurrentState{
			Volumes: []CloudInstanceAPIVolume{{Id: "vol-boot", Role: "BOOT"}, {Id: "vol-data"}},
		},
	}
	m.MergeWith(ctx, resp, m.priorSpec())

	if got := volumesToAPI(m.Volumes); !reflect.DeepEqual(got, []CloudInstanceAPIVolumeRef{{Id: "vol-boot", Role: "BOOT"}, {Id: "vol-data", Role: "DATA"}}) {
		t.Fatalf("volumes = %+v", got)
	}
	if !m.ImageId.IsNull() {
		t.Fatalf("image_id = %q, want null", m.ImageId.ValueString())
	}
	var roles []string
	for _, e := range m.CurrentState.Attributes()["volumes"].(types.List).Elements() {
		roles = append(roles, e.(types.Object).Attributes()["role"].(ovhtypes.TfStringValue).ValueString())
	}
	if !reflect.DeepEqual(roles, []string{"BOOT", "DATA"}) {
		t.Fatalf("current_state.volumes roles = %v, want [BOOT DATA]", roles)
	}
}

// volumeCfg is the volume-related part of a configuration.
type volumeCfg struct {
	image   ovhtypes.TfStringValue
	ids     ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
	volumes types.List
}

func newVolumeCfg() volumeCfg {
	return volumeCfg{
		image:   strNull(),
		ids:     nullCustomStringList(),
		volumes: types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()}),
	}
}

func TestUnitCloudInstanceValidateConfigVolumes(t *testing.T) {
	ctx := context.Background()
	r := &cloudInstanceResource{}
	var sch fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &sch)

	with := func(f func(c *volumeCfg)) volumeCfg {
		c := newVolumeCfg()
		f(&c)
		return c
	}
	cases := []struct {
		name    string
		cfg     volumeCfg
		wantErr bool
	}{
		{"nothing", with(func(c *volumeCfg) { c.image = strVal("img") }), false},
		{"volumes alone", with(func(c *volumeCfg) { c.volumes = volumeRefList(map[string]string{"v1": "BOOT", "v2": ""}) }), false},
		{"volume_ids alone (deprecated)", with(func(c *volumeCfg) { c.ids = customStringList("v1") }), false},
		{"volume_ids with an image", with(func(c *volumeCfg) { c.image = strVal("img"); c.ids = customStringList("d1") }), false},
		{"volumes with volume_ids", with(func(c *volumeCfg) {
			c.volumes = volumeRefList(map[string]string{"v1": ""})
			c.ids = customStringList("d1")
		}), true},
		{"BOOT volume with an image", with(func(c *volumeCfg) {
			c.image = strVal("img")
			c.volumes = volumeRefList(map[string]string{"v1": "BOOT"})
		}), true},
		{"two BOOT volumes", with(func(c *volumeCfg) { c.volumes = volumeRefList(map[string]string{"v1": "BOOT", "v2": "BOOT"}) }), true},
		{"DATA volumes with an image", with(func(c *volumeCfg) {
			c.image = strVal("img")
			c.volumes = volumeRefList(map[string]string{"v1": "DATA"})
		}), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tfsdk.Config{Schema: sch.Schema, Raw: instanceConfigRaw(t, sch.Schema, tc.cfg)}
			var resp fwresource.ValidateConfigResponse
			r.ValidateConfig(ctx, fwresource.ValidateConfigRequest{Config: cfg}, &resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Fatalf("error = %v, want %v: %v", resp.Diagnostics.HasError(), tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

// instanceConfigRaw builds a raw config with only the volume attributes and image_id set.
func instanceConfigRaw(t *testing.T, sch schema.Schema, c volumeCfg) tftypes.Value {
	t.Helper()
	st := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(sch.Type().TerraformType(context.Background()), nil)}
	ctx := context.Background()
	for name, v := range map[string]any{"image_id": c.image, "volume_ids": c.ids, "volumes": c.volumes} {
		if d := st.SetAttribute(ctx, path.Root(name), v); d.HasError() {
			t.Fatal(d)
		}
	}
	return st.Raw
}

// The payload follows whichever attribute the configuration uses.
func TestUnitCloudInstanceVolumesPayload(t *testing.T) {
	ref := func(id, role string) CloudInstanceAPIVolumeRef { return CloudInstanceAPIVolumeRef{Id: id, Role: role} }
	nullVols := types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()})
	cases := []struct {
		name   string
		model  CloudInstanceModel
		create bool
		want   []CloudInstanceAPIVolumeRef
	}{
		{"volumes", CloudInstanceModel{Volumes: volumeRefList(map[string]string{"a": "BOOT", "b": ""})}, true,
			[]CloudInstanceAPIVolumeRef{ref("a", "BOOT"), ref("b", "DATA")}},
		// Deprecated volume_ids: with no image the first volume was the boot
		// volume, and the API refuses an instance with no boot source.
		{"volume_ids, no image, create", CloudInstanceModel{VolumeIds: customStringList("v1", "v2")}, true,
			[]CloudInstanceAPIVolumeRef{ref("v1", "BOOT"), ref("v2", "")}},
		{"volume_ids with an image, create", CloudInstanceModel{ImageId: strVal("img"), VolumeIds: customStringList("v1")}, true,
			[]CloudInstanceAPIVolumeRef{ref("v1", "")}},
		// On update the role is left out: the API keeps the stored boot volume.
		{"volume_ids, no image, update", CloudInstanceModel{VolumeIds: customStringList("v1", "v2")}, false,
			[]CloudInstanceAPIVolumeRef{ref("v1", ""), ref("v2", "")}},
		{"nothing", CloudInstanceModel{}, true, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.model
			if m.Volumes.ElementType(context.Background()) == nil {
				m.Volumes = nullVols
			}
			if got := m.volumesPayload(tc.create); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("payload = %+v, want %+v", got, tc.want)
			}
		})
	}
}

// Volumes and volume_ids are Optional+Computed views of targetSpec.volumes, so
// reading back always fills both, whichever one the configuration uses.
func TestUnitCloudInstanceVolumeFormsRoundTrip(t *testing.T) {
	ctx := context.Background()
	resp := func() *CloudInstanceAPIResponse {
		return &CloudInstanceAPIResponse{
			Id: "i",
			TargetSpec: &CloudInstanceAPITargetSpec{
				Name: "n", Flavor: &CloudInstanceRef{Id: "f"},
				Volumes: []CloudInstanceAPIVolumeRef{{Id: "boot", Role: "BOOT"}, {Id: "d1"}, {Id: "d2", Role: "DATA"}},
			},
		}
	}
	nullVols := types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()})
	ids := func(l ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]) []string {
		var out []string
		for _, e := range l.Elements() {
			out = append(out, e.(ovhtypes.TfStringValue).ValueString())
		}
		return out
	}
	want := []string{"boot", "d1", "d2"}
	models := map[string]*CloudInstanceModel{
		"config uses volumes":    {Volumes: volumeRefList(map[string]string{"boot": "BOOT"}), VolumeIds: nullCustomStringList()},
		"config uses volume_ids": {VolumeIds: customStringList("boot", "d1", "d2"), Volumes: nullVols},
		"import or data source":  {Volumes: nullVols, VolumeIds: nullCustomStringList()},
	}
	for name, m := range models {
		t.Run(name, func(t *testing.T) {
			m.MergeWith(ctx, resp(), m.priorSpec())
			if len(m.Volumes.Elements()) != 3 || !reflect.DeepEqual(ids(m.VolumeIds), want) {
				t.Fatalf("both must be filled in the API order: volumes=%d volume_ids=%v", len(m.Volumes.Elements()), ids(m.VolumeIds))
			}
		})
	}
	t.Run("no volume at all stays null", func(t *testing.T) {
		m := &CloudInstanceModel{Volumes: nullVols, VolumeIds: nullCustomStringList()}
		r := resp()
		r.TargetSpec.Volumes = nil
		m.MergeWith(ctx, r, m.priorSpec())
		if !m.Volumes.IsNull() || !m.VolumeIds.IsNull() {
			t.Fatal("both must stay null")
		}
	})
}

// The boot volume reads the same from a plan and from a state, whichever
// attribute declares the volumes. Moving from volume_ids to volumes keeps it, so
// it must not replace the instance; changing it does.
func TestUnitCloudInstanceBootVolumeId(t *testing.T) {
	nullVols := types.ListNull(types.ObjectType{AttrTypes: instanceVolumeRefAttrTypes()})
	nullIds := nullCustomStringList()
	cases := []struct {
		name    string
		volumes types.List
		ids     ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
		image   ovhtypes.TfStringValue
		want    string
	}{
		{"volumes, BOOT entry", volumeRefList(map[string]string{"b": "BOOT", "d": ""}), nullIds, strNull(), "b"},
		{"volumes, no BOOT", volumeRefList(map[string]string{"d": ""}), nullIds, strNull(), ""},
		{"volume_ids, no image: first volume", nullVols, customStringList("b", "d"), strNull(), "b"},
		{"volume_ids with an image: none", nullVols, customStringList("b", "d"), strVal("img"), ""},
		{"nothing", nullVols, nullIds, strNull(), ""},
		// volumes decides, even next to an image_id read from a state.
		{"volumes win over image_id", volumeRefList(map[string]string{"b": "BOOT"}), nullIds, strVal("img"), "b"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := bootVolumeId(tc.volumes, tc.ids, tc.image); got != tc.want {
				t.Fatalf("bootVolumeId = %q, want %q", got, tc.want)
			}
		})
	}

	// volume_ids = ["b", "d"] -> volumes = [{b BOOT}, {d}] keeps b: no replacement.
	from := bootVolumeId(nullVols, customStringList("b", "d"), strNull())
	to := bootVolumeId(volumeRefList(map[string]string{"b": "BOOT", "d": ""}), nullIds, strNull())
	if from != to {
		t.Fatalf("moving from volume_ids to volumes changed the boot volume: %q -> %q", from, to)
	}
	// Reordering volume_ids changes the first volume, hence the boot volume.
	if bootVolumeId(nullVols, customStringList("d", "b"), strNull()) == from {
		t.Fatal("putting another volume first must change the boot volume")
	}
}

// ModifyPlan end to end: moving from volume_ids to volumes keeps the boot volume
// and must not replace the instance, while changing the boot volume must.
func TestUnitCloudInstanceModifyPlanBootVolume(t *testing.T) {
	ctx := context.Background()
	r := &cloudInstanceResource{}
	var sch fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &sch)

	raw := func(c volumeCfg) tftypes.Value { return instanceConfigRaw(t, sch.Schema, c) }
	with := func(f func(c *volumeCfg)) volumeCfg {
		c := newVolumeCfg()
		f(&c)
		return c
	}
	idsB := with(func(c *volumeCfg) { c.ids = customStringList("b", "d") })
	idsDB := with(func(c *volumeCfg) { c.ids = customStringList("d", "b") })
	volsB := with(func(c *volumeCfg) { c.volumes = volumeRefList(map[string]string{"b": "BOOT", "d": ""}) })
	volsD := with(func(c *volumeCfg) { c.volumes = volumeRefList(map[string]string{"b": "", "d": "BOOT"}) })
	withImage := with(func(c *volumeCfg) { c.image = strVal("img"); c.ids = customStringList("d") })

	// The state holds both views of the volumes, as the resource stores them: the
	// one the configuration did not set is computed from the API.
	stored := func(c volumeCfg) volumeCfg {
		if c.volumes.IsNull() {
			roles := map[string]string{}
			for i, ref := range customStringListToRefs(c.ids) {
				roles[ref.Id] = ""
				if i == 0 && knownString(c.image) == "" {
					roles[ref.Id] = "BOOT"
				}
			}
			c.volumes = volumeRefList(roles)
		}
		if c.ids.IsNull() {
			var list []string
			for _, ref := range volumesToAPI(c.volumes) {
				list = append(list, ref.Id)
			}
			c.ids = customStringList(list...)
		}
		return c
	}
	cases := []struct {
		name        string
		state       volumeCfg
		plan        volumeCfg
		wantReplace bool
	}{
		{"volume_ids -> volumes, same boot", idsB, volsB, false},
		{"volumes -> volume_ids, same boot", volsB, idsB, false},
		{"unchanged", volsB, volsB, false},
		{"volumes, another volume becomes BOOT", volsB, volsD, true},
		{"volume_ids, another volume first", idsB, idsDB, true},
		{"image: adding a data volume", withImage, with(func(c *volumeCfg) { c.image = strVal("img"); c.ids = customStringList("d", "e") }), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := fwresource.ModifyPlanRequest{
				State:  tfsdk.State{Schema: sch.Schema, Raw: raw(stored(tc.state))},
				Plan:   tfsdk.Plan{Schema: sch.Schema, Raw: raw(tc.plan)},
				Config: tfsdk.Config{Schema: sch.Schema, Raw: raw(tc.plan)},
			}
			resp := &fwresource.ModifyPlanResponse{}
			r.ModifyPlan(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			if got := len(resp.RequiresReplace) > 0; got != tc.wantReplace {
				t.Fatalf("replace = %v, want %v", got, tc.wantReplace)
			}
		})
	}
}
