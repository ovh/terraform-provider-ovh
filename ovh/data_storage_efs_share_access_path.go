package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*storageEfsShareAccessPathDataSource)(nil)

func NewStorageEfsShareAccessPathDataSource() datasource.DataSource {
	return &storageEfsShareAccessPathDataSource{}
}

type storageEfsShareAccessPathDataSource struct {
	config *Config
}

func (d *storageEfsShareAccessPathDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_efs_share_access_path"
}

func (d *storageEfsShareAccessPathDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *storageEfsShareAccessPathDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = StorageEfsShareAccessPathDataSourceSchema(ctx)
}

func (d *storageEfsShareAccessPathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StorageEfsShareAccessPathModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/accessPath/" + url.PathEscape(data.Id.ValueString()) + ""

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
