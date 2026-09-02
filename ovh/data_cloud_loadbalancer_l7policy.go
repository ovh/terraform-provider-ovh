package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerL7PolicyDataSource)(nil)

func NewCloudLoadbalancerL7PolicyDataSource() datasource.DataSource {
	return &cloudLoadbalancerL7PolicyDataSource{}
}

type cloudLoadbalancerL7PolicyDataSource struct {
	config *Config
}

func (d *cloudLoadbalancerL7PolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_l7policy"
}

func (d *cloudLoadbalancerL7PolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudLoadbalancerL7PolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get information about an L7 policy on a load balancer listener in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Required lookup keys
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
			"listener_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the listener",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "L7 policy ID",
			},

			// Computed
			"action": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Action of the L7 policy (REDIRECT_PREFIX, REDIRECT_TO_POOL, REDIRECT_TO_URL, REJECT)",
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Name of the L7 policy",
			},
			"description": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Description of the L7 policy",
			},
			"position": schema.Int64Attribute{
				Computed:    true,
				Description: "Position of the L7 policy",
			},
			"redirect_prefix": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Redirect prefix for REDIRECT_PREFIX action",
			},
			"redirect_url": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Redirect URL for REDIRECT_TO_URL action",
			},
			"redirect_http_code": schema.Int64Attribute{
				Computed:    true,
				Description: "HTTP redirect code (301, 302, 303, 307, 308)",
			},
			"redirect_pool_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "ID of the pool for REDIRECT_TO_POOL action",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of L7 rules for this policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Type of the L7 rule (COOKIE, FILE_TYPE, HEADER, HOST_NAME, PATH)",
						},
						"compare_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Comparison type (CONTAINS, ENDS_WITH, EQUAL_TO, REGEX, STARTS_WITH)",
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
							Description: "Whether to invert the rule match",
						},
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
				Description: "Creation date of the L7 policy",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update date of the L7 policy",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "L7 policy readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the L7 policy",
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

func (d *cloudLoadbalancerL7PolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerL7PolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/listener/" + url.PathEscape(data.ListenerId.ValueString()) +
		"/l7policy/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudLoadbalancerL7PolicyAPIResponse
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
