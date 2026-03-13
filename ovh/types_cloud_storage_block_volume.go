package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageBlockVolumeModel represents the Terraform model for the block storage resource
type CloudStorageBlockVolumeModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Size        types.Int64            `tfsdk:"size"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	VolumeType  ovhtypes.TfStringValue `tfsdk:"volume_type"`
	Bootable    types.Bool             `tfsdk:"bootable"`
	CreateFrom  types.Object           `tfsdk:"create_from"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudStorageBlockVolumeAPIResponse struct {
	Id             string                               `json:"id"`
	Checksum       string                               `json:"checksum"`
	CreatedAt      string                               `json:"createdAt"`
	UpdatedAt      string                               `json:"updatedAt"`
	ResourceStatus string                               `json:"resourceStatus"`
	CurrentState   *CloudStorageBlockVolumeCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageBlockVolumeTarget       `json:"targetSpec,omitempty"`
}

type CloudStorageBlockVolumeCurrentState struct {
	Location   *CloudStorageBlockVolumeLocation `json:"location,omitempty"`
	Name       string                           `json:"name,omitempty"`
	Size       int64                            `json:"size,omitempty"`
	VolumeType string                           `json:"volumeType,omitempty"`
	Bootable   *bool                            `json:"bootable,omitempty"`
	Status     string                           `json:"status,omitempty"`
}

type CloudStorageBlockVolumeCreateFrom struct {
	BackupID string `json:"backupId,omitempty"`
}

type CloudStorageBlockVolumeTarget struct {
	Location   *CloudStorageBlockVolumeLocation   `json:"location,omitempty"`
	Name       string                             `json:"name,omitempty"`
	Size       int64                              `json:"size,omitempty"`
	VolumeType string                             `json:"volumeType,omitempty"`
	Bootable   *bool                              `json:"bootable,omitempty"`
	CreateFrom *CloudStorageBlockVolumeCreateFrom `json:"createFrom,omitempty"`
}

type CloudStorageBlockVolumeLocation struct {
	Region string `json:"region,omitempty"`
}

// Create payload
type CloudStorageBlockVolumeCreatePayload struct {
	TargetSpec *CloudStorageBlockVolumeTarget `json:"targetSpec"`
}

// Update payload
type CloudStorageBlockVolumeUpdatePayload struct {
	Checksum   string                         `json:"checksum"`
	TargetSpec *CloudStorageBlockVolumeTarget `json:"targetSpec"`
}

// CreateFromAttrTypes returns the attribute types for the create_from object
func CreateFromAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_id": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageBlockVolumeModel) ToCreate() *CloudStorageBlockVolumeCreatePayload {
	target := &CloudStorageBlockVolumeTarget{
		Location:   &CloudStorageBlockVolumeLocation{Region: m.Region.ValueString()},
		Name:       m.Name.ValueString(),
		Size:       m.Size.ValueInt64(),
		VolumeType: m.VolumeType.ValueString(),
	}

	if !m.Bootable.IsNull() && !m.Bootable.IsUnknown() {
		b := m.Bootable.ValueBool()
		target.Bootable = &b
	}

	if !m.CreateFrom.IsNull() && !m.CreateFrom.IsUnknown() {
		attrs := m.CreateFrom.Attributes()
		if backupIDVal, ok := attrs["backup_id"]; ok {
			if strVal, ok := backupIDVal.(ovhtypes.TfStringValue); ok && !strVal.IsNull() && !strVal.IsUnknown() && strVal.ValueString() != "" {
				target.CreateFrom = &CloudStorageBlockVolumeCreateFrom{
					BackupID: strVal.ValueString(),
				}
			}
		}
	}

	return &CloudStorageBlockVolumeCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageBlockVolumeModel) ToUpdate(checksum string) *CloudStorageBlockVolumeUpdatePayload {
	target := &CloudStorageBlockVolumeTarget{
		Name:       m.Name.ValueString(),
		Size:       m.Size.ValueInt64(),
		VolumeType: m.VolumeType.ValueString(),
	}

	if !m.Bootable.IsNull() && !m.Bootable.IsUnknown() {
		b := m.Bootable.ValueBool()
		target.Bootable = &b
	}

	return &CloudStorageBlockVolumeUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// BlockVolumeCurrentStateAttrTypes returns the attribute types for the current_state object
func BlockVolumeCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"location":    types.ObjectType{AttrTypes: map[string]attr.Type{"region": ovhtypes.TfStringType{}}},
		"name":        ovhtypes.TfStringType{},
		"size":        types.Int64Type,
		"volume_type": ovhtypes.TfStringType{},
		"bootable":    types.BoolType,
		"status":      ovhtypes.TfStringType{},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageBlockVolumeModel) MergeWith(ctx context.Context, response *CloudStorageBlockVolumeAPIResponse) {
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

		bootableVal := types.BoolValue(false)
		if response.CurrentState.Bootable != nil {
			bootableVal = types.BoolValue(*response.CurrentState.Bootable)
		}

		currentStateObj, _ := types.ObjectValue(
			BlockVolumeCurrentStateAttrTypes(),
			map[string]attr.Value{
				"location":    locObj,
				"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Name)},
				"size":        types.Int64Value(response.CurrentState.Size),
				"volume_type": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.VolumeType)},
				"bootable":    bootableVal,
				"status":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Status)},
			},
		)

		m.CurrentState = currentStateObj
	} else {
		m.CurrentState = types.ObjectNull(BlockVolumeCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Size = types.Int64Value(response.TargetSpec.Size)
		if response.TargetSpec.VolumeType != "" {
			m.VolumeType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.VolumeType)}
		}
		if response.TargetSpec.Bootable != nil {
			m.Bootable = types.BoolValue(*response.TargetSpec.Bootable)
		} else if m.Bootable.IsUnknown() {
			m.Bootable = types.BoolNull()
		}

		// Preserve create_from in state if it was set
		if response.TargetSpec.CreateFrom != nil && response.TargetSpec.CreateFrom.BackupID != "" {
			createFromObj, _ := types.ObjectValue(
				CreateFromAttrTypes(),
				map[string]attr.Value{
					"backup_id": ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.CreateFrom.BackupID)},
				},
			)
			m.CreateFrom = createFromObj
		} else if m.CreateFrom.IsNull() || m.CreateFrom.IsUnknown() {
			m.CreateFrom = types.ObjectNull(CreateFromAttrTypes())
		}
	}
}
