package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CloudManagedDatabaseServiceModel is the Terraform model of a managed database service, shared by
// the resource and the data source. The root attributes are the targetSpec of the service, and
// current_state is what runs.
type CloudManagedDatabaseServiceModel struct {
	ServiceName           types.String `tfsdk:"service_name"`
	Engine                types.String `tfsdk:"engine"`
	Region                types.String `tfsdk:"region"`
	Version               types.String `tfsdk:"version"`
	Plan                  types.String `tfsdk:"plan"`
	Flavor                types.String `tfsdk:"flavor"`
	NodesCount            types.Int64  `tfsdk:"nodes_count"`
	Description           types.String `tfsdk:"description"`
	DeletionProtected     types.Bool   `tfsdk:"deletion_protected"`
	PrometheusEnabled     types.Bool   `tfsdk:"prometheus_enabled"`
	AdvancedConfiguration types.Map    `tfsdk:"advanced_configuration"`
	Backup                types.Object `tfsdk:"backup"`
	MaintenanceWindow     types.Object `tfsdk:"maintenance_window"`
	Network               types.Object `tfsdk:"network"`
	Storage               types.Object `tfsdk:"storage"`

	Id             types.String `tfsdk:"id"`
	Checksum       types.String `tfsdk:"checksum"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
	ResourceStatus types.String `tfsdk:"resource_status"`
	CurrentTasks   types.List   `tfsdk:"current_tasks"`
	CurrentState   types.Object `tfsdk:"current_state"`
}

// cloudManagedDatabaseServiceResourceModel adds the arguments that only the creation takes, which
// the API doesn't return.
type cloudManagedDatabaseServiceResourceModel struct {
	CloudManagedDatabaseServiceModel
	ForkSourceServiceId   types.String `tfsdk:"fork_source_service_id"`
	ForkSourceBackupId    types.String `tfsdk:"fork_source_backup_id"`
	ForkSourcePointInTime types.String `tfsdk:"fork_source_point_in_time"`
}

// API types mirror the publicCloud.database models, which every engine shares. Their tfsdk tags map
// them onto the nested attributes of the same shape.

type managedDatabaseBackupSpec struct {
	Regions []string `json:"regions,omitempty" tfsdk:"regions"`
	Time    *string  `json:"time,omitempty" tfsdk:"time"`
}

type managedDatabaseMaintenanceWindow struct {
	StartTimeUtc *string `json:"startTimeUtc,omitempty" tfsdk:"start_time_utc"`
}

type managedDatabaseAuthorizedIpsRule struct {
	Cidr        string  `json:"cidr" tfsdk:"cidr"`
	Description *string `json:"description,omitempty" tfsdk:"description"`
}

type managedDatabaseAuthorizedIps struct {
	// No omitempty: an empty list is how every rule is removed.
	Rules []managedDatabaseAuthorizedIpsRule `json:"rules" tfsdk:"rules"`
}

type managedDatabaseNetwork struct {
	NetworkId     *string                       `json:"networkId,omitempty" tfsdk:"network_id"`
	SubnetId      *string                       `json:"subnetId,omitempty" tfsdk:"subnet_id"`
	AuthorizedIps *managedDatabaseAuthorizedIps `json:"authorizedIps,omitempty" tfsdk:"authorized_ips"`
	Public        *bool                         `json:"public,omitempty" tfsdk:"public"`
}

type managedDatabaseStorageSize struct {
	Unit  string `json:"unit" tfsdk:"unit"`
	Value int64  `json:"value" tfsdk:"value"`
}

type managedDatabaseStorageSpec struct {
	Size *managedDatabaseStorageSize `json:"size,omitempty" tfsdk:"size"`
}

// managedDatabaseServiceSpec is the targetSpec returned by the API, which the POST embeds.
type managedDatabaseServiceSpec struct {
	AdvancedConfiguration *map[string]string                `json:"advancedConfiguration,omitempty"`
	Backup                *managedDatabaseBackupSpec        `json:"backup,omitempty"`
	DeletionProtected     *bool                             `json:"deletionProtected,omitempty"`
	Description           *string                           `json:"description,omitempty"`
	Flavor                *string                           `json:"flavor,omitempty"`
	MaintenanceWindow     *managedDatabaseMaintenanceWindow `json:"maintenanceWindow,omitempty"`
	Network               *managedDatabaseNetwork           `json:"network,omitempty"`
	NodesCount            *int64                            `json:"nodesCount,omitempty"`
	Plan                  *string                           `json:"plan,omitempty"`
	PrometheusEnabled     *bool                             `json:"prometheusEnabled,omitempty"`
	Storage               *managedDatabaseStorageSpec       `json:"storage,omitempty"`
	Version               *string                           `json:"version,omitempty"`
}

type managedDatabaseServiceCreationSpec struct {
	managedDatabaseServiceSpec
	ForkSourceBackupId    *string `json:"forkSourceBackupId,omitempty"`
	ForkSourcePointInTime *string `json:"forkSourcePointInTime,omitempty"`
	ForkSourceServiceId   *string `json:"forkSourceServiceId,omitempty"`
	Region                string  `json:"region"`
}

type managedDatabaseServiceCreation struct {
	TargetSpec managedDatabaseServiceCreationSpec `json:"targetSpec"`
}

// managedDatabaseServiceReplace is the body of a PUT, which replaces the whole targetSpec: every field
// is sent, null where the API allows it.
type managedDatabaseServiceReplace struct {
	TargetSpec managedDatabaseServiceReplaceSpec `json:"targetSpec"`
}

type managedDatabaseServiceReplaceSpec struct {
	AdvancedConfiguration map[string]string                 `json:"advancedConfiguration"`
	Backup                *managedDatabaseBackupReplace     `json:"backup"`
	DeletionProtected     bool                              `json:"deletionProtected"`
	Description           *string                           `json:"description"`
	Flavor                string                            `json:"flavor"`
	MaintenanceWindow     *managedDatabaseMaintenanceWindow `json:"maintenanceWindow"`
	Network               managedDatabaseNetworkReplace     `json:"network"`
	NodesCount            int64                             `json:"nodesCount"`
	Plan                  string                            `json:"plan"`
	PrometheusEnabled     *bool                             `json:"prometheusEnabled"`
	Storage               *managedDatabaseStorageSpec       `json:"storage"`
	Version               string                            `json:"version"`
}

type managedDatabaseBackupReplace struct {
	Regions []string `json:"regions"`
	Time    *string  `json:"time"`
}

// managedDatabaseNetworkReplace only has the authorized IPs: the PUT changes neither the private
// network nor the public exposure.
type managedDatabaseNetworkReplace struct {
	AuthorizedIps managedDatabaseAuthorizedIpsReplace `json:"authorizedIps"`
}

type managedDatabaseAuthorizedIpsReplace struct {
	Rules []managedDatabaseAuthorizedIpsRuleReplace `json:"rules"`
}

type managedDatabaseAuthorizedIpsRuleReplace struct {
	Cidr        string  `json:"cidr"`
	Description *string `json:"description"`
}

type managedDatabaseBackupState struct {
	Pitr          *string  `json:"pitr" tfsdk:"pitr"`
	Regions       []string `json:"regions" tfsdk:"regions"`
	RetentionDays *int64   `json:"retentionDays" tfsdk:"retention_days"`
	Time          *string  `json:"time" tfsdk:"time"`
}

type managedDatabaseStorageState struct {
	Size *managedDatabaseStorageSize `json:"size" tfsdk:"size"`
	Type *string                     `json:"type" tfsdk:"type"`
}

type managedDatabaseEndpoint struct {
	Component string  `json:"component" tfsdk:"component"`
	Domain    string  `json:"domain" tfsdk:"domain"`
	Path      *string `json:"path" tfsdk:"path"`
	Port      *int64  `json:"port" tfsdk:"port"`
	Public    *bool   `json:"public" tfsdk:"public"`
	Scheme    *string `json:"scheme" tfsdk:"scheme"`
	Ssl       bool    `json:"ssl" tfsdk:"ssl"`
	SslMode   *string `json:"sslMode" tfsdk:"ssl_mode"`
	Uri       *string `json:"uri" tfsdk:"uri"`
}

type managedDatabaseNodeLocation struct {
	AvailabilityZone *string `json:"availabilityZone" tfsdk:"availability_zone"`
	Region           string  `json:"region" tfsdk:"region"`
}

type managedDatabaseNode struct {
	CreatedAt string                       `json:"createdAt" tfsdk:"created_at"`
	Flavor    string                       `json:"flavor" tfsdk:"flavor"`
	Id        string                       `json:"id" tfsdk:"id"`
	Location  *managedDatabaseNodeLocation `json:"location" tfsdk:"location"`
	Name      string                       `json:"name" tfsdk:"name"`
	Port      int64                        `json:"port" tfsdk:"port"`
	Role      *string                      `json:"role" tfsdk:"role"`
	Status    string                       `json:"status" tfsdk:"status"`
}

type managedDatabaseAuthorizedIpsRuleState struct {
	Cidr        string  `json:"cidr" tfsdk:"cidr"`
	CreatedAt   *string `json:"createdAt" tfsdk:"created_at"`
	Description *string `json:"description" tfsdk:"description"`
}

type managedDatabaseAuthorizedIpsState struct {
	Rules []managedDatabaseAuthorizedIpsRuleState `json:"rules" tfsdk:"rules"`
}

type managedDatabaseNetworkState struct {
	AuthorizedIps *managedDatabaseAuthorizedIpsState `json:"authorizedIps" tfsdk:"authorized_ips"`
	NetworkId     *string                            `json:"networkId" tfsdk:"network_id"`
	Public        *bool                              `json:"public" tfsdk:"public"`
	SubnetId      *string                            `json:"subnetId" tfsdk:"subnet_id"`
}

// managedDatabaseCapabilityActions tells, for one capability of the service, which actions are
// enabled.
type managedDatabaseCapabilityActions struct {
	Create *string `json:"create" tfsdk:"create"`
	Delete *string `json:"delete" tfsdk:"delete"`
	Read   *string `json:"read" tfsdk:"read"`
	Update *string `json:"update" tfsdk:"update"`
}

type managedDatabaseAdvancedConfigurationProperty struct {
	Default         *string  `json:"default" tfsdk:"default"`
	Description     string   `json:"description" tfsdk:"description"`
	Maximum         *float64 `json:"maximum" tfsdk:"maximum"`
	Minimum         *float64 `json:"minimum" tfsdk:"minimum"`
	Name            string   `json:"name" tfsdk:"name"`
	RequiresRestart *bool    `json:"requiresRestart" tfsdk:"requires_restart"`
	Type            string   `json:"type" tfsdk:"type"`
	Values          []string `json:"values" tfsdk:"values"`
}

type managedDatabaseIntegrationParameter struct {
	Name string `json:"name" tfsdk:"name"`
	Type string `json:"type" tfsdk:"type"`
}

type managedDatabaseIntegrationCapability struct {
	DestinationEngine string                                `json:"destinationEngine" tfsdk:"destination_engine"`
	Parameters        []managedDatabaseIntegrationParameter `json:"parameters" tfsdk:"parameters"`
	SourceEngine      string                                `json:"sourceEngine" tfsdk:"source_engine"`
	Type              string                                `json:"type" tfsdk:"type"`
}

// managedDatabaseEngineCapabilities is the publicCloud.database.ServiceEngineCapabilities model,
// without the connectors of Kafka Connect, which come with its other objects.
type managedDatabaseEngineCapabilities struct {
	AdvancedConfigurationProperties []managedDatabaseAdvancedConfigurationProperty `json:"advancedConfigurationProperties" tfsdk:"advanced_configuration_properties"`
	BackupRegions                   []string                                       `json:"backupRegions" tfsdk:"backup_regions"`
	Categories                      []string                                       `json:"categories" tfsdk:"categories"`
	Commands                        []string                                       `json:"commands" tfsdk:"commands"`
	Integrations                    []managedDatabaseIntegrationCapability         `json:"integrations" tfsdk:"integrations"`
}

type managedDatabaseTaskError struct {
	Message string `tfsdk:"message"`
}

// managedDatabaseTask is a CloudResourceTask, in the shape of the current_tasks attribute.
type managedDatabaseTask struct {
	Errors []managedDatabaseTaskError `tfsdk:"errors"`
	Id     string                     `tfsdk:"id"`
	Link   string                     `tfsdk:"link"`
	Status string                     `tfsdk:"status"`
	Type   string                     `tfsdk:"type"`
}

type managedDatabaseServiceState struct {
	AdvancedConfiguration map[string]string                           `json:"advancedConfiguration" tfsdk:"advanced_configuration"`
	Backup                *managedDatabaseBackupState                 `json:"backup" tfsdk:"backup"`
	Capabilities          map[string]managedDatabaseCapabilityActions `json:"capabilities" tfsdk:"capabilities"`
	DeletionProtected     *bool                                       `json:"deletionProtected" tfsdk:"deletion_protected"`
	Description           *string                                     `json:"description" tfsdk:"description"`
	Endpoints             []managedDatabaseEndpoint                   `json:"endpoints" tfsdk:"endpoints"`
	EngineCapabilities    *managedDatabaseEngineCapabilities          `json:"engineCapabilities" tfsdk:"engine_capabilities"`
	Flavor                *string                                     `json:"flavor" tfsdk:"flavor"`
	MaintenanceWindow     *managedDatabaseMaintenanceWindow           `json:"maintenanceWindow" tfsdk:"maintenance_window"`
	Network               *managedDatabaseNetworkState                `json:"network" tfsdk:"network"`
	Nodes                 []managedDatabaseNode                       `json:"nodes" tfsdk:"nodes"`
	NodesCount            *int64                                      `json:"nodesCount" tfsdk:"nodes_count"`
	Plan                  *string                                     `json:"plan" tfsdk:"plan"`
	PrometheusEnabled     *bool                                       `json:"prometheusEnabled" tfsdk:"prometheus_enabled"`
	Status                *string                                     `json:"status" tfsdk:"status"`
	Storage               *managedDatabaseStorageState                `json:"storage" tfsdk:"storage"`
	Version               *string                                     `json:"version" tfsdk:"version"`
}

type CloudManagedDatabaseServiceAPIResponse struct {
	Id             string                       `json:"id"`
	Checksum       string                       `json:"checksum"`
	CreatedAt      string                       `json:"createdAt"`
	UpdatedAt      string                       `json:"updatedAt"`
	ResourceStatus string                       `json:"resourceStatus"`
	TargetSpec     *managedDatabaseServiceSpec  `json:"targetSpec"`
	CurrentState   *managedDatabaseServiceState `json:"currentState"`
	CurrentTasks   []CloudResourceTask          `json:"currentTasks"`
}

func managedDatabaseBackupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"regions": types.ListType{ElemType: types.StringType},
		"time":    types.StringType,
	}
}

func managedDatabaseMaintenanceWindowAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"start_time_utc": types.StringType,
	}
}

func managedDatabaseNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"network_id": types.StringType,
		"subnet_id":  types.StringType,
		"authorized_ips": types.ObjectType{AttrTypes: map[string]attr.Type{
			"rules": types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"cidr":        types.StringType,
				"description": types.StringType,
			}}},
		}},
		"public": types.BoolType,
	}
}

func managedDatabaseStorageSizeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"unit":  types.StringType,
		"value": types.Int64Type,
	}
}

func managedDatabaseStorageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"size": types.ObjectType{AttrTypes: managedDatabaseStorageSizeAttrTypes()},
	}
}

func managedDatabaseNetworkStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"network_id": types.StringType,
		"subnet_id":  types.StringType,
		"authorized_ips": types.ObjectType{AttrTypes: map[string]attr.Type{
			"rules": types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"cidr":        types.StringType,
				"created_at":  types.StringType,
				"description": types.StringType,
			}}},
		}},
		"public": types.BoolType,
	}
}

func managedDatabaseTaskAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"errors": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"message": types.StringType,
		}}},
		"id":     types.StringType,
		"link":   types.StringType,
		"status": types.StringType,
		"type":   types.StringType,
	}
}

func managedDatabaseCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"advanced_configuration": types.MapType{ElemType: types.StringType},
		"backup": types.ObjectType{AttrTypes: map[string]attr.Type{
			"pitr":           types.StringType,
			"regions":        types.ListType{ElemType: types.StringType},
			"retention_days": types.Int64Type,
			"time":           types.StringType,
		}},
		"capabilities": types.MapType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"create": types.StringType,
			"delete": types.StringType,
			"read":   types.StringType,
			"update": types.StringType,
		}}},
		"engine_capabilities": types.ObjectType{AttrTypes: map[string]attr.Type{
			"advanced_configuration_properties": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"default":          types.StringType,
				"description":      types.StringType,
				"maximum":          types.Float64Type,
				"minimum":          types.Float64Type,
				"name":             types.StringType,
				"requires_restart": types.BoolType,
				"type":             types.StringType,
				"values":           types.ListType{ElemType: types.StringType},
			}}},
			"backup_regions": types.ListType{ElemType: types.StringType},
			"categories":     types.ListType{ElemType: types.StringType},
			"commands":       types.ListType{ElemType: types.StringType},
			"integrations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"destination_engine": types.StringType,
				"parameters": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"name": types.StringType,
					"type": types.StringType,
				}}},
				"source_engine": types.StringType,
				"type":          types.StringType,
			}}},
		}},
		"deletion_protected": types.BoolType,
		"description":        types.StringType,
		"endpoints": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"component": types.StringType,
			"domain":    types.StringType,
			"path":      types.StringType,
			"port":      types.Int64Type,
			"public":    types.BoolType,
			"scheme":    types.StringType,
			"ssl":       types.BoolType,
			"ssl_mode":  types.StringType,
			"uri":       types.StringType,
		}}},
		"flavor":             types.StringType,
		"maintenance_window": types.ObjectType{AttrTypes: managedDatabaseMaintenanceWindowAttrTypes()},
		"network":            types.ObjectType{AttrTypes: managedDatabaseNetworkStateAttrTypes()},
		"nodes": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"created_at": types.StringType,
			"flavor":     types.StringType,
			"id":         types.StringType,
			"location": types.ObjectType{AttrTypes: map[string]attr.Type{
				"availability_zone": types.StringType,
				"region":            types.StringType,
			}},
			"name":   types.StringType,
			"port":   types.Int64Type,
			"role":   types.StringType,
			"status": types.StringType,
		}}},
		"nodes_count":        types.Int64Type,
		"plan":               types.StringType,
		"prometheus_enabled": types.BoolType,
		"status":             types.StringType,
		"storage": types.ObjectType{AttrTypes: map[string]attr.Type{
			"size": types.ObjectType{AttrTypes: managedDatabaseStorageSizeAttrTypes()},
			"type": types.StringType,
		}},
		"version": types.StringType,
	}
}

// managedDatabaseObjectAs decodes a known object, leaving its null and unknown attributes out of the
// API body.
func managedDatabaseObjectAs[T any](ctx context.Context, obj types.Object, diags *diag.Diagnostics) *T {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	target := new(T)
	diags.Append(obj.As(ctx, target, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})...)
	return target
}

func managedDatabaseObjectFrom[T any](ctx context.Context, attrTypes map[string]attr.Type, value *T, diags *diag.Diagnostics) types.Object {
	if value == nil {
		return types.ObjectNull(attrTypes)
	}
	obj, d := types.ObjectValueFrom(ctx, attrTypes, value)
	diags.Append(d...)
	return obj
}

func knownStringPointer(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueStringPointer()
}

func knownInt64Pointer(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueInt64Pointer()
}

func knownBoolPointer(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueBoolPointer()
}

func firstNonNil[T any](values ...*T) *T {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}

// ToCreate builds the POST body. The creation doesn't take the advanced configuration.
func (m *CloudManagedDatabaseServiceModel) ToCreate(ctx context.Context) (*managedDatabaseServiceCreation, diag.Diagnostics) {
	var diags diag.Diagnostics
	spec := managedDatabaseServiceCreationSpec{
		Region: m.Region.ValueString(),
		managedDatabaseServiceSpec: managedDatabaseServiceSpec{
			Backup:            managedDatabaseObjectAs[managedDatabaseBackupSpec](ctx, m.Backup, &diags),
			DeletionProtected: knownBoolPointer(m.DeletionProtected),
			Description:       knownStringPointer(m.Description),
			Flavor:            knownStringPointer(m.Flavor),
			MaintenanceWindow: managedDatabaseObjectAs[managedDatabaseMaintenanceWindow](ctx, m.MaintenanceWindow, &diags),
			Network:           managedDatabaseObjectAs[managedDatabaseNetwork](ctx, m.Network, &diags),
			NodesCount:        knownInt64Pointer(m.NodesCount),
			Plan:              knownStringPointer(m.Plan),
			PrometheusEnabled: knownBoolPointer(m.PrometheusEnabled),
			Storage:           managedDatabaseObjectAs[managedDatabaseStorageSpec](ctx, m.Storage, &diags),
			Version:           knownStringPointer(m.Version),
		},
	}
	return &managedDatabaseServiceCreation{TargetSpec: spec}, diags
}

func (m *cloudManagedDatabaseServiceResourceModel) ToCreate(ctx context.Context) (*managedDatabaseServiceCreation, diag.Diagnostics) {
	payload, diags := m.CloudManagedDatabaseServiceModel.ToCreate(ctx)
	payload.TargetSpec.ForkSourceServiceId = knownStringPointer(m.ForkSourceServiceId)
	payload.TargetSpec.ForkSourceBackupId = knownStringPointer(m.ForkSourceBackupId)
	payload.TargetSpec.ForkSourcePointInTime = knownStringPointer(m.ForkSourcePointInTime)
	return payload, diags
}

// ToPut builds the PUT body from the whole targetSpec of the model.
func (m *CloudManagedDatabaseServiceModel) ToPut(ctx context.Context) (*managedDatabaseServiceReplace, diag.Diagnostics) {
	var diags diag.Diagnostics
	spec := managedDatabaseServiceReplaceSpec{
		DeletionProtected: m.DeletionProtected.ValueBool(),
		Description:       knownStringPointer(m.Description),
		Flavor:            m.Flavor.ValueString(),
		MaintenanceWindow: managedDatabaseObjectAs[managedDatabaseMaintenanceWindow](ctx, m.MaintenanceWindow, &diags),
		Network:           managedDatabaseNetworkReplace{AuthorizedIps: managedDatabaseAuthorizedIpsReplace{Rules: []managedDatabaseAuthorizedIpsRuleReplace{}}},
		NodesCount:        m.NodesCount.ValueInt64(),
		Plan:              m.Plan.ValueString(),
		PrometheusEnabled: knownBoolPointer(m.PrometheusEnabled),
		Storage:           managedDatabaseObjectAs[managedDatabaseStorageSpec](ctx, m.Storage, &diags),
		Version:           m.Version.ValueString(),
	}
	if !m.AdvancedConfiguration.IsNull() && !m.AdvancedConfiguration.IsUnknown() {
		spec.AdvancedConfiguration = map[string]string{}
		diags.Append(m.AdvancedConfiguration.ElementsAs(ctx, &spec.AdvancedConfiguration, false)...)
	}
	if backup := managedDatabaseObjectAs[managedDatabaseBackupSpec](ctx, m.Backup, &diags); backup != nil {
		spec.Backup = &managedDatabaseBackupReplace{Regions: backup.Regions, Time: backup.Time}
	}
	if network := managedDatabaseObjectAs[managedDatabaseNetwork](ctx, m.Network, &diags); network != nil && network.AuthorizedIps != nil {
		for _, rule := range network.AuthorizedIps.Rules {
			spec.Network.AuthorizedIps.Rules = append(spec.Network.AuthorizedIps.Rules, managedDatabaseAuthorizedIpsRuleReplace{Cidr: rule.Cidr, Description: rule.Description})
		}
	}
	return &managedDatabaseServiceReplace{TargetSpec: spec}, diags
}

// MergeWith copies the API response into the model. Root attributes come from the targetSpec, or
// from the current state for the values the targetSpec leaves out.
func (m *CloudManagedDatabaseServiceModel) MergeWith(ctx context.Context, response *CloudManagedDatabaseServiceAPIResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Id = types.StringValue(response.Id)
	m.Checksum = types.StringValue(response.Checksum)
	m.CreatedAt = types.StringValue(response.CreatedAt)
	m.UpdatedAt = types.StringValue(response.UpdatedAt)
	m.ResourceStatus = types.StringValue(response.ResourceStatus)
	m.CurrentState = managedDatabaseObjectFrom(ctx, managedDatabaseCurrentStateAttrTypes(), response.CurrentState, &diags)

	tasks := make([]managedDatabaseTask, 0, len(response.CurrentTasks))
	for _, task := range response.CurrentTasks {
		errs := make([]managedDatabaseTaskError, 0, len(task.Errors))
		for _, e := range task.Errors {
			errs = append(errs, managedDatabaseTaskError{Message: e.Message})
		}
		tasks = append(tasks, managedDatabaseTask{Errors: errs, Id: task.Id, Link: task.Link, Status: task.Status, Type: task.Type})
	}
	var d diag.Diagnostics
	m.CurrentTasks, d = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()}, tasks)
	diags.Append(d...)

	ts, cs := response.TargetSpec, response.CurrentState
	if ts == nil {
		ts = &managedDatabaseServiceSpec{}
	}
	if cs == nil {
		cs = &managedDatabaseServiceState{}
	}

	m.Version = types.StringPointerValue(firstNonNil(ts.Version, cs.Version))
	m.Plan = types.StringPointerValue(firstNonNil(ts.Plan, cs.Plan))
	m.Flavor = types.StringPointerValue(firstNonNil(ts.Flavor, cs.Flavor))
	m.NodesCount = types.Int64PointerValue(firstNonNil(ts.NodesCount, cs.NodesCount))
	m.DeletionProtected = types.BoolPointerValue(firstNonNil(ts.DeletionProtected, cs.DeletionProtected))
	m.PrometheusEnabled = types.BoolPointerValue(firstNonNil(ts.PrometheusEnabled, cs.PrometheusEnabled))

	// A description left out of the configuration stays null when the API returns an empty one.
	description := ""
	if d := firstNonNil(ts.Description, cs.Description); d != nil {
		description = *d
	}
	if description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
		m.Description = types.StringValue(description)
	} else {
		m.Description = types.StringNull()
	}

	// The targetSpec doesn't return the region, which can't change: it is only read for an import or
	// the data source.
	if (m.Region.IsNull() || m.Region.IsUnknown()) && len(cs.Nodes) > 0 && cs.Nodes[0].Location != nil {
		m.Region = types.StringValue(cs.Nodes[0].Location.Region)
	}

	backup := ts.Backup
	if backup == nil && cs.Backup != nil {
		backup = &managedDatabaseBackupSpec{Regions: cs.Backup.Regions, Time: cs.Backup.Time}
	}
	m.Backup = managedDatabaseObjectFrom(ctx, managedDatabaseBackupAttrTypes(), backup, &diags)

	m.MaintenanceWindow = managedDatabaseObjectFrom(ctx, managedDatabaseMaintenanceWindowAttrTypes(), firstNonNil(ts.MaintenanceWindow, cs.MaintenanceWindow), &diags)

	// The targetSpec doesn't return the network and subnet either.
	var network *managedDatabaseNetwork
	if cs.Network != nil {
		network = &managedDatabaseNetwork{NetworkId: cs.Network.NetworkId, SubnetId: cs.Network.SubnetId, Public: cs.Network.Public}
		if ips := cs.Network.AuthorizedIps; ips != nil {
			rules := make([]managedDatabaseAuthorizedIpsRule, 0, len(ips.Rules))
			for _, rule := range ips.Rules {
				rules = append(rules, managedDatabaseAuthorizedIpsRule{Cidr: rule.Cidr, Description: rule.Description})
			}
			network.AuthorizedIps = &managedDatabaseAuthorizedIps{Rules: rules}
		}
	}
	if ts.Network != nil && (ts.Network.AuthorizedIps != nil || ts.Network.Public != nil) {
		if network == nil {
			network = &managedDatabaseNetwork{}
		}
		if ts.Network.AuthorizedIps != nil {
			network.AuthorizedIps = ts.Network.AuthorizedIps
		}
		if ts.Network.Public != nil {
			network.Public = ts.Network.Public
		}
	}
	m.Network = managedDatabaseObjectFrom(ctx, managedDatabaseNetworkAttrTypes(), network, &diags)

	var storage *managedDatabaseStorageSpec
	if ts.Storage != nil && ts.Storage.Size != nil {
		storage = ts.Storage
	} else if cs.Storage != nil && cs.Storage.Size != nil {
		storage = &managedDatabaseStorageSpec{Size: cs.Storage.Size}
	}
	m.Storage = managedDatabaseObjectFrom(ctx, managedDatabaseStorageAttrTypes(), storage, &diags)

	// The whole advanced configuration, which the PUT replaces, with the defaults the API adds. Engines
	// without one return null, which an empty configuration may come back as.
	switch {
	case ts.AdvancedConfiguration != nil:
		var d diag.Diagnostics
		m.AdvancedConfiguration, d = types.MapValueFrom(ctx, types.StringType, *ts.AdvancedConfiguration)
		diags.Append(d...)
	case m.AdvancedConfiguration.IsNull() || m.AdvancedConfiguration.IsUnknown() || len(m.AdvancedConfiguration.Elements()) > 0:
		m.AdvancedConfiguration = types.MapNull(types.StringType)
	}

	return diags
}
