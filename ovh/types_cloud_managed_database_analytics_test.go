package ovh

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const managedDatabasePostgresqlTestId = "6b1b9c7e-2d4f-4b8a-9c3e-5f6a7b8c9d01"

// managedDatabasePostgresqlTestResponse is a READY service, as returned by a GET.
const managedDatabasePostgresqlTestResponse = `{
	"id": "` + managedDatabasePostgresqlTestId + `",
	"checksum": "d41d8cd98f00b204e9800998ecf8427e",
	"createdAt": "2026-11-12T10:00:00Z",
	"updatedAt": "2026-11-12T10:20:00Z",
	"resourceStatus": "READY",
	"currentTasks": [],
	"targetSpec": {
		"advancedConfiguration": {"pg.log_min_duration_statement": "1000"},
		"backup": {"regions": ["eu-west-gra"], "time": "02:00"},
		"deletionProtected": false,
		"description": "testacc-terraform-db",
		"flavor": "b3-8",
		"maintenanceWindow": {"startTimeUtc": "03:00"},
		"network": {"authorizedIps": {"rules": [{"cidr": "203.0.113.0/24", "description": "office"}]}, "public": false},
		"nodesCount": 2,
		"plan": "business",
		"prometheusEnabled": false,
		"storage": {"size": {"unit": "GB", "value": 160}},
		"version": "17"
	},
	"currentState": {
		"advancedConfiguration": {"pg.log_min_duration_statement": "1000", "pg.idle_in_transaction_session_timeout": "0"},
		"backup": {"pitr": "2026-11-12T10:15:00Z", "regions": ["eu-west-gra"], "retentionDays": 2, "time": "02:00"},
		"capabilities": {"BACKUPS": {"create": "ENABLED", "delete": "DISABLED", "read": "ENABLED", "update": "DISABLED"}},
		"deletionProtected": false,
		"description": "testacc-terraform-db",
		"endpoints": [{
			"component": "POSTGRESQL",
			"domain": "postgresql-6b1b9c7e.database.cloud.ovh.net",
			"path": "/defaultdb",
			"port": 20184,
			"public": false,
			"scheme": "postgresql",
			"ssl": true,
			"sslMode": "REQUIRE",
			"uri": "postgresql://ovhadmin@postgresql-6b1b9c7e.database.cloud.ovh.net:20184/defaultdb?sslmode=require"
		}],
		"engineCapabilities": {
			"advancedConfigurationProperties": [{"name": "pg.log_min_duration_statement", "type": "LONG", "description": "Log statements slower than this many milliseconds", "default": "-1", "minimum": -1, "maximum": 2147483647, "requiresRestart": false, "values": null}],
			"backupRegions": ["eu-west-gra", "eu-west-sbg"],
			"categories": null,
			"commands": null,
			"connectors": null,
			"integrations": [{"sourceEngine": "POSTGRESQL", "destinationEngine": "GRAFANA", "type": "GRAFANA_DATASOURCE", "parameters": null}]
		},
		"flavor": "b3-8",
		"maintenanceWindow": {"startTimeUtc": "03:00"},
		"network": {
			"authorizedIps": {"rules": [{"cidr": "203.0.113.0/24", "createdAt": "2026-11-12T10:00:00Z", "description": "office"}]},
			"networkId": null,
			"public": false,
			"subnetId": null
		},
		"nodes": [{
			"createdAt": "2026-11-12T10:00:00Z",
			"flavor": "b3-8",
			"id": "0c9e1f2a-3b4c-4d5e-8f6a-7b8c9d0e1f2a",
			"location": {"availabilityZone": null, "region": "eu-west-gra"},
			"name": "node-1",
			"port": 20184,
			"role": "PRIMARY",
			"status": "READY"
		}],
		"nodesCount": 2,
		"plan": "business",
		"prometheusEnabled": false,
		"status": "READY",
		"storage": {"size": {"unit": "GB", "value": 160}, "type": "HIGH_SPEED"},
		"version": "17"
	}
}`

func managedDatabasePostgresqlTestAPIResponse(t *testing.T, body string) *CloudManagedDatabaseServiceAPIResponse {
	t.Helper()
	var res CloudManagedDatabaseServiceAPIResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return &res
}

// managedDatabasePostgresqlNullModel returns a model whose attributes are all null, as the state of
// an import before its first read.
func managedDatabasePostgresqlNullModel() CloudManagedDatabaseServiceModel {
	return CloudManagedDatabaseServiceModel{
		ServiceName:           types.StringValue("pj"),
		Engine:                types.StringValue("postgresql"),
		Id:                    types.StringValue(managedDatabasePostgresqlTestId),
		AdvancedConfiguration: types.MapNull(types.StringType),
		Backup:                types.ObjectNull(managedDatabaseBackupAttrTypes()),
		MaintenanceWindow:     types.ObjectNull(managedDatabaseMaintenanceWindowAttrTypes()),
		Network:               types.ObjectNull(managedDatabaseNetworkAttrTypes()),
		Storage:               types.ObjectNull(managedDatabaseStorageAttrTypes()),
		CurrentTasks:          types.ListNull(types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()}),
		CurrentState:          types.ObjectNull(managedDatabaseCurrentStateAttrTypes()),
	}
}

// managedDatabasePostgresqlTestState returns the state of the service of the test response, whose
// configuration declares one advanced configuration key.
func managedDatabasePostgresqlTestState(t *testing.T) CloudManagedDatabaseServiceModel {
	t.Helper()
	model := managedDatabasePostgresqlNullModel()
	model.AdvancedConfiguration = types.MapValueMust(types.StringType, map[string]attr.Value{
		"pg.log_min_duration_statement": types.StringValue("1000"),
	})
	if diags := model.MergeWith(context.Background(), managedDatabasePostgresqlTestAPIResponse(t, managedDatabasePostgresqlTestResponse)); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	return model
}

func assertManagedDatabaseJSON(t *testing.T, got any, want string) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var gotValue, wantValue any
	if err := json.Unmarshal(gotJSON, &gotValue); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("invalid expected JSON: %s", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("unexpected body\n got: %s\nwant: %s", gotJSON, want)
	}
}

func TestCloudManagedDatabasePostgresqlToCreate(t *testing.T) {
	ctx := context.Background()
	model := managedDatabasePostgresqlNullModel()
	model.Region = types.StringValue("eu-west-gra")
	model.Version = types.StringValue("17")
	model.Plan = types.StringValue("business")
	model.Flavor = types.StringValue("b3-8")
	model.NodesCount = types.Int64Unknown()
	model.Description = types.StringValue("testacc-terraform-db")
	model.DeletionProtected = types.BoolValue(true)
	model.PrometheusEnabled = types.BoolUnknown()
	model.AdvancedConfiguration = types.MapValueMust(types.StringType, map[string]attr.Value{
		"pg.log_min_duration_statement": types.StringValue("1000"),
	})
	// A planned object can hold unknown attributes, which are left out of the body.
	model.Backup = types.ObjectValueMust(managedDatabaseBackupAttrTypes(), map[string]attr.Value{
		"regions": types.ListUnknown(types.StringType),
		"time":    types.StringValue("02:00"),
	})
	model.Network = types.ObjectValueMust(managedDatabaseNetworkAttrTypes(), map[string]attr.Value{
		"network_id": types.StringValue("net"),
		"subnet_id":  types.StringValue("subnet"),
		"authorized_ips": types.ObjectValueMust(managedDatabaseNetworkAttrTypes()["authorized_ips"].(types.ObjectType).AttrTypes, map[string]attr.Value{
			"rules": types.SetValueMust(managedDatabaseNetworkAttrTypes()["authorized_ips"].(types.ObjectType).AttrTypes["rules"].(types.SetType).ElemType, nil),
		}),
		"public": types.BoolValue(true),
	})
	model.Storage = types.ObjectUnknown(managedDatabaseStorageAttrTypes())
	resourceModel := cloudManagedDatabaseServiceResourceModel{
		CloudManagedDatabaseServiceModel: model,
		ForkSourceServiceId:              types.StringValue("source"),
		ForkSourcePointInTime:            types.StringValue("2026-11-12T10:00:00Z"),
	}

	payload, diags := resourceModel.ToCreate(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	assertManagedDatabaseJSON(t, payload, `{"targetSpec": {
		"region": "eu-west-gra",
		"version": "17",
		"plan": "business",
		"flavor": "b3-8",
		"description": "testacc-terraform-db",
		"deletionProtected": true,
		"backup": {"time": "02:00"},
		"network": {"networkId": "net", "subnetId": "subnet", "authorizedIps": {"rules": []}, "public": true},
		"forkSourceServiceId": "source",
		"forkSourcePointInTime": "2026-11-12T10:00:00Z"
	}}`)
}

// managedDatabasePostgresqlTestReplace is the PUT body of managedDatabasePostgresqlTestResponse, without
// the region, the private network and the public exposure, which the PUT doesn't take.
const managedDatabasePostgresqlTestReplace = `{"targetSpec": {
	"advancedConfiguration": {"pg.log_min_duration_statement": "1000"},
	"backup": {"regions": ["eu-west-gra"], "time": "02:00"},
	"deletionProtected": false,
	"description": "testacc-terraform-db",
	"flavor": "b3-8",
	"maintenanceWindow": {"startTimeUtc": "03:00"},
	"network": {"authorizedIps": {"rules": [{"cidr": "203.0.113.0/24", "description": "office"}]}},
	"nodesCount": 2,
	"plan": "business",
	"prometheusEnabled": false,
	"storage": {"size": {"unit": "GB", "value": 160}},
	"version": "17"
}}`

func TestCloudManagedDatabasePostgresqlToPut(t *testing.T) {
	ctx := context.Background()
	state := managedDatabasePostgresqlTestState(t)
	ipsType := managedDatabaseNetworkAttrTypes()["authorized_ips"].(types.ObjectType)
	ruleType := ipsType.AttrTypes["rules"].(types.SetType).ElemType.(types.ObjectType)
	withRules := func(rules ...attr.Value) types.Object {
		attrs := state.Network.Attributes()
		attrs["authorized_ips"] = types.ObjectValueMust(ipsType.AttrTypes, map[string]attr.Value{"rules": types.SetValueMust(ruleType, rules)})
		return types.ObjectValueMust(managedDatabaseNetworkAttrTypes(), attrs)
	}

	t.Run("whole targetSpec", func(t *testing.T) {
		replace, diags := state.ToPut(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		assertManagedDatabaseJSON(t, replace, managedDatabasePostgresqlTestReplace)
	})

	t.Run("nulls", func(t *testing.T) {
		plan := state
		plan.Description = types.StringNull()
		plan.AdvancedConfiguration = types.MapNull(types.StringType)
		plan.Network = withRules(types.ObjectValueMust(ruleType.AttrTypes, map[string]attr.Value{
			"cidr":        types.StringValue("198.51.100.0/24"),
			"description": types.StringNull(),
		}))
		replace, diags := plan.ToPut(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		spec := replace.TargetSpec
		assertManagedDatabaseJSON(t, []any{spec.Description, spec.AdvancedConfiguration, spec.Network}, `[null, null,
			{"authorizedIps": {"rules": [{"cidr": "198.51.100.0/24", "description": null}]}}]`)
	})

	t.Run("empty", func(t *testing.T) {
		plan := state
		plan.AdvancedConfiguration = types.MapValueMust(types.StringType, map[string]attr.Value{})
		plan.Network = withRules()
		replace, diags := plan.ToPut(ctx)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		spec := replace.TargetSpec
		assertManagedDatabaseJSON(t, []any{spec.AdvancedConfiguration, spec.Network}, `[{}, {"authorizedIps": {"rules": []}}]`)
	})
}

func TestCloudManagedDatabasePostgresqlMergeWith(t *testing.T) {
	ctx := context.Background()
	response := managedDatabasePostgresqlTestAPIResponse(t, managedDatabasePostgresqlTestResponse)

	t.Run("import", func(t *testing.T) {
		model := managedDatabasePostgresqlNullModel()
		if diags := model.MergeWith(ctx, response); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		for name, got := range map[string]attr.Value{
			"region":             model.Region,
			"version":            model.Version,
			"nodes_count":        model.NodesCount,
			"description":        model.Description,
			"deletion_protected": model.DeletionProtected,
			"resource_status":    model.ResourceStatus,
		} {
			want := map[string]attr.Value{
				"region":             types.StringValue("eu-west-gra"),
				"version":            types.StringValue("17"),
				"nodes_count":        types.Int64Value(2),
				"description":        types.StringValue("testacc-terraform-db"),
				"deletion_protected": types.BoolValue(false),
				"resource_status":    types.StringValue("READY"),
			}[name]
			if !got.Equal(want) {
				t.Errorf("unexpected %s %s, want %s", name, got, want)
			}
		}
		if got := len(model.AdvancedConfiguration.Elements()); got != 1 {
			t.Errorf("expected the advanced configuration of the targetSpec, got %s", model.AdvancedConfiguration)
		}
		network := model.Network.Attributes()
		if !network["network_id"].IsNull() || !network["subnet_id"].IsNull() || !network["public"].Equal(types.BoolValue(false)) {
			t.Errorf("expected no private network and no public address, got %s", model.Network)
		}
		if got := network["authorized_ips"].String(); !strings.Contains(got, "203.0.113.0/24") || strings.Contains(got, "created_at") {
			t.Errorf("unexpected authorized_ips %s", got)
		}
		currentState := model.CurrentState.Attributes()
		if got := currentState["endpoints"].String(); !strings.Contains(got, "sslmode=require") {
			t.Errorf("unexpected endpoints %s", got)
		}
		// The current state keeps what the arguments leave out: the creation date of the rules, and the
		// capabilities.
		if got := currentState["network"].String(); !strings.Contains(got, "2026-11-12T10:00:00Z") {
			t.Errorf("expected the creation date of the rules in current_state.network, got %s", got)
		}
		if got := currentState["capabilities"].String(); !strings.Contains(got, `"BACKUPS"`) {
			t.Errorf("unexpected capabilities %s", got)
		}
		if got := currentState["engine_capabilities"].String(); !strings.Contains(got, "pg.log_min_duration_statement") || !strings.Contains(got, "GRAFANA_DATASOURCE") {
			t.Errorf("unexpected engine_capabilities %s", got)
		}
		if !model.CurrentTasks.Equal(types.ListValueMust(types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()}, []attr.Value{})) {
			t.Errorf("expected no current task, got %s", model.CurrentTasks)
		}
		if got := len(currentState["advanced_configuration"].(types.Map).Elements()); got != 2 {
			t.Errorf("expected 2 advanced configuration keys in current_state, got %d", got)
		}
	})

	t.Run("whole advanced configuration", func(t *testing.T) {
		// A key set outside of Terraform shows up, so that the next PUT removes it.
		model := managedDatabasePostgresqlNullModel()
		model.AdvancedConfiguration = types.MapValueMust(types.StringType, map[string]attr.Value{
			"pg.log_temp_files": types.StringValue("1024"),
		})
		if diags := model.MergeWith(ctx, response); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		want := types.MapValueMust(types.StringType, map[string]attr.Value{
			"pg.log_min_duration_statement": types.StringValue("1000"),
		})
		if !model.AdvancedConfiguration.Equal(want) {
			t.Errorf("unexpected advanced configuration %s, want %s", model.AdvancedConfiguration, want)
		}
	})

	t.Run("null advanced configuration", func(t *testing.T) {
		response := managedDatabasePostgresqlTestAPIResponse(t, strings.Replace(managedDatabasePostgresqlTestResponse,
			`"advancedConfiguration": {"pg.log_min_duration_statement": "1000"},`, `"advancedConfiguration": null,`, 1))
		empty := types.MapValueMust(types.StringType, map[string]attr.Value{})
		for name, tc := range map[string]struct{ model, want types.Map }{
			"left out":   {types.MapUnknown(types.StringType), types.MapNull(types.StringType)},
			"empty kept": {empty, empty},
			"removed": {types.MapValueMust(types.StringType, map[string]attr.Value{"pg.log_temp_files": types.StringValue("1024")}),
				types.MapNull(types.StringType)},
		} {
			model := managedDatabasePostgresqlNullModel()
			model.AdvancedConfiguration = tc.model
			if diags := model.MergeWith(ctx, response); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !model.AdvancedConfiguration.Equal(tc.want) {
				t.Errorf("%s: unexpected advanced configuration %s, want %s", name, model.AdvancedConfiguration, tc.want)
			}
		}
	})

	t.Run("data source", func(t *testing.T) {
		model := managedDatabasePostgresqlNullModel()
		model.AdvancedConfiguration = types.MapUnknown(types.StringType)
		if diags := model.MergeWith(ctx, response); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got := len(model.AdvancedConfiguration.Elements()); got != 1 {
			t.Errorf("expected every key of the targetSpec, got %s", model.AdvancedConfiguration)
		}
	})

	t.Run("configured region kept", func(t *testing.T) {
		model := managedDatabasePostgresqlNullModel()
		model.Region = types.StringValue("EU-WEST-GRA")
		if diags := model.MergeWith(ctx, response); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got := model.Region.ValueString(); got != "EU-WEST-GRA" {
			t.Errorf("unexpected region %q", got)
		}
	})

	t.Run("empty description", func(t *testing.T) {
		model := managedDatabasePostgresqlNullModel()
		emptied := managedDatabasePostgresqlTestAPIResponse(t, strings.ReplaceAll(managedDatabasePostgresqlTestResponse, `"testacc-terraform-db"`, `""`))
		if diags := model.MergeWith(ctx, emptied); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if !model.Description.IsNull() {
			t.Errorf("expected a null description, got %s", model.Description)
		}
	})

	// The state saved right after the POST comes from a plan, whose computed attributes are unknown.
	t.Run("no unknown left", func(t *testing.T) {
		model := managedDatabasePostgresqlNullModel()
		model.Region = types.StringValue("eu-west-gra")
		model.NodesCount = types.Int64Unknown()
		model.DeletionProtected = types.BoolUnknown()
		model.PrometheusEnabled = types.BoolUnknown()
		model.Backup = types.ObjectUnknown(managedDatabaseBackupAttrTypes())
		model.MaintenanceWindow = types.ObjectUnknown(managedDatabaseMaintenanceWindowAttrTypes())
		model.Network = types.ObjectUnknown(managedDatabaseNetworkAttrTypes())
		model.Storage = types.ObjectUnknown(managedDatabaseStorageAttrTypes())
		model.Id = types.StringUnknown()
		model.CurrentState = types.ObjectUnknown(managedDatabaseCurrentStateAttrTypes())
		created := managedDatabasePostgresqlTestAPIResponse(t, `{"id": "`+managedDatabasePostgresqlTestId+`", "resourceStatus": "CREATING", "targetSpec": {"version": "17", "plan": "business", "flavor": "b3-8"}}`)
		if diags := model.MergeWith(ctx, created); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		state := managedDatabasePostgresqlTestResourceState(t, &model)
		if !state.Raw.IsFullyKnown() {
			t.Errorf("expected a fully known state, got %v", state.Raw)
		}
	})
}

func managedDatabasePostgresqlTestResourceState(t *testing.T, model *CloudManagedDatabaseServiceModel) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var schemaResp resource.SchemaResponse
	NewCloudManagedDatabaseResource().Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
	if model != nil {
		if diags := state.Set(ctx, &cloudManagedDatabaseServiceResourceModel{CloudManagedDatabaseServiceModel: *model}); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	}
	return state
}

// The model must fit both schemas, which are written separately.
func TestCloudManagedDatabasePostgresqlSchemas(t *testing.T) {
	ctx := context.Background()
	model := managedDatabasePostgresqlTestState(t)

	var resourceSchema resource.SchemaResponse
	NewCloudManagedDatabaseResource().Schema(ctx, resource.SchemaRequest{}, &resourceSchema)
	if diags := resourceSchema.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("invalid resource schema: %v", diags)
	}
	managedDatabasePostgresqlTestResourceState(t, &model)

	var dataSourceSchema datasource.SchemaResponse
	NewCloudManagedDatabaseDataSource().Schema(ctx, datasource.SchemaRequest{}, &dataSourceSchema)
	if diags := dataSourceSchema.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("invalid data source schema: %v", diags)
	}
	state := tfsdk.State{Schema: dataSourceSchema.Schema, Raw: tftypes.NewValue(dataSourceSchema.Schema.Type().TerraformType(ctx), nil)}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Errorf("model doesn't fit the data source schema: %v", diags)
	}
}
