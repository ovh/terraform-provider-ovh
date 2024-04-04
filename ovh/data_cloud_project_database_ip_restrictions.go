package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type cloudProjectDatabaseIPRestrictionsDataSource struct {
	config *Config
}

type ovhCloudProjectDatabaseIpRestrictionsModel struct {
	ServiceName types.String `tfsdk:"service_name"`
	Engine      types.String `tfsdk:"engine"`
	ClusterId   types.String `tfsdk:"cluster_id"`
	IPs         types.Set    `tfsdk:"ips"`
}

func NewCloudProjectDatabaseIPRestrictionsDataSource() datasource.DataSource {
	return &cloudProjectDatabaseIPRestrictionsDataSource{}
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &cloudProjectDatabaseIPRestrictionsDataSource{}
	_ datasource.DataSourceWithConfigure = &cloudProjectDatabaseIPRestrictionsDataSource{}
)

func (d *cloudProjectDatabaseIPRestrictionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_database_ip_restrictions"
}

func (d *cloudProjectDatabaseIPRestrictionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectDatabaseIPRestrictionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage: "Use ip_restrictions field in cloud_project_database datasource instead.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required:    os.Getenv("OVH_CLOUD_PROJECT_SERVICE") == "",
				Optional:    os.Getenv("OVH_CLOUD_PROJECT_SERVICE") != "",
				Description: "Service name",
			},
			"engine": schema.StringAttribute{
				Required:    true,
				Description: "Name of the engine of the service",
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "Cluster ID",
			},
			"ips": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of IP restriction",
			},
		},
	}
}

func (d *cloudProjectDatabaseIPRestrictionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ovhCloudProjectDatabaseIpRestrictionsModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ServiceName.IsNull() {
		config.ServiceName = basetypes.NewStringValue(os.Getenv("OVH_CLOUD_PROJECT_SERVICE"))
	}

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction",
		url.PathEscape(config.ServiceName.ValueString()),
		url.PathEscape(config.Engine.ValueString()),
		url.PathEscape(config.ClusterId.ValueString()),
	)

	log.Printf("[DEBUG] Will read IP restrictions from cluster %s from project %s",
		config.ClusterId.ValueString(), config.ServiceName.ValueString())

	res := make([]string, 0)
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &res); err != nil {
		resp.Diagnostics.AddError("Failed to get ip restrictions", fmt.Sprintf("error calling GET %s: %s", endpoint, err))
		return
	}

	ips := make([]attr.Value, 0, len(res))
	for _, ip := range res {
		ips = append(ips, types.StringValue(ip))
	}
	config.IPs, diags = types.SetValue(types.StringType, ips)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
