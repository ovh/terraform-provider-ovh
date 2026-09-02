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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudLoadbalancerPoolMemberResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudLoadbalancerPoolMemberResource)(nil)
	_ resource.ResourceWithImportState = (*cloudLoadbalancerPoolMemberResource)(nil)
)

func NewCloudLoadbalancerPoolMemberResource() resource.Resource {
	return &cloudLoadbalancerPoolMemberResource{}
}

type cloudLoadbalancerPoolMemberResource struct {
	config *Config
}

func (r *cloudLoadbalancerPoolMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pool_member"
}

func (r *cloudLoadbalancerPoolMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var memberMutableAttrs = MutableAttrs{
	Strings: []string{"name"},
	Int64s:  []string{"weight"},
	Bools:   []string{"backup"},
	Objects: []string{"monitor"},
}

func (r *cloudLoadbalancerPoolMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a member in a pool of a public cloud loadbalancer.",
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
			"pool_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the pool",
				MarkdownDescription: "ID of the pool",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "IP address of the member",
				MarkdownDescription: "IP address of the member",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol_port": schema.Int64Attribute{
				Required:            true,
				Description:         "Port used by the member to receive traffic",
				MarkdownDescription: "Port used by the member to receive traffic",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "ID of the subnet the member is in",
				MarkdownDescription: "ID of the subnet the member is in",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Member name",
				MarkdownDescription: "Member name",
			},
			"weight": schema.Int64Attribute{
				Optional:            true,
				Description:         "Weight of the member in the pool (0-256). Higher weight receives more traffic.",
				MarkdownDescription: "Weight of the member in the pool (0-256). Higher weight receives more traffic.",
			},
			"backup": schema.BoolAttribute{
				Optional:            true,
				Description:         "When true, the member is a backup member and only receives traffic when all non-backup members are down",
				MarkdownDescription: "When true, the member is a backup member and only receives traffic when all non-backup members are down",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"monitor": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Optional health monitor address and port override for this member",
				MarkdownDescription: "Optional health monitor address and port override for this member",
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "IP address used by the health monitor for this member",
						MarkdownDescription: "IP address used by the health monitor for this member",
					},
					"port": schema.Int64Attribute{
						Optional:            true,
						Description:         "Port used by the health monitor for this member",
						MarkdownDescription: "Port used by the health monitor for this member",
					},
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Member ID",
				MarkdownDescription: "Member ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(memberMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the member",
				MarkdownDescription: "Creation date of the member",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the member",
				MarkdownDescription: "Last update date of the member",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(memberMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Member readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Member readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the member",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(memberMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Member name",
					},
					"address": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "IP address of the member",
					},
					"protocol_port": schema.Int64Attribute{
						Computed:    true,
						Description: "Port used by the member",
					},
					"weight": schema.Int64Attribute{
						Computed:    true,
						Description: "Weight of the member",
					},
					"subnet_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "ID of the subnet the member is in",
					},
					"operating_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Operating status of the member",
					},
					"provisioning_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Provisioning status of the member",
					},
					"backup": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this member is a backup member",
					},
					"monitor": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Health monitor address and port override for this member",
						Attributes: map[string]schema.Attribute{
							"address": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "IP address used by the health monitor",
							},
							"port": schema.Int64Attribute{
								Computed:    true,
								Description: "Port used by the health monitor",
							},
						},
					},
				},
			},
		},
	}
}

func (r *cloudLoadbalancerPoolMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 4 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<loadbalancer_id>/<pool_id>/<member_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("loadbalancer_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pool_id"), splits[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[3])...)
}

func (r *cloudLoadbalancerPoolMemberResource) memberEndpoint(serviceName, loadbalancerId, poolId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) +
		"/loadbalancer/" + url.PathEscape(loadbalancerId) +
		"/pool/" + url.PathEscape(poolId) +
		"/member"
}

func (r *cloudLoadbalancerPoolMemberResource) memberEndpointWithId(serviceName, loadbalancerId, poolId, memberId string) string {
	return r.memberEndpoint(serviceName, loadbalancerId, poolId) + "/" + url.PathEscape(memberId)
}

func (r *cloudLoadbalancerPoolMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudLoadbalancerPoolMemberModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := r.memberEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString())

	var responseData CloudLoadbalancerPoolMemberAPIResponse
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

	// Wait for member to be READY
	_, err := r.waitForMemberReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for member to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = r.memberEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), responseData.Id)
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

func (r *cloudLoadbalancerPoolMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudLoadbalancerPoolMemberModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.memberEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerPoolMemberAPIResponse
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

func (r *cloudLoadbalancerPoolMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudLoadbalancerPoolMemberModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := r.memberEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerPoolMemberAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for member to be READY
	_, err := r.waitForMemberReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for member to be ready after update",
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

func (r *cloudLoadbalancerPoolMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudLoadbalancerPoolMemberModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.memberEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), data.Id.ValueString())

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
			res := &CloudLoadbalancerPoolMemberAPIResponse{}
			ep := r.memberEndpointWithId(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.PoolId.ValueString(), data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, ep, res)
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
			"Error waiting for member to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudLoadbalancerPoolMemberResource) waitForMemberReady(ctx context.Context, serviceName, loadbalancerId, poolId, memberId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerPoolMemberAPIResponse{}
			endpoint := r.memberEndpointWithId(serviceName, loadbalancerId, poolId, memberId)
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
