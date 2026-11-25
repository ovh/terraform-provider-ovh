package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = (*vpssDataSource)(nil)

func NewVPSsDataSource() datasource.DataSource {
	return &vpssDataSource{}
}

type vpssDataSource struct {
	config *Config
}

func (d *vpssDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpss"
}

func (d *vpssDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpssDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the VPS of your account.",
		Attributes: map[string]schema.Attribute{
			"vps": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "VPS of your account",
				MarkdownDescription: "VPS of your account",
				NestedObject:        schema.NestedAttributeObject{Attributes: vpsDataSourceAttributes(false)},
			},
		},
	}
}

func (d *vpssDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	servers, diags := newVPSBackend(d.config).List(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := newVPSsDataModel(ctx, servers)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
