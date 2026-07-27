package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareAclsDataSource)(nil)

func NewCloudStorageFileShareAclsDataSource() datasource.DataSource {
	return &cloudStorageFileShareAclsDataSource{}
}

type cloudStorageFileShareAclsDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareAclsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_acls"
}

func (d *cloudStorageFileShareAclsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageFileShareAclsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the access rules (ACLs) of a public cloud file storage share.",
		MarkdownDescription: "List the access rules (ACLs) of a public cloud file storage share.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"share_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the file storage share the access rules apply to",
				MarkdownDescription: "ID of the file storage share the access rules apply to",
			},
			"share_acls": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of access rules of the file storage share",
				MarkdownDescription: "List of access rules of the file storage share",
				NestedObject: schema.NestedAttributeObject{
					Attributes: fileShareAclDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudStorageFileShareAclsDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareAclsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ShareId     ovhtypes.TfStringValue `tfsdk:"share_id"`
	ShareAcls   types.List             `tfsdk:"share_acls"`
}

// fileShareAclListItemAttrTypes returns the attribute types for a single access rule item in the list.
func fileShareAclListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"access_to":       ovhtypes.TfStringType{},
		"access_level":    ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: StorageFileShareAclCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageFileShareAclsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareAclsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString())

	var apiAcls []CloudStorageFileShareAclAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiAcls); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	aclObjs := make([]attr.Value, 0, len(apiAcls))
	for i := range apiAcls {
		v := apiAcls[i]

		var item cloudStorageFileShareAclDataSourceModel
		mapFileShareAclToDataSourceModel(&v, &item)

		obj, diags := types.ObjectValue(
			fileShareAclListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"access_to":       item.AccessTo,
				"access_level":    item.AccessLevel,
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

		aclObjs = append(aclObjs, obj)
	}

	aclsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: fileShareAclListItemAttrTypes()},
		aclObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ShareAcls = aclsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
