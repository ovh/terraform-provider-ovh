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

var _ datasource.DataSourceWithConfigure = (*cloudInstanceGroupsDataSource)(nil)

func NewCloudInstanceGroupsDataSource() datasource.DataSource {
	return &cloudInstanceGroupsDataSource{}
}

type cloudInstanceGroupsDataSource struct {
	config *Config
}

// CloudInstanceGroupsModel is the model for the plural instance groups data source.
type CloudInstanceGroupsModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	InstanceGroups types.List             `tfsdk:"instance_groups"`
}

func (d *cloudInstanceGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_groups"
}

func (d *cloudInstanceGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// InstanceGroupListItemAttrTypes returns the attr types for one element of `instance_groups`.
func InstanceGroupListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"region":          ovhtypes.TfStringType{},
		"policy":          ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: InstanceGroupCurrentStateAttrTypes()},
	}
}

func (d *cloudInstanceGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the instance groups (placement groups) of a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Service name / cloud project id"},
			"instance_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of instance groups",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"name":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"region":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"policy":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"checksum":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"created_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"updated_at":      schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"resource_status": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"current_state":   schema.SingleNestedAttribute{Computed: true, Attributes: instanceGroupCurrentStateDataSourceSchemaAttributes()},
				}},
			},
		},
	}
}

func (d *cloudInstanceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstanceGroupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instanceGroup"

	var responseData []CloudInstanceGroupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	itemObjType := types.ObjectType{AttrTypes: InstanceGroupListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildInstanceGroupListItemObject(ctx, &responseData[i]))
	}
	data.InstanceGroups = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildInstanceGroupListItemObject builds one `instance_groups` element from an API response.
func buildInstanceGroupListItemObject(ctx context.Context, response *CloudInstanceGroupAPIResponse) basetypes.ObjectValue {
	sv := func(v string) ovhtypes.TfStringValue {
		return ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
	}

	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	policyVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}

	if response.TargetSpec != nil {
		ts := response.TargetSpec
		nameVal = sv(ts.Name)
		policyVal = sv(ts.Policy)
		if ts.Location != nil {
			regionVal = sv(ts.Location.Region)
		}
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildInstanceGroupCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(InstanceGroupCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(InstanceGroupListItemAttrTypes(), map[string]attr.Value{
		"id":              sv(response.Id),
		"name":            nameVal,
		"region":          regionVal,
		"policy":          policyVal,
		"checksum":        sv(response.Checksum),
		"created_at":      sv(response.CreatedAt),
		"updated_at":      sv(response.UpdatedAt),
		"resource_status": sv(response.ResourceStatus),
		"current_state":   currentStateVal,
	})
	return obj
}
