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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudLoadbalancerListenerResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudLoadbalancerListenerResource)(nil)
	_ resource.ResourceWithImportState = (*cloudLoadbalancerListenerResource)(nil)
)

func NewCloudLoadbalancerListenerResource() resource.Resource {
	return &cloudLoadbalancerListenerResource{}
}

type cloudLoadbalancerListenerResource struct {
	config *Config
}

func (r *cloudLoadbalancerListenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_listener"
}

func (r *cloudLoadbalancerListenerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var listenerMutableAttrs = MutableAttrs{
	Strings:           []string{"name", "description", "default_tls_container_ref", "default_pool_id"},
	Int64s:            []string{"connection_limit", "timeout_client_data", "timeout_member_data", "timeout_member_connect", "timeout_tcp_inspect"},
	Objects:           []string{"insert_headers"},
	CustomStringLists: []string{"allowed_cidrs", "sni_container_refs", "tls_versions"},
}

func (r *cloudLoadbalancerListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a listener on a load balancer in a public cloud project.",
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
			"protocol": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Protocol of the listener (e.g., HTTP, HTTPS, TCP, UDP)",
				MarkdownDescription: "Protocol of the listener (e.g., `HTTP`, `HTTPS`, `TCP`, `UDP`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol_port": schema.Int64Attribute{
				Required:            true,
				Description:         "Port number the listener listens on",
				MarkdownDescription: "Port number the listener listens on",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},

			// Required — mutable
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the listener",
				MarkdownDescription: "Name of the listener",
			},

			// Optional — mutable
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Description of the listener",
				MarkdownDescription: "Description of the listener",
			},
			"connection_limit": schema.Int64Attribute{
				Optional:            true,
				Description:         "Maximum number of connections allowed",
				MarkdownDescription: "Maximum number of connections allowed",
			},
			"allowed_cidrs": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "List of CIDRs allowed to access the listener",
				MarkdownDescription: "List of CIDRs allowed to access the listener",
			},
			"timeout_client_data": schema.Int64Attribute{
				Optional:            true,
				Description:         "Timeout for client data in milliseconds",
				MarkdownDescription: "Timeout for client data in milliseconds",
			},
			"timeout_member_data": schema.Int64Attribute{
				Optional:            true,
				Description:         "Timeout for member data in milliseconds",
				MarkdownDescription: "Timeout for member data in milliseconds",
			},
			"timeout_member_connect": schema.Int64Attribute{
				Optional:            true,
				Description:         "Timeout for member connection in milliseconds",
				MarkdownDescription: "Timeout for member connection in milliseconds",
			},
			"timeout_tcp_inspect": schema.Int64Attribute{
				Optional:            true,
				Description:         "Timeout for TCP inspect in milliseconds",
				MarkdownDescription: "Timeout for TCP inspect in milliseconds",
			},
			"insert_headers": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Headers to insert into requests",
				MarkdownDescription: "Headers to insert into requests",
				Attributes: map[string]schema.Attribute{
					"x_forwarded_for": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-Forwarded-For header",
					},
					"x_forwarded_port": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-Forwarded-Port header",
					},
					"x_forwarded_proto": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-Forwarded-Proto header",
					},
					"x_ssl_client_verify": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-SSL-Client-Verify header",
					},
					"x_ssl_client_has_cert": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-SSL-Client-Has-Cert header",
					},
					"x_ssl_client_dn": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Insert X-SSL-Client-DN header",
					},
				},
			},
			"default_tls_container_ref": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Reference to the default TLS container",
				MarkdownDescription: "Reference to the default TLS container",
			},
			"default_pool_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "ID of the default pool for this listener",
				MarkdownDescription: "ID of the default pool for this listener",
			},
			"sni_container_refs": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "List of SNI container references",
				MarkdownDescription: "List of SNI container references",
			},
			"tls_versions": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "List of TLS versions allowed",
				MarkdownDescription: "List of TLS versions allowed",
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Listener ID",
				MarkdownDescription: "Listener ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(listenerMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the listener",
				MarkdownDescription: "Creation date of the listener",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the listener",
				MarkdownDescription: "Last update date of the listener",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(listenerMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Listener readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Listener readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the listener",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(listenerMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Listener name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Listener description",
					},
					"protocol": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Listener protocol",
					},
					"protocol_port": schema.Int64Attribute{
						Computed:    true,
						Description: "Port number the listener listens on",
					},
					"connection_limit": schema.Int64Attribute{
						Computed:    true,
						Description: "Maximum number of connections allowed",
					},
					"timeout_client_data": schema.Int64Attribute{
						Computed:    true,
						Description: "Timeout for client data in milliseconds",
					},
					"timeout_member_data": schema.Int64Attribute{
						Computed:    true,
						Description: "Timeout for member data in milliseconds",
					},
					"timeout_member_connect": schema.Int64Attribute{
						Computed:    true,
						Description: "Timeout for member connection in milliseconds",
					},
					"timeout_tcp_inspect": schema.Int64Attribute{
						Computed:    true,
						Description: "Timeout for TCP inspect in milliseconds",
					},
					"operating_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Operating status of the listener",
					},
					"provisioning_status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Provisioning status of the listener",
					},
					"default_tls_container_ref": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Reference to the default TLS container",
					},
					"default_pool_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "ID of the default pool for this listener",
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
					"insert_headers": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Headers inserted into requests",
						Attributes: map[string]schema.Attribute{
							"x_forwarded_for": schema.BoolAttribute{
								Computed:    true,
								Description: "X-Forwarded-For header insertion",
							},
							"x_forwarded_port": schema.BoolAttribute{
								Computed:    true,
								Description: "X-Forwarded-Port header insertion",
							},
							"x_forwarded_proto": schema.BoolAttribute{
								Computed:    true,
								Description: "X-Forwarded-Proto header insertion",
							},
							"x_ssl_client_verify": schema.BoolAttribute{
								Computed:    true,
								Description: "X-SSL-Client-Verify header insertion",
							},
							"x_ssl_client_has_cert": schema.BoolAttribute{
								Computed:    true,
								Description: "X-SSL-Client-Has-Cert header insertion",
							},
							"x_ssl_client_dn": schema.BoolAttribute{
								Computed:    true,
								Description: "X-SSL-Client-DN header insertion",
							},
						},
					},
					"allowed_cidrs": schema.ListAttribute{
						Computed:    true,
						ElementType: ovhtypes.TfStringType{},
						Description: "List of CIDRs allowed to access the listener",
					},
					"sni_container_refs": schema.ListAttribute{
						Computed:    true,
						ElementType: ovhtypes.TfStringType{},
						Description: "List of SNI container references",
					},
					"tls_versions": schema.ListAttribute{
						Computed:    true,
						ElementType: ovhtypes.TfStringType{},
						Description: "List of TLS versions allowed",
					},
				},
			},
		},
	}
}

func (r *cloudLoadbalancerListenerResource) listenerEndpoint(serviceName, loadbalancerId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/loadbalancer/" + url.PathEscape(loadbalancerId) + "/listener"
}

func (r *cloudLoadbalancerListenerResource) listenerItemEndpoint(serviceName, loadbalancerId, listenerId string) string {
	return r.listenerEndpoint(serviceName, loadbalancerId) + "/" + url.PathEscape(listenerId)
}

func (r *cloudLoadbalancerListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<loadbalancer_id>/<listener_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("loadbalancer_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudLoadbalancerListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudLoadbalancerListenerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := r.listenerEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString())

	var responseData CloudLoadbalancerListenerAPIResponse
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

	// Wait for listener to be READY
	_, err := r.waitForListenerReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for listener to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = r.listenerItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), responseData.Id)
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

func (r *cloudLoadbalancerListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudLoadbalancerListenerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.listenerItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerListenerAPIResponse
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

func (r *cloudLoadbalancerListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudLoadbalancerListenerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := r.listenerItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

	var responseData CloudLoadbalancerListenerAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for listener to be READY
	_, err := r.waitForListenerReady(ctx, data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for listener to be ready after update",
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

func (r *cloudLoadbalancerListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudLoadbalancerListenerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.listenerItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())

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
			res := &CloudLoadbalancerListenerAPIResponse{}
			endpoint := r.listenerItemEndpoint(data.ServiceName.ValueString(), data.LoadbalancerId.ValueString(), data.Id.ValueString())
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
			"Error waiting for listener to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudLoadbalancerListenerResource) waitForListenerReady(ctx context.Context, serviceName, loadbalancerId, listenerId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudLoadbalancerListenerAPIResponse{}
			endpoint := r.listenerItemEndpoint(serviceName, loadbalancerId, listenerId)
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
