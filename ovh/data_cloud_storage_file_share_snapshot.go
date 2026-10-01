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

var _ datasource.DataSourceWithConfigure = (*cloudStorageFileShareSnapshotDataSource)(nil)

func NewCloudStorageFileShareSnapshotDataSource() datasource.DataSource {
	return &cloudStorageFileShareSnapshotDataSource{}
}

type cloudStorageFileShareSnapshotDataSource struct {
	config *Config
}

func (d *cloudStorageFileShareSnapshotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_file_share_snapshot"
}

func (d *cloudStorageFileShareSnapshotDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fileShareSnapshotDataSourceAttributes returns the computed attributes describing a
// file storage snapshot, shared between the singular and plural data sources.
func fileShareSnapshotDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Snapshot ID",
		},
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Snapshot name",
		},
		"description": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Snapshot description",
		},
		"share_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "ID of the snapshotted file share",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the snapshot",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the snapshot",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Snapshot readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current state of the file storage snapshot",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Snapshot name",
				},
				"description": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Snapshot description",
				},
				"share_id": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "ID of the snapshotted file share",
				},
				"size": schema.Int64Attribute{
					Computed:    true,
					Description: "Size of the snapshot in GB",
				},
				"location": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Current location",
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region",
						},
						"availability_zone": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Availability zone",
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageFileShareSnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name of the resource representing the id of the cloud project",
			MarkdownDescription: "Service name of the resource representing the id of the cloud project",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Snapshot ID",
			MarkdownDescription: "Snapshot ID",
		},
	}

	for name, attribute := range fileShareSnapshotDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get a file storage snapshot (NFS share snapshot) in a public cloud project.",
		MarkdownDescription: "Get a file storage snapshot (NFS share snapshot) in a public cloud project.",
		Attributes:          attrs,
	}
}

// cloudStorageFileShareSnapshotDataSourceModel is the Terraform state model for this data source.
type cloudStorageFileShareSnapshotDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	ShareId        ovhtypes.TfStringValue `tfsdk:"share_id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageFileShareSnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageFileShareSnapshotDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/file/snapshot/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageFileShareSnapshotAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapFileShareSnapshotToDataSourceModel(&v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildFileShareSnapshotCurrentStateObject builds the current_state object from the API response.
func buildFileShareSnapshotCurrentStateObject(state *CloudStorageFileShareSnapshotCurrentState) types.Object {
	if state == nil {
		return types.ObjectNull(FileShareSnapshotCurrentStateAttrTypes())
	}

	shareId := ""
	if state.Share != nil {
		shareId = state.Share.Id
	}

	region := ""
	az := ""
	if state.Location != nil {
		region = state.Location.Region
		az = state.Location.AvailabilityZone
	}
	locObj, _ := types.ObjectValue(
		map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(az)},
		},
	)

	obj, _ := types.ObjectValue(
		FileShareSnapshotCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"share_id":    ovhtypes.TfStringValue{StringValue: types.StringValue(shareId)},
			"size":        types.Int64Value(state.Size),
			"location":    locObj,
		},
	)
	return obj
}

// mapFileShareSnapshotToDataSourceModel populates the data source model from the API response.
func mapFileShareSnapshotToDataSourceModel(v *CloudStorageFileShareSnapshotAPIResponse, data *cloudStorageFileShareSnapshotDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Name)}
		data.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Description)}
		if v.TargetSpec.Share != nil {
			data.ShareId = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Share.Id)}
		}
	}

	data.CurrentState = buildFileShareSnapshotCurrentStateObject(v.CurrentState)
}
