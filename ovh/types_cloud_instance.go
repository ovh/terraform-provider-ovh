package ovh

import (
	"context"
	"strings"

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
	Id                 string `json:"id,omitempty"`
	SubnetId           string `json:"subnetId,omitempty"`
	IP                 string `json:"ip,omitempty"`
	AutoAssignPublicIp bool   `json:"autoAssignPublicIp,omitempty"`
}

type CloudInstanceAPIShareRef struct {
	Id          string `json:"id"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

// Target spec sent on create (all settable fields incl. immutable location/ssh/group).
type CloudInstanceAPITargetSpec struct {
	Name       string                       `json:"name"`
	Flavor     *CloudInstanceRef            `json:"flavor,omitempty"`
	Image      *CloudInstanceRef            `json:"image,omitempty"`
	Location   *CloudInstanceAPILocation    `json:"location,omitempty"`
	Networks   []CloudInstanceAPINetworkRef `json:"networks,omitempty"`
	Volumes    []CloudInstanceRef           `json:"volumes,omitempty"`
	PowerState string                       `json:"powerState,omitempty"`
	Group      *CloudInstanceRef            `json:"group,omitempty"`
	SSHKeyName string                       `json:"sshKeyName,omitempty"`
	// No omitempty: the API distinguishes null (apply the project's default
	// security group on create, leave the groups unchanged on update) from an
	// explicit empty array (apply no security group at all).
	SecurityGroups []CloudInstanceRef         `json:"securityGroups"`
	Shares         []CloudInstanceAPIShareRef `json:"shares,omitempty"`
}

// Update target spec: mutable-only (no location / sshKeyName / group).
type CloudInstanceAPIUpdateTargetSpec struct {
	Name       string                       `json:"name"`
	Flavor     *CloudInstanceRef            `json:"flavor,omitempty"`
	Image      *CloudInstanceRef            `json:"image,omitempty"`
	Networks   []CloudInstanceAPINetworkRef `json:"networks,omitempty"`
	Volumes    []CloudInstanceRef           `json:"volumes,omitempty"`
	PowerState string                       `json:"powerState,omitempty"`
	// See CloudInstanceAPITargetSpec.SecurityGroups: null and [] differ.
	SecurityGroups []CloudInstanceRef         `json:"securityGroups"`
	Shares         []CloudInstanceAPIShareRef `json:"shares,omitempty"`
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
	Id        string                    `json:"id,omitempty"`
	SubnetId  string                    `json:"subnetId,omitempty"`
	GatewayId string                    `json:"gatewayId,omitempty"`
	Addresses []CloudInstanceAPIAddress `json:"addresses,omitempty"`
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
	CurrentTasks   []CloudResourceTask           `json:"currentTasks,omitempty"`
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

// ---------- shared attribute descriptions ----------

// Shared by the ovh_cloud_instance resource and the ovh_cloud_instance /
// ovh_cloud_instances data sources so a given field never carries two different
// wordings. Only the *Md variants exist where the markdown rendering differs.
const (
	instanceDescServiceName = "Service name of the resource representing the id of the cloud project"
	instanceDescId          = "Unique identifier of the instance"
	instanceDescChecksum    = "Computed hash representing the current target specification value. It implements optimistic concurrency control: the value is echoed back on update and the request is rejected when it no longer matches server-side"
	instanceDescCreatedAt   = "Creation date of the instance, as an RFC 3339 timestamp"
	instanceDescUpdatedAt   = "Last modification date of the instance, as an RFC 3339 timestamp"

	instanceDescResourceStatus   = "Instance readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING). Distinct from current_state.power_state, which carries the lower-level OpenStack administrative power state"
	instanceDescResourceStatusMd = "Instance readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`). Distinct from `current_state.power_state`, which carries the lower-level OpenStack administrative power state"

	instanceDescCurrentState = "Observed state of the instance as reported by the compute backend, as opposed to the requested specification exposed at root level. Null while the instance is still being created and no backend state is available yet"

	instanceDescCsName       = "Observed display name of the instance"
	instanceDescCsPowerState = "Observed administrative power state of the instance as reported by OpenStack. It may transiently differ from the requested power_state while a power transition is in progress"
	instanceDescCsLocked     = "Whether the instance is locked against modifications. While locked, mutating actions are refused until it is unlocked"
	instanceDescCsSSHKeyName = "Name of the SSH key pair injected into the instance at boot, null when none was provided"
	instanceDescCsHostId     = "Opaque identifier of the physical host the instance is running on, as exposed by OpenStack. Null when not available"
	instanceDescCsProjectId  = "Identifier of the Public Cloud project the instance belongs to"
	instanceDescCsUserId     = "Identifier of the OpenStack user that owns the instance"

	instanceDescCsFlavor        = "Observed flavor of the instance, with its full sizing details"
	instanceDescFlavorId        = "Unique identifier of the flavor"
	instanceDescFlavorName      = "Human-readable flavor name (the commercial flavor label)"
	instanceDescFlavorVcpus     = "Number of virtual CPUs provided by the flavor"
	instanceDescFlavorRam       = "Amount of RAM provided by the flavor, in MB"
	instanceDescFlavorDisk      = "Size of the flavor's local root disk, in GB"
	instanceDescFlavorSwap      = "Size of the flavor's swap space, in MB"
	instanceDescFlavorEphemeral = "Size of the flavor's ephemeral disk, in GB"

	instanceDescCsImage         = "Observed image the instance was booted from, null for a boot-from-volume instance which has no image"
	instanceDescImageId         = "Unique identifier of the image"
	instanceDescImageName       = "Human-readable image name, null when the backend does not report it"
	instanceDescImageSize       = "Size of the image, in bytes. Null when the backend does not report it"
	instanceDescImageStatus     = "Lifecycle status of the image as reported by Glance"
	instanceDescImageDeprecated = "Whether the image is flagged as deprecated. A deprecated image still boots existing instances but is no longer recommended for new ones"

	instanceDescCsLocation     = "Observed region and availability zone where the instance is provisioned"
	instanceDescLocationRegion = "Code of the region where the instance is provisioned (for example GRA11, BHS5)"
	instanceDescLocationAZ     = "Availability zone within the region where the instance is placed, null in regions that have none"

	instanceDescCsNetworks         = "Observed network interfaces of the instance: one entry per private network plus at most one entry without a network id for the public (Ext-Net) interface. Entries are ordered by network id, so they do not follow the order of the requested networks"
	instanceDescCsNetworkId        = "Identifier of the network this interface is attached to, null for the public (Ext-Net) interface"
	instanceDescCsNetworkSubnetId  = "Identifier of the subnet this interface draws its fixed address from, null for an entry without a network id"
	instanceDescCsNetworkGatewayId = "Identifier of the gateway providing egress for this interface, null when none applies"
	instanceDescCsAddresses        = "Addresses observed on this interface: its fixed addresses plus, where applicable, its floating IP and any additional IPs routed to it. Each address carries a type of FIXED, FLOATING or ADDITIONAL"
	instanceDescCsAddressesMd      = "Addresses observed on this interface: its fixed addresses plus, where applicable, its floating IP and any additional IPs routed to it. Each address carries a type of `FIXED`, `FLOATING` or `ADDITIONAL`"
	instanceDescAddressIp          = "IP address assigned to the interface (IPv4 or IPv6)"
	instanceDescAddressMac         = "MAC address of the interface this IP is bound to. Null when the backend reports no interface for the address, which happens for an additional IP routed to an instance whose public interface has no visible Ext-Net address yet"
	instanceDescAddressType        = "How this address reaches the instance: FIXED for an address assigned to the interface itself, FLOATING for a floating IP NAT'd onto it, ADDITIONAL for an additional IP routed to the public interface"
	instanceDescAddressTypeMd      = "How this address reaches the instance: `FIXED` for an address assigned to the interface itself, `FLOATING` for a floating IP NAT'd onto it, `ADDITIONAL` for an additional IP routed to the public interface"
	instanceDescAddressVersion     = "IP version of the address (4 for IPv4, 6 for IPv6)"

	instanceDescCsVolumes  = "Observed block volumes attached to the instance"
	instanceDescVolumeId   = "Unique identifier of the attached volume"
	instanceDescVolumeName = "Display name of the attached volume"
	instanceDescVolumeSize = "Size of the attached volume, in GB"

	instanceDescCsShares             = "Observed instance-side share attachments, derived from the Manila access rules that target one of the instance's IPs. Only populated on a single-instance read, never in the list data source"
	instanceDescCsShareId            = "Identifier of the attached file storage share"
	instanceDescCsShareAccessLevel   = "Observed access level of the access rule for this instance (READ_ONLY or READ_WRITE)"
	instanceDescCsShareAccessLevelMd = "Observed access level of the access rule for this instance (`READ_ONLY` or `READ_WRITE`)"
	instanceDescCsShareAccessTo      = "The instance IP address the Manila access rule targets: its fixed IPv4 address on the share's network"
	instanceDescCsShareState         = "Observed state of the underlying access rule (ACTIVE, APPLYING, DENYING, ERROR). Null while no state has been reported yet"
	instanceDescCsShareStateMd       = "Observed state of the underlying access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`). Null while no state has been reported yet"

	instanceDescCsSecurityGroups = "Security groups currently attached to the instance's ports"
	instanceDescSecurityGroupId  = "Security group identifier"

	instanceDescCsGroup = "Instance (placement) group the instance belongs to, null when it is not part of any group"
	instanceDescGroupId = "Identifier of the instance (placement) group"

	instanceDescNetworkRefNetworkId    = "Private network ID. Omit for a public interface"
	instanceDescNetworkRefSubnetId     = "Subnet ID within the private network. Required with network_id"
	instanceDescNetworkRefSubnetIdMd   = "Subnet ID within the private network. Required with `network_id`"
	instanceDescNetworkRefIp           = "IP address of this interface. Without network_id: a public IP the project already owns (additional IP, or an Ext-Net IP of the project in the instance's region). With network_id + subnet_id: pins the port's fixed address when inside the subnet CIDR, otherwise associates the existing floating IP with that address"
	instanceDescNetworkRefIpMd         = "IP address of this interface. Without `network_id`: a public IP the project already owns (additional IP, or an Ext-Net IP of the project in the instance's region). With `network_id` + `subnet_id`: pins the port's fixed address when inside the subnet CIDR, otherwise associates the existing floating IP with that address"
	instanceDescNetworkRefAutoAssign   = "Attach a public interface with a public IP assigned by the platform. Only valid on an entry with no network_id and no ip, and on at most one entry"
	instanceDescNetworkRefAutoAssignMd = "Attach a public interface with a public IP assigned by the platform. Only valid on an entry with no `network_id` and no `ip`, and on at most one entry"

	instanceDescVolumeIds = "IDs of block-storage volumes attached to the instance"

	// Read-only wordings for the data sources, where the root attributes expose
	// the requested target spec rather than a writable input.
	instanceDescDsName             = "Instance name"
	instanceDescDsImageId          = "Identifier of the image the instance boots from, null for a boot-from-volume instance"
	instanceDescDsPowerState       = "Requested power state of the instance (ACTIVE, SHUTOFF or SHELVED)"
	instanceDescDsPowerStateMd     = "Requested power state of the instance (`ACTIVE`, `SHUTOFF` or `SHELVED`)"
	instanceDescDsNetworks         = "Requested network interfaces of the instance, as recorded in its target specification"
	instanceDescDsSecurityGroupIds = "IDs of the security groups applied to the instance's interfaces"
	instanceDescDsShares           = "File storage shares requested to be attached to the instance"

	instanceDescShareRefId            = "Identifier of the file storage share to attach. Adding the reference attaches the share to the instance, removing it detaches it"
	instanceDescShareRefAccessLevel   = "Access level granted to the instance for this share: READ_ONLY or READ_WRITE. Omit it to let the API apply its default (READ_WRITE)"
	instanceDescShareRefAccessLevelMd = "Access level granted to the instance for this share: `READ_ONLY` or `READ_WRITE`. Omit it to let the API apply its default (`READ_WRITE`)"
)

// ---------- attr-type helpers ----------

func instanceNetworkRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"network_id":            ovhtypes.TfStringType{},
		"subnet_id":             ovhtypes.TfStringType{},
		"ip":                    ovhtypes.TfStringType{},
		"auto_assign_public_ip": types.BoolType,
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
		"id":         ovhtypes.TfStringType{},
		"subnet_id":  ovhtypes.TfStringType{},
		"gateway_id": ovhtypes.TfStringType{},
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

// networkRefFromObject reads one `networks` element into an API network ref.
// Unknown and null attributes both collapse to the zero value.
func networkRefFromObject(obj types.Object) CloudInstanceAPINetworkRef {
	attrs := obj.Attributes()
	ref := CloudInstanceAPINetworkRef{}
	if v, ok := attrs["network_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		ref.Id = v.ValueString()
	}
	if v, ok := attrs["subnet_id"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		ref.SubnetId = v.ValueString()
	}
	if v, ok := attrs["ip"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
		ref.IP = v.ValueString()
	}
	if v, ok := attrs["auto_assign_public_ip"].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		ref.AutoAssignPublicIp = v.ValueBool()
	}
	return ref
}

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
		out = append(out, networkRefFromObject(obj))
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

func tfStr(s string) ovhtypes.TfStringValue {
	return ovhtypes.TfStringValue{StringValue: types.StringValue(s)}
}

func buildInstanceFlavorObject(f *CloudInstanceAPIFlavor) basetypes.ObjectValue {
	if f == nil {
		return types.ObjectNull(instanceFlavorAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceFlavorAttrTypes(), map[string]attr.Value{
		"id":        tfStr(f.Id),
		"name":      tfStr(f.Name),
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
		"id":         tfStr(i.Id),
		"name":       tfStr(i.Name),
		"size":       types.Int64Value(i.Size),
		"status":     tfStr(i.Status),
		"deprecated": types.BoolValue(i.Deprecated),
	})
	return obj
}

func buildInstanceLocationObject(l *CloudInstanceAPILocation) basetypes.ObjectValue {
	if l == nil {
		return types.ObjectNull(instanceLocationAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceLocationAttrTypes(), map[string]attr.Value{
		"region":            tfStr(l.Region),
		"availability_zone": tfStr(l.AvailabilityZone),
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
					"ip":      tfStr(a.Ip),
					"mac":     tfStr(a.Mac),
					"type":    tfStr(a.Type),
					"version": types.Int64Value(a.Version),
				})
				addrItems = append(addrItems, addrObj)
			}
			addrs = types.ListValueMust(addrObjType, addrItems)
		}
		netObj, _ := types.ObjectValue(instanceNetworkStateAttrTypes(), map[string]attr.Value{
			"id":         tfStr(n.Id),
			"subnet_id":  tfStr(n.SubnetId),
			"gateway_id": tfStr(n.GatewayId),
			"addresses":  addrs,
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
			"id":   tfStr(v.Id),
			"name": tfStr(v.Name),
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
			"id":           tfStr(s.Id),
			"access_level": tfStr(s.AccessLevel),
			"access_to":    tfStr(s.AccessTo),
			"state":        tfStr(s.State),
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
			"id": tfStr(sg.Id),
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
		"id": tfStr(g.Id),
	})
	return obj
}

// buildInstanceCurrentStateObject assembles the current_state object from the API currentState.
func buildInstanceCurrentStateObject(ctx context.Context, state *CloudInstanceAPICurrentState) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(InstanceCurrentStateAttrTypes(), map[string]attr.Value{
		"name":            tfStr(state.Name),
		"flavor":          buildInstanceFlavorObject(state.Flavor),
		"image":           buildInstanceImageObject(state.Image),
		"location":        buildInstanceLocationObject(state.Location),
		"power_state":     tfStr(state.PowerState),
		"networks":        buildInstanceNetworkStateList(state.Networks),
		"volumes":         buildInstanceVolumeStateList(state.Volumes),
		"shares":          buildInstanceShareStateList(state.Shares),
		"security_groups": buildInstanceSecurityGroupStateList(state.SecurityGroups),
		"group":           buildInstanceGroupObject(state.Group),
		"locked":          types.BoolValue(state.Locked),
		"ssh_key_name":    tfStr(state.SSHKeyName),
		"host_id":         tfStr(state.HostId),
		"project_id":      tfStr(state.ProjectId),
		"user_id":         tfStr(state.UserId),
	})
	return obj
}

func instanceNetworkRefKey(n CloudInstanceAPINetworkRef) string {
	autoAssign := "0"
	if n.AutoAssignPublicIp {
		autoAssign = "1"
	}
	return strings.Join([]string{n.Id, n.SubnetId, n.IP, autoAssign}, "\x00")
}

// orderInstanceNetworksLikePrior re-emits the API networks in the order of the
// prior plan/state list. apiv2 sorts the derived target spec by network UUID
// (id-less public entry first), but `networks` is Optional-only, so Terraform
// requires the applied value to match the config index-for-index — a
// private-first config would otherwise fail with an inconsistent-result error.
// Entries that correlate to nothing in the prior list are genuine out-of-band
// drift and are appended, in API order, at the end.
func orderInstanceNetworksLikePrior(prior types.List, apiNets []CloudInstanceAPINetworkRef) []CloudInstanceAPINetworkRef {
	// Null/unknown prior is the `terraform import` path: keep the API order.
	if prior.IsNull() || prior.IsUnknown() {
		return apiNets
	}
	free := make(map[string][]int, len(apiNets))
	for i, n := range apiNets {
		key := instanceNetworkRefKey(n)
		free[key] = append(free[key], i)
	}
	consumed := make([]bool, len(apiNets))
	out := make([]CloudInstanceAPINetworkRef, 0, len(apiNets))
	for _, elem := range prior.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		key := instanceNetworkRefKey(networkRefFromObject(obj))
		idxs := free[key]
		if len(idxs) == 0 {
			continue
		}
		free[key] = idxs[1:]
		consumed[idxs[0]] = true
		out = append(out, apiNets[idxs[0]])
	}
	for i, n := range apiNets {
		if !consumed[i] {
			out = append(out, n)
		}
	}
	return out
}

// buildInstanceNetworksRootList rebuilds the mutable root `networks` list from
// targetSpec, ordered like the prior plan/state list.
func buildInstanceNetworksRootList(prior types.List, ts *CloudInstanceAPITargetSpec) types.List {
	objType := types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()}
	if ts == nil || len(ts.Networks) == 0 {
		return types.ListNull(objType)
	}
	ordered := orderInstanceNetworksLikePrior(prior, ts.Networks)
	items := make([]attr.Value, 0, len(ordered))
	for _, n := range ordered {
		networkId := types.StringNull()
		if n.Id != "" {
			networkId = types.StringValue(n.Id)
		}
		subnetId := types.StringNull()
		if n.SubnetId != "" {
			subnetId = types.StringValue(n.SubnetId)
		}
		ip := types.StringNull()
		if n.IP != "" {
			ip = types.StringValue(n.IP)
		}
		// Both attributes are Optional-only: a false/empty value must round-trip as
		// null, otherwise a config that omitted it gets an inconsistent-result error.
		autoAssign := types.BoolNull()
		if n.AutoAssignPublicIp {
			autoAssign = types.BoolValue(true)
		}
		obj, _ := types.ObjectValue(instanceNetworkRefAttrTypes(), map[string]attr.Value{
			"network_id":            ovhtypes.TfStringValue{StringValue: networkId},
			"subnet_id":             ovhtypes.TfStringValue{StringValue: subnetId},
			"ip":                    ovhtypes.TfStringValue{StringValue: ip},
			"auto_assign_public_ip": autoAssign,
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

// instanceShareRefDefaultAccessLevel is the value apiv2 normalises an omitted
// accessLevel to before persisting the target spec, and therefore echoes back.
const instanceShareRefDefaultAccessLevel = "READ_WRITE"

// priorShareAccessLevelOmitted reports the share ids whose access_level the
// prior plan/state left unset. Share ids are unique in a target spec (apiv2
// rejects duplicates), so a set keyed by id is enough.
func priorShareAccessLevelOmitted(prior types.List) map[string]bool {
	if prior.IsNull() || prior.IsUnknown() {
		return nil
	}
	out := make(map[string]bool, len(prior.Elements()))
	for _, elem := range prior.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		id, ok := attrs["id"].(ovhtypes.TfStringValue)
		if !ok || id.IsNull() || id.IsUnknown() {
			continue
		}
		level, ok := attrs["access_level"].(ovhtypes.TfStringValue)
		if ok && (level.IsNull() || level.IsUnknown()) {
			out[id.ValueString()] = true
		}
	}
	return out
}

// buildInstanceSharesRootList rebuilds the mutable root `shares` list from targetSpec.
func buildInstanceSharesRootList(prior types.List, ts *CloudInstanceAPITargetSpec) types.List {
	objType := types.ObjectType{AttrTypes: instanceShareRefAttrTypes()}
	if ts == nil || len(ts.Shares) == 0 {
		return types.ListNull(objType)
	}
	omitted := priorShareAccessLevelOmitted(prior)
	items := make([]attr.Value, 0, len(ts.Shares))
	for _, s := range ts.Shares {
		// access_level is Optional-only: apiv2 normalises an omitted level to
		// READ_WRITE and echoes it back, so keep it null when the config left it
		// unset, otherwise the apply result would not match the configuration.
		accessLevel := types.StringNull()
		if s.AccessLevel != "" && !(omitted[s.Id] && s.AccessLevel == instanceShareRefDefaultAccessLevel) {
			accessLevel = types.StringValue(s.AccessLevel)
		}
		obj, _ := types.ObjectValue(instanceShareRefAttrTypes(), map[string]attr.Value{
			"id":           tfStr(s.Id),
			"access_level": ovhtypes.TfStringValue{StringValue: accessLevel},
		})
		items = append(items, obj)
	}
	return types.ListValueMust(objType, items)
}

// instancePriorSpec carries the plan/state values MergeWith needs to keep the
// merged model shaped like the configuration: list order and optionals the user
// omitted but the API resolves and echoes back.
type instancePriorSpec struct {
	Networks types.List
	Shares   types.List
}

// priorSpec snapshots the config-shaped values before a merge overwrites them.
func (m *CloudInstanceModel) priorSpec() instancePriorSpec {
	return instancePriorSpec{Networks: m.Networks, Shares: m.Shares}
}

// MergeWith copies the API response into the Terraform model.
func (m *CloudInstanceModel) MergeWith(ctx context.Context, response *CloudInstanceAPIResponse, prior instancePriorSpec) {
	m.Id = tfStr(response.Id)
	m.Checksum = tfStr(response.Checksum)
	m.CreatedAt = tfStr(response.CreatedAt)
	m.UpdatedAt = tfStr(response.UpdatedAt)
	m.ResourceStatus = tfStr(response.ResourceStatus)

	if response.CurrentState != nil {
		m.CurrentState = buildInstanceCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(InstanceCurrentStateAttrTypes())
	}

	// Typed nulls first: a zero-value list has a nil element type and panics at
	// resp.State.Set when the response carries no target spec.
	m.Networks = types.ListNull(types.ObjectType{AttrTypes: instanceNetworkRefAttrTypes()})
	m.Shares = types.ListNull(types.ObjectType{AttrTypes: instanceShareRefAttrTypes()})
	m.VolumeIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})}
	m.SecurityGroupIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})}

	if response.TargetSpec != nil {
		ts := response.TargetSpec
		m.Name = tfStr(ts.Name)
		if ts.Flavor != nil {
			m.FlavorId = tfStr(ts.Flavor.Id)
		}
		if ts.Image != nil && ts.Image.Id != "" {
			m.ImageId = tfStr(ts.Image.Id)
		} else {
			m.ImageId = ovhtypes.TfStringValue{StringValue: types.StringNull()}
		}
		if ts.Location != nil {
			m.Region = tfStr(ts.Location.Region)
			if ts.Location.AvailabilityZone != "" {
				m.AvailabilityZone = tfStr(ts.Location.AvailabilityZone)
			}
		}
		if ts.PowerState != "" {
			m.PowerState = tfStr(ts.PowerState)
		}
		if ts.SSHKeyName != "" {
			m.SSHKeyName = tfStr(ts.SSHKeyName)
		}
		if ts.Group != nil && ts.Group.Id != "" {
			m.GroupId = tfStr(ts.Group.Id)
		}
		m.Networks = buildInstanceNetworksRootList(prior.Networks, ts)
		m.Shares = buildInstanceSharesRootList(prior.Shares, ts)

		// volume_ids from targetSpec. It is Optional (non-Computed): an omitted
		// attribute is null in config, so the API echoing back an empty array must
		// map to null, not an empty list.
		if len(ts.Volumes) > 0 {
			vals := make([]attr.Value, len(ts.Volumes))
			for i, v := range ts.Volumes {
				vals[i] = tfStr(v.Id)
			}
			m.VolumeIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals)}
		} else {
			m.VolumeIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListNull(ovhtypes.TfStringType{})}
		}

		// security_group_ids is Optional+Computed: an omitted attribute is resolved
		// server-side to the project's default security group, and an explicit empty
		// list means no security group. "No group" therefore round-trips as an empty
		// list and never as null, so refresh and import stay stable whether the API
		// echoes securityGroups as [] or omits it.
		vals := make([]attr.Value, len(ts.SecurityGroups))
		for i, sg := range ts.SecurityGroups {
			vals[i] = tfStr(sg.Id)
		}
		m.SecurityGroupIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, vals)}
	}

	// The flat region/availability_zone come from targetSpec (the requested
	// spec), but the availability zone is often platform-assigned and only
	// surfaces in currentState. Fall back to the observed location so the flat
	// fields don't disagree with current_state.location.
	if response.CurrentState != nil && response.CurrentState.Location != nil {
		loc := response.CurrentState.Location
		if m.Region.IsNull() || m.Region.ValueString() == "" {
			m.Region = tfStr(loc.Region)
		}
		if (m.AvailabilityZone.IsNull() || m.AvailabilityZone.ValueString() == "") && loc.AvailabilityZone != "" {
			m.AvailabilityZone = tfStr(loc.AvailabilityZone)
		}
	}

	// availability_zone is Optional+Computed: when neither the requested spec nor
	// the observed state carries one (e.g. a non-3AZ region), the planned value
	// stays unknown. It must be known after apply, so resolve it to null.
	if m.AvailabilityZone.IsUnknown() {
		m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}

	// power_state is Optional+Computed and defaults server-side. If the target
	// spec didn't echo it back, fall back to the observed state, then null, so
	// the value is always known after apply.
	if m.PowerState.IsUnknown() {
		if response.CurrentState != nil && response.CurrentState.PowerState != "" {
			m.PowerState = tfStr(response.CurrentState.PowerState)
		} else {
			m.PowerState = ovhtypes.TfStringValue{StringValue: types.StringNull()}
		}
	}

	// image_id is Optional+Computed: guard against an unknown leaking through if
	// the target spec omitted it and the else-branch above wasn't reached.
	if m.ImageId.IsUnknown() {
		m.ImageId = ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}

	// security_group_ids is Optional+Computed and must be known after apply, even
	// when the response carried no target spec at all.
	if m.SecurityGroupIds.IsUnknown() {
		m.SecurityGroupIds = ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]{ListValue: basetypes.NewListValueMust(ovhtypes.TfStringType{}, []attr.Value{})}
	}
}
