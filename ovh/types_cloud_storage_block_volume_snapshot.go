package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageBlockVolumeSnapshotModel represents the Terraform model for the block storage snapshot resource
type CloudStorageBlockVolumeSnapshotModel struct {
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
type CloudStorageBlockVolumeSnapshotAPIResponse struct {
	Id             string                                       `json:"id"`
	Checksum       string                                       `json:"checksum"`
	CreatedAt      string                                       `json:"createdAt"`
	UpdatedAt      string                                       `json:"updatedAt"`
	ResourceStatus string                                       `json:"resourceStatus"`
	CurrentState   *CloudStorageBlockVolumeSnapshotCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageBlockVolumeSnapshotTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudStorageBlockVolumeSnapshotCurrentState struct {
	Location    *CloudStorageBlockVolumeSnapshotLocation `json:"location,omitempty"`
	Name        string                                   `json:"name,omitempty"`
	Description string                                   `json:"description,omitempty"`
	VolumeId    string                                   `json:"volumeId,omitempty"`
	Size        int64                                    `json:"size,omitempty"`
}

type CloudStorageBlockVolumeSnapshotTargetSpec struct {
	Location    *CloudStorageBlockVolumeSnapshotLocation `json:"location,omitempty"`
	Name        string                                   `json:"name,omitempty"`
	Description string                                   `json:"description,omitempty"`
	VolumeId    string                                   `json:"volumeId,omitempty"`
}

type CloudStorageBlockVolumeSnapshotLocation struct {
	Region string `json:"region,omitempty"`
}

// Create payload
type CloudStorageBlockVolumeSnapshotCreatePayload struct {
	TargetSpec *CloudStorageBlockVolumeSnapshotTargetSpec `json:"targetSpec"`
}

// Update payload — only mutable fields (name, description)
type CloudStorageBlockVolumeSnapshotUpdatePayload struct {
	Checksum   string                                           `json:"checksum"`
	TargetSpec *CloudStorageBlockVolumeSnapshotUpdateTargetSpec `json:"targetSpec"`
}

type CloudStorageBlockVolumeSnapshotUpdateTargetSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageBlockVolumeSnapshotModel) ToCreate() *CloudStorageBlockVolumeSnapshotCreatePayload {
	return &CloudStorageBlockVolumeSnapshotCreatePayload{
		TargetSpec: &CloudStorageBlockVolumeSnapshotTargetSpec{
			Location:    &CloudStorageBlockVolumeSnapshotLocation{Region: m.Region.ValueString()},
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
			VolumeId:    m.VolumeId.ValueString(),
		},
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageBlockVolumeSnapshotModel) ToUpdate(checksum string) *CloudStorageBlockVolumeSnapshotUpdatePayload {
	return &CloudStorageBlockVolumeSnapshotUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudStorageBlockVolumeSnapshotUpdateTargetSpec{
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
		},
	}
}

// SnapshotCurrentStateAttrTypes returns the attribute types for the current_state object
func SnapshotCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"location":    types.ObjectType{AttrTypes: map[string]attr.Type{"region": ovhtypes.TfStringType{}}},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"volume_id":   ovhtypes.TfStringType{},
		"size":        types.Int64Type,
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageBlockVolumeSnapshotModel) MergeWith(ctx context.Context, response *CloudStorageBlockVolumeSnapshotAPIResponse) {
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
			SnapshotCurrentStateAttrTypes(),
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
		m.CurrentState = types.ObjectNull(SnapshotCurrentStateAttrTypes())
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
