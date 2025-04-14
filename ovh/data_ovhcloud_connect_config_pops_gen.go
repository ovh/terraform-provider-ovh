// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"context"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func OvhcloudConnectConfigPopsDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"pop_configs": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"customer_bgp_area": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Customer Private AS",
						MarkdownDescription: "Customer Private AS",
					},
					"id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "ID of the Pop Configuration",
						MarkdownDescription: "ID of the Pop Configuration",
					},
					"interface_id": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "ID of the interface",
						MarkdownDescription: "ID of the interface",
					},
					"ovh_bgp_area": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "OVH Private AS",
						MarkdownDescription: "OVH Private AS",
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
						Description:         "Subnet should be a /30, first IP for OVH, second IP for customer",
						MarkdownDescription: "Subnet should be a /30, first IP for OVH, second IP for customer",
					},
					"type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Type of the pop configuration",
						MarkdownDescription: "Type of the pop configuration",
					},
				},
			},
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type OvhcloudConnectConfigPopsModel struct {
	PopConfigs  []OvhcloudConnectConfigPopModel `tfsdk:"pop_configs"`
	ServiceName ovhtypes.TfStringValue          `tfsdk:"service_name" json:"serviceName"`
}

type OvhcloudConnectConfigPopModel struct {
	CustomerBgpArea ovhtypes.TfInt64Value  `tfsdk:"customer_bgp_area" json:"customerBgpArea"`
	Id              ovhtypes.TfInt64Value  `tfsdk:"id" json:"id"`
	InterfaceId     ovhtypes.TfInt64Value  `tfsdk:"interface_id" json:"interfaceId"`
	OvhBgpArea      ovhtypes.TfInt64Value  `tfsdk:"ovh_bgp_area" json:"ovhBgpArea"`
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name" json:"serviceName"`
	Status          ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Subnet          ovhtypes.TfStringValue `tfsdk:"subnet" json:"subnet"`
	Type            ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
}

func (v *OvhcloudConnectConfigPopModel) MergeWith(other *OvhcloudConnectConfigPopModel) {

	if (v.CustomerBgpArea.IsUnknown() || v.CustomerBgpArea.IsNull()) && !other.CustomerBgpArea.IsUnknown() {
		v.CustomerBgpArea = other.CustomerBgpArea
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.InterfaceId.IsUnknown() || v.InterfaceId.IsNull()) && !other.InterfaceId.IsUnknown() {
		v.InterfaceId = other.InterfaceId
	}

	if (v.OvhBgpArea.IsUnknown() || v.OvhBgpArea.IsNull()) && !other.OvhBgpArea.IsUnknown() {
		v.OvhBgpArea = other.OvhBgpArea
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
