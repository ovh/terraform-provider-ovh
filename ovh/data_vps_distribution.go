package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ datasource.DataSource              = (*vpsDistributionDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*vpsDistributionDataSource)(nil)
)

func NewVpsDistributionDataSource() datasource.DataSource {
	return &vpsDistributionDataSource{}
}

type vpsDistributionDataSource struct {
	config *Config
}

type vpsDistributionModel struct {
	ServiceName       types.String `tfsdk:"service_name"`
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Distribution      types.String `tfsdk:"distribution"`
	BitFormat         types.Int64  `tfsdk:"bit_format"`
	Locale            types.String `tfsdk:"locale"`
	AvailableLanguage types.List   `tfsdk:"available_language"`
}

func (d *vpsDistributionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_distribution"
}

func (d *vpsDistributionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpsDistributionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns the distribution currently installed on a VPS.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required:    true,
				Description: "The internal name of the VPS (e.g. vps-123456.ovh.net).",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Template id of the installed distribution.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Template name (e.g. Debian 12).",
			},
			"distribution": schema.StringAttribute{
				Computed:    true,
				Description: "Distribution family (e.g. Debian).",
			},
			"bit_format": schema.Int64Attribute{
				Computed:    true,
				Description: "Architecture bit format (32 or 64).",
			},
			"locale": schema.StringAttribute{
				Computed:    true,
				Description: "Default locale of the installed distribution.",
			},
			"available_language": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Languages available for the installed distribution.",
			},
		},
	}
}

func (d *vpsDistributionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vpsDistributionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/vps/%s/distribution", url.PathEscape(data.ServiceName.ValueString()))

	var tpl VPSDistributionTemplate
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &tpl); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				resp.Diagnostics.AddError(
					"VPS API endpoint not available",
					fmt.Sprintf(
						"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
							"This data source may only work on legacy VPS plans, or the endpoint "+
							"may have been deprecated. See the data source's documentation for "+
							"supported VPS generations.",
						endpoint),
				)
				return
			case strings.Contains(msg, "does not exist"):
				resp.Diagnostics.AddError(
					"VPS resource not found",
					fmt.Sprintf(
						"the requested resource at %s does not exist (the VPS may not have "+
							"the required option subscribed, or the resource ID is wrong)",
						endpoint),
				)
				return
			}
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", endpoint),
			err.Error(),
		)
		return
	}

	bit := int64(0)
	if tpl.BitFormat != "" {
		if n, err := strconv.Atoi(tpl.BitFormat); err == nil {
			bit = int64(n)
		}
	}

	langs := make([]attr.Value, 0, len(tpl.AvailableLanguage))
	for _, l := range tpl.AvailableLanguage {
		langs = append(langs, types.StringValue(l))
	}
	langList, diags := types.ListValue(types.StringType, langs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(int64(tpl.ID))
	data.Name = types.StringValue(tpl.Name)
	data.Distribution = types.StringValue(tpl.Distribution)
	data.BitFormat = types.Int64Value(bit)
	data.Locale = types.StringValue(tpl.Locale)
	data.AvailableLanguage = langList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
