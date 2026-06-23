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

var _ datasource.DataSourceWithConfigure = (*cloudSecurityGroupsDataSource)(nil)

func NewCloudSecurityGroupsDataSource() datasource.DataSource {
	return &cloudSecurityGroupsDataSource{}
}

type cloudSecurityGroupsDataSource struct {
	config *Config
}

// CloudSecurityGroupsModel is the model for the plural security groups data source.
type CloudSecurityGroupsModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecurityGroups types.List             `tfsdk:"security_groups"`
}

func (d *cloudSecurityGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_security_groups"
}

func (d *cloudSecurityGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// SecurityGroupListItemAttrTypes returns the attribute types for a single
// security group element of the plural data source list.
func SecurityGroupListItemAttrTypes() map[string]attr.Type {
	ruleObjType := types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()}
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"region":      ovhtypes.TfStringType{},
		"rule": types.ListType{
			ElemType: ruleObjType,
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: SecurityGroupCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudSecurityGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the security groups of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"security_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of security groups",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Security group ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the security group",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Description of the security group",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region of the security group",
						},
						"rule": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of security group rules",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"direction": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Direction of the rule",
									},
									"ethernet_type": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Ethernet type",
									},
									"protocol": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Protocol",
									},
									"port_range_min": schema.Int64Attribute{
										Computed:    true,
										Description: "Minimum port number",
									},
									"port_range_max": schema.Int64Attribute{
										Computed:    true,
										Description: "Maximum port number",
									},
									"remote_group_id": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Remote security group ID",
									},
									"remote_ip_prefix": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Remote IP prefix",
									},
									"description": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Description of the rule",
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
							Description: "Creation date of the security group",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the security group",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Security group readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": securityGroupDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudSecurityGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudSecurityGroupsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup"

	var responseData []CloudSecurityGroupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: SecurityGroupListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildSecurityGroupListItemObject(ctx, &responseData[i]))
	}

	data.SecurityGroups = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildSecurityGroupListItemObject builds a single security group list element
// from an API response, reusing the resource helpers for nested objects.
func buildSecurityGroupListItemObject(ctx context.Context, response *CloudSecurityGroupAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	ruleVal := types.ListNull(types.ObjectType{AttrTypes: SecurityGroupRuleAttrTypes()})

	if response.TargetSpec != nil {
		nameVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}

		if response.TargetSpec.Description != "" {
			descVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		if response.TargetSpec.Location != nil {
			regionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}

		ruleVal = buildSecurityGroupTargetRulesList(response.TargetSpec.Rules)
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildSecurityGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(SecurityGroupCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		SecurityGroupListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":            nameVal,
			"description":     descVal,
			"region":          regionVal,
			"rule":            ruleVal,
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
		},
	)

	return obj
}
