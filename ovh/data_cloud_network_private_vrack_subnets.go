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

var _ datasource.DataSourceWithConfigure = (*cloudNetworkPrivateVrackSubnetsDataSource)(nil)

func NewCloudNetworkPrivateVrackSubnetsDataSource() datasource.DataSource {
	return &cloudNetworkPrivateVrackSubnetsDataSource{}
}

type cloudNetworkPrivateVrackSubnetsDataSource struct {
	config *Config
}

// CloudNetworkPrivateVrackSubnetsModel is the model for the plural subnets data source.
type CloudNetworkPrivateVrackSubnetsModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	NetworkId   ovhtypes.TfStringValue `tfsdk:"network_id"`
	Subnets     types.List             `tfsdk:"subnets"`
}

func (d *cloudNetworkPrivateVrackSubnetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_network_private_vrack_subnets"
}

func (d *cloudNetworkPrivateVrackSubnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// SubnetListItemAttrTypes returns the attribute types for a single subnet
// element of the plural data source list.
func SubnetListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           ovhtypes.TfStringType{},
		"name":         ovhtypes.TfStringType{},
		"cidr":         ovhtypes.TfStringType{},
		"region":       ovhtypes.TfStringType{},
		"description":  ovhtypes.TfStringType{},
		"dhcp_enabled": types.BoolType,
		"dns_nameservers": types.ListType{
			ElemType: types.StringType,
		},
		"gateway_ip": ovhtypes.TfStringType{},
		"allocation_pools": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: subnetAllocationPoolAttrTypes(),
			},
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: SubnetCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudNetworkPrivateVrackSubnetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the subnets of a private network (vRack) in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"network_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Network ID of the parent private network",
			},
			"subnets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of subnets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Subnet ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Subnet name",
						},
						"cidr": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "CIDR address range for the subnet",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the subnet",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Subnet description",
						},
						"dhcp_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether DHCP is enabled on the subnet",
						},
						"dns_nameservers": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "DNS nameservers for the subnet",
						},
						"gateway_ip": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Default gateway IP address",
						},
						"allocation_pools": schema.ListNestedAttribute{
							Computed:    true,
							Description: "IP address allocation pools",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"start": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Start IP address of the pool",
									},
									"end": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "End IP address of the pool",
									},
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
							Description: "Creation date of the subnet",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the subnet",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Subnet readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Current state of the subnet",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Subnet name",
								},
								"cidr": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "CIDR address range",
								},
								"description": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Subnet description",
								},
								"dhcp_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether DHCP is enabled",
								},
								"dns_nameservers": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "DNS nameservers",
								},
								"gateway_ip": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Default gateway IP address",
								},
								"host_routes": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Static host routes",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"destination": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "Destination CIDR",
											},
											"next_hop": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "Next hop IP address",
											},
										},
									},
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

func (d *cloudNetworkPrivateVrackSubnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudNetworkPrivateVrackSubnetsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/" + url.PathEscape(data.NetworkId.ValueString()) + "/subnet"

	var responseData []CloudSubnetAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: SubnetListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildSubnetListItemObject(&responseData[i]))
	}

	data.Subnets = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildSubnetListItemObject builds a single subnet list element from an API
// response, reusing the resource helpers for nested objects.
func buildSubnetListItemObject(response *CloudSubnetAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	cidrVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	gatewayIPVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	dhcpVal := types.BoolNull()
	dnsVal := types.ListNull(types.StringType)

	allocationPoolObjType := types.ObjectType{AttrTypes: subnetAllocationPoolAttrTypes()}
	allocationPoolsVal := types.ListNull(allocationPoolObjType)

	if response.TargetSpec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		cidrVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.CIDR)}
		descVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		gatewayIPVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.GatewayIP)}

		if response.TargetSpec.Location != nil {
			regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}

		if response.TargetSpec.DHCPEnabled != nil {
			dhcpVal = types.BoolValue(*response.TargetSpec.DHCPEnabled)
		}

		if response.TargetSpec.DNSNameservers != nil {
			dnsVals := make([]attr.Value, len(response.TargetSpec.DNSNameservers))
			for i, dns := range response.TargetSpec.DNSNameservers {
				dnsVals[i] = types.StringValue(dns)
			}
			dnsVal = types.ListValueMust(types.StringType, dnsVals)
		}

		if response.TargetSpec.AllocationPools != nil {
			poolObjs := make([]attr.Value, len(response.TargetSpec.AllocationPools))
			for i, pool := range response.TargetSpec.AllocationPools {
				poolObj, _ := types.ObjectValue(
					subnetAllocationPoolAttrTypes(),
					map[string]attr.Value{
						"start": ovhtypes.TfStringValue{StringValue: types.StringValue(pool.Start)},
						"end":   ovhtypes.TfStringValue{StringValue: types.StringValue(pool.End)},
					},
				)
				poolObjs[i] = poolObj
			}
			allocationPoolsVal = types.ListValueMust(allocationPoolObjType, poolObjs)
		}
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildSubnetCurrentStateObject(response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(SubnetCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		SubnetListItemAttrTypes(),
		map[string]attr.Value{
			"id":               ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":             nameVal,
			"cidr":             cidrVal,
			"region":           regionVal,
			"description":      descVal,
			"dhcp_enabled":     dhcpVal,
			"dns_nameservers":  dnsVal,
			"gateway_ip":       gatewayIPVal,
			"allocation_pools": allocationPoolsVal,
			"checksum":         ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":       ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":       ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status":  ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":    currentStateVal,
		},
	)

	return obj
}
