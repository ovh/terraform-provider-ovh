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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithImportState = (*cloudInstanceResource)(nil)
)

func NewCloudInstanceResource() resource.Resource {
	return &cloudInstanceResource{}
}

type cloudInstanceResource struct {
	config *Config
}

var instanceMutableAttrs = MutableAttrs{
	Strings:           []string{"name", "flavor_id", "image_id", "power_state"},
	Lists:             []string{"networks", "shares"},
	CustomStringLists: []string{"volume_ids", "security_group_ids"},
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

func (r *cloudInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<instance_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an instance in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Immutable
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
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
					stringplanmodifier.UseStateForUnknown(),
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
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"power_state": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Desired power state: ACTIVE, SHUTOFF or SHELVED (defaults to ACTIVE)",
				MarkdownDescription: "Desired power state: `ACTIVE`, `SHUTOFF` or `SHELVED` (defaults to `ACTIVE`)",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "SHUTOFF", "SHELVED"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"networks": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Network interfaces attached to the instance, in order (the public interface is placed first by the platform)",
				MarkdownDescription: "Network interfaces attached to the instance, in order (the public interface is placed first by the platform)",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"public": schema.BoolAttribute{
							Optional:    true,
							Description: "Attach the public network on this interface",
						},
						"network_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Optional:    true,
							Description: "Private network ID",
						},
						"subnet_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Optional:    true,
							Description: "Subnet ID within the private network",
						},
						"floating_ip_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Optional:    true,
							Description: "Floating IP ID associated with this interface",
						},
					},
				},
			},
			"volume_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "IDs of block-storage volumes attached to the instance",
				MarkdownDescription: "IDs of block-storage volumes attached to the instance",
			},
			"security_group_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "IDs of security groups applied to all interfaces",
				MarkdownDescription: "IDs of security groups applied to all interfaces",
			},
			"shares": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Filesystem shares attached to the instance. Do NOT also manage the same share's access rules on another resource (the share will show OUT_OF_SYNC)",
				MarkdownDescription: "Filesystem shares attached to the instance. Do NOT also manage the same share's access rules on another resource (the share will show `OUT_OF_SYNC`)",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Share ID",
						},
						"access_level": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Optional:    true,
							Computed:    true,
							Description: "Access level: READ_ONLY or READ_WRITE (defaults to READ_WRITE)",
							Validators: []validator.String{
								stringvalidator.OneOf("READ_ONLY", "READ_WRITE"),
							},
						},
					},
				},
			},
			// Computed envelope
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Instance ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:    ovhtypes.TfStringType{},
				Computed:      true,
				Description:   "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(instanceMutableAttrs)},
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date of the instance",
			},
			"updated_at": schema.StringAttribute{
				CustomType:    ovhtypes.TfStringType{},
				Computed:      true,
				Description:   "Last update date of the instance",
				PlanModifiers: []planmodifier.String{UnknownDuringUpdateStringModifier(instanceMutableAttrs)},
			},
			"resource_status": schema.StringAttribute{
				CustomType:    ovhtypes.TfStringType{},
				Computed:      true,
				Description:   "Instance readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{OutOfSyncPlanModifier()},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:      true,
				Description:   "Observed state of the instance",
				PlanModifiers: []planmodifier.Object{UnknownDuringUpdateObjectModifier(instanceMutableAttrs)},
				Attributes:    instanceCurrentStateSchemaAttributes(),
			},
		},
	}
}

// instanceCurrentStateSchemaAttributes returns the schema attributes for the
// computed current_state object. Shared by the resource and both data sources.
func instanceCurrentStateSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed instance name"},
		"power_state":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Observed power state"},
		"locked":       schema.BoolAttribute{Computed: true, Description: "Whether the instance is locked"},
		"ssh_key_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "SSH key injected at boot"},
		"host_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Opaque physical host ID"},
		"project_id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Owning project ID"},
		"user_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Owning OpenStack user ID"},
		"flavor": schema.SingleNestedAttribute{Computed: true, Description: "Observed flavor details", Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"vcpus":     schema.Int64Attribute{Computed: true},
			"ram":       schema.Int64Attribute{Computed: true},
			"disk":      schema.Int64Attribute{Computed: true},
			"swap":      schema.Int64Attribute{Computed: true},
			"ephemeral": schema.Int64Attribute{Computed: true},
		}},
		"image": schema.SingleNestedAttribute{Computed: true, Description: "Observed image details (null for boot-from-volume)", Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"size":       schema.Int64Attribute{Computed: true},
			"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"deprecated": schema.BoolAttribute{Computed: true},
		}},
		"location": schema.SingleNestedAttribute{Computed: true, Description: "Observed location", Attributes: map[string]schema.Attribute{
			"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}},
		"networks": schema.ListNestedAttribute{Computed: true, Description: "Observed network interfaces", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"public":         schema.BoolAttribute{Computed: true},
			"subnet_id":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"gateway_id":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"floating_ip_id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"addresses": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"ip":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"mac":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"type":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
				"version": schema.Int64Attribute{Computed: true},
			}}},
		}}},
		"volumes": schema.ListNestedAttribute{Computed: true, Description: "Observed attached volumes", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"size": schema.Int64Attribute{Computed: true},
		}}},
		"shares": schema.ListNestedAttribute{Computed: true, Description: "Observed attached shares (only populated on single-instance read)", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"access_level": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"access_to":    schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
			"state":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"security_groups": schema.ListNestedAttribute{Computed: true, Description: "Observed security groups", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}}},
		"group": schema.SingleNestedAttribute{Computed: true, Description: "Placement group membership", Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
		}},
	}
}

func (r *cloudInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance"

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Save state immediately so the ID is tracked even if the workflow fails.
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if _, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), responseData.Id); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance to be ready", err.Error())
		return
	}

	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(responseData.Id)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	// Refresh the checksum right before PUT to avoid a 409 ChecksumMismatch if
	// server-side drift bumped it since the last read.
	var currentData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &currentData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	updatePayload := planData.ToUpdate(currentData.Checksum)

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	if _, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error waiting for instance to be ready after update", err.Error())
		return
	}

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	planData.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudInstanceModel
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
				if reason := res.taskErrorSummary(); reason != "" {
					return res, res.ResourceStatus, fmt.Errorf("instance %s entered ERROR state: %s", instanceId, reason)
				}
				return res, res.ResourceStatus, fmt.Errorf("instance %s entered ERROR state", instanceId)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}
