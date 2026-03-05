package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudFloatingIpModel represents the Terraform model for the floating IP resource
type CloudFloatingIpModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`

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

// API response types
type CloudFloatingIpAPIResponse struct {
	Id             string                          `json:"id"`
	Checksum       string                          `json:"checksum"`
	CreatedAt      string                          `json:"createdAt"`
	UpdatedAt      string                          `json:"updatedAt"`
	ResourceStatus string                          `json:"resourceStatus"`
	CurrentState   *CloudFloatingIpAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudFloatingIpAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudFloatingIpAPICurrentState struct {
	IP          string                      `json:"ip,omitempty"`
	Status      string                      `json:"status,omitempty"`
	Network     *CloudFloatingIpAPINetwork  `json:"network,omitempty"`
	Description string                      `json:"description,omitempty"`
	Location    *CloudFloatingIpAPILocation `json:"location,omitempty"`
}

type CloudFloatingIpAPITargetSpec struct {
	Description string                      `json:"description,omitempty"`
	Location    *CloudFloatingIpAPILocation `json:"location,omitempty"`
}

type CloudFloatingIpAPIUpdateTargetSpec struct {
	Description string `json:"description"`
}

type CloudFloatingIpAPINetwork struct {
	ID string `json:"id"`
}

type CloudFloatingIpAPILocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// Create payload
type CloudFloatingIpCreatePayload struct {
	TargetSpec *CloudFloatingIpAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudFloatingIpUpdatePayload struct {
	Checksum   string                              `json:"checksum"`
	TargetSpec *CloudFloatingIpAPIUpdateTargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudFloatingIpModel) ToCreate() *CloudFloatingIpCreatePayload {
	targetSpec := &CloudFloatingIpAPITargetSpec{
		Description: m.Description.ValueString(),
		Location: &CloudFloatingIpAPILocation{
			Region: m.Region.ValueString(),
		},
	}

	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	return &CloudFloatingIpCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Only description is mutable — network_id, region, availability_zone are immutable
func (m *CloudFloatingIpModel) ToUpdate(checksum string) *CloudFloatingIpUpdatePayload {
	return &CloudFloatingIpUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudFloatingIpAPIUpdateTargetSpec{
			Description: m.Description.ValueString(),
		},
	}
}

// FloatingIpCurrentStateAttrTypes returns the attribute types for the current_state object
func FloatingIpCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":                ovhtypes.TfStringType{},
		"status":            ovhtypes.TfStringType{},
		"network_id":        ovhtypes.TfStringType{},
		"description":       ovhtypes.TfStringType{},
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// buildFloatingIpCurrentStateObject constructs the current_state object from the API response
func buildFloatingIpCurrentStateObject(ctx context.Context, state *CloudFloatingIpAPICurrentState) types.Object {
	networkId := ""
	if state.Network != nil {
		networkId = state.Network.ID
	}

	region := ""
	availabilityZone := ""
	if state.Location != nil {
		region = state.Location.Region
		availabilityZone = state.Location.AvailabilityZone
	}

	obj, _ := types.ObjectValue(
		FloatingIpCurrentStateAttrTypes(),
		map[string]attr.Value{
			"ip":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.IP)},
			"status":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"network_id":        ovhtypes.TfStringValue{StringValue: types.StringValue(networkId)},
			"description":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(availabilityZone)},
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudFloatingIpModel) MergeWith(ctx context.Context, response *CloudFloatingIpAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildFloatingIpCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(FloatingIpCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
		}
	}
}
