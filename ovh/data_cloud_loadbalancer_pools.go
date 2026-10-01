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

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerPoolsDataSource)(nil)

func NewCloudLoadbalancerPoolsDataSource() datasource.DataSource {
	return &cloudLoadbalancerPoolsDataSource{}
}

type cloudLoadbalancerPoolsDataSource struct {
	config *Config
}

// CloudLoadbalancerPoolsModel is the model for the plural pools data source.
type CloudLoadbalancerPoolsModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	Pools          types.List             `tfsdk:"pools"`
}

func (d *cloudLoadbalancerPoolsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pools"
}

func (d *cloudLoadbalancerPoolsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// PoolListItemAttrTypes returns the attribute types for a single pool element
// of the plural data source list.
func PoolListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"protocol":    ovhtypes.TfStringType{},
		"algorithm":   ovhtypes.TfStringType{},
		"persistence": types.ObjectType{
			AttrTypes: poolPersistenceAttrTypes(),
		},
		"health_monitor": types.ObjectType{
			AttrTypes: healthMonitorAttrTypes(),
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: PoolCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudLoadbalancerPoolsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the pools of a load balancer in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"loadbalancer_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the loadbalancer",
			},
			"pools": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of pools",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Pool ID",
						},
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
							Description: "Protocol used by the pool (e.g. HTTP, HTTPS, PROXY, PROXYV2, SCTP, TCP, UDP)",
						},
						"algorithm": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Load balancing algorithm (LEAST_CONNECTIONS, ROUND_ROBIN, SOURCE_IP)",
						},
						"persistence": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Session persistence configuration",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Session persistence type (APP_COOKIE, HTTP_COOKIE, SOURCE_IP)",
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
								"type": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Health monitor type (HTTP, HTTPS, PING, TCP, UDP_CONNECT, SCTP, TLS_HELLO)",
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
									Description: "Number of consecutive health check failures before marking member as unhealthy (1-10)",
								},
								"max_retries_down": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of consecutive health check failures before marking member as ERROR (1-10)",
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
									Description: "HTTP method for health checks (GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS, TRACE)",
								},
								"http_version": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "HTTP version for health checks (1.0 or 1.1)",
								},
								"expected_codes": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Expected HTTP response codes (e.g. 200, 200-202)",
								},
								"domain_name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Domain name for health check requests",
								},
							},
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the pool",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the pool",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Pool readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": loadbalancerPoolDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudLoadbalancerPoolsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerPoolsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/pool"

	var responseData []CloudLoadbalancerPoolAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: PoolListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildPoolListItemObject(ctx, &responseData[i]))
	}

	data.Pools = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildPoolListItemObject builds a single pool list element from an API
// response, reusing the resource helpers for nested objects.
func buildPoolListItemObject(ctx context.Context, response *CloudLoadbalancerPoolAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	protocolVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	algorithmVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	persistenceVal := types.ObjectNull(poolPersistenceAttrTypes())
	healthMonitorVal := types.ObjectNull(healthMonitorAttrTypes())

	if spec := response.TargetSpec; spec != nil {
		nameVal = lbStringOrNull(spec.Name)
		descVal = lbStringOrNull(spec.Description)
		protocolVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Protocol)}
		algorithmVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Algorithm)}

		if spec.Persistence != nil {
			persistenceVal = buildPoolPersistenceObject(spec.Persistence)
		}

		if spec.HealthMonitor != nil {
			healthMonitorVal = buildHealthMonitorObject(spec.HealthMonitor, types.ObjectNull(healthMonitorAttrTypes()))
		}
	}

	currentStateVal := types.ObjectNull(PoolCurrentStateAttrTypes())
	if response.CurrentState != nil {
		currentStateVal = buildPoolCurrentStateObject(ctx, response.CurrentState)
	}

	obj, _ := types.ObjectValue(
		PoolListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":            nameVal,
			"description":     descVal,
			"protocol":        protocolVal,
			"algorithm":       algorithmVal,
			"persistence":     persistenceVal,
			"health_monitor":  healthMonitorVal,
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
		},
	)

	return obj
}
