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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumeBackupsDataSource)(nil)

func NewCloudStorageBlockVolumeBackupsDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumeBackupsDataSource{}
}

type cloudStorageBlockVolumeBackupsDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumeBackupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume_backups"
}

func (d *cloudStorageBlockVolumeBackupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumeBackupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List block storage volume backups for a given volume in a public cloud project.",
		MarkdownDescription: "List block storage volume backups for a given volume in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the volume backups reside",
				MarkdownDescription: "Region where the volume backups reside",
			},
			"volume_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the volume whose backups to list",
				MarkdownDescription: "ID of the volume whose backups to list",
			},
			"backups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of backups for the volume",
				MarkdownDescription: "List of backups for the volume",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Backup ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Backup name",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Backup description",
						},
						"location": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Location of the backup",
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region",
								},
							},
						},
						"volume_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "ID of the backed-up volume",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the backup in GB",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Backup readiness status",
						},
					},
				},
			},
		},
	}
}

// cloudStorageBlockVolumeBackupsDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumeBackupsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	VolumeId    ovhtypes.TfStringValue `tfsdk:"volume_id"`
	Backups     types.List             `tfsdk:"backups"`
}

// backupListItemAttrTypes returns the attribute types for a single backup item in the list.
func backupListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region": ovhtypes.TfStringType{},
		}},
		"volume_id":       ovhtypes.TfStringType{},
		"size":            types.Int64Type,
		"resource_status": ovhtypes.TfStringType{},
	}
}

func (d *cloudStorageBlockVolumeBackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumeBackupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/block/backup?volumeId=" + url.QueryEscape(data.VolumeId.ValueString())

	var apiBackups []CloudStorageBlockVolumeBackupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiBackups); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	backupObjs := make([]attr.Value, 0, len(apiBackups))
	for _, b := range apiBackups {
		// Only include backups matching the requested region (filter client-side)
		if b.CurrentState != nil && b.CurrentState.Location != nil &&
			b.CurrentState.Location.Region != data.Region.ValueString() {
			continue
		}

		name := ""
		description := ""
		volumeId := ""
		if b.TargetSpec != nil {
			name = b.TargetSpec.Name
			description = b.TargetSpec.Description
			volumeId = b.TargetSpec.VolumeId
		}

		size := int64(0)
		region := data.Region.ValueString()
		if b.CurrentState != nil {
			size = b.CurrentState.Size
			if b.CurrentState.VolumeId != "" {
				volumeId = b.CurrentState.VolumeId
			}
			if b.CurrentState.Name != "" {
				name = b.CurrentState.Name
			}
			if b.CurrentState.Description != "" {
				description = b.CurrentState.Description
			}
			if b.CurrentState.Location != nil {
				region = b.CurrentState.Location.Region
			}
		}

		locObj, diags := types.ObjectValue(
			map[string]attr.Type{"region": ovhtypes.TfStringType{}},
			map[string]attr.Value{
				"region": ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		itemObj, diags := types.ObjectValue(
			backupListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(b.Id)},
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"description":     ovhtypes.TfStringValue{StringValue: types.StringValue(description)},
				"location":        locObj,
				"volume_id":       ovhtypes.TfStringValue{StringValue: types.StringValue(volumeId)},
				"size":            types.Int64Value(size),
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(b.ResourceStatus)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		backupObjs = append(backupObjs, itemObj)
	}

	backupsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: backupListItemAttrTypes()},
		backupObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Backups = backupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
