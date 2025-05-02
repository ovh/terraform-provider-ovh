package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectRancherCapabilitiesPlanDataSource)(nil)

func NewCloudProjectRancherCapabilitiesPlanDataSource() datasource.DataSource {
	return &cloudProjectRancherCapabilitiesPlanDataSource{}
}

type cloudProjectRancherCapabilitiesPlanDataSource struct {
	config *Config
}

func (d *cloudProjectRancherCapabilitiesPlanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_rancher_capabilities_plan"
}

func (d *cloudProjectRancherCapabilitiesPlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectRancherCapabilitiesPlanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectRancherCapabilitiesPlanDataSourceSchema(ctx)
}

func (d *cloudProjectRancherCapabilitiesPlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectRancherCapabilitiesPlanModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/rancher/" + url.PathEscape(data.RancherId.ValueString()) + "/capabilities/plan"

	if err := d.config.OVHClient.Get(endpoint, &data.Plans); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
