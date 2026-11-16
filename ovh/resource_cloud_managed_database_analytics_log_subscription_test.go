package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const managedDatabaseLogSubscriptionTestPath = managedDatabasePostgresqlTestItemPath + "/log/subscription"

func newManagedDatabaseLogSubscriptionTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*cloudManagedDatabaseLogSubscriptionResource, func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabaseLogSubscriptionResource().(*cloudManagedDatabaseLogSubscriptionResource)
	r.config = config
	return r, bodies
}

func managedDatabaseLogSubscriptionTestPlan() CloudManagedDatabaseLogSubscriptionModel {
	return CloudManagedDatabaseLogSubscriptionModel{
		ServiceName:    types.StringValue("pj"),
		Engine:         types.StringValue("postgresql"),
		ServiceId:      types.StringValue(managedDatabasePostgresqlTestId),
		Kind:           types.StringValue("customer_logs"),
		StreamId:       types.StringValue("st1"),
		Id:             types.StringUnknown(),
		LdpServiceName: types.StringUnknown(),
		Resource:       types.ObjectUnknown(managedDatabaseLogSubscriptionResourceAttrTypes()),
		CreatedAt:      types.StringUnknown(),
		UpdatedAt:      types.StringUnknown(),
	}
}

const (
	managedDatabaseLogSubscriptionTestList      = "GET " + managedDatabaseLogSubscriptionTestPath + "?kind=CUSTOMER_LOGS"
	managedDatabaseLogSubscriptionTestOperation = "GET /1.0/dbaas/logs/ldp-xx-12345/operation/op1"
	managedDatabaseLogSubscriptionTestPending   = `{"operationId": "op1", "state": "RUNNING"}`
	managedDatabaseLogSubscriptionTestDone      = `{"operationId": "op1", "state": "SUCCESS", "subscriptionId": "s1", "createdAt": "2026-11-12T10:00:00Z", "updatedAt": "2026-11-12T10:00:00Z"}`
)

func managedDatabaseLogSubscriptionTestAnswer(id, stream string) string {
	return `{"subscriptionId": "` + id + `", "kind": "customer_logs", "streamId": "` + stream + `", "serviceName": "ldp-xx-12345",
		"resource": {"name": "` + managedDatabasePostgresqlTestId + `", "type": "postgresql"}, "createdAt": "2026-11-12T10:00:00Z", "updatedAt": "2026-11-12T10:00:00Z"}`
}

func TestCloudManagedDatabasePostgresqlLogSubscriptionCreate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseLogSubscriptionTestResource(t, map[string][]managedDatabaseTestAnswer{
		managedDatabaseLogSubscriptionTestList:           {{http.StatusOK, `[]`}},
		"POST " + managedDatabaseLogSubscriptionTestPath: {{http.StatusOK, `{"operationId": "op1", "serviceName": "ldp-xx-12345"}`}},
		// A 5xx is one more read.
		managedDatabaseLogSubscriptionTestOperation: {
			{http.StatusOK, managedDatabaseLogSubscriptionTestPending},
			{http.StatusBadGateway, `{"message": "bad gateway"}`},
			{http.StatusOK, managedDatabaseLogSubscriptionTestDone},
		},
		"GET " + managedDatabaseLogSubscriptionTestPath + "/s1": {{http.StatusOK, managedDatabaseLogSubscriptionTestAnswer("s1", "st1")}},
	})

	plan := managedDatabaseLogSubscriptionTestPlan()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	start := time.Now()
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("expected the wait to poll at the test pace, took %s", elapsed)
	}

	assertManagedDatabaseJSON(t, json.RawMessage(bodies("POST " + managedDatabaseLogSubscriptionTestPath)[0]), `{"kind": "customer_logs", "streamId": "st1"}`)
	var got CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != "s1" || got.LdpServiceName.ValueString() != "ldp-xx-12345" || !resp.State.Raw.IsFullyKnown() {
		t.Errorf("unexpected state: %v", resp.State.Raw)
	}
}

// A subscription of the kind to the stream, left by a creation that failed after its POST, is adopted
// instead of forwarding the logs twice.
func TestCloudManagedDatabasePostgresqlLogSubscriptionCreateAdopts(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseLogSubscriptionTestResource(t, map[string][]managedDatabaseTestAnswer{
		managedDatabaseLogSubscriptionTestList:                  {{http.StatusOK, `["s0", "s1"]`}},
		"GET " + managedDatabaseLogSubscriptionTestPath + "/s0": {{http.StatusOK, managedDatabaseLogSubscriptionTestAnswer("s0", "other")}},
		"GET " + managedDatabaseLogSubscriptionTestPath + "/s1": {{http.StatusOK, managedDatabaseLogSubscriptionTestAnswer("s1", "st1")}},
	})

	plan := managedDatabaseLogSubscriptionTestPlan()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if posts := bodies("POST " + managedDatabaseLogSubscriptionTestPath); len(posts) != 0 {
		t.Errorf("expected no POST, got %d", len(posts))
	}
	var got CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != "s1" {
		t.Errorf("unexpected state: %v", resp.State.Raw)
	}
}

// When the wait for the operation fails, here for want of rights on the Logs Data Platform service, the
// subscription it created is saved, so that it stays tracked.
func TestCloudManagedDatabasePostgresqlLogSubscriptionCreateWaitFails(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseLogSubscriptionTestResource(t, map[string][]managedDatabaseTestAnswer{
		managedDatabaseLogSubscriptionTestList:                  {{http.StatusOK, `[]`}, {http.StatusOK, `["s1"]`}},
		"POST " + managedDatabaseLogSubscriptionTestPath:        {{http.StatusOK, `{"operationId": "op1", "serviceName": "ldp-xx-12345"}`}},
		managedDatabaseLogSubscriptionTestOperation:             {{http.StatusForbidden, `{"message": "This call has not been granted"}`}},
		"GET " + managedDatabaseLogSubscriptionTestPath + "/s1": {{http.StatusOK, managedDatabaseLogSubscriptionTestAnswer("s1", "st1")}},
	})

	plan := managedDatabaseLogSubscriptionTestPlan()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the failed wait to be reported")
	}
	var got CloudManagedDatabaseLogSubscriptionModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != "s1" {
		t.Errorf("expected the subscription in the state, got %v", resp.State.Raw)
	}
}

func TestCloudManagedDatabasePostgresqlLogSubscriptionDelete(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseLogSubscriptionTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseLogSubscriptionTestPath + "/s1": {{http.StatusOK, `{"operationId": "op1", "serviceName": "ldp-xx-12345"}`}},
		managedDatabaseLogSubscriptionTestOperation:                {{http.StatusOK, managedDatabaseLogSubscriptionTestPending}, {http.StatusOK, managedDatabaseLogSubscriptionTestDone}},
	})

	state := managedDatabaseLogSubscriptionTestPlan()
	if diags := state.MergeWith(ctx, &managedDatabaseLogSubscription{SubscriptionId: "s1", Kind: "customer_logs", StreamId: "st1", ServiceName: "ldp-xx-12345"}); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got := len(bodies(managedDatabaseLogSubscriptionTestOperation)); got != 2 {
		t.Errorf("expected the operation to be read until it succeeds, got %d reads", got)
	}
}

func TestCloudManagedDatabasePostgresqlLogSubscriptionDeleteGone(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseLogSubscriptionTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseLogSubscriptionTestPath + "/s1": {{http.StatusNotFound, `{"message": "subscription not found"}`}},
	})

	state := managedDatabaseLogSubscriptionTestPlan()
	if diags := state.MergeWith(ctx, &managedDatabaseLogSubscription{SubscriptionId: "s1", Kind: "customer_logs", StreamId: "st1", ServiceName: "ldp-xx-12345"}); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

// testAccCloudManagedLogSubscriptionTestCase forwards the logs of the tested service of the engine to a
// stream created in the Logs Data Platform service of OVH_DBAAS_LOGS_SERVICE_TEST.
func testAccCloudManagedLogSubscriptionTestCase(t *testing.T, category managedDatabaseCategory, engine string) resource.TestCase {
	name := acctest.RandomWithPrefix(test_prefix)
	service := "ovh" + category.typeName("")
	subscription := service + "_log_subscription.subscription"

	return resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudManagedService(t, category, engine)
			checkEnvOrSkip(t, "OVH_DBAAS_LOGS_SERVICE_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedServiceConfig(category, engine, name, fmt.Sprintf("  description = %q", name), "") + fmt.Sprintf(`
resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = %[2]q
  title        = %[3]q
  description  = %[3]q
}

resource "%[1]s_log_subscription" "subscription" {
  service_name = %[1]s.db.service_name
  engine       = %[1]s.db.engine
  service_id   = %[1]s.db.id
  stream_id    = ovh_dbaas_logs_output_graylog_stream.stream.id
  kind         = "customer_logs"
}

data "%[1]s_log_subscriptions" "subscriptions" {
  service_name = %[1]s.db.service_name
  engine       = %[1]s.db.engine
  service_id   = %[1]s.db.id
  kind         = "customer_logs"
  depends_on   = [%[1]s_log_subscription.subscription]
}

data "%[1]s_log_subscription" "subscription" {
  service_name = %[1]s.db.service_name
  engine       = %[1]s.db.engine
  service_id   = %[1]s.db.id
  id           = %[1]s_log_subscription.subscription.id
}
`, service, os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST"), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(subscription, "kind", "customer_logs"),
					resource.TestCheckResourceAttr(subscription, "ldp_service_name", os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")),
					resource.TestCheckResourceAttrSet(subscription, "id"),
					resource.TestCheckTypeSetElemAttrPair("data."+service+"_log_subscriptions.subscriptions", "subscription_ids.*", subscription, "id"),
					resource.TestCheckResourceAttr("data."+subscription, "kind", "customer_logs"),
				),
			},
			{
				ResourceName:      subscription,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudManagedImportStateIdFunc(subscription),
			},
		},
	}
}

func TestAccCloudManagedDatabasePostgresqlLogSubscription_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedLogSubscriptionTestCase(t, managedDatabaseFamily, "postgresql"))
}

func TestAccCloudManagedAnalyticsKafkaLogSubscription_basic(t *testing.T) {
	resource.Test(t, testAccCloudManagedLogSubscriptionTestCase(t, managedAnalyticsFamily, "kafka"))
}
