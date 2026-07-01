package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudRegion is the API response for publicCloud.reference.Region (OVHcloud API v2),
// returned by GET /v2/publicCloud/project/{serviceName}/reference/region[/{name}].
type CloudRegion struct {
	Name              string   `json:"name"`
	Status            string   `json:"status"`
	Continent         string   `json:"continent"`
	Country           string   `json:"country"`
	DatacenterName    string   `json:"datacenterName"`
	AvailabilityZones []string `json:"availabilityZones"`
	Services          []string `json:"services"`
}

// cloudRegionDataSourceModel is the Terraform model for the singular ovh_cloud_region data source.
type cloudRegionDataSourceModel struct {
	ServiceName       ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name              ovhtypes.TfStringValue `tfsdk:"name"`
	Status            ovhtypes.TfStringValue `tfsdk:"status"`
	Continent         ovhtypes.TfStringValue `tfsdk:"continent"`
	Country           ovhtypes.TfStringValue `tfsdk:"country"`
	DatacenterName    ovhtypes.TfStringValue `tfsdk:"datacenter_name"`
	AvailabilityZones types.List             `tfsdk:"availability_zones"`
	Services          types.List             `tfsdk:"services"`
}

// cloudRegionsDataSourceModel is the Terraform model for the plural ovh_cloud_regions data source.
type cloudRegionsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Regions     types.List             `tfsdk:"regions"`
}

// cloudRegionDetailAttributes returns the computed schema attributes describing a single
// region. It is shared between the singular data source and the plural data source's nested
// "regions" objects. The "name" attribute is intentionally excluded so each data source can
// declare it as required (singular input) or computed (plural list item).
func cloudRegionDetailAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Region status (ENABLED, DISABLED or MAINTENANCE)",
			MarkdownDescription: "Region status (ENABLED, DISABLED or MAINTENANCE)",
		},
		"continent": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Continent code of the region",
			MarkdownDescription: "Continent code of the region",
		},
		"country": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Country code of the region",
			MarkdownDescription: "Country code of the region",
		},
		"datacenter_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Display name of the datacenter hosting the region",
			MarkdownDescription: "Display name of the datacenter hosting the region",
		},
		"availability_zones": schema.ListAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			Description:         "Availability zones available in the region",
			MarkdownDescription: "Availability zones available in the region",
		},
		"services": schema.ListAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			Description:         "Available OpenStack services in the region",
			MarkdownDescription: "Available OpenStack services in the region",
		},
	}
}

// cloudRegionObjectAttrTypes returns the attribute types of a single region object, matching
// the attributes declared in cloudRegionDetailAttributes plus "name".
func cloudRegionObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               ovhtypes.TfStringType{},
		"status":             ovhtypes.TfStringType{},
		"continent":          ovhtypes.TfStringType{},
		"country":            ovhtypes.TfStringType{},
		"datacenter_name":    ovhtypes.TfStringType{},
		"availability_zones": types.ListType{ElemType: types.StringType},
		"services":           types.ListType{ElemType: types.StringType},
	}
}

// stringSliceToTfList converts a Go string slice into a Terraform list of strings. A nil slice
// becomes a null list, while an empty (but non-nil) slice becomes an empty list.
func stringSliceToTfList(ctx context.Context, in []string) (types.List, diag.Diagnostics) {
	if in == nil {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, in)
}

// MergeWith populates the singular data source model from an API region.
func (m *cloudRegionDataSourceModel) MergeWith(ctx context.Context, region *CloudRegion) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(region.Name)}
	m.Status = ovhtypes.TfStringValue{StringValue: types.StringValue(region.Status)}
	m.Continent = ovhtypes.TfStringValue{StringValue: types.StringValue(region.Continent)}
	m.Country = ovhtypes.TfStringValue{StringValue: types.StringValue(region.Country)}
	m.DatacenterName = ovhtypes.TfStringValue{StringValue: types.StringValue(region.DatacenterName)}

	availabilityZones, d := stringSliceToTfList(ctx, region.AvailabilityZones)
	diags.Append(d...)
	m.AvailabilityZones = availabilityZones

	services, d := stringSliceToTfList(ctx, region.Services)
	diags.Append(d...)
	m.Services = services

	return diags
}

// cloudRegionToObject converts an API region into a Terraform object value matching
// cloudRegionObjectAttrTypes.
func cloudRegionToObject(ctx context.Context, region CloudRegion) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	availabilityZones, d := stringSliceToTfList(ctx, region.AvailabilityZones)
	diags.Append(d...)

	services, d := stringSliceToTfList(ctx, region.Services)
	diags.Append(d...)

	obj, d := types.ObjectValue(cloudRegionObjectAttrTypes(), map[string]attr.Value{
		"name":               ovhtypes.TfStringValue{StringValue: types.StringValue(region.Name)},
		"status":             ovhtypes.TfStringValue{StringValue: types.StringValue(region.Status)},
		"continent":          ovhtypes.TfStringValue{StringValue: types.StringValue(region.Continent)},
		"country":            ovhtypes.TfStringValue{StringValue: types.StringValue(region.Country)},
		"datacenter_name":    ovhtypes.TfStringValue{StringValue: types.StringValue(region.DatacenterName)},
		"availability_zones": availabilityZones,
		"services":           services,
	})
	diags.Append(d...)

	return obj, diags
}
