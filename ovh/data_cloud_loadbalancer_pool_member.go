package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerPoolMemberDataSource)(nil)

func NewCloudLoadbalancerPoolMemberDataSource() datasource.DataSource {
	return &cloudLoadbalancerPoolMemberDataSource{}
}

type cloudLoadbalancerPoolMemberDataSource struct {
	config *Config
}

func (d *cloudLoadbalancerPoolMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pool_member"
}

func (d *cloudLoadbalancerPoolMemberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudLoadbalancerPoolMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about a member in a pool of a public cloud loadbalancer.",
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
			"pool_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the pool",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Member ID",
			},
			"address": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "IP address of the member",
			},
			"protocol_port": schema.Int64Attribute{
				Computed:    true,
				Description: "Port used by the member to receive traffic",
			},
			"subnet_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "ID of the subnet the member is in",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Member name",
			},
			"weight": schema.Int64Attribute{
				Computed:    true,
				Description: "Weight of the member in the pool (0-256). Higher weight receives more traffic.",
			},
			"backup": schema.BoolAttribute{
				Computed:    true,
				Description: "When true, the member is a backup member and only receives traffic when all non-backup members are down",
			},
			"monitor": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Optional health monitor address and port override for this member",
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "IP address used by the health monitor for this member",
					},
					"port": schema.Int64Attribute{
						Computed:    true,
						Description: "Port used by the health monitor for this member",
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
				Description: "Creation date of the member",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the member",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Member readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the member",
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

func (d *cloudLoadbalancerPoolMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerPoolMemberModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/pool/" + url.PathEscape(data.PoolId.ValueString()) +
		"/member/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudLoadbalancerPoolMemberAPIResponse
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
