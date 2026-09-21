package ovh

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Resource and data source share one model per loadbalancer object. A model
// field typed with one of the ovh custom types (TfStringValue,
// TfListNestedValue[…]) only decodes when the schema declares that very type:
// a plain ElementType/no CustomType leaves basetypes.ListType in the schema and
// every read then fails with "Value Conversion Error" — which is what the
// listener data source did on allowed_cidrs. Walk both schemas of every
// loadbalancer object and pin each attribute to its model field's type.
func TestCloudLoadbalancerSchemasMatchModelTypes(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name  string
		model any
		ds    datasource.DataSource
		res   resource.Resource
	}{
		{"loadbalancer", CloudLoadbalancerModel{}, &cloudLoadbalancerDataSource{}, &cloudLoadbalancerResource{}},
		{"listener", CloudLoadbalancerListenerModel{}, &cloudLoadbalancerListenerDataSource{}, &cloudLoadbalancerListenerResource{}},
		{"pool", CloudLoadbalancerPoolModel{}, &cloudLoadbalancerPoolDataSource{}, &cloudLoadbalancerPoolResource{}},
		{"pool_member", CloudLoadbalancerPoolMemberModel{}, &cloudLoadbalancerPoolMemberDataSource{}, &cloudLoadbalancerPoolMemberResource{}},
		{"l7policy", CloudLoadbalancerL7PolicyModel{}, &cloudLoadbalancerL7PolicyDataSource{}, &cloudLoadbalancerL7PolicyResource{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dsResp dsschema.Schema
			{
				var resp datasource.SchemaResponse
				tc.ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
				dsResp = resp.Schema
			}
			var rResp rschema.Schema
			{
				var resp resource.SchemaResponse
				tc.res.Schema(ctx, resource.SchemaRequest{}, &resp)
				rResp = resp.Schema
			}

			for name, want := range modelAttrTypes(ctx, tc.model) {
				if attr, ok := dsResp.Attributes[name]; !ok {
					t.Errorf("data source: missing attribute %q", name)
				} else if got := attr.GetType(); !got.Equal(want) {
					t.Errorf("data source %q: type is %s, model expects %s", name, got, want)
				}

				if attr, ok := rResp.Attributes[name]; !ok {
					t.Errorf("resource: missing attribute %q", name)
				} else if got := attr.GetType(); !got.Equal(want) {
					t.Errorf("resource %q: type is %s, model expects %s", name, got, want)
				}
			}
		})
	}
}

// modelAttrTypes maps each tfsdk field of a model onto the type its zero value
// reports. Fields typed with a framework base type (types.Object, types.Int64…)
// are skipped: a zero types.Object carries no attribute types, so it tells us
// nothing about what the schema should declare, and it accepts whatever the
// schema does declare.
func modelAttrTypes(ctx context.Context, model any) map[string]attr.Type {
	out := map[string]attr.Type{}
	rt := reflect.TypeOf(model)

	for i := range rt.NumField() {
		field := rt.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" || tag == "-" {
			continue
		}
		if field.Type.PkgPath() != "github.com/ovh/terraform-provider-ovh/v2/ovh/types" {
			continue
		}
		value, ok := reflect.New(field.Type).Elem().Interface().(attr.Value)
		if !ok {
			continue
		}
		out[tag] = value.Type(ctx)
	}

	return out
}
