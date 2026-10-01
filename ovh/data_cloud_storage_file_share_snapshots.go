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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareSnapshotsDataSource)(nil)

func NewCloudStorageFileShareSnapshotsDataSource() datasource.DataSource {
	return &cloudStorageFileShareSnapshotsDataSource{}
}

type cloudStorageFileShareSnapshotsDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareSnapshotsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_snapshots"
}

func (d *cloudStorageFileShareSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageFileShareSnapshotsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List file storage snapshots (NFS share snapshots) in a public cloud project.",
		MarkdownDescription: "List file storage snapshots (NFS share snapshots) in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"share_snapshots": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of file storage snapshots",
				MarkdownDescription: "List of file storage snapshots",
				NestedObject: schema.NestedAttributeObject{
					Attributes: fileShareSnapshotDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudStorageFileShareSnapshotsDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareSnapshotsDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	ShareSnapshots types.List             `tfsdk:"share_snapshots"`
}

// fileShareSnapshotListItemAttrTypes returns the attribute types for a single snapshot item in the list.
func fileShareSnapshotListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"description":     ovhtypes.TfStringType{},
		"share_id":        ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: FileShareSnapshotCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageFileShareSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareSnapshotsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/file/snapshot"

	var apiSnapshots []CloudStorageFileShareSnapshotAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiSnapshots); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	snapshotObjs := make([]attr.Value, 0, len(apiSnapshots))
	for i := range apiSnapshots {
		v := apiSnapshots[i]

		var item cloudStorageFileShareSnapshotDataSourceModel
		mapFileShareSnapshotToDataSourceModel(&v, &item)

		obj, diags := types.ObjectValue(
			fileShareSnapshotListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"name":            item.Name,
				"description":     item.Description,
				"share_id":        item.ShareId,
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

		snapshotObjs = append(snapshotObjs, obj)
	}

	snapshotsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: fileShareSnapshotListItemAttrTypes()},
		snapshotObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ShareSnapshots = snapshotsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
