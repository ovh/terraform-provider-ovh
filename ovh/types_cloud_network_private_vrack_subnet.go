package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudSubnetModel represents the Terraform model for the subnet resource
type CloudSubnetModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	NetworkId   ovhtypes.TfStringValue `tfsdk:"network_id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	CIDR        ovhtypes.TfStringValue `tfsdk:"cidr"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional
	Description     ovhtypes.TfStringValue `tfsdk:"description"`
	DHCPEnabled     types.Bool             `tfsdk:"dhcp_enabled"`
	DNSNameservers  types.List             `tfsdk:"dns_nameservers"`
	GatewayIP       ovhtypes.TfStringValue `tfsdk:"gateway_ip"`
	AllocationPools types.List             `tfsdk:"allocation_pools"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types
type CloudSubnetAPIResponse struct {
	Id             string                      `json:"id"`
	Checksum       string                      `json:"checksum"`
	CreatedAt      string                      `json:"createdAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	ResourceStatus string                      `json:"resourceStatus"`
	CurrentState   *CloudSubnetAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudSubnetAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudSubnetAPICurrentState struct {
	Name           string                    `json:"name,omitempty"`
	CIDR           string                    `json:"cidr,omitempty"`
	Description    string                    `json:"description,omitempty"`
	DHCPEnabled    *bool                     `json:"dhcpEnabled,omitempty"`
	DNSNameservers []string                  `json:"dnsNameservers,omitempty"`
	GatewayIP      string                    `json:"gatewayIp,omitempty"`
	HostRoutes     []CloudSubnetAPIHostRoute `json:"hostRoutes,omitempty"`
	Location       *CloudSubnetAPILocation   `json:"location,omitempty"`
}

type CloudSubnetAPILocation struct {
	Region string `json:"region"`
}

type CloudSubnetAPIHostRoute struct {
	Destination string `json:"destination"`
	NextHop     string `json:"nextHop"`
}

type CloudSubnetAPIAllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type CloudSubnetAPITargetSpec struct {
	Name            string                         `json:"name"`
	CIDR            string                         `json:"cidr,omitempty"`
	Description     string                         `json:"description,omitempty"`
	DHCPEnabled     *bool                          `json:"dhcpEnabled,omitempty"`
	DNSNameservers  []string                       `json:"dnsNameservers,omitempty"`
	GatewayIP       string                         `json:"gatewayIp,omitempty"`
	AllocationPools []CloudSubnetAPIAllocationPool `json:"allocationPools,omitempty"`
	Location        *CloudSubnetAPILocation        `json:"location,omitempty"`
}

type CloudSubnetAPIPutTargetSpec struct {
	Name            string                         `json:"name"`
	Description     string                         `json:"description,omitempty"`
	DHCPEnabled     *bool                          `json:"dhcpEnabled,omitempty"`
	DNSNameservers  []string                       `json:"dnsNameservers,omitempty"`
	GatewayIP       string                         `json:"gatewayIp,omitempty"`
	AllocationPools []CloudSubnetAPIAllocationPool `json:"allocationPools,omitempty"`
}

// Create payload
type CloudSubnetCreatePayload struct {
	TargetSpec *CloudSubnetAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudSubnetUpdatePayload struct {
	Checksum   string                       `json:"checksum"`
	TargetSpec *CloudSubnetAPIPutTargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudSubnetModel) ToCreate() *CloudSubnetCreatePayload {
	targetSpec := &CloudSubnetAPITargetSpec{
		Name: m.Name.ValueString(),
		CIDR: m.CIDR.ValueString(),
		Location: &CloudSubnetAPILocation{
			Region: m.Region.ValueString(),
		},
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.DHCPEnabled.IsNull() && !m.DHCPEnabled.IsUnknown() {
		val := m.DHCPEnabled.ValueBool()
		targetSpec.DHCPEnabled = &val
	}

	if !m.DNSNameservers.IsNull() && !m.DNSNameservers.IsUnknown() {
		dns := make([]string, 0, len(m.DNSNameservers.Elements()))
		for _, elem := range m.DNSNameservers.Elements() {
			if strVal, ok := elem.(types.String); ok {
				dns = append(dns, strVal.ValueString())
			}
		}
		targetSpec.DNSNameservers = dns
	}

	if !m.GatewayIP.IsNull() && !m.GatewayIP.IsUnknown() {
		targetSpec.GatewayIP = m.GatewayIP.ValueString()
	}

	if !m.AllocationPools.IsNull() && !m.AllocationPools.IsUnknown() {
		pools := make([]CloudSubnetAPIAllocationPool, 0, len(m.AllocationPools.Elements()))
		for _, elem := range m.AllocationPools.Elements() {
			obj, ok := elem.(basetypes.ObjectValue)
			if !ok {
				continue
			}
			attrs := obj.Attributes()
			pool := CloudSubnetAPIAllocationPool{}
			if v, ok := attrs["start"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				pool.Start = v.ValueString()
			}
			if v, ok := attrs["end"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				pool.End = v.ValueString()
			}
			pools = append(pools, pool)
		}
		targetSpec.AllocationPools = pools
	}

	return &CloudSubnetCreatePayload{
		TargetSpec: targetSpec,
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudSubnetModel) ToUpdate(checksum string) *CloudSubnetUpdatePayload {
	targetSpec := &CloudSubnetAPIPutTargetSpec{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	if !m.DHCPEnabled.IsNull() && !m.DHCPEnabled.IsUnknown() {
		val := m.DHCPEnabled.ValueBool()
		targetSpec.DHCPEnabled = &val
	}

	if !m.DNSNameservers.IsNull() && !m.DNSNameservers.IsUnknown() {
		dns := make([]string, 0, len(m.DNSNameservers.Elements()))
		for _, elem := range m.DNSNameservers.Elements() {
			if strVal, ok := elem.(types.String); ok {
				dns = append(dns, strVal.ValueString())
			}
		}
		targetSpec.DNSNameservers = dns
	}

	if !m.GatewayIP.IsNull() && !m.GatewayIP.IsUnknown() {
		targetSpec.GatewayIP = m.GatewayIP.ValueString()
	}

	if !m.AllocationPools.IsNull() && !m.AllocationPools.IsUnknown() {
		pools := make([]CloudSubnetAPIAllocationPool, 0, len(m.AllocationPools.Elements()))
		for _, elem := range m.AllocationPools.Elements() {
			obj, ok := elem.(basetypes.ObjectValue)
			if !ok {
				continue
			}
			attrs := obj.Attributes()
			pool := CloudSubnetAPIAllocationPool{}
			if v, ok := attrs["start"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				pool.Start = v.ValueString()
			}
			if v, ok := attrs["end"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				pool.End = v.ValueString()
			}
			pools = append(pools, pool)
		}
		targetSpec.AllocationPools = pools
	}

	return &CloudSubnetUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// SubnetCurrentStateAttrTypes returns the attribute types for the current_state object
func SubnetCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         ovhtypes.TfStringType{},
		"cidr":         ovhtypes.TfStringType{},
		"description":  ovhtypes.TfStringType{},
		"dhcp_enabled": types.BoolType,
		"dns_nameservers": types.ListType{
			ElemType: types.StringType,
		},
		"gateway_ip": ovhtypes.TfStringType{},
		"host_routes": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"destination": ovhtypes.TfStringType{},
					"next_hop":    ovhtypes.TfStringType{},
				},
			},
		},
		"location": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"region": ovhtypes.TfStringType{},
			},
		},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudSubnetModel) MergeWith(ctx context.Context, response *CloudSubnetAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildSubnetCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(SubnetCurrentStateAttrTypes())
	}

	// Update region from targetSpec if available
	if response.TargetSpec != nil && response.TargetSpec.Location != nil {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
	}
}

func buildSubnetCurrentStateObject(state *CloudSubnetAPICurrentState) basetypes.ObjectValue {
	// Build location object
	var locationObj basetypes.ObjectValue
	if state.Location != nil {
		locationObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"region": ovhtypes.TfStringType{},
			},
			map[string]attr.Value{
				"region": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.Region)},
			},
		)
	} else {
		locationObj = types.ObjectNull(map[string]attr.Type{
			"region": ovhtypes.TfStringType{},
		})
	}

	// Build dhcp_enabled value
	var dhcpVal attr.Value
	if state.DHCPEnabled != nil {
		dhcpVal = types.BoolValue(*state.DHCPEnabled)
	} else {
		dhcpVal = types.BoolNull()
	}

	// Build dns_nameservers list
	var dnsVal basetypes.ListValue
	if state.DNSNameservers != nil {
		dnsVals := make([]attr.Value, len(state.DNSNameservers))
		for i, dns := range state.DNSNameservers {
			dnsVals[i] = types.StringValue(dns)
		}
		dnsVal, _ = types.ListValue(types.StringType, dnsVals)
	} else {
		dnsVal = types.ListNull(types.StringType)
	}

	// Build host_routes list
	hostRouteAttrTypes := map[string]attr.Type{
		"destination": ovhtypes.TfStringType{},
		"next_hop":    ovhtypes.TfStringType{},
	}

	var hostRoutesVal basetypes.ListValue
	if state.HostRoutes != nil {
		routeObjs := make([]attr.Value, len(state.HostRoutes))
		for i, route := range state.HostRoutes {
			routeObj, _ := types.ObjectValue(
				hostRouteAttrTypes,
				map[string]attr.Value{
					"destination": ovhtypes.TfStringValue{StringValue: types.StringValue(route.Destination)},
					"next_hop":    ovhtypes.TfStringValue{StringValue: types.StringValue(route.NextHop)},
				},
			)
			routeObjs[i] = routeObj
		}
		hostRoutesVal, _ = types.ListValue(types.ObjectType{AttrTypes: hostRouteAttrTypes}, routeObjs)
	} else {
		hostRoutesVal = types.ListNull(types.ObjectType{AttrTypes: hostRouteAttrTypes})
	}

	currentStateObj, _ := types.ObjectValue(
		SubnetCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"cidr":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.CIDR)},
			"description":     ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"dhcp_enabled":    dhcpVal,
			"dns_nameservers": dnsVal,
			"gateway_ip":      ovhtypes.TfStringValue{StringValue: types.StringValue(state.GatewayIP)},
			"host_routes":     hostRoutesVal,
			"location":        locationObj,
		},
	)

	return currentStateObj
}
