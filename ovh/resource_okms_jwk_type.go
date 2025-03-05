package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ basetypes.ObjectTypable = JwkFullType{}

type JwkFullType struct {
	basetypes.ObjectType
}

func (t JwkFullType) Equal(o attr.Type) bool {
	other, ok := o.(JwkFullType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t JwkFullType) String() string {
	return "JwkFullType"
}

func (t JwkFullType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	dAttribute, ok := attributes["d"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`d is missing from object`)

		return nil, diags
	}

	dVal, ok := dAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`d expected to be ovhtypes.TfStringValue, was: %T`, dAttribute))
	}

	dpAttribute, ok := attributes["dp"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dp is missing from object`)

		return nil, diags
	}

	dpVal, ok := dpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dp expected to be ovhtypes.TfStringValue, was: %T`, dpAttribute))
	}

	dqAttribute, ok := attributes["dq"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dq is missing from object`)

		return nil, diags
	}

	dqVal, ok := dqAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dq expected to be ovhtypes.TfStringValue, was: %T`, dqAttribute))
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

	kAttribute, ok := attributes["k"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`k is missing from object`)

		return nil, diags
	}

	kVal, ok := kAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`k expected to be ovhtypes.TfStringValue, was: %T`, kAttribute))
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

	pAttribute, ok := attributes["p"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`p is missing from object`)

		return nil, diags
	}

	pVal, ok := pAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`p expected to be ovhtypes.TfStringValue, was: %T`, pAttribute))
	}

	qAttribute, ok := attributes["q"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`q is missing from object`)

		return nil, diags
	}

	qVal, ok := qAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`q expected to be ovhtypes.TfStringValue, was: %T`, qAttribute))
	}

	qiAttribute, ok := attributes["qi"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`qi is missing from object`)

		return nil, diags
	}

	qiVal, ok := qiAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`qi expected to be ovhtypes.TfStringValue, was: %T`, qiAttribute))
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

	return JwkFullValue{
		Alg:    algVal,
		Crv:    crvVal,
		D:      dVal,
		Dp:     dpVal,
		Dq:     dqVal,
		E:      eVal,
		K:      kVal,
		KeyOps: keyOpsVal,
		Kid:    kidVal,
		Kty:    ktyVal,
		N:      nVal,
		P:      pVal,
		Q:      qVal,
		Qi:     qiVal,
		Use:    useVal,
		X:      xVal,
		Y:      yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func (t JwkFullType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewJwkFullValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewJwkFullValueUnknown(), nil
	}

	if in.IsNull() {
		return NewJwkFullValueNull(), nil
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

	return NewJwkFullValueMust(JwkFullValue{}.AttributeTypes(ctx), attributes), nil
}

func (t JwkFullType) ValueType(ctx context.Context) attr.Value {
	return JwkFullValue{}
}

var _ basetypes.ObjectTypable = JwkFullWritableType{}

type JwkFullWritableType struct {
	basetypes.ObjectType
}

func (t JwkFullWritableType) Equal(o attr.Type) bool {
	other, ok := o.(JwkFullWritableType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t JwkFullWritableType) String() string {
	return "JwkFullWritableType"
}

func (t JwkFullWritableType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	dAttribute, ok := attributes["d"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`d is missing from object`)

		return nil, diags
	}

	dVal, ok := dAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`d expected to be ovhtypes.TfStringValue, was: %T`, dAttribute))
	}

	dpAttribute, ok := attributes["dp"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dp is missing from object`)

		return nil, diags
	}

	dpVal, ok := dpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dp expected to be ovhtypes.TfStringValue, was: %T`, dpAttribute))
	}

	dqAttribute, ok := attributes["dq"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dq is missing from object`)

		return nil, diags
	}

	dqVal, ok := dqAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dq expected to be ovhtypes.TfStringValue, was: %T`, dqAttribute))
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

	kAttribute, ok := attributes["k"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`k is missing from object`)

		return nil, diags
	}

	kVal, ok := kAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`k expected to be ovhtypes.TfStringValue, was: %T`, kAttribute))
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

	pAttribute, ok := attributes["p"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`p is missing from object`)

		return nil, diags
	}

	pVal, ok := pAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`p expected to be ovhtypes.TfStringValue, was: %T`, pAttribute))
	}

	qAttribute, ok := attributes["q"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`q is missing from object`)

		return nil, diags
	}

	qVal, ok := qAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`q expected to be ovhtypes.TfStringValue, was: %T`, qAttribute))
	}

	qiAttribute, ok := attributes["qi"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`qi is missing from object`)

		return nil, diags
	}

	qiVal, ok := qiAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`qi expected to be ovhtypes.TfStringValue, was: %T`, qiAttribute))
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

	return JwkFullWritableValue{
		Alg:    &algVal,
		Crv:    &crvVal,
		D:      &dVal,
		Dp:     &dpVal,
		Dq:     &dqVal,
		E:      &eVal,
		K:      &kVal,
		KeyOps: &keyOpsVal,
		Kid:    &kidVal,
		Kty:    &ktyVal,
		N:      &nVal,
		P:      &pVal,
		Q:      &qVal,
		Qi:     &qiVal,
		Use:    &useVal,
		X:      &xVal,
		Y:      &yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func (t JwkFullWritableType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return JwkFullWritableValue{
			state: attr.ValueStateNull,
		}, nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return JwkFullWritableValue{
			state: attr.ValueStateUnknown,
		}, nil
	}

	if in.IsNull() {
		return JwkFullWritableValue{
			state: attr.ValueStateNull,
		}, nil
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

	return NewJwkFullWritableValueMust(JwkFullWritableValue{}.AttributeTypes(ctx), attributes), nil
}

func (t JwkFullWritableType) ValueType(ctx context.Context) attr.Value {
	return JwkFullValue{}
}

var _ basetypes.ObjectValuable = JwkFullValue{}

type JwkFullValue struct {
	Alg    ovhtypes.TfStringValue                             `tfsdk:"alg" json:"alg"`
	Crv    ovhtypes.TfStringValue                             `tfsdk:"crv" json:"crv"`
	D      ovhtypes.TfStringValue                             `tfsdk:"d" json:"d"`
	Dp     ovhtypes.TfStringValue                             `tfsdk:"dp" json:"dp"`
	Dq     ovhtypes.TfStringValue                             `tfsdk:"dq" json:"dq"`
	E      ovhtypes.TfStringValue                             `tfsdk:"e" json:"e"`
	K      ovhtypes.TfStringValue                             `tfsdk:"k" json:"k"`
	KeyOps ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"key_ops" json:"key_ops"`
	Kid    ovhtypes.TfStringValue                             `tfsdk:"kid" json:"kid"`
	Kty    ovhtypes.TfStringValue                             `tfsdk:"kty" json:"kty"`
	N      ovhtypes.TfStringValue                             `tfsdk:"n" json:"n"`
	P      ovhtypes.TfStringValue                             `tfsdk:"p" json:"p"`
	Q      ovhtypes.TfStringValue                             `tfsdk:"q" json:"q"`
	Qi     ovhtypes.TfStringValue                             `tfsdk:"qi" json:"qi"`
	Use    ovhtypes.TfStringValue                             `tfsdk:"use" json:"use"`
	X      ovhtypes.TfStringValue                             `tfsdk:"x" json:"x"`
	Y      ovhtypes.TfStringValue                             `tfsdk:"y" json:"y"`
	state  attr.ValueState
}

func (v JwkFullValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"alg":    v.Alg,
		"crv":    v.Crv,
		"d":      v.D,
		"dp":     v.Dp,
		"dq":     v.Dq,
		"e":      v.E,
		"k":      v.K,
		"keyOps": v.KeyOps,
		"kid":    v.Kid,
		"kty":    v.Kty,
		"n":      v.N,
		"p":      v.P,
		"q":      v.Q,
		"qi":     v.Qi,
		"use":    v.Use,
		"x":      v.X,
		"y":      v.Y,
	}
}

func (v JwkFullValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"alg":     ovhtypes.TfStringType{},
		"crv":     ovhtypes.TfStringType{},
		"d":       ovhtypes.TfStringType{},
		"dp":      ovhtypes.TfStringType{},
		"dq":      ovhtypes.TfStringType{},
		"e":       ovhtypes.TfStringType{},
		"k":       ovhtypes.TfStringType{},
		"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
		"kid":     ovhtypes.TfStringType{},
		"kty":     ovhtypes.TfStringType{},
		"n":       ovhtypes.TfStringType{},
		"p":       ovhtypes.TfStringType{},
		"q":       ovhtypes.TfStringType{},
		"qi":      ovhtypes.TfStringType{},
		"use":     ovhtypes.TfStringType{},
		"x":       ovhtypes.TfStringType{},
		"y":       ovhtypes.TfStringType{},
	}
}

func (v JwkFullValue) Equal(o attr.Value) bool {
	other, ok := o.(JwkFullValue)

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

	if !v.D.Equal(other.D) {
		return false
	}

	if !v.Dp.Equal(other.Dp) {
		return false
	}

	if !v.Dq.Equal(other.Dq) {
		return false
	}

	if !v.E.Equal(other.E) {
		return false
	}

	if !v.K.Equal(other.K) {
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

	if !v.P.Equal(other.P) {
		return false
	}

	if !v.Q.Equal(other.Q) {
		return false
	}

	if !v.Qi.Equal(other.Qi) {
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

func (v JwkFullValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v JwkFullValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v *JwkFullValue) MergeWith(other *JwkFullValue) {

	if (v.Alg.IsUnknown() || v.Alg.IsNull()) && !other.Alg.IsUnknown() {
		v.Alg = other.Alg
	}

	if (v.Crv.IsUnknown() || v.Crv.IsNull()) && !other.Crv.IsUnknown() {
		v.Crv = other.Crv
	}

	if (v.D.IsUnknown() || v.D.IsNull()) && !other.D.IsUnknown() {
		v.D = other.D
	}

	if (v.Dp.IsUnknown() || v.Dp.IsNull()) && !other.Dp.IsUnknown() {
		v.Dp = other.Dp
	}

	if (v.Dq.IsUnknown() || v.Dq.IsNull()) && !other.Dq.IsUnknown() {
		v.Dq = other.Dq
	}

	if (v.E.IsUnknown() || v.E.IsNull()) && !other.E.IsUnknown() {
		v.E = other.E
	}

	if (v.K.IsUnknown() || v.K.IsNull()) && !other.K.IsUnknown() {
		v.K = other.K
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

	if (v.P.IsUnknown() || v.P.IsNull()) && !other.P.IsUnknown() {
		v.P = other.P
	}

	if (v.Q.IsUnknown() || v.Q.IsNull()) && !other.Q.IsUnknown() {
		v.Q = other.Q
	}

	if (v.Qi.IsUnknown() || v.Qi.IsNull()) && !other.Qi.IsUnknown() {
		v.Qi = other.Qi
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

func (v JwkFullValue) String() string {
	return "JwkFullValue"
}

func (v JwkFullValue) ToCreate() *JwkFullWritableValue {
	res := &JwkFullWritableValue{
		state: v.state,
	}

	if !v.N.IsUnknown() {
		res.N = &v.N
	}

	if !v.Use.IsUnknown() {
		res.Use = &v.Use
	}

	if !v.X.IsUnknown() {
		res.X = &v.X
	}

	if !v.Y.IsUnknown() {
		res.Y = &v.Y
	}

	if !v.Crv.IsUnknown() {
		res.Crv = &v.Crv
	}

	if !v.P.IsUnknown() {
		res.P = &v.P
	}

	if !v.Dp.IsUnknown() {
		res.Dp = &v.Dp
	}

	// Kid should be readonly, but right now it's required
	// send an empty string while the API isn't fixed
	var placeholder = ovhtypes.NewTfStringValue("tf-kid-placeholder")
	res.Kid = &placeholder
	if !v.Kid.IsUnknown() {
		res.Kid = &v.Kid
		panic("Kid should be readonly")
	}

	if !v.Qi.IsUnknown() {
		res.Qi = &v.Qi
	}

	if !v.Kty.IsUnknown() {
		res.Kty = &v.Kty
	}

	if !v.Alg.IsUnknown() {
		res.Alg = &v.Alg
	}
	log.Printf("to create alg %v, source alg %v", res.Alg, v.Alg)

	if !v.K.IsUnknown() {
		res.K = &v.K
	}

	if !v.Q.IsUnknown() {
		res.Q = &v.Q
	}

	if !v.Dq.IsUnknown() {
		res.Dq = &v.Dq
	}

	if !v.D.IsUnknown() {
		res.D = &v.D
	}

	if !v.E.IsUnknown() {
		res.E = &v.E
	}

	if !v.KeyOps.IsUnknown() {
		res.KeyOps = &v.KeyOps
	}

	return res
}

func (v JwkFullValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"alg":     ovhtypes.TfStringType{},
			"crv":     ovhtypes.TfStringType{},
			"d":       ovhtypes.TfStringType{},
			"dp":      ovhtypes.TfStringType{},
			"dq":      ovhtypes.TfStringType{},
			"e":       ovhtypes.TfStringType{},
			"k":       ovhtypes.TfStringType{},
			"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			"kid":     ovhtypes.TfStringType{},
			"kty":     ovhtypes.TfStringType{},
			"n":       ovhtypes.TfStringType{},
			"p":       ovhtypes.TfStringType{},
			"q":       ovhtypes.TfStringType{},
			"qi":      ovhtypes.TfStringType{},
			"use":     ovhtypes.TfStringType{},
			"x":       ovhtypes.TfStringType{},
			"y":       ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"alg":     v.Alg,
			"crv":     v.Crv,
			"d":       v.D,
			"dp":      v.Dp,
			"dq":      v.Dq,
			"e":       v.E,
			"k":       v.K,
			"key_ops": v.KeyOps,
			"kid":     v.Kid,
			"kty":     v.Kty,
			"n":       v.N,
			"p":       v.P,
			"q":       v.Q,
			"qi":      v.Qi,
			"use":     v.Use,
			"x":       v.X,
			"y":       v.Y,
		})

	return objVal, diags
}

func (v JwkFullValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 17)

	var val tftypes.Value
	var err error

	attrTypes["alg"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["crv"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["d"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dp"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dq"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["e"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["k"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["key_ops"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["kid"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["kty"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["n"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["p"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["q"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["qi"] = basetypes.StringType{}.TerraformType(ctx)
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

		val, err = v.D.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["d"] = val

		val, err = v.Dp.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dp"] = val

		val, err = v.Dq.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dq"] = val

		val, err = v.E.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["e"] = val

		val, err = v.K.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["k"] = val

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

		val, err = v.P.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["p"] = val

		val, err = v.Q.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["q"] = val

		val, err = v.Qi.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["qi"] = val

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

func (v JwkFullValue) Type(ctx context.Context) attr.Type {
	return JwkFullType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v *JwkFullValue) UnmarshalJSON(data []byte) error {
	type JsonJwkFullValue JwkFullValue

	var tmp JsonJwkFullValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Alg = tmp.Alg
	v.Crv = tmp.Crv
	v.D = tmp.D
	v.Dp = tmp.Dp
	v.Dq = tmp.Dq
	v.E = tmp.E
	v.K = tmp.K
	v.KeyOps = tmp.KeyOps
	v.Kid = tmp.Kid
	v.Kty = tmp.Kty
	v.N = tmp.N
	v.P = tmp.P
	v.Q = tmp.Q
	v.Qi = tmp.Qi
	v.Use = tmp.Use
	v.X = tmp.X
	v.Y = tmp.Y

	v.state = attr.ValueStateKnown

	return nil
}

var _ basetypes.ObjectValuable = JwkFullWritableValue{}

type JwkFullWritableValue struct {
	*JwkFullValue `json:"-"`
	Alg           *ovhtypes.TfStringValue                             `tfsdk:"alg" json:"alg,omitempty"`
	Crv           *ovhtypes.TfStringValue                             `tfsdk:"crv" json:"crv,omitempty"`
	D             *ovhtypes.TfStringValue                             `tfsdk:"d" json:"d,omitempty"`
	Dp            *ovhtypes.TfStringValue                             `tfsdk:"dp" json:"dp,omitempty"`
	Dq            *ovhtypes.TfStringValue                             `tfsdk:"dq" json:"dq,omitempty"`
	E             *ovhtypes.TfStringValue                             `tfsdk:"e" json:"e,omitempty"`
	K             *ovhtypes.TfStringValue                             `tfsdk:"k" json:"k,omitempty"`
	KeyOps        *ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"key_ops" json:"key_ops,omitempty"`
	Kid           *ovhtypes.TfStringValue                             `tfsdk:"kid" json:"kid,omitempty"`
	Kty           *ovhtypes.TfStringValue                             `tfsdk:"kty" json:"kty,omitempty"`
	N             *ovhtypes.TfStringValue                             `tfsdk:"n" json:"n,omitempty"`
	P             *ovhtypes.TfStringValue                             `tfsdk:"p" json:"p,omitempty"`
	Q             *ovhtypes.TfStringValue                             `tfsdk:"q" json:"q,omitempty"`
	Qi            *ovhtypes.TfStringValue                             `tfsdk:"qi" json:"qi,omitempty"`
	Use           *ovhtypes.TfStringValue                             `tfsdk:"use" json:"use,omitempty"`
	X             *ovhtypes.TfStringValue                             `tfsdk:"x" json:"x,omitempty"`
	Y             *ovhtypes.TfStringValue                             `tfsdk:"y" json:"y,omitempty"`
	state         attr.ValueState
}

func (v JwkFullWritableValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"alg":    v.Alg,
		"crv":    v.Crv,
		"d":      v.D,
		"dp":     v.Dp,
		"dq":     v.Dq,
		"e":      v.E,
		"k":      v.K,
		"keyOps": v.KeyOps,
		"kid":    v.Kid,
		"kty":    v.Kty,
		"n":      v.N,
		"p":      v.P,
		"q":      v.Q,
		"qi":     v.Qi,
		"use":    v.Use,
		"x":      v.X,
		"y":      v.Y,
	}
}

func (v JwkFullWritableValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"alg":     ovhtypes.TfStringType{},
		"crv":     ovhtypes.TfStringType{},
		"d":       ovhtypes.TfStringType{},
		"dp":      ovhtypes.TfStringType{},
		"dq":      ovhtypes.TfStringType{},
		"e":       ovhtypes.TfStringType{},
		"k":       ovhtypes.TfStringType{},
		"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
		"kid":     ovhtypes.TfStringType{},
		"kty":     ovhtypes.TfStringType{},
		"n":       ovhtypes.TfStringType{},
		"p":       ovhtypes.TfStringType{},
		"q":       ovhtypes.TfStringType{},
		"qi":      ovhtypes.TfStringType{},
		"use":     ovhtypes.TfStringType{},
		"x":       ovhtypes.TfStringType{},
		"y":       ovhtypes.TfStringType{},
	}
}

func (v JwkFullWritableValue) Equal(o attr.Value) bool {
	other, ok := o.(JwkFullWritableValue)

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

	if !v.D.Equal(other.D) {
		return false
	}

	if !v.Dp.Equal(other.Dp) {
		return false
	}

	if !v.Dq.Equal(other.Dq) {
		return false
	}

	if !v.E.Equal(other.E) {
		return false
	}

	if !v.K.Equal(other.K) {
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

	if !v.P.Equal(other.P) {
		return false
	}

	if !v.Q.Equal(other.Q) {
		return false
	}

	if !v.Qi.Equal(other.Qi) {
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

func (v JwkFullWritableValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v JwkFullWritableValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v JwkFullWritableValue) String() string {
	return "JwkFullWritableValue"
}

func (v JwkFullWritableValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"alg":     ovhtypes.TfStringType{},
			"crv":     ovhtypes.TfStringType{},
			"d":       ovhtypes.TfStringType{},
			"dp":      ovhtypes.TfStringType{},
			"dq":      ovhtypes.TfStringType{},
			"e":       ovhtypes.TfStringType{},
			"k":       ovhtypes.TfStringType{},
			"key_ops": ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			"kid":     ovhtypes.TfStringType{},
			"kty":     ovhtypes.TfStringType{},
			"n":       ovhtypes.TfStringType{},
			"p":       ovhtypes.TfStringType{},
			"q":       ovhtypes.TfStringType{},
			"qi":      ovhtypes.TfStringType{},
			"use":     ovhtypes.TfStringType{},
			"x":       ovhtypes.TfStringType{},
			"y":       ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"alg":     v.Alg,
			"crv":     v.Crv,
			"d":       v.D,
			"dp":      v.Dp,
			"dq":      v.Dq,
			"e":       v.E,
			"k":       v.K,
			"key_ops": v.KeyOps,
			"kid":     v.Kid,
			"kty":     v.Kty,
			"n":       v.N,
			"p":       v.P,
			"q":       v.Q,
			"qi":      v.Qi,
			"use":     v.Use,
			"x":       v.X,
			"y":       v.Y,
		})

	return objVal, diags
}

func (v JwkFullWritableValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 17)

	var val tftypes.Value
	var err error

	attrTypes["alg"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["crv"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["d"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dp"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dq"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["e"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["k"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["key_ops"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["kid"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["kty"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["n"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["p"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["q"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["qi"] = basetypes.StringType{}.TerraformType(ctx)
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

		val, err = v.D.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["d"] = val

		val, err = v.Dp.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dp"] = val

		val, err = v.Dq.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dq"] = val

		val, err = v.E.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["e"] = val

		val, err = v.K.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["k"] = val

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

		val, err = v.P.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["p"] = val

		val, err = v.Q.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["q"] = val

		val, err = v.Qi.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["qi"] = val

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

func (v JwkFullWritableValue) Type(ctx context.Context) attr.Type {
	return JwkFullWritableType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}
