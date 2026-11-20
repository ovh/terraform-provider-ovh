package ovh

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

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

func managedDatabaseObjectTestPlan(t *testing.T, r fwresource.Resource, model any) tfsdk.Plan {
	t.Helper()
	state := managedDatabaseObjectTestState(t, r, model)
	return tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
}

func TestManagedDatabaseObjectWait(t *testing.T) {
	config, _ := newManagedDatabaseTestConfig(t, map[string][]managedDatabaseTestAnswer{
		// A 5xx or a 429 is one more read.
		"GET /v2/ready":   {{http.StatusOK, `{"status": "CREATING"}`}, {http.StatusServiceUnavailable, `{"message": "unavailable"}`}, {http.StatusOK, `{"status": "READY"}`}},
		"GET /v2/failed":  {{http.StatusOK, `{"status": "ERROR"}`}},
		"GET /v2/deleted": {{http.StatusOK, `{"status": "READY"}`}, {http.StatusOK, `{"status": "DELETING"}`}, {http.StatusNotFound, `{"message": "not found"}`}},
	})
	ctx := context.Background()

	if err := waitForManagedDatabaseObject(ctx, config.OVHClient, "/v2/ready", false, time.Minute); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := waitForManagedDatabaseObject(ctx, config.OVHClient, "/v2/failed", false, time.Minute); err == nil {
		t.Errorf("expected an error for an object in ERROR status")
	}
	if err := waitForManagedDatabaseObject(ctx, config.OVHClient, "/v2/deleted", true, time.Minute); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestManagedDatabaseObjectImport(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		new   func() fwresource.Resource
		id    string
		valid bool
	}{
		{NewCloudManagedDatabaseDatabaseResource, "pj/postgresql/service/object", true},
		{NewCloudManagedDatabaseDatabaseResource, "pj/service/object", false},
		// The type name of a connection pool holds its engine.
		{NewCloudManagedDatabasePostgresqlConnectionPoolResource, "pj/service/object", true},
		{NewCloudManagedDatabasePostgresqlConnectionPoolResource, "pj/postgresql/service/object", false},
	} {
		r := tc.new()
		resp := fwresource.ImportStateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
		r.(fwresource.ResourceWithImportState).ImportState(ctx, fwresource.ImportStateRequest{ID: tc.id}, &resp)
		if resp.Diagnostics.HasError() == tc.valid {
			t.Errorf("import ID %q: unexpected diagnostics %v", tc.id, resp.Diagnostics)
			continue
		}
		var engine types.String
		resp.State.GetAttribute(ctx, path.Root("engine"), &engine)
		if tc.valid && engine.ValueString() != "postgresql" {
			t.Errorf("import ID %q: unexpected engine %s", tc.id, engine)
		}
	}
}

func TestManagedDatabaseCategories(t *testing.T) {
	ctx := context.Background()
	if got := managedAnalyticsFamily.serviceEndpoint("pj", "kafka", "s1"); got != "/v2/publicCloud/project/pj/managedAnalytics/kafka/s1" {
		t.Errorf("unexpected endpoint %s", got)
	}

	// The engine is an argument, unless the type name holds it.
	for _, tc := range []struct {
		new      func() fwresource.Resource
		argument bool
	}{
		{NewCloudManagedDatabaseUserResource, true},
		{NewCloudManagedDatabasePostgresqlUserResource, false},
		{NewCloudManagedAnalyticsUserResource, true},
		{NewCloudManagedAnalyticsClickhouseDatabaseResource, false},
	} {
		r := tc.new()
		var meta fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "ovh"}, &meta)
		var resp fwresource.SchemaResponse
		r.Schema(ctx, fwresource.SchemaRequest{}, &resp)
		if engine := resp.Schema.Attributes["engine"]; engine.IsRequired() != tc.argument || engine.IsComputed() == tc.argument {
			t.Errorf("%s: expected engine as argument %t", meta.TypeName, tc.argument)
		}
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
	if len(names) != 12+31 {
		t.Errorf("expected 12 resources and 31 data sources, got %d", len(names))
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
