package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                   = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithConfigure      = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithImportState    = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithValidateConfig = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithModifyPlan     = (*cloudInstanceResource)(nil)
)

func NewCloudInstanceResource() resource.Resource {
	return &cloudInstanceResource{}
}

type cloudInstanceResource struct {
	config *Config
}

func (r *cloudInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

func (r *cloudInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.config = config
}

var instanceMutableAttrs = MutableAttrs{
	Strings:           []string{"name", "flavor_id", "image_id", "power_state", "user_data"},
	Lists:             []string{"networks", "shares", "volumes"},
	CustomStringLists: []string{"volume_ids", "security_group_ids"},
}

// ValidateConfig enforces the volume rules the API would otherwise reject at
// apply: volumes and the deprecated volume_ids cannot be mixed, at most one volume
// is BOOT, and an image_id excludes a boot volume.
func (r *cloudInstanceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var volumes types.List
	var volumeIds ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
	var imageId ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("volumes"), &volumes)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("volume_ids"), &volumeIds)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("image_id"), &imageId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A null is "not set"; an unknown may still turn out to be set, so only a
	// known value counts as a conflict.
	if !volumes.IsNull() && !volumes.IsUnknown() && !volumeIds.IsNull() && !volumeIds.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("volumes"), "Conflicting volume attributes",
			"volumes cannot be combined with volume_ids. Use volumes alone: it gives each volume a role.")
	}
	if !volumes.IsUnknown() {
		boots := bootVolumeIds(volumes)
		if len(boots) > 1 {
			resp.Diagnostics.AddAttributeError(path.Root("volumes"), "Several boot volumes",
				"At most one entry of volumes can have the role BOOT.")
		}
		if len(boots) > 0 && knownString(imageId) != "" {
			resp.Diagnostics.AddAttributeError(path.Root("image_id"), "Two boot sources",
				"image_id and a volume with the role BOOT are mutually exclusive: an instance boots from one or the other.")
		}
	}
}

// ModifyPlan replaces the instance when volumes gives the BOOT role to another
// volume: the API only lets the boot volume be chosen at creation.
//
// It never plans a replacement it cannot vouch for, since the instance would be
// destroyed before a creation the API may refuse:
//
//   - The deprecated volume_ids never replaces the instance. It carries no role,
//     and on update the API keeps the stored boot volume, as before volumes
//     existed.
//   - The current boot volume must be known: the BOOT volume of the state, or none
//     when the state has an image_id. A state with neither does not tell (written
//     before volumes existed, or an instance whose boot volume the API could not
//     resolve), and the update is left to the API.
//   - A plan with an image and a BOOT volume, or a known boot volume losing its
//     BOOT role, is an error: the API refuses both.
func (r *cloudInstanceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	if req.State.Raw.IsNull() {
		r.modifyCreatePlan(ctx, req, resp)
		return
	}
	var cfgVolumes, planVolumes, stateVolumes types.List
	var planImage, stateImage ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("volumes"), &cfgVolumes)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("volumes"), &planVolumes)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("image_id"), &planImage)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("volumes"), &stateVolumes)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("image_id"), &stateImage)...)
	if resp.Diagnostics.HasError() || cfgVolumes.IsNull() || cfgVolumes.IsUnknown() || planVolumes.IsUnknown() || planImage.IsUnknown() {
		return
	}

	// The roles that count are the planned ones: a role left out is resolved by the
	// plan (BOOT on the current boot volume, DATA otherwise). An unknown role tells
	// nothing, and the update is left to the API. An unknown id with the role BOOT
	// is a volume created or replaced in the same apply: a new boot volume for
	// sure, since the current one has a known id.
	planned, plannedUnknown, idUnknown := "", false, false
	for _, ref := range plannedVolumeRefs(planVolumes) {
		if ref.roleUnknown {
			return
		}
		if ref.idUnknown {
			idUnknown = true
		}
		if ref.role != instanceVolumeRoleBoot || planned != "" || plannedUnknown {
			continue
		}
		if ref.idUnknown {
			plannedUnknown = true
		} else {
			planned = ref.id
		}
	}
	if (planned != "" || plannedUnknown) && knownString(planImage) != "" {
		if plannedUnknown {
			planned = "to be created"
		}
		resp.Diagnostics.AddAttributeError(path.Root("volumes"), "Two boot sources",
			"The instance boots from the image "+knownString(planImage)+" and volumes gives the volume "+planned+" the role BOOT: an instance boots from one or the other, and the API refuses both. "+
				"Set role to DATA on the volume to keep booting from the image. The boot source cannot be switched between an image and a volume in place: "+
				"to boot this instance from the volume, destroy it and create a new instance with no image_id.")
		return
	}

	current, known := "", false
	if ids := bootVolumeIds(stateVolumes); len(ids) > 0 {
		current, known = ids[0], true
	} else if knownString(stateImage) != "" {
		known = true
	}
	if !known {
		return
	}
	if plannedUnknown {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("volumes"))
		return
	}
	if planned == current {
		return
	}
	if planned == "" {
		if idUnknown {
			// The current boot volume may hide behind an unknown id.
			return
		}
		// The API keeps the stored boot volume: it refuses to detach it, and a DATA
		// role on it is silently turned back into BOOT, which would leave a
		// permanent diff. Both are refused at plan, each with its own wording.
		for _, ref := range volumesToAPI(planVolumes) {
			if ref.Id == current {
				resp.Diagnostics.AddAttributeError(path.Root("volumes"), "Boot volume cannot become DATA",
					"The instance boots from the volume "+current+", which keeps the role BOOT: the API ignores a DATA role on it, and the plan would never settle. "+
						"Leave role out on it, or set it to BOOT.")
				return
			}
		}
		resp.Diagnostics.AddAttributeError(path.Root("volumes"), "Boot volume cannot be detached",
			"The instance boots from the volume "+current+", which is missing from volumes: the API refuses to detach the boot volume. "+
				"Keep an entry for it, with role BOOT or no role. To stop booting from it, destroy the instance and create a new one.")
		return
	}
	resp.RequiresReplace = append(resp.RequiresReplace, path.Root("volumes"))
}

// modifyCreatePlan checks the boot source of a creation. With no state, a role
// left out in volumes is DATA, and the only boot sources are image_id and a BOOT
// entry. A creation with neither is refused by the API: ValidateConfig cannot see
// it, since image_id is Optional+Computed and may be null in a configuration that
// relies on the state. It matters most for the create plan of a replacement,
// which Terraform makes with no prior state once the instance is destroyed. The
// deprecated volume_ids keeps its rule: its first volume boots.
//
// The values come from the configuration: with no state, the plan holds unknown
// for every Optional+Computed attribute the configuration leaves out, and an
// unknown image_id would hide the very case this check exists for.
func (r *cloudInstanceResource) modifyCreatePlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var volumes types.List
	var volumeIds ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]
	var imageId ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("volumes"), &volumes)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("volume_ids"), &volumeIds)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("image_id"), &imageId)...)
	if resp.Diagnostics.HasError() || imageId.IsUnknown() || volumes.IsUnknown() || volumeIds.IsUnknown() {
		return
	}
	if !imageId.IsNull() || !volumeIds.IsNull() || volumes.IsNull() {
		return
	}
	for _, ref := range plannedVolumeRefs(volumes) {
		if ref.idUnknown || ref.roleUnknown || ref.role == instanceVolumeRoleBoot {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(path.Root("volumes"), "No boot source",
		"The instance would be created with neither an image_id nor a volume with the role BOOT, which the API refuses. "+
			"Set image_id, or give the volume to boot from the role BOOT: a role left out is DATA.")
}

func (r *cloudInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an instance in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Immutable
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         instanceDescServiceName,
				MarkdownDescription: instanceDescServiceName,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the instance is created",
				MarkdownDescription: "Region where the instance is created",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"availability_zone": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Availability zone of the instance (immutable; assigned by the platform if omitted)",
				MarkdownDescription: "Availability zone of the instance (immutable; assigned by the platform if omitted)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					keepStateUnlessConfigured(),
				},
			},
			"ssh_key_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Name of the SSH key injected at boot (immutable)",
				MarkdownDescription: "Name of the SSH key injected at boot (immutable)",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"group_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "ID of the placement group the instance belongs to (immutable)",
				MarkdownDescription: "ID of the placement group the instance belongs to (immutable)",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			// Mutable
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Instance name",
				MarkdownDescription: "Instance name",
			},
			"flavor_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Flavor ID. Changing it resizes the instance in place",
				MarkdownDescription: "Flavor ID. Changing it resizes the instance in place",
			},
			"image_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Image ID to boot from. Omit for a boot-from-volume instance. WARNING: changing it rebuilds the instance and WIPES the root disk",
				MarkdownDescription: "Image ID to boot from. Omit for a boot-from-volume instance. **WARNING**: changing it rebuilds the instance and **wipes the root disk**",
				PlanModifiers:       []planmodifier.String{keepStateUnlessConfigured()},
			},
			"power_state": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Desired power state: ACTIVE, SHUTOFF, SHELVED, PAUSED or SUSPENDED (defaults to ACTIVE)",
				MarkdownDescription: "Desired power state: `ACTIVE`, `SHUTOFF`, `SHELVED`, `PAUSED` or `SUSPENDED` (defaults to `ACTIVE`)",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "SHUTOFF", "SHELVED", "PAUSED", "SUSPENDED"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			// Mutable, and NOT RequiresReplace: the API applies a new value by
			// rebuilding the instance in place, so its id and ports survive.
			"user_data": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Sensitive:           true,
				Description:         instanceDescUserData,
				MarkdownDescription: instanceDescUserDataMd,
			},
			"networks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Network interfaces attached to the instance. Entries keep the order they are written in; the API returns them sorted by network id and the provider re-orders them back to the configuration. Four shapes: auto_assign_public_ip alone for a platform-assigned public IP (at most one entry); ip alone for a public IP the project already owns (additional IP, or an Ext-Net IP of the project in the instance's region); network_id + subnet_id for a private interface with an IPAM address; network_id + subnet_id + ip to pin the fixed address when ip is inside the subnet CIDR, or to associate an existing floating IP otherwise",
				MarkdownDescription: "Network interfaces attached to the instance. Entries keep the order they are written in; the API returns them sorted by network id and the provider re-orders them back to the configuration. Four shapes:\n" +
					"  * `auto_assign_public_ip` alone — public interface with a platform-assigned public IP (at most one such entry)\n" +
					"  * `ip` alone — a public IP the project already owns: an additional IP, or an Ext-Net IP of the project in the instance's region. Several are allowed and may coexist with `auto_assign_public_ip`\n" +
					"  * `network_id` + `subnet_id` — private interface with an address picked by IPAM\n" +
					"  * `network_id` + `subnet_id` + `ip` — pins the port's fixed address when `ip` is inside the subnet CIDR, otherwise associates the existing floating IP `ip`",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         instanceDescNetworkRefNetworkId,
							MarkdownDescription: instanceDescNetworkRefNetworkId,
						},
						"subnet_id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         instanceDescNetworkRefSubnetId,
							MarkdownDescription: instanceDescNetworkRefSubnetIdMd,
						},
						"ip": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         instanceDescNetworkRefIp,
							MarkdownDescription: instanceDescNetworkRefIpMd,
						},
						"auto_assign_public_ip": schema.BoolAttribute{
							Optional:            true,
							Description:         instanceDescNetworkRefAutoAssign,
							MarkdownDescription: instanceDescNetworkRefAutoAssignMd,
						},
					},
				},
			},
			"volume_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Computed:            true,
				Description:         instanceDescVolumeIds,
				MarkdownDescription: instanceDescVolumeIds,
				PlanModifiers:       []planmodifier.List{derivedVolumeList("volumes")},
				DeprecationMessage:  "Use volumes, which gives each volume a role. volume_ids alone keeps working: without image_id, its first volume is the boot volume.",
			},
			"volumes": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         instanceDescVolumes,
				MarkdownDescription: instanceDescVolumesMd,
				PlanModifiers:       []planmodifier.List{derivedVolumeList("volume_ids")},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         instanceDescVolumeRefId,
							MarkdownDescription: instanceDescVolumeRefId,
						},
						"role": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{roleDefaultsToData()},
							Description:         instanceDescVolumeRefRole,
							MarkdownDescription: instanceDescVolumeRefRole,
							Validators: []validator.String{
								stringvalidator.OneOf(instanceVolumeRoleBoot, instanceVolumeRoleData),
							},
						},
					},
				},
			},
			"security_group_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Computed:            true,
				Description:         "IDs of security groups applied to all interfaces. Omit it to let the platform apply the project's default security group; set an explicit empty list to apply no security group at all (the instance then accepts no inbound traffic)",
				MarkdownDescription: "IDs of security groups applied to all interfaces. Omit it to let the platform apply the project's `default` security group; set an explicit empty list (`[]`) to apply no security group at all (the instance then accepts no inbound traffic)",
				// Optional+Computed: the API resolves an omitted value to the project
				// default and echoes it back, so the server value must be reused when
				// the config stays silent, otherwise every plan touching another
				// attribute would show it as "known after apply".
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"shares": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Filesystem shares attached to the instance. Do NOT also manage the same share's access rules on another resource (the share will show OUT_OF_SYNC)",
				MarkdownDescription: "Filesystem shares attached to the instance. Do NOT also manage the same share's access rules on another resource (the share will show `OUT_OF_SYNC`)",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         instanceDescShareRefId,
							MarkdownDescription: instanceDescShareRefId,
						},
						// Optional-only, like networks[].auto_assign_public_ip: a nested
						// list element cannot carry UseStateForUnknown, so Optional+Computed
						// would re-mark it unknown on every list-length change. MergeWith
						// keeps it null when the config omitted it.
						"access_level": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         instanceDescShareRefAccessLevel,
							MarkdownDescription: instanceDescShareRefAccessLevelMd,
							Validators: []validator.String{
								stringvalidator.OneOf("READ_ONLY", "READ_WRITE"),
							},
						},
					},
				},
			},
			// Computed envelope
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceDescId,
				MarkdownDescription: instanceDescId,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceDescChecksum,
				MarkdownDescription: instanceDescChecksum,
				PlanModifiers:       []planmodifier.String{instanceKeepStateUnlessUpdate()},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceDescCreatedAt,
				MarkdownDescription: instanceDescCreatedAt,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceDescUpdatedAt,
				MarkdownDescription: instanceDescUpdatedAt,
				PlanModifiers:       []planmodifier.String{instanceKeepStateUnlessUpdate()},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         instanceDescResourceStatus,
				MarkdownDescription: instanceDescResourceStatusMd,
				PlanModifiers:       []planmodifier.String{instanceKeepStateUnlessUpdate(), OutOfSyncPlanModifier()},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         instanceDescCurrentState,
				MarkdownDescription: instanceDescCurrentState,
				PlanModifiers:       []planmodifier.Object{instanceKeepStateObject()},
				Attributes:          instanceCurrentStateSchemaAttributes(),
			},
		},
	}
}

// instanceCurrentStateSchemaAttributes returns the schema attributes for the
// computed current_state object. Shared by the resource and both data sources.
func instanceCurrentStateSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsName, MarkdownDescription: instanceDescCsName},
		"power_state":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsPowerState, MarkdownDescription: instanceDescCsPowerState},
		"locked":       schema.BoolAttribute{Computed: true, Description: instanceDescCsLocked, MarkdownDescription: instanceDescCsLocked},
		"ssh_key_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsSSHKeyName, MarkdownDescription: instanceDescCsSSHKeyName},
		"host_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsHostId, MarkdownDescription: instanceDescCsHostId},
		"project_id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsProjectId, MarkdownDescription: instanceDescCsProjectId},
		"user_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsUserId, MarkdownDescription: instanceDescCsUserId},
		"flavor": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsFlavor, MarkdownDescription: instanceDescCsFlavor, Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescFlavorId, MarkdownDescription: instanceDescFlavorId},
			"name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescFlavorName, MarkdownDescription: instanceDescFlavorName},
			"vcpus":     schema.Int64Attribute{Computed: true, Description: instanceDescFlavorVcpus, MarkdownDescription: instanceDescFlavorVcpus},
			"ram":       schema.Int64Attribute{Computed: true, Description: instanceDescFlavorRam, MarkdownDescription: instanceDescFlavorRam},
			"disk":      schema.Int64Attribute{Computed: true, Description: instanceDescFlavorDisk, MarkdownDescription: instanceDescFlavorDisk},
			"swap":      schema.Int64Attribute{Computed: true, Description: instanceDescFlavorSwap, MarkdownDescription: instanceDescFlavorSwap},
			"ephemeral": schema.Int64Attribute{Computed: true, Description: instanceDescFlavorEphemeral, MarkdownDescription: instanceDescFlavorEphemeral},
		}},
		"image": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsImage, MarkdownDescription: instanceDescCsImage, Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageId, MarkdownDescription: instanceDescImageId},
			"name":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageName, MarkdownDescription: instanceDescImageName},
			"size":       schema.Int64Attribute{Computed: true, Description: instanceDescImageSize, MarkdownDescription: instanceDescImageSize},
			"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescImageStatus, MarkdownDescription: instanceDescImageStatus},
			"deprecated": schema.BoolAttribute{Computed: true, Description: instanceDescImageDeprecated, MarkdownDescription: instanceDescImageDeprecated},
		}},
		"location": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsLocation, MarkdownDescription: instanceDescCsLocation, Attributes: map[string]schema.Attribute{
			"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationRegion, MarkdownDescription: instanceDescLocationRegion},
			"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescLocationAZ, MarkdownDescription: instanceDescLocationAZ},
		}},
		"networks": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsNetworks, MarkdownDescription: instanceDescCsNetworks, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkId, MarkdownDescription: instanceDescCsNetworkId},
			"subnet_id":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkSubnetId, MarkdownDescription: instanceDescCsNetworkSubnetId},
			"gateway_id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsNetworkGatewayId, MarkdownDescription: instanceDescCsNetworkGatewayId},
			"addresses": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsAddresses, MarkdownDescription: instanceDescCsAddressesMd, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"ip":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressIp, MarkdownDescription: instanceDescAddressIp},
				"mac":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressMac, MarkdownDescription: instanceDescAddressMac},
				"type":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescAddressType, MarkdownDescription: instanceDescAddressTypeMd},
				"version": schema.Int64Attribute{Computed: true, Description: instanceDescAddressVersion, MarkdownDescription: instanceDescAddressVersion},
			}}},
		}}},
		"volumes": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsVolumes, MarkdownDescription: instanceDescCsVolumes, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescVolumeId, MarkdownDescription: instanceDescVolumeId},
			"name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescVolumeName, MarkdownDescription: instanceDescVolumeName},
			"size": schema.Int64Attribute{Computed: true, Description: instanceDescVolumeSize, MarkdownDescription: instanceDescVolumeSize},
			"role": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescVolumeRole, MarkdownDescription: instanceDescVolumeRole},
		}}},
		"shares": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsShares, MarkdownDescription: instanceDescCsShares, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareId, MarkdownDescription: instanceDescCsShareId},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareAccessLevel, MarkdownDescription: instanceDescCsShareAccessLevelMd},
			"access_to":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareAccessTo, MarkdownDescription: instanceDescCsShareAccessTo},
			"state":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescCsShareState, MarkdownDescription: instanceDescCsShareStateMd},
		}}},
		"security_groups": schema.ListNestedAttribute{Computed: true, Description: instanceDescCsSecurityGroups, MarkdownDescription: instanceDescCsSecurityGroups, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescSecurityGroupId, MarkdownDescription: instanceDescSecurityGroupId},
		}}},
		"group": schema.SingleNestedAttribute{Computed: true, Description: instanceDescCsGroup, MarkdownDescription: instanceDescCsGroup, Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: instanceDescGroupId, MarkdownDescription: instanceDescGroupId},
		}},
	}
}

func (r *cloudInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<instance_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Snapshot the planned list shapes: the merge below overwrites them, and the
	// second merge must still order networks like the configuration.
	prior := data.priorSpec()
	createPayload := data.ToCreate()
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance"

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Save state immediately so the ID is tracked even if the workflow fails.
	data.MergeWith(ctx, &responseData, prior)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if _, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), responseData.Id); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance to be ready", err.Error())
		return
	}

	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(responseData.Id)
	// Reset: json.Unmarshal keeps old values for keys the server omits (omitempty).
	responseData = CloudInstanceAPIResponse{}
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, prior)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prior := data.priorSpec()
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData, prior)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prior := planData.priorSpec()
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	// Refresh the checksum right before PUT to avoid a 409 ChecksumMismatch if
	// server-side drift bumped it since the last read.
	var currentData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &currentData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	// data holds the prior state: ToUpdate needs it to tell an unchanged user_data
	// (absent key) from one the config actually moved (rebuild).
	updatePayload := planData.ToUpdate(currentData.Checksum, data.UserData)

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	if _, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance to be ready after update", err.Error())
		return
	}

	// Reset: json.Unmarshal keeps old values for keys the server omits (omitempty).
	responseData = CloudInstanceAPIResponse{}
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	planData.MergeWith(ctx, &responseData, prior)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudInstanceAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance to be deleted", err.Error())
	}
}

func (r *cloudInstanceResource) waitForInstanceReady(ctx context.Context, serviceName, instanceId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudInstanceAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/compute/instance/" + url.PathEscape(instanceId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			// ERROR is terminal: stop polling and surface the reason reported by
			// the failed task(s) instead of letting the SDK emit a generic
			// "unexpected state 'ERROR'. last error: %!s(<nil>)".
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("instance", instanceId, res.CurrentTasks)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}
