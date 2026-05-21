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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumeSnapshotsDataSource)(nil)

func NewCloudStorageBlockVolumeSnapshotsDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumeSnapshotsDataSource{}
}

type cloudStorageBlockVolumeSnapshotsDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumeSnapshotsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume_snapshots"
}

func (d *cloudStorageBlockVolumeSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumeSnapshotsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List block storage volume snapshots for a given volume in a public cloud project.",
		MarkdownDescription: "List block storage volume snapshots for a given volume in a public cloud project.",
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
				Description:         "Region where the volume snapshots reside",
				MarkdownDescription: "Region where the volume snapshots reside",
			},
			"volume_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the volume whose snapshots to list",
				MarkdownDescription: "ID of the volume whose snapshots to list",
			},
			"snapshots": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of snapshots for the volume",
				MarkdownDescription: "List of snapshots for the volume",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
				},
			},
		},
	}
}

// cloudStorageBlockVolumeSnapshotsDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumeSnapshotsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	VolumeId    ovhtypes.TfStringValue `tfsdk:"volume_id"`
	Snapshots   types.List             `tfsdk:"snapshots"`
}

// snapshotListItemAttrTypes returns the attribute types for a single snapshot item in the list.
func snapshotListItemAttrTypes() map[string]attr.Type {
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

func (d *cloudStorageBlockVolumeSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumeSnapshotsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/block/snapshot?volumeId=" + url.QueryEscape(data.VolumeId.ValueString())

	var apiSnapshots []CloudStorageBlockVolumeSnapshotAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiSnapshots); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	snapshotObjs := make([]attr.Value, 0, len(apiSnapshots))
	for _, s := range apiSnapshots {
		// Only include snapshots matching the requested region (filter client-side)
		if s.CurrentState != nil && s.CurrentState.Location != nil &&
			s.CurrentState.Location.Region != data.Region.ValueString() {
			continue
		}

		name := ""
		description := ""
		volumeId := ""
		if s.TargetSpec != nil {
			name = s.TargetSpec.Name
			description = s.TargetSpec.Description
			volumeId = s.TargetSpec.VolumeId
		}

		size := int64(0)
		region := data.Region.ValueString()
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

		itemObj, diags := types.ObjectValue(
			snapshotListItemAttrTypes(),
			map[string]attr.Value{
				"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(s.Id)},
				"name":            ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"description":     ovhtypes.TfStringValue{StringValue: types.StringValue(description)},
				"location":        locObj,
				"volume_id":       ovhtypes.TfStringValue{StringValue: types.StringValue(volumeId)},
				"size":            types.Int64Value(size),
				"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(s.ResourceStatus)},
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		snapshotObjs = append(snapshotObjs, itemObj)
	}

	snapshotsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: snapshotListItemAttrTypes()},
		snapshotObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Snapshots = snapshotsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
