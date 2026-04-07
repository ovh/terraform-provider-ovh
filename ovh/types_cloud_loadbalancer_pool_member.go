package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudLoadbalancerPoolMemberModel represents the Terraform model for the pool member resource.
type CloudLoadbalancerPoolMemberModel struct {
	// Required — immutable
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	PoolId         ovhtypes.TfStringValue `tfsdk:"pool_id"`
	Address        ovhtypes.TfStringValue `tfsdk:"address"`
	ProtocolPort   types.Int64            `tfsdk:"protocol_port"`

	// Optional — immutable
	SubnetId ovhtypes.TfStringValue `tfsdk:"subnet_id"`

	// Optional — mutable
	Name    ovhtypes.TfStringValue `tfsdk:"name"`
	Weight  types.Int64            `tfsdk:"weight"`
	Backup  types.Bool             `tfsdk:"backup"`
	Monitor types.Object           `tfsdk:"monitor"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types

type CloudLoadbalancerPoolMemberAPIMonitor struct {
	Address string `json:"address,omitempty"`
	Port    int64  `json:"port,omitempty"`
}

type CloudLoadbalancerPoolMemberAPISubnetRef struct {
	Id string `json:"id"`
}

type CloudLoadbalancerPoolMemberAPIResponse struct {
	Id             string                                        `json:"id"`
	Checksum       string                                        `json:"checksum"`
	CreatedAt      string                                        `json:"createdAt"`
	UpdatedAt      string                                        `json:"updatedAt"`
	ResourceStatus string                                        `json:"resourceStatus"`
	CurrentState   *CloudLoadbalancerPoolMemberAPICurrentState   `json:"currentState,omitempty"`
	TargetSpec     *CloudLoadbalancerPoolMemberAPITargetSpec     `json:"targetSpec,omitempty"`
}

type CloudLoadbalancerPoolMemberAPICurrentState struct {
	Name               string                                    `json:"name,omitempty"`
	Address            string                                    `json:"address"`
	ProtocolPort       int64                                     `json:"protocolPort"`
	Weight             int64                                     `json:"weight"`
	Subnet             *CloudLoadbalancerPoolMemberAPISubnetRef  `json:"subnet,omitempty"`
	OperatingStatus    string                                    `json:"operatingStatus,omitempty"`
	ProvisioningStatus string                                    `json:"provisioningStatus,omitempty"`
	Backup             bool                                      `json:"backup,omitempty"`
	Monitor            *CloudLoadbalancerPoolMemberAPIMonitor    `json:"monitor,omitempty"`
}

// TargetSpec for POST (all fields including immutable)
type CloudLoadbalancerPoolMemberAPITargetSpec struct {
	Name         string                                  `json:"name,omitempty"`
	Address      string                                  `json:"address"`
	ProtocolPort int64                                   `json:"protocolPort"`
	Weight       *int64                                  `json:"weight,omitempty"`
	Subnet       *CloudLoadbalancerPoolMemberAPISubnetRef `json:"subnet,omitempty"`
	Monitor      *CloudLoadbalancerPoolMemberAPIMonitor  `json:"monitor,omitempty"`
	Backup       *bool                                   `json:"backup,omitempty"`
}

// UpdateTargetSpec for PUT (only mutable fields — no address, no protocolPort, no subnet)
type CloudLoadbalancerPoolMemberAPIUpdateTargetSpec struct {
	Name    string                                 `json:"name,omitempty"`
	Weight  *int64                                 `json:"weight,omitempty"`
	Monitor *CloudLoadbalancerPoolMemberAPIMonitor `json:"monitor,omitempty"`
	Backup  *bool                                  `json:"backup,omitempty"`
}

// Create payload
type CloudLoadbalancerPoolMemberCreatePayload struct {
	TargetSpec *CloudLoadbalancerPoolMemberAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudLoadbalancerPoolMemberUpdatePayload struct {
	Checksum   string                                        `json:"checksum"`
	TargetSpec *CloudLoadbalancerPoolMemberAPIUpdateTargetSpec `json:"targetSpec"`
}

// memberMonitorAttrTypes returns the attribute types for the monitor object.
func memberMonitorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": ovhtypes.TfStringType{},
		"port":    types.Int64Type,
	}
}

// ToCreate converts the Terraform model to the API create payload.
func (m *CloudLoadbalancerPoolMemberModel) ToCreate() *CloudLoadbalancerPoolMemberCreatePayload {
	targetSpec := &CloudLoadbalancerPoolMemberAPITargetSpec{
		Address:      m.Address.ValueString(),
		ProtocolPort: m.ProtocolPort.ValueInt64(),
	}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	if !m.Weight.IsNull() && !m.Weight.IsUnknown() {
		w := m.Weight.ValueInt64()
		targetSpec.Weight = &w
	}

	if !m.Backup.IsNull() && !m.Backup.IsUnknown() {
		b := m.Backup.ValueBool()
		targetSpec.Backup = &b
	}

	if !m.SubnetId.IsNull() && !m.SubnetId.IsUnknown() {
		targetSpec.Subnet = &CloudLoadbalancerPoolMemberAPISubnetRef{
			Id: m.SubnetId.ValueString(),
		}
	}

	if !m.Monitor.IsNull() && !m.Monitor.IsUnknown() {
		attrs := m.Monitor.Attributes()
		mon := &CloudLoadbalancerPoolMemberAPIMonitor{}
		if addrVal, ok := attrs["address"].(ovhtypes.TfStringValue); ok && !addrVal.IsNull() && !addrVal.IsUnknown() {
			mon.Address = addrVal.ValueString()
		}
		if portVal, ok := attrs["port"].(types.Int64); ok && !portVal.IsNull() && !portVal.IsUnknown() {
			mon.Port = portVal.ValueInt64()
		}
		targetSpec.Monitor = mon
	}

	return &CloudLoadbalancerPoolMemberCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload.
// Only mutable fields are included — no address, no protocolPort, no subnet.
func (m *CloudLoadbalancerPoolMemberModel) ToUpdate(checksum string) *CloudLoadbalancerPoolMemberUpdatePayload {
	targetSpec := &CloudLoadbalancerPoolMemberAPIUpdateTargetSpec{}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		targetSpec.Name = m.Name.ValueString()
	}

	if !m.Weight.IsNull() && !m.Weight.IsUnknown() {
		w := m.Weight.ValueInt64()
		targetSpec.Weight = &w
	}

	if !m.Backup.IsNull() && !m.Backup.IsUnknown() {
		b := m.Backup.ValueBool()
		targetSpec.Backup = &b
	}

	if !m.Monitor.IsNull() && !m.Monitor.IsUnknown() {
		attrs := m.Monitor.Attributes()
		mon := &CloudLoadbalancerPoolMemberAPIMonitor{}
		if addrVal, ok := attrs["address"].(ovhtypes.TfStringValue); ok && !addrVal.IsNull() && !addrVal.IsUnknown() {
			mon.Address = addrVal.ValueString()
		}
		if portVal, ok := attrs["port"].(types.Int64); ok && !portVal.IsNull() && !portVal.IsUnknown() {
			mon.Port = portVal.ValueInt64()
		}
		targetSpec.Monitor = mon
	}

	return &CloudLoadbalancerPoolMemberUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// MemberCurrentStateAttrTypes returns the attribute types for the current_state object.
func MemberCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                ovhtypes.TfStringType{},
		"address":             ovhtypes.TfStringType{},
		"protocol_port":       types.Int64Type,
		"weight":              types.Int64Type,
		"subnet_id":           ovhtypes.TfStringType{},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
		"backup":              types.BoolType,
		"monitor": types.ObjectType{
			AttrTypes: memberMonitorAttrTypes(),
		},
	}
}

// buildMemberMonitorObject constructs the monitor object from API response.
func buildMemberMonitorObject(mon *CloudLoadbalancerPoolMemberAPIMonitor) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(
		memberMonitorAttrTypes(),
		map[string]attr.Value{
			"address": ovhtypes.TfStringValue{StringValue: types.StringValue(mon.Address)},
			"port":    types.Int64Value(mon.Port),
		},
	)
	return obj
}

// buildMemberCurrentStateObject constructs the current_state object from API response.
func buildMemberCurrentStateObject(state *CloudLoadbalancerPoolMemberAPICurrentState) basetypes.ObjectValue {
	subnetIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if state.Subnet != nil {
		subnetIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Subnet.Id)}
	}

	var monitorVal basetypes.ObjectValue
	if state.Monitor != nil {
		monitorVal = buildMemberMonitorObject(state.Monitor)
	} else {
		monitorVal = types.ObjectNull(memberMonitorAttrTypes())
	}

	obj, _ := types.ObjectValue(
		MemberCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"address":             ovhtypes.TfStringValue{StringValue: types.StringValue(state.Address)},
			"protocol_port":       types.Int64Value(state.ProtocolPort),
			"weight":              types.Int64Value(state.Weight),
			"subnet_id":           subnetIdVal,
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
			"backup":              types.BoolValue(state.Backup),
			"monitor":             monitorVal,
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model.
func (m *CloudLoadbalancerPoolMemberModel) MergeWith(ctx context.Context, response *CloudLoadbalancerPoolMemberAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildMemberCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(MemberCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Address = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Address)}
		m.ProtocolPort = types.Int64Value(response.TargetSpec.ProtocolPort)

		if response.TargetSpec.Name != "" || (!m.Name.IsNull() && !m.Name.IsUnknown()) {
			m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		}

		if response.TargetSpec.Weight != nil {
			m.Weight = types.Int64Value(*response.TargetSpec.Weight)
		} else if m.Weight.IsNull() || m.Weight.IsUnknown() {
			m.Weight = types.Int64Null()
		}

		if response.TargetSpec.Backup != nil {
			m.Backup = types.BoolValue(*response.TargetSpec.Backup)
		} else if m.Backup.IsNull() || m.Backup.IsUnknown() {
			m.Backup = types.BoolNull()
		}

		if response.TargetSpec.Subnet != nil {
			m.SubnetId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Subnet.Id)}
		} else if m.SubnetId.IsNull() || m.SubnetId.IsUnknown() {
			m.SubnetId = ovhtypes.TfStringValue{StringValue: types.StringNull()}
		}

		if response.TargetSpec.Monitor != nil {
			m.Monitor = buildMemberMonitorObject(response.TargetSpec.Monitor)
		} else {
			m.Monitor = types.ObjectNull(memberMonitorAttrTypes())
		}
	}
}
