package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	testAccCloudManagedDatabasePostgresqlPrometheusResource = "ovh_cloud_managed_database_prometheus.prometheus"
	managedDatabasePrometheusTestPath                       = managedDatabasePostgresqlTestItemPath + "/prometheus"
	managedDatabasePrometheusTestResetPath                  = managedDatabasePrometheusTestPath + "/credentials/reset"
	managedDatabasePrometheusTestEndpoint                   = `{"username": "prometheus", "targets": [{"host": "postgresql-6b1b9c7e.database.cloud.ovh.net", "port": 9273}]}`
)

// managedDatabasePostgresqlTestEnabled is the test service with Prometheus enabled in its targetSpec.
var managedDatabasePostgresqlTestEnabled = strings.Replace(managedDatabasePostgresqlTestResponse, `"prometheusEnabled": false`, `"prometheusEnabled": true`, 1)

func newManagedDatabasePrometheusTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel], func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabasePrometheusResource().(*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel])
	r.config = config
	return r, bodies
}

// managedDatabasePrometheusTestModel returns the state of the Prometheus endpoint of the test service.
func managedDatabasePrometheusTestModel() CloudManagedDatabasePrometheusModel {
	return CloudManagedDatabasePrometheusModel{
		ServiceName:   types.StringValue("pj"),
		Engine:        types.StringValue("postgresql"),
		ServiceId:     types.StringValue(managedDatabasePostgresqlTestId),
		PasswordReset: types.StringNull(),
		Id:            types.StringValue(managedDatabasePostgresqlTestId),
		Username:      types.StringValue("prometheus"),
		Password:      types.StringValue("first-password"),
		Targets:       types.ListNull(types.ObjectType{AttrTypes: managedDatabasePrometheusTargetAttrTypes()}),
		Timeouts:      types.ObjectNull(managedDatabasePrometheusTimeoutsAttrTypes()),
	}
}

// managedDatabasePrometheusTestUnknown marks the computed attributes unknown, as in a plan.
func managedDatabasePrometheusTestUnknown(model CloudManagedDatabasePrometheusModel) CloudManagedDatabasePrometheusModel {
	model.Username = types.StringUnknown()
	model.Password = types.StringUnknown()
	model.Targets = types.ListUnknown(types.ObjectType{AttrTypes: managedDatabasePrometheusTargetAttrTypes()})
	return model
}

func managedDatabasePrometheusTestState(t *testing.T, model *CloudManagedDatabasePrometheusModel) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	NewCloudManagedDatabasePrometheusResource().Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
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

func managedDatabasePrometheusTestPlan(t *testing.T, model *CloudManagedDatabasePrometheusModel) tfsdk.Plan {
	t.Helper()
	state := managedDatabasePrometheusTestState(t, model)
	return tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
}

func TestCloudManagedDatabasePostgresqlPrometheusCreate(t *testing.T) {
	ctx := context.Background()
	// Disabled for the wait before the PUT and its two attempts, then enabled.
	disabled := managedDatabaseTestAnswer{http.StatusOK, managedDatabasePostgresqlTestResponse}
	r, bodies := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath:   {disabled, disabled, disabled, {http.StatusOK, managedDatabasePostgresqlTestEnabled}},
		"PUT " + managedDatabasePostgresqlTestItemPath:   {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, managedDatabasePostgresqlTestEnabled}},
		"POST " + managedDatabasePrometheusTestResetPath: {{http.StatusOK, `{"username": "prometheus", "password": "first-password"}`}},
		"GET " + managedDatabasePrometheusTestPath:       {{http.StatusOK, managedDatabasePrometheusTestEndpoint}},
	})

	plan := managedDatabasePrometheusTestUnknown(managedDatabasePrometheusTestModel())
	plan.Id = types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabasePrometheusTestState(t, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabasePrometheusTestPlan(t, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	puts := bodies("PUT " + managedDatabasePostgresqlTestItemPath)
	if len(puts) != 2 {
		t.Fatalf("expected the PUT to be retried once after the 409, got %d PUT", len(puts))
	}
	assertManagedDatabaseJSON(t, json.RawMessage(puts[1]), strings.Replace(managedDatabasePostgresqlTestPut, `"prometheusEnabled": false`, `"prometheusEnabled": true`, 1))

	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
	var got CloudManagedDatabasePrometheusModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != managedDatabasePostgresqlTestId || got.Password.ValueString() != "first-password" || len(got.Targets.Elements()) != 1 {
		t.Errorf("unexpected state: id %s, password %s, targets %s", got.Id, got.Password, got.Targets)
	}
}

func TestCloudManagedDatabasePostgresqlPrometheusRead(t *testing.T) {
	ctx := context.Background()

	t.Run("enabled", func(t *testing.T) {
		r, _ := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
			"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, managedDatabasePostgresqlTestEnabled}},
			"GET " + managedDatabasePrometheusTestPath:     {{http.StatusOK, managedDatabasePrometheusTestEndpoint}},
		})
		state := managedDatabasePrometheusTestModel()
		resp := fwresource.ReadResponse{State: managedDatabasePrometheusTestState(t, &state)}
		r.Read(ctx, fwresource.ReadRequest{State: managedDatabasePrometheusTestState(t, &state)}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
		var got CloudManagedDatabasePrometheusModel
		resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
		// The API never returns the password: the one of the state is kept.
		if got.Password.ValueString() != "first-password" || len(got.Targets.Elements()) != 1 {
			t.Errorf("unexpected state: password %s, targets %s", got.Password, got.Targets)
		}
	})

	for name, answer := range map[string]managedDatabaseTestAnswer{
		"disabled":        {http.StatusOK, managedDatabasePostgresqlTestResponse},
		"service deleted": {http.StatusNotFound, `{"message": "service not found"}`},
	} {
		t.Run(name, func(t *testing.T) {
			r, _ := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
				"GET " + managedDatabasePostgresqlTestItemPath: {answer},
			})
			state := managedDatabasePrometheusTestModel()
			resp := fwresource.ReadResponse{State: managedDatabasePrometheusTestState(t, &state)}
			r.Read(ctx, fwresource.ReadRequest{State: managedDatabasePrometheusTestState(t, &state)}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.State.Raw.IsNull() {
				t.Errorf("expected the resource to be removed from the state, got %v", resp.State.Raw)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlPrometheusUpdate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabasePrometheusTestResetPath: {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, `{"username": "prometheus", "password": "second-password"}`}},
		"GET " + managedDatabasePrometheusTestPath:       {{http.StatusOK, managedDatabasePrometheusTestEndpoint}},
	})

	state := managedDatabasePrometheusTestModel()
	plan := managedDatabasePrometheusTestUnknown(state)
	plan.PasswordReset = types.StringValue("v2")
	resp := fwresource.UpdateResponse{State: managedDatabasePrometheusTestState(t, &state)}
	r.Update(ctx, fwresource.UpdateRequest{
		Plan:  managedDatabasePrometheusTestPlan(t, &plan),
		State: managedDatabasePrometheusTestState(t, &state),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if got := len(bodies("POST " + managedDatabasePrometheusTestResetPath)); got != 2 {
		t.Errorf("expected the credentials reset to be retried once after the 409, got %d", got)
	}
	var got CloudManagedDatabasePrometheusModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Password.ValueString() != "second-password" || got.PasswordReset.ValueString() != "v2" {
		t.Errorf("unexpected state: password %s, password_reset %s", got.Password, got.PasswordReset)
	}
}

func TestCloudManagedDatabasePostgresqlPrometheusDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("disabled", func(t *testing.T) {
		enabled := managedDatabaseTestAnswer{http.StatusOK, managedDatabasePostgresqlTestEnabled}
		r, bodies := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
			"GET " + managedDatabasePostgresqlTestItemPath: {enabled, enabled, {http.StatusOK, managedDatabasePostgresqlTestResponse}},
			"PUT " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, managedDatabasePostgresqlTestResponse}},
		})
		state := managedDatabasePrometheusTestModel()
		resp := fwresource.DeleteResponse{State: managedDatabasePrometheusTestState(t, &state)}
		r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabasePrometheusTestState(t, &state)}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
		puts := bodies("PUT " + managedDatabasePostgresqlTestItemPath)
		if len(puts) != 1 {
			t.Fatalf("expected 1 PUT, got %d", len(puts))
		}
		assertManagedDatabaseJSON(t, json.RawMessage(puts[0]), managedDatabasePostgresqlTestPut)
	})

	t.Run("service deleted", func(t *testing.T) {
		r, _ := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
			"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusNotFound, `{"message": "service not found"}`}},
		})
		state := managedDatabasePrometheusTestModel()
		resp := fwresource.DeleteResponse{State: managedDatabasePrometheusTestState(t, &state)}
		r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabasePrometheusTestState(t, &state)}, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
		}
	})

	// The API refuses any PUT of a locked service: the endpoint stays enabled, with a warning, so that the
	// service can still be destroyed.
	t.Run("service locked", func(t *testing.T) {
		locked := strings.Replace(managedDatabasePostgresqlTestEnabled, `"resourceStatus": "READY"`, `"resourceStatus": "SUSPENDED"`, 1)
		r, _ := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
			"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, locked}},
			"PUT " + managedDatabasePostgresqlTestItemPath: {{http.StatusConflict, `{"class": "Client::Conflict::ServiceLocked", "message": "service is LOCKED"}`}},
		})
		state := managedDatabasePrometheusTestModel()
		resp := fwresource.DeleteResponse{State: managedDatabasePrometheusTestState(t, &state)}
		r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabasePrometheusTestState(t, &state)}, &resp)
		if resp.Diagnostics.HasError() || resp.Diagnostics.WarningsCount() != 1 {
			t.Errorf("expected one warning, got %v", resp.Diagnostics)
		}
	})
}

// Only a new password_reset resets the password: not its removal, nor a change of timeouts, nor a
// password_reset unknown in the plan that turns out null.
func TestCloudManagedDatabasePostgresqlPrometheusUpdateKeepsPassword(t *testing.T) {
	ctx := context.Background()
	for name, change := range map[string]func(*CloudManagedDatabasePrometheusModel){
		"password_reset removed": func(m *CloudManagedDatabasePrometheusModel) { m.PasswordReset = types.StringNull() },
		"timeouts changed": func(m *CloudManagedDatabasePrometheusModel) {
			m.Timeouts = types.ObjectValueMust(managedDatabasePrometheusTimeoutsAttrTypes(), map[string]attr.Value{"create": types.StringValue("1h"), "delete": types.StringNull()})
		},
		"unknown password_reset turned null": func(m *CloudManagedDatabasePrometheusModel) {
			m.PasswordReset, m.Password = types.StringNull(), types.StringUnknown()
		},
	} {
		t.Run(name, func(t *testing.T) {
			r, bodies := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
				"POST " + managedDatabasePrometheusTestResetPath: {{http.StatusOK, `{"username": "prometheus", "password": "new-password"}`}},
				"GET " + managedDatabasePrometheusTestPath:       {{http.StatusOK, managedDatabasePrometheusTestEndpoint}},
			})
			state := managedDatabasePrometheusTestModel()
			state.PasswordReset = types.StringValue("v1")
			plan := state
			change(&plan)
			resp := fwresource.UpdateResponse{State: managedDatabasePrometheusTestState(t, &state)}
			r.Update(ctx, fwresource.UpdateRequest{Plan: managedDatabasePrometheusTestPlan(t, &plan), State: managedDatabasePrometheusTestState(t, &state)}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if n := len(bodies("POST " + managedDatabasePrometheusTestResetPath)); n != 0 {
				t.Errorf("expected no credentials reset, got %d", n)
			}
			var got CloudManagedDatabasePrometheusModel
			resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
			if got.Password.ValueString() != "first-password" {
				t.Errorf("expected the password kept, got %s", got.Password)
			}
		})
	}
}

// The plan announces a new password only when password_reset resets it.
func TestCloudManagedDatabasePostgresqlPrometheusModifyPlan(t *testing.T) {
	ctx := context.Background()
	r := NewCloudManagedDatabasePrometheusResource().(*managedDatabasePrometheusResource[CloudManagedDatabasePrometheusModel, *CloudManagedDatabasePrometheusModel])
	state := managedDatabasePrometheusTestModel()
	state.PasswordReset = types.StringValue("v1")
	for name, tc := range map[string]struct {
		passwordReset types.String
		resets        bool
	}{
		"new value": {types.StringValue("v2"), true},
		"unknown":   {types.StringUnknown(), true},
		"removed":   {types.StringNull(), false},
		"unchanged": {types.StringValue("v1"), false},
	} {
		t.Run(name, func(t *testing.T) {
			plan := state
			plan.PasswordReset = tc.passwordReset
			resp := fwresource.ModifyPlanResponse{Plan: managedDatabasePrometheusTestPlan(t, &plan)}
			r.ModifyPlan(ctx, fwresource.ModifyPlanRequest{
				Plan:  managedDatabasePrometheusTestPlan(t, &plan),
				State: managedDatabasePrometheusTestState(t, &state),
			}, &resp)
			var got types.String
			resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("password"), &got)...)
			if resp.Diagnostics.HasError() || got.IsUnknown() != tc.resets {
				t.Errorf("unexpected password %s, diagnostics %v", got, resp.Diagnostics)
			}
		})
	}
}

// The wait for the service follows timeouts.create.
func TestCloudManagedDatabasePostgresqlPrometheusCreateTimeout(t *testing.T) {
	ctx := context.Background()
	updating := strings.Replace(managedDatabasePostgresqlTestResponse, `"resourceStatus": "READY"`, `"resourceStatus": "UPDATING"`, 1)
	r, _ := newManagedDatabasePrometheusTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, updating}},
	})

	plan := managedDatabasePrometheusTestUnknown(managedDatabasePrometheusTestModel())
	plan.Id = types.StringUnknown()
	plan.Timeouts = types.ObjectValueMust(managedDatabasePrometheusTimeoutsAttrTypes(), map[string]attr.Value{"create": types.StringValue("50ms"), "delete": types.StringNull()})
	resp := fwresource.CreateResponse{State: managedDatabasePrometheusTestState(t, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabasePrometheusTestPlan(t, &plan)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the wait for the service to time out")
	}
}

// The plan of an OUT_OF_SYNC service, which Terraform can see as no change, leaves prometheus_enabled
// unknown, as in any update, unless the configuration sets it.
func TestCloudManagedDatabasePostgresqlOutOfSyncPlanPrometheus(t *testing.T) {
	ctx := context.Background()
	r := &cloudManagedDatabaseServiceResource{category: managedDatabaseFamily}
	model := managedDatabasePostgresqlTestState(t)
	model.ResourceStatus = types.StringValue("OUT_OF_SYNC")
	model.PrometheusEnabled = types.BoolValue(true)
	state := managedDatabasePostgresqlTestResourceState(t, &model)

	for name, configured := range map[string]types.Bool{"unset": types.BoolNull(), "set": types.BoolValue(true)} {
		t.Run(name, func(t *testing.T) {
			config := model
			config.PrometheusEnabled = configured
			raw := managedDatabasePostgresqlTestResourceState(t, &config).Raw
			resp := fwresource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}
			r.ModifyPlan(ctx, fwresource.ModifyPlanRequest{
				State:  state,
				Plan:   tfsdk.Plan{Schema: state.Schema, Raw: state.Raw},
				Config: tfsdk.Config{Schema: state.Schema, Raw: raw},
			}, &resp)
			var got types.Bool
			resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("prometheus_enabled"), &got)...)
			if resp.Diagnostics.HasError() || got.IsUnknown() != configured.IsNull() {
				t.Errorf("unexpected prometheus_enabled %s, diagnostics %v", got, resp.Diagnostics)
			}
		})
	}
}

// newManagedDatabasePrometheusLiveLikeTestAPI adds to the live-like API the Prometheus endpoint of its
// service, which exists while the prometheusEnabled flag of the targetSpec is on.
func newManagedDatabasePrometheusLiveLikeTestAPI(t *testing.T) (string, func(func(spec map[string]any)), func()) {
	t.Helper()
	endpoint, edit, outOfSync := newManagedDatabaseLiveLikeTestAPI(t)
	target, err := url.Parse(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	enabled := func() bool {
		var on bool
		edit(func(spec map[string]any) { on, _ = spec["prometheusEnabled"].(bool) })
		return on
	}
	resets := 0
	mux := http.NewServeMux()
	mux.Handle("/", httputil.NewSingleHostReverseProxy(&url.URL{Scheme: target.Scheme, Host: target.Host}))
	mux.HandleFunc("GET "+managedDatabasePrometheusTestPath, func(w http.ResponseWriter, r *http.Request) {
		if !enabled() {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "prometheus is disabled"}`)
			return
		}
		fmt.Fprint(w, managedDatabasePrometheusTestEndpoint)
	})
	mux.HandleFunc("POST "+managedDatabasePrometheusTestResetPath, func(w http.ResponseWriter, r *http.Request) {
		resets++
		fmt.Fprintf(w, `{"username": "prometheus", "password": "password-%d"}`, resets)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server.URL + "/1.0", edit, outOfSync
}

// Terraform destroys the Prometheus resource before it updates the service it depends on: the update
// must not turn the endpoint back on, even when the service is OUT_OF_SYNC. A password reset keeps the
// username and the targets in the plan.
func TestCloudManagedDatabasePostgresqlPrometheusLiveLikePlan(t *testing.T) {
	endpoint, edit, outOfSync := newManagedDatabasePrometheusLiveLikeTestAPI(t)
	config := func(description string, prometheus int, arguments string) string {
		return fmt.Sprintf(`
provider "ovh" {
  endpoint           = %q
  application_key    = "key"
  application_secret = "secret"
  consumer_key       = "consumer"
}

resource "ovh_cloud_managed_database" "db" {
  service_name = "pj"
  engine       = "postgresql"
  region       = "eu-west-gra"
  version      = "17"
  plan         = "business"
  flavor       = "b3-8"
  description  = %q
}

resource "ovh_cloud_managed_database_prometheus" "prometheus" {
  count          = %d
  service_name   = ovh_cloud_managed_database.db.service_name
  engine         = ovh_cloud_managed_database.db.engine
  service_id     = ovh_cloud_managed_database.db.id
  %s
}
`, endpoint, description, prometheus, arguments)
	}
	enabled := func(want bool) resource.TestCheckFunc {
		return func(*terraform.State) error {
			var got any
			edit(func(spec map[string]any) { got = spec["prometheusEnabled"] })
			if got != want {
				return fmt.Errorf("expected prometheusEnabled %t in the API, got %v", want, got)
			}
			return nil
		}
	}
	// The address of the instance in the plan, and in the state.
	planned, prometheus := testAccCloudManagedDatabasePostgresqlPrometheusResource+"[0]", testAccCloudManagedDatabasePostgresqlPrometheusResource+".0"
	keptPassword := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(planned, plancheck.ResourceActionUpdate),
			plancheck.ExpectKnownValue(planned, tfjsonpath.New("password"), knownvalue.StringExact("password-3")),
		},
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("first", 1, `password_reset = "v1"`),
				Check:  enabled(true),
			},
			{
				Config: config("second", 0, `password_reset = "v1"`),
				Check: resource.ComposeTestCheckFunc(
					enabled(false),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "prometheus_enabled", "false"),
				),
			},
			{
				Config: config("second", 1, `password_reset = "v1"`),
				Check:  enabled(true),
			},
			{
				Config: config("second", 1, `password_reset = "v2"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(planned, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(planned, tfjsonpath.New("username"), knownvalue.StringExact("prometheus")),
						plancheck.ExpectKnownValue(planned, tfjsonpath.New("targets"), knownvalue.ListSizeExact(1)),
					},
				},
				Check: resource.TestCheckResourceAttr(prometheus, "password", "password-3"),
			},
			// Neither a change of timeouts nor the removal of password_reset resets the password.
			{
				Config:           config("second", 1, `password_reset = "v2"`+"\n  timeouts = { delete = \"1h\" }"),
				ConfigPlanChecks: keptPassword,
				Check:            resource.TestCheckResourceAttr(prometheus, "password", "password-3"),
			},
			{
				Config:           config("second", 1, ""),
				ConfigPlanChecks: keptPassword,
				Check:            resource.TestCheckResourceAttr(prometheus, "password", "password-3"),
			},
			{
				PreConfig: outOfSync,
				Config:    config("second", 0, ""),
				Check:     enabled(false),
			},
		},
	})
}

// testAccCloudManagedDatabasePostgresqlPrometheusConfig renders a service whose Prometheus endpoint
// the resource enables, with password_reset = null, as if removed from the configuration, when
// passwordReset is "".
func testAccCloudManagedDatabasePostgresqlPrometheusConfig(name, passwordReset string) string {
	value := "null"
	if passwordReset != "" {
		value = fmt.Sprintf("%q", passwordReset)
	}
	return testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), "") + fmt.Sprintf(`
resource "ovh_cloud_managed_database_prometheus" "prometheus" {
  service_name   = ovh_cloud_managed_database.db.service_name
  engine         = ovh_cloud_managed_database.db.engine
  service_id     = ovh_cloud_managed_database.db.id
  password_reset = %s
}
`, value)
}

func TestAccCloudManagedDatabasePostgresqlPrometheus_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlPrometheusConfig(name, "v1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(testAccCloudManagedDatabasePostgresqlPrometheusResource, "id", testAccCloudManagedDatabasePostgresqlResource, "id"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlPrometheusResource, "username"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlPrometheusResource, "password"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlPrometheusResource, "targets.0.host"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlPrometheusResource, "targets.0.port"),
				),
			},
			{
				Config: testAccCloudManagedDatabasePostgresqlPrometheusConfig(name, "v2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlPrometheusResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "prometheus_enabled", "true"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlPrometheusResource, "password"),
				),
			},
			{
				// Removing password_reset keeps the password.
				Config: testAccCloudManagedDatabasePostgresqlPrometheusConfig(name, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlPrometheusResource, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(testAccCloudManagedDatabasePostgresqlPrometheusResource, tfjsonpath.New("password"), knownvalue.NotNull()),
					},
				},
			},
			{
				ResourceName:            testAccCloudManagedDatabasePostgresqlPrometheusResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "password_reset"},
				ImportStateIdFunc:       testAccCloudManagedImportStateIdFunc(testAccCloudManagedDatabasePostgresqlPrometheusResource),
			},
		},
	})
}

func TestCloudManagedDatabaseMongodbPrometheusCreate(t *testing.T) {
	ctx := context.Background()
	service := "/v2/publicCloud/project/pj/managedDatabase/mongodb/" + managedDatabasePostgresqlTestId
	// Already enabled: no PUT.
	config, _ := newManagedDatabaseTestConfig(t, map[string][]managedDatabaseTestAnswer{
		"GET " + service: {{http.StatusOK, managedDatabasePostgresqlTestEnabled}},
		"POST " + service + "/prometheus/credentials/reset": {{http.StatusOK, `{"username": "prometheus", "password": "secret"}`}},
		"GET " + service + "/prometheus":                    {{http.StatusOK, `{"username": "prometheus", "srvDomain": "_prometheus._tcp.mongodb-6b1b9c7e.database.cloud.ovh.net"}`}},
	})
	r := NewCloudManagedDatabaseMongodbPrometheusResource()
	r.(fwresource.ResourceWithConfigure).Configure(ctx, fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})

	plan := cloudManagedDatabaseMongodbPrometheusModel{
		managedDatabasePrometheusModel: managedDatabasePrometheusModel{
			ServiceName:   types.StringValue("pj"),
			Engine:        types.StringValue("mongodb"),
			ServiceId:     types.StringValue(managedDatabasePostgresqlTestId),
			PasswordReset: types.StringNull(),
			Id:            types.StringUnknown(),
			Username:      types.StringUnknown(),
			Password:      types.StringUnknown(),
			Timeouts:      types.ObjectNull(managedDatabasePrometheusTimeoutsAttrTypes()),
		},
		SrvDomain: types.StringUnknown(),
	}
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var got cloudManagedDatabaseMongodbPrometheusModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Engine.ValueString() != "mongodb" || got.SrvDomain.IsNull() || got.Password.ValueString() != "secret" {
		t.Errorf("unexpected state: %v", resp.State.Raw)
	}
}
