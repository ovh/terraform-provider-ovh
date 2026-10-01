package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSourceWithConfigure = (*dedicatedCloudUserObjectRightDataSource)(nil)

func NewDedicatedCloudUserObjectRightDataSource() datasource.DataSource {
	return &dedicatedCloudUserObjectRightDataSource{}
}

type dedicatedCloudUserObjectRightDataSource struct {
	config *Config
}

func (d *dedicatedCloudUserObjectRightDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_cloud_user_object_right"
}

func (d *dedicatedCloudUserObjectRightDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dedicatedCloudUserObjectRightDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DedicatedCloudUserObjectRightDataSourceSchema(ctx)
}

func (d *dedicatedCloudUserObjectRightDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DedicatedCloudUserObjectRightDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := dedicatedCloudUserObjectRightEndpoint(data.ServiceName.ValueString(), data.UserId.ValueInt64(), data.ObjectRightId.ValueInt64())

	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
