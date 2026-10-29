package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageSwiftContainersDataSource)(nil)

func NewCloudStorageSwiftContainersDataSource() datasource.DataSource {
	return &cloudStorageSwiftContainersDataSource{}
}

type cloudStorageSwiftContainersDataSource struct {
	config *Config
}

func (d *cloudStorageSwiftContainersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_swift_containers"
}

func (d *cloudStorageSwiftContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageSwiftContainersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the Swift object storage containers in a public cloud project.",
		MarkdownDescription: "List the Swift object storage containers in a public cloud project.",
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
				Description:         "If set, only containers located in this region are returned",
				MarkdownDescription: "If set, only containers located in this region are returned",
			},
			"containers": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of containers",
				MarkdownDescription: "List of containers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: storageSwiftContainerDataSourceAttributes(),
				},
			},
		},
	}
}

type cloudStorageSwiftContainersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Containers  types.List             `tfsdk:"containers"`
}

func storageSwiftContainerListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"location":        types.ObjectType{AttrTypes: StorageSwiftContainerLocationAttrTypes()},
		"archive":         types.BoolType,
		"type":            ovhtypes.TfStringType{},
		"cors":            types.SetType{ElemType: ovhtypes.TfStringType{}},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: StorageSwiftContainerCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageSwiftContainersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageSwiftContainersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString())

	apiContainers, err := helpers.GetAllPagesV2[CloudStorageSwiftContainerAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	regionFilter := ""
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		regionFilter = data.Region.ValueString()
	}

	containerObjs := make([]attr.Value, 0, len(apiContainers))
	for i := range apiContainers {
		v := apiContainers[i]

		region := ""
		if v.TargetSpec != nil && v.TargetSpec.Location != nil {
			region = v.TargetSpec.Location.Region
		}
		if v.CurrentState != nil && v.CurrentState.Location != nil && v.CurrentState.Location.Region != "" {
			region = v.CurrentState.Location.Region
		}
		if regionFilter != "" && !strings.EqualFold(region, regionFilter) {
			continue
		}

		var item cloudStorageSwiftContainerDataSourceModel
		mapStorageSwiftContainerToDataSourceModel(ctx, &v, &item)

		obj, diags := types.ObjectValue(
			storageSwiftContainerListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"name":            item.Name,
				"location":        item.Location,
				"archive":         item.Archive,
				"type":            item.Type,
				"cors":            item.Cors,
				"checksum":        item.Checksum,
				"created_at":      item.CreatedAt,
				"updated_at":      item.UpdatedAt,
				"resource_status": item.ResourceStatus,
				"current_state":   item.CurrentState,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		containerObjs = append(containerObjs, obj)
	}

	containersList, diags := types.ListValue(
		types.ObjectType{AttrTypes: storageSwiftContainerListItemAttrTypes()},
		containerObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Containers = containersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
