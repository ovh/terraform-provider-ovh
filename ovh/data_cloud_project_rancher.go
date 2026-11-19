package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/datasources/cloud_project_rancher"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
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
	resp.Schema = cloud_project_rancher.CloudProjectRancherDataSourceSchema(ctx)
}

func (d *cloudProjectRancherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloud_project_rancher.CloudProjectRancherModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectId = cloudProjectRancherProjectIdOrEnv(data.ProjectId, &resp.Diagnostics)
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

// cloudProjectRancherProjectIdOrEnv returns the project_id of an
// ovh_cloud_project_rancher* data source: the configured value, or the
// OVH_CLOUD_PROJECT_SERVICE environment variable when the attribute is omitted
// or empty. An error diagnostic is added when neither is set.
func cloudProjectRancherProjectIdOrEnv(value ovhtypes.TfStringValue, diags *diag.Diagnostics) ovhtypes.TfStringValue {
	if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
		return value
	}
	if env := os.Getenv("OVH_CLOUD_PROJECT_SERVICE"); env != "" {
		return ovhtypes.NewTfStringValue(env)
	}
	diags.AddError(
		"Missing project_id",
		"The project_id attribute is required. Provide it in the data source "+
			"configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
	)
	return value
}
