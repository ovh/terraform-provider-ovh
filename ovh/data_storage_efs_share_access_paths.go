package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*storageEfsShareAccessPathsDataSource)(nil)

func NewStorageEfsShareAccessPathsDataSource() datasource.DataSource {
	return &storageEfsShareAccessPathsDataSource{}
}

type storageEfsShareAccessPathsDataSource struct {
	config *Config
}

func (d *storageEfsShareAccessPathsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_efs_share_access_paths"
}

func (d *storageEfsShareAccessPathsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *storageEfsShareAccessPathsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = StorageEfsShareAccessPathsDataSourceSchema(ctx)
}

func (d *storageEfsShareAccessPathsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StorageEfsShareAccessPathsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/accessPath"

	if err := d.config.OVHClient.Get(endpoint, &data.AccessPaths); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
