package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*ovhcloudConnectConfigPopDatacenterExtrasDataSource)(nil)

func NewOvhcloudConnectConfigPopDatacenterExtrasDataSource() datasource.DataSource {
	return &ovhcloudConnectConfigPopDatacenterExtrasDataSource{}
}

type ovhcloudConnectConfigPopDatacenterExtrasDataSource struct {
	config *Config
}

func (d *ovhcloudConnectConfigPopDatacenterExtrasDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connect_config_pop_datacenter_extras"
}

func (d *ovhcloudConnectConfigPopDatacenterExtrasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ovhcloudConnectConfigPopDatacenterExtrasDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OvhcloudConnectConfigPopDatacenterExtrasDataSourceSchema(ctx)
}

func (d *ovhcloudConnectConfigPopDatacenterExtrasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data               OvhcloudConnectConfigPopDatacenterExtrasModel
		datacenterExtraIDs []ovhtypes.TfInt64Value
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/" + url.PathEscape(data.ConfigDatacenterId.String()) + "/extra/"

	if err := d.config.OVHClient.Get(endpoint, &datacenterExtraIDs); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Fetch each occ datacenter config data
	for _, extraID := range datacenterExtraIDs {
		var extraConfig OvhcloudConnectConfigPopDatacenterExtraModel
		endpoint = "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/config/pop/" + url.PathEscape(data.ConfigPopId.String()) + "/datacenter/" + url.PathEscape(data.ConfigDatacenterId.String()) + "/extra/" + url.PathEscape(extraID.String())
		if err := d.config.OVHClient.Get(endpoint, &extraConfig); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		data.ExtraConfigs = append(data.ExtraConfigs, extraConfig)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
