// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CloudProjectVolumesDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"cloud_project_volumes": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
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
					"size": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Volume size",
						MarkdownDescription: "Volume size",
					},
				},
				CustomType: CloudProjectVolumesType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudProjectVolumesValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType: ovhtypes.NewTfListNestedType[CloudProjectVolumesValue](ctx),
			Computed:   true,
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
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudProjectVolumesModel struct {
	CloudProjectVolumes ovhtypes.TfListNestedValue[CloudProjectVolumesValue] `tfsdk:"cloud_project_volumes" json:"cloudProjectVolumes"`
	RegionName          ovhtypes.TfStringValue                               `tfsdk:"region_name" json:"regionName"`
	ServiceName         ovhtypes.TfStringValue                               `tfsdk:"service_name" json:"serviceName"`
}

func (v *CloudProjectVolumesModel) MergeWith(other *CloudProjectVolumesModel) {

	if (v.CloudProjectVolumes.IsUnknown() || v.CloudProjectVolumes.IsNull()) && !other.CloudProjectVolumes.IsUnknown() {
		v.CloudProjectVolumes = other.CloudProjectVolumes
	}

	if (v.RegionName.IsUnknown() || v.RegionName.IsNull()) && !other.RegionName.IsUnknown() {
		v.RegionName = other.RegionName
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

}

var _ basetypes.ObjectTypable = CloudProjectVolumesType{}

type CloudProjectVolumesType struct {
	basetypes.ObjectType
}

func (t CloudProjectVolumesType) Equal(o attr.Type) bool {
	other, ok := o.(CloudProjectVolumesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudProjectVolumesType) String() string {
	return "CloudProjectVolumesType"
}

func (t CloudProjectVolumesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	sizeAttribute, ok := attributes["size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`size is missing from object`)

		return nil, diags
	}

	sizeVal, ok := sizeAttribute.(ovhtypes.TfInt64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`size expected to be ovhtypes.TfInt64Value, was: %T`, sizeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudProjectVolumesValue{
		Id:    idVal,
		Name:  nameVal,
		Size:  sizeVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCloudProjectVolumesValueNull() CloudProjectVolumesValue {
	return CloudProjectVolumesValue{
		state: attr.ValueStateNull,
	}
}

func NewCloudProjectVolumesValueUnknown() CloudProjectVolumesValue {
	return CloudProjectVolumesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCloudProjectVolumesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudProjectVolumesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CloudProjectVolumesValue Attribute Value",
				"While creating a CloudProjectVolumesValue value, a missing attribute value was detected. "+
					"A CloudProjectVolumesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudProjectVolumesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CloudProjectVolumesValue Attribute Type",
				"While creating a CloudProjectVolumesValue value, an invalid attribute value was detected. "+
					"A CloudProjectVolumesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudProjectVolumesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CloudProjectVolumesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CloudProjectVolumesValue Attribute Value",
				"While creating a CloudProjectVolumesValue value, an extra attribute value was detected. "+
					"A CloudProjectVolumesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CloudProjectVolumesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCloudProjectVolumesValueUnknown(), diags
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCloudProjectVolumesValueUnknown(), diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewCloudProjectVolumesValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	sizeAttribute, ok := attributes["size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`size is missing from object`)

		return NewCloudProjectVolumesValueUnknown(), diags
	}

	sizeVal, ok := sizeAttribute.(ovhtypes.TfInt64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`size expected to be ovhtypes.TfInt64Value, was: %T`, sizeAttribute))
	}

	if diags.HasError() {
		return NewCloudProjectVolumesValueUnknown(), diags
	}

	return CloudProjectVolumesValue{
		Id:    idVal,
		Name:  nameVal,
		Size:  sizeVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCloudProjectVolumesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CloudProjectVolumesValue {
	object, diags := NewCloudProjectVolumesValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCloudProjectVolumesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CloudProjectVolumesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCloudProjectVolumesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCloudProjectVolumesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCloudProjectVolumesValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCloudProjectVolumesValueMust(CloudProjectVolumesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CloudProjectVolumesType) ValueType(ctx context.Context) attr.Value {
	return CloudProjectVolumesValue{}
}

var _ basetypes.ObjectValuable = CloudProjectVolumesValue{}

type CloudProjectVolumesValue struct {
	Id    ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name  ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	Size  ovhtypes.TfInt64Value  `tfsdk:"size" json:"size"`
	state attr.ValueState
}

func (v *CloudProjectVolumesValue) UnmarshalJSON(data []byte) error {
	type JsonCloudProjectVolumesValue CloudProjectVolumesValue

	var tmp JsonCloudProjectVolumesValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Id = tmp.Id
	v.Name = tmp.Name
	v.Size = tmp.Size

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CloudProjectVolumesValue) MergeWith(other *CloudProjectVolumesValue) {

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.Size.IsUnknown() || v.Size.IsNull()) && !other.Size.IsUnknown() {
		v.Size = other.Size
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CloudProjectVolumesValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"id":   v.Id,
		"name": v.Name,
		"size": v.Size,
	}
}
func (v CloudProjectVolumesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["size"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.Size.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["size"] = val

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

func (v CloudProjectVolumesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudProjectVolumesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudProjectVolumesValue) String() string {
	return "CloudProjectVolumesValue"
}

func (v CloudProjectVolumesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"id":   ovhtypes.TfStringType{},
			"name": ovhtypes.TfStringType{},
			"size": ovhtypes.TfInt64Type{},
		},
		map[string]attr.Value{
			"id":   v.Id,
			"name": v.Name,
			"size": v.Size,
		})

	return objVal, diags
}

func (v CloudProjectVolumesValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudProjectVolumesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Size.Equal(other.Size) {
		return false
	}

	return true
}

func (v CloudProjectVolumesValue) Type(ctx context.Context) attr.Type {
	return CloudProjectVolumesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CloudProjectVolumesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"name": ovhtypes.TfStringType{},
		"size": ovhtypes.TfInt64Type{},
	}
}
