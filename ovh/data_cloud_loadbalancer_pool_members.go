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

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerPoolMembersDataSource)(nil)

func NewCloudLoadbalancerPoolMembersDataSource() datasource.DataSource {
	return &cloudLoadbalancerPoolMembersDataSource{}
}

type cloudLoadbalancerPoolMembersDataSource struct {
	config *Config
}

// CloudLoadbalancerPoolMembersModel is the model for the plural pool members data source.
type CloudLoadbalancerPoolMembersModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	PoolId         ovhtypes.TfStringValue `tfsdk:"pool_id"`
	Members        types.List             `tfsdk:"members"`
}

func (d *cloudLoadbalancerPoolMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_pool_members"
}

func (d *cloudLoadbalancerPoolMembersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// MemberListItemAttrTypes returns the attribute types for a single member
// element of the plural data source list.
func MemberListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"name":          ovhtypes.TfStringType{},
		"address":       ovhtypes.TfStringType{},
		"protocol_port": types.Int64Type,
		"subnet_id":     ovhtypes.TfStringType{},
		"weight":        types.Int64Type,
		"backup":        types.BoolType,
		"monitor": types.ObjectType{
			AttrTypes: memberMonitorAttrTypes(),
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: MemberCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudLoadbalancerPoolMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the members of a pool in a public cloud loadbalancer.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"loadbalancer_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the loadbalancer",
			},
			"pool_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the pool",
			},
			"members": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of members",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Member ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Member name",
						},
						"address": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "IP address of the member",
						},
						"protocol_port": schema.Int64Attribute{
							Computed:    true,
							Description: "Port used by the member to receive traffic",
						},
						"subnet_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "ID of the subnet the member is in",
						},
						"weight": schema.Int64Attribute{
							Computed:    true,
							Description: "Weight of the member in the pool (0-256). Higher weight receives more traffic.",
						},
						"backup": schema.BoolAttribute{
							Computed:    true,
							Description: "When true, the member is a backup member and only receives traffic when all non-backup members are down",
						},
						"monitor": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Optional health monitor address and port override for this member",
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "IP address used by the health monitor for this member",
								},
								"port": schema.Int64Attribute{
									Computed:    true,
									Description: "Port used by the health monitor for this member",
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
							Description: "Creation date of the member",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the member",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Member readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": loadbalancerPoolMemberDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudLoadbalancerPoolMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerPoolMembersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/pool/" + url.PathEscape(data.PoolId.ValueString()) +
		"/member"

	var responseData []CloudLoadbalancerPoolMemberAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: MemberListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildMemberListItemObject(&responseData[i]))
	}

	data.Members = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildMemberListItemObject builds a single member list element from an API
// response, reusing the resource helpers for nested objects.
func buildMemberListItemObject(response *CloudLoadbalancerPoolMemberAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	addressVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	protocolPortVal := types.Int64Null()
	subnetIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	weightVal := types.Int64Null()
	backupVal := types.BoolNull()
	monitorVal := types.ObjectNull(memberMonitorAttrTypes())

	if spec := response.TargetSpec; spec != nil {
		nameVal = lbStringOrNull(spec.Name)
		addressVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Address)}
		protocolPortVal = types.Int64Value(spec.ProtocolPort)

		if spec.Subnet != nil && spec.Subnet.Id != "" {
			subnetIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Subnet.Id)}
		}

		if spec.Weight != nil {
			weightVal = types.Int64Value(*spec.Weight)
		}

		if spec.Backup != nil {
			backupVal = types.BoolValue(*spec.Backup)
		}

		if spec.Monitor != nil {
			monitorVal = buildMemberMonitorObject(spec.Monitor)
		}
	}

	currentStateVal := types.ObjectNull(MemberCurrentStateAttrTypes())
	if response.CurrentState != nil {
		currentStateVal = buildMemberCurrentStateObject(response.CurrentState)
	}

	obj, _ := types.ObjectValue(
		MemberListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":            nameVal,
			"address":         addressVal,
			"protocol_port":   protocolPortVal,
			"subnet_id":       subnetIdVal,
			"weight":          weightVal,
			"backup":          backupVal,
			"monitor":         monitorVal,
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
		},
	)

	return obj
}
