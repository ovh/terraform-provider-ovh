package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudGatewayModel represents the Terraform model for the gateway resource
type CloudGatewayModel struct {
	// Required — immutable
	ServiceName      ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region           ovhtypes.TfStringValue `tfsdk:"region"`
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`

	// Optional — mutable
	Description     ovhtypes.TfStringValue                             `tfsdk:"description"`
	ExternalGateway types.Object                                       `tfsdk:"external_gateway"`
	SubnetIds       ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"subnet_ids"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types

type CloudGatewayAPIExternalGateway struct {
	Enabled bool   `json:"enabled"`
	Model   string `json:"model,omitempty"`
}

type CloudGatewayAPIResponse struct {
	Id             string                       `json:"id"`
	Checksum       string                       `json:"checksum"`
	CreatedAt      string                       `json:"createdAt"`
	UpdatedAt      string                       `json:"updatedAt"`
	ResourceStatus string                       `json:"resourceStatus"`
	CurrentState   *CloudGatewayAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudGatewayAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudGatewayAPICurrentState struct {
	Name            string                          `json:"name,omitempty"`
	Description     string                          `json:"description,omitempty"`
	Status          string                          `json:"status,omitempty"`
	ExternalGateway *CloudGatewayAPIExternalGateway `json:"externalGateway,omitempty"`
	ExternalIP      string                          `json:"externalIp,omitempty"`
	Subnets         []CloudGatewayAPISubnet         `json:"subnets,omitempty"`
	Location        *CloudGatewayAPILocation        `json:"location,omitempty"`
}

type CloudGatewayAPITargetSpec struct {
	Name            string                          `json:"name"`
	Description     string                          `json:"description,omitempty"`
	ExternalGateway *CloudGatewayAPIExternalGateway `json:"externalGateway,omitempty"`
	Subnets         []CloudGatewayAPISubnet         `json:"subnets,omitempty"`
	Location        *CloudGatewayAPILocation        `json:"location,omitempty"`
}

type CloudGatewayAPISubnet struct {
	Id string `json:"id"`
}

type CloudGatewayAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// Create payload
type CloudGatewayCreatePayload struct {
	TargetSpec *CloudGatewayAPITargetSpec `json:"targetSpec"`
}

// Update payload — uses a separate struct without Location (immutable)
type CloudGatewayUpdateTargetSpec struct {
	Name            string                          `json:"name"`
	Description     string                          `json:"description,omitempty"`
	ExternalGateway *CloudGatewayAPIExternalGateway `json:"externalGateway,omitempty"`
	Subnets         []CloudGatewayAPISubnet         `json:"subnets,omitempty"`
}

type CloudGatewayUpdatePayload struct {
	Checksum   string                        `json:"checksum"`
	TargetSpec *CloudGatewayUpdateTargetSpec `json:"targetSpec"`
}

// ExternalGatewayAttrTypes returns the attribute types for the external_gateway object
func ExternalGatewayAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"model":   ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudGatewayModel) ToCreate() *CloudGatewayCreatePayload {
	targetSpec := &CloudGatewayAPITargetSpec{
		Name: m.Name.ValueString(),
		Location: &CloudGatewayAPILocation{
			Region: m.Region.ValueString(),
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

	// Handle external_gateway
	if !m.ExternalGateway.IsNull() && !m.ExternalGateway.IsUnknown() {
		egAttrs := m.ExternalGateway.Attributes()
		eg := &CloudGatewayAPIExternalGateway{}
		if enabledVal, ok := egAttrs["enabled"].(types.Bool); ok {
			eg.Enabled = enabledVal.ValueBool()
		}
		if modelVal, ok := egAttrs["model"].(ovhtypes.TfStringValue); ok && !modelVal.IsNull() && !modelVal.IsUnknown() {
			eg.Model = modelVal.ValueString()
		}
		targetSpec.ExternalGateway = eg
	}

	// Handle optional subnet_ids
	if !m.SubnetIds.IsNull() && !m.SubnetIds.IsUnknown() {
		subnets := make([]CloudGatewayAPISubnet, 0)
		for _, elem := range m.SubnetIds.Elements() {
			if strVal, ok := elem.(ovhtypes.TfStringValue); ok {
				subnets = append(subnets, CloudGatewayAPISubnet{Id: strVal.ValueString()})
			}
		}
		targetSpec.Subnets = subnets
	}

	return &CloudGatewayCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Note: location is immutable and not included in update payload
func (m *CloudGatewayModel) ToUpdate(checksum string) *CloudGatewayUpdatePayload {
	targetSpec := &CloudGatewayUpdateTargetSpec{
		Name: m.Name.ValueString(),
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	// Handle external_gateway
	if !m.ExternalGateway.IsNull() && !m.ExternalGateway.IsUnknown() {
		egAttrs := m.ExternalGateway.Attributes()
		eg := &CloudGatewayAPIExternalGateway{}
		if enabledVal, ok := egAttrs["enabled"].(types.Bool); ok {
			eg.Enabled = enabledVal.ValueBool()
		}
		if modelVal, ok := egAttrs["model"].(ovhtypes.TfStringValue); ok && !modelVal.IsNull() && !modelVal.IsUnknown() {
			eg.Model = modelVal.ValueString()
		}
		targetSpec.ExternalGateway = eg
	}

	// Handle optional subnet_ids
	if !m.SubnetIds.IsNull() && !m.SubnetIds.IsUnknown() {
		subnets := make([]CloudGatewayAPISubnet, 0)
		for _, elem := range m.SubnetIds.Elements() {
			if strVal, ok := elem.(ovhtypes.TfStringValue); ok {
				subnets = append(subnets, CloudGatewayAPISubnet{Id: strVal.ValueString()})
			}
		}
		targetSpec.Subnets = subnets
	}

	return &CloudGatewayUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// gatewaySubnetAttrTypes returns the attr types for a single subnet object
func gatewaySubnetAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

// GatewayCurrentStateAttrTypes returns the attribute types for the current_state object
func GatewayCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"status":      ovhtypes.TfStringType{},
		"external_gateway": types.ObjectType{
			AttrTypes: ExternalGatewayAttrTypes(),
		},
		"external_ip": ovhtypes.TfStringType{},
		"subnets": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: gatewaySubnetAttrTypes(),
			},
		},
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// buildGatewayCurrentStateObject constructs the current_state object from API response
func buildGatewayCurrentStateObject(ctx context.Context, state *CloudGatewayAPICurrentState) basetypes.ObjectValue {
	// Build subnets list
	subnetObjType := types.ObjectType{AttrTypes: gatewaySubnetAttrTypes()}

	var subnetsVal basetypes.ListValue
	if state.Subnets != nil {
		subnetObjs := make([]attr.Value, len(state.Subnets))
		for i, subnet := range state.Subnets {
			subnetObj, _ := types.ObjectValue(
				gatewaySubnetAttrTypes(),
				map[string]attr.Value{
					"id": ovhtypes.TfStringValue{StringValue: types.StringValue(subnet.Id)},
				},
			)
			subnetObjs[i] = subnetObj
		}
		subnetsVal, _ = types.ListValue(subnetObjType, subnetObjs)
	} else {
		subnetsVal = types.ListNull(subnetObjType)
	}

	// Build region and availability_zone from location
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringValue("")}
	azVal := ovhtypes.TfStringValue{StringValue: types.StringValue("")}
	if state.Location != nil {
		regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.Region)}
		if state.Location.AvailabilityZone != "" {
			azVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.AvailabilityZone)}
		}
	}

	// Build external_gateway object
	var externalGatewayVal basetypes.ObjectValue
	if state.ExternalGateway != nil {
		modelVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
		if state.ExternalGateway.Model != "" {
			modelVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.ExternalGateway.Model)}
		}
		externalGatewayVal, _ = types.ObjectValue(
			ExternalGatewayAttrTypes(),
			map[string]attr.Value{
				"enabled": types.BoolValue(state.ExternalGateway.Enabled),
				"model":   modelVal,
			},
		)
	} else {
		externalGatewayVal, _ = types.ObjectValue(
			ExternalGatewayAttrTypes(),
			map[string]attr.Value{
				"enabled": types.BoolValue(false),
				"model":   ovhtypes.TfStringValue{StringValue: types.StringNull()},
			},
		)
	}

	currentStateObj, _ := types.ObjectValue(
		GatewayCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":              ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"status":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"external_gateway":  externalGatewayVal,
			"external_ip":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.ExternalIP)},
			"subnets":           subnetsVal,
			"region":            regionVal,
			"availability_zone": azVal,
		},
	)

	return currentStateObj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudGatewayModel) MergeWith(ctx context.Context, response *CloudGatewayAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildGatewayCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(GatewayCurrentStateAttrTypes())
	}

	// Set flattened root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		// Set external_gateway from targetSpec
		if response.TargetSpec.ExternalGateway != nil {
			modelVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
			if response.TargetSpec.ExternalGateway.Model != "" {
				modelVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.ExternalGateway.Model)}
			}
			m.ExternalGateway, _ = types.ObjectValue(
				ExternalGatewayAttrTypes(),
				map[string]attr.Value{
					"enabled": types.BoolValue(response.TargetSpec.ExternalGateway.Enabled),
					"model":   modelVal,
				},
			)
		} else {
			m.ExternalGateway = types.ObjectNull(ExternalGatewayAttrTypes())
		}

		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
		}

		// Set subnet_ids from targetSpec
		if response.TargetSpec.Subnets != nil {
			subnetVals := make([]attr.Value, len(response.TargetSpec.Subnets))
			for i, subnet := range response.TargetSpec.Subnets {
				subnetVals[i] = ovhtypes.TfStringValue{StringValue: types.StringValue(subnet.Id)}
			}
			m.SubnetIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, subnetVals),
			}
		} else {
			m.SubnetIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{
				ListValue: basetypes.NewListNull(ovhtypes.TfStringType{}),
			}
		}
	}
}
