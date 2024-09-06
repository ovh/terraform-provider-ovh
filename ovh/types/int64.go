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

type TfInt64Type struct {
	basetypes.Int64Type
}

var _ basetypes.Int64Typable = TfInt64Type{}

func (t TfInt64Type) Equal(o attr.Type) bool {
	other, ok := o.(TfInt64Type)

	if !ok {
		return false
	}

	return t.Int64Type.Equal(other.Int64Type)
}

func (t TfInt64Type) String() string {
	return "TfInt64Type"
}

func (t TfInt64Type) ValueFromInt64(ctx context.Context, in basetypes.Int64Value) (basetypes.Int64Valuable, diag.Diagnostics) {
	value := TfInt64Value{
		Int64Value: in,
	}

	return value, nil
}

func (t TfInt64Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.Int64Type.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	intValue, ok := attrValue.(basetypes.Int64Value)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	intValuable, diags := t.ValueFromInt64(ctx, intValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting Int64Value to Int64Valuable: %v", diags)
	}

	return intValuable, nil
}

func (t TfInt64Type) ValueType(ctx context.Context) attr.Value {
	return TfInt64Value{}
}

type TfInt64Value struct {
	basetypes.Int64Value
}

var _ basetypes.Int64Valuable = TfInt64Value{}

func (t *TfInt64Value) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.Int64Value = basetypes.NewInt64Value(v)

	return nil
}

func (t TfInt64Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Int64Value.ValueInt64())
}

func (v TfInt64Value) Type(ctx context.Context) attr.Type {
	return TfInt64Type{}
}

func NewTfInt64Value(value int64) TfInt64Value {
	return TfInt64Value{
		Int64Value: basetypes.NewInt64Value(value),
	}
}

func NewTfInt64ValueNull() TfInt64Value {
	return TfInt64Value{
		Int64Value: basetypes.NewInt64Null(),
	}
}
