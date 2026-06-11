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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumesDataSource)(nil)

func NewCloudStorageBlockVolumesDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumesDataSource{}
}

type cloudStorageBlockVolumesDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volumes"
}

func (d *cloudStorageBlockVolumesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List block storage volumes in a public cloud project.",
		MarkdownDescription: "List block storage volumes in a public cloud project.",
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
				Description:         "Region where the volumes reside",
				MarkdownDescription: "Region where the volumes reside",
			},
			"volumes": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of volumes",
				MarkdownDescription: "List of volumes",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Volume ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Volume name",
						},
						"location": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Location of the volume",
							Attributes: map[string]schema.Attribute{
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region",
								},
							},
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the volume in GB",
						},
						"volume_type": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Volume type (CLASSIC, HIGH_SPEED, HIGH_SPEED_GEN2)",
						},
						"bootable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the volume is bootable",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Volume status",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Volume readiness status",
						},
						"encryption": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Encryption configuration of the volume",
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the volume is encrypted at rest with LUKS",
								},
							},
						},
						"attached_instances": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Instances the volume is attached to",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Instance ID",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// cloudStorageBlockVolumesDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumesDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Volumes     types.List             `tfsdk:"volumes"`
}

// volumeListItemAttrTypes returns the attribute types for a single volume item in the list.
func volumeListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"name": ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region": ovhtypes.TfStringType{},
		}},
		"size":               types.Int64Type,
		"volume_type":        ovhtypes.TfStringType{},
		"bootable":           types.BoolType,
		"status":             ovhtypes.TfStringType{},
		"resource_status":    ovhtypes.TfStringType{},
		"encryption":         types.ObjectType{AttrTypes: BlockVolumeEncryptionAttrTypes()},
		"attached_instances": types.ListType{ElemType: types.ObjectType{AttrTypes: BlockVolumeAttachedInstanceAttrTypes()}},
	}
}

func (d *cloudStorageBlockVolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume"

	var apiVolumes []CloudStorageBlockVolumeAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &apiVolumes); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	volumeObjs := make([]attr.Value, 0, len(apiVolumes))
	for _, v := range apiVolumes {
		// Only include volumes matching the requested region (filter client-side)
		if v.CurrentState != nil && v.CurrentState.Location != nil &&
			v.CurrentState.Location.Region != data.Region.ValueString() {
			continue
		}

		name := ""
		size := int64(0)
		volumeType := ""
		region := data.Region.ValueString()
		var encryption *CloudStorageBlockVolumeEncryption
		if v.TargetSpec != nil {
			name = v.TargetSpec.Name
			size = v.TargetSpec.Size
			volumeType = v.TargetSpec.VolumeType
			encryption = v.TargetSpec.Encryption
		}

		bootable := false
		status := ""
		var attachedInstances []CloudStorageBlockVolumeAttachedInstance
		if v.CurrentState != nil {
			if v.CurrentState.Name != "" {
				name = v.CurrentState.Name
			}
			if v.CurrentState.Size != 0 {
				size = v.CurrentState.Size
			}
			if v.CurrentState.VolumeType != "" {
				volumeType = v.CurrentState.VolumeType
			}
			if v.CurrentState.Location != nil {
				region = v.CurrentState.Location.Region
			}
			if v.CurrentState.Bootable != nil {
				bootable = *v.CurrentState.Bootable
			}
			if v.CurrentState.Encryption != nil {
				encryption = v.CurrentState.Encryption
			}
			status = v.CurrentState.Status
			attachedInstances = v.CurrentState.AttachedInstances
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

		encryptionObj := types.ObjectNull(BlockVolumeEncryptionAttrTypes())
		if encryption != nil {
			encryptionObj, diags = types.ObjectValue(
				BlockVolumeEncryptionAttrTypes(),
				map[string]attr.Value{
					"enabled": types.BoolValue(encryption.Enabled),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		attachedInstanceObjType := types.ObjectType{AttrTypes: BlockVolumeAttachedInstanceAttrTypes()}
		attachedInstanceObjs := make([]attr.Value, 0, len(attachedInstances))
		for _, ai := range attachedInstances {
			aiObj, diags := types.ObjectValue(
				BlockVolumeAttachedInstanceAttrTypes(),
				map[string]attr.Value{
					"id": ovhtypes.TfStringValue{StringValue: types.StringValue(ai.Id)},
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			attachedInstanceObjs = append(attachedInstanceObjs, aiObj)
		}
		attachedInstancesList, diags := types.ListValue(attachedInstanceObjType, attachedInstanceObjs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		itemObj, diags := types.ObjectValue(
			volumeListItemAttrTypes(),
			map[string]attr.Value{
				"id":                 ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)},
				"name":               ovhtypes.TfStringValue{StringValue: types.StringValue(name)},
				"location":           locObj,
				"size":               types.Int64Value(size),
				"volume_type":        ovhtypes.TfStringValue{StringValue: types.StringValue(volumeType)},
				"bootable":           types.BoolValue(bootable),
				"status":             ovhtypes.TfStringValue{StringValue: types.StringValue(status)},
				"resource_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)},
				"encryption":         encryptionObj,
				"attached_instances": attachedInstancesList,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		volumeObjs = append(volumeObjs, itemObj)
	}

	volumesList, diags := types.ListValue(
		types.ObjectType{AttrTypes: volumeListItemAttrTypes()},
		volumeObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Volumes = volumesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
