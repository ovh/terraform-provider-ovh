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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileSharesDataSource)(nil)

func NewCloudStorageFileSharesDataSource() datasource.DataSource {
	return &cloudStorageFileSharesDataSource{}
}

type cloudStorageFileSharesDataSource struct {
	config *Config
}

func (d *cloudStorageFileSharesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_shares"
}

func (d *cloudStorageFileSharesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageFileSharesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List file storage shares (NFS) in a public cloud project.",
		MarkdownDescription: "List file storage shares (NFS) in a public cloud project.",
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
				Description:         "If set, only file shares located in this region are returned",
				MarkdownDescription: "If set, only file shares located in this region are returned",
			},
			"file_shares": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of file storage shares",
				MarkdownDescription: "List of file storage shares",
				NestedObject: schema.NestedAttributeObject{
					Attributes: fileShareDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudStorageFileSharesDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileSharesDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	FileShares  types.List             `tfsdk:"file_shares"`
}

// fileShareListItemAttrTypes returns the attribute types for a single file share item in the list.
func fileShareListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               ovhtypes.TfStringType{},
		"name":             ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
		"size":             types.Int64Type,
		"protocol":         ovhtypes.TfStringType{},
		"share_type":       ovhtypes.TfStringType{},
		"location":         types.ObjectType{AttrTypes: fileShareLocationAttrTypes()},
		"share_network_id": ovhtypes.TfStringType{},
		"access_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: FileShareAccessRuleAttrTypes()}},
		"checksum":         ovhtypes.TfStringType{},
		"created_at":       ovhtypes.TfStringType{},
		"updated_at":       ovhtypes.TfStringType{},
		"resource_status":  ovhtypes.TfStringType{},
		"current_state":    types.ObjectType{AttrTypes: FileShareCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageFileSharesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileSharesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/share"

	var apiShares []CloudStorageFileShareAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiShares); err != nil {
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

	shareObjs := make([]attr.Value, 0, len(apiShares))
	for i := range apiShares {
		v := apiShares[i]

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

		// Reuse the singular mapping to keep the shape consistent.
		var item cloudStorageFileShareDataSourceModel
		mapFileShareToDataSourceModel(ctx, &v, &item)

		obj, diags := types.ObjectValue(
			fileShareListItemAttrTypes(),
			map[string]attr.Value{
				"id":               item.Id,
				"name":             item.Name,
				"description":      item.Description,
				"size":             item.Size,
				"protocol":         item.Protocol,
				"share_type":       item.ShareType,
				"location":         item.Location,
				"share_network_id": item.ShareNetworkId,
				"access_rules":     item.AccessRules,
				"checksum":         item.Checksum,
				"created_at":       item.CreatedAt,
				"updated_at":       item.UpdatedAt,
				"resource_status":  item.ResourceStatus,
				"current_state":    item.CurrentState,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		shareObjs = append(shareObjs, obj)
	}

	sharesList, diags := types.ListValue(
		types.ObjectType{AttrTypes: fileShareListItemAttrTypes()},
		shareObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FileShares = sharesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
