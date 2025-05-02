package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*ovhcloudConnectConfigPopDatacentersDataSource)(nil)

func NewOvhcloudConnectConfigPopDatacentersDataSource() datasource.DataSource {
	return &ovhcloudConnectConfigPopDatacentersDataSource{}
}

type ovhcloudConnectConfigPopDatacentersDataSource struct {
	config *Config
}

func (d *ovhcloudConnectConfigPopDatacentersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connect_config_pop_datacenters"
}

func (d *ovhcloudConnectConfigPopDatacentersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ovhcloudConnectConfigPopDatacentersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OvhcloudConnectConfigPopDatacentersDataSourceSchema(ctx)
}

func (d *ovhcloudConnectConfigPopDatacentersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		datacenterIDs []ovhtypes.TfInt64Value
		data          OvhcloudConnectConfigPopDatacentersModel
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/"

	if err := d.config.OVHClient.Get(endpoint, &datacenterIDs); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Fetch each occ datacenter config data
	for _, datacenterID := range datacenterIDs {
		var datacenterConfig OvhcloudConnectConfigPopDatacenterModel
		endpoint = "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/" + url.PathEscape(datacenterID.String())
		if err := d.config.OVHClient.Get(endpoint, &datacenterConfig); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		data.DatacenterConfigs = append(data.DatacenterConfigs, datacenterConfig)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
