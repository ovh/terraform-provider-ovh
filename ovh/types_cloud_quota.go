package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudQuotaModel is the Terraform model for the quota data source.
//
// Following the apiv2 convention used by other resources (e.g.
// ovh_cloud_instance, ovh_cloud_loadbalancer), the API's targetSpec envelope
// is flattened to top-level attributes (`regions`,
// `prevent_automatic_quota_upgrade`) and is NOT exposed as a nested
// `target_spec` block.
type CloudQuotaModel struct {
	// Input
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Flattened targetSpec (computed for the data source)
	PreventAutomaticQuotaUpgrade types.Bool `tfsdk:"prevent_automatic_quota_upgrade"`
	Regions                      types.List `tfsdk:"regions"`

	// Envelope (computed)
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// CloudQuotaResourceModel is the Terraform model for the writable quota
// resource. It mirrors CloudQuotaModel but omits the data-source-only
// `region` filter so that the resource schema and the model line up.
type CloudQuotaResourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`

	// Flattened targetSpec (required for the resource)
	PreventAutomaticQuotaUpgrade types.Bool `tfsdk:"prevent_automatic_quota_upgrade"`
	Regions                      types.List `tfsdk:"regions"`

	// Envelope (computed)
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// MergeWith populates the resource model from the API response.
//
// The quota envelope is a singleton whose targetSpec.regions lists EVERY region
// the project has a quota in, while this resource manages only the regions the
// user configured. The input `regions` attribute is therefore reconciled
// against the regions already present in the model (the prior state on Read, the
// planned regions on Create/Update): server regions are filtered down to the
// managed set and their profiles refreshed from the response. When the model
// has no regions yet (import), the full server list is adopted.
func (m *CloudQuotaResourceModel) MergeWith(ctx context.Context, response *CloudQuotaAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}

	if response.TargetSpec != nil {
		m.PreventAutomaticQuotaUpgrade = types.BoolValue(response.TargetSpec.PreventAutomaticQuotaUpgrade)
		m.Regions = reconcileManagedQuotaRegions(ctx, m.Regions, response.TargetSpec.Regions)
	} else {
		m.PreventAutomaticQuotaUpgrade = types.BoolNull()
		m.Regions = types.ListNull(quotaTargetSpecRegionElementType())
	}

	m.CurrentState = buildQuotaCurrentStateObject(response.CurrentState)
}

// reconcileManagedQuotaRegions narrows the singleton's full targetSpec.regions
// down to the set the resource manages.
//
//   - prior has regions  → keep only server regions whose region matches one in
//     prior, refreshing each profile from the server (drift detection). A managed
//     region absent from the server is dropped.
//   - prior is null/empty → adopt the full server list (import: the managed
//     subset is unknown).
func reconcileManagedQuotaRegions(ctx context.Context, prior basetypes.ListValue, serverRegions []CloudQuotaRegionTargetSpecAPI) basetypes.ListValue {
	managed := managedQuotaRegionNames(ctx, prior)
	if len(managed) == 0 {
		return buildQuotaTargetSpecRegionsList(serverRegions)
	}

	elems := make([]attr.Value, 0, len(managed))
	for _, r := range serverRegions {
		region := ""
		if r.Location != nil {
			region = r.Location.Region
		}
		if managed[region] {
			elems = append(elems, buildQuotaTargetSpecRegionObject(r))
		}
	}
	val, _ := types.ListValue(quotaTargetSpecRegionElementType(), elems)
	return val
}

// managedQuotaRegionNames returns the set of region names currently held in the
// `regions` list value (prior state or plan).
func managedQuotaRegionNames(ctx context.Context, regions basetypes.ListValue) map[string]bool {
	if regions.IsNull() || regions.IsUnknown() {
		return nil
	}
	var plan []quotaRegionPlan
	if diags := regions.ElementsAs(ctx, &plan, false); diags.HasError() {
		return nil
	}
	names := make(map[string]bool, len(plan))
	for _, p := range plan {
		names[p.Region.ValueString()] = true
	}
	return names
}

// ------------------------------------------------------------------
// API response structs — mirror the public-cloud-apiv2 model.Quota
// ------------------------------------------------------------------

type CloudQuotaLocationAPI struct {
	Region string `json:"region"`
}

type CloudQuotaUsageAPI struct {
	Limit int    `json:"limit"`
	Used  *int   `json:"used"`
	Unit  string `json:"unit"`
}

type CloudQuotaLimitAPI struct {
	Limit int    `json:"limit"`
	Unit  string `json:"unit"`
}

// ----- Profile (available_profiles) -----

type CloudQuotaProfileComputeAPI struct {
	Cores     int `json:"cores"`
	Instances int `json:"instances"`
	Memory    int `json:"memory"`
}

type CloudQuotaProfileVolumeAPI struct {
	BackupSizeTotal int `json:"backupSizeTotal"`
	Backups         int `json:"backups"`
	SizeTotal       int `json:"sizeTotal"`
	Snapshots       int `json:"snapshots"`
	Volumes         int `json:"volumes"`
}

type CloudQuotaProfileNetworkAPI struct {
	FloatingIps        int `json:"floatingIps"`
	Gateways           int `json:"gateways"`
	Networks           int `json:"networks"`
	SecurityGroupRules int `json:"securityGroupRules"`
	SecurityGroups     int `json:"securityGroups"`
	Subnets            int `json:"subnets"`
}

type CloudQuotaProfileLoadbalancerAPI struct {
	HealthMonitors int `json:"healthMonitors"`
	L7Policies     int `json:"l7Policies"`
	L7Rules        int `json:"l7Rules"`
	Listeners      int `json:"listeners"`
	Loadbalancers  int `json:"loadbalancers"`
	Members        int `json:"members"`
	Pools          int `json:"pools"`
}

type CloudQuotaProfileKeyManagerAPI struct {
	Containers int `json:"containers"`
	Secrets    int `json:"secrets"`
}

type CloudQuotaProfileShareAPI struct {
	BackupSizeTotal int `json:"backupSizeTotal"`
	Backups         int `json:"backups"`
	Shares          int `json:"shares"`
	SizeTotal       int `json:"sizeTotal"`
	Snapshots       int `json:"snapshots"`
}

type CloudQuotaProfileKeypairAPI struct {
	Keypairs int `json:"keypairs"`
}

type CloudQuotaAvailableProfileAPI struct {
	Name         string                            `json:"name"`
	Compute      *CloudQuotaProfileComputeAPI      `json:"compute,omitempty"`
	Volume       *CloudQuotaProfileVolumeAPI       `json:"volume,omitempty"`
	Network      *CloudQuotaProfileNetworkAPI      `json:"network,omitempty"`
	Loadbalancer *CloudQuotaProfileLoadbalancerAPI `json:"loadbalancer,omitempty"`
	KeyManager   *CloudQuotaProfileKeyManagerAPI   `json:"keyManager,omitempty"`
	Share        *CloudQuotaProfileShareAPI        `json:"share,omitempty"`
	Keypair      *CloudQuotaProfileKeypairAPI      `json:"keypair,omitempty"`
}

// ----- Region usage (current_state.regions[].usage.*) -----

type CloudQuotaRegionComputeAPI struct {
	Cores     CloudQuotaUsageAPI `json:"cores"`
	Instances CloudQuotaUsageAPI `json:"instances"`
	Memory    CloudQuotaUsageAPI `json:"memory"`
}

type CloudQuotaRegionVolumeAPI struct {
	BackupSizeTotal CloudQuotaUsageAPI `json:"backupSizeTotal"`
	Backups         CloudQuotaUsageAPI `json:"backups"`
	PerVolumeSize   CloudQuotaLimitAPI `json:"perVolumeSize"`
	SizeTotal       CloudQuotaUsageAPI `json:"sizeTotal"`
	Snapshots       CloudQuotaUsageAPI `json:"snapshots"`
	Volumes         CloudQuotaUsageAPI `json:"volumes"`
}

type CloudQuotaRegionNetworkAPI struct {
	FloatingIps        CloudQuotaUsageAPI `json:"floatingIps"`
	Gateways           CloudQuotaUsageAPI `json:"gateways"`
	Networks           CloudQuotaUsageAPI `json:"networks"`
	SecurityGroupRules CloudQuotaUsageAPI `json:"securityGroupRules"`
	SecurityGroups     CloudQuotaUsageAPI `json:"securityGroups"`
	Subnets            CloudQuotaUsageAPI `json:"subnets"`
}

type CloudQuotaRegionLoadbalancerAPI struct {
	HealthMonitors CloudQuotaUsageAPI `json:"healthMonitors"`
	L7Policies     CloudQuotaUsageAPI `json:"l7Policies"`
	L7Rules        CloudQuotaUsageAPI `json:"l7Rules"`
	Listeners      CloudQuotaUsageAPI `json:"listeners"`
	Loadbalancers  CloudQuotaUsageAPI `json:"loadbalancers"`
	Members        CloudQuotaUsageAPI `json:"members"`
	Pools          CloudQuotaUsageAPI `json:"pools"`
}

type CloudQuotaRegionKeyManagerAPI struct {
	Containers CloudQuotaUsageAPI `json:"containers"`
	Secrets    CloudQuotaUsageAPI `json:"secrets"`
}

type CloudQuotaRegionShareAPI struct {
	BackupSizeTotal   CloudQuotaUsageAPI `json:"backupSizeTotal"`
	Backups           CloudQuotaUsageAPI `json:"backups"`
	PerShareSize      CloudQuotaLimitAPI `json:"perShareSize"`
	ShareNetworks     CloudQuotaUsageAPI `json:"shareNetworks"`
	Shares            CloudQuotaUsageAPI `json:"shares"`
	SizeTotal         CloudQuotaUsageAPI `json:"sizeTotal"`
	SnapshotSizeTotal CloudQuotaUsageAPI `json:"snapshotSizeTotal"`
	Snapshots         CloudQuotaUsageAPI `json:"snapshots"`
}

type CloudQuotaRegionKeypairAPI struct {
	Keypairs CloudQuotaUsageAPI `json:"keypairs"`
}

type CloudQuotaUsageDetailsAPI struct {
	Compute      *CloudQuotaRegionComputeAPI      `json:"compute,omitempty"`
	KeyManager   *CloudQuotaRegionKeyManagerAPI   `json:"keyManager,omitempty"`
	Keypair      *CloudQuotaRegionKeypairAPI      `json:"keypair,omitempty"`
	Loadbalancer *CloudQuotaRegionLoadbalancerAPI `json:"loadbalancer,omitempty"`
	Network      *CloudQuotaRegionNetworkAPI      `json:"network,omitempty"`
	Share        *CloudQuotaRegionShareAPI        `json:"share,omitempty"`
	Volume       *CloudQuotaRegionVolumeAPI       `json:"volume,omitempty"`
}

type CloudQuotaRegionCurrentStateAPI struct {
	Location *CloudQuotaLocationAPI     `json:"location,omitempty"`
	Profile  string                     `json:"profile"`
	Usage    *CloudQuotaUsageDetailsAPI `json:"usage,omitempty"`
}

type CloudQuotaRegionTargetSpecAPI struct {
	Location *CloudQuotaLocationAPI `json:"location,omitempty"`
	Profile  string                 `json:"profile"`
}

type CloudQuotaTargetSpecAPI struct {
	PreventAutomaticQuotaUpgrade bool                            `json:"preventAutomaticQuotaUpgrade"`
	Regions                      []CloudQuotaRegionTargetSpecAPI `json:"regions,omitempty"`
}

type CloudQuotaCurrentStateAPI struct {
	PreventAutomaticQuotaUpgrade bool                              `json:"preventAutomaticQuotaUpgrade"`
	AvailableProfiles            []CloudQuotaAvailableProfileAPI   `json:"availableProfiles,omitempty"`
	Regions                      []CloudQuotaRegionCurrentStateAPI `json:"regions,omitempty"`
}

// ----- currentTasks (common.CurrentTask[]) -----

type CloudQuotaTaskErrorAPI struct {
	Message string `json:"message"`
}

type CloudQuotaCurrentTaskAPI struct {
	Id     string                   `json:"id"`
	Type   string                   `json:"type"`
	Status string                   `json:"status"`
	Link   string                   `json:"link"`
	Errors []CloudQuotaTaskErrorAPI `json:"errors,omitempty"`
}

type CloudQuotaAPIResponse struct {
	Id             string                     `json:"id"`
	ResourceStatus string                     `json:"resourceStatus"`
	Checksum       string                     `json:"checksum"`
	CreatedAt      string                     `json:"createdAt"`
	UpdatedAt      string                     `json:"updatedAt"`
	TargetSpec     *CloudQuotaTargetSpecAPI   `json:"targetSpec,omitempty"`
	CurrentState   *CloudQuotaCurrentStateAPI `json:"currentState,omitempty"`
	CurrentTasks   []CloudQuotaCurrentTaskAPI `json:"currentTasks,omitempty"`
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

// ----- flattened top-level `regions` (mirrors targetSpec.regions[]) -----

func quotaTargetSpecRegionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":  ovhtypes.TfStringType{},
		"profile": ovhtypes.TfStringType{},
	}
}

func quotaTargetSpecRegionElementType() attr.Type {
	return types.ObjectType{AttrTypes: quotaTargetSpecRegionAttrTypes()}
}

// ----- available_profiles -----

func quotaProfileComputeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cores":     types.Int64Type,
		"instances": types.Int64Type,
		"memory":    types.Int64Type,
	}
}

func quotaProfileVolumeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_size_total": types.Int64Type,
		"backups":           types.Int64Type,
		"size_total":        types.Int64Type,
		"snapshots":         types.Int64Type,
		"volumes":           types.Int64Type,
	}
}

func quotaProfileNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"floating_ips":         types.Int64Type,
		"gateways":             types.Int64Type,
		"networks":             types.Int64Type,
		"security_group_rules": types.Int64Type,
		"security_groups":      types.Int64Type,
		"subnets":              types.Int64Type,
	}
}

func quotaProfileLoadbalancerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_monitors": types.Int64Type,
		"l7_policies":     types.Int64Type,
		"l7_rules":        types.Int64Type,
		"listeners":       types.Int64Type,
		"loadbalancers":   types.Int64Type,
		"members":         types.Int64Type,
		"pools":           types.Int64Type,
	}
}

func quotaProfileKeyManagerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"containers": types.Int64Type,
		"secrets":    types.Int64Type,
	}
}

func quotaProfileShareAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_size_total": types.Int64Type,
		"backups":           types.Int64Type,
		"shares":            types.Int64Type,
		"size_total":        types.Int64Type,
		"snapshots":         types.Int64Type,
	}
}

func quotaProfileKeypairAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keypairs": types.Int64Type,
	}
}

func quotaAvailableProfileAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         ovhtypes.TfStringType{},
		"compute":      types.ObjectType{AttrTypes: quotaProfileComputeAttrTypes()},
		"volume":       types.ObjectType{AttrTypes: quotaProfileVolumeAttrTypes()},
		"network":      types.ObjectType{AttrTypes: quotaProfileNetworkAttrTypes()},
		"loadbalancer": types.ObjectType{AttrTypes: quotaProfileLoadbalancerAttrTypes()},
		"key_manager":  types.ObjectType{AttrTypes: quotaProfileKeyManagerAttrTypes()},
		"share":        types.ObjectType{AttrTypes: quotaProfileShareAttrTypes()},
		"keypair":      types.ObjectType{AttrTypes: quotaProfileKeypairAttrTypes()},
	}
}

func quotaAvailableProfileElementType() attr.Type {
	return types.ObjectType{AttrTypes: quotaAvailableProfileAttrTypes()}
}

// ----- regions[].usage -----

func quotaRegionComputeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cores":     types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"instances": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"memory":    types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionVolumeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_size_total": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backups":           types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"per_volume_size":   types.ObjectType{AttrTypes: quotaLimitAttrTypes()},
		"size_total":        types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshots":         types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"volumes":           types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"floating_ips":         types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"gateways":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"networks":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"security_group_rules": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"security_groups":      types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"subnets":              types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionLoadbalancerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_monitors": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"l7_policies":     types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"l7_rules":        types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"listeners":       types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"loadbalancers":   types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"members":         types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"pools":           types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionKeyManagerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"containers": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"secrets":    types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
	}
}

func quotaRegionShareAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"backup_size_total":   types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"backups":             types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"per_share_size":      types.ObjectType{AttrTypes: quotaLimitAttrTypes()},
		"share_networks":      types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"shares":              types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"size_total":          types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshot_size_total": types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
		"snapshots":           types.ObjectType{AttrTypes: quotaUsageAttrTypes()},
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
		"profile":      ovhtypes.TfStringType{},
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
		"prevent_automatic_quota_upgrade": types.BoolType,
		"available_profiles":              types.ListType{ElemType: quotaAvailableProfileElementType()},
		"regions":                         types.ListType{ElemType: quotaRegionElementType()},
	}
}

// ------------------------------------------------------------------
// Builders — API -> Terraform values
// ------------------------------------------------------------------

func buildQuotaUsageObject(u CloudQuotaUsageAPI) basetypes.ObjectValue {
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

func buildQuotaLimitObject(l CloudQuotaLimitAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaLimitAttrTypes(), map[string]attr.Value{
		"limit": types.Int64Value(int64(l.Limit)),
		"unit":  ovhtypes.TfStringValue{StringValue: types.StringValue(l.Unit)},
	})
	return obj
}

// ----- flattened top-level regions (mirrors targetSpec.regions[]) -----

func buildQuotaTargetSpecRegionObject(r CloudQuotaRegionTargetSpecAPI) basetypes.ObjectValue {
	region := ""
	if r.Location != nil {
		region = r.Location.Region
	}
	obj, _ := types.ObjectValue(quotaTargetSpecRegionAttrTypes(), map[string]attr.Value{
		"region":  ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
		"profile": ovhtypes.TfStringValue{StringValue: types.StringValue(r.Profile)},
	})
	return obj
}

func buildQuotaTargetSpecRegionsList(regions []CloudQuotaRegionTargetSpecAPI) basetypes.ListValue {
	if regions == nil {
		return types.ListNull(quotaTargetSpecRegionElementType())
	}
	elems := make([]attr.Value, len(regions))
	for i, r := range regions {
		elems[i] = buildQuotaTargetSpecRegionObject(r)
	}
	val, _ := types.ListValue(quotaTargetSpecRegionElementType(), elems)
	return val
}

// ----- available_profiles -----

func buildQuotaProfileComputeObject(p *CloudQuotaProfileComputeAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileComputeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileComputeAttrTypes(), map[string]attr.Value{
		"cores":     types.Int64Value(int64(p.Cores)),
		"instances": types.Int64Value(int64(p.Instances)),
		"memory":    types.Int64Value(int64(p.Memory)),
	})
	return obj
}

func buildQuotaProfileVolumeObject(p *CloudQuotaProfileVolumeAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileVolumeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileVolumeAttrTypes(), map[string]attr.Value{
		"backup_size_total": types.Int64Value(int64(p.BackupSizeTotal)),
		"backups":           types.Int64Value(int64(p.Backups)),
		"size_total":        types.Int64Value(int64(p.SizeTotal)),
		"snapshots":         types.Int64Value(int64(p.Snapshots)),
		"volumes":           types.Int64Value(int64(p.Volumes)),
	})
	return obj
}

func buildQuotaProfileNetworkObject(p *CloudQuotaProfileNetworkAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileNetworkAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileNetworkAttrTypes(), map[string]attr.Value{
		"floating_ips":         types.Int64Value(int64(p.FloatingIps)),
		"gateways":             types.Int64Value(int64(p.Gateways)),
		"networks":             types.Int64Value(int64(p.Networks)),
		"security_group_rules": types.Int64Value(int64(p.SecurityGroupRules)),
		"security_groups":      types.Int64Value(int64(p.SecurityGroups)),
		"subnets":              types.Int64Value(int64(p.Subnets)),
	})
	return obj
}

func buildQuotaProfileLoadbalancerObject(p *CloudQuotaProfileLoadbalancerAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileLoadbalancerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileLoadbalancerAttrTypes(), map[string]attr.Value{
		"health_monitors": types.Int64Value(int64(p.HealthMonitors)),
		"l7_policies":     types.Int64Value(int64(p.L7Policies)),
		"l7_rules":        types.Int64Value(int64(p.L7Rules)),
		"listeners":       types.Int64Value(int64(p.Listeners)),
		"loadbalancers":   types.Int64Value(int64(p.Loadbalancers)),
		"members":         types.Int64Value(int64(p.Members)),
		"pools":           types.Int64Value(int64(p.Pools)),
	})
	return obj
}

func buildQuotaProfileKeyManagerObject(p *CloudQuotaProfileKeyManagerAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileKeyManagerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileKeyManagerAttrTypes(), map[string]attr.Value{
		"containers": types.Int64Value(int64(p.Containers)),
		"secrets":    types.Int64Value(int64(p.Secrets)),
	})
	return obj
}

func buildQuotaProfileShareObject(p *CloudQuotaProfileShareAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileShareAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileShareAttrTypes(), map[string]attr.Value{
		"backup_size_total": types.Int64Value(int64(p.BackupSizeTotal)),
		"backups":           types.Int64Value(int64(p.Backups)),
		"shares":            types.Int64Value(int64(p.Shares)),
		"size_total":        types.Int64Value(int64(p.SizeTotal)),
		"snapshots":         types.Int64Value(int64(p.Snapshots)),
	})
	return obj
}

func buildQuotaProfileKeypairObject(p *CloudQuotaProfileKeypairAPI) basetypes.ObjectValue {
	if p == nil {
		return types.ObjectNull(quotaProfileKeypairAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaProfileKeypairAttrTypes(), map[string]attr.Value{
		"keypairs": types.Int64Value(int64(p.Keypairs)),
	})
	return obj
}

func buildQuotaAvailableProfileObject(p CloudQuotaAvailableProfileAPI) basetypes.ObjectValue {
	obj, _ := types.ObjectValue(quotaAvailableProfileAttrTypes(), map[string]attr.Value{
		"name":         ovhtypes.TfStringValue{StringValue: types.StringValue(p.Name)},
		"compute":      buildQuotaProfileComputeObject(p.Compute),
		"volume":       buildQuotaProfileVolumeObject(p.Volume),
		"network":      buildQuotaProfileNetworkObject(p.Network),
		"loadbalancer": buildQuotaProfileLoadbalancerObject(p.Loadbalancer),
		"key_manager":  buildQuotaProfileKeyManagerObject(p.KeyManager),
		"share":        buildQuotaProfileShareObject(p.Share),
		"keypair":      buildQuotaProfileKeypairObject(p.Keypair),
	})
	return obj
}

func buildQuotaAvailableProfilesList(profiles []CloudQuotaAvailableProfileAPI) basetypes.ListValue {
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

// ----- regions[].usage -----

func buildQuotaRegionComputeObject(c *CloudQuotaRegionComputeAPI) basetypes.ObjectValue {
	if c == nil {
		return types.ObjectNull(quotaRegionComputeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionComputeAttrTypes(), map[string]attr.Value{
		"cores":     buildQuotaUsageObject(c.Cores),
		"instances": buildQuotaUsageObject(c.Instances),
		"memory":    buildQuotaUsageObject(c.Memory),
	})
	return obj
}

func buildQuotaRegionVolumeObject(v *CloudQuotaRegionVolumeAPI) basetypes.ObjectValue {
	if v == nil {
		return types.ObjectNull(quotaRegionVolumeAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionVolumeAttrTypes(), map[string]attr.Value{
		"backup_size_total": buildQuotaUsageObject(v.BackupSizeTotal),
		"backups":           buildQuotaUsageObject(v.Backups),
		"per_volume_size":   buildQuotaLimitObject(v.PerVolumeSize),
		"size_total":        buildQuotaUsageObject(v.SizeTotal),
		"snapshots":         buildQuotaUsageObject(v.Snapshots),
		"volumes":           buildQuotaUsageObject(v.Volumes),
	})
	return obj
}

func buildQuotaRegionNetworkObject(n *CloudQuotaRegionNetworkAPI) basetypes.ObjectValue {
	if n == nil {
		return types.ObjectNull(quotaRegionNetworkAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionNetworkAttrTypes(), map[string]attr.Value{
		"floating_ips":         buildQuotaUsageObject(n.FloatingIps),
		"gateways":             buildQuotaUsageObject(n.Gateways),
		"networks":             buildQuotaUsageObject(n.Networks),
		"security_group_rules": buildQuotaUsageObject(n.SecurityGroupRules),
		"security_groups":      buildQuotaUsageObject(n.SecurityGroups),
		"subnets":              buildQuotaUsageObject(n.Subnets),
	})
	return obj
}

func buildQuotaRegionLoadbalancerObject(l *CloudQuotaRegionLoadbalancerAPI) basetypes.ObjectValue {
	if l == nil {
		return types.ObjectNull(quotaRegionLoadbalancerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionLoadbalancerAttrTypes(), map[string]attr.Value{
		"health_monitors": buildQuotaUsageObject(l.HealthMonitors),
		"l7_policies":     buildQuotaUsageObject(l.L7Policies),
		"l7_rules":        buildQuotaUsageObject(l.L7Rules),
		"listeners":       buildQuotaUsageObject(l.Listeners),
		"loadbalancers":   buildQuotaUsageObject(l.Loadbalancers),
		"members":         buildQuotaUsageObject(l.Members),
		"pools":           buildQuotaUsageObject(l.Pools),
	})
	return obj
}

func buildQuotaRegionKeyManagerObject(k *CloudQuotaRegionKeyManagerAPI) basetypes.ObjectValue {
	if k == nil {
		return types.ObjectNull(quotaRegionKeyManagerAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionKeyManagerAttrTypes(), map[string]attr.Value{
		"containers": buildQuotaUsageObject(k.Containers),
		"secrets":    buildQuotaUsageObject(k.Secrets),
	})
	return obj
}

func buildQuotaRegionShareObject(s *CloudQuotaRegionShareAPI) basetypes.ObjectValue {
	if s == nil {
		return types.ObjectNull(quotaRegionShareAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionShareAttrTypes(), map[string]attr.Value{
		"backup_size_total":   buildQuotaUsageObject(s.BackupSizeTotal),
		"backups":             buildQuotaUsageObject(s.Backups),
		"per_share_size":      buildQuotaLimitObject(s.PerShareSize),
		"share_networks":      buildQuotaUsageObject(s.ShareNetworks),
		"shares":              buildQuotaUsageObject(s.Shares),
		"size_total":          buildQuotaUsageObject(s.SizeTotal),
		"snapshot_size_total": buildQuotaUsageObject(s.SnapshotSizeTotal),
		"snapshots":           buildQuotaUsageObject(s.Snapshots),
	})
	return obj
}

func buildQuotaRegionKeypairObject(k *CloudQuotaRegionKeypairAPI) basetypes.ObjectValue {
	if k == nil {
		return types.ObjectNull(quotaRegionKeypairAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaRegionKeypairAttrTypes(), map[string]attr.Value{
		"keypairs": buildQuotaUsageObject(k.Keypairs),
	})
	return obj
}

func buildQuotaRegionObject(r CloudQuotaRegionCurrentStateAPI) basetypes.ObjectValue {
	region := ""
	if r.Location != nil {
		region = r.Location.Region
	}

	var compute, volume, network, loadbalancer, keyManager, share, keypair basetypes.ObjectValue
	if r.Usage != nil {
		compute = buildQuotaRegionComputeObject(r.Usage.Compute)
		volume = buildQuotaRegionVolumeObject(r.Usage.Volume)
		network = buildQuotaRegionNetworkObject(r.Usage.Network)
		loadbalancer = buildQuotaRegionLoadbalancerObject(r.Usage.Loadbalancer)
		keyManager = buildQuotaRegionKeyManagerObject(r.Usage.KeyManager)
		share = buildQuotaRegionShareObject(r.Usage.Share)
		keypair = buildQuotaRegionKeypairObject(r.Usage.Keypair)
	} else {
		compute = types.ObjectNull(quotaRegionComputeAttrTypes())
		volume = types.ObjectNull(quotaRegionVolumeAttrTypes())
		network = types.ObjectNull(quotaRegionNetworkAttrTypes())
		loadbalancer = types.ObjectNull(quotaRegionLoadbalancerAttrTypes())
		keyManager = types.ObjectNull(quotaRegionKeyManagerAttrTypes())
		share = types.ObjectNull(quotaRegionShareAttrTypes())
		keypair = types.ObjectNull(quotaRegionKeypairAttrTypes())
	}

	obj, _ := types.ObjectValue(quotaRegionAttrTypes(), map[string]attr.Value{
		"region":       ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
		"profile":      ovhtypes.TfStringValue{StringValue: types.StringValue(r.Profile)},
		"compute":      compute,
		"volume":       volume,
		"network":      network,
		"loadbalancer": loadbalancer,
		"key_manager":  keyManager,
		"share":        share,
		"keypair":      keypair,
	})
	return obj
}

func buildQuotaRegionsList(regions []CloudQuotaRegionCurrentStateAPI) basetypes.ListValue {
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

func buildQuotaCurrentStateObject(state *CloudQuotaCurrentStateAPI) basetypes.ObjectValue {
	if state == nil {
		return types.ObjectNull(quotaCurrentStateAttrTypes())
	}
	obj, _ := types.ObjectValue(quotaCurrentStateAttrTypes(), map[string]attr.Value{
		"prevent_automatic_quota_upgrade": types.BoolValue(state.PreventAutomaticQuotaUpgrade),
		"available_profiles":              buildQuotaAvailableProfilesList(state.AvailableProfiles),
		"regions":                         buildQuotaRegionsList(state.Regions),
	})
	return obj
}

// MergeWith copies API response fields into the Terraform model. The flattened
// targetSpec attributes (`prevent_automatic_quota_upgrade`, `regions`) are
// populated from response.TargetSpec to match the apiv2 convention used by
// other resources like ovh_cloud_instance.
func (m *CloudQuotaModel) MergeWith(_ context.Context, response *CloudQuotaAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}

	if response.TargetSpec != nil {
		m.PreventAutomaticQuotaUpgrade = types.BoolValue(response.TargetSpec.PreventAutomaticQuotaUpgrade)
		m.Regions = buildQuotaTargetSpecRegionsList(response.TargetSpec.Regions)
	} else {
		m.PreventAutomaticQuotaUpgrade = types.BoolNull()
		m.Regions = types.ListNull(quotaTargetSpecRegionElementType())
	}

	m.CurrentState = buildQuotaCurrentStateObject(response.CurrentState)
}
