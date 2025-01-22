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

type TfStringValue struct {
	basetypes.StringValue
}

var _ basetypes.StringValuable = TfStringValue{}

func (t *TfStringValue) UnmarshalJSON(data []byte) error {
	var v *string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		t.StringValue = basetypes.NewStringNull()
	} else {
		t.StringValue = basetypes.NewStringValue(*v)
	}

	return nil
}

func (t TfStringValue) MarshalJSON() ([]byte, error) {
	if t.IsNull() || t.IsUnknown() {
		return []byte("null"), nil
	}
	return json.Marshal(t.StringValue.ValueString())
}

func (v TfStringValue) Equal(o attr.Value) bool {
	other, ok := o.(TfStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v TfStringValue) Type(ctx context.Context) attr.Type {
	return TfStringType{}
}

func NewTfStringValue(value string) TfStringValue {
	return TfStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewTfStringNull() TfStringValue {
	return TfStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

type TfStringType struct {
	basetypes.StringType
}

var _ basetypes.StringTypable = TfStringType{}

func (t TfStringType) Equal(o attr.Type) bool {
	other, ok := o.(TfStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t TfStringType) String() string {
	return "TfStringType"
}

func (t TfStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := TfStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t TfStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t TfStringType) ValueType(ctx context.Context) attr.Value {
	return TfStringValue{}
}
