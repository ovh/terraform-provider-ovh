// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type TfNumberType struct {
	basetypes.NumberType
}

var _ basetypes.NumberTypable = TfNumberType{}

func (t TfNumberType) Equal(o attr.Type) bool {
	other, ok := o.(TfNumberType)

	if !ok {
		return false
	}

	return t.NumberType.Equal(other.NumberType)
}

func (t TfNumberType) String() string {
	return "TfNumberType"
}

func (t TfNumberType) ValueFromNumber(ctx context.Context, in basetypes.NumberValue) (basetypes.NumberValuable, diag.Diagnostics) {
	return TfNumberValue{
		NumberValue: in,
	}, nil
}

func (t TfNumberType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.NumberType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	numberValue, ok := attrValue.(basetypes.NumberValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	numberValuable, diags := t.ValueFromNumber(ctx, numberValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting NumberValue to NumberValuable: %v", diags)
	}

	return numberValuable, nil
}

func (t TfNumberType) ValueType(ctx context.Context) attr.Value {
	return TfNumberValue{}
}

type TfNumberValue struct {
	basetypes.NumberValue
}

var _ basetypes.NumberValuable = TfNumberValue{}

func (t *TfNumberValue) UnmarshalJSON(data []byte) error {
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.NumberValue = basetypes.NewNumberValue(big.NewFloat(v))

	return nil
}

func (t TfNumberValue) MarshalJSON() ([]byte, error) {
	floatVal, _ := t.ValueBigFloat().Float64()
	return json.Marshal(floatVal)
}

func (t TfNumberValue) ToNumberValue(ctx context.Context) (basetypes.NumberValue, diag.Diagnostics) {
	return t.NumberValue, nil
}

func (v TfNumberValue) Type(ctx context.Context) attr.Type {
	return TfNumberType{}
}
