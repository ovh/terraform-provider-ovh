package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudInstanceFlavorsDataSource)(nil)

func NewCloudInstanceFlavorsDataSource() datasource.DataSource {
	return &cloudInstanceFlavorsDataSource{}
}

type cloudInstanceFlavorsDataSource struct {
	config *Config
}

func (d *cloudInstanceFlavorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_flavors"
}

func (d *cloudInstanceFlavorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudInstanceFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CloudInstanceFlavorsDataSourceSchema(ctx)
}

func (d *cloudInstanceFlavorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceFlavorsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/cloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/reference/instance/flavor?region=" + url.QueryEscape(data.RegionName.ValueString())

	var arr []CloudInstanceFlavorsValue
	if err := d.config.OVHClient.Get(endpoint, &arr); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Filter by name regexp if provided
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		namePattern := data.Name.ValueString()
		re, err := regexp.Compile(namePattern)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid name filter",
				fmt.Sprintf("The name filter %q is not a valid regular expression: %s", namePattern, err.Error()),
			)
			return
		}

		filtered := make([]CloudInstanceFlavorsValue, 0, len(arr))
		for _, f := range arr {
			if re.MatchString(f.Name.ValueString()) {
				filtered = append(filtered, f)
			}
		}
		arr = filtered
	}

	var b []attr.Value
	for _, a := range arr {
		b = append(b, a)
	}

	data.Flavors = ovhtypes.TfListNestedValue[CloudInstanceFlavorsValue]{
		ListValue: basetypes.NewListValueMust(CloudInstanceFlavorsValue{}.Type(ctx), b),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
