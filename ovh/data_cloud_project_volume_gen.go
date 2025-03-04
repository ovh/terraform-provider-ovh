// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"context"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CloudProjectVolumeDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Volume ID",
			MarkdownDescription: "Volume ID",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Volume name",
			MarkdownDescription: "Volume name",
		},
		"region_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Region name",
			MarkdownDescription: "Region name",
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"size": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Volume size",
			MarkdownDescription: "Volume size",
		},
		"volume_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Volume ID",
			MarkdownDescription: "Volume ID",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudProjectVolumeModel struct {
	Id          ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	RegionName  ovhtypes.TfStringValue `tfsdk:"region_name" json:"regionName"`
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name" json:"serviceName"`
	Size        ovhtypes.TfInt64Value  `tfsdk:"size" json:"size"`
	VolumeId    ovhtypes.TfStringValue `tfsdk:"volume_id" json:"volumeId"`
}

func (v *CloudProjectVolumeModel) MergeWith(other *CloudProjectVolumeModel) {

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.RegionName.IsUnknown() || v.RegionName.IsNull()) && !other.RegionName.IsUnknown() {
		v.RegionName = other.RegionName
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

	if (v.Size.IsUnknown() || v.Size.IsNull()) && !other.Size.IsUnknown() {
		v.Size = other.Size
	}

	if (v.VolumeId.IsUnknown() || v.VolumeId.IsNull()) && !other.VolumeId.IsUnknown() {
		v.VolumeId = other.VolumeId
	}

}
