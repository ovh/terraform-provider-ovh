package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

const (
	testAccCloudManagedDatabasePostgresqlResource = "ovh_cloud_managed_database.db"
	managedDatabasePostgresqlTestPath             = "/v2/publicCloud/project/pj/managedDatabase/postgresql"
	managedDatabasePostgresqlTestItemPath         = managedDatabasePostgresqlTestPath + "/" + managedDatabasePostgresqlTestId
	managedDatabasePostgresqlTestBusy             = `{"class": "Client::Conflict::ServiceOperationInProgress", "message": "An operation is in progress"}`
	managedDatabasePostgresqlTestCreating         = `{"id": "` + managedDatabasePostgresqlTestId + `", "resourceStatus": "CREATING", "targetSpec": {"version": "17", "plan": "business", "flavor": "b3-8"}}`
)

func init() {
	for _, category := range []managedDatabaseCategory{managedDatabaseFamily, managedAnalyticsFamily} {
		name := "ovh" + category.typeName("")
		resource.AddTestSweepers(name, &resource.Sweeper{
			Name: name,
			F: func(region string) error {
				for _, engine := range category.engines {
					if err := testSweepCloudManagedServices(region, category, engine); err != nil {
						return err
					}
				}
				return nil
			},
		})
	}
}

type managedDatabaseTestAnswer struct {
	status int
	body   string
}

// newManagedDatabaseServiceTestResource returns the service resource of the family, talking to the fake API
// of newManagedDatabaseTestConfig.
func newManagedDatabaseServiceTestResource(t *testing.T, category managedDatabaseCategory, answers map[string][]managedDatabaseTestAnswer) (*cloudManagedDatabaseServiceResource, func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := &cloudManagedDatabaseServiceResource{category: category}
	r.config = config
	return r, bodies
}

// newManagedDatabaseTestConfig returns a provider configuration talking to a fake API that serves
// the answers of each "METHOD path" in turn, the last one repeating, and those of "*" to the requests
// that have none. The returned function gives the bodies received for a "METHOD path", and for "" the
// "METHOD path" of every request, in order.
func newManagedDatabaseTestConfig(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*Config, func(string) []string) {
	t.Helper()

	delay, minTimeout := managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout
	managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout = time.Millisecond, time.Millisecond
	t.Cleanup(func() {
		managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout = delay, minTimeout
	})

	var mu sync.Mutex
	bodies := map[string][]string{}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	answer := func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		key := r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			key += "?" + r.URL.RawQuery
		}
		body, _ := io.ReadAll(r.Body)
		bodies[key] = append(bodies[key], string(body))
		bodies[""] = append(bodies[""], key)
		queue := answers[key]
		if len(queue) == 0 {
			key, queue = "*", answers["*"]
		}
		if len(queue) == 0 {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		if len(queue) > 1 {
			answers[key] = queue[1:]
		}
		w.WriteHeader(queue[0].status)
		fmt.Fprint(w, queue[0].body)
	}
	mux.HandleFunc("/v2/", answer)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	config := &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}

	return config, func(key string) []string {
		mu.Lock()
		defer mu.Unlock()
		return bodies[key]
	}
}

func managedDatabasePostgresqlTestPlan(t *testing.T, model *CloudManagedDatabaseServiceModel) tfsdk.Plan {
	t.Helper()
	state := managedDatabasePostgresqlTestResourceState(t, model)
	return tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
}

// managedDatabaseServiceTestEngines are services of both families, which share the schema and the bodies
// of the service: only the path, which holds the family and the engine, differs.
var managedDatabaseServiceTestEngines = []struct {
	category managedDatabaseCategory
	engine   string
}{
	{managedDatabaseFamily, "postgresql"},
	{managedDatabaseFamily, "mysql"},
	{managedAnalyticsFamily, "kafka"},
}

func TestCloudManagedDatabaseServiceCreate(t *testing.T) {
	ctx := context.Background()
	for _, tc := range managedDatabaseServiceTestEngines {
		t.Run(tc.engine, func(t *testing.T) {
			path := tc.category.endpoint("pj", tc.engine)
			item := path + "/" + managedDatabasePostgresqlTestId
			protected := strings.Replace(managedDatabasePostgresqlTestResponse, `"deletionProtected": false`, `"deletionProtected": true`, 1)
			r, bodies := newManagedDatabaseServiceTestResource(t, tc.category, map[string][]managedDatabaseTestAnswer{
				"POST " + path: {{http.StatusOK, managedDatabasePostgresqlTestCreating}},
				"GET " + item:  {{http.StatusOK, managedDatabasePostgresqlTestCreating}, {http.StatusOK, protected}},
			})

			// The plan of a configuration that sets the required attributes, a description and the deletion
			// protection.
			plan := managedDatabasePostgresqlNullModel()
			plan.Engine = types.StringValue(tc.engine)
			plan.Region = types.StringValue("eu-west-gra")
			plan.Version = types.StringValue("17")
			plan.Plan = types.StringValue("business")
			plan.Flavor = types.StringValue("b3-8")
			plan.Description = types.StringValue("testacc-terraform-db")
			plan.DeletionProtected = types.BoolValue(true)
			plan.AdvancedConfiguration = types.MapUnknown(types.StringType)
			plan.NodesCount = types.Int64Unknown()
			plan.PrometheusEnabled = types.BoolUnknown()
			plan.Backup = types.ObjectUnknown(managedDatabaseBackupAttrTypes())
			plan.MaintenanceWindow = types.ObjectUnknown(managedDatabaseMaintenanceWindowAttrTypes())
			plan.Network = types.ObjectUnknown(managedDatabaseNetworkAttrTypes())
			plan.Storage = types.ObjectUnknown(managedDatabaseStorageAttrTypes())
			plan.Id = types.StringUnknown()
			plan.Checksum = types.StringUnknown()
			plan.CreatedAt = types.StringUnknown()
			plan.UpdatedAt = types.StringUnknown()
			plan.ResourceStatus = types.StringUnknown()
			plan.CurrentTasks = types.ListUnknown(types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()})
			plan.CurrentState = types.ObjectUnknown(managedDatabaseCurrentStateAttrTypes())

			resp := fwresource.CreateResponse{State: managedDatabasePostgresqlTestResourceState(t, nil)}
			r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabasePostgresqlTestPlan(t, &plan)}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			posts := bodies("POST " + path)
			if len(posts) != 1 {
				t.Fatalf("expected 1 POST, got %d", len(posts))
			}
			assertManagedDatabaseJSON(t, json.RawMessage(posts[0]), `{"targetSpec": {
				"region": "eu-west-gra",
				"version": "17",
				"plan": "business",
				"flavor": "b3-8",
				"description": "testacc-terraform-db",
				"deletionProtected": true
			}}`)

			if !resp.State.Raw.IsFullyKnown() {
				t.Errorf("expected a fully known state, got %v", resp.State.Raw)
			}
			var got cloudManagedDatabaseServiceResourceModel
			resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
			if got.Id.ValueString() != managedDatabasePostgresqlTestId || got.Engine.ValueString() != tc.engine || !got.DeletionProtected.ValueBool() {
				t.Errorf("unexpected state: id %s, engine %s, deletion_protected %s", got.Id, got.Engine, got.DeletionProtected)
			}
		})
	}
}

func TestCloudManagedDatabaseServiceUpdate(t *testing.T) {
	ctx := context.Background()
	for _, tc := range managedDatabaseServiceTestEngines {
		t.Run(tc.engine, func(t *testing.T) {
			item := tc.category.serviceEndpoint("pj", tc.engine, managedDatabasePostgresqlTestId)
			renamed := strings.ReplaceAll(managedDatabasePostgresqlTestResponse, "testacc-terraform-db", "testacc-terraform-renamed")
			// A 409 retries the PUT after a new read of the service.
			current := managedDatabaseTestAnswer{http.StatusOK, managedDatabasePostgresqlTestResponse}
			r, bodies := newManagedDatabaseServiceTestResource(t, tc.category, map[string][]managedDatabaseTestAnswer{
				"GET " + item: {current, current, current, {http.StatusOK, renamed}},
				"PUT " + item: {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, renamed}},
			})

			// advanced_configuration is read-only: the PUT keeps the current one.
			state := managedDatabasePostgresqlTestState(t)
			state.Engine = types.StringValue(tc.engine)
			plan := state
			plan.Description = types.StringValue("testacc-terraform-renamed")
			plan.AdvancedConfiguration = types.MapUnknown(types.StringType)
			plan.Checksum = types.StringUnknown()
			plan.UpdatedAt = types.StringUnknown()
			plan.ResourceStatus = types.StringUnknown()
			plan.CurrentTasks = types.ListUnknown(types.ObjectType{AttrTypes: managedDatabaseTaskAttrTypes()})
			plan.CurrentState = types.ObjectUnknown(managedDatabaseCurrentStateAttrTypes())

			resp := fwresource.UpdateResponse{State: managedDatabasePostgresqlTestResourceState(t, &state)}
			r.Update(ctx, fwresource.UpdateRequest{
				Plan:  managedDatabasePostgresqlTestPlan(t, &plan),
				State: managedDatabasePostgresqlTestResourceState(t, &state),
			}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			puts := bodies("PUT " + item)
			if len(puts) != 2 {
				t.Fatalf("expected the PUT to be retried once after the 409, got %d PUT", len(puts))
			}
			assertManagedDatabaseJSON(t, json.RawMessage(puts[1]), strings.Replace(managedDatabasePostgresqlTestReplace, "testacc-terraform-db", "testacc-terraform-renamed", 1))

			var got cloudManagedDatabaseServiceResourceModel
			resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
			if got.Description.ValueString() != "testacc-terraform-renamed" {
				t.Errorf("unexpected description %s", got.Description)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlUpdatePublic(t *testing.T) {
	ctx := context.Background()
	private := managedDatabasePostgresqlTestResponse
	public := strings.ReplaceAll(private, `"public": false`, `"public": true`)
	renamed := func(service string) string {
		return strings.ReplaceAll(service, "testacc-terraform-db", "testacc-terraform-renamed")
	}
	for name, tc := range map[string]struct {
		before, after string
		public        bool
		calls         string
	}{
		"publish":   {private, renamed(public), true, "PUT PATCH"},
		"unpublish": {public, renamed(private), false, "PATCH PUT"},
	} {
		t.Run(name, func(t *testing.T) {
			// The waits and the PUT read the service as it was, the read after the update as it is.
			before := managedDatabaseTestAnswer{http.StatusOK, tc.before}
			r, bodies := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
				"GET " + managedDatabasePostgresqlTestItemPath:   {before, before, before, before, before, {http.StatusOK, tc.after}},
				"PUT " + managedDatabasePostgresqlTestItemPath:   {{http.StatusOK, tc.after}},
				"PATCH " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, tc.after}},
			})

			state := managedDatabasePostgresqlTestState(t)
			if diags := state.MergeWith(ctx, managedDatabasePostgresqlTestAPIResponse(t, tc.before)); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			// The description changes too, so that a PUT is sent.
			plan := state
			plan.Description = types.StringValue("testacc-terraform-renamed")
			network := state.Network.Attributes()
			network["public"] = types.BoolValue(tc.public)
			plan.Network = types.ObjectValueMust(managedDatabaseNetworkAttrTypes(), network)

			resp := fwresource.UpdateResponse{State: managedDatabasePostgresqlTestResourceState(t, &state)}
			r.Update(ctx, fwresource.UpdateRequest{
				Plan:  managedDatabasePostgresqlTestPlan(t, &plan),
				State: managedDatabasePostgresqlTestResourceState(t, &state),
			}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			var calls []string
			for _, request := range bodies("") {
				if method, _, _ := strings.Cut(request, " "); method != http.MethodGet {
					calls = append(calls, method)
				}
			}
			if got := strings.Join(calls, " "); got != tc.calls {
				t.Errorf("expected %s, got %s", tc.calls, got)
			}
			patches := bodies("PATCH " + managedDatabasePostgresqlTestItemPath)
			if len(patches) != 1 {
				t.Fatalf("expected 1 PATCH, got %d", len(patches))
			}
			assertManagedDatabaseJSON(t, json.RawMessage(patches[0]), fmt.Sprintf(`{"targetSpec": {"network": {"public": %t}}}`, tc.public))

			var got cloudManagedDatabaseServiceResourceModel
			resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
			if !got.Network.Attributes()["public"].Equal(types.BoolValue(tc.public)) {
				t.Errorf("unexpected network %s", got.Network)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlPutNoChange(t *testing.T) {
	// The fake API fails the test on a PUT, which it has no answer for.
	r, _ := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, managedDatabasePostgresqlTestResponse}},
	})
	disable := func(service *CloudManagedDatabaseServiceModel) {
		service.PrometheusEnabled = types.BoolValue(false)
	}
	if err := putManagedDatabaseServiceWhenReady(context.Background(), r.config.OVHClient, managedDatabasePostgresqlTestItemPath, disable, time.Minute); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestCloudManagedDatabasePostgresqlPutOutOfSync(t *testing.T) {
	// The PUT of the same targetSpec brings the OUT_OF_SYNC service back in sync, without waiting for
	// it to be READY first.
	outOfSync := managedDatabaseTestAnswer{http.StatusOK, strings.Replace(managedDatabasePostgresqlTestResponse, `"resourceStatus": "READY"`, `"resourceStatus": "OUT_OF_SYNC"`, 1)}
	r, bodies := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {outOfSync, outOfSync, {http.StatusOK, managedDatabasePostgresqlTestResponse}},
		"PUT " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, managedDatabasePostgresqlTestResponse}},
	})
	unchanged := func(*CloudManagedDatabaseServiceModel) {}
	if err := putManagedDatabaseServiceWhenReady(context.Background(), r.config.OVHClient, managedDatabasePostgresqlTestItemPath, unchanged, time.Minute); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	puts := bodies("PUT " + managedDatabasePostgresqlTestItemPath)
	if len(puts) != 1 {
		t.Fatalf("expected 1 PUT, got %d", len(puts))
	}
	assertManagedDatabaseJSON(t, json.RawMessage(puts[0]), managedDatabasePostgresqlTestReplace)
}

func TestCloudManagedDatabasePostgresqlDelete(t *testing.T) {
	ctx := context.Background()
	deleting := strings.Replace(managedDatabasePostgresqlTestResponse, `"resourceStatus": "READY"`, `"resourceStatus": "DELETING"`, 1)
	r, bodies := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabasePostgresqlTestItemPath: {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, deleting}},
		"GET " + managedDatabasePostgresqlTestItemPath:    {{http.StatusOK, deleting}, {http.StatusNotFound, `{"message": "service not found"}`}},
	})

	state := managedDatabasePostgresqlTestState(t)
	resp := fwresource.DeleteResponse{State: managedDatabasePostgresqlTestResourceState(t, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabasePostgresqlTestResourceState(t, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got := len(bodies("DELETE " + managedDatabasePostgresqlTestItemPath)); got != 2 {
		t.Errorf("expected the DELETE to be retried once after the 409, got %d DELETE", got)
	}
}

func TestCloudManagedDatabasePostgresqlReadNotFound(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusNotFound, `{"message": "service not found"}`}},
	})

	state := managedDatabasePostgresqlTestState(t)
	resp := fwresource.ReadResponse{State: managedDatabasePostgresqlTestResourceState(t, &state)}
	r.Read(ctx, fwresource.ReadRequest{State: managedDatabasePostgresqlTestResourceState(t, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected the service to be removed from the state, got %v", resp.State.Raw)
	}
}

func TestCloudManagedDatabasePostgresqlCreationOnly(t *testing.T) {
	ctx := context.Background()
	model := managedDatabasePostgresqlTestState(t)
	state := managedDatabasePostgresqlTestResourceState(t, &model)
	for name, tc := range map[string]struct {
		state, plan types.String
		replace     bool
	}{
		"set after an import": {types.StringNull(), types.StringValue("b"), false},
		"removed":             {types.StringValue("a"), types.StringNull(), false},
		"changed":             {types.StringValue("a"), types.StringValue("b"), true},
	} {
		t.Run(name, func(t *testing.T) {
			req := planmodifier.StringRequest{
				State:       state,
				Plan:        tfsdk.Plan{Schema: state.Schema, Raw: state.Raw},
				StateValue:  tc.state,
				PlanValue:   tc.plan,
				ConfigValue: tc.plan,
			}
			resp := planmodifier.StringResponse{PlanValue: tc.plan}
			managedDatabaseCreationOnly.PlanModifyString(ctx, req, &resp)
			if resp.RequiresReplace != tc.replace {
				t.Errorf("expected RequiresReplace %t, got %t", tc.replace, resp.RequiresReplace)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlWaitError(t *testing.T) {
	r, _ := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {{http.StatusOK, `{
			"id": "` + managedDatabasePostgresqlTestId + `",
			"resourceStatus": "ERROR",
			"currentTasks": [{"id": "1", "type": "CREATE", "status": "ERROR", "errors": [{"message": "flavor b3-8 is not available in eu-west-gra"}]}]
		}`}},
	})

	err := waitForManagedDatabaseServiceReady(context.Background(), r.config.OVHClient, managedDatabasePostgresqlTestItemPath, time.Minute)
	if err == nil || !strings.Contains(err.Error(), "flavor b3-8 is not available in eu-west-gra") {
		t.Errorf("expected the error of the failed task, got %v", err)
	}
}

func TestCloudManagedDatabasePostgresqlWaitEngineReady(t *testing.T) {
	// A new service answered READY while its engine was still provisioning its second node, then UPDATING.
	ctx := context.Background()
	provisioning := strings.ReplaceAll(managedDatabasePostgresqlTestResponse, `"status": "READY"`, `"status": "UPDATING"`)
	updating := strings.Replace(provisioning, `"resourceStatus": "READY"`, `"resourceStatus": "UPDATING"`, 1)
	r, _ := newManagedDatabaseServiceTestResource(t, managedDatabaseFamily, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabasePostgresqlTestItemPath: {
			{http.StatusOK, managedDatabasePostgresqlTestCreating},
			{http.StatusOK, provisioning},
			{http.StatusOK, updating},
			{http.StatusOK, managedDatabasePostgresqlTestResponse},
		},
	})

	if err := waitForManagedDatabaseServiceReady(ctx, r.config.OVHClient, managedDatabasePostgresqlTestItemPath, time.Minute); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var res CloudManagedDatabaseServiceAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, managedDatabasePostgresqlTestItemPath, &res); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if res.ResourceStatus != "READY" {
		t.Errorf("expected the wait to last until the engine is READY, the service is %s", res.ResourceStatus)
	}
}

// Only the 409s of an operation running on the service are retried: the others would answer the same
// until the timeout.
func TestCloudManagedDatabaseServiceRetryBusy(t *testing.T) {
	for class, retried := range map[string]bool{
		"Client::Conflict::ServiceOperationInProgress": true,
		"Client::Conflict::ServiceNotReady":            true,
		"Client::Conflict::ServiceOnlyOneModification": false,
		"Client::Conflict::ServiceLocked":              false,
	} {
		t.Run(class, func(t *testing.T) {
			calls := 0
			err := retryWhileManagedDatabaseServiceBusy(context.Background(), time.Minute, func() error {
				calls++
				if calls == 1 {
					return &ovh.APIError{Code: http.StatusConflict, Class: class}
				}
				return nil
			})
			if retried && (err != nil || calls != 2) {
				t.Errorf("expected one retry, got %d calls and error %v", calls, err)
			}
			if !retried && (err == nil || calls != 1) {
				t.Errorf("expected the error at once, got %d calls and error %v", calls, err)
			}
		})
	}
}

// newManagedDatabaseLiveLikeTestAPI returns the endpoint of a fake API holding one PostgreSQL service of
// project pj as the live API does: its targetSpec gets a maintenance window, a storage and the
// network when the request leaves them out, and the defaults of the engine in its advanced configuration.
// A PUT changes updatedAt but keeps the public exposure, which a PATCH changes.
// It also returns a function that changes the targetSpec, as done outside of Terraform, and one that
// makes the service OUT_OF_SYNC until the next PUT.
func newManagedDatabaseLiveLikeTestAPI(t *testing.T) (string, func(func(spec map[string]any)), func()) {
	t.Helper()

	delay, minTimeout := managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout
	managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout = time.Millisecond, time.Millisecond
	t.Cleanup(func() {
		managedDatabaseServicePollDelay, managedDatabaseServicePollMinTimeout = delay, minTimeout
	})

	var mu sync.Mutex
	var spec map[string]any
	status, updatedAt := "READY", "2026-11-12T10:20:00Z"
	store := func(w http.ResponseWriter, r *http.Request) bool {
		var body struct {
			TargetSpec map[string]any `json:"targetSpec"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("%s %s: %s", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return false
		}
		advanced := map[string]any{"pg.io_method": "worker", "pg.password_encryption": "scram-sha-256"}
		if requested, ok := body.TargetSpec["advancedConfiguration"].(map[string]any); ok {
			for key, value := range requested {
				advanced[key] = value
			}
		}
		spec = body.TargetSpec
		spec["advancedConfiguration"] = advanced
		delete(spec, "region")
		for key, value := range map[string]any{
			"maintenanceWindow": map[string]any{"startTimeUtc": "07:52"},
			"storage":           map[string]any{"size": map[string]any{"unit": "GB", "value": 160}},
			"network":           map[string]any{"authorizedIps": map[string]any{"rules": []any{}}},
		} {
			if spec[key] == nil {
				spec[key] = value
			}
		}
		return true
	}
	respond := func(w http.ResponseWriter) {
		var service map[string]any
		if err := json.Unmarshal([]byte(managedDatabasePostgresqlTestResponse), &service); err != nil {
			t.Fatal(err)
		}
		state := service["currentState"].(map[string]any)
		for key, value := range spec {
			state[key] = value
		}
		service["targetSpec"] = spec
		service["resourceStatus"] = status
		service["updatedAt"] = updatedAt
		if err := json.NewEncoder(w).Encode(service); err != nil {
			t.Error(err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc("/1.0/auth/details", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"account": "xx1111-ovh", "method": "provider"}`)
	})
	mux.HandleFunc("POST "+managedDatabasePostgresqlTestPath, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if store(w, r) {
			respond(w)
		}
	})
	mux.HandleFunc(managedDatabasePostgresqlTestItemPath, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case spec == nil:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "service not found"}`)
		case r.Method == http.MethodGet:
			respond(w)
		case r.Method == http.MethodPut:
			public := spec["network"].(map[string]any)["public"]
			if store(w, r) {
				spec["network"].(map[string]any)["public"] = public
				status, updatedAt = "READY", time.Now().UTC().Format(time.RFC3339Nano)
				respond(w)
			}
		case r.Method == http.MethodPatch:
			var body struct {
				TargetSpec struct {
					Network struct {
						Public bool `json:"public"`
					} `json:"network"`
				} `json:"targetSpec"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("%s %s: %s", r.Method, r.URL.Path, err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			spec["network"].(map[string]any)["public"] = body.TargetSpec.Network.Public
			respond(w)
		case r.Method == http.MethodDelete:
			respond(w)
			spec = nil
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	edit := func(change func(spec map[string]any)) {
		mu.Lock()
		defer mu.Unlock()
		change(spec)
	}
	outOfSync := func() {
		mu.Lock()
		defer mu.Unlock()
		status = "OUT_OF_SYNC"
	}
	return server.URL + "/1.0", edit, outOfSync
}

// The plan after an apply is empty although the API fills in settings left out of the configuration
// and adds defaults to the advanced configuration. advanced_configuration is read-only: a property set
// outside of Terraform shows without a change, and an update sends it back. Publishing the service and
// taking it off its public address update it in place. An OUT_OF_SYNC service is updated with a PUT,
// which brings it back in sync.
func TestCloudManagedDatabasePostgresqlLiveLikePlan(t *testing.T) {
	endpoint, edit, outOfSync := newManagedDatabaseLiveLikeTestAPI(t)
	config := func(attributes string) string {
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
  flavor             = "b3-8"
  description        = "testacc-terraform-db"
  deletion_protected = true
%s
}
`, endpoint, attributes)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(""),
				Check:  resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "advanced_configuration.pg.io_method", "worker"),
			},
			{
				PreConfig: func() {
					edit(func(spec map[string]any) {
						spec["advancedConfiguration"].(map[string]any)["pg.log_temp_files"] = "1024"
					})
				},
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "advanced_configuration.pg.log_temp_files", "1024"),
			},
			{
				Config: config("  prometheus_enabled = true"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "advanced_configuration.pg.log_temp_files", "1024"),
			},
			{
				Config: config("  prometheus_enabled = true\n  network = { public = true }"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.public", "true"),
			},
			{
				Config: config("  prometheus_enabled = false\n  network = { public = false }"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.public", "false"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "prometheus_enabled", "false"),
				),
			},
			{
				PreConfig: outOfSync,
				Config:    config("  prometheus_enabled = false\n  network = { public = false }"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "resource_status", "READY"),
			},
			// A configuration that sets what the API fills in, so that Terraform proposes no change of its own.
			{
				PreConfig: outOfSync,
				Config: config(`  prometheus_enabled = false
  maintenance_window = { start_time_utc = "07:52" }
  storage            = { size = { unit = "GB", value = 160 } }
  network            = { authorized_ips = { rules = [] }, public = false }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "resource_status", "READY"),
			},
		},
	})
}

// A value of the configuration unknown until the apply, such as the ID of a subnet being replaced,
// stays unknown in the plan: the service is recreated on the new subnet.
func TestCloudManagedDatabasePostgresqlUnknownConfigPlan(t *testing.T) {
	endpoint, _, _ := newManagedDatabaseLiveLikeTestAPI(t)
	config := func(subnet string) string {
		return fmt.Sprintf(`
provider "ovh" {
  endpoint           = %q
  application_key    = "key"
  application_secret = "secret"
  consumer_key       = "consumer"
}

resource "terraform_data" "subnet" {
  triggers_replace = %q
}

resource "ovh_cloud_managed_database" "db" {
  service_name = "pj"
  engine       = "postgresql"
  region       = "eu-west-gra"
  version      = "17"
  plan         = "business"
  flavor       = "b3-8"
  description  = "testacc-terraform-db"
  network = {
    network_id = "net"
    subnet_id  = terraform_data.subnet.id
  }
}
`, endpoint, subnet)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("subnet-1")},
			{
				Config: config("subnet-2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttrPair(testAccCloudManagedDatabasePostgresqlResource, "network.subnet_id", "terraform_data.subnet", "id"),
			},
		},
	})
}

// testSweepCloudManagedServices deletes the services of the engine left by an aborted acceptance run.
func testSweepCloudManagedServices(region string, category managedDatabaseCategory, engine string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		return nil
	}

	ctx := context.Background()
	endpoint := category.endpoint(serviceName, engine)
	services, err := helpers.GetAllPagesV2[CloudManagedDatabaseServiceAPIResponse](ctx, client, endpoint)
	if err != nil {
		return fmt.Errorf("error listing %s services: %s", engine, err)
	}

	for _, service := range services {
		if service.TargetSpec == nil || service.TargetSpec.Description == nil || !strings.HasPrefix(*service.TargetSpec.Description, test_prefix) {
			continue
		}
		item := endpoint + "/" + url.PathEscape(service.Id)
		if service.TargetSpec.DeletionProtected != nil && *service.TargetSpec.DeletionProtected {
			unprotect := func(service *CloudManagedDatabaseServiceModel) {
				service.DeletionProtected = types.BoolValue(false)
			}
			if err := putManagedDatabaseServiceWhenReady(ctx, client, item, unprotect, managedDatabaseServiceDeleteTimeout); err != nil {
				return fmt.Errorf("error removing the deletion protection of %s service %s: %s", engine, service.Id, err)
			}
		}
		if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseServiceDeleteTimeout, func() error {
			return client.DeleteWithContext(ctx, item, nil)
		}); err != nil && !helpers.IsNotFound(err) {
			return fmt.Errorf("error deleting %s service %s: %s", engine, service.Id, err)
		}
	}

	return nil
}

func testAccCheckCloudManagedServiceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		var category managedDatabaseCategory
		switch rs.Type {
		case "ovh_cloud_managed_database":
			category = managedDatabaseFamily
		case "ovh_cloud_managed_analytics":
			category = managedAnalyticsFamily
		default:
			continue
		}

		engine, id := rs.Primary.Attributes["engine"], rs.Primary.Attributes["id"]
		err := testAccOVHClient.Get(category.serviceEndpoint(rs.Primary.Attributes["service_name"], engine, id), &CloudManagedDatabaseServiceAPIResponse{})
		if err == nil {
			return fmt.Errorf("%s service %s still exists", engine, id)
		}
		if !helpers.IsNotFound(err) {
			return fmt.Errorf("error checking that %s service %s was destroyed: %w", engine, id, err)
		}
	}
	return nil
}

// testAccCloudManagedServiceEnv returns the variable OVH_CLOUD_MANAGED_<FAMILY>_<ENGINE>_<SETTING>_TEST
// that configures the tested service of the engine, such as KAFKA_CONNECT for kafkaConnect.
func testAccCloudManagedServiceEnv(category managedDatabaseCategory, engine, setting string) string {
	engine = regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(engine, "${1}_${2}")
	return strings.ToUpper("OVH_CLOUD_MANAGED_" + category.family + "_" + engine + "_" + setting + "_TEST")
}

// testAccCloudManagedServiceSetting returns the value of the variable of the setting, or its default:
// MongoDB runs on the production plan, as in the v1 MongoDB tests, since db2-free, the flavor of the
// discovery plan, is end of sale.
func testAccCloudManagedServiceSetting(category managedDatabaseCategory, engine, setting string) string {
	if value := os.Getenv(testAccCloudManagedServiceEnv(category, engine, setting)); value != "" {
		return value
	}
	if engine == "mongodb" && setting == "PLAN" {
		return "production"
	}
	return ""
}

func testAccPreCheckCloudManagedService(t *testing.T, category managedDatabaseCategory, engine string) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	for _, setting := range []string{"VERSION", "PLAN", "FLAVOR", "REGION"} {
		if testAccCloudManagedServiceSetting(category, engine, setting) == "" {
			checkEnvOrSkip(t, testAccCloudManagedServiceEnv(category, engine, setting))
		}
	}
}

func testAccPreCheckCloudManagedDatabasePostgresql(t *testing.T) {
	testAccPreCheckCloudManagedService(t, managedDatabaseFamily, "postgresql")
}

func testAccCloudManagedDatabasePostgresqlConfig(name, attributes, networkAttributes string) string {
	return testAccCloudManagedServiceConfig(managedDatabaseFamily, "postgresql", name, attributes, networkAttributes)
}

// testAccCloudManagedServiceConfig renders the service ovh_cloud_managed_<family>.db of the engine with the
// given attributes and network attributes. When the NETWORK_REGION setting of the engine is set, the
// service is attached to a private network created in that region, which native services require. The
// NODES_COUNT setting sets nodes_count, which the API requires on a plan of more than one node,
// although its spec says it applies the minimum of the plan.
func testAccCloudManagedServiceConfig(category managedDatabaseCategory, engine, name, attributes, networkAttributes string) string {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	setting := func(key string) string {
		return testAccCloudManagedServiceSetting(category, engine, key)
	}

	var config, network string
	if networkRegion := setting("NETWORK_REGION"); networkRegion != "" {
		config = fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = %[1]q
  name         = %[2]q
  region       = %[3]q
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = %[2]q
  cidr         = "10.0.0.0/24"
  region       = %[3]q
}
`, serviceName, name, networkRegion)
		network = `
    network_id = ovh_cloud_network_private_vrack.network.id
    subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id`
	}
	if nodesCount := setting("NODES_COUNT"); nodesCount != "" {
		attributes += "\n  nodes_count = " + nodesCount
	}
	if networkAttributes != "" {
		network += "\n    " + networkAttributes
	}
	if network != "" {
		attributes += "\n  network = {" + network + "\n  }"
	}

	return config + fmt.Sprintf(`
resource "ovh%s" "db" {
  service_name = %q
  engine       = %q
  region       = %q
  version      = %q
  plan         = %q
  flavor       = %q
%s
}
`,
		category.typeName(""),
		serviceName,
		engine,
		setting("REGION"),
		setting("VERSION"),
		setting("PLAN"),
		setting("FLAVOR"),
		attributes,
	)
}

// testAccCloudManagedEngineCapabilities is left out of the import checks of a service, since the API
// changes them between two reads: it drops the region of the service from the backup regions shortly
// after the service is ready, and sometimes leaves out the advanced configuration properties.
const testAccCloudManagedEngineCapabilities = "current_state.engine_capabilities"

// testAccCloudManagedImportStateIdFunc returns the import ID of a service, <service_name>/<engine>/<id>,
// or of an object of a service, <service_name>/<engine>/<service_id>/<id>: without the engine when the
// type name holds it, and without the id for a Prometheus endpoint, which has none.
func testAccCloudManagedImportStateIdFunc(address string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[address]
		if !ok {
			return "", fmt.Errorf("not found: %s", address)
		}
		attributes := rs.Primary.Attributes
		id := []string{attributes["service_name"]}
		if !strings.Contains(rs.Type, "_"+attributes["engine"]+"_") {
			id = append(id, attributes["engine"])
		}
		if serviceId, ok := attributes["service_id"]; ok {
			id = append(id, serviceId)
		}
		if !strings.HasSuffix(rs.Type, "_prometheus") {
			id = append(id, attributes["id"])
		}
		return strings.Join(id, "/"), nil
	}
}

func TestAccCloudManagedDatabasePostgresql_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "description", name),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "region", os.Getenv("OVH_CLOUD_MANAGED_DATABASE_POSTGRESQL_REGION_TEST")),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "resource_status", "READY"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlResource, "id"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlResource, "nodes_count"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlResource, "storage.size.value"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "current_state.endpoints.0.component", "POSTGRESQL"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlResource, "current_state.endpoints.0.uri"),
				),
			},
			{
				ResourceName:            testAccCloudManagedDatabasePostgresqlResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{testAccCloudManagedEngineCapabilities},
				ImportStateIdFunc:       testAccCloudManagedImportStateIdFunc(testAccCloudManagedDatabasePostgresqlResource),
			},
		},
	})
}

func TestAccCloudManagedDatabasePostgresql_update(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	expectUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlResource, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf(`
  description        = %q
  deletion_protected = true
  prometheus_enabled = false
  backup             = { time = "02:00" }
  maintenance_window = { start_time_utc = "03:00" }`, name), `authorized_ips = { rules = [{ cidr = "203.0.113.0/24", description = "office" }] }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "deletion_protected", "true"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "backup.time", "02:00"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "maintenance_window.start_time_utc", "03:00"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.authorized_ips.rules.#", "1"),
				),
			},
			{
				Config: testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf(`
  description        = "%s-updated"
  deletion_protected = false
  prometheus_enabled = true
  backup             = { time = "04:00" }
  maintenance_window = { start_time_utc = "05:00" }`, name), `authorized_ips = { rules = [{ cidr = "203.0.113.0/24", description = "office" }, { cidr = "198.51.100.0/24" }] }
    public         = true`),
				ConfigPlanChecks: expectUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "description", name+"-updated"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "deletion_protected", "false"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "prometheus_enabled", "true"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "backup.time", "04:00"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "maintenance_window.start_time_utc", "05:00"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.authorized_ips.rules.#", "2"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.public", "true"),
				),
			},
			{
				Config: testAccCloudManagedDatabasePostgresqlConfig(name, "", `authorized_ips = { rules = [{ cidr = "198.51.100.0/24" }] }
    public         = false`),
				ConfigPlanChecks: expectUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "description"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.authorized_ips.rules.#", "1"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlResource, "network.public", "false"),
				),
			},
		},
	})
}

// testAccCloudManagedEngineTestCase tests a service of the engine. config(name, step) renders, at step
// v1 then v2, the service ovh_cloud_managed_<family>.db described as <name>-<step>, and the data source
// of the service, named db. checks are the engine's own checks of v1, and v2 must update the service in
// place. The import of the service is verified last.
func testAccCloudManagedEngineTestCase(t *testing.T, category managedDatabaseCategory, engine string, config func(name, step string) string, checks ...resource.TestCheckFunc) resource.TestCase {
	name := acctest.RandomWithPrefix(test_prefix)
	service := "ovh" + category.typeName("") + ".db"
	checks = append(checks, resource.TestCheckResourceAttr(service, "resource_status", "READY"))
	for _, attribute := range []string{"region", "version", "plan", "flavor", "nodes_count", "description"} {
		checks = append(checks, resource.TestCheckResourceAttrPair("data."+service, attribute, service, attribute))
	}

	return resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedService(t, category, engine) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config(name, "v1"),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: config(name, "v2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(service, plancheck.ResourceActionUpdate)},
				},
				Check: resource.TestCheckResourceAttr(service, "description", name+"-v2"),
			},
			{
				ResourceName:            service,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{testAccCloudManagedEngineCapabilities},
				ImportStateIdFunc:       testAccCloudManagedImportStateIdFunc(service),
			},
		},
	}
}

func testAccCloudManagedDatabaseMongodbConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedDatabaseFamily, "mongodb", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_database" "db" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  id           = ovh_cloud_managed_database.db.id
}
`
}

func TestAccCloudManagedDatabaseMongodb_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedDatabaseFamily, "mongodb", testAccCloudManagedDatabaseMongodbConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_database.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedDatabaseValkeyConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedDatabaseFamily, "valkey", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_database" "db" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  id           = ovh_cloud_managed_database.db.id
}
`
}

func TestAccCloudManagedDatabaseValkey_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedDatabaseFamily, "valkey", testAccCloudManagedDatabaseValkeyConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_database.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsKafkaConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "kafka", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsKafka_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "kafka", testAccCloudManagedAnalyticsKafkaConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_analytics.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsOpensearchConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "opensearch", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsOpensearch_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "opensearch", testAccCloudManagedAnalyticsOpensearchConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_analytics.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsClickhouseConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "clickhouse", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsClickhouse_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "clickhouse", testAccCloudManagedAnalyticsClickhouseConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_analytics.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedDatabaseMysqlConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedDatabaseFamily, "mysql", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_database" "db" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  id           = ovh_cloud_managed_database.db.id
}
`
}

func TestAccCloudManagedDatabaseMysql_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedDatabaseFamily, "mysql", testAccCloudManagedDatabaseMysqlConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_database.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsGrafanaConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "grafana", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsGrafana_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "grafana", testAccCloudManagedAnalyticsGrafanaConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_analytics.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsKafkaConnectConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "kafkaConnect", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsKafkaConnect_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "kafkaConnect", testAccCloudManagedAnalyticsKafkaConnectConfig,
		resource.TestCheckResourceAttrSet("ovh_cloud_managed_analytics.db", "current_state.endpoints.0.uri"),
	))
}

func testAccCloudManagedAnalyticsKafkaMirrorMakerConfig(name, step string) string {
	return testAccCloudManagedServiceConfig(managedAnalyticsFamily, "kafkaMirrorMaker", name, fmt.Sprintf("  description = %q", name+"-"+step), "") + `
data "ovh_cloud_managed_analytics" "db" {
  service_name = ovh_cloud_managed_analytics.db.service_name
  engine       = ovh_cloud_managed_analytics.db.engine
  id           = ovh_cloud_managed_analytics.db.id
}
`
}

func TestAccCloudManagedAnalyticsKafkaMirrorMaker_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedEngineTestCase(t, managedAnalyticsFamily, "kafkaMirrorMaker", testAccCloudManagedAnalyticsKafkaMirrorMakerConfig))
}
