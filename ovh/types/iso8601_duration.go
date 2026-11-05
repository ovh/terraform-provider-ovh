package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/peterhellberg/duration"
)

// TfISO8601DurationType is a string holding an ISO 8601 duration whose equivalent spellings
// are semantically equal, as the API normalizes them (PT120S is returned as PT2M).
type TfISO8601DurationType struct {
	TfStringType
}

var _ basetypes.StringTypable = TfISO8601DurationType{}

func (t TfISO8601DurationType) Equal(o attr.Type) bool {
	_, ok := o.(TfISO8601DurationType)
	return ok
}

func (t TfISO8601DurationType) String() string {
	return "TfISO8601DurationType"
}

func (t TfISO8601DurationType) ValueType(context.Context) attr.Value {
	return TfISO8601DurationValue{}
}

func (t TfISO8601DurationType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return TfISO8601DurationValue{TfStringValue{StringValue: in}}, nil
}

func (t TfISO8601DurationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}
	return TfISO8601DurationValue{TfStringValue{StringValue: stringValue}}, nil
}

type TfISO8601DurationValue struct {
	TfStringValue
}

var _ basetypes.StringValuableWithSemanticEquals = TfISO8601DurationValue{}

func (v TfISO8601DurationValue) Type(context.Context) attr.Type {
	return TfISO8601DurationType{}
}

func (v TfISO8601DurationValue) Equal(o attr.Value) bool {
	other, ok := o.(TfISO8601DurationValue)
	return ok && v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals compares the parsed durations, keeping the prior spelling.
func (v TfISO8601DurationValue) StringSemanticEquals(_ context.Context, o basetypes.StringValuable) (bool, diag.Diagnostics) {
	other, ok := o.(TfISO8601DurationValue)
	if !ok {
		return false, nil
	}
	if v.Equal(other) {
		return true, nil
	}
	a, errA := duration.Parse(v.ValueString())
	b, errB := duration.Parse(other.ValueString())
	return errA == nil && errB == nil && a == b, nil
}
