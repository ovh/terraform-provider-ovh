package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectNetworkPrivateDataSource)(nil)

func NewCloudProjectNetworkPrivateDataSource() datasource.DataSource {
	return &cloudProjectNetworkPrivateDataSource{}
}

type cloudProjectNetworkPrivateDataSource struct {
	config *Config
}

func (d *cloudProjectNetworkPrivateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_network_private"
}

func (d *cloudProjectNetworkPrivateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectNetworkPrivateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectNetworkPrivateDataSourceSchema(ctx)
}

func (d *cloudProjectNetworkPrivateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectNetworkPrivateModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/private/" + url.PathEscape(data.NetworkId.ValueString())

	if err := d.config.OVHClient.Get(endpoint, &data); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Compute RegionsOpenstackIds map from Regions list
	internal_openstack_ids := map[string]attr.Value{}
	for _, region_attrs := range data.Regions.Elements() {
		internal_openstack_ids[region_attrs.(RegionsValue).Region.ValueString()] = region_attrs.(RegionsValue).OpenstackId
	}

	regions_openstack_ids, diags := ovhtypes.NewTfMapNestedValue[ovhtypes.TfStringValue](ctx, internal_openstack_ids)
	resp.Diagnostics.Append(diags...)

	data.RegionsOpenstackIds = regions_openstack_ids

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
