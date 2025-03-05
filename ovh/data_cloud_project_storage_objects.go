package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectStorageObjectsDataSource)(nil)

func NewCloudProjectStorageObjectsDataSource() datasource.DataSource {
	return &cloudProjectStorageObjectsDataSource{}
}

type cloudProjectStorageObjectsDataSource struct {
	config *Config
}

func (d *cloudProjectStorageObjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_storage_objects"
}

func (d *cloudProjectStorageObjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectStorageObjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectStorageObjectsDataSourceSchema(ctx)
}

func (d *cloudProjectStorageObjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectStorageObjectsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	queryParams := url.Values{}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		queryParams.Add("limit", strconv.FormatInt(data.Limit.ValueInt64(), 10))
	}
	if !data.KeyMarker.IsNull() && !data.KeyMarker.IsUnknown() {
		queryParams.Add("keyMarker", data.KeyMarker.ValueString())
	}
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		queryParams.Add("prefix", data.Prefix.ValueString())
	}
	if !data.VersionIdMarker.IsNull() && !data.VersionIdMarker.IsUnknown() {
		queryParams.Add("versionIdMarker", data.VersionIdMarker.ValueString())
	}
	if !data.WithVersions.IsNull() && !data.WithVersions.IsUnknown() {
		queryParams.Add("withVersions", strconv.FormatBool(data.WithVersions.ValueBool()))
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/region/" + url.PathEscape(data.RegionName.ValueString()) +
		"/storage/" + url.PathEscape(data.Name.ValueString()) + "/object?" +
		queryParams.Encode()

	if err := d.config.OVHClient.Get(endpoint, &data.Objects); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
