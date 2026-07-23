package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudInstanceGroupModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Policy      ovhtypes.TfStringValue `tfsdk:"policy"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

type CloudInstanceGroupAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudInstanceGroupAPIMemberRef struct {
	Id string `json:"id"`
}

type CloudInstanceGroupAPITargetSpec struct {
	Name     string                         `json:"name"`
	Policy   string                         `json:"policy"`
	Location *CloudInstanceGroupAPILocation `json:"location,omitempty"`
}

type CloudInstanceGroupAPICurrentState struct {
	Name     string                           `json:"name,omitempty"`
	Policy   string                           `json:"policy,omitempty"`
	Location *CloudInstanceGroupAPILocation   `json:"location,omitempty"`
	Members  []CloudInstanceGroupAPIMemberRef `json:"members,omitempty"`
}

type CloudInstanceGroupAPIResponse struct {
	Id             string                             `json:"id"`
	Checksum       string                             `json:"checksum"`
	CreatedAt      string                             `json:"createdAt"`
	UpdatedAt      string                             `json:"updatedAt"`
	ResourceStatus string                             `json:"resourceStatus"`
	TargetSpec     *CloudInstanceGroupAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudInstanceGroupAPICurrentState `json:"currentState,omitempty"`
}

type CloudInstanceGroupCreatePayload struct {
	TargetSpec *CloudInstanceGroupAPITargetSpec `json:"targetSpec"`
}

func instanceGroupLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

func instanceGroupMemberAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"id": ovhtypes.TfStringType{}}
}

func InstanceGroupCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":     ovhtypes.TfStringType{},
		"policy":   ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: instanceGroupLocationAttrTypes()},
		"members":  types.ListType{ElemType: types.ObjectType{AttrTypes: instanceGroupMemberAttrTypes()}},
	}
}

func (m *CloudInstanceGroupModel) ToCreate() *CloudInstanceGroupCreatePayload {
	return &CloudInstanceGroupCreatePayload{
		TargetSpec: &CloudInstanceGroupAPITargetSpec{
			Name:     m.Name.ValueString(),
			Policy:   m.Policy.ValueString(),
			Location: &CloudInstanceGroupAPILocation{Region: m.Region.ValueString()},
		},
	}
}

func buildInstanceGroupCurrentStateObject(ctx context.Context, s *CloudInstanceGroupAPICurrentState) basetypes.ObjectValue {
	sv := func(v string) ovhtypes.TfStringValue {
		return ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
	}

	locVal := types.ObjectNull(instanceGroupLocationAttrTypes())
	if s.Location != nil {
		locVal, _ = types.ObjectValue(instanceGroupLocationAttrTypes(), map[string]attr.Value{
			"region":            sv(s.Location.Region),
			"availability_zone": sv(s.Location.AvailabilityZone),
		})
	}

	memberObjType := types.ObjectType{AttrTypes: instanceGroupMemberAttrTypes()}
	var membersVal basetypes.ListValue
	if s.Members == nil {
		membersVal = types.ListNull(memberObjType)
	} else {
		items := make([]attr.Value, 0, len(s.Members))
		for _, mem := range s.Members {
			obj, _ := types.ObjectValue(instanceGroupMemberAttrTypes(), map[string]attr.Value{"id": sv(mem.Id)})
			items = append(items, obj)
		}
		membersVal = types.ListValueMust(memberObjType, items)
	}

	obj, _ := types.ObjectValue(InstanceGroupCurrentStateAttrTypes(), map[string]attr.Value{
		"name":     sv(s.Name),
		"policy":   sv(s.Policy),
		"location": locVal,
		"members":  membersVal,
	})
	return obj
}

func (m *CloudInstanceGroupModel) MergeWith(ctx context.Context, response *CloudInstanceGroupAPIResponse) {
	sv := func(v string) ovhtypes.TfStringValue {
		return ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
	}
	m.Id = sv(response.Id)
	m.Checksum = sv(response.Checksum)
	m.CreatedAt = sv(response.CreatedAt)
	m.UpdatedAt = sv(response.UpdatedAt)
	m.ResourceStatus = sv(response.ResourceStatus)

	if response.CurrentState != nil {
		m.CurrentState = buildInstanceGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceGroupCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = sv(response.TargetSpec.Name)
		m.Policy = sv(response.TargetSpec.Policy)
		if response.TargetSpec.Location != nil {
			m.Region = sv(response.TargetSpec.Location.Region)
		}
	}
}
