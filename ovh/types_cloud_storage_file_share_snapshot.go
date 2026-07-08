package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageFileShareSnapshotModel represents the Terraform model for the file share snapshot resource
type CloudStorageFileShareSnapshotModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ShareId     ovhtypes.TfStringValue `tfsdk:"share_id"`
	// Optional — mutable
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
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
type CloudStorageFileShareSnapshotAPIResponse struct {
	Id             string                                     `json:"id"`
	Checksum       string                                     `json:"checksum"`
	CreatedAt      string                                     `json:"createdAt"`
	UpdatedAt      string                                     `json:"updatedAt"`
	ResourceStatus string                                     `json:"resourceStatus"`
	CurrentState   *CloudStorageFileShareSnapshotCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageFileShareSnapshotTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudStorageFileShareSnapshotShareRef struct {
	Id string `json:"id"`
}

type CloudStorageFileShareSnapshotCurrentState struct {
	Name        string                                 `json:"name,omitempty"`
	Description string                                 `json:"description,omitempty"`
	Share       *CloudStorageFileShareSnapshotShareRef `json:"share,omitempty"`
	Size        int64                                  `json:"size,omitempty"`
	Location    *CloudStorageFileShareAPILocation      `json:"location,omitempty"`
}

type CloudStorageFileShareSnapshotTargetSpec struct {
	Name        string                                 `json:"name,omitempty"`
	Description string                                 `json:"description,omitempty"`
	Share       *CloudStorageFileShareSnapshotShareRef `json:"share,omitempty"`
}

// Create payload
type CloudStorageFileShareSnapshotCreatePayload struct {
	TargetSpec *CloudStorageFileShareSnapshotTargetSpec `json:"targetSpec"`
}

// Update payload — only mutable fields (name, description)
type CloudStorageFileShareSnapshotUpdatePayload struct {
	Checksum   string                                         `json:"checksum"`
	TargetSpec *CloudStorageFileShareSnapshotUpdateTargetSpec `json:"targetSpec"`
}

type CloudStorageFileShareSnapshotUpdateTargetSpec struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageFileShareSnapshotModel) ToCreate() *CloudStorageFileShareSnapshotCreatePayload {
	target := &CloudStorageFileShareSnapshotTargetSpec{
		Share: &CloudStorageFileShareSnapshotShareRef{Id: m.ShareId.ValueString()},
	}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		target.Name = m.Name.ValueString()
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	return &CloudStorageFileShareSnapshotCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageFileShareSnapshotModel) ToUpdate(checksum string) *CloudStorageFileShareSnapshotUpdatePayload {
	return &CloudStorageFileShareSnapshotUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudStorageFileShareSnapshotUpdateTargetSpec{
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
		},
	}
}

// FileShareSnapshotCurrentStateAttrTypes returns the attribute types for the current_state object
func FileShareSnapshotCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"share_id":    ovhtypes.TfStringType{},
		"size":        types.Int64Type,
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		}},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageFileShareSnapshotModel) MergeWith(ctx context.Context, response *CloudStorageFileShareSnapshotAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		shareId := ""
		if response.CurrentState.Share != nil {
			shareId = response.CurrentState.Share.Id
		}

		region := ""
		az := ""
		if response.CurrentState.Location != nil {
			region = response.CurrentState.Location.Region
			az = response.CurrentState.Location.AvailabilityZone
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

		currentStateObj, _ := types.ObjectValue(
			FileShareSnapshotCurrentStateAttrTypes(),
			map[string]attr.Value{
				"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Name)},
				"description": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Description)},
				"share_id":    ovhtypes.TfStringValue{StringValue: types.StringValue(shareId)},
				"size":        types.Int64Value(response.CurrentState.Size),
				"location":    locObj,
			},
		)

		m.CurrentState = currentStateObj
	} else {
		m.CurrentState = types.ObjectNull(FileShareSnapshotCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		if response.TargetSpec.Share != nil && response.TargetSpec.Share.Id != "" {
			m.ShareId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Share.Id)}
		}
	}
}
