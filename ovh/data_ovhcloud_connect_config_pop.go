package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*ovhcloudConnectConfigPopDataSource)(nil)

func NewOvhcloudConnectConfigPopDataSource() datasource.DataSource {
	return &ovhcloudConnectConfigPopDataSource{}
}

type ovhcloudConnectConfigPopDataSource struct {
	config *Config
}

func (d *ovhcloudConnectConfigPopDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connect_config_pops"
}

func (d *ovhcloudConnectConfigPopDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ovhcloudConnectConfigPopDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OvhcloudConnectConfigPopDataSourceSchema(ctx)
}

func (d *ovhcloudConnectConfigPopDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		popIDs []ovhtypes.TfInt64Value
		data   OvhcloudConnectConfigsPopModel
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/"

	// Retrieve list of pop config IDs
	if err := d.config.OVHClient.Get(endpoint, &popIDs); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Fetch each occ pop config data
	for _, popID := range popIDs {
		var popConfig OvhcloudConnectConfigPopModel
		endpoint = "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(popID.String())
		if err := d.config.OVHClient.Get(endpoint, &popConfig); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		popConfig.PopId = popID
		popConfig.ServiceName = data.ServiceName

		data.PopConfigs = append(data.PopConfigs, popConfig)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
