package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectNetworkPrivateSubnetsDataSource)(nil)

func NewCloudProjectNetworkPrivateSubnetsDataSource() datasource.DataSource {
	return &cloudProjectNetworkPrivateSubnetsDataSource{}
}

type cloudProjectNetworkPrivateSubnetsDataSource struct {
	config *Config
}

func (d *cloudProjectNetworkPrivateSubnetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_network_private_subnets"
}

func (d *cloudProjectNetworkPrivateSubnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectNetworkPrivateSubnetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectNetworkPrivateSubnetsDataSourceSchema(ctx)
}

func (d *cloudProjectNetworkPrivateSubnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectNetworkPrivateSubnetsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/private/" + url.PathEscape(data.NetworkId.ValueString()) + "/subnet"

	var arr []CloudProjectNetworkPrivateSubnetsValue

	if err := d.config.OVHClient.Get(endpoint, &arr); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	var b []attr.Value
	for _, a := range arr {
		b = append(b, a)
	}

	data.Subnets = ovhtypes.TfListNestedValue[CloudProjectNetworkPrivateSubnetsValue]{
		ListValue: basetypes.NewListValueMust(CloudProjectNetworkPrivateSubnetsValue{}.Type(ctx), b),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
