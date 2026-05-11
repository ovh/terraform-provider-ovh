package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudProjectQuotaModel is the Terraform model for the quota data source.
type CloudProjectQuotaModel struct {
	// Input
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Envelope (computed)
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	TargetSpec     types.Object           `tfsdk:"target_spec"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// ------------------------------------------------------------------
// API response structs — mirror the public-cloud-apiv2 model.Quota
// ------------------------------------------------------------------

type CloudProjectQuotaUsageAPI struct {
	Limit int    `json:"limit"`
	Used  *int   `json:"used"`
	Unit  string `json:"unit"`
}

type CloudProjectQuotaLimitAPI struct {
	Limit int    `json:"limit"`
	Unit  string `json:"unit"`
}

type CloudProjectQuotaTargetSpecAPI struct {
	Profile string `json:"profile"`
}

type CloudProjectQuotaProfileComputeAPI struct {
	Instances          int `json:"instances"`
	Cores              int `json:"cores"`
	Ram                int `json:"ram"`
	SecurityGroups     int `json:"securityGroups"`
	SecurityGroupRules int `json:"securityGroupRules"`
	ServerGroups       int `json:"serverGroups"`
	ServerGroupMembers int `json:"serverGroupMembers"`
}

type CloudProjectQuotaProfileVolumeAPI struct {
	Volumes         int `json:"volumes"`
	Gigabytes       int `json:"gigabytes"`
	Snapshots       int `json:"snapshots"`
	Backups         int `json:"backups"`
	BackupGigabytes int `json:"backupGigabytes"`
}

type CloudProjectQuotaProfileNetworkAPI struct {
	Networks           int `json:"networks"`
	Subnets            int `json:"subnets"`
	FloatingIps        int `json:"floatingIps"`
	Gateways           int `json:"gateways"`
	SecurityGroups     int `json:"securityGroups"`
	SecurityGroupRules int `json:"securityGroupRules"`
}

type CloudProjectQuotaProfileLoadbalancerAPI struct {
	Loadbalancers  int `json:"loadbalancers"`
	Listeners      int `json:"listeners"`
	Pools          int `json:"pools"`
	Members        int `json:"members"`
	Healthmonitors int `json:"healthmonitors"`
	L7Policies     int `json:"l7Policies"`
	L7Rules        int `json:"l7Rules"`
}

type CloudProjectQuotaProfileKeyManagerAPI struct {
	Secrets    int `json:"secrets"`
	Containers int `json:"containers"`
}

type CloudProjectQuotaProfileShareAPI struct {
	Shares          int `json:"shares"`
	Gigabytes       int `json:"gigabytes"`
	Snapshots       int `json:"snapshots"`
	Backups         int `json:"backups"`
	BackupGigabytes int `json:"backupGigabytes"`
}

type CloudProjectQuotaProfileKeypairAPI struct {
	Keypairs int `json:"keypairs"`
}

type CloudProjectQuotaAvailableProfileAPI struct {
	Name         string                                  `json:"name"`
	Compute      CloudProjectQuotaProfileComputeAPI      `json:"compute"`
	BlockStorage CloudProjectQuotaProfileVolumeAPI       `json:"blockStorage"`
	Network      CloudProjectQuotaProfileNetworkAPI      `json:"network"`
	Loadbalancer CloudProjectQuotaProfileLoadbalancerAPI `json:"loadbalancer"`
	KeyManager   CloudProjectQuotaProfileKeyManagerAPI   `json:"keyManager"`
	Share        CloudProjectQuotaProfileShareAPI        `json:"share"`
	Keypair      CloudProjectQuotaProfileKeypairAPI      `json:"keypair"`
}

type CloudProjectQuotaRegionComputeAPI struct {
	Instances CloudProjectQuotaUsageAPI `json:"instances"`
	Cores     CloudProjectQuotaUsageAPI `json:"cores"`
	Memory    CloudProjectQuotaUsageAPI `json:"memory"`
}

type CloudProjectQuotaRegionVolumeAPI struct {
	Volumes         CloudProjectQuotaUsageAPI `json:"volumes"`
	Gigabytes       CloudProjectQuotaUsageAPI `json:"gigabytes"`
	Snapshots       CloudProjectQuotaUsageAPI `json:"snapshots"`
	Backups         CloudProjectQuotaUsageAPI `json:"backups"`
	BackupGigabytes CloudProjectQuotaUsageAPI `json:"backupGigabytes"`
	PerVolumeSize   CloudProjectQuotaLimitAPI `json:"perVolumeSize"`
}

type CloudProjectQuotaRegionNetworkAPI struct {
	Networks           CloudProjectQuotaUsageAPI `json:"networks"`
	Subnets            CloudProjectQuotaUsageAPI `json:"subnets"`
	FloatingIps        CloudProjectQuotaUsageAPI `json:"floatingIps"`
	Gateways           CloudProjectQuotaUsageAPI `json:"gateways"`
	SecurityGroups     CloudProjectQuotaUsageAPI `json:"securityGroups"`
	SecurityGroupRules CloudProjectQuotaUsageAPI `json:"securityGroupRules"`
}

type CloudProjectQuotaRegionLoadbalancerAPI struct {
	Loadbalancers  CloudProjectQuotaUsageAPI `json:"loadbalancers"`
	Listeners      CloudProjectQuotaUsageAPI `json:"listeners"`
	Pools          CloudProjectQuotaUsageAPI `json:"pools"`
	Members        CloudProjectQuotaUsageAPI `json:"members"`
	Healthmonitors CloudProjectQuotaUsageAPI `json:"healthmonitors"`
	L7Policies     CloudProjectQuotaUsageAPI `json:"l7Policies"`
	L7Rules        CloudProjectQuotaUsageAPI `json:"l7Rules"`
}

type CloudProjectQuotaRegionKeyManagerAPI struct {
	Secrets    CloudProjectQuotaUsageAPI `json:"secrets"`
	Containers CloudProjectQuotaUsageAPI `json:"containers"`
}

type CloudProjectQuotaRegionShareAPI struct {
	Shares            CloudProjectQuotaUsageAPI `json:"shares"`
	SizeTotal         CloudProjectQuotaUsageAPI `json:"sizeTotal"`
	Snapshots         CloudProjectQuotaUsageAPI `json:"snapshots"`
	SnapshotGigabytes CloudProjectQuotaUsageAPI `json:"snapshotGigabytes"`
	Backups           CloudProjectQuotaUsageAPI `json:"backups"`
	BackupGigabytes   CloudProjectQuotaUsageAPI `json:"backupGigabytes"`
	ShareNetworks     CloudProjectQuotaUsageAPI `json:"shareNetworks"`
	PerShareSize      CloudProjectQuotaLimitAPI `json:"perShareSize"`
}

type CloudProjectQuotaRegionKeypairAPI struct {
	Keypairs CloudProjectQuotaUsageAPI `json:"keypairs"`
}

type CloudProjectQuotaRegionAPI struct {
	Region       string                                  `json:"region"`
	Compute      *CloudProjectQuotaRegionComputeAPI      `json:"compute,omitempty"`
	Volume       *CloudProjectQuotaRegionVolumeAPI       `json:"volume,omitempty"`
	Network      *CloudProjectQuotaRegionNetworkAPI      `json:"network,omitempty"`
	Loadbalancer *CloudProjectQuotaRegionLoadbalancerAPI `json:"loadbalancer,omitempty"`
	KeyManager   *CloudProjectQuotaRegionKeyManagerAPI   `json:"keyManager,omitempty"`
	Share        *CloudProjectQuotaRegionShareAPI        `json:"share,omitempty"`
	Keypair      *CloudProjectQuotaRegionKeypairAPI      `json:"keypair,omitempty"`
}

type CloudProjectQuotaCurrentStateAPI struct {
	Profile           string                                 `json:"profile"`
	AvailableProfiles []CloudProjectQuotaAvailableProfileAPI `json:"availableProfiles,omitempty"`
	Regions           []CloudProjectQuotaRegionAPI           `json:"regions,omitempty"`
}

type CloudProjectQuotaAPIResponse struct {
	Id             string                            `json:"id"`
	ResourceStatus string                            `json:"resourceStatus"`
	Checksum       string                            `json:"checksum"`
	CreatedAt      string                            `json:"createdAt"`
	UpdatedAt      string                            `json:"updatedAt"`
	TargetSpec     *CloudProjectQuotaTargetSpecAPI   `json:"targetSpec,omitempty"`
	CurrentState   *CloudProjectQuotaCurrentStateAPI `json:"currentState,omitempty"`
}

// ------------------------------------------------------------------
// Attribute-type helpers
// ------------------------------------------------------------------

func quotaUsageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"limit": types.Int64Type,
		"used":  types.Int64Type,
		"unit":  ovhtypes.TfStringType{},
	}
}

func quotaLimitAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"limit": types.Int64Type,
		"unit":  ovhtypes.TfStringType{},
	}
}

func quotaTargetSpecAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"profile": ovhtypes.TfStringType{},
	}
}

func quotaProfileComputeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instances":            types.Int64Type,
		"cores":                types.Int64Type,
		"ram":                  types.Int64Type,
		"security_groups":      types.Int64Type,
		"security_group_rules": types.Int64Type,
		"server_groups":        types.Int64Type,
		"server_group_members": types.Int64Type,
	}
}

func quotaProfileVolumeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"volumes":          types.Int64Type,
		"gigabytes":        types.Int64Type,
		"snapshots":        types.Int64Type,
		"backups":          types.Int64Type,
		"backup_gigabytes": types.Int64Type,
	}
}

func quotaProfileNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"networks":             types.Int64Type,
		"subnets":              types.Int64Type,
		"floating_ips":         types.Int64Type,
		"gateways":             types.Int64Type,
		"security_groups":      types.Int64Type,
		"security_group_rules": types.Int64Type,
	}
}

func quotaProfileLoadbalancerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"loadbalancers":  types.Int64Type,
		"listeners":      types.Int64Type,
		"pools":          types.Int64Type,
		"members":        types.Int64Type,
		"healthmonitors": types.Int64Type,
		"l7_policies":    types.Int64Type,
		"l7_rules":       types.Int64Type,
	}
}

func quotaProfileKeyManagerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"secrets":    types.Int64Type,
		"containers": types.Int64Type,
	}
}

func quotaProfileShareAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"shares":           types.Int64Type,
		"gigabytes":        types.Int64Type,
		"snapshots":        types.Int64Type,
		"backups":          types.Int64Type,
		"backup_gigabytes": types.Int64Type,
	}
}

func quotaProfileKeypairAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keypairs": types.Int64Type,
	}
}

func quotaAvailableProfileAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"compute":       types.ObjectType{AttrTypes: quotaProfileComputeAttrTypes()},
		"block_storage": types.ObjectType{AttrTypes: quotaProfileVolumeAttrTypes()},
		"network":       types.ObjectType{AttrTypes: quotaProfileNetworkAttrTypes()},
		"loadbalancer":  types.ObjectType{AttrTypes: quotaProfileLoadbalancerAttrTypes()},
		"key_manager":   types.ObjectType{AttrTypes: quotaProfileKeyManagerAttrTypes()},
		"share":         types.ObjectType{AttrTypes: quotaProfileShareAttrTypes()},
		"keypair":       types.ObjectType{AttrTypes: quotaProfileKeypairAttrTypes()},
	}
}

func quotaAvailableProfileElementType() attr.Type {
	return types.ObjectType{AttrTypes: quotaAvailableProfileAttrTypes()}
}

func quotaRegionComputeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instances": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"cores":     types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"memory":    types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionVolumeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"volumes":          types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"gigabytes":        types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshots":        types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backups":          types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backup_gigabytes": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"per_volume_size":  types.ObjectType{AttrTypes: quotaLimitAttrTypes()},
	}
}

func quotaRegionNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"networks":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"subnets":              types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"floating_ips":         types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"gateways":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"security_groups":      types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"security_group_rules": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionLoadbalancerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"loadbalancers":  types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"listeners":      types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"pools":          types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"members":        types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"healthmonitors": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"l7_policies":    types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"l7_rules":       types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionKeyManagerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"secrets":    types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"containers": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionShareAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"shares":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"size_total":         types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshots":          types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshot_gigabytes": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backups":            types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backup_gigabytes":   types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"share_networks":     types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"per_share_size":     types.ObjectType{AttrTypes: quotaLimitAttrTypes()},
	}
}

func quotaRegionKeypairAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keypairs": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":       ovhtypes.TfStringType{},
		"compute":      types.ObjectType{AttrTypes: quotaRegionComputeAttrTypes()},
		"volume":       types.ObjectType{AttrTypes: quotaRegionVolumeAttrTypes()},
		"network":      types.ObjectType{AttrTypes: quotaRegionNetworkAttrTypes()},
		"loadbalancer": types.ObjectType{AttrTypes: quotaRegionLoadbalancerAttrTypes()},
		"key_manager":  types.ObjectType{AttrTypes: quotaRegionKeyManagerAttrTypes()},
		"share":        types.ObjectType{AttrTypes: quotaRegionShareAttrTypes()},
		"keypair":      types.ObjectType{AttrTypes: quotaRegionKeypairAttrTypes()},
	}
}

func quotaRegionElementType() attr.Type {
	return types.ObjectType{AttrTypes: quotaRegionAttrTypes()}
}

func quotaCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"profile":            ovhtypes.TfStringType{},
		"available_profiles": types.ListType{ElemType: quotaAvailableProfileElementType()},
		"regions":            types.ListType{ElemType: quotaRegionElementType()},
	}
}

// ------------------------------------------------------------------
// Builders — API -> Terraform values
// ------------------------------------------------------------------

func buildQuotaUsageObject(u CloudProjectQuotaUsageAPI) basetypes.ObjectValue {
	usedVal := types.Int64Null()
	if u.Used != nil {
		usedVal = types.Int64Value(int64(*u.Used))
	}
	obj, _ := types.ObjectValue(quotaUsageAttrTypes(), map[string]attr.Value{
		"limit": types.Int64Value(int64(u.Limit)),
		"used":  usedVal,
		"unit":  ovhtypes.TfStringValue{StringValue: types.StringValue(u.Unit)},
	})
	return obj
}

func buildQuotaLimitObject(l CloudProjectQuotaLimitAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaLimitAttrTypes(), map[string]attr.Value{
		"limit": types.Int64Value(int64(l.Limit)),
		"unit":  ovhtypes.TfStringValue{StringValue: types.StringValue(l.Unit)},
	})
	return obj
}

func buildQuotaTargetSpecObject(t *CloudProjectQuotaTargetSpecAPI) basetypes.ObjectValue {
	if t == nil {
		return types.ObjectNull(quotaTargetSpecAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaTargetSpecAttrTypes(), map[string]attr.Value{
		"profile": ovhtypes.TfStringValue{StringValue: types.StringValue(t.Profile)},
	})
	return obj
}

func buildQuotaProfileComputeObject(p CloudProjectQuotaProfileComputeAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileComputeAttrTypes(), map[string]attr.Value{
		"instances":            types.Int64Value(int64(p.Instances)),
		"cores":                types.Int64Value(int64(p.Cores)),
		"ram":                  types.Int64Value(int64(p.Ram)),
		"security_groups":      types.Int64Value(int64(p.SecurityGroups)),
		"security_group_rules": types.Int64Value(int64(p.SecurityGroupRules)),
		"server_groups":        types.Int64Value(int64(p.ServerGroups)),
		"server_group_members": types.Int64Value(int64(p.ServerGroupMembers)),
	})
	return obj
}

func buildQuotaProfileVolumeObject(p CloudProjectQuotaProfileVolumeAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileVolumeAttrTypes(), map[string]attr.Value{
		"volumes":          types.Int64Value(int64(p.Volumes)),
		"gigabytes":        types.Int64Value(int64(p.Gigabytes)),
		"snapshots":        types.Int64Value(int64(p.Snapshots)),
		"backups":          types.Int64Value(int64(p.Backups)),
		"backup_gigabytes": types.Int64Value(int64(p.BackupGigabytes)),
	})
	return obj
}

func buildQuotaProfileNetworkObject(p CloudProjectQuotaProfileNetworkAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileNetworkAttrTypes(), map[string]attr.Value{
		"networks":             types.Int64Value(int64(p.Networks)),
		"subnets":              types.Int64Value(int64(p.Subnets)),
		"floating_ips":         types.Int64Value(int64(p.FloatingIps)),
		"gateways":             types.Int64Value(int64(p.Gateways)),
		"security_groups":      types.Int64Value(int64(p.SecurityGroups)),
		"security_group_rules": types.Int64Value(int64(p.SecurityGroupRules)),
	})
	return obj
}

func buildQuotaProfileLoadbalancerObject(p CloudProjectQuotaProfileLoadbalancerAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileLoadbalancerAttrTypes(), map[string]attr.Value{
		"loadbalancers":  types.Int64Value(int64(p.Loadbalancers)),
		"listeners":      types.Int64Value(int64(p.Listeners)),
		"pools":          types.Int64Value(int64(p.Pools)),
		"members":        types.Int64Value(int64(p.Members)),
		"healthmonitors": types.Int64Value(int64(p.Healthmonitors)),
		"l7_policies":    types.Int64Value(int64(p.L7Policies)),
		"l7_rules":       types.Int64Value(int64(p.L7Rules)),
	})
	return obj
}

func buildQuotaProfileKeyManagerObject(p CloudProjectQuotaProfileKeyManagerAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileKeyManagerAttrTypes(), map[string]attr.Value{
		"secrets":    types.Int64Value(int64(p.Secrets)),
		"containers": types.Int64Value(int64(p.Containers)),
	})
	return obj
}

func buildQuotaProfileShareObject(p CloudProjectQuotaProfileShareAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileShareAttrTypes(), map[string]attr.Value{
		"shares":           types.Int64Value(int64(p.Shares)),
		"gigabytes":        types.Int64Value(int64(p.Gigabytes)),
		"snapshots":        types.Int64Value(int64(p.Snapshots)),
		"backups":          types.Int64Value(int64(p.Backups)),
		"backup_gigabytes": types.Int64Value(int64(p.BackupGigabytes)),
	})
	return obj
}

func buildQuotaProfileKeypairObject(p CloudProjectQuotaProfileKeypairAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaProfileKeypairAttrTypes(), map[string]attr.Value{
		"keypairs": types.Int64Value(int64(p.Keypairs)),
	})
	return obj
}

func buildQuotaAvailableProfileObject(p CloudProjectQuotaAvailableProfileAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaAvailableProfileAttrTypes(), map[string]attr.Value{
		"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(p.Name)},
		"compute":       buildQuotaProfileComputeObject(p.Compute),
		"block_storage": buildQuotaProfileVolumeObject(p.BlockStorage),
		"network":       buildQuotaProfileNetworkObject(p.Network),
		"loadbalancer":  buildQuotaProfileLoadbalancerObject(p.Loadbalancer),
		"key_manager":   buildQuotaProfileKeyManagerObject(p.KeyManager),
		"share":         buildQuotaProfileShareObject(p.Share),
		"keypair":       buildQuotaProfileKeypairObject(p.Keypair),
	})
	return obj
}

func buildQuotaAvailableProfilesList(profiles []CloudProjectQuotaAvailableProfileAPI) basetypes.ListValue {
	if profiles == nil {
		return types.ListNull(quotaAvailableProfileElementType())
	}
	elems := make([]attr.Value, len(profiles))
	for i, p := range profiles {
		elems[i] = buildQuotaAvailableProfileObject(p)
	}
	val, _ := types.ListValue(quotaAvailableProfileElementType(), elems)
	return val
}

func buildQuotaRegionComputeObject(c *CloudProjectQuotaRegionComputeAPI) basetypes.ObjectValue {
	if c == nil {
		return types.ObjectNull(quotaRegionComputeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionComputeAttrTypes(), map[string]attr.Value{
		"instances": buildQuotaUsageObject(c.Instances),
		"cores":     buildQuotaUsageObject(c.Cores),
		"memory":    buildQuotaUsageObject(c.Memory),
	})
	return obj
}

func buildQuotaRegionVolumeObject(v *CloudProjectQuotaRegionVolumeAPI) basetypes.ObjectValue {
	if v == nil {
		return types.ObjectNull(quotaRegionVolumeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionVolumeAttrTypes(), map[string]attr.Value{
		"volumes":          buildQuotaUsageObject(v.Volumes),
		"gigabytes":        buildQuotaUsageObject(v.Gigabytes),
		"snapshots":        buildQuotaUsageObject(v.Snapshots),
		"backups":          buildQuotaUsageObject(v.Backups),
		"backup_gigabytes": buildQuotaUsageObject(v.BackupGigabytes),
		"per_volume_size":  buildQuotaLimitObject(v.PerVolumeSize),
	})
	return obj
}

func buildQuotaRegionNetworkObject(n *CloudProjectQuotaRegionNetworkAPI) basetypes.ObjectValue {
	if n == nil {
		return types.ObjectNull(quotaRegionNetworkAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionNetworkAttrTypes(), map[string]attr.Value{
		"networks":             buildQuotaUsageObject(n.Networks),
		"subnets":              buildQuotaUsageObject(n.Subnets),
		"floating_ips":         buildQuotaUsageObject(n.FloatingIps),
		"gateways":             buildQuotaUsageObject(n.Gateways),
		"security_groups":      buildQuotaUsageObject(n.SecurityGroups),
		"security_group_rules": buildQuotaUsageObject(n.SecurityGroupRules),
	})
	return obj
}

func buildQuotaRegionLoadbalancerObject(l *CloudProjectQuotaRegionLoadbalancerAPI) basetypes.ObjectValue {
	if l == nil {
		return types.ObjectNull(quotaRegionLoadbalancerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionLoadbalancerAttrTypes(), map[string]attr.Value{
		"loadbalancers":  buildQuotaUsageObject(l.Loadbalancers),
		"listeners":      buildQuotaUsageObject(l.Listeners),
		"pools":          buildQuotaUsageObject(l.Pools),
		"members":        buildQuotaUsageObject(l.Members),
		"healthmonitors": buildQuotaUsageObject(l.Healthmonitors),
		"l7_policies":    buildQuotaUsageObject(l.L7Policies),
		"l7_rules":       buildQuotaUsageObject(l.L7Rules),
	})
	return obj
}

func buildQuotaRegionKeyManagerObject(k *CloudProjectQuotaRegionKeyManagerAPI) basetypes.ObjectValue {
	if k == nil {
		return types.ObjectNull(quotaRegionKeyManagerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionKeyManagerAttrTypes(), map[string]attr.Value{
		"secrets":    buildQuotaUsageObject(k.Secrets),
		"containers": buildQuotaUsageObject(k.Containers),
	})
	return obj
}

func buildQuotaRegionShareObject(s *CloudProjectQuotaRegionShareAPI) basetypes.ObjectValue {
	if s == nil {
		return types.ObjectNull(quotaRegionShareAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionShareAttrTypes(), map[string]attr.Value{
		"shares":             buildQuotaUsageObject(s.Shares),
		"size_total":         buildQuotaUsageObject(s.SizeTotal),
		"snapshots":          buildQuotaUsageObject(s.Snapshots),
		"snapshot_gigabytes": buildQuotaUsageObject(s.SnapshotGigabytes),
		"backups":            buildQuotaUsageObject(s.Backups),
		"backup_gigabytes":   buildQuotaUsageObject(s.BackupGigabytes),
		"share_networks":     buildQuotaUsageObject(s.ShareNetworks),
		"per_share_size":     buildQuotaLimitObject(s.PerShareSize),
	})
	return obj
}

func buildQuotaRegionKeypairObject(k *CloudProjectQuotaRegionKeypairAPI) basetypes.ObjectValue {
	if k == nil {
		return types.ObjectNull(quotaRegionKeypairAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionKeypairAttrTypes(), map[string]attr.Value{
		"keypairs": buildQuotaUsageObject(k.Keypairs),
	})
	return obj
}

func buildQuotaRegionObject(r CloudProjectQuotaRegionAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaRegionAttrTypes(), map[string]attr.Value{
		"region":       ovhtypes.TfStringValue{StringValue: types.StringValue(r.Region)},
		"compute":      buildQuotaRegionComputeObject(r.Compute),
		"volume":       buildQuotaRegionVolumeObject(r.Volume),
		"network":      buildQuotaRegionNetworkObject(r.Network),
		"loadbalancer": buildQuotaRegionLoadbalancerObject(r.Loadbalancer),
		"key_manager":  buildQuotaRegionKeyManagerObject(r.KeyManager),
		"share":        buildQuotaRegionShareObject(r.Share),
		"keypair":      buildQuotaRegionKeypairObject(r.Keypair),
	})
	return obj
}

func buildQuotaRegionsList(regions []CloudProjectQuotaRegionAPI) basetypes.ListValue {
	if regions == nil {
		return types.ListNull(quotaRegionElementType())
	}
	elems := make([]attr.Value, len(regions))
	for i, r := range regions {
		elems[i] = buildQuotaRegionObject(r)
	}
	val, _ := types.ListValue(quotaRegionElementType(), elems)
	return val
}

func buildQuotaCurrentStateObject(state *CloudProjectQuotaCurrentStateAPI) basetypes.ObjectValue {
	if state == nil {
		return types.ObjectNull(quotaCurrentStateAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaCurrentStateAttrTypes(), map[string]attr.Value{
		"profile":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Profile)},
		"available_profiles": buildQuotaAvailableProfilesList(state.AvailableProfiles),
		"regions":            buildQuotaRegionsList(state.Regions),
	})
	return obj
}

// MergeWith copies API response fields into the Terraform model.
func (m *CloudProjectQuotaModel) MergeWith(_ context.Context, response *CloudProjectQuotaAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.TargetSpec = buildQuotaTargetSpecObject(response.TargetSpec)
	m.CurrentState = buildQuotaCurrentStateObject(response.CurrentState)
}
