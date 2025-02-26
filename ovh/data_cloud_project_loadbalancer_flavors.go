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

var _ datasource.DataSourceWithConfigure = (*cloudProjectLoadbalancerFlavorsDataSource)(nil)

func NewCloudProjectLoadbalancerFlavorsDataSource() datasource.DataSource {
	return &cloudProjectLoadbalancerFlavorsDataSource{}
}

type cloudProjectLoadbalancerFlavorsDataSource struct {
	config *Config
}

func (d *cloudProjectLoadbalancerFlavorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_loadbalancer_flavors"
}

func (d *cloudProjectLoadbalancerFlavorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectLoadbalancerFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectLoadbalancerFlavorsDataSourceSchema(ctx)
}

func (d *cloudProjectLoadbalancerFlavorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectLoadbalancerFlavorsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/region/" + url.PathEscape(data.RegionName.ValueString()) + "/loadbalancing/flavor"

	var arr []CloudProjectLoadbalancerFlavorsValue
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

	data.Flavors = ovhtypes.TfListNestedValue[CloudProjectLoadbalancerFlavorsValue]{
		ListValue: basetypes.NewListValueMust(CloudProjectLoadbalancerFlavorsValue{}.Type(ctx), b),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
