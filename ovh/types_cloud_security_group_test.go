package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func sgTfStr(s string) ovhtypes.TfStringValue {
	if s == "" {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return ovhtypes.TfStringValue{StringValue: types.StringValue(s)}
}

func sgTcpRule(port int64) CloudSecurityGroupAPITargetRule {
	return CloudSecurityGroupAPITargetRule{
		Direction:      "INGRESS",
		EthernetType:   "IPV4",
		Protocol:       "TCP",
		PortRangeMin:   &port,
		PortRangeMax:   &port,
		RemoteIpPrefix: "0.0.0.0/0",
	}
}

func sgRuleObj(r CloudSecurityGroupAPITargetRule) attr.Value {
	return sgRuleObjWith(r, nil)
}

func sgRuleObjWith(r CloudSecurityGroupAPITargetRule, overrides map[string]attr.Value) attr.Value {
	portMin, portMax := types.Int64Null(), types.Int64Null()
	if r.PortRangeMin != nil {
		portMin = types.Int64Value(*r.PortRangeMin)
	}
	if r.PortRangeMax != nil {
		portMax = types.Int64Value(*r.PortRangeMax)
	}
	remoteGroupId := ""
	if r.RemoteGroup != nil {
		remoteGroupId = r.RemoteGroup.Id
	}
	attrs := map[string]attr.Value{
		"direction":        sgTfStr(r.Direction),
		"ethernet_type":    sgTfStr(r.EthernetType),
		"protocol":         sgTfStr(r.Protocol),
		"port_range_min":   portMin,
		"port_range_max":   portMax,
		"remote_group_id":  sgTfStr(remoteGroupId),
		"remote_ip_prefix": sgTfStr(r.RemoteIpPrefix),
		"description":      sgTfStr(r.Description),
	}
	for k, v := range overrides {
		attrs[k] = v
	}
	obj, _ := types.ObjectValue(SecurityGroupRuleAttrTypes(), attrs)
	return obj
}

func sgRuleList(rules ...CloudSecurityGroupAPITargetRule) types.List {
	elems := make([]attr.Value, len(rules))
	for i, r := range rules {
		elems[i] = sgRuleObj(r)
	}
	return sgRuleListOf(elems...)
}

func sgRuleListOf(elems ...attr.Value) types.List {
	l, _ := types.ListValue(types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()}, elems)
	return l
}

func sgMergeRules(prior types.List, apiRules ...CloudSecurityGroupAPITargetRule) types.List {
	m := &CloudSecurityGroupModel{Rule: prior}
	m.MergeWith(context.Background(), &CloudSecurityGroupAPIResponse{
		Id:         "sg-1",
		TargetSpec: &CloudSecurityGroupAPITargetSpec{Name: "sg", Rules: apiRules},
	}, prior)
	return m.Rule
}

func assertSgRules(t *testing.T, got types.List, want ...CloudSecurityGroupAPITargetRule) {
	t.Helper()
	if expected := sgRuleList(want...); !got.Equal(expected) {
		t.Fatalf("rule = %v, want %v", got, expected)
	}
}

// Issue #1612: config 80/443/22, the API answers in Neutron order 80/22/443.
func TestUnitCloudSecurityGroupRulesKeepPriorOrder(t *testing.T) {
	http, https, ssh := sgTcpRule(80), sgTcpRule(443), sgTcpRule(22)

	got := sgMergeRules(sgRuleList(http, https, ssh), http, ssh, https)
	assertSgRules(t, got, http, https, ssh)
}

func TestUnitCloudSecurityGroupRulesExtraAPIRuleAppended(t *testing.T) {
	http, https, ssh, rdp := sgTcpRule(80), sgTcpRule(443), sgTcpRule(22), sgTcpRule(3389)

	got := sgMergeRules(sgRuleList(http, https, ssh), rdp, http, ssh, https)
	assertSgRules(t, got, http, https, ssh, rdp)

	got = sgMergeRules(sgRuleList(http, https, ssh), ssh, http)
	assertSgRules(t, got, http, ssh)
}

func TestUnitCloudSecurityGroupRulesNullPriorKeepsAPIOrder(t *testing.T) {
	http, https, ssh := sgTcpRule(80), sgTcpRule(443), sgTcpRule(22)
	null := types.ListNull(types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()})

	got := sgMergeRules(null, http, ssh, https)
	assertSgRules(t, got, http, ssh, https)

	got = sgMergeRules(types.ListUnknown(types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()}), http, ssh, https)
	assertSgRules(t, got, http, ssh, https)
}

func TestUnitCloudSecurityGroupRulesNullOptionalFieldsMatch(t *testing.T) {
	egressAll := CloudSecurityGroupAPITargetRule{Direction: "EGRESS", EthernetType: "IPV6"}
	fromGroup := CloudSecurityGroupAPITargetRule{
		Direction:    "INGRESS",
		EthernetType: "IPV4",
		RemoteGroup:  &CloudSecurityGroupAPIRemoteGroup{Id: "sg-peer"},
		Description:  "peers",
	}
	http := sgTcpRule(80)

	got := sgMergeRules(sgRuleList(egressAll, http, fromGroup), fromGroup, http, egressAll)
	assertSgRules(t, got, egressAll, http, fromGroup)
}

func TestUnitCloudSecurityGroupRulesDuplicatesConsumedOnce(t *testing.T) {
	http, ssh := sgTcpRule(80), sgTcpRule(22)

	got := sgMergeRules(sgRuleList(http, http, ssh), ssh, http)
	assertSgRules(t, got, http, ssh)

	got = sgMergeRules(sgRuleList(ssh, http), http, ssh, http)
	assertSgRules(t, got, ssh, http, http)
}

// remote_group_id points at a group created in the same apply: unknown in the plan.
func TestUnitCloudSecurityGroupRulesUnknownPriorFieldMatchesAnyValue(t *testing.T) {
	peer := CloudSecurityGroupAPITargetRule{
		Direction:    "INGRESS",
		EthernetType: "IPV4",
		Protocol:     "TCP",
		RemoteGroup:  &CloudSecurityGroupAPIRemoteGroup{Id: "sg-created-in-same-apply"},
	}
	http, ssh := sgTcpRule(80), sgTcpRule(22)
	unknownRemoteGroup := map[string]attr.Value{"remote_group_id": ovhtypes.TfStringValue{StringValue: types.StringUnknown()}}

	prior := sgRuleListOf(sgRuleObjWith(peer, unknownRemoteGroup), sgRuleObj(http), sgRuleObj(ssh))
	got := sgMergeRules(prior, http, ssh, peer)
	assertSgRules(t, got, peer, http, ssh)
}

func TestUnitCloudSecurityGroupRulesExactMatchesBeforeUnknown(t *testing.T) {
	http, alt := sgTcpRule(80), sgTcpRule(8080)
	unknownPorts := map[string]attr.Value{"port_range_min": types.Int64Unknown(), "port_range_max": types.Int64Unknown()}

	prior := sgRuleListOf(sgRuleObjWith(alt, unknownPorts), sgRuleObj(http))
	got := sgMergeRules(prior, http, alt)
	assertSgRules(t, got, alt, http)
}

func TestUnitCloudSecurityGroupRulesUnknownPriorElement(t *testing.T) {
	http, ssh := sgTcpRule(80), sgTcpRule(22)

	prior := sgRuleListOf(types.ObjectUnknown(SecurityGroupRuleAttrTypes()), sgRuleObj(http))
	got := sgMergeRules(prior, http, ssh)
	assertSgRules(t, got, ssh, http)
}
