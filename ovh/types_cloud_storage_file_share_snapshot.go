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
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
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

type CloudStorageFileShareSnapshotCurrentState struct {
	Location     *CloudStorageFileShareSnapshotLocation `json:"location,omitempty"`
	Name         string                                 `json:"name,omitempty"`
	Description  string                                 `json:"description,omitempty"`
	ShareId      string                                 `json:"shareId,omitempty"`
	SnapshotSize int64                                  `json:"snapshotSize,omitempty"`
	ShareSize    int64                                  `json:"shareSize,omitempty"`
	ShareProto   string                                 `json:"shareProto,omitempty"`
}

type CloudStorageFileShareSnapshotTargetSpec struct {
	Location    *CloudStorageFileShareSnapshotLocation `json:"location,omitempty"`
	Name        string                                 `json:"name,omitempty"`
	Description string                                 `json:"description,omitempty"`
	ShareId     string                                 `json:"shareId,omitempty"`
}

type CloudStorageFileShareSnapshotLocation struct {
	Region string `json:"region,omitempty"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageFileShareSnapshotModel) ToCreate() *CloudStorageFileShareSnapshotCreatePayload {
	return &CloudStorageFileShareSnapshotCreatePayload{
		TargetSpec: &CloudStorageFileShareSnapshotTargetSpec{
			Location:    &CloudStorageFileShareSnapshotLocation{Region: m.Region.ValueString()},
			Name:        m.Name.ValueString(),
			Description: m.Description.ValueString(),
			ShareId:     m.ShareId.ValueString(),
		},
	}
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
		"location":      types.ObjectType{AttrTypes: map[string]attr.Type{"region": ovhtypes.TfStringType{}}},
		"name":          ovhtypes.TfStringType{},
		"description":   ovhtypes.TfStringType{},
		"share_id":      ovhtypes.TfStringType{},
		"snapshot_size": types.Int64Type,
		"share_size":    types.Int64Type,
		"share_proto":   ovhtypes.TfStringType{},
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
		locObj, _ := types.ObjectValue(
			map[string]attr.Type{"region": ovhtypes.TfStringType{}},
			map[string]attr.Value{"region": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Location.Region)}},
		)

		currentStateObj, _ := types.ObjectValue(
			FileShareSnapshotCurrentStateAttrTypes(),
			map[string]attr.Value{
				"location":      locObj,
				"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Name)},
				"description":   ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Description)},
				"share_id":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.ShareId)},
				"snapshot_size": types.Int64Value(response.CurrentState.SnapshotSize),
				"share_size":    types.Int64Value(response.CurrentState.ShareSize),
				"share_proto":   ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.ShareProto)},
			},
		)

		m.CurrentState = currentStateObj
	} else {
		m.CurrentState = types.ObjectNull(FileShareSnapshotCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		if response.TargetSpec.ShareId != "" {
			m.ShareId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.ShareId)}
		}
	}
}
