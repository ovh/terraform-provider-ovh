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

type TfBoolValue struct {
	basetypes.BoolValue
}

var _ basetypes.BoolValuable = TfBoolValue{}

func (t *TfBoolValue) UnmarshalJSON(data []byte) error {
	var v *bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		t.BoolValue = basetypes.NewBoolNull()
	} else {
		t.BoolValue = basetypes.NewBoolValue(*v)
	}

	return nil
}

func (t TfBoolValue) MarshalJSON() ([]byte, error) {
	if t.IsNull() || t.IsUnknown() {
		return []byte("null"), nil
	}
	return json.Marshal(t.BoolValue.ValueBool())
}

func (v TfBoolValue) Equal(o attr.Value) bool {
	other, ok := o.(TfBoolValue)

	if !ok {
		return false
	}

	return v.BoolValue.Equal(other.BoolValue)
}

func (v TfBoolValue) Type(ctx context.Context) attr.Type {
	return TfBoolType{}
}

func NewTfBoolValue(v bool) TfBoolValue {
	return TfBoolValue{
		BoolValue: basetypes.NewBoolValue(v),
	}
}

type TfBoolType struct {
	basetypes.BoolType
}

var _ basetypes.BoolTypable = TfBoolType{}

func (t TfBoolType) Equal(o attr.Type) bool {
	other, ok := o.(TfBoolType)

	if !ok {
		return false
	}

	return t.BoolType.Equal(other.BoolType)
}

func (t TfBoolType) String() string {
	return "TfBoolType"
}

func (t TfBoolType) ValueFromBool(ctx context.Context, in basetypes.BoolValue) (basetypes.BoolValuable, diag.Diagnostics) {
	value := TfBoolValue{
		BoolValue: in,
	}

	return value, nil
}

func (t TfBoolType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.BoolType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	boolValue, ok := attrValue.(basetypes.BoolValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	boolValuable, diags := t.ValueFromBool(ctx, boolValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting BoolValue to BoolValuable: %v", diags)
	}

	return boolValuable, nil
}

func (t TfBoolType) ValueType(ctx context.Context) attr.Value {
	return TfBoolValue{}
}
