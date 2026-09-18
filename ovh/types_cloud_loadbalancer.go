package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudLoadbalancerModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Network     types.Object           `tfsdk:"network"`
	FlavorName  ovhtypes.TfStringValue `tfsdk:"flavor_name"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`

	// Optional — mutable
	Description ovhtypes.TfStringValue `tfsdk:"description"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types

type CloudLoadbalancerAPINetworkRef struct {
	ID       string `json:"id"`
	SubnetID string `json:"subnetId"`
	IP       string `json:"ip,omitempty"`
}

type CloudLoadbalancerAPIAddress struct {
	IP   string `json:"ip"`
	Type string `json:"type"`
}

type CloudLoadbalancerAPINetwork struct {
	ID        string                        `json:"id"`
	SubnetID  string                        `json:"subnetId"`
	Addresses []CloudLoadbalancerAPIAddress `json:"addresses,omitempty"`
}

type CloudLoadbalancerAPIFlavorRef struct {
	Name string `json:"name"`
}

type CloudLoadbalancerAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudLoadbalancerAPIResponse struct {
	Id             string                            `json:"id"`
	Checksum       string                            `json:"checksum"`
	CreatedAt      string                            `json:"createdAt"`
	UpdatedAt      string                            `json:"updatedAt"`
	ResourceStatus string                            `json:"resourceStatus"`
	CurrentState   *CloudLoadbalancerAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudLoadbalancerAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudLoadbalancerAPICurrentState struct {
	Name               string                         `json:"name,omitempty"`
	Description        string                         `json:"description,omitempty"`
	Location           *CloudLoadbalancerAPILocation  `json:"location,omitempty"`
	Network            *CloudLoadbalancerAPINetwork   `json:"network,omitempty"`
	OperatingStatus    string                         `json:"operatingStatus,omitempty"`
	ProvisioningStatus string                         `json:"provisioningStatus,omitempty"`
	Flavor             *CloudLoadbalancerAPIFlavorRef `json:"flavor,omitempty"`
}

type CloudLoadbalancerAPITargetSpec struct {
	Name        string                          `json:"name"`
	Description string                          `json:"description,omitempty"`
	Location    *CloudLoadbalancerAPILocation   `json:"location,omitempty"`
	Network     *CloudLoadbalancerAPINetworkRef `json:"network,omitempty"`
	Flavor      *CloudLoadbalancerAPIFlavorRef  `json:"flavor,omitempty"`
}

// Create payload
type CloudLoadbalancerCreatePayload struct {
	TargetSpec *CloudLoadbalancerAPITargetSpec `json:"targetSpec"`
}

// Update payload — uses a separate struct without immutable fields
type CloudLoadbalancerUpdateTargetSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CloudLoadbalancerUpdatePayload struct {
	Checksum   string                             `json:"checksum"`
	TargetSpec *CloudLoadbalancerUpdateTargetSpec `json:"targetSpec"`
}

func loadbalancerObjectString(attrs map[string]attr.Value, key string) string {
	v, ok := attrs[key]
	if !ok || v == nil || v.IsNull() || v.IsUnknown() {
		return ""
	}

	switch s := v.(type) {
	case ovhtypes.TfStringValue:
		return s.ValueString()
	case basetypes.StringValue:
		return s.ValueString()
	}

	return ""
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudLoadbalancerModel) ToCreate() *CloudLoadbalancerCreatePayload {
	targetSpec := &CloudLoadbalancerAPITargetSpec{
		Name: m.Name.ValueString(),
		Location: &CloudLoadbalancerAPILocation{
			Region: m.Region.ValueString(),
		},
		Flavor: &CloudLoadbalancerAPIFlavorRef{
			Name: m.FlavorName.ValueString(),
		},
	}

	if !m.Network.IsNull() && !m.Network.IsUnknown() {
		attrs := m.Network.Attributes()
		targetSpec.Network = &CloudLoadbalancerAPINetworkRef{
			ID:       loadbalancerObjectString(attrs, "id"),
			SubnetID: loadbalancerObjectString(attrs, "subnet_id"),
			IP:       loadbalancerObjectString(attrs, "ip"),
		}
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	// Handle optional availability zone
	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	return &CloudLoadbalancerCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Note: location, network and flavor are immutable and not included in update payload
func (m *CloudLoadbalancerModel) ToUpdate(checksum string) *CloudLoadbalancerUpdatePayload {
	targetSpec := &CloudLoadbalancerUpdateTargetSpec{
		Name: m.Name.ValueString(),
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	return &CloudLoadbalancerUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// loadbalancerRefAttrTypes returns the attr types for a ref object (flavor)
func loadbalancerRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

// LoadbalancerNetworkAttrTypes returns the attribute types for the target spec network object
func LoadbalancerNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        ovhtypes.TfStringType{},
		"subnet_id": ovhtypes.TfStringType{},
		"ip":        ovhtypes.TfStringType{},
	}
}

// LoadbalancerAddressAttrTypes returns the attribute types for a current state VIP address
func LoadbalancerAddressAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":   ovhtypes.TfStringType{},
		"type": ovhtypes.TfStringType{},
	}
}

// LoadbalancerCurrentStateNetworkAttrTypes returns the attribute types for the current state network object
func LoadbalancerCurrentStateNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        ovhtypes.TfStringType{},
		"subnet_id": ovhtypes.TfStringType{},
		"addresses": types.ListType{
			ElemType: types.ObjectType{AttrTypes: LoadbalancerAddressAttrTypes()},
		},
	}
}

// LoadbalancerCurrentStateAttrTypes returns the attribute types for the current_state object
func LoadbalancerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                ovhtypes.TfStringType{},
		"description":         ovhtypes.TfStringType{},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
		"region":              ovhtypes.TfStringType{},
		"availability_zone":   ovhtypes.TfStringType{},
		"network": types.ObjectType{
			AttrTypes: LoadbalancerCurrentStateNetworkAttrTypes(),
		},
		"flavor": types.ObjectType{
			AttrTypes: loadbalancerRefAttrTypes(),
		},
	}
}

// buildLoadbalancerRefObject constructs a ref object value from an ID string
func buildLoadbalancerRefObject(id string) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(
		loadbalancerRefAttrTypes(),
		map[string]attr.Value{
			"id": ovhtypes.TfStringValue{StringValue: types.StringValue(id)},
		},
	)
	return obj
}

func buildLoadbalancerCurrentStateNetworkObject(network *CloudLoadbalancerAPINetwork) basetypes.ObjectValue {
	if network == nil {
		return types.ObjectNull(LoadbalancerCurrentStateNetworkAttrTypes())
	}

	addressObjType := types.ObjectType{AttrTypes: LoadbalancerAddressAttrTypes()}

	addressesVal := types.ListNull(addressObjType)
	if network.Addresses != nil {
		elems := make([]attr.Value, len(network.Addresses))
		for i, address := range network.Addresses {
			elems[i], _ = types.ObjectValue(
				LoadbalancerAddressAttrTypes(),
				map[string]attr.Value{
					"ip":   ovhtypes.TfStringValue{StringValue: types.StringValue(address.IP)},
					"type": ovhtypes.TfStringValue{StringValue: types.StringValue(address.Type)},
				},
			)
		}
		addressesVal, _ = types.ListValue(addressObjType, elems)
	}

	obj, _ := types.ObjectValue(
		LoadbalancerCurrentStateNetworkAttrTypes(),
		map[string]attr.Value{
			"id":        ovhtypes.TfStringValue{StringValue: types.StringValue(network.ID)},
			"subnet_id": ovhtypes.TfStringValue{StringValue: types.StringValue(network.SubnetID)},
			"addresses": addressesVal,
		},
	)

	return obj
}

// buildLoadbalancerCurrentStateObject constructs the current_state object from API response
func buildLoadbalancerCurrentStateObject(ctx context.Context, state *CloudLoadbalancerAPICurrentState) basetypes.ObjectValue {
	// Build region and availability_zone from location
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringValue("")}
	azVal := ovhtypes.TfStringValue{StringValue: types.StringValue("")}
	if state.Location != nil {
		regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.Region)}
		if state.Location.AvailabilityZone != "" {
			azVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.AvailabilityZone)}
		}
	}

	// Build flavor object
	var flavorVal basetypes.ObjectValue
	if state.Flavor != nil {
		flavorVal = buildLoadbalancerRefObject(state.Flavor.Name)
	} else {
		flavorVal = types.ObjectNull(loadbalancerRefAttrTypes())
	}

	currentStateObj, _ := types.ObjectValue(
		LoadbalancerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
			"region":              regionVal,
			"availability_zone":   azVal,
			"network":             buildLoadbalancerCurrentStateNetworkObject(state.Network),
			"flavor":              flavorVal,
		},
	)

	return currentStateObj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudLoadbalancerModel) MergeWith(ctx context.Context, response *CloudLoadbalancerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildLoadbalancerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(LoadbalancerCurrentStateAttrTypes())
	}

	// Set flattened root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}

		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		// Set region from targetSpec location (immutable, always from targetSpec)
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
		}

		if response.TargetSpec.Network != nil {
			// The API never invents an ip: keep the configured one when it echoes back empty
			configuredIP := ""
			if !m.Network.IsNull() && !m.Network.IsUnknown() {
				configuredIP = loadbalancerObjectString(m.Network.Attributes(), "ip")
			}

			ip := response.TargetSpec.Network.IP
			if ip == "" {
				ip = configuredIP
			}

			ipVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
			if ip != "" {
				ipVal = ovhtypes.TfStringValue{StringValue: types.StringValue(ip)}
			}

			m.Network, _ = types.ObjectValue(
				LoadbalancerNetworkAttrTypes(),
				map[string]attr.Value{
					"id":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.ID)},
					"subnet_id": ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.SubnetID)},
					"ip":        ipVal,
				},
			)
		} else if m.Network.IsUnknown() {
			m.Network = types.ObjectNull(LoadbalancerNetworkAttrTypes())
		}

		if response.TargetSpec.Flavor != nil {
			m.FlavorName = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Flavor.Name)}
		}
	}
}
