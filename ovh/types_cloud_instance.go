package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudInstanceModel represents the Terraform model for the instance resource
type CloudInstanceModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	FlavorId    ovhtypes.TfStringValue `tfsdk:"flavor_id"`
	ImageId     ovhtypes.TfStringValue `tfsdk:"image_id"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional
	AvailabilityZone ovhtypes.TfStringValue                             `tfsdk:"availability_zone"`
	VolumeIds        ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"volume_ids"`
	Networks         types.List                                         `tfsdk:"networks"`
	SSHKeyName       ovhtypes.TfStringValue                             `tfsdk:"ssh_key_name"`
	GroupId          ovhtypes.TfStringValue                             `tfsdk:"group_id"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types
type CloudInstanceAPIResponse struct {
	Id             string                        `json:"id"`
	Checksum       string                        `json:"checksum"`
	CreatedAt      string                        `json:"createdAt"`
	UpdatedAt      string                        `json:"updatedAt"`
	ResourceStatus string                        `json:"resourceStatus"`
	CurrentState   *CloudInstanceAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudInstanceAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudInstanceAPICurrentState struct {
	Flavor         *CloudInstanceAPIFlavor   `json:"flavor,omitempty"`
	Image          *CloudInstanceAPIImage    `json:"image,omitempty"`
	Name           string                    `json:"name,omitempty"`
	HostId         string                    `json:"hostId,omitempty"`
	SSHKeyName     string                    `json:"sshKeyName,omitempty"`
	ProjectId      string                    `json:"projectId,omitempty"`
	UserId         string                    `json:"userId,omitempty"`
	Networks       []CloudInstanceAPINetwork `json:"networks,omitempty"`
	Volumes        []CloudInstanceAPIVolume  `json:"volumes,omitempty"`
	SecurityGroups []string                  `json:"securityGroups,omitempty"`
	Group          *CloudInstanceAPIGroupRef `json:"group,omitempty"`
}

type CloudInstanceAPIGroupRef struct {
	Id string `json:"id"`
}

type CloudInstanceAPIFlavor struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Vcpus int64  `json:"vcpus,omitempty"`
	Ram   int64  `json:"ram,omitempty"`
	Disk  int64  `json:"disk,omitempty"`
}

type CloudInstanceAPIImage struct {
	Id     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type CloudInstanceAPINetwork struct {
	Id           string                       `json:"id"`
	Name         string                       `json:"name,omitempty"`
	Public       *bool                        `json:"public,omitempty"`
	SubnetId     string                       `json:"subnetId,omitempty"`
	GatewayId    string                       `json:"gatewayId,omitempty"`
	FloatingIpId string                       `json:"floatingIpId,omitempty"`
	Addresses    []CloudInstanceAPINetAddress `json:"addresses,omitempty"`
}

type CloudInstanceAPINetAddress struct {
	Ip      string `json:"ip,omitempty"`
	Mac     string `json:"mac,omitempty"`
	Type    string `json:"type,omitempty"`
	Version int64  `json:"version,omitempty"`
}

type CloudInstanceAPIVolume struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type CloudInstanceAPITargetSpec struct {
	Flavor     *CloudInstanceAPIFlavorRef   `json:"flavor"`
	Image      *CloudInstanceAPIImageRef    `json:"image"`
	Location   *CloudInstanceAPILocation    `json:"location,omitempty"`
	Name       string                       `json:"name"`
	SSHKeyName string                       `json:"sshKeyName,omitempty"`
	Networks   []CloudInstanceAPINetworkRef `json:"networks"`
	Volumes    []CloudInstanceAPIVolumeRef  `json:"volumes"`
	Group      *CloudInstanceAPIGroupRef    `json:"group,omitempty"`
}

type CloudInstanceAPIFlavorRef struct {
	Id string `json:"id"`
}

type CloudInstanceAPIImageRef struct {
	Id string `json:"id"`
}

type CloudInstanceAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudInstanceAPINetworkRef struct {
	Id           *string `json:"id,omitempty"`
	Public       *bool   `json:"public,omitempty"`
	SubnetId     *string `json:"subnetId,omitempty"`
	FloatingIpId *string `json:"floatingIpId,omitempty"`
}

type CloudInstanceAPIVolumeRef struct {
	Id string `json:"id"`
}

// Create payload
type CloudInstanceCreatePayload struct {
	TargetSpec *CloudInstanceAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudInstanceUpdatePayload struct {
	Checksum   string                      `json:"checksum"`
	TargetSpec *CloudInstanceAPITargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudInstanceModel) ToCreate() *CloudInstanceCreatePayload {
	targetSpec := &CloudInstanceAPITargetSpec{
		Flavor: &CloudInstanceAPIFlavorRef{
			Id: m.FlavorId.ValueString(),
		},
		Image: &CloudInstanceAPIImageRef{
			Id: m.ImageId.ValueString(),
		},
		Location: &CloudInstanceAPILocation{
			Region: m.Region.ValueString(),
		},
		Name:       m.Name.ValueString(),
		SSHKeyName: m.SSHKeyName.ValueString(),
	}

	// Set availability zone if provided
	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	// Set group if provided (immutable, only on create)
	if !m.GroupId.IsNull() && !m.GroupId.IsUnknown() {
		targetSpec.Group = &CloudInstanceAPIGroupRef{Id: m.GroupId.ValueString()}
	}

	// Build networks from structured networks if provided
	if !m.Networks.IsNull() && !m.Networks.IsUnknown() && len(m.Networks.Elements()) > 0 {
		networks := make([]CloudInstanceAPINetworkRef, 0, len(m.Networks.Elements()))
		for _, n := range m.Networks.Elements() {
			obj, ok := n.(basetypes.ObjectValue)
			if !ok {
				continue
			}

			attrs := obj.Attributes()

			var idPtr *string
			if v, ok := attrs["id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				id := v.ValueString()
				idPtr = &id
			}

			var publicPtr *bool
			if v, ok := attrs["public"].(basetypes.BoolValue); ok && !v.IsNull() && !v.IsUnknown() {
				val := v.ValueBool()
				publicPtr = &val
			}

			var subnetIdPtr *string
			if v, ok := attrs["subnet_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				val := v.ValueString()
				subnetIdPtr = &val
			}

			var floatingIpIdPtr *string
			if v, ok := attrs["floating_ip_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				val := v.ValueString()
				floatingIpIdPtr = &val
			}

			if idPtr != nil || publicPtr != nil || subnetIdPtr != nil || floatingIpIdPtr != nil {
				networks = append(networks, CloudInstanceAPINetworkRef{
					Id:           idPtr,
					Public:       publicPtr,
					SubnetId:     subnetIdPtr,
					FloatingIpId: floatingIpIdPtr,
				})
			}
		}
		if len(networks) > 0 {
			targetSpec.Networks = networks
		}
	}

	// Add volumes if specified
	if !m.VolumeIds.IsNull() && !m.VolumeIds.IsUnknown() {
		volumes := make([]CloudInstanceAPIVolumeRef, 0)
		for _, volId := range m.VolumeIds.Elements() {
			if strVal, ok := volId.(ovhtypes.TfStringValue); ok {
				volumes = append(volumes, CloudInstanceAPIVolumeRef{
					Id: strVal.ValueString(),
				})
			}
		}
		targetSpec.Volumes = volumes
	}

	return &CloudInstanceCreatePayload{
		TargetSpec: targetSpec,
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudInstanceModel) ToUpdate(checksum string) *CloudInstanceUpdatePayload {
	targetSpec := &CloudInstanceAPITargetSpec{
		Flavor: &CloudInstanceAPIFlavorRef{
			Id: m.FlavorId.ValueString(),
		},
		Image: &CloudInstanceAPIImageRef{
			Id: m.ImageId.ValueString(),
		},
		Name: m.Name.ValueString(),
	}

	// Build networks from structured networks if provided
	if !m.Networks.IsNull() && !m.Networks.IsUnknown() && len(m.Networks.Elements()) > 0 {
		networks := make([]CloudInstanceAPINetworkRef, 0, len(m.Networks.Elements()))
		for _, n := range m.Networks.Elements() {
			obj, ok := n.(basetypes.ObjectValue)
			if !ok {
				continue
			}

			attrs := obj.Attributes()

			var idPtr *string
			if v, ok := attrs["id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				id := v.ValueString()
				idPtr = &id
			}

			var publicPtr *bool
			if v, ok := attrs["public"].(basetypes.BoolValue); ok && !v.IsNull() && !v.IsUnknown() {
				b := v.ValueBool()
				publicPtr = &b
			}

			var subnetIdPtr *string
			if v, ok := attrs["subnet_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				s := v.ValueString()
				subnetIdPtr = &s
			}

			var floatingIpIdPtr *string
			if v, ok := attrs["floating_ip_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
				s := v.ValueString()
				floatingIpIdPtr = &s
			}

			if idPtr != nil || publicPtr != nil || subnetIdPtr != nil || floatingIpIdPtr != nil {
				networks = append(networks, CloudInstanceAPINetworkRef{
					Id:           idPtr,
					Public:       publicPtr,
					SubnetId:     subnetIdPtr,
					FloatingIpId: floatingIpIdPtr,
				})
			}
		}
		if len(networks) > 0 {
			targetSpec.Networks = networks
		}
	}

	// Add volumes if specified
	if !m.VolumeIds.IsNull() && !m.VolumeIds.IsUnknown() {
		volumes := make([]CloudInstanceAPIVolumeRef, 0)
		for _, volId := range m.VolumeIds.Elements() {
			if strVal, ok := volId.(ovhtypes.TfStringValue); ok {
				volumes = append(volumes, CloudInstanceAPIVolumeRef{
					Id: strVal.ValueString(),
				})
			}
		}
		targetSpec.Volumes = volumes
	}

	return &CloudInstanceUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// CurrentStateAttrTypes returns the attribute types for the current_state object
func CurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"flavor": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":    ovhtypes.TfStringType{},
				"name":  ovhtypes.TfStringType{},
				"vcpus": types.Int64Type,
				"ram":   types.Int64Type,
				"disk":  types.Int64Type,
			},
		},
		"image": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":     ovhtypes.TfStringType{},
				"name":   ovhtypes.TfStringType{},
				"status": ovhtypes.TfStringType{},
			},
		},
		"name":         ovhtypes.TfStringType{},
		"host_id":      ovhtypes.TfStringType{},
		"ssh_key_name": ovhtypes.TfStringType{},
		"project_id":   ovhtypes.TfStringType{},
		"user_id":      ovhtypes.TfStringType{},
		"networks": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":             ovhtypes.TfStringType{},
					"public":         types.BoolType,
					"subnet_id":      ovhtypes.TfStringType{},
					"gateway_id":     ovhtypes.TfStringType{},
					"floating_ip_id": ovhtypes.TfStringType{},
					"addresses": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"ip":      ovhtypes.TfStringType{},
								"mac":     ovhtypes.TfStringType{},
								"type":    ovhtypes.TfStringType{},
								"version": types.Int64Type,
							},
						},
					},
				},
			},
		},
		"volumes": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   ovhtypes.TfStringType{},
					"name": ovhtypes.TfStringType{},
					"size": types.Int64Type,
				},
			},
		},
		"security_groups": types.ListType{
			ElemType: types.StringType,
		},
		"group_id": ovhtypes.TfStringType{},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudInstanceModel) MergeWith(ctx context.Context, response *CloudInstanceAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(CurrentStateAttrTypes())
	}

	// Update region from targetSpec if available
	if response.TargetSpec != nil && response.TargetSpec.Location != nil {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
	}

	// Update group_id from targetSpec if available
	if response.TargetSpec != nil && response.TargetSpec.Group != nil {
		m.GroupId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Group.Id)}
	}
}

func buildCurrentStateObject(ctx context.Context, state *CloudInstanceAPICurrentState) basetypes.ObjectValue {
	// Build flavor object
	var flavorObj basetypes.ObjectValue
	if state.Flavor != nil {
		flavorObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"id":    ovhtypes.TfStringType{},
				"name":  ovhtypes.TfStringType{},
				"vcpus": types.Int64Type,
				"ram":   types.Int64Type,
				"disk":  types.Int64Type,
			},
			map[string]attr.Value{
				"id":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.Flavor.Id)},
				"name":  ovhtypes.TfStringValue{StringValue: types.StringValue(state.Flavor.Name)},
				"vcpus": types.Int64Value(state.Flavor.Vcpus),
				"ram":   types.Int64Value(state.Flavor.Ram),
				"disk":  types.Int64Value(state.Flavor.Disk),
			},
		)
	} else {
		flavorObj = types.ObjectNull(map[string]attr.Type{
			"id":    ovhtypes.TfStringType{},
			"name":  ovhtypes.TfStringType{},
			"vcpus": types.Int64Type,
			"ram":   types.Int64Type,
			"disk":  types.Int64Type,
		})
	}

	// Build image object
	var imageObj basetypes.ObjectValue
	if state.Image != nil {
		imageObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"id":     ovhtypes.TfStringType{},
				"name":   ovhtypes.TfStringType{},
				"status": ovhtypes.TfStringType{},
			},
			map[string]attr.Value{
				"id":     ovhtypes.TfStringValue{StringValue: types.StringValue(state.Image.Id)},
				"name":   ovhtypes.TfStringValue{StringValue: types.StringValue(state.Image.Name)},
				"status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Image.Status)},
			},
		)
	} else {
		imageObj = types.ObjectNull(map[string]attr.Type{
			"id":     ovhtypes.TfStringType{},
			"name":   ovhtypes.TfStringType{},
			"status": ovhtypes.TfStringType{},
		})
	}

	// Build networks list
	networkAttrTypes := map[string]attr.Type{
		"id":             ovhtypes.TfStringType{},
		"public":         types.BoolType,
		"subnet_id":      ovhtypes.TfStringType{},
		"gateway_id":     ovhtypes.TfStringType{},
		"floating_ip_id": ovhtypes.TfStringType{},
		"addresses": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"ip":      ovhtypes.TfStringType{},
					"mac":     ovhtypes.TfStringType{},
					"type":    ovhtypes.TfStringType{},
					"version": types.Int64Type,
				},
			},
		},
	}

	var networksVal basetypes.ListValue
	if state.Networks != nil {
		networkObjs := make([]attr.Value, len(state.Networks))
		for i, net := range state.Networks {
			// Build addresses list
			addressObjs := make([]attr.Value, len(net.Addresses))
			for j, addr := range net.Addresses {
				addrObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"ip":      ovhtypes.TfStringType{},
						"mac":     ovhtypes.TfStringType{},
						"type":    ovhtypes.TfStringType{},
						"version": types.Int64Type,
					},
					map[string]attr.Value{
						"ip":      ovhtypes.TfStringValue{StringValue: types.StringValue(addr.Ip)},
						"mac":     ovhtypes.TfStringValue{StringValue: types.StringValue(addr.Mac)},
						"type":    ovhtypes.TfStringValue{StringValue: types.StringValue(addr.Type)},
						"version": types.Int64Value(addr.Version),
					},
				)
				addressObjs[j] = addrObj
			}

			addressesList, _ := types.ListValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"ip":      ovhtypes.TfStringType{},
						"mac":     ovhtypes.TfStringType{},
						"type":    ovhtypes.TfStringType{},
						"version": types.Int64Type,
					},
				},
				addressObjs,
			)

			var publicVal attr.Value
			if net.Public != nil {
				publicVal = types.BoolValue(*net.Public)
			} else {
				publicVal = types.BoolNull()
			}

			netObj, _ := types.ObjectValue(
				networkAttrTypes,
				map[string]attr.Value{
					"id":             ovhtypes.TfStringValue{StringValue: types.StringValue(net.Id)},
					"public":         publicVal,
					"subnet_id":      ovhtypes.TfStringValue{StringValue: types.StringValue(net.SubnetId)},
					"gateway_id":     ovhtypes.TfStringValue{StringValue: types.StringValue(net.GatewayId)},
					"floating_ip_id": ovhtypes.TfStringValue{StringValue: types.StringValue(net.FloatingIpId)},
					"addresses":      addressesList,
				},
			)
			networkObjs[i] = netObj
		}
		networksVal, _ = types.ListValue(types.ObjectType{AttrTypes: networkAttrTypes}, networkObjs)
	} else {
		networksVal = types.ListNull(types.ObjectType{AttrTypes: networkAttrTypes})
	}

	// Build volumes list
	volumeAttrTypes := map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"name": ovhtypes.TfStringType{},
		"size": types.Int64Type,
	}

	var volumesVal basetypes.ListValue
	if state.Volumes != nil {
		volumeObjs := make([]attr.Value, len(state.Volumes))
		for i, vol := range state.Volumes {
			volObj, _ := types.ObjectValue(
				volumeAttrTypes,
				map[string]attr.Value{
					"id":   ovhtypes.TfStringValue{StringValue: types.StringValue(vol.Id)},
					"name": ovhtypes.TfStringValue{StringValue: types.StringValue(vol.Name)},
					"size": types.Int64Value(vol.Size),
				},
			)
			volumeObjs[i] = volObj
		}
		volumesVal, _ = types.ListValue(types.ObjectType{AttrTypes: volumeAttrTypes}, volumeObjs)
	} else {
		volumesVal = types.ListNull(types.ObjectType{AttrTypes: volumeAttrTypes})
	}

	// Build security groups list
	var securityGroupsVal basetypes.ListValue
	if state.SecurityGroups != nil {
		sgVals := make([]attr.Value, len(state.SecurityGroups))
		for i, sg := range state.SecurityGroups {
			sgVals[i] = types.StringValue(sg)
		}
		securityGroupsVal, _ = types.ListValue(types.StringType, sgVals)
	} else {
		securityGroupsVal = types.ListNull(types.StringType)
	}

	// Build group_id value
	var groupIdVal attr.Value
	if state.Group != nil && state.Group.Id != "" {
		groupIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(state.Group.Id)}
	} else {
		groupIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue("")}
	}

	// Build the complete current_state object
	currentStateObj, _ := types.ObjectValue(
		CurrentStateAttrTypes(),
		map[string]attr.Value{
			"flavor":          flavorObj,
			"image":           imageObj,
			"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"host_id":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.HostId)},
			"ssh_key_name":    ovhtypes.TfStringValue{StringValue: types.StringValue(state.SSHKeyName)},
			"project_id":      ovhtypes.TfStringValue{StringValue: types.StringValue(state.ProjectId)},
			"user_id":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.UserId)},
			"networks":        networksVal,
			"volumes":         volumesVal,
			"security_groups": securityGroupsVal,
			"group_id":        groupIdVal,
		},
	)

	return currentStateObj
}
