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
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

func OkmsServiceKeyPemDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := OkmsServiceKeyAttributes(ctx)
	attrs["okms_id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Okms ID",
		MarkdownDescription: "Okms ID",
	}
	attrs["keys_pem"] = schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"pem": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "The key in base64 encoded PEM format",
					MarkdownDescription: "The key in base64 encoded PEM format",
				},
			},
			CustomType: PemType{
				ObjectType: types.ObjectType{
					AttrTypes: PemValue{}.AttributeTypes(ctx),
				},
			},
		},
		Computed:            true,
		Description:         "The keys in PEM format",
		MarkdownDescription: "The keys in PEM format",
		CustomType:          ovhtypes.NewTfListNestedType[PemValue](ctx),
	}

	appendIamSchema(attrs, ctx)
	return schema.Schema{
		Attributes:  attrs,
		Description: "Use this data source to retrieve information about a KMS service key, in the PEM format.",
	}
}

type OkmsServiceKeyPemModel struct {
	CreatedAt  ovhtypes.TfStringValue                             `tfsdk:"created_at" json:"createdAt"`
	Curve      ovhtypes.TfStringValue                             `tfsdk:"curve" json:"curve"`
	Iam        IamValue                                           `tfsdk:"iam" json:"iam"`
	Id         ovhtypes.TfStringValue                             `tfsdk:"id" json:"id"`
	Keys       ovhtypes.TfListNestedValue[PemValue]               `tfsdk:"keys_pem" json:"keysPEM"`
	Name       ovhtypes.TfStringValue                             `tfsdk:"name" json:"name"`
	OkmsId     ovhtypes.TfStringValue                             `tfsdk:"okms_id" json:"okmsId"`
	Operations ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"operations" json:"operations"`
	Size       ovhtypes.TfInt64Value                              `tfsdk:"size" json:"size"`
	State      ovhtypes.TfStringValue                             `tfsdk:"state" json:"state"`
	Type       ovhtypes.TfStringValue                             `tfsdk:"type" json:"type"`
}

func (v *OkmsServiceKeyPemModel) MergeWith(other *OkmsServiceKeyPemModel) {
	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}

	if (v.Curve.IsUnknown() || v.Curve.IsNull()) && !other.Curve.IsUnknown() {
		v.Curve = other.Curve
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Keys.IsUnknown() || v.Keys.IsNull()) && !other.Keys.IsUnknown() {
		v.Keys = other.Keys
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.OkmsId.IsUnknown() || v.OkmsId.IsNull()) && !other.OkmsId.IsUnknown() {
		v.OkmsId = other.OkmsId
	}

	if (v.Operations.IsUnknown() || v.Operations.IsNull()) && !other.Operations.IsUnknown() {
		v.Operations = other.Operations
	}

	if (v.Size.IsUnknown() || v.Size.IsNull()) && !other.Size.IsUnknown() {
		v.Size = other.Size
	}

	if (v.State.IsUnknown() || v.State.IsNull()) && !other.State.IsUnknown() {
		v.State = other.State
	}

	if (v.Type.IsUnknown() || v.Type.IsNull()) && !other.Type.IsUnknown() {
		v.Type = other.Type
	}
}

var _ basetypes.ObjectTypable = PemType{}

type PemType struct {
	basetypes.ObjectType
}

func (t PemType) Equal(o attr.Type) bool {
	other, ok := o.(PemType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PemType) String() string {
	return "PemType"
}

func (t PemType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	pemAttribute, ok := attributes["pem"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`pem is missing from object`)

		return nil, diags
	}

	pemVal, ok := pemAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`pem expected to be ovhtypes.TfStringValue, was: %T`, pemAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return PemValue{
		Pem:   pemVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewPemValueNull() PemValue {
	return PemValue{
		state: attr.ValueStateNull,
	}
}

func NewPemValueUnknown() PemValue {
	return PemValue{
		state: attr.ValueStateUnknown,
	}
}

func NewPemValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PemValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PemValue Attribute Value",
				"While creating a PemValue value, a missing attribute value was detected. "+
					"A PemValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PemValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PemValue Attribute Type",
				"While creating a PemValue value, an invalid attribute value was detected. "+
					"A PemValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PemValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PemValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PemValue Attribute Value",
				"While creating a PemValue value, an extra attribute value was detected. "+
					"A PemValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PemValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPemValueUnknown(), diags
	}

	pemAttribute, ok := attributes["pem"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`pem is missing from object`)

		return NewPemValueUnknown(), diags
	}

	pemVal, ok := pemAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`pem expected to be ovhtypes.TfStringValue, was: %T`, pemAttribute))
	}

	if diags.HasError() {
		return NewPemValueUnknown(), diags
	}

	return PemValue{
		Pem:   pemVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewPemValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PemValue {
	object, diags := NewPemValue(attributeTypes, attributes)

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

		panic("NewPemValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PemType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPemValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPemValueUnknown(), nil
	}

	if in.IsNull() {
		return NewPemValueNull(), nil
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

	return NewPemValueMust(PemValue{}.AttributeTypes(ctx), attributes), nil
}

func (t PemType) ValueType(ctx context.Context) attr.Value {
	return PemValue{}
}

var _ basetypes.ObjectValuable = PemValue{}

type PemValue struct {
	Pem   ovhtypes.TfStringValue `tfsdk:"pem" json:"pem"`
	state attr.ValueState
}

func (v *PemValue) UnmarshalJSON(data []byte) error {
	type JsonPemValue PemValue

	var tmp JsonPemValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	v.Pem = tmp.Pem
	v.state = attr.ValueStateKnown

	return nil
}

func (v *PemValue) MergeWith(other *PemValue) {
	if (v.Pem.IsUnknown() || v.Pem.IsNull()) && !other.Pem.IsUnknown() {
		v.Pem = other.Pem
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v PemValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"pem": v.Pem,
	}
}
func (v PemValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 17)

	var val tftypes.Value
	var err error

	attrTypes["pem"] = basetypes.StringType{}.TerraformType(ctx)
	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Pem.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["pem"] = val

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

func (v PemValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PemValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PemValue) String() string {
	return "PemValue"
}

func (v PemValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"pem": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"pem": v.Pem,
		})

	return objVal, diags
}

func (v PemValue) Equal(o attr.Value) bool {
	other, ok := o.(PemValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Pem.Equal(other.Pem) {
		return false
	}

	return true
}

func (v PemValue) Type(ctx context.Context) attr.Type {
	return PemType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PemValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"pem": ovhtypes.TfStringType{},
	}
}
