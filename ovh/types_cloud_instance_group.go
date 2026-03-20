package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudInstanceGroupModel represents the Terraform model for the instance group resource
type CloudInstanceGroupModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Policy      ovhtypes.TfStringValue `tfsdk:"policy"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// CloudInstanceGroupDataSourceModel represents the Terraform model for the instance group data source
type CloudInstanceGroupDataSourceModel struct {
	// Required
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name"`
	InstanceGroupId ovhtypes.TfStringValue `tfsdk:"instance_group_id"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
	TargetSpec     types.Object           `tfsdk:"target_spec"`
}

// API response types
type CloudInstanceGroupAPIResponse struct {
	Id             string                             `json:"id"`
	Checksum       string                             `json:"checksum"`
	CreatedAt      string                             `json:"createdAt"`
	UpdatedAt      string                             `json:"updatedAt"`
	ResourceStatus string                             `json:"resourceStatus"`
	CurrentState   *CloudInstanceGroupAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudInstanceGroupAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudInstanceGroupAPICurrentState struct {
	Name     string                           `json:"name,omitempty"`
	Policy   string                           `json:"policy,omitempty"`
	Location *CloudInstanceGroupAPILocation   `json:"location,omitempty"`
	Members  []CloudInstanceGroupAPIMemberRef `json:"members,omitempty"`
}

type CloudInstanceGroupAPITargetSpec struct {
	Name     string                         `json:"name"`
	Policy   string                         `json:"policy"`
	Location *CloudInstanceGroupAPILocation `json:"location"`
}

type CloudInstanceGroupAPILocation struct {
	Region string `json:"region"`
}

type CloudInstanceGroupAPIMemberRef struct {
	Id string `json:"id"`
}

// Create payload
type CloudInstanceGroupCreatePayload struct {
	TargetSpec *CloudInstanceGroupAPITargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudInstanceGroupModel) ToCreate() *CloudInstanceGroupCreatePayload {
	return &CloudInstanceGroupCreatePayload{
		TargetSpec: &CloudInstanceGroupAPITargetSpec{
			Name:   m.Name.ValueString(),
			Policy: m.Policy.ValueString(),
			Location: &CloudInstanceGroupAPILocation{
				Region: m.Region.ValueString(),
			},
		},
	}
}

// InstanceGroupCurrentStateAttrTypes returns the attribute types for the current_state object
func InstanceGroupCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":   ovhtypes.TfStringType{},
		"policy": ovhtypes.TfStringType{},
		"region": ovhtypes.TfStringType{},
		"members": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id": ovhtypes.TfStringType{},
				},
			},
		},
	}
}

// InstanceGroupTargetSpecAttrTypes returns the attribute types for the target_spec object (data source)
func InstanceGroupTargetSpecAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":   ovhtypes.TfStringType{},
		"policy": ovhtypes.TfStringType{},
		"region": ovhtypes.TfStringType{},
	}
}

// buildInstanceGroupCurrentStateObject constructs the current_state object from the API response
func buildInstanceGroupCurrentStateObject(ctx context.Context, state *CloudInstanceGroupAPICurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	// Build members list
	memberAttrTypes := map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}

	var membersVal types.List
	if state.Members != nil {
		memberObjs := make([]attr.Value, len(state.Members))
		for i, member := range state.Members {
			memberObj, _ := types.ObjectValue(
				memberAttrTypes,
				map[string]attr.Value{
					"id": ovhtypes.TfStringValue{StringValue: types.StringValue(member.Id)},
				},
			)
			memberObjs[i] = memberObj
		}
		membersVal, _ = types.ListValue(types.ObjectType{AttrTypes: memberAttrTypes}, memberObjs)
	} else {
		membersVal = types.ListNull(types.ObjectType{AttrTypes: memberAttrTypes})
	}

	obj, _ := types.ObjectValue(
		InstanceGroupCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"policy":  ovhtypes.TfStringValue{StringValue: types.StringValue(state.Policy)},
			"region":  ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"members": membersVal,
		},
	)

	return obj
}

// buildInstanceGroupTargetSpecObject constructs the target_spec object from the API response (data source)
func buildInstanceGroupTargetSpecObject(ctx context.Context, spec *CloudInstanceGroupAPITargetSpec) types.Object {
	region := ""
	if spec.Location != nil {
		region = spec.Location.Region
	}

	obj, _ := types.ObjectValue(
		InstanceGroupTargetSpecAttrTypes(),
		map[string]attr.Value{
			"name":   ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Name)},
			"policy": ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Policy)},
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform model (resource)
func (m *CloudInstanceGroupModel) MergeWith(ctx context.Context, response *CloudInstanceGroupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildInstanceGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceGroupCurrentStateAttrTypes())
	}

	// Update fields from targetSpec if available
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Policy = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Policy)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
	}
}

// MergeWith merges API response data into the data source model
func (m *CloudInstanceGroupDataSourceModel) MergeWith(ctx context.Context, response *CloudInstanceGroupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildInstanceGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceGroupCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.TargetSpec = buildInstanceGroupTargetSpecObject(ctx, response.TargetSpec)
	} else {
		m.TargetSpec = types.ObjectNull(InstanceGroupTargetSpecAttrTypes())
	}
}
