package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudInstanceBackupModel represents the Terraform model for the instance backup resource
type CloudInstanceBackupModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudInstanceBackupAPIResponse struct {
	Id             string                              `json:"id"`
	Checksum       string                              `json:"checksum"`
	CreatedAt      string                              `json:"createdAt"`
	UpdatedAt      string                              `json:"updatedAt"`
	ResourceStatus string                              `json:"resourceStatus"`
	CurrentState   *CloudInstanceBackupAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudInstanceBackupAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudInstanceBackupAPICurrentState struct {
	Name       string                             `json:"name,omitempty"`
	Location   *CloudInstanceBackupAPILocation    `json:"location,omitempty"`
	Instance   *CloudInstanceBackupAPIInstanceRef `json:"instance,omitempty"`
	Size       int64                              `json:"size,omitempty"`
	MinDisk    int                                `json:"minDisk,omitempty"`
	MinRam     int                                `json:"minRam,omitempty"`
	Visibility string                             `json:"visibility,omitempty"`
	Status     string                             `json:"status,omitempty"`
}

type CloudInstanceBackupAPITargetSpec struct {
	Name     string                             `json:"name"`
	Location *CloudInstanceBackupAPILocation    `json:"location"`
	Instance *CloudInstanceBackupAPIInstanceRef `json:"instance"`
}

type CloudInstanceBackupAPIInstanceRef struct {
	ID string `json:"id"`
}

type CloudInstanceBackupAPILocation struct {
	Region string `json:"region"`
}

// Create payload
type CloudInstanceBackupCreatePayload struct {
	TargetSpec *CloudInstanceBackupAPITargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudInstanceBackupModel) ToCreate() *CloudInstanceBackupCreatePayload {
	return &CloudInstanceBackupCreatePayload{
		TargetSpec: &CloudInstanceBackupAPITargetSpec{
			Name: m.Name.ValueString(),
			Location: &CloudInstanceBackupAPILocation{
				Region: m.Region.ValueString(),
			},
			Instance: &CloudInstanceBackupAPIInstanceRef{
				ID: m.InstanceId.ValueString(),
			},
		},
	}
}

// InstanceBackupCurrentStateAttrTypes returns the attribute types for the current_state object
func InstanceBackupCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"region":      ovhtypes.TfStringType{},
		"instance_id": ovhtypes.TfStringType{},
		"size":        types.Int64Type,
		"min_disk":    types.Int64Type,
		"min_ram":     types.Int64Type,
		"visibility":  ovhtypes.TfStringType{},
		"status":      ovhtypes.TfStringType{},
	}
}

// buildInstanceBackupCurrentStateObject constructs the current_state object from the API response
func buildInstanceBackupCurrentStateObject(ctx context.Context, state *CloudInstanceBackupAPICurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	instanceId := ""
	if state.Instance != nil {
		instanceId = state.Instance.ID
	}

	obj, _ := types.ObjectValue(
		InstanceBackupCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"region":      ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"instance_id": ovhtypes.TfStringValue{StringValue: types.StringValue(instanceId)},
			"size":        types.Int64Value(state.Size),
			"min_disk":    types.Int64Value(int64(state.MinDisk)),
			"min_ram":     types.Int64Value(int64(state.MinRam)),
			"visibility":  ovhtypes.TfStringValue{StringValue: types.StringValue(state.Visibility)},
			"status":      ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudInstanceBackupModel) MergeWith(ctx context.Context, response *CloudInstanceBackupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildInstanceBackupCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceBackupCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		if response.TargetSpec.Instance != nil {
			m.InstanceId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Instance.ID)}
		}
	}
}
