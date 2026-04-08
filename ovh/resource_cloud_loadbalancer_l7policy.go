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

var (
	_ resource.Resource                = (*cloudLoadbalancerL7PolicyResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudLoadbalancerL7PolicyResource)(nil)
	_ resource.ResourceWithImportState = (*cloudLoadbalancerL7PolicyResource)(nil)
)

func NewCloudLoadbalancerL7PolicyResource() resource.Resource {
	return &cloudLoadbalancerL7PolicyResource{}
}

type cloudLoadbalancerL7PolicyResource struct {
	config *Config
}

func (r *cloudLoadbalancerL7PolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_l7policy"
}

func (r *cloudLoadbalancerL7PolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var l7PolicyMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description", "action", "redirect_prefix", "redirect_url", "redirect_pool_id"},
	Int64s:  []string{"position", "redirect_http_code"},
}

func (r *cloudLoadbalancerL7PolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an L7 policy on a load balancer listener in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Required — immutable
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"loadbalancer_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the load balancer",
				MarkdownDescription: "ID of the load balancer",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the listener",
				MarkdownDescription: "ID of the listener",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Required — mutable
			"action": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Action of the L7 policy (REDIRECT_PREFIX, REDIRECT_TO_POOL, REDIRECT_TO_URL, REJECT)",
				MarkdownDescription: "Action of the L7 policy (`REDIRECT_PREFIX`, `REDIRECT_TO_POOL`, `REDIRECT_TO_URL`, `REJECT`)",
			},

			// Optional — mutable
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Name of the L7 policy",
				MarkdownDescription: "Name of the L7 policy",
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Description of the L7 policy",
				MarkdownDescription: "Description of the L7 policy",
			},
			"position": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Position of the L7 policy",
				MarkdownDescription: "Position of the L7 policy",
			},
			"redirect_prefix": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Redirect prefix for REDIRECT_PREFIX action",
				MarkdownDescription: "Redirect prefix for `REDIRECT_PREFIX` action",
			},
			"redirect_url": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Redirect URL for REDIRECT_TO_URL action",
				MarkdownDescription: "Redirect URL for `REDIRECT_TO_URL` action",
			},
			"redirect_http_code": schema.Int64Attribute{
				Optional:            true,
				Description:         "HTTP redirect code (301, 302, 303, 307, 308)",
				MarkdownDescription: "HTTP redirect code (`301`, `302`, `303`, `307`, `308`)",
			},
			"redirect_pool_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "ID of the pool for REDIRECT_TO_POOL action",
				MarkdownDescription: "ID of the pool for `REDIRECT_TO_POOL` action",
			},
			"rules": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "List of L7 rules for this policy",
				MarkdownDescription: "List of L7 rules for this policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Type of the L7 rule (COOKIE, FILE_TYPE, HEADER, HOST_NAME, PATH)",
							MarkdownDescription: "Type of the L7 rule (`COOKIE`, `FILE_TYPE`, `HEADER`, `HOST_NAME`, `PATH`)",
						},
						"compare_type": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Comparison type (CONTAINS, ENDS_WITH, EQUAL_TO, REGEX, STARTS_WITH)",
							MarkdownDescription: "Comparison type (`CONTAINS`, `ENDS_WITH`, `EQUAL_TO`, `REGEX`, `STARTS_WITH`)",
						},
						"value": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Value to compare against",
							MarkdownDescription: "Value to compare against",
						},
						"key": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Key for COOKIE and HEADER rule types",
							MarkdownDescription: "Key for `COOKIE` and `HEADER` rule types",
						},
						"invert": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Whether to invert the rule match",
							MarkdownDescription: "Whether to invert the rule match",
						},
					},
				},
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "L7 policy ID",
				MarkdownDescription: "L7 policy ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(l7PolicyMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the L7 policy",
				MarkdownDescription: "Creation date of the L7 policy",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the L7 policy",
				MarkdownDescription: "Last update date of the L7 policy",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(l7PolicyMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "L7 policy readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "L7 policy readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the L7 policy",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(l7PolicyMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "L7 policy name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "L7 policy description",
					},
					"action": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "L7 policy action",
					},
					"position": schema.Int64Attribute{
						Computed:    true,
						Description: "L7 policy position",
					},
					"redirect_prefix": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Redirect prefix",
					},
					"redirect_url": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Redirect URL",
					},
					"redirect_http_code": schema.Int64Attribute{
						Computed:    true,
						Description: "HTTP redirect code",
					},
					"redirect_pool_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Redirect pool ID",
					},
					"operating_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Operating status of the L7 policy",
					},
					"provisioning_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Provisioning status of the L7 policy",
					},
					"rules": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Current state of the L7 rules",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Rule ID",
								},
								"type": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Rule type",
								},
								"compare_type": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Comparison type",
								},
								"value": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Value to compare against",
								},
								"key": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Key for COOKIE and HEADER rule types",
								},
								"invert": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the rule match is inverted",
								},
								"operating_status": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Operating status of the rule",
								},
								"provisioning_status": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Provisioning status of the rule",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *cloudLoadbalancerL7PolicyResource) l7PolicyEndpoint(serviceName, lbId, listenerId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/loadbalancer/" + url.PathEscape(lbId) + "/listener/" + url.PathEscape(listenerId) + "/l7policy"
}

func (r *cloudLoadbalancerL7PolicyResource) l7PolicyItemEndpoint(serviceName, lbId, listenerId, l7PolicyId string) string {
	return r.l7PolicyEndpoint(serviceName, lbId, listenerId) + "/" + url.PathEscape(l7PolicyId)
}

func (r *cloudLoadbalancerL7PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 4 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<loadbalancer_id>/<listener_id>/<l7policy_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("loadbalancer_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), splits[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[3])...)
}

func (r *cloudLoadbalancerL7PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudLoadbalancerL7PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := r.l7PolicyEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString())

	var responseData CloudLoadbalancerL7PolicyAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the resource ID is tracked even if the workflow fails
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Wait for L7 policy to be READY
	_, err := r.waitForL7PolicyReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for L7 policy to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = r.l7PolicyItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), responseData.Id)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudLoadbalancerL7PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudLoadbalancerL7PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.l7PolicyItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerL7PolicyAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudLoadbalancerL7PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudLoadbalancerL7PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := r.l7PolicyItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerL7PolicyAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for L7 policy to be READY
	_, err := r.waitForL7PolicyReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for L7 policy to be ready after update",
			err.Error(),
		)
		return
	}

	// Read final state
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	planData.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudLoadbalancerL7PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudLoadbalancerL7PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.l7PolicyItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for deletion to complete
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerL7PolicyAPIResponse{}
			endpoint := r.l7PolicyItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for L7 policy to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudLoadbalancerL7PolicyResource) waitForL7PolicyReady(ctx context.Context, serviceName, lbId, listenerId, l7PolicyId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerL7PolicyAPIResponse{}
			endpoint := r.l7PolicyItemEndpoint(serviceName, lbId, listenerId, l7PolicyId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
