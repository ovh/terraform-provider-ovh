package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectsDataSource)(nil)

func NewCloudProjectsDataSource() datasource.DataSource {
	return &cloudProjectsDataSource{}
}

type cloudProjectsDataSource struct {
	config *Config
}

func (d *cloudProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_projects"
}

func (d *cloudProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectsDataSourceSchema(ctx)
}

func (d *cloudProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		projectsIDs []string
		data        CloudProjectsModel
	)

	// Retrieve list of projects
	if err := d.config.OVHClient.Get("/cloud/project", &projectsIDs); err != nil {
		resp.Diagnostics.AddError("Error calling Get /cloud/project", err.Error())
		return
	}

	// Fetch each project data
	for _, projectID := range projectsIDs {
		var projectData CloudProjectModel
		endpoint := "/cloud/project/" + url.PathEscape(projectID)
		if err := d.config.OVHClient.Get(endpoint, &projectData); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		projectData.ServiceName = ovhtypes.NewTfStringValue(projectID)

		data.Projects = append(data.Projects, projectData)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
