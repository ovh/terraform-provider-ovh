package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudFileStorageAclModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ShareId     ovhtypes.TfStringValue `tfsdk:"share_id"`
	AccessTo    ovhtypes.TfStringValue `tfsdk:"access_to"`

	// Required — mutable
	AccessLevel ovhtypes.TfStringValue `tfsdk:"access_level"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

type CloudFileStorageAclAPIResponse struct {
	Id             string                           `json:"id"`
	Checksum       string                           `json:"checksum"`
	CreatedAt      string                           `json:"createdAt"`
	UpdatedAt      string                           `json:"updatedAt"`
	ResourceStatus string                           `json:"resourceStatus"`
	CurrentState   *CloudFileStorageAclCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudFileStorageAclTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudFileStorageAclCurrentState struct {
	AccessTo    string `json:"accessTo,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
	State       string `json:"state,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type CloudFileStorageAclTargetSpec struct {
	AccessTo    string `json:"accessTo,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

type CloudFileStorageAclUpdateTargetSpec struct {
	AccessLevel string `json:"accessLevel,omitempty"`
}

type CloudFileStorageAclCreatePayload struct {
	TargetSpec *CloudFileStorageAclTargetSpec `json:"targetSpec"`
}

// Update payload — only accessLevel is mutable
type CloudFileStorageAclUpdatePayload struct {
	Checksum   string                               `json:"checksum"`
	TargetSpec *CloudFileStorageAclUpdateTargetSpec `json:"targetSpec"`
}

func (m *CloudFileStorageAclModel) ToCreate() *CloudFileStorageAclCreatePayload {
	return &CloudFileStorageAclCreatePayload{
		TargetSpec: &CloudFileStorageAclTargetSpec{
			AccessTo:    m.AccessTo.ValueString(),
			AccessLevel: m.AccessLevel.ValueString(),
		},
	}
}

func (m *CloudFileStorageAclModel) ToUpdate(checksum string) *CloudFileStorageAclUpdatePayload {
	return &CloudFileStorageAclUpdatePayload{
		Checksum: checksum,
		TargetSpec: &CloudFileStorageAclUpdateTargetSpec{
			AccessLevel: m.AccessLevel.ValueString(),
		},
	}
}

func FileStorageAclCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_to":    ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
		"state":        ovhtypes.TfStringType{},
		"created_at":   ovhtypes.TfStringType{},
	}
}

func buildFileStorageAclCurrentStateObject(state *CloudFileStorageAclCurrentState) types.Object {
	obj, _ := types.ObjectValue(
		FileStorageAclCurrentStateAttrTypes(),
		map[string]attr.Value{
			"access_to":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.AccessTo)},
			"access_level": ovhtypes.TfStringValue{StringValue: types.StringValue(state.AccessLevel)},
			"state":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.State)},
			"created_at":   ovhtypes.TfStringValue{StringValue: types.StringValue(state.CreatedAt)},
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model. The share_id
// (parent) is never part of the API response — it stays whatever is already
// set in the model from plan/state.
func (m *CloudFileStorageAclModel) MergeWith(ctx context.Context, response *CloudFileStorageAclAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildFileStorageAclCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(FileStorageAclCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.AccessTo = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.AccessTo)}
		m.AccessLevel = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.AccessLevel)}
	}
}
