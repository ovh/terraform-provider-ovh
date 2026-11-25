package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// VPSDataModel is a VPS as the ovh_vps data source and the elements of the
// ovh_vpss data source expose it: the target specification of the API v2
// contract (vps.server.ReadTargetSpec) flat at the root, plus location.region,
// and the current state nested in current_state.
type VPSDataModel struct {
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	DisplayName  ovhtypes.TfStringValue `tfsdk:"display_name"`
	PowerState   ovhtypes.TfStringValue `tfsdk:"power_state"`
	Location     types.Object           `tfsdk:"location"`
	Installation types.Object           `tfsdk:"installation"`
	Options      types.Object           `tfsdk:"options"`
	CurrentState types.Object           `tfsdk:"current_state"`
}

// VPSsDataModel is the model of the ovh_vpss data source.
type VPSsDataModel struct {
	Vps types.List `tfsdk:"vps"`
}

// Descriptions, taken from the API v2 contract (vps.server.*).
const (
	vpsDescId                  = "Identifier of the VPS, equal to its service name"
	vpsDescDisplayName         = "Name you gave to the VPS"
	vpsDescPowerState          = "Power state you requested: ACTIVE, RESCUE or SHUTOFF. Until the VPS API v2 is available, it is read from the current power state, and is null while the VPS is in a transitional state (installing, rebooting...)"
	vpsDescLocation            = "Location of the VPS"
	vpsDescRegion              = "OVHcloud region code, such as eu-west-gra. Null when the datacenter of the VPS matches no known region"
	vpsDescInstallation        = "Installation you requested"
	vpsDescInstallationFrom    = "Source the VPS is installed from. Until the VPS API v2 is available, image_id and image_name are both set from the image installed on the VPS, which current_state.installation.from holds"
	vpsDescImageId             = "Fixed image build"
	vpsDescImageName           = "Image name, such as \"Debian 12\". It follows the latest build: a new build does not reinstall the VPS"
	vpsDescOptions             = "Settings you requested for the options of the VPS, matching current_state.options"
	vpsDescAutoBackup          = "Schedule of the automated backup you requested. Null when you did not request one"
	vpsDescAutoBackupCron      = "Daily schedule of the automated backup, as a cron expression \"M H * * *\" in UTC, such as \"0 2 * * *\" for 2:00 AM UTC"
	vpsDescCurrentState        = "Current state of the VPS, read-only"
	vpsDescCsDisplayName       = "Name of the VPS"
	vpsDescCsPowerState        = "Observed power state of the VPS: ACTIVE, RESCUE, SHUTOFF, or UNKNOWN when it cannot be read"
	vpsDescCsInstallation      = "Current installation of the VPS"
	vpsDescCsInstallationFrom  = "Source the VPS is installed from. All the fields are null when the source is unknown, such as an image that was removed"
	vpsDescCsImageId           = "Identifier of the image build the VPS is installed from. Null when the source is unknown"
	vpsDescCsImageName         = "Name of the image the VPS is installed from, such as \"Debian 12\". Null when the source is unknown"
	vpsDescCsModel             = "Hardware model of the VPS"
	vpsDescCsVcpus             = "Number of virtual CPUs of the VPS"
	vpsDescCsRam               = "Amount of RAM of the VPS, in MB"
	vpsDescCsDisk              = "Size of the main disk of the VPS, in GB"
	vpsDescCsIpAddresses       = "IP addresses of the VPS"
	vpsDescCsIp                = "IP address"
	vpsDescCsIpVersion         = "IP version: 4 for IPv4, 6 for IPv6"
	vpsDescCsOptions           = "Options of the VPS and their current settings. They are subscribed and terminated through an order"
	vpsDescCsAutoBackup        = "Automated backup option"
	vpsDescCsAutoBackupEnabled = "Whether you subscribed to the automated backup option"
	vpsDescCsAutoBackupCron    = "Daily schedule applied to the automated backup, as a cron expression \"M H * * *\" in UTC. Null when you did not subscribe to the option or when the automated backup is disabled"
	vpsDescCsAutoBackupRotate  = "Number of automated backups kept. Null when you did not subscribe to the option"
	vpsDescCsSnapshot          = "Manual snapshot option"
	vpsDescCsSnapshotEnabled   = "Whether you subscribed to the manual snapshot option"
	vpsDescCsAdditionalDisks   = "Additional disk option"
	vpsDescCsAdditionalEnabled = "Whether you subscribed to the additional disk option"
	vpsDescCsAdditionalSize    = "Size of the additional disk, in GB. Null when you did not subscribe to the option"
)

// vpsDataSourceAttributes is the schema of a VPS, shared by ovh_vps (id
// required) and the elements of ovh_vpss (id computed).
func vpsDataSourceAttributes(idRequired bool) map[string]schema.Attribute {
	id := schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: vpsDescId, MarkdownDescription: vpsDescId}
	if idRequired {
		id.Computed = false
		id.Required = true
	}

	return map[string]schema.Attribute{
		"id":           id,
		"display_name": vpsStringAttribute(vpsDescDisplayName),
		"power_state":  vpsStringAttribute(vpsDescPowerState),
		"location": schema.SingleNestedAttribute{Computed: true, Description: vpsDescLocation, MarkdownDescription: vpsDescLocation, Attributes: map[string]schema.Attribute{
			"region": vpsStringAttribute(vpsDescRegion),
		}},
		"installation": schema.SingleNestedAttribute{Computed: true, Description: vpsDescInstallation, MarkdownDescription: vpsDescInstallation, Attributes: map[string]schema.Attribute{
			"from": schema.SingleNestedAttribute{Computed: true, Description: vpsDescInstallationFrom, MarkdownDescription: vpsDescInstallationFrom, Attributes: map[string]schema.Attribute{
				"image_id":   vpsStringAttribute(vpsDescImageId),
				"image_name": vpsStringAttribute(vpsDescImageName),
			}},
		}},
		"options": schema.SingleNestedAttribute{Computed: true, Description: vpsDescOptions, MarkdownDescription: vpsDescOptions, Attributes: map[string]schema.Attribute{
			"auto_backup": schema.SingleNestedAttribute{Computed: true, Description: vpsDescAutoBackup, MarkdownDescription: vpsDescAutoBackup, Attributes: map[string]schema.Attribute{
				"cron": vpsStringAttribute(vpsDescAutoBackupCron),
			}},
		}},
		"current_state": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCurrentState, MarkdownDescription: vpsDescCurrentState, Attributes: vpsCurrentStateAttributes()},
	}
}

func vpsCurrentStateAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"display_name": vpsStringAttribute(vpsDescCsDisplayName),
		"power_state":  vpsStringAttribute(vpsDescCsPowerState),
		"installation": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsInstallation, MarkdownDescription: vpsDescCsInstallation, Attributes: map[string]schema.Attribute{
			"from": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsInstallationFrom, MarkdownDescription: vpsDescCsInstallationFrom, Attributes: map[string]schema.Attribute{
				"image_id":   vpsStringAttribute(vpsDescCsImageId),
				"image_name": vpsStringAttribute(vpsDescCsImageName),
			}},
		}},
		"location": schema.SingleNestedAttribute{Computed: true, Description: vpsDescLocation, MarkdownDescription: vpsDescLocation, Attributes: map[string]schema.Attribute{
			"region": vpsStringAttribute(vpsDescRegion),
		}},
		"model": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsModel, MarkdownDescription: vpsDescCsModel, Attributes: map[string]schema.Attribute{
			"vcpus": vpsInt64Attribute(vpsDescCsVcpus),
			"ram":   vpsInt64Attribute(vpsDescCsRam),
			"disk":  vpsInt64Attribute(vpsDescCsDisk),
		}},
		"ip_addresses": schema.ListNestedAttribute{Computed: true, Description: vpsDescCsIpAddresses, MarkdownDescription: vpsDescCsIpAddresses, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"ip":      vpsStringAttribute(vpsDescCsIp),
			"version": vpsInt64Attribute(vpsDescCsIpVersion),
		}}},
		"options": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsOptions, MarkdownDescription: vpsDescCsOptions, Attributes: map[string]schema.Attribute{
			"auto_backup": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsAutoBackup, MarkdownDescription: vpsDescCsAutoBackup, Attributes: map[string]schema.Attribute{
				"enabled":  vpsBoolAttribute(vpsDescCsAutoBackupEnabled),
				"cron":     vpsStringAttribute(vpsDescCsAutoBackupCron),
				"rotation": vpsInt64Attribute(vpsDescCsAutoBackupRotate),
			}},
			"snapshot": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsSnapshot, MarkdownDescription: vpsDescCsSnapshot, Attributes: map[string]schema.Attribute{
				"enabled": vpsBoolAttribute(vpsDescCsSnapshotEnabled),
			}},
			"additional_disks": schema.SingleNestedAttribute{Computed: true, Description: vpsDescCsAdditionalDisks, MarkdownDescription: vpsDescCsAdditionalDisks, Attributes: map[string]schema.Attribute{
				"enabled": vpsBoolAttribute(vpsDescCsAdditionalEnabled),
				"size":    vpsInt64Attribute(vpsDescCsAdditionalSize),
			}},
		}},
	}
}

func vpsStringAttribute(description string) schema.StringAttribute {
	return schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: description, MarkdownDescription: description}
}

func vpsInt64Attribute(description string) schema.Int64Attribute {
	return schema.Int64Attribute{Computed: true, Description: description, MarkdownDescription: description}
}

func vpsBoolAttribute(description string) schema.BoolAttribute {
	return schema.BoolAttribute{Computed: true, Description: description, MarkdownDescription: description}
}

// Attribute types, matching vpsDataSourceAttributes.

func vpsDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"display_name":  ovhtypes.TfStringType{},
		"power_state":   ovhtypes.TfStringType{},
		"location":      types.ObjectType{AttrTypes: vpsLocationAttrTypes()},
		"installation":  types.ObjectType{AttrTypes: vpsInstallationAttrTypes()},
		"options":       types.ObjectType{AttrTypes: vpsOptionsAttrTypes()},
		"current_state": types.ObjectType{AttrTypes: vpsCurrentStateAttrTypes()},
	}
}

func vpsLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"region": ovhtypes.TfStringType{}}
}

func vpsInstallationFromAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"image_id":   ovhtypes.TfStringType{},
		"image_name": ovhtypes.TfStringType{},
	}
}

func vpsInstallationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"from": types.ObjectType{AttrTypes: vpsInstallationFromAttrTypes()}}
}

func vpsAutoBackupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"cron": ovhtypes.TfStringType{}}
}

func vpsOptionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"auto_backup": types.ObjectType{AttrTypes: vpsAutoBackupAttrTypes()}}
}

func vpsModelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vcpus": types.Int64Type,
		"ram":   types.Int64Type,
		"disk":  types.Int64Type,
	}
}

func vpsIPAddressAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":      ovhtypes.TfStringType{},
		"version": types.Int64Type,
	}
}

func vpsAutoBackupOptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":  types.BoolType,
		"cron":     ovhtypes.TfStringType{},
		"rotation": types.Int64Type,
	}
}

func vpsSnapshotOptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"enabled": types.BoolType}
}

func vpsAdditionalDisksOptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"size":    types.Int64Type,
	}
}

func vpsCurrentOptionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auto_backup":      types.ObjectType{AttrTypes: vpsAutoBackupOptionAttrTypes()},
		"snapshot":         types.ObjectType{AttrTypes: vpsSnapshotOptionAttrTypes()},
		"additional_disks": types.ObjectType{AttrTypes: vpsAdditionalDisksOptionAttrTypes()},
	}
}

func vpsCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": ovhtypes.TfStringType{},
		"power_state":  ovhtypes.TfStringType{},
		"installation": types.ObjectType{AttrTypes: vpsInstallationAttrTypes()},
		"location":     types.ObjectType{AttrTypes: vpsLocationAttrTypes()},
		"model":        types.ObjectType{AttrTypes: vpsModelAttrTypes()},
		"ip_addresses": types.ListType{ElemType: types.ObjectType{AttrTypes: vpsIPAddressAttrTypes()}},
		"options":      types.ObjectType{AttrTypes: vpsCurrentOptionsAttrTypes()},
	}
}

// MergeWith sets the model from a VPS read through the backend.
func (m *VPSDataModel) MergeWith(ctx context.Context, srv *VPSAPIServer) {
	m.Id = tfStr(srv.Id)
	m.DisplayName = tfStr(srv.TargetSpec.DisplayName)
	m.PowerState = vpsTfStringPtr(srv.TargetSpec.PowerState)
	m.Installation = vpsInstallationObject(srv.TargetSpec.Installation.From)

	m.Options = types.ObjectValueMust(vpsOptionsAttrTypes(), map[string]attr.Value{
		"auto_backup": types.ObjectNull(vpsAutoBackupAttrTypes()),
	})
	if autoBackup := srv.TargetSpec.Options.AutoBackup; autoBackup != nil {
		m.Options = types.ObjectValueMust(vpsOptionsAttrTypes(), map[string]attr.Value{
			"auto_backup": types.ObjectValueMust(vpsAutoBackupAttrTypes(), map[string]attr.Value{
				"cron": tfStr(autoBackup.Cron),
			}),
		})
	}

	// The location is no part of the target specification: it is set at the
	// order and read from the current state.
	if srv.CurrentState == nil {
		m.Location = types.ObjectNull(vpsLocationAttrTypes())
		m.CurrentState = types.ObjectNull(vpsCurrentStateAttrTypes())
		return
	}
	cs := srv.CurrentState
	m.Location = vpsLocationObject(cs.Location)

	ipAddresses := make([]attr.Value, 0, len(cs.IpAddresses))
	for _, ip := range cs.IpAddresses {
		ipAddresses = append(ipAddresses, types.ObjectValueMust(vpsIPAddressAttrTypes(), map[string]attr.Value{
			"ip":      tfStr(ip.Ip),
			"version": types.Int64Value(ip.Version),
		}))
	}

	m.CurrentState = types.ObjectValueMust(vpsCurrentStateAttrTypes(), map[string]attr.Value{
		"display_name": tfStr(cs.DisplayName),
		"power_state":  tfStr(cs.PowerState),
		"installation": vpsInstallationObject(cs.Installation.From),
		"location":     vpsLocationObject(cs.Location),
		"model": types.ObjectValueMust(vpsModelAttrTypes(), map[string]attr.Value{
			"vcpus": types.Int64Value(cs.Model.Vcpus),
			"ram":   types.Int64Value(cs.Model.Ram),
			"disk":  types.Int64Value(cs.Model.Disk),
		}),
		"ip_addresses": types.ListValueMust(types.ObjectType{AttrTypes: vpsIPAddressAttrTypes()}, ipAddresses),
		"options": types.ObjectValueMust(vpsCurrentOptionsAttrTypes(), map[string]attr.Value{
			"auto_backup": types.ObjectValueMust(vpsAutoBackupOptionAttrTypes(), map[string]attr.Value{
				"enabled":  types.BoolValue(cs.Options.AutoBackup.Enabled),
				"cron":     vpsTfStringPtr(cs.Options.AutoBackup.Cron),
				"rotation": types.Int64PointerValue(cs.Options.AutoBackup.Rotation),
			}),
			"snapshot": types.ObjectValueMust(vpsSnapshotOptionAttrTypes(), map[string]attr.Value{
				"enabled": types.BoolValue(cs.Options.Snapshot.Enabled),
			}),
			"additional_disks": types.ObjectValueMust(vpsAdditionalDisksOptionAttrTypes(), map[string]attr.Value{
				"enabled": types.BoolValue(cs.Options.AdditionalDisks.Enabled),
				"size":    types.Int64PointerValue(cs.Options.AdditionalDisks.Size),
			}),
		}),
	})
}

// newVPSsDataModel builds the ovh_vpss model from the VPS read through the
// backend.
func newVPSsDataModel(ctx context.Context, servers []VPSAPIServer) VPSsDataModel {
	elements := make([]attr.Value, 0, len(servers))
	for i := range servers {
		var m VPSDataModel
		m.MergeWith(ctx, &servers[i])
		elements = append(elements, types.ObjectValueMust(vpsDataAttrTypes(), map[string]attr.Value{
			"id":            m.Id,
			"display_name":  m.DisplayName,
			"power_state":   m.PowerState,
			"location":      m.Location,
			"installation":  m.Installation,
			"options":       m.Options,
			"current_state": m.CurrentState,
		}))
	}
	return VPSsDataModel{Vps: types.ListValueMust(types.ObjectType{AttrTypes: vpsDataAttrTypes()}, elements)}
}

func vpsInstallationObject(from VPSAPIInstallationFrom) basetypes.ObjectValue {
	return types.ObjectValueMust(vpsInstallationAttrTypes(), map[string]attr.Value{
		"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
			"image_id":   vpsTfStringPtr(from.ImageId),
			"image_name": vpsTfStringPtr(from.ImageName),
		}),
	})
}

func vpsLocationObject(location VPSAPILocation) basetypes.ObjectValue {
	return types.ObjectValueMust(vpsLocationAttrTypes(), map[string]attr.Value{
		"region": vpsTfStringPtr(location.Region),
	})
}

// vpsTfStringPtr returns a null string for nil.
func vpsTfStringPtr(s *string) ovhtypes.TfStringValue {
	if s == nil {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return tfStr(*s)
}
