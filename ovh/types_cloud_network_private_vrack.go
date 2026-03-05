package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudNetworkPrivateVrackModel represents the Terraform model for the network resource
type CloudNetworkPrivateVrackModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional
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
type CloudNetworkAPIResponse struct {
	Id             string                       `json:"id"`
	Checksum       string                       `json:"checksum"`
	CreatedAt      string                       `json:"createdAt"`
	UpdatedAt      string                       `json:"updatedAt"`
	ResourceStatus string                       `json:"resourceStatus"`
	CurrentState   *CloudNetworkAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudNetworkAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudNetworkAPICurrentState struct {
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Location    *CloudNetworkAPILocation `json:"location,omitempty"`
}

type CloudNetworkAPILocation struct {
	Region string `json:"region"`
}

type CloudNetworkAPITargetSpec struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Location    *CloudNetworkAPILocation `json:"location,omitempty"`
}

type CloudNetworkAPIPutTargetSpec struct {
	Name string `json:"name"`
}

// Create payload
type CloudNetworkCreatePayload struct {
	TargetSpec *CloudNetworkAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudNetworkUpdatePayload struct {
	Checksum   string                        `json:"checksum"`
	TargetSpec *CloudNetworkAPIPutTargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudNetworkPrivateVrackModel) ToCreate() *CloudNetworkCreatePayload {
	targetSpec := &CloudNetworkAPITargetSpec{
		Name: m.Name.ValueString(),
		Location: &CloudNetworkAPILocation{
			Region: m.Region.ValueString(),
		},
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	return &CloudNetworkCreatePayload{
		TargetSpec: targetSpec,
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudNetworkPrivateVrackModel) ToUpdate(checksum string) *CloudNetworkUpdatePayload {
	return &CloudNetworkUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudNetworkAPIPutTargetSpec{
			Name: m.Name.ValueString(),
		},
	}
}

// NetworkCurrentStateAttrTypes returns the attribute types for the current_state object
func NetworkCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"region": ovhtypes.TfStringType{},
			},
		},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudNetworkPrivateVrackModel) MergeWith(ctx context.Context, response *CloudNetworkAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildNetworkCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(NetworkCurrentStateAttrTypes())
	}

	// Update region from targetSpec if available
	if response.TargetSpec != nil && response.TargetSpec.Location != nil {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
	}
}

func buildNetworkCurrentStateObject(state *CloudNetworkAPICurrentState) basetypes.ObjectValue {
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

	currentStateObj, _ := types.ObjectValue(
		NetworkCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"location":    locationObj,
		},
	)

	return currentStateObj
}
