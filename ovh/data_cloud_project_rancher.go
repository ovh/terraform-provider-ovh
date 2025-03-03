package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectRancherDataSource)(nil)

func NewCloudProjectRancherDataSource() datasource.DataSource {
	return &cloudProjectRancherDataSource{}
}

type cloudProjectRancherDataSource struct {
	config *Config
}

func (d *cloudProjectRancherDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_rancher"
}

func (d *cloudProjectRancherDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectRancherDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectRancherDataSourceSchema(ctx)
}

func (d *cloudProjectRancherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectRancherModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.Id.ValueString())
	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
