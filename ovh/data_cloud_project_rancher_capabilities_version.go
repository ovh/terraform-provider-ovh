package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectRancherCapabilitiesVersionDataSource)(nil)

func NewCloudProjectRancherCapabilitiesVersionDataSource() datasource.DataSource {
	return &cloudProjectRancherCapabilitiesVersionDataSource{}
}

type cloudProjectRancherCapabilitiesVersionDataSource struct {
	config *Config
}

func (d *cloudProjectRancherCapabilitiesVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_rancher_capabilities_version"
}

func (d *cloudProjectRancherCapabilitiesVersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectRancherCapabilitiesVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectRancherCapabilitiesVersionDataSourceSchema(ctx)
}

func (d *cloudProjectRancherCapabilitiesVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectRancherCapabilitiesVersionModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.RancherId.ValueString()) + "/capabilities/version"

	if err := d.config.OVHClient.Get(endpoint, &data.Versions); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
