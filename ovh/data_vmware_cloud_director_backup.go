package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*vmwareCloudDirectorBackupDataSource)(nil)

func NewVmwareCloudDirectorBackupDataSource() datasource.DataSource {
	return &vmwareCloudDirectorBackupDataSource{}
}

type vmwareCloudDirectorBackupDataSource struct {
	config *Config
}

func (d *vmwareCloudDirectorBackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vmware_cloud_director_backup"
}

func (d *vmwareCloudDirectorBackupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vmwareCloudDirectorBackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VmwareCloudDirectorBackupDataSourceSchema(ctx)
}

func (d *vmwareCloudDirectorBackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmwareCloudDirectorBackupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/v2/vmwareCloudDirector/backup/" + url.PathEscape(data.BackupId.ValueString())

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
