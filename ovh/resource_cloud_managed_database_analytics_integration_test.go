package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	managedDatabaseIntegrationTestPath     = managedDatabasePostgresqlTestItemPath + "/integration"
	managedDatabaseIntegrationTestItemPath = managedDatabaseIntegrationTestPath + "/i1"
)

// managedDatabaseIntegrationTestAnswer is an integration as the API returns it: without its parameters,
// and with the dates apiv2 puts in when v1 gives none.
func managedDatabaseIntegrationTestAnswer(status string) string {
	return `{"id": "i1", "sourceServiceId": "` + managedDatabasePostgresqlTestId + `", "destinationServiceId": "dst", "type": "GRAFANA_DATASOURCE",
		"status": "` + status + `", "checksum": "c1", "createdAt": "1970-01-01T00:00:00Z", "updatedAt": "1970-01-01T00:00:00Z"}`
}

func newManagedDatabaseIntegrationTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*cloudManagedDatabaseIntegrationResource, func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabaseIntegrationResource().(*cloudManagedDatabaseIntegrationResource)
	r.config = config
	return r, bodies
}

func managedDatabaseIntegrationTestParameters(value string) types.Map {
	return types.MapValueMust(types.StringType, map[string]attr.Value{"datasource_name": types.StringValue(value)})
}

func managedDatabaseIntegrationTestPlan() CloudManagedDatabaseIntegrationModel {
	return CloudManagedDatabaseIntegrationModel{
		ServiceName:          types.StringValue("pj"),
		Engine:               types.StringValue("postgresql"),
		ServiceId:            types.StringValue(managedDatabasePostgresqlTestId),
		SourceServiceId:      types.StringValue(managedDatabasePostgresqlTestId),
		DestinationServiceId: types.StringValue("dst"),
		Type:                 types.StringValue("GRAFANA_DATASOURCE"),
		Parameters:           managedDatabaseIntegrationTestParameters("pg"),
		Id:                   types.StringUnknown(),
		Status:               types.StringUnknown(),
		Checksum:             types.StringUnknown(),
	}
}

// managedDatabaseIntegrationTestState is the integration once created, with the parameters of the plan.
func managedDatabaseIntegrationTestState() CloudManagedDatabaseIntegrationModel {
	state := managedDatabaseIntegrationTestPlan()
	state.Id = types.StringValue("i1")
	state.Status = types.StringValue("READY")
	state.Checksum = types.StringValue("c1")
	return state
}

func TestCloudManagedDatabasePostgresqlIntegrationCreate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseIntegrationTestResource(t, map[string][]managedDatabaseTestAnswer{
		// v1 answers a 5xx of Aiven with this 409.
		"POST " + managedDatabaseIntegrationTestPath:    {{http.StatusConflict, `{"class": "Client::Conflict::SystemUnavailable", "message": "system is unavailable"}`}, {http.StatusOK, managedDatabaseIntegrationTestAnswer("CREATING")}},
		"GET " + managedDatabaseIntegrationTestItemPath: {{http.StatusOK, managedDatabaseIntegrationTestAnswer("CREATING")}, {http.StatusOK, managedDatabaseIntegrationTestAnswer("READY")}},
	})

	plan := managedDatabaseIntegrationTestPlan()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	posts := bodies("POST " + managedDatabaseIntegrationTestPath)
	if len(posts) != 2 {
		t.Fatalf("expected the POST to be retried once after the 409, got %d POST", len(posts))
	}
	assertManagedDatabaseJSON(t, json.RawMessage(posts[1]), `{
		"sourceServiceId": "`+managedDatabasePostgresqlTestId+`",
		"destinationServiceId": "dst",
		"type": "GRAFANA_DATASOURCE",
		"parameters": {"datasource_name": "pg"}
	}`)
	var got CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Status.ValueString() != "READY" || !got.Parameters.Equal(plan.Parameters) || !resp.State.Raw.IsFullyKnown() {
		t.Errorf("unexpected state: %v", resp.State.Raw)
	}
}

func TestCloudManagedDatabasePostgresqlIntegrationRead(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseIntegrationTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabaseIntegrationTestItemPath: {{http.StatusOK, managedDatabaseIntegrationTestAnswer("READY")}},
	})

	state := managedDatabaseIntegrationTestState()
	resp := fwresource.ReadResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Read(ctx, fwresource.ReadRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	// The API doesn't return the parameters: the state keeps them.
	var got CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if !got.Parameters.Equal(state.Parameters) {
		t.Errorf("unexpected parameters %s", got.Parameters)
	}
}

func TestCloudManagedDatabasePostgresqlIntegrationModifyPlan(t *testing.T) {
	ctx := context.Background()
	r := NewCloudManagedDatabaseIntegrationResource().(*cloudManagedDatabaseIntegrationResource)

	for name, tc := range map[string]struct {
		current, planned types.Map
		replace          bool
	}{
		"unchanged": {managedDatabaseIntegrationTestParameters("pg"), managedDatabaseIntegrationTestParameters("pg"), false},
		"changed":   {managedDatabaseIntegrationTestParameters("pg"), managedDatabaseIntegrationTestParameters("other"), true},
		"removed":   {managedDatabaseIntegrationTestParameters("pg"), types.MapNull(types.StringType), true},
		// An import has no parameters: setting them adopts the integration.
		"imported": {types.MapNull(types.StringType), managedDatabaseIntegrationTestParameters("pg"), false},
	} {
		t.Run(name, func(t *testing.T) {
			state := managedDatabaseIntegrationTestState()
			state.Parameters = tc.current
			plan := state
			plan.Parameters = tc.planned
			resp := fwresource.ModifyPlanResponse{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}
			r.ModifyPlan(ctx, fwresource.ModifyPlanRequest{
				Plan:  managedDatabaseObjectTestPlan(t, r, &plan),
				State: managedDatabaseObjectTestState(t, r, &state),
			}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if replace := len(resp.RequiresReplace) > 0; replace != tc.replace {
				t.Errorf("unexpected RequiresReplace %v", resp.RequiresReplace)
			}
		})
	}
}

// After an import, setting the parameters only records them.
func TestCloudManagedDatabasePostgresqlIntegrationUpdate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseIntegrationTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabaseIntegrationTestItemPath: {{http.StatusOK, managedDatabaseIntegrationTestAnswer("READY")}},
	})

	state := managedDatabaseIntegrationTestState()
	state.Parameters = types.MapNull(types.StringType)
	plan := managedDatabaseIntegrationTestPlan()
	plan.Id = state.Id
	resp := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Update(ctx, fwresource.UpdateRequest{
		Plan:  managedDatabaseObjectTestPlan(t, r, &plan),
		State: managedDatabaseObjectTestState(t, r, &state),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	for _, request := range bodies("") {
		if method, _, _ := strings.Cut(request, " "); method != http.MethodGet {
			t.Errorf("unexpected request %s", request)
		}
	}
	var got CloudManagedDatabaseIntegrationModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if !got.Parameters.Equal(plan.Parameters) || got.Status.ValueString() != "READY" || !resp.State.Raw.IsFullyKnown() {
		t.Errorf("unexpected state: %v", resp.State.Raw)
	}
}

// v1 refuses an integration of a service that is neither its source nor its destination.
func TestCloudManagedDatabasePostgresqlIntegrationValidateConfig(t *testing.T) {
	ctx := context.Background()
	r := NewCloudManagedDatabaseIntegrationResource().(*cloudManagedDatabaseIntegrationResource)

	for name, tc := range map[string]struct {
		source, destination string
		valid               bool
	}{
		"source":      {managedDatabasePostgresqlTestId, "dst", true},
		"destination": {"src", managedDatabasePostgresqlTestId, true},
		"neither":     {"src", "dst", false},
	} {
		t.Run(name, func(t *testing.T) {
			model := managedDatabaseIntegrationTestPlan()
			model.SourceServiceId = types.StringValue(tc.source)
			model.DestinationServiceId = types.StringValue(tc.destination)
			plan := managedDatabaseObjectTestPlan(t, r, &model)
			var resp fwresource.ValidateConfigResponse
			r.ValidateConfig(ctx, fwresource.ValidateConfigRequest{Config: tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}}, &resp)
			if resp.Diagnostics.HasError() == tc.valid {
				t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlIntegrationDelete(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseIntegrationTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseIntegrationTestItemPath: {{http.StatusOK, managedDatabaseIntegrationTestAnswer("DELETING")}},
		"GET " + managedDatabaseIntegrationTestItemPath:    {{http.StatusOK, managedDatabaseIntegrationTestAnswer("DELETING")}, {http.StatusNotFound, `{"message": "integration not found"}`}},
	})

	state := managedDatabaseIntegrationTestState()
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

// testAccCloudManagedIntegrationTestCase integrates the tested service of the engine with the existing
// service of its INTEGRATION_DESTINATION setting, by an integration of its INTEGRATION_TYPE setting,
// which must link the engine to the engine of the destination, with the parameters of its optional
// INTEGRATION_PARAMETERS setting, a map such as { topic = "logs" }.
func testAccCloudManagedIntegrationTestCase(t *testing.T, category managedDatabaseCategory, engine string) resource.TestCase {
	name := acctest.RandomWithPrefix(test_prefix)
	destination, integrationType := testAccCloudManagedServiceEnv(category, engine, "INTEGRATION_DESTINATION"), testAccCloudManagedServiceEnv(category, engine, "INTEGRATION_TYPE")
	service := "ovh" + category.typeName("")
	integration := service + "_integration.integration"
	parameters := "null"
	if value := os.Getenv(testAccCloudManagedServiceEnv(category, engine, "INTEGRATION_PARAMETERS")); value != "" {
		parameters = value
	}

	return resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudManagedService(t, category, engine)
			checkEnvOrSkip(t, destination)
			checkEnvOrSkip(t, integrationType)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedServiceConfig(category, engine, name, fmt.Sprintf("  description = %q", name), "") + fmt.Sprintf(`
resource "%[1]s_integration" "integration" {
  service_name           = %[1]s.db.service_name
  engine                 = %[1]s.db.engine
  service_id             = %[1]s.db.id
  source_service_id      = %[1]s.db.id
  destination_service_id = %[2]q
  type                   = %[3]q
  parameters             = %[4]s
}
`, service, os.Getenv(destination), os.Getenv(integrationType), parameters),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integration, "status", "READY"),
					resource.TestCheckResourceAttrSet(integration, "id"),
				),
			},
			{
				ResourceName:      integration,
				ImportState:       true,
				ImportStateVerify: true,
				// The API doesn't return the parameters.
				ImportStateVerifyIgnore: []string{"parameters"},
				ImportStateIdFunc:       testAccCloudManagedImportStateIdFunc(integration),
			},
		},
	}
}

func TestAccCloudManagedDatabasePostgresqlIntegration_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedIntegrationTestCase(t, managedDatabaseFamily, "postgresql"))
}

func TestAccCloudManagedAnalyticsKafkaIntegration_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedIntegrationTestCase(t, managedAnalyticsFamily, "kafka"))
}
