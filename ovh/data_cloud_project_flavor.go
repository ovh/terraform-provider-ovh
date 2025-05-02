package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectFlavorDataSource)(nil)

func NewCloudProjectFlavorDataSource() datasource.DataSource {
	return &cloudProjectFlavorDataSource{}
}

type cloudProjectFlavorDataSource struct {
	config *Config
}

func (d *cloudProjectFlavorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_flavor"
}

func (d *cloudProjectFlavorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectFlavorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectFlavorDataSourceSchema(ctx)
}

func (d *cloudProjectFlavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectFlavorModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/flavor/" + url.PathEscape(data.Id.ValueString())

	var flavor CloudProjectFlavorModel
	if err := d.config.OVHClient.Get(endpoint, &flavor); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	flavor.ServiceName = data.ServiceName

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &flavor)...)
}
