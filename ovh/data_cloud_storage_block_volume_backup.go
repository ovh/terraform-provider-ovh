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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumeBackupDataSource)(nil)

func NewCloudStorageBlockVolumeBackupDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumeBackupDataSource{}
}

type cloudStorageBlockVolumeBackupDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumeBackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume_backup"
}

func (d *cloudStorageBlockVolumeBackupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumeBackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get a block storage volume backup in a public cloud project.",
		MarkdownDescription: "Get a block storage volume backup in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Backup ID",
				MarkdownDescription: "Backup ID",
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
	}
}

// cloudStorageBlockVolumeBackupDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumeBackupDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	Location       types.Object           `tfsdk:"location"`
	VolumeId       ovhtypes.TfStringValue `tfsdk:"volume_id"`
	Size           types.Int64            `tfsdk:"size"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
}

func (d *cloudStorageBlockVolumeBackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumeBackupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/block/backup/" + url.PathEscape(data.Id.ValueString())

	var b CloudStorageBlockVolumeBackupAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &b); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	name := ""
	description := ""
	volumeId := ""
	region := ""
	if b.TargetSpec != nil {
		name = b.TargetSpec.Name
		description = b.TargetSpec.Description
		volumeId = b.TargetSpec.VolumeId
		if b.TargetSpec.Location != nil {
			region = b.TargetSpec.Location.Region
		}
	}

	size := int64(0)
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

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(b.Id)}
	data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(name)}
	data.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(description)}
	data.Location = locObj
	data.VolumeId = ovhtypes.TfStringValue{StringValue: types.StringValue(volumeId)}
	data.Size = types.Int64Value(size)
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(b.ResourceStatus)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
