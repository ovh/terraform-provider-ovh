package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudRegionDataSource)(nil)

func NewCloudRegionDataSource() datasource.DataSource {
	return &cloudRegionDataSource{}
}

type cloudRegionDataSource struct {
	config *Config
}

func (d *cloudRegionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_region"
}

func (d *cloudRegionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudRegionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := cloudRegionDetailAttributes()
	attributes["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Service name of the resource representing the id of the cloud project",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project",
	}
	attributes["name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Name of the region (e.g. GRA11)",
		MarkdownDescription: "Name of the region (e.g. GRA11)",
	}

	resp.Schema = schema.Schema{
		Description:         "Retrieve information about a single region of a public cloud project.",
		MarkdownDescription: "Retrieve information about a single region of a public cloud project.",
		Attributes:          attributes,
	}
}

func (d *cloudRegionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudRegionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/reference/region/" + url.PathEscape(data.Name.ValueString())

	var region CloudRegion
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &region); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.MergeWith(ctx, &region)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
