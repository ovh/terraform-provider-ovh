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

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancersDataSource)(nil)

func NewCloudLoadbalancersDataSource() datasource.DataSource {
	return &cloudLoadbalancersDataSource{}
}

type cloudLoadbalancersDataSource struct {
	config *Config
}

// CloudLoadbalancersModel is the model for the plural loadbalancers data source.
type CloudLoadbalancersModel struct {
	ServiceName   ovhtypes.TfStringValue `tfsdk:"service_name"`
	Loadbalancers types.List             `tfsdk:"loadbalancers"`
}

func (d *cloudLoadbalancersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancers"
}

func (d *cloudLoadbalancersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// LoadbalancerListItemAttrTypes returns the attribute types for a single
// loadbalancer element of the plural data source list.
func LoadbalancerListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                ovhtypes.TfStringType{},
		"name":              ovhtypes.TfStringType{},
		"description":       ovhtypes.TfStringType{},
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
		"network": types.ObjectType{
			AttrTypes: LoadbalancerNetworkAttrTypes(),
		},
		"flavor_name":     ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: LoadbalancerCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudLoadbalancersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the loadbalancers of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"loadbalancers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of loadbalancers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Loadbalancer ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Loadbalancer name",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Loadbalancer description",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region where the loadbalancer is located",
						},
						"availability_zone": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Availability zone for the loadbalancer",
						},
						"network": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Network of the VIP",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "ID of the network for the VIP",
								},
								"subnet_id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "ID of the subnet for the VIP",
								},
								"ip": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "IP requested for the VIP",
								},
							},
						},
						"flavor_name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the loadbalancer flavor",
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the loadbalancer",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the loadbalancer",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Loadbalancer readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": loadbalancerDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudLoadbalancersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/loadbalancer"

	var responseData []CloudLoadbalancerAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: LoadbalancerListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildLoadbalancerListItemObject(ctx, &responseData[i]))
	}

	data.Loadbalancers = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildLoadbalancerListItemObject builds a single loadbalancer list element from
// an API response, reusing the resource helpers for nested objects.
func buildLoadbalancerListItemObject(ctx context.Context, response *CloudLoadbalancerAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	azVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	flavorNameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	networkVal := types.ObjectNull(LoadbalancerNetworkAttrTypes())

	if response.TargetSpec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}

		if response.TargetSpec.Description != "" {
			descVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		if response.TargetSpec.Location != nil {
			regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				azVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
		}

		if response.TargetSpec.Network != nil {
			ipVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
			if response.TargetSpec.Network.IP != "" {
				ipVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.IP)}
			}

			networkVal, _ = types.ObjectValue(
				LoadbalancerNetworkAttrTypes(),
				map[string]attr.Value{
					"id":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.ID)},
					"subnet_id": ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Network.SubnetID)},
					"ip":        ipVal,
				},
			)
		}

		if response.TargetSpec.Flavor != nil {
			flavorNameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Flavor.Name)}
		}
	}

	currentStateVal := types.ObjectNull(LoadbalancerCurrentStateAttrTypes())
	if response.CurrentState != nil {
		currentStateVal = buildLoadbalancerCurrentStateObject(ctx, response.CurrentState)
	}

	obj, _ := types.ObjectValue(
		LoadbalancerListItemAttrTypes(),
		map[string]attr.Value{
			"id":                ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":              nameVal,
			"description":       descVal,
			"region":            regionVal,
			"availability_zone": azVal,
			"network":           networkVal,
			"flavor_name":       flavorNameVal,
			"checksum":          ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status":   ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":     currentStateVal,
		},
	)

	return obj
}
