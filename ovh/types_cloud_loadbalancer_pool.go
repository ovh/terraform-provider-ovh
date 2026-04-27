package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudLoadbalancerPoolModel represents the Terraform model for the pool resource
type CloudLoadbalancerPoolModel struct {
	// Required — immutable
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	Protocol       ovhtypes.TfStringValue `tfsdk:"protocol"`

	// Required — mutable
	Algorithm ovhtypes.TfStringValue `tfsdk:"algorithm"`

	// Optional — mutable
	Name          ovhtypes.TfStringValue `tfsdk:"name"`
	Description   ovhtypes.TfStringValue `tfsdk:"description"`
	Persistence   types.Object           `tfsdk:"persistence"`
	HealthMonitor types.Object           `tfsdk:"health_monitor"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types

type CloudLoadbalancerPoolAPISessionPersistence struct {
	Type       string `json:"type"`
	CookieName string `json:"cookieName,omitempty"`
}

type CloudLoadbalancerPoolAPIHealthMonitor struct {
	Type           string `json:"type"`
	Delay          int    `json:"delay"`
	Timeout        int    `json:"timeout"`
	MaxRetries     int    `json:"maxRetries"`
	MaxRetriesDown *int   `json:"maxRetriesDown,omitempty"`
	Name           string `json:"name,omitempty"`
	URLPath        string `json:"urlPath,omitempty"`
	HTTPMethod     string `json:"httpMethod,omitempty"`
	HTTPVersion    string `json:"httpVersion,omitempty"`
	ExpectedCodes  string `json:"expectedCodes,omitempty"`
	DomainName     string `json:"domainName,omitempty"`
}

type CloudLoadbalancerPoolAPIHealthMonitorCurrentState struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Delay              int    `json:"delay"`
	Timeout            int    `json:"timeout"`
	MaxRetries         int    `json:"maxRetries"`
	MaxRetriesDown     int    `json:"maxRetriesDown,omitempty"`
	Name               string `json:"name,omitempty"`
	URLPath            string `json:"urlPath,omitempty"`
	HTTPMethod         string `json:"httpMethod,omitempty"`
	HTTPVersion        string `json:"httpVersion,omitempty"`
	ExpectedCodes      string `json:"expectedCodes,omitempty"`
	DomainName         string `json:"domainName,omitempty"`
	OperatingStatus    string `json:"operatingStatus,omitempty"`
	ProvisioningStatus string `json:"provisioningStatus,omitempty"`
}

type CloudLoadbalancerPoolAPIResponse struct {
	Id             string                                `json:"id"`
	Checksum       string                                `json:"checksum"`
	CreatedAt      string                                `json:"createdAt"`
	UpdatedAt      string                                `json:"updatedAt"`
	ResourceStatus string                                `json:"resourceStatus"`
	CurrentState   *CloudLoadbalancerPoolAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudLoadbalancerPoolAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudLoadbalancerPoolAPICurrentState struct {
	Name               string                                             `json:"name,omitempty"`
	Description        string                                             `json:"description,omitempty"`
	Protocol           string                                             `json:"protocol"`
	Algorithm          string                                             `json:"algorithm"`
	Persistence        *CloudLoadbalancerPoolAPISessionPersistence        `json:"persistence,omitempty"`
	HealthMonitor      *CloudLoadbalancerPoolAPIHealthMonitorCurrentState `json:"healthMonitor,omitempty"`
	OperatingStatus    string                                             `json:"operatingStatus,omitempty"`
	ProvisioningStatus string                                             `json:"provisioningStatus,omitempty"`
}

type CloudLoadbalancerPoolAPITargetSpec struct {
	Name          string                                      `json:"name,omitempty"`
	Description   string                                      `json:"description"`
	Protocol      string                                      `json:"protocol"`
	Algorithm     string                                      `json:"algorithm"`
	Persistence   *CloudLoadbalancerPoolAPISessionPersistence `json:"persistence,omitempty"`
	HealthMonitor *CloudLoadbalancerPoolAPIHealthMonitor      `json:"healthMonitor,omitempty"`
}

// Create payload
type CloudLoadbalancerPoolCreatePayload struct {
	TargetSpec *CloudLoadbalancerPoolAPITargetSpec `json:"targetSpec"`
}

// Update payload — uses a separate struct without protocol (immutable)
type CloudLoadbalancerPoolUpdateTargetSpec struct {
	Name          string                                      `json:"name,omitempty"`
	Description   string                                      `json:"description"`
	Algorithm     string                                      `json:"algorithm"`
	Persistence   *CloudLoadbalancerPoolAPISessionPersistence `json:"persistence,omitempty"`
	HealthMonitor *CloudLoadbalancerPoolAPIHealthMonitor      `json:"healthMonitor,omitempty"`
}

type CloudLoadbalancerPoolUpdatePayload struct {
	Checksum   string                                 `json:"checksum"`
	TargetSpec *CloudLoadbalancerPoolUpdateTargetSpec `json:"targetSpec"`
}

// poolPersistenceAttrTypes returns the attribute types for the persistence object
func poolPersistenceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":        ovhtypes.TfStringType{},
		"cookie_name": ovhtypes.TfStringType{},
	}
}

// healthMonitorAttrTypes returns the attribute types for the health_monitor object (targetSpec / user input)
func healthMonitorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":             ovhtypes.TfStringType{},
		"delay":            types.Int64Type,
		"timeout":          types.Int64Type,
		"max_retries":      types.Int64Type,
		"max_retries_down": types.Int64Type,
		"name":             ovhtypes.TfStringType{},
		"url_path":         ovhtypes.TfStringType{},
		"http_method":      ovhtypes.TfStringType{},
		"http_version":     ovhtypes.TfStringType{},
		"expected_codes":   ovhtypes.TfStringType{},
		"domain_name":      ovhtypes.TfStringType{},
	}
}

// healthMonitorCurrentStateAttrTypes returns the attribute types for the health_monitor in current_state
func healthMonitorCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  ovhtypes.TfStringType{},
		"type":                ovhtypes.TfStringType{},
		"delay":               types.Int64Type,
		"timeout":             types.Int64Type,
		"max_retries":         types.Int64Type,
		"max_retries_down":    types.Int64Type,
		"name":                ovhtypes.TfStringType{},
		"url_path":            ovhtypes.TfStringType{},
		"http_method":         ovhtypes.TfStringType{},
		"http_version":        ovhtypes.TfStringType{},
		"expected_codes":      ovhtypes.TfStringType{},
		"domain_name":         ovhtypes.TfStringType{},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
	}
}

// extractHealthMonitorFromModel extracts the health_monitor object from the Terraform model into an API struct
func extractHealthMonitorFromModel(obj types.Object) *CloudLoadbalancerPoolAPIHealthMonitor {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	attrs := obj.Attributes()
	hm := &CloudLoadbalancerPoolAPIHealthMonitor{}

	if v, ok := attrs["type"].(ovhtypes.TfStringValue); ok {
		hm.Type = v.ValueString()
	}
	if v, ok := attrs["delay"].(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
		hm.Delay = int(v.ValueInt64())
	}
	if v, ok := attrs["timeout"].(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
		hm.Timeout = int(v.ValueInt64())
	}
	if v, ok := attrs["max_retries"].(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
		hm.MaxRetries = int(v.ValueInt64())
	}
	if v, ok := attrs["max_retries_down"].(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
		val := int(v.ValueInt64())
		hm.MaxRetriesDown = &val
	}
	if v, ok := attrs["name"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.Name = v.ValueString()
	}
	if v, ok := attrs["url_path"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.URLPath = v.ValueString()
	}
	if v, ok := attrs["http_method"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.HTTPMethod = v.ValueString()
	}
	if v, ok := attrs["http_version"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.HTTPVersion = v.ValueString()
	}
	if v, ok := attrs["expected_codes"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.ExpectedCodes = v.ValueString()
	}
	if v, ok := attrs["domain_name"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		hm.DomainName = v.ValueString()
	}

	return hm
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudLoadbalancerPoolModel) ToCreate() *CloudLoadbalancerPoolCreatePayload {
	targetSpec := &CloudLoadbalancerPoolAPITargetSpec{
		Protocol:  m.Protocol.ValueString(),
		Algorithm: m.Algorithm.ValueString(),
	}

	// Handle optional name
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	// Handle optional persistence
	if !m.Persistence.IsNull() && !m.Persistence.IsUnknown() {
		attrs := m.Persistence.Attributes()
		p := &CloudLoadbalancerPoolAPISessionPersistence{}
		if typeVal, ok := attrs["type"].(ovhtypes.TfStringValue); ok {
			p.Type = typeVal.ValueString()
		}
		if cookieVal, ok := attrs["cookie_name"].(ovhtypes.TfStringValue); ok && !cookieVal.IsNull() && !cookieVal.IsUnknown() {
			p.CookieName = cookieVal.ValueString()
		}
		targetSpec.Persistence = p
	}

	// Handle optional health_monitor
	targetSpec.HealthMonitor = extractHealthMonitorFromModel(m.HealthMonitor)

	return &CloudLoadbalancerPoolCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Note: protocol is immutable and not included in update payload
func (m *CloudLoadbalancerPoolModel) ToUpdate(checksum string) *CloudLoadbalancerPoolUpdatePayload {
	targetSpec := &CloudLoadbalancerPoolUpdateTargetSpec{
		Algorithm: m.Algorithm.ValueString(),
	}

	// Handle optional name
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	// Handle optional persistence
	if !m.Persistence.IsNull() && !m.Persistence.IsUnknown() {
		attrs := m.Persistence.Attributes()
		p := &CloudLoadbalancerPoolAPISessionPersistence{}
		if typeVal, ok := attrs["type"].(ovhtypes.TfStringValue); ok {
			p.Type = typeVal.ValueString()
		}
		if cookieVal, ok := attrs["cookie_name"].(ovhtypes.TfStringValue); ok && !cookieVal.IsNull() && !cookieVal.IsUnknown() {
			p.CookieName = cookieVal.ValueString()
		}
		targetSpec.Persistence = p
	}

	// Handle optional health_monitor
	targetSpec.HealthMonitor = extractHealthMonitorFromModel(m.HealthMonitor)

	return &CloudLoadbalancerPoolUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// PoolCurrentStateAttrTypes returns the attribute types for the current_state object
func PoolCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"protocol":    ovhtypes.TfStringType{},
		"algorithm":   ovhtypes.TfStringType{},
		"persistence": types.ObjectType{
			AttrTypes: poolPersistenceAttrTypes(),
		},
		"health_monitor": types.ObjectType{
			AttrTypes: healthMonitorCurrentStateAttrTypes(),
		},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
	}
}

// buildPoolPersistenceObject constructs the persistence object from API response
func buildPoolPersistenceObject(p *CloudLoadbalancerPoolAPISessionPersistence) basetypes.ObjectValue {
	cookieNameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if p.CookieName != "" {
		cookieNameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(p.CookieName)}
	}

	obj, _ := types.ObjectValue(
		poolPersistenceAttrTypes(),
		map[string]attr.Value{
			"type":        ovhtypes.TfStringValue{StringValue: types.StringValue(p.Type)},
			"cookie_name": cookieNameVal,
		},
	)
	return obj
}

// buildHealthMonitorCurrentStateObject constructs the health_monitor object for current_state from API response
func buildHealthMonitorCurrentStateObject(hm *CloudLoadbalancerPoolAPIHealthMonitorCurrentState) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.Name)}
	urlPathVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.URLPath)}
	httpMethodVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.HTTPMethod)}
	httpVersionVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.HTTPVersion)}
	expectedCodesVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.ExpectedCodes)}
	domainNameVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.DomainName)}

	obj, _ := types.ObjectValue(
		healthMonitorCurrentStateAttrTypes(),
		map[string]attr.Value{
			"id":                  ovhtypes.TfStringValue{StringValue: types.StringValue(hm.ID)},
			"type":                ovhtypes.TfStringValue{StringValue: types.StringValue(hm.Type)},
			"delay":               types.Int64Value(int64(hm.Delay)),
			"timeout":             types.Int64Value(int64(hm.Timeout)),
			"max_retries":         types.Int64Value(int64(hm.MaxRetries)),
			"max_retries_down":    types.Int64Value(int64(hm.MaxRetriesDown)),
			"name":                nameVal,
			"url_path":            urlPathVal,
			"http_method":         httpMethodVal,
			"http_version":        httpVersionVal,
			"expected_codes":      expectedCodesVal,
			"domain_name":         domainNameVal,
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(hm.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(hm.ProvisioningStatus)},
		},
	)
	return obj
}

// buildHealthMonitorObject constructs the health_monitor object for targetSpec from API response
func buildHealthMonitorObject(hm *CloudLoadbalancerPoolAPIHealthMonitor) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.Name)}
	urlPathVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.URLPath)}
	httpMethodVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.HTTPMethod)}
	httpVersionVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.HTTPVersion)}
	expectedCodesVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.ExpectedCodes)}
	domainNameVal := ovhtypes.TfStringValue{StringValue: types.StringValue(hm.DomainName)}

	var maxRetriesDownVal attr.Value
	if hm.MaxRetriesDown != nil {
		maxRetriesDownVal = types.Int64Value(int64(*hm.MaxRetriesDown))
	} else {
		maxRetriesDownVal = types.Int64Null()
	}

	obj, _ := types.ObjectValue(
		healthMonitorAttrTypes(),
		map[string]attr.Value{
			"type":             ovhtypes.TfStringValue{StringValue: types.StringValue(hm.Type)},
			"delay":            types.Int64Value(int64(hm.Delay)),
			"timeout":          types.Int64Value(int64(hm.Timeout)),
			"max_retries":      types.Int64Value(int64(hm.MaxRetries)),
			"max_retries_down": maxRetriesDownVal,
			"name":             nameVal,
			"url_path":         urlPathVal,
			"http_method":      httpMethodVal,
			"http_version":     httpVersionVal,
			"expected_codes":   expectedCodesVal,
			"domain_name":      domainNameVal,
		},
	)
	return obj
}

// buildPoolCurrentStateObject constructs the current_state object from API response
func buildPoolCurrentStateObject(ctx context.Context, state *CloudLoadbalancerPoolAPICurrentState) basetypes.ObjectValue {
	// Build persistence object
	var persistenceVal basetypes.ObjectValue
	if state.Persistence != nil {
		persistenceVal = buildPoolPersistenceObject(state.Persistence)
	} else {
		persistenceVal = types.ObjectNull(poolPersistenceAttrTypes())
	}

	// Build health_monitor object
	var healthMonitorVal basetypes.ObjectValue
	if state.HealthMonitor != nil {
		healthMonitorVal = buildHealthMonitorCurrentStateObject(state.HealthMonitor)
	} else {
		healthMonitorVal = types.ObjectNull(healthMonitorCurrentStateAttrTypes())
	}

	currentStateObj, _ := types.ObjectValue(
		PoolCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"protocol":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Protocol)},
			"algorithm":           ovhtypes.TfStringValue{StringValue: types.StringValue(state.Algorithm)},
			"persistence":         persistenceVal,
			"health_monitor":      healthMonitorVal,
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
		},
	)

	return currentStateObj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudLoadbalancerPoolModel) MergeWith(ctx context.Context, response *CloudLoadbalancerPoolAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildPoolCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(PoolCurrentStateAttrTypes())
	}

	// Set flattened root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Protocol = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Protocol)}
		m.Algorithm = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Algorithm)}

		// Keep name null if user didn't set it and API returns empty
		if response.TargetSpec.Name != "" || (!m.Name.IsNull() && !m.Name.IsUnknown()) {
			m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		}

		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		// Set persistence from targetSpec
		if response.TargetSpec.Persistence != nil {
			m.Persistence = buildPoolPersistenceObject(response.TargetSpec.Persistence)
		} else {
			// Only set null if user didn't configure persistence
			if m.Persistence.IsNull() || m.Persistence.IsUnknown() {
				m.Persistence = types.ObjectNull(poolPersistenceAttrTypes())
			} else {
				m.Persistence = types.ObjectNull(poolPersistenceAttrTypes())
			}
		}

		// Set health_monitor from targetSpec
		if response.TargetSpec.HealthMonitor != nil {
			m.HealthMonitor = buildHealthMonitorObject(response.TargetSpec.HealthMonitor)
		} else {
			m.HealthMonitor = types.ObjectNull(healthMonitorAttrTypes())
		}
	}
}
