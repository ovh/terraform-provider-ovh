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

var _ datasource.DataSourceWithConfigure = (*cloudPrivateEndpointServiceAvailabilitiesDataSource)(nil)

func NewCloudPrivateEndpointServiceAvailabilitiesDataSource() datasource.DataSource {
	return &cloudPrivateEndpointServiceAvailabilitiesDataSource{}
}

type cloudPrivateEndpointServiceAvailabilitiesDataSource struct {
	config *Config
}

func (d *cloudPrivateEndpointServiceAvailabilitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_private_endpoint_service_availabilities"
}

func (d *cloudPrivateEndpointServiceAvailabilitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudPrivateEndpointServiceAvailabilitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List services reachable through a private endpoint in a public cloud project.",
		MarkdownDescription: "List services reachable through a private endpoint in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Optional region filter. When set, only tiers available in that region are returned.",
				MarkdownDescription: "Optional region filter. When set, only tiers available in that region are returned.",
			},
			"service_availabilities": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of discoverable service availabilities",
				MarkdownDescription: "List of discoverable service availabilities",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Service name",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Service description",
						},
						"tiers": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Enabled tiers of the service",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Tier name",
									},
									"region": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Region where this tier is available",
									},
									"transport": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Transport protocol of the service",
									},
									"port": schema.Int64Attribute{
										Computed:    true,
										Description: "Port the service listens on through the endpoint",
									},
									"bandwidth": schema.Int64Attribute{
										Computed:    true,
										Description: "Bandwidth of the tier, in Mbit/s. Null when the catalog does not declare it",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type cloudPrivateEndpointServiceAvailabilitiesDataSourceModel struct {
	ServiceName           ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region                ovhtypes.TfStringValue `tfsdk:"region"`
	ServiceAvailabilities types.List             `tfsdk:"service_availabilities"`
}

func privateEndpointServiceAvailabilityTierAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      ovhtypes.TfStringType{},
		"region":    ovhtypes.TfStringType{},
		"transport": ovhtypes.TfStringType{},
		"port":      types.Int64Type,
		"bandwidth": types.Int64Type,
	}
}

func privateEndpointServiceAvailabilityAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"tiers": types.ListType{
			ElemType: types.ObjectType{AttrTypes: privateEndpointServiceAvailabilityTierAttrTypes()},
		},
	}
}

func (d *cloudPrivateEndpointServiceAvailabilitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudPrivateEndpointServiceAvailabilitiesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/reference/privateEndpoint/serviceAvailability"
	if !data.Region.IsNull() && !data.Region.IsUnknown() && data.Region.ValueString() != "" {
		endpoint += "?region=" + url.QueryEscape(data.Region.ValueString())
	}

	var apiAvailabilities []CloudPrivateEndpointServiceAvailability
	if err := d.config.OVHClient.Get(endpoint, &apiAvailabilities); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	tierObjType := types.ObjectType{AttrTypes: privateEndpointServiceAvailabilityTierAttrTypes()}
	availabilityObjs := make([]attr.Value, 0, len(apiAvailabilities))
	for _, svc := range apiAvailabilities {
		tierObjs := make([]attr.Value, 0, len(svc.Tiers))
		for _, tier := range svc.Tiers {
			region := ""
			if tier.Location != nil {
				region = tier.Location.Region
			}
			tierObj, diags := types.ObjectValue(privateEndpointServiceAvailabilityTierAttrTypes(), map[string]attr.Value{
				"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(tier.Name)},
				"region":    ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
				"transport": ovhtypes.TfStringValue{StringValue: types.StringValue(tier.Transport)},
				"port":      types.Int64Value(tier.Port),
				"bandwidth": types.Int64PointerValue(tier.Bandwidth),
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			tierObjs = append(tierObjs, tierObj)
		}

		tiersList, diags := types.ListValue(tierObjType, tierObjs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		svcObj, diags := types.ObjectValue(privateEndpointServiceAvailabilityAttrTypes(), map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(svc.Name)},
			"description": ovhtypes.TfStringValue{StringValue: types.StringValue(svc.Description)},
			"tiers":       tiersList,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		availabilityObjs = append(availabilityObjs, svcObj)
	}

	availabilitiesList, diags := types.ListValue(types.ObjectType{AttrTypes: privateEndpointServiceAvailabilityAttrTypes()}, availabilityObjs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ServiceAvailabilities = availabilitiesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
