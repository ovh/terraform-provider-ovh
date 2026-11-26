package types

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func iso(s string) TfISO8601DurationValue {
	return TfISO8601DurationValue{TfStringValue{StringValue: basetypes.NewStringValue(s)}}
}

func TestTfISO8601DurationValue_StringSemanticEquals(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"PT0S", "P0D", true},
		{"P0D", "PT0S", true},
		{"PT0S", "PT0S", true},
		{"PT30S", "PT30S", true},
		{"PT0S", "PT30S", false},
		{"PT1M", "PT60S", true},
		{"PT120S", "PT2M", true},
		{"PT900S", "PT15M", true},
		{"PT1M", "PT61S", false},
		{"bogus", "PT0S", false},
	}
	for _, c := range cases {
		got, diags := iso(c.a).StringSemanticEquals(context.Background(), iso(c.b))
		if diags.HasError() || got != c.want {
			t.Errorf("%q ~ %q: got %v, want %v (diags: %v)", c.a, c.b, got, c.want, diags)
		}
	}
}

func TestTfISO8601DurationValue_Equal(t *testing.T) {
	if !iso("PT30S").Equal(iso("PT30S")) {
		t.Error("same spelling must be Equal")
	}
	if iso("PT0S").Equal(iso("P0D")) {
		t.Error("Equal must stay strict; only semantic equality tolerates zero spellings")
	}
	if iso("PT30S").Equal(basetypes.NewStringValue("PT30S")) {
		t.Error("plain StringValue must not be Equal to a duration value")
	}
}

func TestTfISO8601DurationType_ValueFromTerraform(t *testing.T) {
	ctx := context.Background()
	v, err := TfISO8601DurationType{}.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.String, "PT5M"))
	if err != nil {
		t.Fatal(err)
	}
	d, ok := v.(TfISO8601DurationValue)
	if !ok || d.ValueString() != "PT5M" {
		t.Fatalf("got %#v, want duration value PT5M", v)
	}
	if !d.Type(ctx).Equal(TfISO8601DurationType{}) {
		t.Error("value must report its own type")
	}
}
