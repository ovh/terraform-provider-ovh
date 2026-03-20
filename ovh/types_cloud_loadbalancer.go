package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudLoadbalancerModel represents the Terraform model for the loadbalancer resource
type CloudLoadbalancerModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	VipNetworkId ovhtypes.TfStringValue `tfsdk:"vip_network_id"`
	VipSubnetId  ovhtypes.TfStringValue `tfsdk:"vip_subnet_id"`
	FlavorId     ovhtypes.TfStringValue `tfsdk:"flavor_id"`

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
	ID string `json:"id"`
}

type CloudLoadbalancerAPISubnetRef struct {
	ID string `json:"id"`
}

type CloudLoadbalancerAPIFlavorRef struct {
	ID string `json:"id"`
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
	Name               string                          `json:"name,omitempty"`
	Description        string                          `json:"description,omitempty"`
	Location           *CloudLoadbalancerAPILocation   `json:"location,omitempty"`
	VipAddress         string                          `json:"vipAddress,omitempty"`
	VipNetwork         *CloudLoadbalancerAPINetworkRef `json:"vipNetwork,omitempty"`
	VipSubnet          *CloudLoadbalancerAPISubnetRef  `json:"vipSubnet,omitempty"`
	OperatingStatus    string                          `json:"operatingStatus,omitempty"`
	ProvisioningStatus string                          `json:"provisioningStatus,omitempty"`
	Flavor             *CloudLoadbalancerAPIFlavorRef  `json:"flavor,omitempty"`
}

type CloudLoadbalancerAPITargetSpec struct {
	Name        string                          `json:"name"`
	Description string                          `json:"description,omitempty"`
	Location    *CloudLoadbalancerAPILocation   `json:"location,omitempty"`
	VipNetwork  *CloudLoadbalancerAPINetworkRef `json:"vipNetwork,omitempty"`
	VipSubnet   *CloudLoadbalancerAPISubnetRef  `json:"vipSubnet,omitempty"`
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

// ToCreate converts the Terraform model to the API create payload
func (m *CloudLoadbalancerModel) ToCreate() *CloudLoadbalancerCreatePayload {
	targetSpec := &CloudLoadbalancerAPITargetSpec{
		Name: m.Name.ValueString(),
		Location: &CloudLoadbalancerAPILocation{
			Region: m.Region.ValueString(),
		},
		VipNetwork: &CloudLoadbalancerAPINetworkRef{
			ID: m.VipNetworkId.ValueString(),
		},
		VipSubnet: &CloudLoadbalancerAPISubnetRef{
			ID: m.VipSubnetId.ValueString(),
		},
		Flavor: &CloudLoadbalancerAPIFlavorRef{
			ID: m.FlavorId.ValueString(),
		},
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
// Note: location and refs are immutable and not included in update payload
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

// loadbalancerRefAttrTypes returns the attr types for a ref object (vip_network, vip_subnet, flavor)
func loadbalancerRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

// LoadbalancerCurrentStateAttrTypes returns the attribute types for the current_state object
func LoadbalancerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                ovhtypes.TfStringType{},
		"description":         ovhtypes.TfStringType{},
		"vip_address":         ovhtypes.TfStringType{},
		"operating_status":    ovhtypes.TfStringType{},
		"provisioning_status": ovhtypes.TfStringType{},
		"region":              ovhtypes.TfStringType{},
		"availability_zone":   ovhtypes.TfStringType{},
		"vip_network": types.ObjectType{
			AttrTypes: loadbalancerRefAttrTypes(),
		},
		"vip_subnet": types.ObjectType{
			AttrTypes: loadbalancerRefAttrTypes(),
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

	// Build vip_network object
	var vipNetworkVal basetypes.ObjectValue
	if state.VipNetwork != nil {
		vipNetworkVal = buildLoadbalancerRefObject(state.VipNetwork.ID)
	} else {
		vipNetworkVal = types.ObjectNull(loadbalancerRefAttrTypes())
	}

	// Build vip_subnet object
	var vipSubnetVal basetypes.ObjectValue
	if state.VipSubnet != nil {
		vipSubnetVal = buildLoadbalancerRefObject(state.VipSubnet.ID)
	} else {
		vipSubnetVal = types.ObjectNull(loadbalancerRefAttrTypes())
	}

	// Build flavor object
	var flavorVal basetypes.ObjectValue
	if state.Flavor != nil {
		flavorVal = buildLoadbalancerRefObject(state.Flavor.ID)
	} else {
		flavorVal = types.ObjectNull(loadbalancerRefAttrTypes())
	}

	currentStateObj, _ := types.ObjectValue(
		LoadbalancerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"vip_address":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.VipAddress)},
			"operating_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.OperatingStatus)},
			"provisioning_status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProvisioningStatus)},
			"region":              regionVal,
			"availability_zone":   azVal,
			"vip_network":         vipNetworkVal,
			"vip_subnet":          vipSubnetVal,
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

		// Set flattened ref IDs from targetSpec
		if response.TargetSpec.VipNetwork != nil {
			m.VipNetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.VipNetwork.ID)}
		}
		if response.TargetSpec.VipSubnet != nil {
			m.VipSubnetId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.VipSubnet.ID)}
		}
		if response.TargetSpec.Flavor != nil {
			m.FlavorId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Flavor.ID)}
		}
	}
}
