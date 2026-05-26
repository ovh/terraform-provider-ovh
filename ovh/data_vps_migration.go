package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*vpsMigrationDataSource)(nil)

func NewVpsMigrationDataSource() datasource.DataSource {
	return &vpsMigrationDataSource{}
}

type vpsMigrationDataSource struct {
	config *Config
}

func (d *vpsMigrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_migration"
}

func (d *vpsMigrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpsMigrationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VpsMigrationDataSourceSchema(ctx)
}

func (d *vpsMigrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VpsMigrationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	endpoint := "/vps/" + url.PathEscape(serviceName) + "/migration2020"
	api := &vpsMigrationAPIResponse{}
	if err := d.config.OVHClient.Get(endpoint, api); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	applyAPIResponse(ctx, &data, api)
	data.ID = ovhtypes.NewTfStringValue(serviceName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
