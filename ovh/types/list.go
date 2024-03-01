// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// tfListNestedType is the attribute type of a TfListNestedValue.
type tfListNestedType[T attr.Value] struct {
	basetypes.ListType
}

var (
	_ basetypes.ListTypable  = (*tfListNestedType[attr.Value])(nil)
	_ basetypes.ListValuable = (*TfListNestedValue[attr.Value])(nil)
)

func NewTfListNestedType[T attr.Value](ctx context.Context) tfListNestedType[T] {
	var zero T
	return tfListNestedType[T]{basetypes.ListType{ElemType: zero.Type(ctx)}}
}

func (t tfListNestedType[T]) Equal(o attr.Type) bool {
	other, ok := o.(tfListNestedType[T])

	if !ok {
		return false
	}

	return t.ListType.Equal(other.ListType)
}

func (t tfListNestedType[T]) String() string {
	var zero T
	return fmt.Sprintf("TfListNestedType[%T]", zero)
}

func (t tfListNestedType[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewListNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	var zero T
	listValue, d := basetypes.NewListValue(zero.Type(ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	value := TfListNestedValue[T]{
		ListValue: listValue,
	}

	return value, diags
}

func (t tfListNestedType[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewListNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	var zero T
	listValue, d := basetypes.NewListValue(zero.Type(ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	value := TfListNestedValue[T]{
		ListValue: listValue,
	}

	return value, diags
}

func (t tfListNestedType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func (t tfListNestedType[T]) ValueType(ctx context.Context) attr.Value {
	return TfListNestedValue[T]{}
}

func (t tfListNestedType[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	return NewListNestedObjectValueOfNull[T](ctx), diags
}

// TfListNestedValue represents a Terraform Plugin Framework List value`.
type TfListNestedValue[T attr.Value] struct {
	basetypes.ListValue
}

var _ basetypes.SetValuable = (*TfListNestedValue[attr.Value])(nil)

func (t *TfListNestedValue[T]) UnmarshalJSON(data []byte) error {
	var v []T
	var d []attr.Value

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, c := range v {
		d = append(d, c)
	}

	var zero T
	t.ListValue = basetypes.NewListValueMust(zero.Type(context.Background()), d)

	return nil
}

func (t TfListNestedValue[T]) MarshalJSON() ([]byte, error) {
	var toMarshal []T

	for _, elem := range t.Elements() {
		toMarshal = append(toMarshal, elem.(T))
	}

	return json.Marshal(toMarshal)
}

func (t TfListNestedValue[T]) ToSetValue(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	return basetypes.NewSetValueMust(t.ElementType(ctx), t.Elements()), nil
}

func (v TfListNestedValue[T]) ToListValue(ctx context.Context) (basetypes.ListValue, diag.Diagnostics) {
	return basetypes.NewListValue(v.ElementType(ctx), v.Elements())
}

func (t TfListNestedValue[T]) ElementType(ctx context.Context) attr.Type {
	var zero T
	return zero.Type(ctx)
}

func (t TfListNestedValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	var zero T
	listType := tftypes.List{ElementType: zero.Type(ctx).TerraformType(ctx)}

	switch {
	case t.IsNull():
		return tftypes.NewValue(listType, nil), nil
	case t.IsUnknown():
		return tftypes.NewValue(listType, tftypes.UnknownValue), nil
	default:
		vals := make([]tftypes.Value, 0, len(t.Elements()))

		for _, elem := range t.Elements() {
			val, err := elem.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(listType, tftypes.UnknownValue), err
			}

			vals = append(vals, val)
		}

		if err := tftypes.ValidateValue(listType, vals); err != nil {
			return tftypes.NewValue(listType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(listType, vals), nil
	}
}

func (v TfListNestedValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(TfListNestedValue[T])

	if !ok {
		return false
	}

	return v.ListValue.Equal(other.ListValue)
}

func (v TfListNestedValue[T]) Type(ctx context.Context) attr.Type {
	return NewTfListNestedType[T](ctx)
}

func NewListNestedObjectValueOfNull[T attr.Value](ctx context.Context) TfListNestedValue[T] {
	var zero T
	return TfListNestedValue[T]{ListValue: basetypes.NewListNull(zero.Type(ctx))}
}

func NewListNestedObjectValueOfUnknown[T attr.Value](ctx context.Context) TfListNestedValue[T] {
	var zero T
	return TfListNestedValue[T]{ListValue: basetypes.NewListUnknown(zero.Type(ctx))}
}
