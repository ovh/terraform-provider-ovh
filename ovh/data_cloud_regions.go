package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudRegionsDataSource)(nil)

func NewCloudRegionsDataSource() datasource.DataSource {
	return &cloudRegionsDataSource{}
}

type cloudRegionsDataSource struct {
	config *Config
}

func (d *cloudRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_regions"
}

func (d *cloudRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	regionAttributes := cloudRegionDetailAttributes()
	regionAttributes["name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Computed:            true,
		Description:         "Name of the region (e.g. GRA11)",
		MarkdownDescription: "Name of the region (e.g. GRA11)",
	}

	resp.Schema = schema.Schema{
		Description:         "List the regions available for a public cloud project.",
		MarkdownDescription: "List the regions available for a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"regions": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of regions available for the cloud project",
				MarkdownDescription: "List of regions available for the cloud project",
				NestedObject: schema.NestedAttributeObject{
					Attributes: regionAttributes,
				},
			},
		},
	}
}

func (d *cloudRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudRegionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/region"

	regions, err := helpers.GetAllPagesV2[CloudRegion](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	regionObjs := make([]attr.Value, 0, len(regions))
	for _, region := range regions {
		obj, diags := cloudRegionToObject(ctx, region)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionObjs = append(regionObjs, obj)
	}

	regionList, diags := types.ListValue(
		types.ObjectType{AttrTypes: cloudRegionObjectAttrTypes()},
		regionObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Regions = regionList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
