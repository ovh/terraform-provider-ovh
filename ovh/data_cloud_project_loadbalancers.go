package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudProjectLoadbalancersDataSource)(nil)

func NewCloudProjectLoadbalancersDataSource() datasource.DataSource {
	return &cloudProjectLoadbalancersDataSource{}
}

type cloudProjectLoadbalancersDataSource struct {
	config *Config
}

func (d *cloudProjectLoadbalancersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_project_loadbalancers"
}

func (d *cloudProjectLoadbalancersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudProjectLoadbalancersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudProjectLoadbalancersDataSourceSchema(ctx)
}

func (d *cloudProjectLoadbalancersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudProjectLoadbalancersModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() {
		data.ServiceName.StringValue = basetypes.NewStringValue(os.Getenv("OVH_CLOUD_PROJECT_SERVICE"))
	}

	// Read API call logic
	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/loadbalancing/loadbalancer",
		url.PathEscape(data.ServiceName.ValueString()),
		url.PathEscape(data.RegionName.ValueString()),
	)

	var arr []CloudProjectLoadbalancersValue

	if err := d.config.OVHClient.Get(endpoint, &arr); err != nil {
		resp.Diagnostics.AddError("Failed to list loadbalancers", fmt.Sprintf("error calling GET %s: %s", endpoint, err))
		return
	}

	var b []attr.Value
	for _, a := range arr {
		b = append(b, a)
	}

	data.CloudProjectLoadbalancers = ovhtypes.TfListNestedValue[CloudProjectLoadbalancersValue]{
		ListValue: basetypes.NewListValueMust(CloudProjectLoadbalancersValue{}.Type(ctx), b),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
