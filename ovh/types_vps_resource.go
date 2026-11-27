package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// VPSResourceModel is the model of the ovh_vps resource: the root of the
// ovh_vps data source (VPSDataModel), with the installation settings of the
// API v2 contract (vps.server.UpdateInstallation), plus the order parameters
// the contract does not hold (plan_code, location.region, order) and timeouts.
type VPSResourceModel struct {
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	PlanCode     ovhtypes.TfStringValue `tfsdk:"plan_code"`
	DisplayName  ovhtypes.TfStringValue `tfsdk:"display_name"`
	PowerState   ovhtypes.TfStringValue `tfsdk:"power_state"`
	Location     types.Object           `tfsdk:"location"`
	Installation types.Object           `tfsdk:"installation"`
	Options      types.Object           `tfsdk:"options"`
	Order        types.Object           `tfsdk:"order"`
	Timeouts     types.Object           `tfsdk:"timeouts"`
	CurrentState types.Object           `tfsdk:"current_state"`
}

// Defaults of the resource.
const (
	vpsDefaultPowerState  = "ACTIVE"
	vpsDefaultDuration    = "P1M"
	vpsDefaultPricingMode = "default"

	vpsCreateTimeout = 60 * time.Minute
	vpsUpdateTimeout = 60 * time.Minute
	vpsDeleteTimeout = 30 * time.Minute
)

// Descriptions of the resource. The ones the data source shares are in
// types_vps_data.go; these differ from them, or are absent from the data source.
const (
	vpsReinstallDescription       = "Changing it reinstalls the VPS, which erases its disk"
	vpsDescPlanCode               = "Plan code of the VPS offer, such as vps-2025-model3, from the VPS order catalog. " + vpsRequiresReplaceDescription
	vpsDescResourceRegion         = "OVHcloud region code, such as eu-west-gra, among the regions of GET /v2/location. " + vpsRequiresReplaceDescription
	vpsDescResourcePowerState     = "Power state you request: ACTIVE, RESCUE or SHUTOFF. Defaults to ACTIVE. Until the VPS API v2 is available, it is read from the current power state, and keeps its value while the VPS is in a transitional state (installing, rebooting...)"
	vpsDescResourceInstallFrom    = "Source to install the VPS from. Set exactly one of image_id and image_name. After an import or the migration of the state, both are read from the image installed on the VPS. " + vpsReinstallDescription
	vpsDescSshKey                 = "Public SSH key to inject, in OpenSSH public key format. " + vpsReinstallDescription
	vpsDescPostInstallationScript = "Script run after the installation. Write-only: it is never read back. " + vpsReinstallDescription
	vpsDescSendPasswordByEmail    = "Send the generated root password to you by email. Default: true"
	vpsDescOrder                  = "Order of the VPS. Only used when the VPS is ordered"
	vpsDescOrderDuration          = "Duration of the order, in ISO 8601 format, such as P1M for one month. Defaults to P1M"
	vpsDescOrderPricingMode       = "Pricing mode of the order. Defaults to default"
	vpsDescOrderId                = "Identifier of the order of the VPS"
	vpsDescOrderOptions           = "Plan codes of the options ordered with the VPS, such as the automated backup or an additional disk, from the addons of plan_code in the VPS order catalog. Only used when the VPS is ordered: the options of the VPS are read in current_state.options"
	vpsDescTimeouts               = "How long to wait for the VPS, as durations such as 1h30m"
	vpsDescTimeoutCreate          = "Wait for the creation of the VPS, from the order to its delivery. Defaults to 60m"
	vpsDescTimeoutUpdate          = "Wait for the update of the VPS. Defaults to 60m"
	vpsDescTimeoutDelete          = "Wait for the termination of the VPS. Defaults to 30m"
)

// vpsResourceAttributes is the schema of the ovh_vps resource.
func vpsResourceAttributes() map[string]schema.Attribute {
	str := func(description string) schema.StringAttribute {
		return schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Description: description, MarkdownDescription: description}
	}
	object := func(description string, attributes map[string]schema.Attribute) schema.SingleNestedAttribute {
		return schema.SingleNestedAttribute{Description: description, MarkdownDescription: description, Attributes: attributes}
	}
	timeout := func(description string) schema.StringAttribute {
		a := str(description)
		a.Optional = true
		a.Validators = []validator.String{managedDatabaseDuration}
		return a
	}

	id := str(vpsDescId)
	id.Computed = true
	id.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}

	planCode := str(vpsDescPlanCode)
	planCode.Required = true
	planCode.PlanModifiers = []planmodifier.String{vpsRequiresReplaceIfPriorKnown()}

	region := str(vpsDescResourceRegion)
	region.Required = true
	region.PlanModifiers = []planmodifier.String{vpsRequiresReplaceIfPriorKnown()}
	location := object(vpsDescLocation, map[string]schema.Attribute{"region": region})
	location.Required = true

	displayName := str(vpsDescDisplayName)
	displayName.Optional = true
	displayName.Computed = true
	displayName.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
	displayName.Validators = []validator.String{stringvalidator.LengthAtMost(50)}

	powerState := str(vpsDescResourcePowerState)
	powerState.Optional = true
	powerState.Computed = true
	powerState.Default = stringdefault.StaticString(vpsDefaultPowerState)
	powerState.Validators = []validator.String{stringvalidator.OneOf("ACTIVE", "RESCUE", "SHUTOFF")}

	imageId := str(vpsDescImageId)
	imageId.Optional = true
	imageId.Computed = true
	imageId.PlanModifiers = []planmodifier.String{keepStateUnlessConfigured()}
	imageName := str(vpsDescImageName)
	imageName.Optional = true
	imageName.Computed = true
	imageName.PlanModifiers = []planmodifier.String{keepStateUnlessConfigured()}
	imageName.Validators = []validator.String{stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("image_id"))}
	from := object(vpsDescResourceInstallFrom, map[string]schema.Attribute{"image_id": imageId, "image_name": imageName})
	from.Required = true

	sshKey := str(vpsDescSshKey)
	sshKey.Optional = true
	script := str(vpsDescPostInstallationScript)
	script.Optional = true
	script.Sensitive = true
	installation := object(vpsDescInstallation, map[string]schema.Attribute{
		"from":                     from,
		"ssh_key":                  sshKey,
		"post_installation_script": script,
		"send_password_by_email": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
			Description:         vpsDescSendPasswordByEmail,
			MarkdownDescription: vpsDescSendPasswordByEmail,
		},
	})
	installation.Required = true

	cron := str(vpsDescAutoBackupCron)
	cron.Optional = true
	cron.Computed = true
	cron.Validators = []validator.String{stringvalidator.RegexMatches(vpsCronPattern,
		`must be a daily cron expression "M H * * *", with M from 0 to 59 and H from 0 to 23, without leading zeros, such as "0 2 * * *"`)}
	autoBackup := object(vpsDescAutoBackup, map[string]schema.Attribute{"cron": cron})
	autoBackup.Optional = true
	autoBackup.Computed = true
	options := object(vpsDescOptions, map[string]schema.Attribute{"auto_backup": autoBackup})
	options.Optional = true
	options.Computed = true
	options.PlanModifiers = []planmodifier.Object{objectplanmodifier.UseStateForUnknown()}

	duration := str(vpsDescOrderDuration)
	duration.Optional = true
	duration.Computed = true
	duration.Default = stringdefault.StaticString(vpsDefaultDuration)
	pricingMode := str(vpsDescOrderPricingMode)
	pricingMode.Optional = true
	pricingMode.Computed = true
	pricingMode.Default = stringdefault.StaticString(vpsDefaultPricingMode)
	order := object(vpsDescOrder, map[string]schema.Attribute{
		"duration":     duration,
		"pricing_mode": pricingMode,
		"order_id": schema.Int64Attribute{
			Computed:            true,
			Description:         vpsDescOrderId,
			MarkdownDescription: vpsDescOrderId,
			PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"options": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Description:         vpsDescOrderOptions,
			MarkdownDescription: vpsDescOrderOptions,
		},
	})
	order.Optional = true
	order.Computed = true
	order.PlanModifiers = []planmodifier.Object{objectplanmodifier.UseStateForUnknown()}

	timeouts := object(vpsDescTimeouts, map[string]schema.Attribute{
		"create": timeout(vpsDescTimeoutCreate),
		"update": timeout(vpsDescTimeoutUpdate),
		"delete": timeout(vpsDescTimeoutDelete),
	})
	timeouts.Optional = true

	currentState := object(vpsDescCurrentState, vpsResourceCurrentStateAttributes())
	currentState.Computed = true
	currentState.PlanModifiers = []planmodifier.Object{vpsCurrentStatePlanModifier()}

	return map[string]schema.Attribute{
		"id":            id,
		"plan_code":     planCode,
		"location":      location,
		"display_name":  displayName,
		"power_state":   powerState,
		"installation":  installation,
		"options":       options,
		"order":         order,
		"timeouts":      timeouts,
		"current_state": currentState,
	}
}

// vpsResourceCurrentStateAttributes is current_state of the data source
// (vpsCurrentStateAttributes), as resource attributes.
func vpsResourceCurrentStateAttributes() map[string]schema.Attribute {
	return vpsResourceComputedAttributes(vpsCurrentStateAttributes())
}

// vpsResourceComputedAttributes converts computed data source attributes into
// resource attributes, keeping their types and descriptions.
func vpsResourceComputedAttributes(attrs map[string]dsschema.Attribute) map[string]schema.Attribute {
	out := make(map[string]schema.Attribute, len(attrs))
	for name, a := range attrs {
		switch a := a.(type) {
		case dsschema.StringAttribute:
			out[name] = schema.StringAttribute{CustomType: a.CustomType, Computed: true, Description: a.Description, MarkdownDescription: a.MarkdownDescription}
		case dsschema.Int64Attribute:
			out[name] = schema.Int64Attribute{CustomType: a.CustomType, Computed: true, Description: a.Description, MarkdownDescription: a.MarkdownDescription}
		case dsschema.BoolAttribute:
			out[name] = schema.BoolAttribute{CustomType: a.CustomType, Computed: true, Description: a.Description, MarkdownDescription: a.MarkdownDescription}
		case dsschema.SingleNestedAttribute:
			out[name] = schema.SingleNestedAttribute{Computed: true, Description: a.Description, MarkdownDescription: a.MarkdownDescription,
				Attributes: vpsResourceComputedAttributes(a.Attributes)}
		case dsschema.ListNestedAttribute:
			out[name] = schema.ListNestedAttribute{Computed: true, Description: a.Description, MarkdownDescription: a.MarkdownDescription,
				NestedObject: schema.NestedAttributeObject{Attributes: vpsResourceComputedAttributes(a.NestedObject.Attributes)}}
		default:
			panic(fmt.Sprintf("attribute %s: %T is not converted to a resource attribute", name, a))
		}
	}
	return out
}

// Attribute types, matching vpsResourceAttributes.

func vpsResourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"plan_code":     ovhtypes.TfStringType{},
		"display_name":  ovhtypes.TfStringType{},
		"power_state":   ovhtypes.TfStringType{},
		"location":      types.ObjectType{AttrTypes: vpsLocationAttrTypes()},
		"installation":  types.ObjectType{AttrTypes: vpsResourceInstallationAttrTypes()},
		"options":       types.ObjectType{AttrTypes: vpsOptionsAttrTypes()},
		"order":         types.ObjectType{AttrTypes: vpsOrderAttrTypes()},
		"timeouts":      types.ObjectType{AttrTypes: vpsTimeoutsAttrTypes()},
		"current_state": types.ObjectType{AttrTypes: vpsCurrentStateAttrTypes()},
	}
}

func vpsResourceInstallationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"from":                     types.ObjectType{AttrTypes: vpsInstallationFromAttrTypes()},
		"ssh_key":                  ovhtypes.TfStringType{},
		"post_installation_script": ovhtypes.TfStringType{},
		"send_password_by_email":   types.BoolType,
	}
}

func vpsOrderAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"duration":     ovhtypes.TfStringType{},
		"pricing_mode": ovhtypes.TfStringType{},
		"order_id":     types.Int64Type,
		"options":      types.ListType{ElemType: types.StringType},
	}
}

func vpsTimeoutsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create": ovhtypes.TfStringType{},
		"update": ovhtypes.TfStringType{},
		"delete": ovhtypes.TfStringType{},
	}
}

// timeout returns the duration that timeouts sets for the operation (create,
// update or delete), or def.
func (m *VPSResourceModel) timeout(operation string, def time.Duration) time.Duration {
	value, _ := m.Timeouts.Attributes()[operation].(ovhtypes.TfStringValue)
	if value.IsNull() || value.IsUnknown() {
		return def
	}
	duration, err := time.ParseDuration(value.ValueString()) // checked by the schema
	if err != nil {
		return def
	}
	return duration
}
