// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"context"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func OvhcloudConnectConfigPopDatacenterExtraDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"pop_id": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Required:            true,
			Description:         "Pop ID",
			MarkdownDescription: "Pop ID",
		},
		"datacenter_id": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Required:            true,
			Description:         "Datacenter ID",
			MarkdownDescription: "Datacenter ID",
		},
		"extra_configs": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"bgp_neighbor_area": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "BGP AS number",
						MarkdownDescription: "BGP AS number",
					},
					"bgp_neighbor_ip": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Router IP for BGP",
						MarkdownDescription: "Router IP for BGP",
					},
					"datacenter_id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Required:            true,
						Description:         "Datacenter ID",
						MarkdownDescription: "Datacenter ID",
					},
					"extra_id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Required:            true,
						Description:         "Extra ID",
						MarkdownDescription: "Extra ID",
					},
					"id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "ID of the extra configuration ",
						MarkdownDescription: "ID of the extra configuration ",
					},
					"next_hop": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Static route next hop",
						MarkdownDescription: "Static route next hop",
					},
					"pop_id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Required:            true,
						Description:         "Pop ID",
						MarkdownDescription: "Pop ID",
					},
					"service_name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Service name",
						MarkdownDescription: "Service name",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Status of the pop configuration",
						MarkdownDescription: "Status of the pop configuration",
					},
					"subnet": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Static route ip",
						MarkdownDescription: "Static route ip",
					},
					"type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Type of the configuration",
						MarkdownDescription: "Type of the configuration",
					},
				},
			},
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type OvhcloudConnectConfigsPopDatacenterExtraModel struct {
	ExtraConfigs []OvhcloudConnectConfigPopDatacenterExtraModel `tfsdk:"extra_configs"`
	DatacenterId ovhtypes.TfInt64Value                          `tfsdk:"datacenter_id" json:"datacenterId"`
	PopId        ovhtypes.TfInt64Value                          `tfsdk:"pop_id" json:"popId"`
	ServiceName  ovhtypes.TfStringValue                         `tfsdk:"service_name" json:"serviceName"`
}

type OvhcloudConnectConfigPopDatacenterExtraModel struct {
	BgpNeighborArea ovhtypes.TfInt64Value  `tfsdk:"bgp_neighbor_area" json:"bgpNeighborArea"`
	BgpNeighborIp   ovhtypes.TfStringValue `tfsdk:"bgp_neighbor_ip" json:"bgpNeighborIp"`
	DatacenterId    ovhtypes.TfInt64Value  `tfsdk:"datacenter_id" json:"datacenterId"`
	ExtraId         ovhtypes.TfInt64Value  `tfsdk:"extra_id" json:"extraId"`
	Id              ovhtypes.TfInt64Value  `tfsdk:"id" json:"id"`
	NextHop         ovhtypes.TfStringValue `tfsdk:"next_hop" json:"nextHop"`
	PopId           ovhtypes.TfInt64Value  `tfsdk:"pop_id" json:"popId"`
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name" json:"serviceName"`
	Status          ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Subnet          ovhtypes.TfStringValue `tfsdk:"subnet" json:"subnet"`
	Type            ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
}

func (v *OvhcloudConnectConfigPopDatacenterExtraModel) MergeWith(other *OvhcloudConnectConfigPopDatacenterExtraModel) {

	if (v.BgpNeighborArea.IsUnknown() || v.BgpNeighborArea.IsNull()) && !other.BgpNeighborArea.IsUnknown() {
		v.BgpNeighborArea = other.BgpNeighborArea
	}

	if (v.BgpNeighborIp.IsUnknown() || v.BgpNeighborIp.IsNull()) && !other.BgpNeighborIp.IsUnknown() {
		v.BgpNeighborIp = other.BgpNeighborIp
	}

	if (v.DatacenterId.IsUnknown() || v.DatacenterId.IsNull()) && !other.DatacenterId.IsUnknown() {
		v.DatacenterId = other.DatacenterId
	}

	if (v.ExtraId.IsUnknown() || v.ExtraId.IsNull()) && !other.ExtraId.IsUnknown() {
		v.ExtraId = other.ExtraId
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.NextHop.IsUnknown() || v.NextHop.IsNull()) && !other.NextHop.IsUnknown() {
		v.NextHop = other.NextHop
	}

	if (v.PopId.IsUnknown() || v.PopId.IsNull()) && !other.PopId.IsUnknown() {
		v.PopId = other.PopId
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

	if (v.Status.IsUnknown() || v.Status.IsNull()) && !other.Status.IsUnknown() {
		v.Status = other.Status
	}

	if (v.Subnet.IsUnknown() || v.Subnet.IsNull()) && !other.Subnet.IsUnknown() {
		v.Subnet = other.Subnet
	}

	if (v.Type.IsUnknown() || v.Type.IsNull()) && !other.Type.IsUnknown() {
		v.Type = other.Type
	}

}
