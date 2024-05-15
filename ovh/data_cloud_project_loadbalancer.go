package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectLoadbalancerDataSource)(nil)

func NewCloudProjectLoadbalancerDataSource() datasource.DataSource {
	return &cloudProjectLoadbalancerDataSource{}
}

type cloudProjectLoadbalancerDataSource struct {
	config *Config
}

func (d *cloudProjectLoadbalancerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_loadbalancer"
}

func (d *cloudProjectLoadbalancerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectLoadbalancerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectLoadbalancerDataSourceSchema(ctx)
}

func (d *cloudProjectLoadbalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectLoadbalancerModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() {
		data.ServiceName.StringValue = basetypes.NewStringValue(os.Getenv("OVH_CLOUD_PROJECT_SERVICE"))
	}

	// Read API call logic
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer/%s",
		url.PathEscape(data.ServiceName.ValueString()),
		url.PathEscape(data.RegionName.ValueString()),
		url.PathEscape(data.Id.ValueString()),
	)

	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError("Failed to get loadbalancer details", fmt.Sprintf("error calling GET %s: %s", endpoint, err))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
