// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ovh

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type okmsRegionStringValue struct {
	basetypes.StringValue
}

var _ basetypes.StringValuable = okmsRegionStringValue{}

func (t *okmsRegionStringValue) UnmarshalJSON(data []byte) error {
	var v *string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		t.StringValue = basetypes.NewStringNull()
	} else {
		// Convert values received from the API to the new region format
		t.StringValue = basetypes.NewStringValue(okmsReformatRegion(*v))
	}

	return nil
}

func (t okmsRegionStringValue) MarshalJSON() ([]byte, error) {
	if t.IsNull() || t.IsUnknown() {
		return []byte("null"), nil
	}
	return json.Marshal(t.StringValue.ValueString())
}

func (v okmsRegionStringValue) Equal(o attr.Value) bool {
	other, ok := o.(okmsRegionStringValue)

	if !ok {
		return false
	}

	if v.StringValue.IsUnknown() && !other.StringValue.IsUnknown() {
		return false
	}

	if v.StringValue.IsNull() && !other.StringValue.IsNull() {
		return false
	}

	if v.IsUnknown() || v.IsNull() {
		return true
	}

	return okmsReformatRegion(v.StringValue.ValueString()) == okmsReformatRegion(other.StringValue.ValueString())
}

func (v okmsRegionStringValue) Type(ctx context.Context) attr.Type {
	return okmsRegionStringType{}
}

func NewokmsRegionStringValue(value string) okmsRegionStringValue {
	return okmsRegionStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

type okmsRegionStringType struct {
	basetypes.StringType
}

var _ basetypes.StringTypable = okmsRegionStringType{}

func (t okmsRegionStringType) Equal(o attr.Type) bool {
	other, ok := o.(okmsRegionStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t okmsRegionStringType) String() string {
	return "okmsRegionStringType"
}

func (t okmsRegionStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := okmsRegionStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t okmsRegionStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t okmsRegionStringType) ValueType(ctx context.Context) attr.Value {
	return okmsRegionStringValue{}
}
