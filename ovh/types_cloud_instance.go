package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudInstanceModel is the Terraform model for the ovh_cloud_instance resource.
type CloudInstanceModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`
	SSHKeyName       ovhtypes.TfStringValue `tfsdk:"ssh_key_name"`
	GroupId          ovhtypes.TfStringValue `tfsdk:"group_id"`

	// Required — mutable
	Name     ovhtypes.TfStringValue `tfsdk:"name"`
	FlavorId ovhtypes.TfStringValue `tfsdk:"flavor_id"`

	// Optional — mutable
	ImageId          ovhtypes.TfStringValue                             `tfsdk:"image_id"`
	PowerState       ovhtypes.TfStringValue                             `tfsdk:"power_state"`
	Networks         types.List                                         `tfsdk:"networks"`
	VolumeIds        ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"volume_ids"`
	SecurityGroupIds ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"security_group_ids"`
	Shares           types.List                                         `tfsdk:"shares"`

	// Computed envelope
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// ---------- API DTOs (camelCase JSON tags, mirror internal/model/instance.go) ----------

type CloudInstanceRef struct {
	Id string `json:"id"`
}

type CloudInstanceAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudInstanceAPINetworkRef struct {
	Public       bool   `json:"public"`
	Id           string `json:"id,omitempty"`
	SubnetId     string `json:"subnetId,omitempty"`
	FloatingIpId string `json:"floatingIpId,omitempty"`
}

type CloudInstanceAPIShareRef struct {
	Id          string `json:"id"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

// Target spec sent on create (all settable fields incl. immutable location/ssh/group).
type CloudInstanceAPITargetSpec struct {
	Name           string                       `json:"name"`
	Flavor         *CloudInstanceRef            `json:"flavor,omitempty"`
	Image          *CloudInstanceRef            `json:"image,omitempty"`
	Location       *CloudInstanceAPILocation    `json:"location,omitempty"`
	Networks       []CloudInstanceAPINetworkRef `json:"networks,omitempty"`
	Volumes        []CloudInstanceRef           `json:"volumes,omitempty"`
	PowerState     string                       `json:"powerState,omitempty"`
	Group          *CloudInstanceRef            `json:"group,omitempty"`
	SSHKeyName     string                       `json:"sshKeyName,omitempty"`
	SecurityGroups []CloudInstanceRef           `json:"securityGroups,omitempty"`
	Shares         []CloudInstanceAPIShareRef   `json:"shares,omitempty"`
}

// Update target spec: mutable-only (no location / sshKeyName / group).
type CloudInstanceAPIUpdateTargetSpec struct {
	Name           string                       `json:"name"`
	Flavor         *CloudInstanceRef            `json:"flavor,omitempty"`
	Image          *CloudInstanceRef            `json:"image,omitempty"`
	Networks       []CloudInstanceAPINetworkRef `json:"networks,omitempty"`
	Volumes        []CloudInstanceRef           `json:"volumes,omitempty"`
	PowerState     string                       `json:"powerState,omitempty"`
	SecurityGroups []CloudInstanceRef           `json:"securityGroups,omitempty"`
	Shares         []CloudInstanceAPIShareRef   `json:"shares,omitempty"`
}

// Observed nested objects.
type CloudInstanceAPIFlavor struct {
	Id        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Vcpus     int64  `json:"vcpus,omitempty"`
	Ram       int64  `json:"ram,omitempty"`
	Disk      int64  `json:"disk,omitempty"`
	Swap      int64  `json:"swap,omitempty"`
	Ephemeral int64  `json:"ephemeral,omitempty"`
}

type CloudInstanceAPIImage struct {
	Id         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Size       int64  `json:"size,omitempty"`
	Status     string `json:"status,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
}

type CloudInstanceAPIAddress struct {
	Ip      string `json:"ip,omitempty"`
	Mac     string `json:"mac,omitempty"`
	Type    string `json:"type,omitempty"`
	Version int64  `json:"version,omitempty"`
}

type CloudInstanceAPINetworkState struct {
	Id           string                    `json:"id,omitempty"`
	Public       bool                      `json:"public"`
	SubnetId     string                    `json:"subnetId,omitempty"`
	GatewayId    string                    `json:"gatewayId,omitempty"`
	FloatingIpId string                    `json:"floatingIpId,omitempty"`
	Addresses    []CloudInstanceAPIAddress `json:"addresses,omitempty"`
}

type CloudInstanceAPIVolume struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type CloudInstanceAPIShareState struct {
	Id          string `json:"id"`
	AccessLevel string `json:"accessLevel,omitempty"`
	AccessTo    string `json:"accessTo,omitempty"`
	State       string `json:"state,omitempty"`
}

type CloudInstanceAPICurrentState struct {
	Name           string                         `json:"name,omitempty"`
	Flavor         *CloudInstanceAPIFlavor        `json:"flavor,omitempty"`
	Image          *CloudInstanceAPIImage         `json:"image,omitempty"`
	Location       *CloudInstanceAPILocation      `json:"location,omitempty"`
	PowerState     string                         `json:"powerState,omitempty"`
	Networks       []CloudInstanceAPINetworkState `json:"networks,omitempty"`
	Volumes        []CloudInstanceAPIVolume       `json:"volumes,omitempty"`
	Shares         []CloudInstanceAPIShareState   `json:"shares,omitempty"`
	SecurityGroups []CloudInstanceRef             `json:"securityGroups,omitempty"`
	Group          *CloudInstanceRef              `json:"group,omitempty"`
	Locked         bool                           `json:"locked,omitempty"`
	SSHKeyName     string                         `json:"sshKeyName,omitempty"`
	HostId         string                         `json:"hostId,omitempty"`
	ProjectId      string                         `json:"projectId,omitempty"`
	UserId         string                         `json:"userId,omitempty"`
}

type CloudInstanceAPIResponse struct {
	Id             string                        `json:"id"`
	Checksum       string                        `json:"checksum"`
	CreatedAt      string                        `json:"createdAt"`
	UpdatedAt      string                        `json:"updatedAt"`
	ResourceStatus string                        `json:"resourceStatus"`
	TargetSpec     *CloudInstanceAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudInstanceAPICurrentState `json:"currentState,omitempty"`
}

type CloudInstanceCreatePayload struct {
	TargetSpec *CloudInstanceAPITargetSpec `json:"targetSpec"`
}

type CloudInstanceUpdatePayload struct {
	Checksum   string                            `json:"checksum"`
	TargetSpec *CloudInstanceAPIUpdateTargetSpec `json:"targetSpec"`
}

// ---------- attr-type helpers ----------

func instanceNetworkRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"public":         types.BoolType,
		"network_id":     ovhtypes.TfStringType{},
		"subnet_id":      ovhtypes.TfStringType{},
		"floating_ip_id": ovhtypes.TfStringType{},
	}
}

func instanceShareRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
	}
}

func instanceFlavorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        ovhtypes.TfStringType{},
		"name":      ovhtypes.TfStringType{},
		"vcpus":     types.Int64Type,
		"ram":       types.Int64Type,
		"disk":      types.Int64Type,
		"swap":      types.Int64Type,
		"ephemeral": types.Int64Type,
	}
}

func instanceImageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         ovhtypes.TfStringType{},
		"name":       ovhtypes.TfStringType{},
		"size":       types.Int64Type,
		"status":     ovhtypes.TfStringType{},
		"deprecated": types.BoolType,
	}
}

func instanceLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

func instanceAddressAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":      ovhtypes.TfStringType{},
		"mac":     ovhtypes.TfStringType{},
		"type":    ovhtypes.TfStringType{},
		"version": types.Int64Type,
	}
}

func instanceNetworkStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             ovhtypes.TfStringType{},
		"public":         types.BoolType,
		"subnet_id":      ovhtypes.TfStringType{},
		"gateway_id":     ovhtypes.TfStringType{},
		"floating_ip_id": ovhtypes.TfStringType{},
		"addresses": types.ListType{
			ElemType: types.ObjectType{AttrTypes: instanceAddressAttrTypes()},
		},
	}
}

func instanceVolumeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"name": ovhtypes.TfStringType{},
		"size": types.Int64Type,
	}
}

func instanceShareStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
		"access_to":    ovhtypes.TfStringType{},
		"state":        ovhtypes.TfStringType{},
	}
}

func instanceSecurityGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

func instanceGroupRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

// InstanceCurrentStateAttrTypes returns the attribute types for current_state.
func InstanceCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":            ovhtypes.TfStringType{},
		"flavor":          types.ObjectType{AttrTypes: instanceFlavorAttrTypes()},
		"image":           types.ObjectType{AttrTypes: instanceImageAttrTypes()},
		"location":        types.ObjectType{AttrTypes: instanceLocationAttrTypes()},
		"power_state":     ovhtypes.TfStringType{},
		"networks":        types.ListType{ElemType: types.ObjectType{AttrTypes: instanceNetworkStateAttrTypes()}},
		"volumes":         types.ListType{ElemType: types.ObjectType{AttrTypes: instanceVolumeAttrTypes()}},
		"shares":          types.ListType{ElemType: types.ObjectType{AttrTypes: instanceShareStateAttrTypes()}},
		"security_groups": types.ListType{ElemType: types.ObjectType{AttrTypes: instanceSecurityGroupAttrTypes()}},
		"group":           types.ObjectType{AttrTypes: instanceGroupRefAttrTypes()},
		"locked":          types.BoolType,
		"ssh_key_name":    ovhtypes.TfStringType{},
		"host_id":         ovhtypes.TfStringType{},
		"project_id":      ovhtypes.TfStringType{},
		"user_id":         ovhtypes.TfStringType{},
	}
}

// ---------- ToCreate ----------

// networksToAPI converts the root networks list into API network refs.
func networksToAPI(list types.List) []CloudInstanceAPINetworkRef {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	out := make([]CloudInstanceAPINetworkRef, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		ref := CloudInstanceAPINetworkRef{}
		if v, ok := attrs["public"].(types.Bool); ok && !v.IsNull() {
			ref.Public = v.ValueBool()
		}
		if v, ok := attrs["network_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			ref.Id = v.ValueString()
		}
		if v, ok := attrs["subnet_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			ref.SubnetId = v.ValueString()
		}
		if v, ok := attrs["floating_ip_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			ref.FloatingIpId = v.ValueString()
		}
		out = append(out, ref)
	}
	return out
}

// sharesToAPI converts the root shares list into API share refs.
func sharesToAPI(list types.List) []CloudInstanceAPIShareRef {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	out := make([]CloudInstanceAPIShareRef, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		ref := CloudInstanceAPIShareRef{}
		if v, ok := attrs["id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			ref.Id = v.ValueString()
		}
		if v, ok := attrs["access_level"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			ref.AccessLevel = v.ValueString()
		}
		out = append(out, ref)
	}
	return out
}

// customStringListToRefs converts an ovhtypes string list into []CloudInstanceRef.
func customStringListToRefs(list ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]) []CloudInstanceRef {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	out := make([]CloudInstanceRef, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		if v, ok := elem.(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			out = append(out, CloudInstanceRef{Id: v.ValueString()})
		}
	}
	return out
}

// ToCreate builds the create payload including all fields (mutable + immutable).
func (m *CloudInstanceModel) ToCreate() *CloudInstanceCreatePayload {
	ts := &CloudInstanceAPITargetSpec{
		Name:   m.Name.ValueString(),
		Flavor: &CloudInstanceRef{Id: m.FlavorId.ValueString()},
		Location: &CloudInstanceAPILocation{
			Region: m.Region.ValueString(),
		},
	}

	if !m.ImageId.IsNull() && !m.ImageId.IsUnknown() && m.ImageId.ValueString() != "" {
		ts.Image = &CloudInstanceRef{Id: m.ImageId.ValueString()}
	}
	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		ts.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}
	if !m.PowerState.IsNull() && !m.PowerState.IsUnknown() {
		ts.PowerState = m.PowerState.ValueString()
	}
	if !m.SSHKeyName.IsNull() && !m.SSHKeyName.IsUnknown() {
		ts.SSHKeyName = m.SSHKeyName.ValueString()
	}
	if !m.GroupId.IsNull() && !m.GroupId.IsUnknown() && m.GroupId.ValueString() != "" {
		ts.Group = &CloudInstanceRef{Id: m.GroupId.ValueString()}
	}
	ts.Networks = networksToAPI(m.Networks)
	ts.Volumes = customStringListToRefs(m.VolumeIds)
	ts.SecurityGroups = customStringListToRefs(m.SecurityGroupIds)
	ts.Shares = sharesToAPI(m.Shares)

	return &CloudInstanceCreatePayload{TargetSpec: ts}
}

// ToUpdate builds the update payload with mutable fields only, plus checksum.
// Location, sshKeyName and group are immutable and intentionally excluded.
func (m *CloudInstanceModel) ToUpdate(checksum string) *CloudInstanceUpdatePayload {
	ts := &CloudInstanceAPIUpdateTargetSpec{
		Name:   m.Name.ValueString(),
		Flavor: &CloudInstanceRef{Id: m.FlavorId.ValueString()},
	}

	if !m.ImageId.IsNull() && !m.ImageId.IsUnknown() && m.ImageId.ValueString() != "" {
		ts.Image = &CloudInstanceRef{Id: m.ImageId.ValueString()}
	}
	if !m.PowerState.IsNull() && !m.PowerState.IsUnknown() {
		ts.PowerState = m.PowerState.ValueString()
	}
	ts.Networks = networksToAPI(m.Networks)
	ts.Volumes = customStringListToRefs(m.VolumeIds)
	ts.SecurityGroups = customStringListToRefs(m.SecurityGroupIds)
	ts.Shares = sharesToAPI(m.Shares)

	return &CloudInstanceUpdatePayload{Checksum: checksum, TargetSpec: ts}
}

func str(s string) ovhtypes.TfStringValue {
	return ovhtypes.TfStringValue{StringValue: types.StringValue(s)}
}

func buildInstanceFlavorObject(f *CloudInstanceAPIFlavor) basetypes.ObjectValue {
	if f == nil {
		return types.ObjectNull(instanceFlavorAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceFlavorAttrTypes(), map[string]attr.Value{
		"id":        str(f.Id),
		"name":      str(f.Name),
		"vcpus":     types.Int64Value(f.Vcpus),
		"ram":       types.Int64Value(f.Ram),
		"disk":      types.Int64Value(f.Disk),
		"swap":      types.Int64Value(f.Swap),
		"ephemeral": types.Int64Value(f.Ephemeral),
	})
	return obj
}

func buildInstanceImageObject(i *CloudInstanceAPIImage) basetypes.ObjectValue {
	if i == nil {
		return types.ObjectNull(instanceImageAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceImageAttrTypes(), map[string]attr.Value{
		"id":         str(i.Id),
		"name":       str(i.Name),
		"size":       types.Int64Value(i.Size),
		"status":     str(i.Status),
		"deprecated": types.BoolValue(i.Deprecated),
	})
	return obj
}

func buildInstanceLocationObject(l *CloudInstanceAPILocation) basetypes.ObjectValue {
	if l == nil {
		return types.ObjectNull(instanceLocationAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceLocationAttrTypes(), map[string]attr.Value{
		"region":            str(l.Region),
		"availability_zone": str(l.AvailabilityZone),
	})
	return obj
}

func buildInstanceNetworkStateList(networks []CloudInstanceAPINetworkState) basetypes.ListValue {
	netObjType := types.ObjectType{AttrTypes: instanceNetworkStateAttrTypes()}
	addrObjType := types.ObjectType{AttrTypes: instanceAddressAttrTypes()}
	if networks == nil {
		return types.ListNull(netObjType)
	}
	items := make([]attr.Value, 0, len(networks))
	for _, n := range networks {
		var addrs basetypes.ListValue
		if n.Addresses == nil {
			addrs = types.ListNull(addrObjType)
		} else {
			addrItems := make([]attr.Value, 0, len(n.Addresses))
			for _, a := range n.Addresses {
				addrObj, _ := types.ObjectValue(instanceAddressAttrTypes(), map[string]attr.Value{
					"ip":      str(a.Ip),
					"mac":     str(a.Mac),
					"type":    str(a.Type),
					"version": types.Int64Value(a.Version),
				})
				addrItems = append(addrItems, addrObj)
			}
			addrs = types.ListValueMust(addrObjType, addrItems)
		}
		netObj, _ := types.ObjectValue(instanceNetworkStateAttrTypes(), map[string]attr.Value{
			"id":             str(n.Id),
			"public":         types.BoolValue(n.Public),
			"subnet_id":      str(n.SubnetId),
			"gateway_id":     str(n.GatewayId),
			"floating_ip_id": str(n.FloatingIpId),
			"addresses":      addrs,
		})
		items = append(items, netObj)
	}
	return types.ListValueMust(netObjType, items)
}

func buildInstanceVolumeStateList(volumes []CloudInstanceAPIVolume) basetypes.ListValue {
	objType := types.ObjectType{AttrTypes: instanceVolumeAttrTypes()}
	if volumes == nil {
		return types.ListNull(objType)
	}
	items := make([]attr.Value, 0, len(volumes))
	for _, v := range volumes {
		obj, _ := types.ObjectValue(instanceVolumeAttrTypes(), map[string]attr.Value{
			"id":   str(v.Id),
			"name": str(v.Name),
			"size": types.Int64Value(v.Size),
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

func buildInstanceShareStateList(shares []CloudInstanceAPIShareState) basetypes.ListValue {
	objType := types.ObjectType{AttrTypes: instanceShareStateAttrTypes()}
	if shares == nil {
		return types.ListNull(objType)
	}
	items := make([]attr.Value, 0, len(shares))
	for _, s := range shares {
		obj, _ := types.ObjectValue(instanceShareStateAttrTypes(), map[string]attr.Value{
			"id":           str(s.Id),
			"access_level": str(s.AccessLevel),
			"access_to":    str(s.AccessTo),
			"state":        str(s.State),
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

func buildInstanceSecurityGroupStateList(sgs []CloudInstanceRef) basetypes.ListValue {
	objType := types.ObjectType{AttrTypes: instanceSecurityGroupAttrTypes()}
	if sgs == nil {
		return types.ListNull(objType)
	}
	items := make([]attr.Value, 0, len(sgs))
	for _, sg := range sgs {
		obj, _ := types.ObjectValue(instanceSecurityGroupAttrTypes(), map[string]attr.Value{
			"id": str(sg.Id),
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

func buildInstanceGroupObject(g *CloudInstanceRef) basetypes.ObjectValue {
	if g == nil {
		return types.ObjectNull(instanceGroupRefAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceGroupRefAttrTypes(), map[string]attr.Value{
		"id": str(g.Id),
	})
	return obj
}

// buildInstanceCurrentStateObject assembles the current_state object from the API currentState.
func buildInstanceCurrentStateObject(ctx context.Context, state *CloudInstanceAPICurrentState) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(InstanceCurrentStateAttrTypes(), map[string]attr.Value{
		"name":            str(state.Name),
		"flavor":          buildInstanceFlavorObject(state.Flavor),
		"image":           buildInstanceImageObject(state.Image),
		"location":        buildInstanceLocationObject(state.Location),
		"power_state":     str(state.PowerState),
		"networks":        buildInstanceNetworkStateList(state.Networks),
		"volumes":         buildInstanceVolumeStateList(state.Volumes),
		"shares":          buildInstanceShareStateList(state.Shares),
		"security_groups": buildInstanceSecurityGroupStateList(state.SecurityGroups),
		"group":           buildInstanceGroupObject(state.Group),
		"locked":          types.BoolValue(state.Locked),
		"ssh_key_name":    str(state.SSHKeyName),
		"host_id":         str(state.HostId),
		"project_id":      str(state.ProjectId),
		"user_id":         str(state.UserId),
	})
	return obj
}

// buildInstanceNetworksRootList rebuilds the mutable root `networks` list from targetSpec.
func buildInstanceNetworksRootList(ts *CloudInstanceAPITargetSpec) types.List {
	objType := types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}
	if ts == nil || ts.Networks == nil {
		return types.ListNull(objType)
	}
	items := make([]attr.Value, 0, len(ts.Networks))
	for _, n := range ts.Networks {
		networkId := types.StringNull()
		if n.Id != "" {
			networkId = types.StringValue(n.Id)
		}
		subnetId := types.StringNull()
		if n.SubnetId != "" {
			subnetId = types.StringValue(n.SubnetId)
		}
		floatingIpId := types.StringNull()
		if n.FloatingIpId != "" {
			floatingIpId = types.StringValue(n.FloatingIpId)
		}
		obj, _ := types.ObjectValue(instanceNetworkRefAttrTypes(), map[string]attr.Value{
			"public":         types.BoolValue(n.Public),
			"network_id":     ovhtypes.TfStringValue{StringValue: networkId},
			"subnet_id":      ovhtypes.TfStringValue{StringValue: subnetId},
			"floating_ip_id": ovhtypes.TfStringValue{StringValue: floatingIpId},
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

// buildInstanceSharesRootList rebuilds the mutable root `shares` list from targetSpec.
func buildInstanceSharesRootList(ts *CloudInstanceAPITargetSpec) types.List {
	objType := types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}
	if ts == nil || ts.Shares == nil {
		return types.ListNull(objType)
	}
	items := make([]attr.Value, 0, len(ts.Shares))
	for _, s := range ts.Shares {
		accessLevel := types.StringNull()
		if s.AccessLevel != "" {
			accessLevel = types.StringValue(s.AccessLevel)
		}
		obj, _ := types.ObjectValue(instanceShareRefAttrTypes(), map[string]attr.Value{
			"id":           str(s.Id),
			"access_level": ovhtypes.TfStringValue{StringValue: accessLevel},
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

// MergeWith copies the API response into the Terraform model.
func (m *CloudInstanceModel) MergeWith(ctx context.Context, response *CloudInstanceAPIResponse) {
	m.Id = str(response.Id)
	m.Checksum = str(response.Checksum)
	m.CreatedAt = str(response.CreatedAt)
	m.UpdatedAt = str(response.UpdatedAt)
	m.ResourceStatus = str(response.ResourceStatus)

	if response.CurrentState != nil {
		m.CurrentState = buildInstanceCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		ts := response.TargetSpec
		m.Name = str(ts.Name)
		if ts.Flavor != nil {
			m.FlavorId = str(ts.Flavor.Id)
		}
		if ts.Image != nil && ts.Image.Id != "" {
			m.ImageId = str(ts.Image.Id)
		} else {
			m.ImageId = ovhtypes.TfStringValue{StringValue: types.StringNull()}
		}
		if ts.Location != nil {
			m.Region = str(ts.Location.Region)
			if ts.Location.AvailabilityZone != "" {
				m.AvailabilityZone = str(ts.Location.AvailabilityZone)
			}
		}
		if ts.PowerState != "" {
			m.PowerState = str(ts.PowerState)
		}
		if ts.SSHKeyName != "" {
			m.SSHKeyName = str(ts.SSHKeyName)
		}
		if ts.Group != nil && ts.Group.Id != "" {
			m.GroupId = str(ts.Group.Id)
		}
		m.Networks = buildInstanceNetworksRootList(ts)
		m.Shares = buildInstanceSharesRootList(ts)

		// volume_ids / security_group_ids from targetSpec
		if ts.Volumes != nil {
			vals := make([]attr.Value, len(ts.Volumes))
			for i, v := range ts.Volumes {
				vals[i] = str(v.Id)
			}
			m.VolumeIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals)}
		} else {
			m.VolumeIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})}
		}
		if ts.SecurityGroups != nil {
			vals := make([]attr.Value, len(ts.SecurityGroups))
			for i, sg := range ts.SecurityGroups {
				vals[i] = str(sg.Id)
			}
			m.SecurityGroupIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals)}
		} else {
			m.SecurityGroupIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})}
		}
	}
}
