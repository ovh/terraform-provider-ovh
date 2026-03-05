package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageBlockVolumeBackupModel represents the Terraform model for the block storage backup resource
type CloudStorageBlockVolumeBackupModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	VolumeId    ovhtypes.TfStringValue `tfsdk:"volume_id"`
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

// API response types
type CloudStorageBlockVolumeBackupAPIResponse struct {
	Id             string                                        `json:"id"`
	Checksum       string                                        `json:"checksum"`
	CreatedAt      string                                        `json:"createdAt"`
	UpdatedAt      string                                        `json:"updatedAt"`
	ResourceStatus string                                        `json:"resourceStatus"`
	CurrentState   *CloudStorageBlockVolumeBackupCurrentState    `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageBlockVolumeBackupTargetSpec      `json:"targetSpec,omitempty"`
}

type CloudStorageBlockVolumeBackupCurrentState struct {
	Location    *CloudStorageBlockVolumeBackupLocation `json:"location,omitempty"`
	Name        string                                 `json:"name,omitempty"`
	Description string                                 `json:"description,omitempty"`
	VolumeId    string                                 `json:"volumeId,omitempty"`
	Size        int64                                  `json:"size,omitempty"`
}

type CloudStorageBlockVolumeBackupTargetSpec struct {
	Location    *CloudStorageBlockVolumeBackupLocation `json:"location,omitempty"`
	Name        string                                 `json:"name,omitempty"`
	Description string                                 `json:"description,omitempty"`
	VolumeId    string                                 `json:"volumeId,omitempty"`
}

type CloudStorageBlockVolumeBackupLocation struct {
	Region string `json:"region,omitempty"`
}

// Create payload
type CloudStorageBlockVolumeBackupCreatePayload struct {
	TargetSpec *CloudStorageBlockVolumeBackupTargetSpec `json:"targetSpec"`
}

// Update payload — only mutable fields (name, description)
type CloudStorageBlockVolumeBackupUpdatePayload struct {
	Checksum   string                                        `json:"checksum"`
	TargetSpec *CloudStorageBlockVolumeBackupUpdateTargetSpec `json:"targetSpec"`
}

type CloudStorageBlockVolumeBackupUpdateTargetSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageBlockVolumeBackupModel) ToCreate() *CloudStorageBlockVolumeBackupCreatePayload {
	return &CloudStorageBlockVolumeBackupCreatePayload{
		TargetSpec: &CloudStorageBlockVolumeBackupTargetSpec{
			Location:    &CloudStorageBlockVolumeBackupLocation{Region: m.Region.ValueString()},
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
			VolumeId:    m.VolumeId.ValueString(),
		},
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageBlockVolumeBackupModel) ToUpdate(checksum string) *CloudStorageBlockVolumeBackupUpdatePayload {
	return &CloudStorageBlockVolumeBackupUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudStorageBlockVolumeBackupUpdateTargetSpec{
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
		},
	}
}

// BackupCurrentStateAttrTypes returns the attribute types for the current_state object
func BackupCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"location":    types.ObjectType{AttrTypes: map[string]attr.Type{"region": ovhtypes.TfStringType{}}},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"volume_id":   ovhtypes.TfStringType{},
		"size":        types.Int64Type,
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageBlockVolumeBackupModel) MergeWith(ctx context.Context, response *CloudStorageBlockVolumeBackupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		locObj, _ := types.ObjectValue(
			map[string]attr.Type{"region": ovhtypes.TfStringType{}},
			map[string]attr.Value{"region": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Location.Region)}},
		)

		currentStateObj, _ := types.ObjectValue(
			BackupCurrentStateAttrTypes(),
			map[string]attr.Value{
				"location":    locObj,
				"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Name)},
				"description": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Description)},
				"volume_id":   ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.VolumeId)},
				"size":        types.Int64Value(response.CurrentState.Size),
			},
		)

		m.CurrentState = currentStateObj
	} else {
		m.CurrentState = types.ObjectNull(BackupCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		if response.TargetSpec.VolumeId != "" {
			m.VolumeId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.VolumeId)}
		}
	}
}
