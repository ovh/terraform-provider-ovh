package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*ovhcloudConnectDatacentersDataSource)(nil)

func NewOvhcloudConnectDatacentersDataSource() datasource.DataSource {
	return &ovhcloudConnectDatacentersDataSource{}
}

type ovhcloudConnectDatacentersDataSource struct {
	config *Config
}

func (d *ovhcloudConnectDatacentersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connect_datacenters"
}

func (d *ovhcloudConnectDatacentersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ovhcloudConnectDatacentersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OvhcloudConnectDatacentersDataSourceSchema(ctx)
}

func (d *ovhcloudConnectDatacentersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		datacenterIDs []ovhtypes.TfInt64Value
		data          OvhcloudConnectDatacentersModel
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/datacenter/"

	if err := d.config.OVHClient.Get(endpoint, &datacenterIDs); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Fetch each occ datacenter config data
	for _, datacenterID := range datacenterIDs {
		var datacenter OvhcloudConnectDatacenterModel
		endpoint = "/ovhCloudConnect/" + url.PathEscape(data.ServiceName.ValueString()) + "/datacenter/" + url.PathEscape(datacenterID.String())
		if err := d.config.OVHClient.Get(endpoint, &datacenter); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		data.Datacenters = append(data.Datacenters, datacenter)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
