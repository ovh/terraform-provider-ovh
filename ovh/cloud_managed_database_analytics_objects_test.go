package ovh

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// managedDatabaseObjectTestState returns the state of r holding model, or a null state when model is
// nil. It also checks the schema of r.
func managedDatabaseObjectTestState(t *testing.T, r fwresource.Resource, model any) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if diags := schemaResp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("invalid schema: %v", diags)
	}
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
	if model != nil {
		if diags := state.Set(ctx, model); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	}
	return state
}

func TestManagedDatabaseCategories(t *testing.T) {
	if got := managedAnalyticsFamily.serviceEndpoint("pj", "kafka", "s1"); got != "/v2/publicCloud/project/pj/managedAnalytics/kafka/s1" {
		t.Errorf("unexpected endpoint %s", got)
	}
}

// Every resource and data source of the families has a valid schema, a name of its own, and an engine
// argument, if any, that takes engines of its family only.
func TestManagedDatabaseSchemas(t *testing.T) {
	ctx := context.Background()
	names := map[string]bool{}
	for _, new := range (&OvhProvider{}).Resources(ctx) {
		r := new()
		var meta fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "ovh"}, &meta)
		if !strings.HasPrefix(meta.TypeName, "ovh_cloud_managed_") {
			continue
		}
		var resp fwresource.SchemaResponse
		r.Schema(ctx, fwresource.SchemaRequest{}, &resp)
		if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() || names["resource "+meta.TypeName] {
			t.Errorf("resource %s: duplicate or invalid schema: %v", meta.TypeName, diags)
		}
		names["resource "+meta.TypeName] = true
		testManagedDatabaseEngineArgument(t, meta.TypeName, resp.Schema.Attributes["engine"])
	}
	for _, new := range (&OvhProvider{}).DataSources(ctx) {
		d := new()
		var meta datasource.MetadataResponse
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "ovh"}, &meta)
		if !strings.HasPrefix(meta.TypeName, "ovh_cloud_managed_") {
			continue
		}
		var resp datasource.SchemaResponse
		d.Schema(ctx, datasource.SchemaRequest{}, &resp)
		if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() || names["data "+meta.TypeName] {
			t.Errorf("data source %s: duplicate or invalid schema: %v", meta.TypeName, diags)
		}
		names["data "+meta.TypeName] = true
		testManagedDatabaseEngineArgument(t, "data "+meta.TypeName, resp.Schema.Attributes["engine"])
	}
	if len(names) != 2+4 {
		t.Errorf("expected 2 resources and 4 data sources, got %d", len(names))
	}
}

// testManagedDatabaseEngineArgument checks that an engine argument rejects an unknown engine and the
// engines of the other family, and takes at least one engine of the family of the type.
func testManagedDatabaseEngineArgument(t *testing.T, typeName string, engine any) {
	t.Helper()
	argument, ok := engine.(interface {
		IsRequired() bool
		StringValidators() []validator.String
	})
	if !ok || !argument.IsRequired() {
		return
	}
	family := managedDatabaseFamily
	if strings.Contains(typeName, "ovh"+managedAnalyticsFamily.typeName("")) {
		family = managedAnalyticsFamily
	}

	taken := 0
	for _, name := range append(append([]string{"unknown"}, managedDatabaseFamily.engines...), managedAnalyticsFamily.engines...) {
		var diags diag.Diagnostics
		for _, v := range argument.StringValidators() {
			var resp validator.StringResponse
			v.ValidateString(context.Background(), validator.StringRequest{Path: path.Root("engine"), ConfigValue: types.StringValue(name)}, &resp)
			diags.Append(resp.Diagnostics...)
		}
		if diags.HasError() {
			continue
		}
		taken++
		if !slices.Contains(family.engines, name) {
			t.Errorf("%s: takes the engine %s, outside the %s family", typeName, name, family.family)
		}
	}
	if taken == 0 {
		t.Errorf("%s: takes no engine", typeName)
	}
}

// A 403 on a read is an error, unlike a 404: every resource stays in the state.
func TestManagedDatabaseReadForbidden(t *testing.T) {
	ctx := context.Background()
	config, _ := newManagedDatabaseTestConfig(t, map[string][]managedDatabaseTestAnswer{
		"*": {{http.StatusForbidden, `{"class": "Client::Forbidden", "message": "This call has not been granted"}`}},
	})
	for _, new := range (&OvhProvider{}).Resources(ctx) {
		r := new()
		var meta fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "ovh"}, &meta)
		if !strings.HasPrefix(meta.TypeName, "ovh_cloud_managed_") {
			continue
		}
		r.(fwresource.ResourceWithConfigure).Configure(ctx, fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})

		state := managedDatabaseObjectTestState(t, r, nil)
		for name, value := range map[string]string{"service_name": "pj", "engine": "e", "service_id": "s1", "id": "o1"} {
			if _, ok := state.Schema.GetAttributes()[name]; ok {
				if diags := state.SetAttribute(ctx, path.Root(name), value); diags.HasError() {
					t.Fatalf("%s: unexpected diagnostics: %v", meta.TypeName, diags)
				}
			}
		}
		resp := fwresource.ReadResponse{State: state}
		r.Read(ctx, fwresource.ReadRequest{State: state}, &resp)
		if !resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
			t.Errorf("%s: expected an error and the resource kept in the state, got %v", meta.TypeName, resp.Diagnostics)
		}
	}
}

func TestManagedDatabaseStringFormats(t *testing.T) {
	for _, tc := range []struct {
		format managedDatabaseStringFormat
		value  string
		valid  bool
	}{
		{managedDatabaseDuration, "1h30m", true},
		{managedDatabaseDuration, "40 minutes", false},
	} {
		var resp validator.StringResponse
		tc.format.ValidateString(context.Background(), validator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(tc.value)}, &resp)
		if resp.Diagnostics.HasError() == tc.valid {
			t.Errorf("%s %q: expected valid %t, got %v", tc.format.name, tc.value, tc.valid, resp.Diagnostics)
		}
	}
}
