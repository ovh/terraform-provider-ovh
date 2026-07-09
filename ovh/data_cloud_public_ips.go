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

var _ datasource.DataSourceWithConfigure = (*cloudPublicIPsDataSource)(nil)

func NewCloudPublicIPsDataSource() datasource.DataSource {
	return &cloudPublicIPsDataSource{}
}

type cloudPublicIPsDataSource struct {
	config *Config
}

func (d *cloudPublicIPsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_public_ips"
}

func (d *cloudPublicIPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudPublicIPsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list all public IPs (additional, external network and floating IPs) of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"public_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of public IPs of the project",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Public IP address",
						},
						"type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Type of the public IP (ADDITIONAL_IP, EXT_NET_IP, FLOATING_IP)",
						},
					},
				},
			},
		},
	}
}

// cloudPublicIPsDataSourceModel is the Terraform state model for this data source.
type cloudPublicIPsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	PublicIPs   types.List             `tfsdk:"public_ips"`
}

// cloudPublicIPSummaryAttrTypes returns the attribute types for a single public IP
// item of the list.
func cloudPublicIPSummaryAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":   ovhtypes.TfStringType{},
		"type": ovhtypes.TfStringType{},
	}
}

func (d *cloudPublicIPsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudPublicIPsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/publicIp"

	apiPublicIPs, err := helpers.GetAllPagesV2[CloudPublicIPSummary](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	publicIPObjs := make([]attr.Value, 0, len(apiPublicIPs))
	for _, publicIP := range apiPublicIPs {
		itemObj, diags := types.ObjectValue(
			cloudPublicIPSummaryAttrTypes(),
			map[string]attr.Value{
				"ip":   ovhtypes.TfStringValue{StringValue: types.StringValue(publicIP.IP)},
				"type": ovhtypes.TfStringValue{StringValue: types.StringValue(publicIP.Type)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		publicIPObjs = append(publicIPObjs, itemObj)
	}

	publicIPsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: cloudPublicIPSummaryAttrTypes()},
		publicIPObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PublicIPs = publicIPsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
