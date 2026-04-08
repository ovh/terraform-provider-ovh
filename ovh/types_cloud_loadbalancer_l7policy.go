package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudLoadbalancerL7PolicyModel represents the Terraform model for the L7 policy resource.
type CloudLoadbalancerL7PolicyModel struct {
	// Required — immutable
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	ListenerId     ovhtypes.TfStringValue `tfsdk:"listener_id"`

	// Required — mutable
	Action ovhtypes.TfStringValue `tfsdk:"action"`

	// Optional — mutable
	Name             ovhtypes.TfStringValue `tfsdk:"name"`
	Description      ovhtypes.TfStringValue `tfsdk:"description"`
	Position         types.Int64            `tfsdk:"position"`
	RedirectPrefix   ovhtypes.TfStringValue `tfsdk:"redirect_prefix"`
	RedirectUrl      ovhtypes.TfStringValue `tfsdk:"redirect_url"`
	RedirectHttpCode types.Int64            `tfsdk:"redirect_http_code"`
	RedirectPoolId   ovhtypes.TfStringValue `tfsdk:"redirect_pool_id"`
	Rules            types.List             `tfsdk:"rules"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API types

type CloudLoadbalancerL7PolicyAPIRedirectPool struct {
	ID string `json:"id"`
}

type CloudLoadbalancerL7PolicyAPIRuleSpec struct {
	RuleType    string `json:"type"`
	CompareType string `json:"compareType"`
	Value       string `json:"value"`
	Key         string `json:"key,omitempty"`
	Invert      bool   `json:"invert,omitempty"`
}

type CloudLoadbalancerL7PolicyAPIRuleState struct {
	ID                 string `json:"id"`
	RuleType           string `json:"type"`
	CompareType        string `json:"compareType"`
	Value              string `json:"value"`
	Key                string `json:"key,omitempty"`
	Invert             bool   `json:"invert,omitempty"`
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`
	OperatingStatus    string `json:"operatingStatus,omitempty"`
}

type CloudLoadbalancerL7PolicyAPITargetSpec struct {
	Name             string                                    `json:"name,omitempty"`
	Description      string                                    `json:"description,omitempty"`
	Action           string                                    `json:"action"`
	Position         *int32                                    `json:"position,omitempty"`
	RedirectPrefix   string                                    `json:"redirectPrefix,omitempty"`
	RedirectUrl      string                                    `json:"redirectUrl,omitempty"`
	RedirectHttpCode *int32                                    `json:"redirectHttpCode,omitempty"`
	RedirectPool     *CloudLoadbalancerL7PolicyAPIRedirectPool `json:"redirectPool,omitempty"`
	Rules            []CloudLoadbalancerL7PolicyAPIRuleSpec    `json:"rules,omitempty"`
}

type CloudLoadbalancerL7PolicyAPICurrentState struct {
	Name               string                                    `json:"name,omitempty"`
	Description        string                                    `json:"description,omitempty"`
	Action             string                                    `json:"action"`
	Position           int32                                     `json:"position,omitempty"`
	RedirectPrefix     string                                    `json:"redirectPrefix,omitempty"`
	RedirectUrl        string                                    `json:"redirectUrl,omitempty"`
	RedirectHttpCode   int32                                     `json:"redirectHttpCode,omitempty"`
	RedirectPool       *CloudLoadbalancerL7PolicyAPIRedirectPool `json:"redirectPool,omitempty"`
	Rules              []CloudLoadbalancerL7PolicyAPIRuleState   `json:"rules,omitempty"`
	OperatingStatus    string                                    `json:"operatingStatus,omitempty"`
	ProvisioningStatus string                                    `json:"provisioningStatus,omitempty"`
}

type CloudLoadbalancerL7PolicyAPIResponse struct {
	Id             string                                     `json:"id"`
	Checksum       string                                     `json:"checksum"`
	CreatedAt      string                                     `json:"createdAt"`
	UpdatedAt      string                                     `json:"updatedAt"`
	ResourceStatus string                                     `json:"resourceStatus"`
	CurrentState   *CloudLoadbalancerL7PolicyAPICurrentState  `json:"currentState,omitempty"`
	TargetSpec     *CloudLoadbalancerL7PolicyAPITargetSpec    `json:"targetSpec,omitempty"`
}

// Create payload
type CloudLoadbalancerL7PolicyCreatePayload struct {
	TargetSpec *CloudLoadbalancerL7PolicyAPITargetSpec `json:"targetSpec"`
}

// Update payload — all fields are mutable for L7 policy
type CloudLoadbalancerL7PolicyUpdateTargetSpec struct {
	Name             string                                    `json:"name,omitempty"`
	Description      string                                    `json:"description,omitempty"`
	Action           string                                    `json:"action"`
	Position         *int32                                    `json:"position,omitempty"`
	RedirectPrefix   string                                    `json:"redirectPrefix,omitempty"`
	RedirectUrl      string                                    `json:"redirectUrl,omitempty"`
	RedirectHttpCode *int32                                    `json:"redirectHttpCode,omitempty"`
	RedirectPool     *CloudLoadbalancerL7PolicyAPIRedirectPool `json:"redirectPool,omitempty"`
	Rules            []CloudLoadbalancerL7PolicyAPIRuleSpec    `json:"rules,omitempty"`
}

type CloudLoadbalancerL7PolicyUpdatePayload struct {
	Checksum   string                                     `json:"checksum"`
	TargetSpec *CloudLoadbalancerL7PolicyUpdateTargetSpec  `json:"targetSpec"`
}

// l7PolicyRuleAttrTypes returns the attribute types for a single rule in the rules list.
func l7PolicyRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":         ovhtypes.TfStringType{},
		"compare_type": ovhtypes.TfStringType{},
		"value":        ovhtypes.TfStringType{},
		"key":          ovhtypes.TfStringType{},
		"invert":       types.BoolType,
	}
}

// l7PolicyRuleElementType returns the object type for the rules list element.
func l7PolicyRuleElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: l7PolicyRuleAttrTypes(),
	}
}

// l7PolicyCurrentStateRuleAttrTypes returns the attribute types for a rule in current_state.
func l7PolicyCurrentStateRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   ovhtypes.TfStringType{},
		"type":                 ovhtypes.TfStringType{},
		"compare_type":         ovhtypes.TfStringType{},
		"value":                ovhtypes.TfStringType{},
		"key":                  ovhtypes.TfStringType{},
		"invert":               types.BoolType,
		"operating_status":     ovhtypes.TfStringType{},
		"provisioning_status":  ovhtypes.TfStringType{},
	}
}

// l7PolicyCurrentStateRuleElementType returns the object type for the current_state rules list element.
func l7PolicyCurrentStateRuleElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: l7PolicyCurrentStateRuleAttrTypes(),
	}
}

// L7PolicyCurrentStateAttrTypes returns the attribute types for the current_state object.
func L7PolicyCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                ovhtypes.TfStringType{},
		"description":         ovhtypes.TfStringType{},
		"action":              ovhtypes.TfStringType{},
		"position":            types.Int64Type,
		"redirect_prefix":     ovhtypes.TfStringType{},
		"redirect_url":        ovhtypes.TfStringType{},
		"redirect_http_code":  types.Int64Type,
		"redirect_pool_id":    ovhtypes.TfStringType{},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
		"rules": types.ListType{
			ElemType: l7PolicyCurrentStateRuleElementType(),
		},
	}
}

// extractRulesFromModel extracts rules from the Terraform model list into API rule specs.
func extractRulesFromModel(rulesList types.List) []CloudLoadbalancerL7PolicyAPIRuleSpec {
	if rulesList.IsNull() || rulesList.IsUnknown() {
		return nil
	}

	elements := rulesList.Elements()
	if len(elements) == 0 {
		return nil
	}

	rules := make([]CloudLoadbalancerL7PolicyAPIRuleSpec, 0, len(elements))
	for _, elem := range elements {
		objVal, ok := elem.(types.Object)
		if !ok {
			continue
		}
		attrs := objVal.Attributes()

		rule := CloudLoadbalancerL7PolicyAPIRuleSpec{}

		if v, ok := attrs["type"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.RuleType = v.ValueString()
		}
		if v, ok := attrs["compare_type"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.CompareType = v.ValueString()
		}
		if v, ok := attrs["value"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Value = v.ValueString()
		}
		if v, ok := attrs["key"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Key = v.ValueString()
		}
		if v, ok := attrs["invert"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
			rule.Invert = v.ValueBool()
		}

		rules = append(rules, rule)
	}

	return rules
}

// ToCreate converts the Terraform model to the API create payload.
func (m *CloudLoadbalancerL7PolicyModel) ToCreate() *CloudLoadbalancerL7PolicyCreatePayload {
	targetSpec := &CloudLoadbalancerL7PolicyAPITargetSpec{
		Action: m.Action.ValueString(),
	}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.Position.IsNull() && !m.Position.IsUnknown() {
		v := int32(m.Position.ValueInt64())
		targetSpec.Position = &v
	}

	if !m.RedirectPrefix.IsNull() && !m.RedirectPrefix.IsUnknown() {
		targetSpec.RedirectPrefix = m.RedirectPrefix.ValueString()
	}

	if !m.RedirectUrl.IsNull() && !m.RedirectUrl.IsUnknown() {
		targetSpec.RedirectUrl = m.RedirectUrl.ValueString()
	}

	if !m.RedirectHttpCode.IsNull() && !m.RedirectHttpCode.IsUnknown() {
		v := int32(m.RedirectHttpCode.ValueInt64())
		targetSpec.RedirectHttpCode = &v
	}

	if !m.RedirectPoolId.IsNull() && !m.RedirectPoolId.IsUnknown() {
		targetSpec.RedirectPool = &CloudLoadbalancerL7PolicyAPIRedirectPool{
			ID: m.RedirectPoolId.ValueString(),
		}
	}

	targetSpec.Rules = extractRulesFromModel(m.Rules)

	return &CloudLoadbalancerL7PolicyCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload.
func (m *CloudLoadbalancerL7PolicyModel) ToUpdate(checksum string) *CloudLoadbalancerL7PolicyUpdatePayload {
	targetSpec := &CloudLoadbalancerL7PolicyUpdateTargetSpec{
		Action: m.Action.ValueString(),
	}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.Position.IsNull() && !m.Position.IsUnknown() {
		v := int32(m.Position.ValueInt64())
		targetSpec.Position = &v
	}

	if !m.RedirectPrefix.IsNull() && !m.RedirectPrefix.IsUnknown() {
		targetSpec.RedirectPrefix = m.RedirectPrefix.ValueString()
	}

	if !m.RedirectUrl.IsNull() && !m.RedirectUrl.IsUnknown() {
		targetSpec.RedirectUrl = m.RedirectUrl.ValueString()
	}

	if !m.RedirectHttpCode.IsNull() && !m.RedirectHttpCode.IsUnknown() {
		v := int32(m.RedirectHttpCode.ValueInt64())
		targetSpec.RedirectHttpCode = &v
	}

	if !m.RedirectPoolId.IsNull() && !m.RedirectPoolId.IsUnknown() {
		targetSpec.RedirectPool = &CloudLoadbalancerL7PolicyAPIRedirectPool{
			ID: m.RedirectPoolId.ValueString(),
		}
	}

	targetSpec.Rules = extractRulesFromModel(m.Rules)

	return &CloudLoadbalancerL7PolicyUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// buildL7PolicyRuleObject constructs a single rule object from API rule spec (for targetSpec).
func buildL7PolicyRuleObject(rule CloudLoadbalancerL7PolicyAPIRuleSpec) basetypes.ObjectValue {
	keyVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if rule.Key != "" {
		keyVal = ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Key)}
	}

	obj, _ := types.ObjectValue(
		l7PolicyRuleAttrTypes(),
		map[string]attr.Value{
			"type":         ovhtypes.TfStringValue{StringValue: types.StringValue(rule.RuleType)},
			"compare_type": ovhtypes.TfStringValue{StringValue: types.StringValue(rule.CompareType)},
			"value":        ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Value)},
			"key":          keyVal,
			"invert":       types.BoolValue(rule.Invert),
		},
	)
	return obj
}

// buildL7PolicyCurrentStateRuleObject constructs a single rule object from API rule state (for currentState).
func buildL7PolicyCurrentStateRuleObject(rule CloudLoadbalancerL7PolicyAPIRuleState) basetypes.ObjectValue {
	keyVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if rule.Key != "" {
		keyVal = ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Key)}
	}

	obj, _ := types.ObjectValue(
		l7PolicyCurrentStateRuleAttrTypes(),
		map[string]attr.Value{
			"id":                  ovhtypes.TfStringValue{StringValue: types.StringValue(rule.ID)},
			"type":                ovhtypes.TfStringValue{StringValue: types.StringValue(rule.RuleType)},
			"compare_type":        ovhtypes.TfStringValue{StringValue: types.StringValue(rule.CompareType)},
			"value":               ovhtypes.TfStringValue{StringValue: types.StringValue(rule.Value)},
			"key":                 keyVal,
			"invert":              types.BoolValue(rule.Invert),
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(rule.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(rule.ProvisioningStatus)},
		},
	)
	return obj
}

// buildL7PolicyRulesListFromTargetSpec builds the rules list from API targetSpec rules.
func buildL7PolicyRulesListFromTargetSpec(rules []CloudLoadbalancerL7PolicyAPIRuleSpec) basetypes.ListValue {
	if rules == nil {
		return types.ListNull(l7PolicyRuleElementType())
	}

	elems := make([]attr.Value, len(rules))
	for i, rule := range rules {
		elems[i] = buildL7PolicyRuleObject(rule)
	}

	val, _ := types.ListValue(l7PolicyRuleElementType(), elems)
	return val
}

// buildL7PolicyCurrentStateRulesList builds the rules list from API currentState rules.
func buildL7PolicyCurrentStateRulesList(rules []CloudLoadbalancerL7PolicyAPIRuleState) basetypes.ListValue {
	if rules == nil {
		return types.ListNull(l7PolicyCurrentStateRuleElementType())
	}

	elems := make([]attr.Value, len(rules))
	for i, rule := range rules {
		elems[i] = buildL7PolicyCurrentStateRuleObject(rule)
	}

	val, _ := types.ListValue(l7PolicyCurrentStateRuleElementType(), elems)
	return val
}

// buildL7PolicyCurrentStateObject constructs the current_state object from API response.
func buildL7PolicyCurrentStateObject(ctx context.Context, state *CloudLoadbalancerL7PolicyAPICurrentState) basetypes.ObjectValue {
	redirectPoolIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if state.RedirectPool != nil && state.RedirectPool.ID != "" {
		redirectPoolIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.RedirectPool.ID)}
	}

	obj, _ := types.ObjectValue(
		L7PolicyCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"action":              ovhtypes.TfStringValue{StringValue: types.StringValue(state.Action)},
			"position":            types.Int64Value(int64(state.Position)),
			"redirect_prefix":     ovhtypes.TfStringValue{StringValue: types.StringValue(state.RedirectPrefix)},
			"redirect_url":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.RedirectUrl)},
			"redirect_http_code":  types.Int64Value(int64(state.RedirectHttpCode)),
			"redirect_pool_id":    redirectPoolIdVal,
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
			"rules":               buildL7PolicyCurrentStateRulesList(state.Rules),
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model.
func (m *CloudLoadbalancerL7PolicyModel) MergeWith(ctx context.Context, response *CloudLoadbalancerL7PolicyAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildL7PolicyCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(L7PolicyCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Action = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Action)}

		// Keep name null if user didn't set it and API returns empty
		if response.TargetSpec.Name != "" || (!m.Name.IsNull() && !m.Name.IsUnknown()) {
			m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		}

		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		// Handle position
		if response.TargetSpec.Position != nil {
			m.Position = types.Int64Value(int64(*response.TargetSpec.Position))
		} else if m.Position.IsNull() || m.Position.IsUnknown() {
			m.Position = types.Int64Null()
		}

		// Handle redirect_prefix
		if response.TargetSpec.RedirectPrefix != "" || (!m.RedirectPrefix.IsNull() && !m.RedirectPrefix.IsUnknown()) {
			m.RedirectPrefix = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.RedirectPrefix)}
		}

		// Handle redirect_url
		if response.TargetSpec.RedirectUrl != "" || (!m.RedirectUrl.IsNull() && !m.RedirectUrl.IsUnknown()) {
			m.RedirectUrl = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.RedirectUrl)}
		}

		// Handle redirect_http_code
		if response.TargetSpec.RedirectHttpCode != nil {
			m.RedirectHttpCode = types.Int64Value(int64(*response.TargetSpec.RedirectHttpCode))
		} else if m.RedirectHttpCode.IsNull() || m.RedirectHttpCode.IsUnknown() {
			m.RedirectHttpCode = types.Int64Null()
		}

		// Handle redirect_pool_id (unwrap meta object)
		if response.TargetSpec.RedirectPool != nil && response.TargetSpec.RedirectPool.ID != "" {
			m.RedirectPoolId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.RedirectPool.ID)}
		} else if m.RedirectPoolId.IsNull() || m.RedirectPoolId.IsUnknown() {
			m.RedirectPoolId = ovhtypes.TfStringValue{StringValue: types.StringNull()}
		}

		// Handle rules
		if response.TargetSpec.Rules != nil {
			m.Rules = buildL7PolicyRulesListFromTargetSpec(response.TargetSpec.Rules)
		} else {
			if m.Rules.IsNull() || m.Rules.IsUnknown() {
				m.Rules = types.ListNull(l7PolicyRuleElementType())
			} else {
				m.Rules = types.ListNull(l7PolicyRuleElementType())
			}
		}
	}
}
