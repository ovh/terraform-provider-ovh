package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ datasource.DataSource              = (*vpsDistributionSoftwareItemDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*vpsDistributionSoftwareItemDataSource)(nil)
)

func NewVpsDistributionSoftwareItemDataSource() datasource.DataSource {
	return &vpsDistributionSoftwareItemDataSource{}
}

type vpsDistributionSoftwareItemDataSource struct {
	config *Config
}

type vpsDistributionSoftwareItemModel struct {
	ServiceName types.String `tfsdk:"service_name"`
	SoftwareID  types.Int64  `tfsdk:"software_id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Status      types.String `tfsdk:"status"`
}

func (d *vpsDistributionSoftwareItemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_distribution_software_item"
}

func (d *vpsDistributionSoftwareItemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpsDistributionSoftwareItemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns details for a single piece of software installed on a VPS.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required:    true,
				Description: "The internal name of the VPS.",
			},
			"software_id": schema.Int64Attribute{
				Required:    true,
				Description: "Software id as exposed by /vps/{serviceName}/distribution/software.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Software name.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Software type (database, environment, webserver).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Software status (stable, testing, deprecated).",
			},
		},
	}
}

func (d *vpsDistributionSoftwareItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vpsDistributionSoftwareItemModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf(
		"/vps/%s/distribution/software/%d",
		url.PathEscape(data.ServiceName.ValueString()),
		data.SoftwareID.ValueInt64(),
	)

	var sw VPSDistributionSoftware
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &sw); err != nil {
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

	data.Name = types.StringValue(sw.Name)
	data.Type = types.StringValue(sw.Type)
	data.Status = types.StringValue(sw.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
