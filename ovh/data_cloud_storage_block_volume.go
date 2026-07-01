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

var _ datasource.DataSourceWithConfigure = (*cloudStorageBlockVolumeDataSource)(nil)

func NewCloudStorageBlockVolumeDataSource() datasource.DataSource {
	return &cloudStorageBlockVolumeDataSource{}
}

type cloudStorageBlockVolumeDataSource struct {
	config *Config
}

func (d *cloudStorageBlockVolumeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume"
}

func (d *cloudStorageBlockVolumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageBlockVolumeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get a block storage volume in a public cloud project.",
		MarkdownDescription: "Get a block storage volume in a public cloud project.",
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
				Description:         "Volume ID",
				MarkdownDescription: "Volume ID",
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
	}
}

// cloudStorageBlockVolumeDataSourceModel is the Terraform state model for this data source.
type cloudStorageBlockVolumeDataSourceModel struct {
	ServiceName       ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id                ovhtypes.TfStringValue `tfsdk:"id"`
	Name              ovhtypes.TfStringValue `tfsdk:"name"`
	Location          types.Object           `tfsdk:"location"`
	Size              types.Int64            `tfsdk:"size"`
	VolumeType        ovhtypes.TfStringValue `tfsdk:"volume_type"`
	Bootable          types.Bool             `tfsdk:"bootable"`
	Status            ovhtypes.TfStringValue `tfsdk:"status"`
	ResourceStatus    ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Encryption        types.Object           `tfsdk:"encryption"`
	AttachedInstances types.List             `tfsdk:"attached_instances"`
}

func (d *cloudStorageBlockVolumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageBlockVolumeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/block/volume/" + url.PathEscape(data.Id.ValueString())

	var v CloudStorageBlockVolumeAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	name := ""
	size := int64(0)
	volumeType := ""
	region := ""
	var encryption *CloudStorageBlockVolumeEncryption
	if v.TargetSpec != nil {
		name = v.TargetSpec.Name
		size = v.TargetSpec.Size
		volumeType = v.TargetSpec.VolumeType
		encryption = v.TargetSpec.Encryption
		if v.TargetSpec.Location != nil {
			region = v.TargetSpec.Location.Region
		}
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

	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(name)}
	data.Location = locObj
	data.Size = types.Int64Value(size)
	data.VolumeType = ovhtypes.TfStringValue{StringValue: types.StringValue(volumeType)}
	data.Bootable = types.BoolValue(bootable)
	data.Status = ovhtypes.TfStringValue{StringValue: types.StringValue(status)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}
	data.Encryption = encryptionObj
	data.AttachedInstances = attachedInstancesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
