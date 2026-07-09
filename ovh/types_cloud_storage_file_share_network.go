package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageFileShareNetworkModel represents the Terraform model for the file storage share network resource.
// Every settable attribute is immutable: the resource has no update route.
type CloudStorageFileShareNetworkModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	NetworkId   ovhtypes.TfStringValue `tfsdk:"network_id"`
	SubnetId    ovhtypes.TfStringValue `tfsdk:"subnet_id"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — immutable
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
type CloudStorageFileShareNetworkAPIResponse struct {
	Id             string                                       `json:"id"`
	Checksum       string                                       `json:"checksum"`
	CreatedAt      string                                       `json:"createdAt"`
	UpdatedAt      string                                       `json:"updatedAt"`
	ResourceStatus string                                       `json:"resourceStatus"`
	CurrentState   *CloudStorageFileShareNetworkAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageFileShareNetworkAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudStorageFileShareNetworkAPINetworkRef struct {
	Id string `json:"id"`
}

type CloudStorageFileShareNetworkAPISubnetRef struct {
	Id string `json:"id"`
}

type CloudStorageFileShareNetworkAPILocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudStorageFileShareNetworkAPITargetSpec struct {
	Name        string                                     `json:"name,omitempty"`
	Description string                                     `json:"description,omitempty"`
	Network     *CloudStorageFileShareNetworkAPINetworkRef `json:"network,omitempty"`
	Subnet      *CloudStorageFileShareNetworkAPISubnetRef  `json:"subnet,omitempty"`
	Location    *CloudStorageFileShareNetworkAPILocation   `json:"location,omitempty"`
}

type CloudStorageFileShareNetworkAPICurrentState struct {
	Name        string                                     `json:"name,omitempty"`
	Description string                                     `json:"description,omitempty"`
	Network     *CloudStorageFileShareNetworkAPINetworkRef `json:"network,omitempty"`
	Subnet      *CloudStorageFileShareNetworkAPISubnetRef  `json:"subnet,omitempty"`
	Location    *CloudStorageFileShareNetworkAPILocation   `json:"location,omitempty"`
}

// Create payload
type CloudStorageFileShareNetworkCreatePayload struct {
	TargetSpec *CloudStorageFileShareNetworkAPITargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageFileShareNetworkModel) ToCreate() *CloudStorageFileShareNetworkCreatePayload {
	target := &CloudStorageFileShareNetworkAPITargetSpec{
		Name:     m.Name.ValueString(),
		Network:  &CloudStorageFileShareNetworkAPINetworkRef{Id: m.NetworkId.ValueString()},
		Subnet:   &CloudStorageFileShareNetworkAPISubnetRef{Id: m.SubnetId.ValueString()},
		Location: &CloudStorageFileShareNetworkAPILocation{Region: m.Region.ValueString()},
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	return &CloudStorageFileShareNetworkCreatePayload{TargetSpec: target}
}

// fileShareNetworkLocationAttrTypes returns the attr types for the root-level
// location object exposed by the file share network data sources.
func fileShareNetworkLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// FileShareNetworkCurrentStateAttrTypes returns the attribute types for the current_state object
func FileShareNetworkCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"network_id":  ovhtypes.TfStringType{},
		"subnet_id":   ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		}},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageFileShareNetworkModel) MergeWith(ctx context.Context, response *CloudStorageFileShareNetworkAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildFileShareNetworkCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(FileShareNetworkCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		if response.TargetSpec.Network != nil && response.TargetSpec.Network.Id != "" {
			m.NetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.Id)}
		}
		if response.TargetSpec.Subnet != nil && response.TargetSpec.Subnet.Id != "" {
			m.SubnetId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Subnet.Id)}
		}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
	}
}

// buildFileShareNetworkCurrentStateObject constructs the current_state object from the API response
func buildFileShareNetworkCurrentStateObject(state *CloudStorageFileShareNetworkAPICurrentState) types.Object {
	region := ""
	az := ""
	if state.Location != nil {
		region = state.Location.Region
		az = state.Location.AvailabilityZone
	}

	locObj, _ := types.ObjectValue(
		map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(az)},
		},
	)

	networkId := ""
	if state.Network != nil {
		networkId = state.Network.Id
	}
	subnetId := ""
	if state.Subnet != nil {
		subnetId = state.Subnet.Id
	}

	obj, _ := types.ObjectValue(
		FileShareNetworkCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"network_id":  ovhtypes.TfStringValue{StringValue: types.StringValue(networkId)},
			"subnet_id":   ovhtypes.TfStringValue{StringValue: types.StringValue(subnetId)},
			"location":    locObj,
		},
	)

	return obj
}
