package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudInstanceSnapshotModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

type CloudInstanceSnapshotAPIResponse struct {
	Id             string                             `json:"id"`
	Checksum       string                             `json:"checksum"`
	CreatedAt      string                             `json:"createdAt"`
	UpdatedAt      string                             `json:"updatedAt"`
	ResourceStatus string                             `json:"resourceStatus"`
	CurrentState   *CloudInstanceSnapshotCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudInstanceSnapshotTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudInstanceSnapshotCurrentState struct {
	Instance   *CloudInstanceSnapshotInstanceRef `json:"instance,omitempty"`
	Location   *CloudInstanceSnapshotLocation    `json:"location,omitempty"`
	MinDisk    int64                             `json:"minDisk,omitempty"`
	MinRam     int64                             `json:"minRam,omitempty"`
	Name       string                            `json:"name,omitempty"`
	Size       int64                             `json:"size,omitempty"`
	Status     string                            `json:"status,omitempty"`
	Visibility string                            `json:"visibility,omitempty"`
}

type CloudInstanceSnapshotTargetSpec struct {
	Instance *CloudInstanceSnapshotInstanceRef `json:"instance,omitempty"`
	Location *CloudInstanceSnapshotLocation    `json:"location,omitempty"`
	Name     string                            `json:"name,omitempty"`
}

type CloudInstanceSnapshotInstanceRef struct {
	Id string `json:"id,omitempty"`
}

type CloudInstanceSnapshotLocation struct {
	Region string `json:"region,omitempty"`
}

type CloudInstanceSnapshotCreatePayload struct {
	TargetSpec *CloudInstanceSnapshotTargetSpec `json:"targetSpec"`
}

func (m *CloudInstanceSnapshotModel) ToCreate() *CloudInstanceSnapshotCreatePayload {
	return &CloudInstanceSnapshotCreatePayload{
		TargetSpec: &CloudInstanceSnapshotTargetSpec{
			Instance: &CloudInstanceSnapshotInstanceRef{Id: m.InstanceId.ValueString()},
			Location: &CloudInstanceSnapshotLocation{Region: m.Region.ValueString()},
			Name:     m.Name.ValueString(),
		},
	}
}

func InstanceSnapshotCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instance":   types.ObjectType{AttrTypes: map[string]attr.Type{"id": ovhtypes.TfStringType{}}},
		"location":   types.ObjectType{AttrTypes: map[string]attr.Type{"region": ovhtypes.TfStringType{}}},
		"min_disk":   types.Int64Type,
		"min_ram":    types.Int64Type,
		"name":       ovhtypes.TfStringType{},
		"size":       types.Int64Type,
		"status":     ovhtypes.TfStringType{},
		"visibility": ovhtypes.TfStringType{},
	}
}

func (m *CloudInstanceSnapshotModel) MergeWith(ctx context.Context, response *CloudInstanceSnapshotAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		instanceAttrTypes := map[string]attr.Type{"id": ovhtypes.TfStringType{}}
		var instanceObj types.Object
		if response.CurrentState.Instance != nil {
			instanceObj, _ = types.ObjectValue(
				instanceAttrTypes,
				map[string]attr.Value{"id": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Instance.Id)}},
			)
		} else {
			instanceObj = types.ObjectNull(instanceAttrTypes)
		}

		locationAttrTypes := map[string]attr.Type{"region": ovhtypes.TfStringType{}}
		var locObj types.Object
		if response.CurrentState.Location != nil {
			locObj, _ = types.ObjectValue(
				locationAttrTypes,
				map[string]attr.Value{"region": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Location.Region)}},
			)
		} else {
			locObj = types.ObjectNull(locationAttrTypes)
		}

		currentStateObj, _ := types.ObjectValue(
			InstanceSnapshotCurrentStateAttrTypes(),
			map[string]attr.Value{
				"instance":   instanceObj,
				"location":   locObj,
				"min_disk":   types.Int64Value(response.CurrentState.MinDisk),
				"min_ram":    types.Int64Value(response.CurrentState.MinRam),
				"name":       ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Name)},
				"size":       types.Int64Value(response.CurrentState.Size),
				"status":     ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Status)},
				"visibility": ovhtypes.TfStringValue{StringValue: types.StringValue(response.CurrentState.Visibility)},
			},
		)

		m.CurrentState = currentStateObj
	} else {
		m.CurrentState = types.ObjectNull(InstanceSnapshotCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		if response.TargetSpec.Instance != nil {
			m.InstanceId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Instance.Id)}
		}
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
	}
}
