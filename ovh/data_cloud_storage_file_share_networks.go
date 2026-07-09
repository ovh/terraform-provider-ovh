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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareNetworksDataSource)(nil)

func NewCloudStorageFileShareNetworksDataSource() datasource.DataSource {
	return &cloudStorageFileShareNetworksDataSource{}
}

type cloudStorageFileShareNetworksDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_networks"
}

func (d *cloudStorageFileShareNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageFileShareNetworksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List file storage share networks in a public cloud project.",
		MarkdownDescription: "List file storage share networks in a public cloud project.",
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
				Description:         "If set, only share networks located in this region are returned",
				MarkdownDescription: "If set, only share networks located in this region are returned",
			},
			"share_networks": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of file storage share networks",
				MarkdownDescription: "List of file storage share networks",
				NestedObject: schema.NestedAttributeObject{
					Attributes: fileShareNetworkDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudStorageFileShareNetworksDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareNetworksDataSourceModel struct {
	ServiceName   ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region        ovhtypes.TfStringValue `tfsdk:"region"`
	ShareNetworks types.List             `tfsdk:"share_networks"`
}

// fileShareNetworkListItemAttrTypes returns the attribute types for a single share network item in the list.
func fileShareNetworkListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"description":     ovhtypes.TfStringType{},
		"network_id":      ovhtypes.TfStringType{},
		"subnet_id":       ovhtypes.TfStringType{},
		"location":        types.ObjectType{AttrTypes: fileShareNetworkLocationAttrTypes()},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: FileShareNetworkCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageFileShareNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareNetworksDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/network"

	var apiNetworks []CloudStorageFileShareNetworkAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiNetworks); err != nil {
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

	networkObjs := make([]attr.Value, 0, len(apiNetworks))
	for i := range apiNetworks {
		v := apiNetworks[i]

		region := ""
		if v.TargetSpec != nil && v.TargetSpec.Location != nil {
			region = v.TargetSpec.Location.Region
		}
		if v.CurrentState != nil && v.CurrentState.Location != nil && v.CurrentState.Location.Region != "" {
			region = v.CurrentState.Location.Region
		}
		if regionFilter != "" && region != regionFilter {
			continue
		}

		var item cloudStorageFileShareNetworkDataSourceModel
		mapFileShareNetworkToDataSourceModel(&v, &item)

		obj, diags := types.ObjectValue(
			fileShareNetworkListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"name":            item.Name,
				"description":     item.Description,
				"network_id":      item.NetworkId,
				"subnet_id":       item.SubnetId,
				"location":        item.Location,
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

		networkObjs = append(networkObjs, obj)
	}

	networksList, diags := types.ListValue(
		types.ObjectType{AttrTypes: fileShareNetworkListItemAttrTypes()},
		networkObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ShareNetworks = networksList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
