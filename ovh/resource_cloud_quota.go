package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// ------------------------------------------------------------------
// Resource model — distinct from CloudQuotaModel (data source) since
// target_spec is a user-provided Required block here.
// ------------------------------------------------------------------

// CloudQuotaResourceModel is the Terraform model for the writable quota resource.
type CloudQuotaResourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`

	// User input
	TargetSpec types.Object `tfsdk:"target_spec"`

	// Computed envelope
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// ------------------------------------------------------------------
// PUT payload types
// ------------------------------------------------------------------

// CloudQuotaUpdateTargetSpecAPI is the wire payload for the update target spec.
type CloudQuotaUpdateTargetSpecAPI struct {
	ManualQuota bool                            `json:"manualQuota"`
	Regions     []CloudQuotaRegionTargetSpecAPI `json:"regions"`
}

// CloudQuotaUpdateAPI is the wire payload for PUT /quota.
type CloudQuotaUpdateAPI struct {
	Checksum   string                        `json:"checksum"`
	TargetSpec CloudQuotaUpdateTargetSpecAPI `json:"targetSpec"`
}

// ------------------------------------------------------------------
// Helpers — extract target_spec from the Terraform model into the API payload
// ------------------------------------------------------------------

type quotaTargetSpecRegionPlan struct {
	Region  ovhtypes.TfStringValue `tfsdk:"region"`
	Profile ovhtypes.TfStringValue `tfsdk:"profile"`
}

type quotaTargetSpecPlan struct {
	ManualQuota types.Bool                  `tfsdk:"manual_quota"`
	Regions     []quotaTargetSpecRegionPlan `tfsdk:"regions"`
}

func (m *CloudQuotaResourceModel) targetSpecToAPI(ctx context.Context) (CloudQuotaUpdateTargetSpecAPI, error) {
	var plan quotaTargetSpecPlan
	diags := m.TargetSpec.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if diags.HasError() {
		return CloudQuotaUpdateTargetSpecAPI{}, fmt.Errorf("invalid target_spec: %s", diags)
	}

	regions := make([]CloudQuotaRegionTargetSpecAPI, 0, len(plan.Regions))
	for _, r := range plan.Regions {
		regions = append(regions, CloudQuotaRegionTargetSpecAPI{
			Location: &CloudQuotaLocationAPI{Region: r.Region.ValueString()},
			Profile:  r.Profile.ValueString(),
		})
	}

	return CloudQuotaUpdateTargetSpecAPI{
		ManualQuota: plan.ManualQuota.ValueBool(),
		Regions:     regions,
	}, nil
}

// MergeWith copies API response into the resource model. The user-provided
// target_spec block must remain stable, so we rebuild it from
// response.TargetSpec exactly like the data source does.
func (m *CloudQuotaResourceModel) MergeWith(_ context.Context, response *CloudQuotaAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.TargetSpec = buildQuotaTargetSpecObject(response.TargetSpec)
	m.CurrentState = buildQuotaCurrentStateObject(response.CurrentState)
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
	Objects: []string{"target_spec"},
}

func (r *cloudQuotaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the public cloud project quota: applied profile per region and manual-quota toggle. There is exactly one quota envelope per project; importing a project's quota uses the service name as id.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the cloud project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_spec": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Desired quota specification for the project. Regions omitted are left unchanged upstream.",
				Attributes: map[string]schema.Attribute{
					"manual_quota": schema.BoolAttribute{
						Required:    true,
						Description: "When true, automatic quota upgrades are disabled for this project.",
					},
					"regions": schema.ListNestedAttribute{
						Required:    true,
						Description: "Target quota profile per region.",
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
				Description: "Quota readiness (CREATING, UPDATING, PENDING, ERROR, OUT_OF_SYNC, READY).",
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
		"manual_quota": schema.BoolAttribute{Computed: true},
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

// putQuota applies the desired target spec, polls for terminal status, and
// returns the final envelope. The quota envelope is always present per
// project so this is also used for "create".
func (r *cloudQuotaResource) putQuota(ctx context.Context, serviceName string, spec CloudQuotaUpdateTargetSpecAPI) (*CloudQuotaAPIResponse, error) {
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/quota"

	// Fetch current envelope to get the latest checksum.
	var current CloudQuotaAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &current); err != nil {
		return nil, fmt.Errorf("error fetching current quota envelope: %w", err)
	}

	payload := CloudQuotaUpdateAPI{
		Checksum:   current.Checksum,
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

func (r *cloudQuotaResource) waitForQuotaReady(ctx context.Context, serviceName string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING"},
		Target:  []string{"READY", "OUT_OF_SYNC", "ERROR"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudQuotaAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/quota"
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func (r *cloudQuotaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudQuotaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := data.targetSpecToAPI(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Invalid target_spec", err.Error())
		return
	}

	final, err := r.putQuota(ctx, data.ServiceName.ValueString(), spec)
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
	var planData CloudQuotaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := planData.targetSpecToAPI(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Invalid target_spec", err.Error())
		return
	}

	final, err := r.putQuota(ctx, planData.ServiceName.ValueString(), spec)
	if err != nil {
		resp.Diagnostics.AddError("Error applying quota target spec", err.Error())
		return
	}

	planData.MergeWith(ctx, final)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudQuotaResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// The quota envelope is a singleton owned by the project and cannot be
	// deleted. Removing the resource from Terraform state is the only action.
	_ = ctx
}
