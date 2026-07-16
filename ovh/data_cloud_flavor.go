package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudFlavorDataSource)(nil)

func NewCloudFlavorDataSource() datasource.DataSource {
	return &cloudFlavorDataSource{}
}

type cloudFlavorDataSource struct {
	config *Config
}

func (d *cloudFlavorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_flavor"
}

func (d *cloudFlavorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudFlavorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a flavor available in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
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
	}
}

func (d *cloudFlavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudFlavorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/instance/flavor/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudFlavorAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	sn := data.ServiceName
	id := data.Id
	data = buildFlavorModel(&responseData)
	data.ServiceName = sn
	data.Id = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
