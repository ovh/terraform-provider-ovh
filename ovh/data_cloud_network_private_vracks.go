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

var _ datasource.DataSourceWithConfigure = (*cloudNetworkPrivateVracksDataSource)(nil)

func NewCloudNetworkPrivateVracksDataSource() datasource.DataSource {
	return &cloudNetworkPrivateVracksDataSource{}
}

type cloudNetworkPrivateVracksDataSource struct {
	config *Config
}

// CloudNetworkPrivateVracksModel is the model for the plural networks data source.
type CloudNetworkPrivateVracksModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Networks    types.List             `tfsdk:"networks"`
}

func (d *cloudNetworkPrivateVracksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_network_private_vracks"
}

func (d *cloudNetworkPrivateVracksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// NetworkListItemAttrTypes returns the attribute types for a single network
// element of the plural data source list.
func NetworkListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"description":     ovhtypes.TfStringType{},
		"region":          ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: NetworkCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudNetworkPrivateVracksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the private networks (vRack) of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"networks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of private networks",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Network ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Network name",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Network description",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the network",
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the network",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the network",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Network readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the network",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Network name",
								},
								"description": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Network description",
								},
								"location": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Location details",
									Attributes: map[string]schema.Attribute{
										"region": schema.StringAttribute{
											CustomType:  ovhtypes.TfStringType{},
											Computed:    true,
											Description: "Region code",
										},
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

func (d *cloudNetworkPrivateVracksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudNetworkPrivateVracksModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network"

	var responseData []CloudNetworkAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: NetworkListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildNetworkListItemObject(&responseData[i]))
	}

	data.Networks = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildNetworkListItemObject builds a single network list element from an API
// response, reusing the resource helpers for nested objects.
func buildNetworkListItemObject(response *CloudNetworkAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}

	if response.TargetSpec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		descVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		if response.TargetSpec.Location != nil {
			regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildNetworkCurrentStateObject(response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(NetworkCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		NetworkListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":            nameVal,
			"description":     descVal,
			"region":          regionVal,
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
		},
	)

	return obj
}
