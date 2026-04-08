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
	_ resource.Resource                = (*cloudLoadbalancerPoolResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudLoadbalancerPoolResource)(nil)
	_ resource.ResourceWithImportState = (*cloudLoadbalancerPoolResource)(nil)
)

func NewCloudLoadbalancerPoolResource() resource.Resource {
	return &cloudLoadbalancerPoolResource{}
}

type cloudLoadbalancerPoolResource struct {
	config *Config
}

func (r *cloudLoadbalancerPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pool"
}

func (r *cloudLoadbalancerPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	Objects: []string{"persistence", "health_monitor"},
}

func (r *cloudLoadbalancerPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a pool in a public cloud loadbalancer.",
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
			"health_monitor": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Health monitor configuration",
				MarkdownDescription: "Health monitor configuration",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Health monitor type (HTTP, HTTPS, PING, TCP, UDP_CONNECT, SCTP, TLS_HELLO)",
						MarkdownDescription: "Health monitor type (`HTTP`, `HTTPS`, `PING`, `TCP`, `UDP_CONNECT`, `SCTP`, `TLS_HELLO`)",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"delay": schema.Int64Attribute{
						Required:            true,
						Description:         "Seconds between health checks",
						MarkdownDescription: "Seconds between health checks",
					},
					"timeout": schema.Int64Attribute{
						Required:            true,
						Description:         "Seconds to wait for a health check response",
						MarkdownDescription: "Seconds to wait for a health check response",
					},
					"max_retries": schema.Int64Attribute{
						Required:            true,
						Description:         "Number of consecutive health check failures before marking member as unhealthy (1-10)",
						MarkdownDescription: "Number of consecutive health check failures before marking member as unhealthy (1-10)",
					},
					"max_retries_down": schema.Int64Attribute{
						Optional:            true,
						Description:         "Number of consecutive health check failures before marking member as ERROR (1-10)",
						MarkdownDescription: "Number of consecutive health check failures before marking member as ERROR (1-10)",
					},
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "Health monitor name",
						MarkdownDescription: "Health monitor name",
					},
					"url_path": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "URL path for HTTP/HTTPS health checks",
						MarkdownDescription: "URL path for HTTP/HTTPS health checks",
					},
					"http_method": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "HTTP method for health checks (GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS, TRACE)",
						MarkdownDescription: "HTTP method for health checks (`GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `TRACE`)",
					},
					"http_version": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "HTTP version for health checks (1.0 or 1.1)",
						MarkdownDescription: "HTTP version for health checks (`1.0` or `1.1`)",
					},
					"expected_codes": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "Expected HTTP response codes (e.g. 200, 200-202)",
						MarkdownDescription: "Expected HTTP response codes (e.g. `200`, `200-202`)",
					},
					"domain_name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "Domain name for health check requests",
						MarkdownDescription: "Domain name for health check requests",
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
					"health_monitor": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Health monitor configuration",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Health monitor ID",
							},
							"type": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Health monitor type",
							},
							"delay": schema.Int64Attribute{
								Computed:    true,
								Description: "Seconds between health checks",
							},
							"timeout": schema.Int64Attribute{
								Computed:    true,
								Description: "Seconds to wait for a health check response",
							},
							"max_retries": schema.Int64Attribute{
								Computed:    true,
								Description: "Number of consecutive health check failures before marking member as unhealthy",
							},
							"max_retries_down": schema.Int64Attribute{
								Computed:    true,
								Description: "Number of consecutive health check failures before marking member as ERROR",
							},
							"name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Health monitor name",
							},
							"url_path": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "URL path for HTTP/HTTPS health checks",
							},
							"http_method": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "HTTP method for health checks",
							},
							"http_version": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "HTTP version for health checks",
							},
							"expected_codes": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Expected HTTP response codes",
							},
							"domain_name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Domain name for health check requests",
							},
							"operating_status": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Operating status of the health monitor",
							},
							"provisioning_status": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Provisioning status of the health monitor",
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
				},
			},
		},
	}
}

func (r *cloudLoadbalancerPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<loadbalancer_id>/<pool_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("loadbalancer_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudLoadbalancerPoolResource) poolEndpoint(serviceName, loadbalancerId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/loadbalancer/" + url.PathEscape(loadbalancerId) + "/pool"
}

func (r *cloudLoadbalancerPoolResource) poolEndpointWithId(serviceName, loadbalancerId, poolId string) string {
	return r.poolEndpoint(serviceName, loadbalancerId) + "/" + url.PathEscape(poolId)
}

func (r *cloudLoadbalancerPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudLoadbalancerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := r.poolEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString())

	var responseData CloudLoadbalancerPoolAPIResponse
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
	_, err := r.waitForPoolReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for pool to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), responseData.Id)
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

func (r *cloudLoadbalancerPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudLoadbalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerPoolAPIResponse
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

func (r *cloudLoadbalancerPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudLoadbalancerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerPoolAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for pool to be READY
	_, err := r.waitForPoolReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())
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

func (r *cloudLoadbalancerPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudLoadbalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

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
			res := &CloudLoadbalancerPoolAPIResponse{}
			endpoint := r.poolEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())
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

func (r *cloudLoadbalancerPoolResource) waitForPoolReady(ctx context.Context, serviceName, loadbalancerId, poolId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerPoolAPIResponse{}
			endpoint := r.poolEndpointWithId(serviceName, loadbalancerId, poolId)
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
