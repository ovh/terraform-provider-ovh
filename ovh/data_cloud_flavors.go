package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudFlavorsDataSource)(nil)

func NewCloudFlavorsDataSource() datasource.DataSource {
	return &cloudFlavorsDataSource{}
}

type cloudFlavorsDataSource struct {
	config *Config
}

// CloudFlavorsModel is the model for the plural flavors data source.
type CloudFlavorsModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Flavors     types.List             `tfsdk:"flavors"`
}

func (d *cloudFlavorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_flavors"
}

func (d *cloudFlavorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudFlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the flavors available in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Filter the flavors by region",
			},
			"flavors": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of flavors",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Flavor ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Flavor name",
						},
						"vcpus": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of virtual CPUs",
						},
						"ram": schema.Int64Attribute{
							Computed:    true,
							Description: "Amount of RAM (MB)",
						},
						"disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Root disk size (GB)",
						},
						"swap": schema.Int64Attribute{
							Computed:    true,
							Description: "Swap size (MB)",
						},
						"ephemeral": schema.Int64Attribute{
							Computed:    true,
							Description: "Ephemeral disk size (GB)",
						},
						"is_public": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the flavor is publicly available to the project",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Flavor description",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region where the flavor is offered",
						},
					},
				},
			},
		},
	}
}

func (d *cloudFlavorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudFlavorsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/flavor"
	if !data.Region.IsNull() && data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var responseData []CloudFlavorAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	itemObjType := types.ObjectType{AttrTypes: flavorAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		m := buildFlavorModel(&responseData[i])
		obj, diags := types.ObjectValue(flavorAttrTypes(), map[string]attr.Value{
			"id":          m.Id,
			"name":        m.Name,
			"vcpus":       m.Vcpus,
			"ram":         m.Ram,
			"disk":        m.Disk,
			"swap":        m.Swap,
			"ephemeral":   m.Ephemeral,
			"is_public":   m.IsPublic,
			"description": m.Description,
			"region":      m.Region,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, obj)
	}

	data.Flavors = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
