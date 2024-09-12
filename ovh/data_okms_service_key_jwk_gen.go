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

func OkmsServiceKeyJwkDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation time of the key",
			MarkdownDescription: "Creation time of the key",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Key ID",
			MarkdownDescription: "Key ID",
		},
		"keys": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"alg": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The algorithm intended to be used with the key",
						MarkdownDescription: "The algorithm intended to be used with the key",
					},
					"crv": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The cryptographic curve used with the key",
						MarkdownDescription: "The cryptographic curve used with the key",
					},
					"e": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The exponent value for the RSA public key",
						MarkdownDescription: "The exponent value for the RSA public key",
					},
					"key_ops": schema.ListAttribute{
						CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
						Computed:            true,
						Description:         "The operation for which the key is intended to be used",
						MarkdownDescription: "The operation for which the key is intended to be used",
					},
					"kid": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "key ID parameter used to match a specific key",
						MarkdownDescription: "key ID parameter used to match a specific key",
					},
					"kty": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Key type parameter identifies the cryptographic algorithm family used with the key, such as RSA or EC",
						MarkdownDescription: "Key type parameter identifies the cryptographic algorithm family used with the key, such as RSA or EC",
					},
					"n": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The modulus value for the RSA public key",
						MarkdownDescription: "The modulus value for the RSA public key",
					},
					"use": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The intended use of the public key",
						MarkdownDescription: "The intended use of the public key",
					},
					"x": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The x coordinate for the Elliptic Curve point",
						MarkdownDescription: "The x coordinate for the Elliptic Curve point",
					},
					"y": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "The y coordinate for the Elliptic Curve point",
						MarkdownDescription: "The y coordinate for the Elliptic Curve point",
					},
				},
				CustomType: JwkType{
					ObjectType: types.ObjectType{
						AttrTypes: JwkValue{}.AttributeTypes(ctx),
					},
				},
			},
			Computed:            true,
			Description:         "The key in JWK format",
			MarkdownDescription: "The key in JWK format",
			CustomType:          ovhtypes.NewTfListNestedType[JwkValue](ctx),
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key name",
			MarkdownDescription: "Key name",
		},
		"okms_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Okms ID",
			MarkdownDescription: "Okms ID",
		},
		"size": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Size of the key",
			MarkdownDescription: "Size of the key",
		},
		"state": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "State of the key",
			MarkdownDescription: "State of the key",
		},
		"type": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key type",
			MarkdownDescription: "Key type",
		},
	}

	return schema.Schema{
		Attributes:  attrs,
		Description: "Use this data source to retrieve information about a KMS service key, in the JWK format.",
	}
}

type OkmsServiceKeyJwkModel struct {
	CreatedAt ovhtypes.TfStringValue               `tfsdk:"created_at" json:"createdAt"`
	Id        ovhtypes.TfStringValue               `tfsdk:"id" json:"id"`
	Keys      ovhtypes.TfListNestedValue[JwkValue] `tfsdk:"keys" json:"keys"`
	Name      ovhtypes.TfStringValue               `tfsdk:"name" json:"name"`
	OkmsId    ovhtypes.TfStringValue               `tfsdk:"okms_id" json:"okmsId"`
	Size      ovhtypes.TfInt64Value                `tfsdk:"size" json:"size"`
	State     ovhtypes.TfStringValue               `tfsdk:"state" json:"state"`
	Type      ovhtypes.TfStringValue               `tfsdk:"type" json:"type"`
}

func (v *OkmsServiceKeyJwkModel) MergeWith(other *OkmsServiceKeyJwkModel) {

	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
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

var _ basetypes.ObjectTypable = JwkType{}

type JwkType struct {
	basetypes.ObjectType
}

func (t JwkType) Equal(o attr.Type) bool {
	other, ok := o.(JwkType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t JwkType) String() string {
	return "JwkType"
}

func (t JwkType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	algAttribute, ok := attributes["alg"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`alg is missing from object`)

		return nil, diags
	}

	algVal, ok := algAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`alg expected to be ovhtypes.TfStringValue, was: %T`, algAttribute))
	}

	crvAttribute, ok := attributes["crv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`crv is missing from object`)

		return nil, diags
	}

	crvVal, ok := crvAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`crv expected to be ovhtypes.TfStringValue, was: %T`, crvAttribute))
	}

	eAttribute, ok := attributes["e"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`e is missing from object`)

		return nil, diags
	}

	eVal, ok := eAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`e expected to be ovhtypes.TfStringValue, was: %T`, eAttribute))
	}

	keyOpsAttribute, ok := attributes["key_ops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key_ops is missing from object`)

		return nil, diags
	}

	keyOpsVal, ok := keyOpsAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key_ops expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, keyOpsAttribute))
	}

	kidAttribute, ok := attributes["kid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kid is missing from object`)

		return nil, diags
	}

	kidVal, ok := kidAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kid expected to be ovhtypes.TfStringValue, was: %T`, kidAttribute))
	}

	ktyAttribute, ok := attributes["kty"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kty is missing from object`)

		return nil, diags
	}

	ktyVal, ok := ktyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kty expected to be ovhtypes.TfStringValue, was: %T`, ktyAttribute))
	}

	nAttribute, ok := attributes["n"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`n is missing from object`)

		return nil, diags
	}

	nVal, ok := nAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`n expected to be ovhtypes.TfStringValue, was: %T`, nAttribute))
	}

	useAttribute, ok := attributes["use"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`use is missing from object`)

		return nil, diags
	}

	useVal, ok := useAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`use expected to be ovhtypes.TfStringValue, was: %T`, useAttribute))
	}

	xAttribute, ok := attributes["x"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`x is missing from object`)

		return nil, diags
	}

	xVal, ok := xAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`x expected to be ovhtypes.TfStringValue, was: %T`, xAttribute))
	}

	yAttribute, ok := attributes["y"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`y is missing from object`)

		return nil, diags
	}

	yVal, ok := yAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`y expected to be ovhtypes.TfStringValue, was: %T`, yAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return JwkValue{
		Alg:    algVal,
		Crv:    crvVal,
		E:      eVal,
		KeyOps: keyOpsVal,
		Kid:    kidVal,
		Kty:    ktyVal,
		N:      nVal,
		Use:    useVal,
		X:      xVal,
		Y:      yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewJwkValueNull() JwkValue {
	return JwkValue{
		state: attr.ValueStateNull,
	}
}

func NewJwkValueUnknown() JwkValue {
	return JwkValue{
		state: attr.ValueStateUnknown,
	}
}

func NewJwkValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (JwkValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing JwkValue Attribute Value",
				"While creating a JwkValue value, a missing attribute value was detected. "+
					"A JwkValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid JwkValue Attribute Type",
				"While creating a JwkValue value, an invalid attribute value was detected. "+
					"A JwkValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("JwkValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra JwkValue Attribute Value",
				"While creating a JwkValue value, an extra attribute value was detected. "+
					"A JwkValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra JwkValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewJwkValueUnknown(), diags
	}

	algAttribute, ok := attributes["alg"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`alg is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	algVal, ok := algAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`alg expected to be ovhtypes.TfStringValue, was: %T`, algAttribute))
	}

	crvAttribute, ok := attributes["crv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`crv is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	crvVal, ok := crvAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`crv expected to be ovhtypes.TfStringValue, was: %T`, crvAttribute))
	}

	eAttribute, ok := attributes["e"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`e is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	eVal, ok := eAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`e expected to be ovhtypes.TfStringValue, was: %T`, eAttribute))
	}

	keyOpsAttribute, ok := attributes["key_ops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key_ops is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	keyOpsVal, ok := keyOpsAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key_ops expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, keyOpsAttribute))
	}

	kidAttribute, ok := attributes["kid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kid is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	kidVal, ok := kidAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kid expected to be ovhtypes.TfStringValue, was: %T`, kidAttribute))
	}

	ktyAttribute, ok := attributes["kty"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kty is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	ktyVal, ok := ktyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kty expected to be ovhtypes.TfStringValue, was: %T`, ktyAttribute))
	}

	nAttribute, ok := attributes["n"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`n is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	nVal, ok := nAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`n expected to be ovhtypes.TfStringValue, was: %T`, nAttribute))
	}

	useAttribute, ok := attributes["use"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`use is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	useVal, ok := useAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`use expected to be ovhtypes.TfStringValue, was: %T`, useAttribute))
	}

	xAttribute, ok := attributes["x"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`x is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	xVal, ok := xAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`x expected to be ovhtypes.TfStringValue, was: %T`, xAttribute))
	}

	yAttribute, ok := attributes["y"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`y is missing from object`)

		return NewJwkValueUnknown(), diags
	}

	yVal, ok := yAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`y expected to be ovhtypes.TfStringValue, was: %T`, yAttribute))
	}

	if diags.HasError() {
		return NewJwkValueUnknown(), diags
	}

	return JwkValue{
		Alg:    algVal,
		Crv:    crvVal,
		E:      eVal,
		KeyOps: keyOpsVal,
		Kid:    kidVal,
		Kty:    ktyVal,
		N:      nVal,
		Use:    useVal,
		X:      xVal,
		Y:      yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewJwkValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) JwkValue {
	object, diags := NewJwkValue(attributeTypes, attributes)

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

		panic("NewJwkValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t JwkType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewJwkValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewJwkValueUnknown(), nil
	}

	if in.IsNull() {
		return NewJwkValueNull(), nil
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

	return NewJwkValueMust(JwkValue{}.AttributeTypes(ctx), attributes), nil
}

func (t JwkType) ValueType(ctx context.Context) attr.Value {
	return JwkValue{}
}

var _ basetypes.ObjectValuable = JwkValue{}

type JwkValue struct {
	Alg    ovhtypes.TfStringValue                             `tfsdk:"alg" json:"alg"`
	Crv    ovhtypes.TfStringValue                             `tfsdk:"crv" json:"crv"`
	E      ovhtypes.TfStringValue                             `tfsdk:"e" json:"e"`
	KeyOps ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"key_ops" json:"key_ops"`
	Kid    ovhtypes.TfStringValue                             `tfsdk:"kid" json:"kid"`
	Kty    ovhtypes.TfStringValue                             `tfsdk:"kty" json:"kty"`
	N      ovhtypes.TfStringValue                             `tfsdk:"n" json:"n"`
	Use    ovhtypes.TfStringValue                             `tfsdk:"use" json:"use"`
	X      ovhtypes.TfStringValue                             `tfsdk:"x" json:"x"`
	Y      ovhtypes.TfStringValue                             `tfsdk:"y" json:"y"`
	state  attr.ValueState
}

func (v *JwkValue) UnmarshalJSON(data []byte) error {
	type JsonJwkValue JwkValue

	var tmp JsonJwkValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Alg = tmp.Alg
	v.Crv = tmp.Crv
	v.E = tmp.E
	v.KeyOps = tmp.KeyOps
	v.Kid = tmp.Kid
	v.Kty = tmp.Kty
	v.N = tmp.N
	v.Use = tmp.Use
	v.X = tmp.X
	v.Y = tmp.Y

	v.state = attr.ValueStateKnown

	return nil
}

func (v *JwkValue) MergeWith(other *JwkValue) {

	if (v.Alg.IsUnknown() || v.Alg.IsNull()) && !other.Alg.IsUnknown() {
		v.Alg = other.Alg
	}

	if (v.Crv.IsUnknown() || v.Crv.IsNull()) && !other.Crv.IsUnknown() {
		v.Crv = other.Crv
	}

	if (v.E.IsUnknown() || v.E.IsNull()) && !other.E.IsUnknown() {
		v.E = other.E
	}

	if (v.KeyOps.IsUnknown() || v.KeyOps.IsNull()) && !other.KeyOps.IsUnknown() {
		v.KeyOps = other.KeyOps
	}

	if (v.Kid.IsUnknown() || v.Kid.IsNull()) && !other.Kid.IsUnknown() {
		v.Kid = other.Kid
	}

	if (v.Kty.IsUnknown() || v.Kty.IsNull()) && !other.Kty.IsUnknown() {
		v.Kty = other.Kty
	}

	if (v.N.IsUnknown() || v.N.IsNull()) && !other.N.IsUnknown() {
		v.N = other.N
	}

	if (v.Use.IsUnknown() || v.Use.IsNull()) && !other.Use.IsUnknown() {
		v.Use = other.Use
	}

	if (v.X.IsUnknown() || v.X.IsNull()) && !other.X.IsUnknown() {
		v.X = other.X
	}

	if (v.Y.IsUnknown() || v.Y.IsNull()) && !other.Y.IsUnknown() {
		v.Y = other.Y
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v JwkValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"alg":    v.Alg,
		"crv":    v.Crv,
		"e":      v.E,
		"keyOps": v.KeyOps,
		"kid":    v.Kid,
		"kty":    v.Kty,
		"n":      v.N,
		"use":    v.Use,
		"x":      v.X,
		"y":      v.Y,
	}
}
func (v JwkValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 17)

	var val tftypes.Value
	var err error

	attrTypes["alg"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["crv"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["e"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["key_ops"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["kid"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["kty"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["n"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["use"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["x"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["y"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 17)

		val, err = v.Alg.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["alg"] = val

		val, err = v.Crv.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["crv"] = val

		val, err = v.E.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["e"] = val

		val, err = v.KeyOps.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["key_ops"] = val

		val, err = v.Kid.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["kid"] = val

		val, err = v.Kty.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["kty"] = val

		val, err = v.N.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["n"] = val

		val, err = v.Use.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["use"] = val

		val, err = v.X.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["x"] = val

		val, err = v.Y.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["y"] = val

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

func (v JwkValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v JwkValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v JwkValue) String() string {
	return "JwkValue"
}

func (v JwkValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"alg":     ovhtypes.TfStringType{},
			"crv":     ovhtypes.TfStringType{},
			"e":       ovhtypes.TfStringType{},
			"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			"kid":     ovhtypes.TfStringType{},
			"kty":     ovhtypes.TfStringType{},
			"n":       ovhtypes.TfStringType{},
			"use":     ovhtypes.TfStringType{},
			"x":       ovhtypes.TfStringType{},
			"y":       ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"alg":     v.Alg,
			"crv":     v.Crv,
			"e":       v.E,
			"key_ops": v.KeyOps,
			"kid":     v.Kid,
			"kty":     v.Kty,
			"n":       v.N,
			"use":     v.Use,
			"x":       v.X,
			"y":       v.Y,
		})

	return objVal, diags
}

func (v JwkValue) Equal(o attr.Value) bool {
	other, ok := o.(JwkValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Alg.Equal(other.Alg) {
		return false
	}

	if !v.Crv.Equal(other.Crv) {
		return false
	}

	if !v.E.Equal(other.E) {
		return false
	}

	if !v.KeyOps.Equal(other.KeyOps) {
		return false
	}

	if !v.Kid.Equal(other.Kid) {
		return false
	}

	if !v.Kty.Equal(other.Kty) {
		return false
	}

	if !v.N.Equal(other.N) {
		return false
	}

	if !v.Use.Equal(other.Use) {
		return false
	}

	if !v.X.Equal(other.X) {
		return false
	}

	if !v.Y.Equal(other.Y) {
		return false
	}

	return true
}

func (v JwkValue) Type(ctx context.Context) attr.Type {
	return JwkType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v JwkValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"alg":     ovhtypes.TfStringType{},
		"crv":     ovhtypes.TfStringType{},
		"e":       ovhtypes.TfStringType{},
		"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
		"kid":     ovhtypes.TfStringType{},
		"kty":     ovhtypes.TfStringType{},
		"n":       ovhtypes.TfStringType{},
		"use":     ovhtypes.TfStringType{},
		"x":       ovhtypes.TfStringType{},
		"y":       ovhtypes.TfStringType{},
	}
}
