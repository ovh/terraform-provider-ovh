package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudGatewaysDataSource)(nil)

func NewCloudGatewaysDataSource() datasource.DataSource {
	return &cloudGatewaysDataSource{}
}

type cloudGatewaysDataSource struct {
	config *Config
}

// CloudGatewaysModel is the model for the plural gateways data source.
type CloudGatewaysModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Gateways    types.List             `tfsdk:"gateways"`
}

func (d *cloudGatewaysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_gateways"
}

func (d *cloudGatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// GatewayListItemAttrTypes returns the attribute types for a single gateway
// element of the plural data source list.
func GatewayListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"region":      ovhtypes.TfStringType{},
		"external_gateway": types.ObjectType{
			AttrTypes: ExternalGatewayAttrTypes(),
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: GatewayCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudGatewaysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the gateways of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"gateways": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of gateways",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Gateway ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Gateway name",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Gateway description",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the gateway",
						},
						"external_gateway": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "External gateway configuration",
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the external gateway is enabled",
								},
								"model": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "External gateway sizing model",
								},
							},
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the gateway",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the gateway",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Gateway readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the gateway",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Gateway name",
								},
								"description": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Gateway description",
								},
								"status": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "OpenStack router status (ACTIVE, BUILD, DOWN, ERROR)",
								},
								"external_gateway": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "External gateway configuration",
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Computed:    true,
											Description: "Whether the external gateway is enabled",
										},
										"model": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "External gateway sizing model",
										},
									},
								},
								"external_ip": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "External IP address assigned to the gateway",
								},
								"subnets": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Currently attached subnets",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "Subnet ID",
											},
										},
									},
								},
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region",
								},
								"availability_zone": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Availability zone",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudGatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudGatewaysModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway"

	var responseData []CloudGatewayAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: GatewayListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildGatewayListItemObject(ctx, &responseData[i]))
	}

	data.Gateways = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildGatewayListItemObject builds a single gateway list element from an API
// response, reusing the resource helpers for nested objects.
func buildGatewayListItemObject(ctx context.Context, response *CloudGatewayAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	externalGatewayVal := types.ObjectNull(ExternalGatewayAttrTypes())

	if response.TargetSpec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		descVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}

		if response.TargetSpec.Location != nil {
			regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}

		if response.TargetSpec.ExternalGateway != nil {
			modelVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
			if response.TargetSpec.ExternalGateway.Model != "" {
				modelVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.ExternalGateway.Model)}
			}
			externalGatewayVal, _ = types.ObjectValue(
				ExternalGatewayAttrTypes(),
				map[string]attr.Value{
					"enabled": types.BoolValue(response.TargetSpec.ExternalGateway.Enabled),
					"model":   modelVal,
				},
			)
		}
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildGatewayCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(GatewayCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		GatewayListItemAttrTypes(),
		map[string]attr.Value{
			"id":               ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":             nameVal,
			"description":      descVal,
			"region":           regionVal,
			"external_gateway": externalGatewayVal,
			"checksum":         ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":       ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":       ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status":  ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":    currentStateVal,
		},
	)

	return obj
}
