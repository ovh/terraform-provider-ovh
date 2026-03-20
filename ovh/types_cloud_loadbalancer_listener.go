package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudLoadbalancerListenerModel represents the Terraform model for the listener resource
type CloudLoadbalancerListenerModel struct {
	// Required — immutable
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	Protocol       ovhtypes.TfStringValue `tfsdk:"protocol"`
	ProtocolPort   types.Int64            `tfsdk:"protocol_port"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`

	// Optional — mutable
	Description            ovhtypes.TfStringValue                             `tfsdk:"description"`
	ConnectionLimit        types.Int64                                        `tfsdk:"connection_limit"`
	AllowedCidrs           ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"allowed_cidrs"`
	TimeoutClientData      types.Int64                                        `tfsdk:"timeout_client_data"`
	TimeoutMemberData      types.Int64                                        `tfsdk:"timeout_member_data"`
	TimeoutMemberConnect   types.Int64                                        `tfsdk:"timeout_member_connect"`
	TimeoutTcpInspect      types.Int64                                        `tfsdk:"timeout_tcp_inspect"`
	InsertHeaders          types.Object                                       `tfsdk:"insert_headers"`
	DefaultTlsContainerRef ovhtypes.TfStringValue                             `tfsdk:"default_tls_container_ref"`
	SniContainerRefs       ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"sni_container_refs"`
	TlsVersions            ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"tls_versions"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types

type CloudLoadbalancerListenerAPIResponse struct {
	Id             string                                     `json:"id"`
	Checksum       string                                     `json:"checksum"`
	CreatedAt      string                                     `json:"createdAt"`
	UpdatedAt      string                                     `json:"updatedAt"`
	ResourceStatus string                                     `json:"resourceStatus"`
	CurrentState   *CloudLoadbalancerListenerAPICurrentState  `json:"currentState,omitempty"`
	TargetSpec     *CloudLoadbalancerListenerAPITargetSpec    `json:"targetSpec,omitempty"`
}

type CloudLoadbalancerListenerAPIInsertHeaders struct {
	XForwardedFor     bool `json:"xForwardedFor,omitempty"`
	XForwardedPort    bool `json:"xForwardedPort,omitempty"`
	XForwardedProto   bool `json:"xForwardedProto,omitempty"`
	XSslClientVerify  bool `json:"xSslClientVerify,omitempty"`
	XSslClientHasCert bool `json:"xSslClientHasCert,omitempty"`
	XSslClientDn      bool `json:"xSslClientDn,omitempty"`
}

type CloudLoadbalancerListenerAPILocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudLoadbalancerListenerAPICurrentState struct {
	Name                   string                                     `json:"name,omitempty"`
	Description            string                                     `json:"description,omitempty"`
	Protocol               string                                     `json:"protocol,omitempty"`
	ProtocolPort           int                                        `json:"protocolPort,omitempty"`
	ConnectionLimit        *int                                       `json:"connectionLimit,omitempty"`
	AllowedCidrs           []string                                   `json:"allowedCidrs,omitempty"`
	TimeoutClientData      *int                                       `json:"timeoutClientData,omitempty"`
	TimeoutMemberData      *int                                       `json:"timeoutMemberData,omitempty"`
	TimeoutMemberConnect   *int                                       `json:"timeoutMemberConnect,omitempty"`
	TimeoutTcpInspect      *int                                       `json:"timeoutTcpInspect,omitempty"`
	InsertHeaders          *CloudLoadbalancerListenerAPIInsertHeaders `json:"insertHeaders,omitempty"`
	DefaultTlsContainerRef string                                     `json:"defaultTlsContainerRef,omitempty"`
	SniContainerRefs       []string                                   `json:"sniContainerRefs,omitempty"`
	TlsVersions            []string                                   `json:"tlsVersions,omitempty"`
	OperatingStatus        string                                     `json:"operatingStatus,omitempty"`
	ProvisioningStatus     string                                     `json:"provisioningStatus,omitempty"`
	Location               *CloudLoadbalancerListenerAPILocation      `json:"location,omitempty"`
}

type CloudLoadbalancerListenerAPITargetSpec struct {
	Name                   string                                     `json:"name"`
	Description            string                                     `json:"description,omitempty"`
	Protocol               string                                     `json:"protocol"`
	ProtocolPort           int                                        `json:"protocolPort"`
	ConnectionLimit        *int                                       `json:"connectionLimit,omitempty"`
	AllowedCidrs           []string                                   `json:"allowedCidrs,omitempty"`
	TimeoutClientData      *int                                       `json:"timeoutClientData,omitempty"`
	TimeoutMemberData      *int                                       `json:"timeoutMemberData,omitempty"`
	TimeoutMemberConnect   *int                                       `json:"timeoutMemberConnect,omitempty"`
	TimeoutTcpInspect      *int                                       `json:"timeoutTcpInspect,omitempty"`
	InsertHeaders          *CloudLoadbalancerListenerAPIInsertHeaders `json:"insertHeaders,omitempty"`
	DefaultTlsContainerRef string                                     `json:"defaultTlsContainerRef,omitempty"`
	SniContainerRefs       []string                                   `json:"sniContainerRefs,omitempty"`
	TlsVersions            []string                                   `json:"tlsVersions,omitempty"`
}

type CloudLoadbalancerListenerAPIUpdateTargetSpec struct {
	Name                   string                                     `json:"name"`
	Description            string                                     `json:"description,omitempty"`
	ConnectionLimit        *int                                       `json:"connectionLimit,omitempty"`
	AllowedCidrs           []string                                   `json:"allowedCidrs,omitempty"`
	TimeoutClientData      *int                                       `json:"timeoutClientData,omitempty"`
	TimeoutMemberData      *int                                       `json:"timeoutMemberData,omitempty"`
	TimeoutMemberConnect   *int                                       `json:"timeoutMemberConnect,omitempty"`
	TimeoutTcpInspect      *int                                       `json:"timeoutTcpInspect,omitempty"`
	InsertHeaders          *CloudLoadbalancerListenerAPIInsertHeaders `json:"insertHeaders,omitempty"`
	DefaultTlsContainerRef string                                     `json:"defaultTlsContainerRef,omitempty"`
	SniContainerRefs       []string                                   `json:"sniContainerRefs,omitempty"`
	TlsVersions            []string                                   `json:"tlsVersions,omitempty"`
}

// Create payload
type CloudLoadbalancerListenerCreatePayload struct {
	TargetSpec *CloudLoadbalancerListenerAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudLoadbalancerListenerUpdatePayload struct {
	Checksum   string                                        `json:"checksum"`
	TargetSpec *CloudLoadbalancerListenerAPIUpdateTargetSpec `json:"targetSpec"`
}

func intPtr(v int) *int { return &v }

// InsertHeadersAttrTypes returns the attribute types for the insert_headers object
func ListenerInsertHeadersAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"x_forwarded_for":       types.BoolType,
		"x_forwarded_port":      types.BoolType,
		"x_forwarded_proto":     types.BoolType,
		"x_ssl_client_verify":   types.BoolType,
		"x_ssl_client_has_cert": types.BoolType,
		"x_ssl_client_dn":       types.BoolType,
	}
}

// ListenerCurrentStateAttrTypes returns the attribute types for the current_state object
func ListenerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                      ovhtypes.TfStringType{},
		"description":               ovhtypes.TfStringType{},
		"protocol":                  ovhtypes.TfStringType{},
		"protocol_port":             types.Int64Type,
		"connection_limit":          types.Int64Type,
		"timeout_client_data":       types.Int64Type,
		"timeout_member_data":       types.Int64Type,
		"timeout_member_connect":    types.Int64Type,
		"timeout_tcp_inspect":       types.Int64Type,
		"operating_status":          ovhtypes.TfStringType{},
		"provisioning_status":       ovhtypes.TfStringType{},
		"default_tls_container_ref": ovhtypes.TfStringType{},
		"region":                    ovhtypes.TfStringType{},
		"availability_zone":         ovhtypes.TfStringType{},
		"insert_headers": types.ObjectType{
			AttrTypes: ListenerInsertHeadersAttrTypes(),
		},
		"allowed_cidrs": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
		"sni_container_refs": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
		"tls_versions": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
	}
}

func buildInsertHeadersFromAPI(ih *CloudLoadbalancerListenerAPIInsertHeaders) basetypes.ObjectValue {
	if ih == nil {
		obj, _ := types.ObjectValue(
			ListenerInsertHeadersAttrTypes(),
			map[string]attr.Value{
				"x_forwarded_for":       types.BoolValue(false),
				"x_forwarded_port":      types.BoolValue(false),
				"x_forwarded_proto":     types.BoolValue(false),
				"x_ssl_client_verify":   types.BoolValue(false),
				"x_ssl_client_has_cert": types.BoolValue(false),
				"x_ssl_client_dn":       types.BoolValue(false),
			},
		)
		return obj
	}

	obj, _ := types.ObjectValue(
		ListenerInsertHeadersAttrTypes(),
		map[string]attr.Value{
			"x_forwarded_for":       types.BoolValue(ih.XForwardedFor),
			"x_forwarded_port":      types.BoolValue(ih.XForwardedPort),
			"x_forwarded_proto":     types.BoolValue(ih.XForwardedProto),
			"x_ssl_client_verify":   types.BoolValue(ih.XSslClientVerify),
			"x_ssl_client_has_cert": types.BoolValue(ih.XSslClientHasCert),
			"x_ssl_client_dn":       types.BoolValue(ih.XSslClientDn),
		},
	)
	return obj
}

func buildStringListFromAPI(values []string) basetypes.ListValue {
	if values == nil {
		return types.ListNull(ovhtypes.TfStringType{})
	}
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
	}
	val, _ := types.ListValue(ovhtypes.TfStringType{}, elems)
	return val
}

func buildListenerCurrentStateObject(ctx context.Context, state *CloudLoadbalancerListenerAPICurrentState) types.Object {
	region := ""
	availabilityZone := ""
	if state.Location != nil {
		region = state.Location.Region
		availabilityZone = state.Location.AvailabilityZone
	}

	var connectionLimit attr.Value
	if state.ConnectionLimit != nil {
		connectionLimit = types.Int64Value(int64(*state.ConnectionLimit))
	} else {
		connectionLimit = types.Int64Null()
	}

	var timeoutClientData attr.Value
	if state.TimeoutClientData != nil {
		timeoutClientData = types.Int64Value(int64(*state.TimeoutClientData))
	} else {
		timeoutClientData = types.Int64Null()
	}

	var timeoutMemberData attr.Value
	if state.TimeoutMemberData != nil {
		timeoutMemberData = types.Int64Value(int64(*state.TimeoutMemberData))
	} else {
		timeoutMemberData = types.Int64Null()
	}

	var timeoutMemberConnect attr.Value
	if state.TimeoutMemberConnect != nil {
		timeoutMemberConnect = types.Int64Value(int64(*state.TimeoutMemberConnect))
	} else {
		timeoutMemberConnect = types.Int64Null()
	}

	var timeoutTcpInspect attr.Value
	if state.TimeoutTcpInspect != nil {
		timeoutTcpInspect = types.Int64Value(int64(*state.TimeoutTcpInspect))
	} else {
		timeoutTcpInspect = types.Int64Null()
	}

	obj, _ := types.ObjectValue(
		ListenerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                      ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":               ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"protocol":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.Protocol)},
			"protocol_port":             types.Int64Value(int64(state.ProtocolPort)),
			"connection_limit":          connectionLimit,
			"timeout_client_data":       timeoutClientData,
			"timeout_member_data":       timeoutMemberData,
			"timeout_member_connect":    timeoutMemberConnect,
			"timeout_tcp_inspect":       timeoutTcpInspect,
			"operating_status":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
			"default_tls_container_ref": ovhtypes.TfStringValue{StringValue: types.StringValue(state.DefaultTlsContainerRef)},
			"region":                    ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"availability_zone":         ovhtypes.TfStringValue{StringValue: types.StringValue(availabilityZone)},
			"insert_headers":            buildInsertHeadersFromAPI(state.InsertHeaders),
			"allowed_cidrs":             buildStringListFromAPI(state.AllowedCidrs),
			"sni_container_refs":        buildStringListFromAPI(state.SniContainerRefs),
			"tls_versions":              buildStringListFromAPI(state.TlsVersions),
		},
	)

	return obj
}

func extractInsertHeadersFromModel(m types.Object) *CloudLoadbalancerListenerAPIInsertHeaders {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}

	attrs := m.Attributes()
	ih := &CloudLoadbalancerListenerAPIInsertHeaders{}

	if v, ok := attrs["x_forwarded_for"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XForwardedFor = v.ValueBool()
	}
	if v, ok := attrs["x_forwarded_port"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XForwardedPort = v.ValueBool()
	}
	if v, ok := attrs["x_forwarded_proto"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XForwardedProto = v.ValueBool()
	}
	if v, ok := attrs["x_ssl_client_verify"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XSslClientVerify = v.ValueBool()
	}
	if v, ok := attrs["x_ssl_client_has_cert"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XSslClientHasCert = v.ValueBool()
	}
	if v, ok := attrs["x_ssl_client_dn"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ih.XSslClientDn = v.ValueBool()
	}

	return ih
}

func extractStringListFromModel(list ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	result := make([]string, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		if strVal, ok := elem.(ovhtypes.TfStringValue); ok {
			result = append(result, strVal.ValueString())
		}
	}
	return result
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudLoadbalancerListenerModel) ToCreate() *CloudLoadbalancerListenerCreatePayload {
	targetSpec := &CloudLoadbalancerListenerAPITargetSpec{
		Name:         m.Name.ValueString(),
		Protocol:     m.Protocol.ValueString(),
		ProtocolPort: int(m.ProtocolPort.ValueInt64()),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.ConnectionLimit.IsNull() && !m.ConnectionLimit.IsUnknown() {
		targetSpec.ConnectionLimit = intPtr(int(m.ConnectionLimit.ValueInt64()))
	}

	if !m.TimeoutClientData.IsNull() && !m.TimeoutClientData.IsUnknown() {
		targetSpec.TimeoutClientData = intPtr(int(m.TimeoutClientData.ValueInt64()))
	}

	if !m.TimeoutMemberData.IsNull() && !m.TimeoutMemberData.IsUnknown() {
		targetSpec.TimeoutMemberData = intPtr(int(m.TimeoutMemberData.ValueInt64()))
	}

	if !m.TimeoutMemberConnect.IsNull() && !m.TimeoutMemberConnect.IsUnknown() {
		targetSpec.TimeoutMemberConnect = intPtr(int(m.TimeoutMemberConnect.ValueInt64()))
	}

	if !m.TimeoutTcpInspect.IsNull() && !m.TimeoutTcpInspect.IsUnknown() {
		targetSpec.TimeoutTcpInspect = intPtr(int(m.TimeoutTcpInspect.ValueInt64()))
	}

	if !m.DefaultTlsContainerRef.IsNull() && !m.DefaultTlsContainerRef.IsUnknown() {
		targetSpec.DefaultTlsContainerRef = m.DefaultTlsContainerRef.ValueString()
	}

	targetSpec.InsertHeaders = extractInsertHeadersFromModel(m.InsertHeaders)
	targetSpec.AllowedCidrs = extractStringListFromModel(m.AllowedCidrs)
	targetSpec.SniContainerRefs = extractStringListFromModel(m.SniContainerRefs)
	targetSpec.TlsVersions = extractStringListFromModel(m.TlsVersions)

	return &CloudLoadbalancerListenerCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Note: protocol and protocolPort are immutable and not included
func (m *CloudLoadbalancerListenerModel) ToUpdate(checksum string) *CloudLoadbalancerListenerUpdatePayload {
	targetSpec := &CloudLoadbalancerListenerAPIUpdateTargetSpec{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.ConnectionLimit.IsNull() && !m.ConnectionLimit.IsUnknown() {
		targetSpec.ConnectionLimit = intPtr(int(m.ConnectionLimit.ValueInt64()))
	}

	if !m.TimeoutClientData.IsNull() && !m.TimeoutClientData.IsUnknown() {
		targetSpec.TimeoutClientData = intPtr(int(m.TimeoutClientData.ValueInt64()))
	}

	if !m.TimeoutMemberData.IsNull() && !m.TimeoutMemberData.IsUnknown() {
		targetSpec.TimeoutMemberData = intPtr(int(m.TimeoutMemberData.ValueInt64()))
	}

	if !m.TimeoutMemberConnect.IsNull() && !m.TimeoutMemberConnect.IsUnknown() {
		targetSpec.TimeoutMemberConnect = intPtr(int(m.TimeoutMemberConnect.ValueInt64()))
	}

	if !m.TimeoutTcpInspect.IsNull() && !m.TimeoutTcpInspect.IsUnknown() {
		targetSpec.TimeoutTcpInspect = intPtr(int(m.TimeoutTcpInspect.ValueInt64()))
	}

	if !m.DefaultTlsContainerRef.IsNull() && !m.DefaultTlsContainerRef.IsUnknown() {
		targetSpec.DefaultTlsContainerRef = m.DefaultTlsContainerRef.ValueString()
	}

	targetSpec.InsertHeaders = extractInsertHeadersFromModel(m.InsertHeaders)
	targetSpec.AllowedCidrs = extractStringListFromModel(m.AllowedCidrs)
	targetSpec.SniContainerRefs = extractStringListFromModel(m.SniContainerRefs)
	targetSpec.TlsVersions = extractStringListFromModel(m.TlsVersions)

	return &CloudLoadbalancerListenerUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudLoadbalancerListenerModel) MergeWith(ctx context.Context, response *CloudLoadbalancerListenerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildListenerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(ListenerCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Protocol = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Protocol)}
		m.ProtocolPort = types.Int64Value(int64(response.TargetSpec.ProtocolPort))

		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		// Handle nullable int fields from targetSpec
		if response.TargetSpec.ConnectionLimit != nil {
			m.ConnectionLimit = types.Int64Value(int64(*response.TargetSpec.ConnectionLimit))
		} else if m.ConnectionLimit.IsUnknown() {
			m.ConnectionLimit = types.Int64Null()
		}

		if response.TargetSpec.TimeoutClientData != nil {
			m.TimeoutClientData = types.Int64Value(int64(*response.TargetSpec.TimeoutClientData))
		} else if m.TimeoutClientData.IsUnknown() {
			m.TimeoutClientData = types.Int64Null()
		}

		if response.TargetSpec.TimeoutMemberData != nil {
			m.TimeoutMemberData = types.Int64Value(int64(*response.TargetSpec.TimeoutMemberData))
		} else if m.TimeoutMemberData.IsUnknown() {
			m.TimeoutMemberData = types.Int64Null()
		}

		if response.TargetSpec.TimeoutMemberConnect != nil {
			m.TimeoutMemberConnect = types.Int64Value(int64(*response.TargetSpec.TimeoutMemberConnect))
		} else if m.TimeoutMemberConnect.IsUnknown() {
			m.TimeoutMemberConnect = types.Int64Null()
		}

		if response.TargetSpec.TimeoutTcpInspect != nil {
			m.TimeoutTcpInspect = types.Int64Value(int64(*response.TargetSpec.TimeoutTcpInspect))
		} else if m.TimeoutTcpInspect.IsUnknown() {
			m.TimeoutTcpInspect = types.Int64Null()
		}

		// Handle defaultTlsContainerRef
		if response.TargetSpec.DefaultTlsContainerRef != "" || (!m.DefaultTlsContainerRef.IsNull() && !m.DefaultTlsContainerRef.IsUnknown()) {
			m.DefaultTlsContainerRef = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.DefaultTlsContainerRef)}
		}

		// Handle insert_headers from targetSpec
		if response.TargetSpec.InsertHeaders != nil {
			m.InsertHeaders = buildInsertHeadersFromAPI(response.TargetSpec.InsertHeaders)
		} else if m.InsertHeaders.IsNull() || m.InsertHeaders.IsUnknown() {
			m.InsertHeaders = types.ObjectNull(ListenerInsertHeadersAttrTypes())
		}

		// Handle list fields from targetSpec
		if response.TargetSpec.AllowedCidrs != nil {
			vals := make([]attr.Value, len(response.TargetSpec.AllowedCidrs))
			for i, v := range response.TargetSpec.AllowedCidrs {
				vals[i] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
			}
			m.AllowedCidrs = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals),
			}
		} else {
			m.AllowedCidrs = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListNull(ovhtypes.TfStringType{}),
			}
		}

		if response.TargetSpec.SniContainerRefs != nil {
			vals := make([]attr.Value, len(response.TargetSpec.SniContainerRefs))
			for i, v := range response.TargetSpec.SniContainerRefs {
				vals[i] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
			}
			m.SniContainerRefs = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals),
			}
		} else {
			m.SniContainerRefs = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListNull(ovhtypes.TfStringType{}),
			}
		}

		if response.TargetSpec.TlsVersions != nil {
			vals := make([]attr.Value, len(response.TargetSpec.TlsVersions))
			for i, v := range response.TargetSpec.TlsVersions {
				vals[i] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
			}
			m.TlsVersions = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals),
			}
		} else {
			m.TlsVersions = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListNull(ovhtypes.TfStringType{}),
			}
		}
	}
}
