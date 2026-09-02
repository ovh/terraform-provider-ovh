package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerListenerDataSource)(nil)

func NewCloudLoadbalancerListenerDataSource() datasource.DataSource {
	return &cloudLoadbalancerListenerDataSource{}
}

type cloudLoadbalancerListenerDataSource struct {
	config *Config
}

func (d *cloudLoadbalancerListenerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_listener"
}

func (d *cloudLoadbalancerListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudLoadbalancerListenerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about a listener on a load balancer in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Required — lookup keys
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
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Listener ID",
			},

			// Computed — from targetSpec
			"protocol": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Protocol of the listener (e.g., HTTP, HTTPS, TCP, UDP)",
			},
			"protocol_port": schema.Int64Attribute{
				Computed:    true,
				Description: "Port number the listener listens on",
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

			// Computed — metadata
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
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the listener",
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

func (d *cloudLoadbalancerListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerListenerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/listener/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudLoadbalancerListenerAPIResponse
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
