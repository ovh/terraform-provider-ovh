package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerPoolDataSource)(nil)

func NewCloudLoadbalancerPoolDataSource() datasource.DataSource {
	return &cloudLoadbalancerPoolDataSource{}
}

type cloudLoadbalancerPoolDataSource struct {
	config *Config
}

func (d *cloudLoadbalancerPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pool"
}

func (d *cloudLoadbalancerPoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudLoadbalancerPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about a pool in a public cloud loadbalancer.",
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
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Pool ID",
			},
			"protocol": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Protocol used by the pool (e.g. HTTP, HTTPS, PROXY, PROXYV2, SCTP, TCP, UDP)",
			},
			"algorithm": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Load balancing algorithm (e.g. LEAST_CONNECTIONS, ROUND_ROBIN, SOURCE_IP)",
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
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the pool",
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

func (d *cloudLoadbalancerPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerPoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/pool/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudLoadbalancerPoolAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
