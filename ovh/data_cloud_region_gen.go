package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func CloudRegionDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"region_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Region name",
			MarkdownDescription: "Region name",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Region identifier",
			MarkdownDescription: "Region identifier",
		},
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Region status",
			MarkdownDescription: "Region status",
		},
		"continent": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Continent code",
			MarkdownDescription: "Continent code",
		},
		"country": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Country code",
			MarkdownDescription: "Country code",
		},
		"datacenter_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Datacenter display name",
			MarkdownDescription: "Datacenter display name",
		},
		"services": schema.ListAttribute{
			ElementType:         ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Available OpenStack services",
			MarkdownDescription: "Available OpenStack services",
		},
		"availability_zones": schema.ListAttribute{
			ElementType:         ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Availability zones in this region",
			MarkdownDescription: "Availability zones in this region",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudRegionModel struct {
	ServiceName       ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
	RegionName        ovhtypes.TfStringValue `tfsdk:"region_name" json:"-"`
	Name              ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	Status            ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Continent         ovhtypes.TfStringValue `tfsdk:"continent" json:"continent"`
	Country           ovhtypes.TfStringValue `tfsdk:"country" json:"country"`
	DatacenterName    ovhtypes.TfStringValue `tfsdk:"datacenter_name" json:"datacenterName"`
	Services          types.List             `tfsdk:"services" json:"services"`
	AvailabilityZones types.List             `tfsdk:"availability_zones" json:"availabilityZones"`
}

func (m *CloudRegionModel) UnmarshalJSON(data []byte) error {
	type raw struct {
		Name              string   `json:"name"`
		Status            string   `json:"status"`
		Continent         string   `json:"continent"`
		Country           string   `json:"country"`
		DatacenterName    string   `json:"datacenterName"`
		Services          []string `json:"services"`
		AvailabilityZones []string `json:"availabilityZones"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	m.Name = ovhtypes.NewTfStringValue(r.Name)
	m.Status = ovhtypes.NewTfStringValue(r.Status)
	m.Continent = ovhtypes.NewTfStringValue(r.Continent)
	m.Country = ovhtypes.NewTfStringValue(r.Country)
	m.DatacenterName = ovhtypes.NewTfStringValue(r.DatacenterName)

	svcElems := make([]attr.Value, len(r.Services))
	for i, s := range r.Services {
		svcElems[i] = ovhtypes.NewTfStringValue(s)
	}
	m.Services = types.ListValueMust(ovhtypes.TfStringType{}, svcElems)

	azElems := make([]attr.Value, len(r.AvailabilityZones))
	for i, az := range r.AvailabilityZones {
		azElems[i] = ovhtypes.NewTfStringValue(az)
	}
	m.AvailabilityZones = types.ListValueMust(ovhtypes.TfStringType{}, azElems)

	return nil
}

// --- CloudRegions (list datasource) ---

func CloudRegionsDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Filter regions by name (regexp match)",
			MarkdownDescription: "Filter regions by name (regexp match)",
		},
		"regions": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Region identifier",
						MarkdownDescription: "Region identifier",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Region status",
						MarkdownDescription: "Region status",
					},
					"continent": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Continent code",
						MarkdownDescription: "Continent code",
					},
					"country": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Country code",
						MarkdownDescription: "Country code",
					},
					"datacenter_name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Datacenter display name",
						MarkdownDescription: "Datacenter display name",
					},
					"services": schema.ListAttribute{
						ElementType:         ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Available OpenStack services",
						MarkdownDescription: "Available OpenStack services",
					},
					"availability_zones": schema.ListAttribute{
						ElementType:         ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Availability zones in this region",
						MarkdownDescription: "Availability zones in this region",
					},
				},
				CustomType: CloudRegionsValueType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudRegionsValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType: ovhtypes.NewTfListNestedType[CloudRegionsValue](ctx),
			Computed:   true,
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudRegionsModel struct {
	ServiceName ovhtypes.TfStringValue                        `tfsdk:"service_name" json:"-"`
	Name        ovhtypes.TfStringValue                        `tfsdk:"name" json:"-"`
	Regions     ovhtypes.TfListNestedValue[CloudRegionsValue] `tfsdk:"regions" json:"regions"`
}

// --- CloudRegionsValue (nested object for list items) ---

var _ basetypes.ObjectTypable = CloudRegionsValueType{}

type CloudRegionsValueType struct {
	basetypes.ObjectType
}

func (t CloudRegionsValueType) Equal(o attr.Type) bool {
	other, ok := o.(CloudRegionsValueType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudRegionsValueType) String() string {
	return "CloudRegionsValueType"
}

func (t CloudRegionsValueType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics
	attributes := in.Attributes()

	nameVal, ok := attributes["name"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, attributes["name"]))
		return nil, diags
	}
	statusVal, ok := attributes["status"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`status expected to be ovhtypes.TfStringValue, was: %T`, attributes["status"]))
		return nil, diags
	}
	continentVal, ok := attributes["continent"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`continent expected to be ovhtypes.TfStringValue, was: %T`, attributes["continent"]))
		return nil, diags
	}
	countryVal, ok := attributes["country"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`country expected to be ovhtypes.TfStringValue, was: %T`, attributes["country"]))
		return nil, diags
	}
	datacenterNameVal, ok := attributes["datacenter_name"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`datacenter_name expected to be ovhtypes.TfStringValue, was: %T`, attributes["datacenter_name"]))
		return nil, diags
	}
	servicesVal, ok := attributes["services"].(basetypes.ListValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`services expected to be basetypes.ListValue, was: %T`, attributes["services"]))
		return nil, diags
	}
	availabilityZonesVal, ok := attributes["availability_zones"].(basetypes.ListValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`availability_zones expected to be basetypes.ListValue, was: %T`, attributes["availability_zones"]))
		return nil, diags
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudRegionsValue{
		Name:              nameVal,
		Status:            statusVal,
		Continent:         continentVal,
		Country:           countryVal,
		DatacenterName:    datacenterNameVal,
		Services:          servicesVal,
		AvailabilityZones: availabilityZonesVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func (t CloudRegionsValueType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return CloudRegionsValue{state: attr.ValueStateNull}, nil
	}
	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return CloudRegionsValue{state: attr.ValueStateUnknown}, nil
	}
	if in.IsNull() {
		return CloudRegionsValue{state: attr.ValueStateNull}, nil
	}

	attributes := map[string]attr.Value{}
	val := map[string]tftypes.Value{}
	if err := in.As(&val); err != nil {
		return nil, err
	}
	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}
		attributes[k] = a
	}

	result, diags := NewCloudRegionsValue(CloudRegionsValue{}.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		diagsStrings := make([]string, 0, len(diags))
		for _, d := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf("%s | %s | %s", d.Severity(), d.Summary(), d.Detail()))
		}
		return nil, fmt.Errorf("error creating CloudRegionsValue: %s", strings.Join(diagsStrings, "\n"))
	}
	return result, nil
}

func (t CloudRegionsValueType) ValueType(ctx context.Context) attr.Value {
	return CloudRegionsValue{}
}

var _ basetypes.ObjectValuable = CloudRegionsValue{}

type CloudRegionsValue struct {
	Name              ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	Status            ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Continent         ovhtypes.TfStringValue `tfsdk:"continent" json:"continent"`
	Country           ovhtypes.TfStringValue `tfsdk:"country" json:"country"`
	DatacenterName    ovhtypes.TfStringValue `tfsdk:"datacenter_name" json:"datacenterName"`
	Services          basetypes.ListValue    `tfsdk:"services" json:"services"`
	AvailabilityZones basetypes.ListValue    `tfsdk:"availability_zones" json:"availabilityZones"`
	state             attr.ValueState
}

func (v *CloudRegionsValue) UnmarshalJSON(data []byte) error {
	type raw struct {
		Name              string   `json:"name"`
		Status            string   `json:"status"`
		Continent         string   `json:"continent"`
		Country           string   `json:"country"`
		DatacenterName    string   `json:"datacenterName"`
		Services          []string `json:"services"`
		AvailabilityZones []string `json:"availabilityZones"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	v.Name = ovhtypes.NewTfStringValue(r.Name)
	v.Status = ovhtypes.NewTfStringValue(r.Status)
	v.Continent = ovhtypes.NewTfStringValue(r.Continent)
	v.Country = ovhtypes.NewTfStringValue(r.Country)
	v.DatacenterName = ovhtypes.NewTfStringValue(r.DatacenterName)

	svcElems := make([]attr.Value, len(r.Services))
	for i, s := range r.Services {
		svcElems[i] = ovhtypes.NewTfStringValue(s)
	}
	v.Services = types.ListValueMust(ovhtypes.TfStringType{}, svcElems)

	azElems := make([]attr.Value, len(r.AvailabilityZones))
	for i, az := range r.AvailabilityZones {
		azElems[i] = ovhtypes.NewTfStringValue(az)
	}
	v.AvailabilityZones = types.ListValueMust(ovhtypes.TfStringType{}, azElems)

	v.state = attr.ValueStateKnown
	return nil
}

func NewCloudRegionsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudRegionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]
		if !ok {
			diags.AddError("Missing Attribute", fmt.Sprintf("CloudRegionsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()))
			continue
		}
		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError("Invalid Attribute Type", fmt.Sprintf("Expected: %s, Given: %s", attributeType.String(), attribute.Type(ctx)))
		}
	}
	if diags.HasError() {
		return CloudRegionsValue{state: attr.ValueStateUnknown}, diags
	}

	return CloudRegionsValue{
		Name:              attributes["name"].(ovhtypes.TfStringValue),
		Status:            attributes["status"].(ovhtypes.TfStringValue),
		Continent:         attributes["continent"].(ovhtypes.TfStringValue),
		Country:           attributes["country"].(ovhtypes.TfStringValue),
		DatacenterName:    attributes["datacenter_name"].(ovhtypes.TfStringValue),
		Services:          attributes["services"].(basetypes.ListValue),
		AvailabilityZones: attributes["availability_zones"].(basetypes.ListValue),
		state:             attr.ValueStateKnown,
	}, diags
}

func (v CloudRegionsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"name":               ovhtypes.TfStringType{},
		"status":             ovhtypes.TfStringType{},
		"continent":          ovhtypes.TfStringType{},
		"country":            ovhtypes.TfStringType{},
		"datacenter_name":    ovhtypes.TfStringType{},
		"services":           types.ListType{ElemType: ovhtypes.TfStringType{}},
		"availability_zones": types.ListType{ElemType: ovhtypes.TfStringType{}},
	}
}

func (v CloudRegionsValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"name":               v.Name,
		"status":             v.Status,
		"continent":          v.Continent,
		"country":            v.Country,
		"datacenter_name":    v.DatacenterName,
		"services":           v.Services,
		"availability_zones": v.AvailabilityZones,
	}
}

func (v CloudRegionsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 7)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["status"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["continent"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["country"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["datacenter_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["services"] = basetypes.ListType{ElemType: basetypes.StringType{}}.TerraformType(ctx)
	attrTypes["availability_zones"] = basetypes.ListType{ElemType: basetypes.StringType{}}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 7)
		var val tftypes.Value
		var err error

		val, err = v.Name.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["name"] = val

		val, err = v.Status.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["status"] = val

		val, err = v.Continent.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["continent"] = val

		val, err = v.Country.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["country"] = val

		val, err = v.DatacenterName.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["datacenter_name"] = val

		val, err = v.Services.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["services"] = val

		val, err = v.AvailabilityZones.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["availability_zones"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CloudRegionsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudRegionsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudRegionsValue) String() string {
	return "CloudRegionsValue"
}

func (v CloudRegionsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), v.Attributes())
}

func (v CloudRegionsValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudRegionsValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Name.Equal(other.Name) && v.Status.Equal(other.Status) && v.Continent.Equal(other.Continent) &&
		v.Country.Equal(other.Country) && v.DatacenterName.Equal(other.DatacenterName) &&
		v.Services.Equal(other.Services) && v.AvailabilityZones.Equal(other.AvailabilityZones)
}

func (v CloudRegionsValue) Type(ctx context.Context) attr.Type {
	return CloudRegionsValueType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}
