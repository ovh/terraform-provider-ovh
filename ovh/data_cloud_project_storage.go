package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectStorageDataSource)(nil)

func NewCloudProjectStorageDataSource() datasource.DataSource {
	return &cloudProjectStorageDataSource{}
}

type cloudProjectStorageDataSource struct {
	config *Config
}

func (d *cloudProjectStorageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage"
}

func (d *cloudProjectStorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectStorageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectStorageDataSourceSchema(ctx)
}

func (d *cloudProjectStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectStorageModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	queryParams := url.Values{}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(data.Limit.ValueInt64(), 10))
	}
	if !data.Marker.IsNull() && !data.Marker.IsUnknown() {
		queryParams.Add("marker", data.Marker.ValueString())
	}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		queryParams.Add("prefix", data.Prefix.ValueString())
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/storage/" + url.PathEscape(data.Name.ValueString()) + "?" + queryParams.Encode()
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
