package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*dedicatedServerSpecificationsHardwareDataSource)(nil)

func NewDedicatedServerSpecificationsHardwareDataSource() datasource.DataSource {
	return &dedicatedServerSpecificationsHardwareDataSource{}
}

type dedicatedServerSpecificationsHardwareDataSource struct {
	config *Config
}

func (d *dedicatedServerSpecificationsHardwareDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_specifications_hardware"
}

func (d *dedicatedServerSpecificationsHardwareDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dedicatedServerSpecificationsHardwareDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DedicatedServerSpecificationsHardwareDataSourceSchema(ctx)
}

func (d *dedicatedServerSpecificationsHardwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DedicatedServerSpecificationsHardwareModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/dedicated/server/" + url.PathEscape(data.ServiceName.ValueString()) + "/specifications/hardware"
	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
