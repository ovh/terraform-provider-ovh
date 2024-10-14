// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.MapTypable  = tfMapNestedType[basetypes.StringValue]{}
	_ basetypes.MapValuable = TfMapNestedValue[basetypes.StringValue]{}
)

type tfMapNestedType[T attr.Value] struct {
	basetypes.MapType
}

func NewTfMapNestedType[T attr.Value](ctx context.Context) tfMapNestedType[T] {
	var zero T
	return tfMapNestedType[T]{basetypes.MapType{ElemType: zero.Type(ctx)}}
}

func (t tfMapNestedType[T]) Equal(o attr.Type) bool {
	other, ok := o.(tfMapNestedType[T])

	if !ok {
		return false
	}

	return t.MapType.Equal(other.MapType)
}

func (t tfMapNestedType[T]) String() string {
	var zero T
	return fmt.Sprintf("%T", zero)
}

func (t tfMapNestedType[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	var diags diag.Diagnostics
	var zero T

	if in.IsNull() {
		return NewNullTfMapNestedValue[T](ctx), diags
	}

	if in.IsUnknown() {
		return NewUnknownTfMapNestedValue[T](ctx), diags
	}

	mapValue, d := basetypes.NewMapValue(zero.Type(ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.NewMapUnknown(types.StringType), diags
	}

	value := TfMapNestedValue[T]{
		MapValue: mapValue,
	}

	return value, diags
}

func (t tfMapNestedType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.MapType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	mapValue, ok := attrValue.(basetypes.MapValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	mapValuable, diags := t.ValueFromMap(ctx, mapValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting MapValue to MapValuable: %v", diags)
	}

	return mapValuable, nil
}

func (t tfMapNestedType[T]) ValueType(ctx context.Context) attr.Value {
	return TfMapNestedValue[T]{}
}

type TfMapNestedValue[T attr.Value] struct {
	basetypes.MapValue
}

func (v *TfMapNestedValue[T]) UnmarshalJSON(data []byte) error {
	var mm map[string]T

	if err := json.Unmarshal(data, &mm); err != nil {
		return err
	}

	var zero T
	if mm == nil {
		v.MapValue = basetypes.NewMapNull(zero.Type(context.Background()))
	} else {
		d := make(map[string]attr.Value, len(mm))
		for k, v := range mm {
			d[k] = v
		}
		v.MapValue = basetypes.NewMapValueMust(zero.Type(context.Background()), d)
	}

	return nil
}

func (t TfMapNestedValue[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() || t.IsUnknown() {
		return []byte("null"), nil
	}

	elems := t.Elements()
	toMarshal := make(map[string]T, len(elems))
	for key, elem := range elems {
		toMarshal[key] = elem.(T)
	}

	return json.Marshal(toMarshal)
}

// ToTerraformValue returns the data contained in the Map as a tftypes.Value.
func (v TfMapNestedValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {

	if v.MapValue.IsNull() {
		var zero T
		mapType := tftypes.Map{ElementType: zero.Type(ctx).TerraformType(ctx)}
		return tftypes.NewValue(mapType, nil), nil
	}

	return v.MapValue.ToTerraformValue(ctx)
}

func (v TfMapNestedValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(TfMapNestedValue[T])

	if !ok {
		return false
	}

	return v.MapValue.Equal(other.MapValue)
}

func (v TfMapNestedValue[T]) Type(ctx context.Context) attr.Type {
	return NewTfMapNestedType[T](ctx)
}

func NewTfMapNestedValue[T attr.Value](ctx context.Context, elements map[string]attr.Value) (TfMapNestedValue[T], diag.Diagnostics) {
	var zero T

	mapValue, diags := basetypes.NewMapValue(zero.Type(ctx), elements)
	if diags.HasError() {
		return NewUnknownTfMapNestedValue[T](ctx), diags
	}

	return TfMapNestedValue[T]{MapValue: mapValue}, diags
}

func NewNullTfMapNestedValue[T attr.Value](ctx context.Context) TfMapNestedValue[T] {
	var zero T
	return TfMapNestedValue[T]{MapValue: basetypes.NewMapNull(zero.Type(ctx))}
}

func NewUnknownTfMapNestedValue[T attr.Value](ctx context.Context) TfMapNestedValue[T] {
	var zero T
	return TfMapNestedValue[T]{MapValue: basetypes.NewMapUnknown(zero.Type(ctx))}
}

func NewTfMapNestedValueMust[T attr.Value](ctx context.Context, elements map[string]attr.Value) TfMapNestedValue[T] {
	mapVal, _ := NewTfMapNestedValue[T](ctx, elements)
	return mapVal
}
