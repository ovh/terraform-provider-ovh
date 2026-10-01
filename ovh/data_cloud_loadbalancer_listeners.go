package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerListenersDataSource)(nil)

func NewCloudLoadbalancerListenersDataSource() datasource.DataSource {
	return &cloudLoadbalancerListenersDataSource{}
}

type cloudLoadbalancerListenersDataSource struct {
	config *Config
}

// CloudLoadbalancerListenersModel is the model for the plural listeners data source.
type CloudLoadbalancerListenersModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	Listeners      types.List             `tfsdk:"listeners"`
}

func (d *cloudLoadbalancerListenersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_listeners"
}

func (d *cloudLoadbalancerListenersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

// ListenerListItemAttrTypes returns the attribute types for a single listener
// element of the plural data source list.
func ListenerListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               ovhtypes.TfStringType{},
		"name":             ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
		"protocol":         ovhtypes.TfStringType{},
		"protocol_port":    types.Int64Type,
		"connection_limit": types.Int64Type,
		"allowed_cidrs": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
		"timeout_client_data":    types.Int64Type,
		"timeout_member_data":    types.Int64Type,
		"timeout_member_connect": types.Int64Type,
		"timeout_tcp_inspect":    types.Int64Type,
		"insert_headers": types.ObjectType{
			AttrTypes: ListenerInsertHeadersAttrTypes(),
		},
		"default_tls_container_ref": ovhtypes.TfStringType{},
		"default_pool_id":           ovhtypes.TfStringType{},
		"sni_container_refs": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
		"tls_versions": types.ListType{
			ElemType: ovhtypes.TfStringType{},
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: ListenerCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudLoadbalancerListenersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the listeners of a load balancer in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"loadbalancer_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the load balancer",
			},
			"listeners": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of listeners",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Listener ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the listener",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Description of the listener",
						},
						"protocol": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Protocol of the listener (e.g., HTTP, HTTPS, TCP, UDP)",
						},
						"protocol_port": schema.Int64Attribute{
							Computed:    true,
							Description: "Port number the listener listens on",
						},
						"connection_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of connections allowed",
						},
						"allowed_cidrs": schema.ListAttribute{
							Computed:    true,
							ElementType: ovhtypes.TfStringType{},
							Description: "List of CIDRs allowed to access the listener",
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
						"insert_headers": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Headers to insert into requests",
							Attributes: map[string]schema.Attribute{
								"x_forwarded_for": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-Forwarded-For header",
								},
								"x_forwarded_port": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-Forwarded-Port header",
								},
								"x_forwarded_proto": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-Forwarded-Proto header",
								},
								"x_ssl_client_verify": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-SSL-Client-Verify header",
								},
								"x_ssl_client_has_cert": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-SSL-Client-Has-Cert header",
								},
								"x_ssl_client_dn": schema.BoolAttribute{
									Computed:    true,
									Description: "Insert X-SSL-Client-DN header",
								},
							},
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
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the listener",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the listener",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Listener readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": loadbalancerListenerDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudLoadbalancerListenersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerListenersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/listener"

	var responseData []CloudLoadbalancerListenerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: ListenerListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildListenerListItemObject(ctx, &responseData[i]))
	}

	data.Listeners = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// lbInt64OrNull turns an optional API int into an Int64 attribute value.
func lbInt64OrNull(v *int) basetypes.Int64Value {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

// lbStringOrNull turns an optional API string into a string attribute value.
func lbStringOrNull(v string) ovhtypes.TfStringValue {
	if v == "" {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
}

// buildListenerListItemObject builds a single listener list element from an API
// response, reusing the resource helpers for nested objects.
func buildListenerListItemObject(ctx context.Context, response *CloudLoadbalancerListenerAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	protocolVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	protocolPortVal := types.Int64Null()
	connectionLimitVal := types.Int64Null()
	timeoutClientDataVal := types.Int64Null()
	timeoutMemberDataVal := types.Int64Null()
	timeoutMemberConnectVal := types.Int64Null()
	timeoutTcpInspectVal := types.Int64Null()
	insertHeadersVal := types.ObjectNull(ListenerInsertHeadersAttrTypes())
	defaultTlsContainerRefVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	defaultPoolIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	allowedCidrsVal := types.ListNull(ovhtypes.TfStringType{})
	sniContainerRefsVal := types.ListNull(ovhtypes.TfStringType{})
	tlsVersionsVal := types.ListNull(ovhtypes.TfStringType{})

	if spec := response.TargetSpec; spec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Name)}
		descVal = lbStringOrNull(spec.Description)
		protocolVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Protocol)}
		protocolPortVal = types.Int64Value(int64(spec.ProtocolPort))
		connectionLimitVal = lbInt64OrNull(spec.ConnectionLimit)
		timeoutClientDataVal = lbInt64OrNull(spec.TimeoutClientData)
		timeoutMemberDataVal = lbInt64OrNull(spec.TimeoutMemberData)
		timeoutMemberConnectVal = lbInt64OrNull(spec.TimeoutMemberConnect)
		timeoutTcpInspectVal = lbInt64OrNull(spec.TimeoutTcpInspect)
		defaultTlsContainerRefVal = lbStringOrNull(spec.DefaultTlsContainerRef)

		if spec.InsertHeaders != nil {
			insertHeadersVal = buildInsertHeadersFromAPI(spec.InsertHeaders)
		}

		if spec.DefaultPool != nil && spec.DefaultPool.Id != "" {
			defaultPoolIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.DefaultPool.Id)}
		}

		allowedCidrsVal = buildStringListFromAPI(spec.AllowedCidrs)
		sniContainerRefsVal = buildStringListFromAPI(spec.SniContainerRefs)
		tlsVersionsVal = buildStringListFromAPI(spec.TlsVersions)
	}

	currentStateVal := types.ObjectNull(ListenerCurrentStateAttrTypes())
	if response.CurrentState != nil {
		currentStateVal = buildListenerCurrentStateObject(ctx, response.CurrentState)
	}

	obj, _ := types.ObjectValue(
		ListenerListItemAttrTypes(),
		map[string]attr.Value{
			"id":                        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":                      nameVal,
			"description":               descVal,
			"protocol":                  protocolVal,
			"protocol_port":             protocolPortVal,
			"connection_limit":          connectionLimitVal,
			"allowed_cidrs":             allowedCidrsVal,
			"timeout_client_data":       timeoutClientDataVal,
			"timeout_member_data":       timeoutMemberDataVal,
			"timeout_member_connect":    timeoutMemberConnectVal,
			"timeout_tcp_inspect":       timeoutTcpInspectVal,
			"insert_headers":            insertHeadersVal,
			"default_tls_container_ref": defaultTlsContainerRefVal,
			"default_pool_id":           defaultPoolIdVal,
			"sni_container_refs":        sniContainerRefsVal,
			"tls_versions":              tlsVersionsVal,
			"checksum":                  ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":                ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":                ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status":           ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":             currentStateVal,
		},
	)

	return obj
}
