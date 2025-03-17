package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectImagesDataSource)(nil)

func NewCloudProjectImagesDataSource() datasource.DataSource {
	return &cloudProjectImagesDataSource{}
}

type cloudProjectImagesDataSource struct {
	config *Config
}

func (d *cloudProjectImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_images"
}

func (d *cloudProjectImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectImagesDataSourceSchema(ctx)
}

func (d *cloudProjectImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectImagesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Add query parameters
	queryParams := url.Values{}
	if !data.FlavorType.IsNull() && !data.FlavorType.IsUnknown() {
		queryParams.Add("flavorType", data.FlavorType.ValueString())
	}
	if !data.OsType.IsNull() && !data.OsType.IsUnknown() {
		queryParams.Add("osType", data.OsType.ValueString())
	}
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		queryParams.Add("region", data.Region.ValueString())
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/image?" + queryParams.Encode()
	if err := d.config.OVHClient.Get(endpoint, &data.Images); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
