package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// ------------------------------------------------------------------
// Resource model — top-level mutable attributes mirror targetSpec.
// We share CloudQuotaModel from types_cloud_quota.go so that the
// data source and resource expose the same flattened shape.
// ------------------------------------------------------------------

// ------------------------------------------------------------------
// PUT payload types
// ------------------------------------------------------------------

// CloudQuotaUpdateTargetSpecAPI is the wire payload for the update target spec.
type CloudQuotaUpdateTargetSpecAPI struct {
	PreventAutomaticQuotaUpgrade bool                            `json:"preventAutomaticQuotaUpgrade"`
	Regions                      []CloudQuotaRegionTargetSpecAPI `json:"regions"`
}

// CloudQuotaUpdateAPI is the wire payload for PUT /quota.
type CloudQuotaUpdateAPI struct {
	Checksum   string                        `json:"checksum"`
	TargetSpec CloudQuotaUpdateTargetSpecAPI `json:"targetSpec"`
}

// ------------------------------------------------------------------
// Helpers — extract top-level fields from the Terraform model into
// the API payload.
// ------------------------------------------------------------------

type quotaRegionPlan struct {
	Region  ovhtypes.TfStringValue `tfsdk:"region"`
	Profile ovhtypes.TfStringValue `tfsdk:"profile"`
}

func (m *CloudQuotaResourceModel) toTargetSpecAPI(ctx context.Context) (CloudQuotaUpdateTargetSpecAPI, error) {
	regions := make([]CloudQuotaRegionTargetSpecAPI, 0)

	if !m.Regions.IsNull() && !m.Regions.IsUnknown() {
		var plan []quotaRegionPlan
		diags := m.Regions.ElementsAs(ctx, &plan, false)
		if diags.HasError() {
			return CloudQuotaUpdateTargetSpecAPI{}, fmt.Errorf("invalid regions: %s", diags)
		}
		for _, r := range plan {
			regions = append(regions, CloudQuotaRegionTargetSpecAPI{
				Location: &CloudQuotaLocationAPI{
					Region: r.Region.ValueString(),
				},
				Profile: r.Profile.ValueString(),
			})
		}
	}

	return CloudQuotaUpdateTargetSpecAPI{
		PreventAutomaticQuotaUpgrade: m.PreventAutomaticQuotaUpgrade.ValueBool(),
		Regions:                      regions,
	}, nil
}

// ------------------------------------------------------------------
// Resource boilerplate
// ------------------------------------------------------------------

var (
	_ resource.Resource                = (*cloudQuotaResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudQuotaResource)(nil)
	_ resource.ResourceWithImportState = (*cloudQuotaResource)(nil)
)

// NewCloudQuotaResource returns a new writable cloud quota resource.
func NewCloudQuotaResource() resource.Resource {
	return &cloudQuotaResource{}
}

type cloudQuotaResource struct {
	config *Config
}

func (r *cloudQuotaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_quota"
}

func (r *cloudQuotaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var quotaMutableAttrs = MutableAttrs{
	Bools: []string{"prevent_automatic_quota_upgrade"},
	Lists: []string{"regions"},
}

func (r *cloudQuotaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the public cloud project quota: applied profile per region and automatic-quota-upgrade toggle. There is exactly one quota envelope per project; importing a project's quota uses the service name as id.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the cloud project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Flattened targetSpec
			"prevent_automatic_quota_upgrade": schema.BoolAttribute{
				Required:    true,
				Description: "When true, automatic quota upgrades are disabled for this project.",
			},
			"regions": schema.ListNestedAttribute{
				Required:    true,
				Description: "Target quota profile per region. Regions omitted are left unchanged upstream.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Region where the profile applies (e.g. GRA11).",
						},
						"profile": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Required:    true,
							Description: "Quota profile to apply in this region. Available values are in current_state.available_profiles.",
						},
					},
				},
			},

			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource identifier (the project id).",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Quota readiness in the system (CREATING, UPDATING, DELETING, OUT_OF_SYNC, READY, ERROR, SUSPENDED, UNKNOWN).",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Computed hash representing the current target specification value.",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(quotaMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation date (RFC3339).",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date (RFC3339).",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(quotaMutableAttrs),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current quota state of the project (live usage and available profiles).",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(quotaMutableAttrs),
				},
				Attributes: cloudQuotaCurrentStateSchemaAttrs(),
			},
		},
	}
}

// cloudQuotaCurrentStateSchemaAttrs mirrors the data source's current_state shape.
// Kept verbatim so resource and data source produce identical plan output.
func cloudQuotaCurrentStateSchemaAttrs() map[string]schema.Attribute {
	usageAttrs := map[string]schema.Attribute{
		"limit": schema.Int64Attribute{Computed: true},
		"used":  schema.Int64Attribute{Computed: true},
		"unit":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
	}
	limitAttrs := map[string]schema.Attribute{
		"limit": schema.Int64Attribute{Computed: true},
		"unit":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
	}

	return map[string]schema.Attribute{
		"prevent_automatic_quota_upgrade": schema.BoolAttribute{Computed: true},
		"available_profiles": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"compute": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"cores":     schema.Int64Attribute{Computed: true},
							"instances": schema.Int64Attribute{Computed: true},
							"memory":    schema.Int64Attribute{Computed: true},
						},
					},
					"volume": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"backup_size_total": schema.Int64Attribute{Computed: true},
							"backups":           schema.Int64Attribute{Computed: true},
							"size_total":        schema.Int64Attribute{Computed: true},
							"snapshots":         schema.Int64Attribute{Computed: true},
							"volumes":           schema.Int64Attribute{Computed: true},
						},
					},
					"network": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"floating_ips":         schema.Int64Attribute{Computed: true},
							"gateways":             schema.Int64Attribute{Computed: true},
							"networks":             schema.Int64Attribute{Computed: true},
							"security_group_rules": schema.Int64Attribute{Computed: true},
							"security_groups":      schema.Int64Attribute{Computed: true},
							"subnets":              schema.Int64Attribute{Computed: true},
						},
					},
					"loadbalancer": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"health_monitors": schema.Int64Attribute{Computed: true},
							"l7_policies":     schema.Int64Attribute{Computed: true},
							"l7_rules":        schema.Int64Attribute{Computed: true},
							"listeners":       schema.Int64Attribute{Computed: true},
							"loadbalancers":   schema.Int64Attribute{Computed: true},
							"members":         schema.Int64Attribute{Computed: true},
							"pools":           schema.Int64Attribute{Computed: true},
						},
					},
					"key_manager": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"containers": schema.Int64Attribute{Computed: true},
							"secrets":    schema.Int64Attribute{Computed: true},
						},
					},
					"share": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"backup_size_total": schema.Int64Attribute{Computed: true},
							"backups":           schema.Int64Attribute{Computed: true},
							"shares":            schema.Int64Attribute{Computed: true},
							"size_total":        schema.Int64Attribute{Computed: true},
							"snapshots":         schema.Int64Attribute{Computed: true},
						},
					},
					"keypair": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"keypairs": schema.Int64Attribute{Computed: true},
						},
					},
				},
			},
		},
		"regions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"region":  schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"profile": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"compute": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"cores":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"instances": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"memory":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"volume": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"backup_size_total": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"backups":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"per_volume_size":   schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
							"size_total":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"snapshots":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"volumes":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"network": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"floating_ips":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"gateways":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"networks":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"security_group_rules": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"security_groups":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"subnets":              schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"loadbalancer": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"health_monitors": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"l7_policies":     schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"l7_rules":        schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"listeners":       schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"loadbalancers":   schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"members":         schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"pools":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"key_manager": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"containers": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"secrets":    schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"share": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"backup_size_total":   schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"backups":             schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"per_share_size":      schema.SingleNestedAttribute{Computed: true, Attributes: limitAttrs},
							"share_networks":      schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"shares":              schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"size_total":          schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"snapshot_size_total": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
							"snapshots":           schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
					"keypair": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"keypairs": schema.SingleNestedAttribute{Computed: true, Attributes: usageAttrs},
						},
					},
				},
			},
		},
	}
}

func (r *cloudQuotaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The quota envelope is a singleton per project — import by service_name only.
	if req.ID == "" {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be the project service name")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), req.ID)...)
}

// putQuota applies the desired target spec under the given optimistic-concurrency
// checksum, polls for terminal status, and returns the final envelope. A
// concurrent change invalidates the checksum and surfaces as an error (no retry).
// The quota envelope is always present per project so this is also used for "create".
func (r *cloudQuotaResource) putQuota(ctx context.Context, serviceName, checksum string, spec CloudQuotaUpdateTargetSpecAPI) (*CloudQuotaAPIResponse, error) {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/quota"

	payload := CloudQuotaUpdateAPI{
		Checksum:   checksum,
		TargetSpec: spec,
	}

	var afterPut CloudQuotaAPIResponse
	if err := r.config.OVHClient.PutWithContext(ctx, endpoint, payload, &afterPut); err != nil {
		return nil, fmt.Errorf("error calling Put %s: %w", endpoint, err)
	}

	if _, err := r.waitForQuotaReady(ctx, serviceName); err != nil {
		return nil, err
	}

	// Final GET so caller works with the post-reconcile state.
	final := CloudQuotaAPIResponse{}
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &final); err != nil {
		return nil, fmt.Errorf("error fetching final quota envelope: %w", err)
	}
	return &final, nil
}

// quotaTaskErrorSuffix collects any task error messages from the envelope so a
// failed reconcile can surface the real cause instead of a bare status.
func quotaTaskErrorSuffix(res *CloudQuotaAPIResponse) string {
	var msgs []string
	for _, t := range res.CurrentTasks {
		for _, e := range t.Errors {
			if e.Message != "" {
				msgs = append(msgs, e.Message)
			}
		}
	}
	if len(msgs) == 0 {
		return ""
	}
	return ": " + strings.Join(msgs, "; ")
}

// waitForQuotaReady polls the quota envelope until ResourceStatus reaches READY.
// The pending states mirror common.ResourceStatusEnum's non-terminal values; a
// reconcile that lands in ERROR is surfaced immediately (with any task error
// messages) rather than spun on until the timeout.
func (r *cloudQuotaResource) waitForQuotaReady(ctx context.Context, serviceName string) (any, error) {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/quota"
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "OUT_OF_SYNC", "SUSPENDED", "UNKNOWN"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudQuotaAPIResponse{}
			if err := r.config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
				return res, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, fmt.Errorf("quota reconcile failed (status ERROR)%s", quotaTaskErrorSuffix(res))
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func (r *cloudQuotaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudQuotaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := data.toTargetSpecAPI(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Invalid quota target spec", err.Error())
		return
	}

	// Create has no prior state, so read the singleton envelope once to obtain
	// the current checksum for the optimistic-concurrency guard.
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/quota"
	var current CloudQuotaAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &current); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	final, err := r.putQuota(ctx, data.ServiceName.ValueString(), current.Checksum, spec)
	if err != nil {
		resp.Diagnostics.AddError("Error applying quota target spec", err.Error())
		return
	}

	data.MergeWith(ctx, final)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudQuotaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudQuotaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/quota"
	var responseData CloudQuotaAPIResponse
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

func (r *cloudQuotaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData CloudQuotaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := planData.toTargetSpecAPI(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Invalid quota target spec", err.Error())
		return
	}

	final, err := r.putQuota(ctx, planData.ServiceName.ValueString(), stateData.Checksum.ValueString(), spec)
	if err != nil {
		resp.Diagnostics.AddError("Error applying quota target spec", err.Error())
		return
	}

	planData.MergeWith(ctx, final)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudQuotaResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// The quota envelope is a singleton owned by the project and cannot be
	// deleted. Removing the resource from Terraform state is the only action.
}
