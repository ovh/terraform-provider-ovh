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

var _ datasource.DataSourceWithConfigure = (*cloudInstancesDataSource)(nil)

func NewCloudInstancesDataSource() datasource.DataSource {
	return &cloudInstancesDataSource{}
}

type cloudInstancesDataSource struct {
	config *Config
}

// CloudInstancesModel is the model for the plural instances data source.
type CloudInstancesModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Instances   types.List             `tfsdk:"instances"`
}

func (d *cloudInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instances"
}

func (d *cloudInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// InstanceListItemAttrTypes returns the attr types for one element of `instances`.
func InstanceListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                ovhtypes.TfStringType{},
		"name":              ovhtypes.TfStringType{},
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
		"flavor_id":         ovhtypes.TfStringType{},
		"image_id":          ovhtypes.TfStringType{},
		"power_state":       ovhtypes.TfStringType{},
		"checksum":          ovhtypes.TfStringType{},
		"created_at":        ovhtypes.TfStringType{},
		"updated_at":        ovhtypes.TfStringType{},
		"resource_status":   ovhtypes.TfStringType{},
		"current_state":     types.ObjectType{AttrTypes: InstanceCurrentStateAttrTypes()},
	}
}

func (d *cloudInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the instances of a public cloud project. Note: current_state.shares is not populated by the list endpoint; use ovh_cloud_instance for share visibility.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: "Service name / cloud project id"},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of instances",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":                schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"name":              schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"region":            schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"availability_zone": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"flavor_id":         schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"image_id":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"power_state":       schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"checksum":          schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"created_at":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"updated_at":        schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"resource_status":   schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true},
					"current_state":     schema.SingleNestedAttribute{Computed: true, Attributes: instanceCurrentStateDataSourceSchemaAttributes()},
				}},
			},
		},
	}
}

func (d *cloudInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudInstancesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance"

	var responseData []CloudInstanceAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	itemObjType := types.ObjectType{AttrTypes: InstanceListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildInstanceListItemObject(ctx, &responseData[i]))
	}
	data.Instances = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildInstanceListItemObject builds one `instances` element from an API response.
func buildInstanceListItemObject(ctx context.Context, response *CloudInstanceAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	azVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	flavorVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	imageVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	powerVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}

	if response.TargetSpec != nil {
		ts := response.TargetSpec
		nameVal = str(ts.Name)
		if ts.Location != nil {
			regionVal = str(ts.Location.Region)
			if ts.Location.AvailabilityZone != "" {
				azVal = str(ts.Location.AvailabilityZone)
			}
		}
		if ts.Flavor != nil {
			flavorVal = str(ts.Flavor.Id)
		}
		if ts.Image != nil && ts.Image.Id != "" {
			imageVal = str(ts.Image.Id)
		}
		if ts.PowerState != "" {
			powerVal = str(ts.PowerState)
		}
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildInstanceCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(InstanceCurrentStateAttrTypes())
	}

	// Fall back to the observed location for region/availability_zone when the
	// requested spec doesn't carry them (AZ is often platform-assigned).
	if response.CurrentState != nil && response.CurrentState.Location != nil {
		loc := response.CurrentState.Location
		if regionVal.IsNull() || regionVal.ValueString() == "" {
			regionVal = str(loc.Region)
		}
		if (azVal.IsNull() || azVal.ValueString() == "") && loc.AvailabilityZone != "" {
			azVal = str(loc.AvailabilityZone)
		}
	}

	obj, _ := types.ObjectValue(InstanceListItemAttrTypes(), map[string]attr.Value{
		"id":                str(response.Id),
		"name":              nameVal,
		"region":            regionVal,
		"availability_zone": azVal,
		"flavor_id":         flavorVal,
		"image_id":          imageVal,
		"power_state":       powerVal,
		"checksum":          str(response.Checksum),
		"created_at":        str(response.CreatedAt),
		"updated_at":        str(response.UpdatedAt),
		"resource_status":   str(response.ResourceStatus),
		"current_state":     currentStateVal,
	})
	return obj
}
