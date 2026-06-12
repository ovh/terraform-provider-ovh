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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumeSnapshotDataSource)(nil)

func NewCloudStorageBlockVolumeSnapshotDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumeSnapshotDataSource{}
}

type cloudStorageBlockVolumeSnapshotDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumeSnapshotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume_snapshot"
}

func (d *cloudStorageBlockVolumeSnapshotDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumeSnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get a block storage volume snapshot in a public cloud project.",
		MarkdownDescription: "Get a block storage volume snapshot in a public cloud project.",
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
				Description:         "Snapshot ID",
				MarkdownDescription: "Snapshot ID",
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
			"location": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Location of the snapshot",
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
				Description: "ID of the snapshotted volume",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Size of the snapshot in GB",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Snapshot readiness status",
			},
		},
	}
}

// cloudStorageBlockVolumeSnapshotDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumeSnapshotDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Description    ovhtypes.TfStringValue `tfsdk:"description"`
	Location       types.Object           `tfsdk:"location"`
	VolumeId       ovhtypes.TfStringValue `tfsdk:"volume_id"`
	Size           types.Int64            `tfsdk:"size"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
}

func (d *cloudStorageBlockVolumeSnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumeSnapshotDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/block/snapshot/" + url.PathEscape(data.Id.ValueString())

	var s CloudStorageBlockVolumeSnapshotAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &s); err != nil {
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
	if s.TargetSpec != nil {
		name = s.TargetSpec.Name
		description = s.TargetSpec.Description
		volumeId = s.TargetSpec.VolumeId
		if s.TargetSpec.Location != nil {
			region = s.TargetSpec.Location.Region
		}
	}

	size := int64(0)
	if s.CurrentState != nil {
		size = s.CurrentState.Size
		if s.CurrentState.VolumeId != "" {
			volumeId = s.CurrentState.VolumeId
		}
		if s.CurrentState.Name != "" {
			name = s.CurrentState.Name
		}
		if s.CurrentState.Description != "" {
			description = s.CurrentState.Description
		}
		if s.CurrentState.Location != nil {
			region = s.CurrentState.Location.Region
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

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(s.Id)}
	data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(name)}
	data.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(description)}
	data.Location = locObj
	data.VolumeId = ovhtypes.TfStringValue{StringValue: types.StringValue(volumeId)}
	data.Size = types.Int64Value(size)
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(s.ResourceStatus)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
