package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudSshKeysDataSource)(nil)

func NewCloudSshKeysDataSource() datasource.DataSource {
	return &cloudSshKeysDataSource{}
}

type cloudSshKeysDataSource struct {
	config *Config
}

func (d *cloudSshKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_ssh_keys"
}

func (d *cloudSshKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudSshKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudSshKeysDataSourceSchema(ctx)
}

func (d *cloudSshKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudSshKeysModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := resolveCloudSshKeyServiceName(&data.ServiceName, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := cloudSshKeyBaseEndpoint(serviceName)

	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &data.SshKeys); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
