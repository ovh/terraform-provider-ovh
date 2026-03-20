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
	_ resource.Resource                = (*cloudLoadbalancerListenerPoolResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudLoadbalancerListenerPoolResource)(nil)
	_ resource.ResourceWithImportState = (*cloudLoadbalancerListenerPoolResource)(nil)
)

func NewCloudLoadbalancerListenerPoolResource() resource.Resource {
	return &cloudLoadbalancerListenerPoolResource{}
}

type cloudLoadbalancerListenerPoolResource struct {
	config *Config
}

func (r *cloudLoadbalancerListenerPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_listener_pool"
}

func (r *cloudLoadbalancerListenerPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var poolMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description", "algorithm"},
	Objects: []string{"persistence"},
}

func (r *cloudLoadbalancerListenerPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a pool for a listener in a public cloud loadbalancer.",
		Attributes: map[string]schema.Attribute{
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
				Description:         "ID of the loadbalancer",
				MarkdownDescription: "ID of the loadbalancer",
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
			"protocol": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Protocol used by the pool (e.g. HTTP, HTTPS, PROXY, PROXYV2, SCTP, TCP, UDP)",
				MarkdownDescription: "Protocol used by the pool (e.g. `HTTP`, `HTTPS`, `PROXY`, `PROXYV2`, `SCTP`, `TCP`, `UDP`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Load balancing algorithm (e.g. LEAST_CONNECTIONS, ROUND_ROBIN, SOURCE_IP)",
				MarkdownDescription: "Load balancing algorithm (e.g. `LEAST_CONNECTIONS`, `ROUND_ROBIN`, `SOURCE_IP`)",
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Pool name",
				MarkdownDescription: "Pool name",
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Pool description",
				MarkdownDescription: "Pool description",
			},
			"persistence": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Session persistence configuration",
				MarkdownDescription: "Session persistence configuration",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Session persistence type (APP_COOKIE, HTTP_COOKIE, SOURCE_IP)",
						MarkdownDescription: "Session persistence type (`APP_COOKIE`, `HTTP_COOKIE`, `SOURCE_IP`)",
					},
					"cookie_name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "Cookie name for APP_COOKIE persistence type",
						MarkdownDescription: "Cookie name for `APP_COOKIE` persistence type",
					},
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Pool ID",
				MarkdownDescription: "Pool ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(poolMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the pool",
				MarkdownDescription: "Creation date of the pool",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the pool",
				MarkdownDescription: "Last update date of the pool",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(poolMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Pool readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Pool readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the pool",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(poolMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Pool name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Pool description",
					},
					"protocol": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Protocol used by the pool",
					},
					"algorithm": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Load balancing algorithm",
					},
					"persistence": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Session persistence configuration",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Session persistence type",
							},
							"cookie_name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Cookie name for APP_COOKIE persistence type",
							},
						},
					},
					"operating_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Operating status of the pool",
					},
					"provisioning_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Provisioning status of the pool",
					},
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Region",
					},
					"availability_zone": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Availability zone",
					},
				},
			},
		},
	}
}

func (r *cloudLoadbalancerListenerPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 4 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<loadbalancer_id>/<listener_id>/<pool_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("loadbalancer_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), splits[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[3])...)
}

func (r *cloudLoadbalancerListenerPoolResource) poolEndpoint(serviceName, loadbalancerId, listenerId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/loadbalancer/" + url.PathEscape(loadbalancerId) + "/listener/" + url.PathEscape(listenerId) + "/pool"
}

func (r *cloudLoadbalancerListenerPoolResource) poolEndpointWithId(serviceName, loadbalancerId, listenerId, poolId string) string {
	return r.poolEndpoint(serviceName, loadbalancerId, listenerId) + "/" + url.PathEscape(poolId)
}

func (r *cloudLoadbalancerListenerPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudLoadbalancerListenerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := r.poolEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString())

	var responseData CloudLoadbalancerListenerPoolAPIResponse
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

	// Wait for pool to be READY
	_, err := r.waitForPoolReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for pool to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), responseData.Id)
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

func (r *cloudLoadbalancerListenerPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudLoadbalancerListenerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerListenerPoolAPIResponse
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

func (r *cloudLoadbalancerListenerPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudLoadbalancerListenerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerListenerPoolAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for pool to be READY
	_, err := r.waitForPoolReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for pool to be ready after update",
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

func (r *cloudLoadbalancerListenerPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudLoadbalancerListenerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())

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
			res := &CloudLoadbalancerListenerPoolAPIResponse{}
			endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.ListenerId.ValueString(), data.Id.ValueString())
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
			"Error waiting for pool to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudLoadbalancerListenerPoolResource) waitForPoolReady(ctx context.Context, serviceName, loadbalancerId, listenerId, poolId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerListenerPoolAPIResponse{}
			endpoint := r.poolEndpointWithId(serviceName, loadbalancerId, listenerId, poolId)
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
