package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestLoadbalancerPluralDataSourcesSchemaMatchesAttrTypes checks that the object
// type advertised by each plural data source schema is exactly the one used to
// build the list elements, so state conversion cannot fail at runtime.
func TestLoadbalancerPluralDataSourcesSchemaMatchesAttrTypes(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name      string
		ds        datasource.DataSource
		attribute string
		attrTypes map[string]attr.Type
	}{
		{"loadbalancers", NewCloudLoadbalancersDataSource(), "loadbalancers", LoadbalancerListItemAttrTypes()},
		{"listeners", NewCloudLoadbalancerListenersDataSource(), "listeners", ListenerListItemAttrTypes()},
		{"pools", NewCloudLoadbalancerPoolsDataSource(), "pools", PoolListItemAttrTypes()},
		{"members", NewCloudLoadbalancerPoolMembersDataSource(), "members", MemberListItemAttrTypes()},
		{"l7policies", NewCloudLoadbalancerL7PoliciesDataSource(), "l7policies", L7PolicyListItemAttrTypes()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var resp datasource.SchemaResponse
			c.ds.(interface {
				Schema(context.Context, datasource.SchemaRequest, *datasource.SchemaResponse)
			}).Schema(ctx, datasource.SchemaRequest{}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
			}

			att, ok := resp.Schema.Attributes[c.attribute]
			if !ok {
				t.Fatalf("attribute %q missing from schema", c.attribute)
			}

			listType, ok := att.GetType().(types.ListType)
			if !ok {
				t.Fatalf("attribute %q is not a ListType but %T", c.attribute, att.GetType())
			}

			want := types.ObjectType{AttrTypes: c.attrTypes}
			if !listType.ElemType.Equal(want) {
				t.Fatalf("element type mismatch for %q:\n schema: %v\n  built: %v", c.attribute, listType.ElemType, want)
			}
		})
	}
}
