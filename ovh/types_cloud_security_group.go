package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudSecurityGroupModel represents the Terraform model for the security group resource
type CloudSecurityGroupModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`

	// Optional — mutable
	Description ovhtypes.TfStringValue `tfsdk:"description"`
	Rule        types.List             `tfsdk:"rule"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudSecurityGroupAPIResponse struct {
	Id             string                             `json:"id"`
	Checksum       string                             `json:"checksum"`
	CreatedAt      string                             `json:"createdAt"`
	UpdatedAt      string                             `json:"updatedAt"`
	ResourceStatus string                             `json:"resourceStatus"`
	CurrentState   *CloudSecurityGroupAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudSecurityGroupAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudSecurityGroupAPICurrentState struct {
	Name         string                           `json:"name,omitempty"`
	Description  string                           `json:"description,omitempty"`
	Location     *CloudSecurityGroupAPILocation   `json:"location,omitempty"`
	Rules        []CloudSecurityGroupAPIStateRule `json:"rules,omitempty"`
	DefaultRules []CloudSecurityGroupAPIStateRule `json:"defaultRules,omitempty"`
}

type CloudSecurityGroupAPITargetSpec struct {
	Name        string                            `json:"name"`
	Description string                            `json:"description,omitempty"`
	Location    *CloudSecurityGroupAPILocation    `json:"location,omitempty"`
	Rules       []CloudSecurityGroupAPITargetRule `json:"rules,omitempty"`
}

type CloudSecurityGroupAPIUpdateTargetSpec struct {
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Rules       []CloudSecurityGroupAPITargetRule `json:"rules,omitempty"`
}

type CloudSecurityGroupAPILocation struct {
	Region string `json:"region,omitempty"`
}

type CloudSecurityGroupAPITargetRule struct {
	Direction      string                            `json:"direction"`
	EthernetType   string                            `json:"ethernetType"`
	Protocol       string                            `json:"protocol,omitempty"`
	PortRangeMin   *int64                            `json:"portRangeMin,omitempty"`
	PortRangeMax   *int64                            `json:"portRangeMax,omitempty"`
	RemoteGroup    *CloudSecurityGroupAPIRemoteGroup `json:"remoteGroup,omitempty"`
	RemoteIpPrefix string                            `json:"remoteIpPrefix,omitempty"`
	Description    string                            `json:"description,omitempty"`
}

type CloudSecurityGroupAPIStateRule struct {
	Id             string                            `json:"id,omitempty"`
	Direction      string                            `json:"direction"`
	EthernetType   string                            `json:"ethernetType"`
	Protocol       string                            `json:"protocol,omitempty"`
	PortRangeMin   *int64                            `json:"portRangeMin,omitempty"`
	PortRangeMax   *int64                            `json:"portRangeMax,omitempty"`
	RemoteGroup    *CloudSecurityGroupAPIRemoteGroup `json:"remoteGroup,omitempty"`
	RemoteIpPrefix string                            `json:"remoteIpPrefix,omitempty"`
	Description    string                            `json:"description,omitempty"`
}

type CloudSecurityGroupAPIRemoteGroup struct {
	Id string `json:"id"`
}

// Create payload
type CloudSecurityGroupCreatePayload struct {
	TargetSpec *CloudSecurityGroupAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudSecurityGroupUpdatePayload struct {
	Checksum   string                                 `json:"checksum"`
	TargetSpec *CloudSecurityGroupAPIUpdateTargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudSecurityGroupModel) ToCreate() *CloudSecurityGroupCreatePayload {
	targetSpec := &CloudSecurityGroupAPITargetSpec{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
		Location: &CloudSecurityGroupAPILocation{
			Region: m.Region.ValueString(),
		},
		Rules: buildSecurityGroupTargetRules(m.Rule),
	}

	return &CloudSecurityGroupCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudSecurityGroupModel) ToUpdate(checksum string) *CloudSecurityGroupUpdatePayload {
	return &CloudSecurityGroupUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudSecurityGroupAPIUpdateTargetSpec{
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
			Rules:       buildSecurityGroupTargetRules(m.Rule),
		},
	}
}

func buildSecurityGroupTargetRules(ruleList types.List) []CloudSecurityGroupAPITargetRule {
	if ruleList.IsNull() || ruleList.IsUnknown() {
		return nil
	}

	rules := make([]CloudSecurityGroupAPITargetRule, 0, len(ruleList.Elements()))
	for _, elem := range ruleList.Elements() {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		attrs := obj.Attributes()

		rule := CloudSecurityGroupAPITargetRule{}

		if v, ok := attrs["direction"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Direction = v.ValueString()
		}
		if v, ok := attrs["ethernet_type"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.EthernetType = v.ValueString()
		}
		if v, ok := attrs["protocol"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Protocol = v.ValueString()
		}
		if v, ok := attrs["port_range_min"].(basetypes.Int64Value); ok && !v.IsNull() && !v.IsUnknown() {
			val := v.ValueInt64()
			rule.PortRangeMin = &val
		}
		if v, ok := attrs["port_range_max"].(basetypes.Int64Value); ok && !v.IsNull() && !v.IsUnknown() {
			val := v.ValueInt64()
			rule.PortRangeMax = &val
		}
		if v, ok := attrs["remote_group_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.RemoteGroup = &CloudSecurityGroupAPIRemoteGroup{
				Id: v.ValueString(),
			}
		}
		if v, ok := attrs["remote_ip_prefix"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.RemoteIpPrefix = v.ValueString()
		}
		if v, ok := attrs["description"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Description = v.ValueString()
		}

		rules = append(rules, rule)
	}

	return rules
}

// SecurityGroupRuleAttrTypes returns the attribute types for a rule nested object
func SecurityGroupRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"direction":        ovhtypes.TfStringType{},
		"ethernet_type":    ovhtypes.TfStringType{},
		"protocol":         ovhtypes.TfStringType{},
		"port_range_min":   types.Int64Type,
		"port_range_max":   types.Int64Type,
		"remote_group_id":  ovhtypes.TfStringType{},
		"remote_ip_prefix": ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
	}
}

// SecurityGroupCurrentStateRuleAttrTypes returns the attribute types for a current_state rule
func SecurityGroupCurrentStateRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               ovhtypes.TfStringType{},
		"direction":        ovhtypes.TfStringType{},
		"ethernet_type":    ovhtypes.TfStringType{},
		"protocol":         ovhtypes.TfStringType{},
		"port_range_min":   types.Int64Type,
		"port_range_max":   types.Int64Type,
		"remote_group_id":  ovhtypes.TfStringType{},
		"remote_ip_prefix": ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
	}
}

// SecurityGroupCurrentStateAttrTypes returns the attribute types for current_state
func SecurityGroupCurrentStateAttrTypes() map[string]attr.Type {
	ruleObjType := types.ObjectType{
		AttrTypes: SecurityGroupCurrentStateRuleAttrTypes(),
	}
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"region":      ovhtypes.TfStringType{},
		"rules": types.ListType{
			ElemType: ruleObjType,
		},
		"default_rules": types.ListType{
			ElemType: ruleObjType,
		},
	}
}

// buildStateRulesList converts a slice of API state rules into a Terraform list value.
func buildStateRulesList(apiRules []CloudSecurityGroupAPIStateRule) basetypes.ListValue {
	ruleAttrTypes := SecurityGroupCurrentStateRuleAttrTypes()
	ruleObjType := types.ObjectType{AttrTypes: ruleAttrTypes}

	if apiRules == nil {
		return types.ListNull(ruleObjType)
	}

	ruleObjs := make([]attr.Value, len(apiRules))
	for i, rule := range apiRules {
		var portRangeMin, portRangeMax attr.Value
		if rule.PortRangeMin != nil {
			portRangeMin = types.Int64Value(*rule.PortRangeMin)
		} else {
			portRangeMin = types.Int64Null()
		}
		if rule.PortRangeMax != nil {
			portRangeMax = types.Int64Value(*rule.PortRangeMax)
		} else {
			portRangeMax = types.Int64Null()
		}

		remoteGroupId := ""
		if rule.RemoteGroup != nil {
			remoteGroupId = rule.RemoteGroup.Id
		}

		ruleObj, _ := types.ObjectValue(
			ruleAttrTypes,
			map[string]attr.Value{
				"id":               ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Id)},
				"direction":        ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Direction)},
				"ethernet_type":    ovhtypes.TfStringValue{StringValue: types.StringValue(rule.EthernetType)},
				"protocol":         ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Protocol)},
				"port_range_min":   portRangeMin,
				"port_range_max":   portRangeMax,
				"remote_group_id":  ovhtypes.TfStringValue{StringValue: types.StringValue(remoteGroupId)},
				"remote_ip_prefix": ovhtypes.TfStringValue{StringValue: types.StringValue(rule.RemoteIpPrefix)},
				"description":      ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Description)},
			},
		)
		ruleObjs[i] = ruleObj
	}
	val, _ := types.ListValue(ruleObjType, ruleObjs)
	return val
}

func buildSecurityGroupCurrentStateObject(ctx context.Context, state *CloudSecurityGroupAPICurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	obj, _ := types.ObjectValue(
		SecurityGroupCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":   ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"region":        ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"rules":         buildStateRulesList(state.Rules),
			"default_rules": buildStateRulesList(state.DefaultRules),
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudSecurityGroupModel) MergeWith(ctx context.Context, response *CloudSecurityGroupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildSecurityGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(SecurityGroupCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}

		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}

		// Rebuild rule list from targetSpec rules to keep Terraform state in sync
		if response.TargetSpec.Rules != nil {
			ruleAttrTypes := SecurityGroupRuleAttrTypes()
			ruleObjs := make([]attr.Value, len(response.TargetSpec.Rules))
			for i, rule := range response.TargetSpec.Rules {
				var portRangeMin, portRangeMax attr.Value
				if rule.PortRangeMin != nil {
					portRangeMin = types.Int64Value(*rule.PortRangeMin)
				} else {
					portRangeMin = types.Int64Null()
				}
				if rule.PortRangeMax != nil {
					portRangeMax = types.Int64Value(*rule.PortRangeMax)
				} else {
					portRangeMax = types.Int64Null()
				}

				remoteGroupId := ""
				if rule.RemoteGroup != nil {
					remoteGroupId = rule.RemoteGroup.Id
				}

				var remoteGroupIdVal attr.Value
				if remoteGroupId != "" {
					remoteGroupIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(remoteGroupId)}
				} else {
					remoteGroupIdVal = ovhtypes.TfStringValue{StringValue: types.StringNull()}
				}

				var protocolVal attr.Value
				if rule.Protocol != "" {
					protocolVal = ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Protocol)}
				} else {
					protocolVal = ovhtypes.TfStringValue{StringValue: types.StringNull()}
				}

				var remoteIpPrefixVal attr.Value
				if rule.RemoteIpPrefix != "" {
					remoteIpPrefixVal = ovhtypes.TfStringValue{StringValue: types.StringValue(rule.RemoteIpPrefix)}
				} else {
					remoteIpPrefixVal = ovhtypes.TfStringValue{StringValue: types.StringNull()}
				}

				var descriptionVal attr.Value
				if rule.Description != "" {
					descriptionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Description)}
				} else {
					descriptionVal = ovhtypes.TfStringValue{StringValue: types.StringNull()}
				}

				ruleObj, _ := types.ObjectValue(
					ruleAttrTypes,
					map[string]attr.Value{
						"direction":        ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Direction)},
						"ethernet_type":    ovhtypes.TfStringValue{StringValue: types.StringValue(rule.EthernetType)},
						"protocol":         protocolVal,
						"port_range_min":   portRangeMin,
						"port_range_max":   portRangeMax,
						"remote_group_id":  remoteGroupIdVal,
						"remote_ip_prefix": remoteIpPrefixVal,
						"description":      descriptionVal,
					},
				)
				ruleObjs[i] = ruleObj
			}
			m.Rule, _ = types.ListValue(types.ObjectType{AttrTypes: ruleAttrTypes}, ruleObjs)
		} else {
			m.Rule = types.ListNull(types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()})
		}
	}
}
