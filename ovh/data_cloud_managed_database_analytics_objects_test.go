package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// managedDatabaseDataSourceTestRead reads d, configured with the given string arguments, from the fake
// API, and returns its state. It fails when an attribute and the decoded model don't match.
func managedDatabaseDataSourceTestRead(t *testing.T, d datasource.DataSource, arguments map[string]string, answers map[string][]managedDatabaseTestAnswer) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	providerConfig, _ := newManagedDatabaseTestConfig(t, answers)
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: providerConfig}, &datasource.ConfigureResponse{})

	var schemaResp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if diags := schemaResp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("invalid schema: %v", diags)
	}
	// The tests read PostgreSQL services: the engine is postgresql wherever it is an argument.
	if engine, ok := schemaResp.Schema.Attributes["engine"]; ok && engine.IsRequired() && arguments["engine"] == "" {
		arguments["engine"] = "postgresql"
	}
	config := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
	for name, value := range arguments {
		if diags := config.SetAttribute(ctx, path.Root(name), types.StringValue(value)); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	}

	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: config.Raw.Copy()}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: config.Raw}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
	return resp.State
}
