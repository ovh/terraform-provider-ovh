package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageFileShareAclModel struct {
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

type CloudStorageFileShareAclAPIResponse struct {
	Id             string                                `json:"id"`
	Checksum       string                                `json:"checksum"`
	CreatedAt      string                                `json:"createdAt"`
	UpdatedAt      string                                `json:"updatedAt"`
	ResourceStatus string                                `json:"resourceStatus"`
	CurrentState   *CloudStorageFileShareAclCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageFileShareAclTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudStorageFileShareAclCurrentState struct {
	AccessTo    string `json:"accessTo,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
	State       string `json:"state,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type CloudStorageFileShareAclTargetSpec struct {
	AccessTo    string `json:"accessTo,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

type CloudStorageFileShareAclCreatePayload struct {
	TargetSpec *CloudStorageFileShareAclTargetSpec `json:"targetSpec"`
}

func (m *CloudStorageFileShareAclModel) ToCreate() *CloudStorageFileShareAclCreatePayload {
	return &CloudStorageFileShareAclCreatePayload{
		TargetSpec: &CloudStorageFileShareAclTargetSpec{
			AccessTo:    m.AccessTo.ValueString(),
			AccessLevel: m.AccessLevel.ValueString(),
		},
	}
}

func StorageFileShareAclCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_to":    ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
		"state":        ovhtypes.TfStringType{},
		"created_at":   ovhtypes.TfStringType{},
	}
}

func buildStorageFileShareAclCurrentStateObject(state *CloudStorageFileShareAclCurrentState) types.Object {
	obj, _ := types.ObjectValue(
		StorageFileShareAclCurrentStateAttrTypes(),
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
func (m *CloudStorageFileShareAclModel) MergeWith(ctx context.Context, response *CloudStorageFileShareAclAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildStorageFileShareAclCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(StorageFileShareAclCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.AccessTo = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.AccessTo)}
		m.AccessLevel = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.AccessLevel)}
	}
}
