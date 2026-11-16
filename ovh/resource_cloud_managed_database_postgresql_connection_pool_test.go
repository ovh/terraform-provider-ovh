package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const (
	testAccCloudManagedDatabasePostgresqlConnectionPoolResource = "ovh_cloud_managed_database_postgresql_connection_pool.pool"
	managedDatabaseConnectionPoolTestPath                       = managedDatabasePostgresqlTestItemPath + "/connectionPool"
	managedDatabaseConnectionPoolTestItemPath                   = managedDatabaseConnectionPoolTestPath + "/p1"
	managedDatabaseConnectionPoolTestNotReady                   = `{"class": "Client::Conflict::ConnectionPoolStatusNotReady", "message": "connection pool \"pool\" status is not ready"}`
)

func managedDatabaseConnectionPoolTestAnswer(size int, userId string) string {
	return fmt.Sprintf(`{"id": "p1", "name": "pool", "databaseId": "d1", "mode": "TRANSACTION", "size": %d, "userId": %s, "port": 20185, "sslMode": "REQUIRE",
		"uri": "postgresql://postgresql-6b1b9c7e.database.cloud.ovh.net:20185/pool?sslmode=require", "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "updatedAt": "2026-11-12T10:00:00Z"}`, size, userId)
}

func newManagedDatabaseConnectionPoolTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*cloudManagedDatabasePostgresqlConnectionPoolResource, func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabasePostgresqlConnectionPoolResource().(*cloudManagedDatabasePostgresqlConnectionPoolResource)
	r.config = config
	return r, bodies
}

// managedDatabaseConnectionPoolTestModel returns the state of a pool of 10 connections, with user u1.
func managedDatabaseConnectionPoolTestModel() CloudManagedDatabaseConnectionPoolModel {
	var model CloudManagedDatabaseConnectionPoolModel
	model.ServiceName, model.Engine, model.ServiceId = types.StringValue("pj"), types.StringValue("postgresql"), types.StringValue(managedDatabasePostgresqlTestId)
	var pool managedDatabaseConnectionPool
	if err := json.Unmarshal([]byte(managedDatabaseConnectionPoolTestAnswer(10, `"u1"`)), &pool); err != nil {
		panic(err)
	}
	model.MergeWith(&pool)
	return model
}

func TestCloudManagedDatabasePostgresqlConnectionPoolCreate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseConnectionPoolTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseConnectionPoolTestPath: {{http.StatusOK, managedDatabaseConnectionPoolTestAnswer(10, `"u1"`)}},
	})

	plan := managedDatabaseConnectionPoolTestModel()
	plan.Id, plan.Port, plan.SslMode, plan.Uri = types.StringUnknown(), types.Int64Unknown(), types.StringUnknown(), types.StringUnknown()
	plan.Checksum, plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	assertManagedDatabaseJSON(t, json.RawMessage(bodies("POST " + managedDatabaseConnectionPoolTestPath)[0]),
		`{"databaseId": "d1", "mode": "TRANSACTION", "name": "pool", "size": 10, "userId": "u1"}`)
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
}

// The forward path can refuse a pool made right after its database, which Aiven doesn't know yet.
func TestCloudManagedDatabasePostgresqlConnectionPoolCreateUnknownDatabase(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseConnectionPoolTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseConnectionPoolTestPath: {
			{http.StatusForbidden, `{"class": "Client::Forbidden::BadRequest", "message": "Bad Request", "details": {"message": "Service database 'tfacc' does not exist"}}`},
			{http.StatusOK, managedDatabaseConnectionPoolTestAnswer(10, "null")},
		},
	})

	plan := managedDatabaseConnectionPoolTestModel()
	plan.UserId = types.StringNull()
	plan.Id, plan.Port, plan.SslMode, plan.Uri = types.StringUnknown(), types.Int64Unknown(), types.StringUnknown(), types.StringUnknown()
	plan.Checksum, plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	posts := bodies("POST " + managedDatabaseConnectionPoolTestPath)
	if len(posts) != 2 {
		t.Fatalf("expected the POST to be retried once after the 403, got %d POST", len(posts))
	}
	// A pool without user sends no userId.
	assertManagedDatabaseJSON(t, json.RawMessage(posts[1]), `{"databaseId": "d1", "mode": "TRANSACTION", "name": "pool", "size": 10}`)
}

// On apiv2's native path, a pool gets its URI and its port once PgBouncer serves it.
func TestCloudManagedDatabasePostgresqlConnectionPoolCreateWaitsForURI(t *testing.T) {
	ctx := context.Background()
	pending := strings.Replace(managedDatabaseConnectionPoolTestAnswer(10, `"u1"`), `"postgresql://postgresql-6b1b9c7e.database.cloud.ovh.net:20185/pool?sslmode=require"`, `""`, 1)
	ready := strings.ReplaceAll(managedDatabaseConnectionPoolTestAnswer(10, `"u1"`), "20185", "5433")
	r, _ := newManagedDatabaseConnectionPoolTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseConnectionPoolTestPath:    {{http.StatusOK, pending}},
		"GET " + managedDatabaseConnectionPoolTestItemPath: {{http.StatusOK, pending}, {http.StatusOK, ready}},
	})

	plan := managedDatabaseConnectionPoolTestModel()
	plan.Id, plan.Port, plan.SslMode, plan.Uri = types.StringUnknown(), types.Int64Unknown(), types.StringUnknown(), types.StringUnknown()
	plan.Checksum, plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var got CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Uri.ValueString() == "" || got.Port.ValueInt64() != 5433 {
		t.Errorf("expected the URI and the port of the pool once served, got %s and %s", got.Uri, got.Port)
	}
}

func TestCloudManagedDatabasePostgresqlConnectionPoolUpdate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseConnectionPoolTestResource(t, map[string][]managedDatabaseTestAnswer{
		// On apiv2's native path, a pool refuses a change while it applies the previous one.
		"PUT " + managedDatabaseConnectionPoolTestItemPath: {{http.StatusConflict, managedDatabaseConnectionPoolTestNotReady}, {http.StatusOK, managedDatabaseConnectionPoolTestAnswer(20, "null")}},
	})

	state := managedDatabaseConnectionPoolTestModel()
	plan := state
	plan.Size = types.Int64Value(20)
	plan.UserId = types.StringNull()
	resp := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Update(ctx, fwresource.UpdateRequest{
		Plan:  managedDatabaseObjectTestPlan(t, r, &plan),
		State: managedDatabaseObjectTestState(t, r, &state),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	puts := bodies("PUT " + managedDatabaseConnectionPoolTestItemPath)
	if len(puts) != 2 {
		t.Fatalf("expected the PUT to be retried once after the 409, got %d PUT", len(puts))
	}
	// Every argument, and the removed user as null.
	assertManagedDatabaseJSON(t, json.RawMessage(puts[1]), `{"databaseId": "d1", "mode": "TRANSACTION", "size": 20, "userId": null}`)
	var got CloudManagedDatabaseConnectionPoolModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Size.ValueInt64() != 20 || !got.UserId.IsNull() {
		t.Errorf("unexpected state: size %s, user_id %s", got.Size, got.UserId)
	}
}

// On apiv2's native path, a pool refuses a DELETE while it applies a change, and stays until its workflow
// removes it: deleting its database or its user before would be refused.
func TestCloudManagedDatabasePostgresqlConnectionPoolDelete(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseConnectionPoolTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseConnectionPoolTestItemPath: {{http.StatusConflict, managedDatabaseConnectionPoolTestNotReady}, {http.StatusOK, managedDatabaseConnectionPoolTestAnswer(10, `"u1"`)}},
		"GET " + managedDatabaseConnectionPoolTestItemPath:    {{http.StatusOK, managedDatabaseConnectionPoolTestAnswer(10, `"u1"`)}, {http.StatusNotFound, `{"message": "connection pool not found"}`}},
	})

	state := managedDatabaseConnectionPoolTestModel()
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got := len(bodies("DELETE " + managedDatabaseConnectionPoolTestItemPath)); got != 2 {
		t.Errorf("expected the DELETE to be retried once after the 409, got %d DELETE", got)
	}
	if got := len(bodies("GET " + managedDatabaseConnectionPoolTestItemPath)); got != 2 {
		t.Errorf("expected the delete to wait for the 404, got %d GET", got)
	}
}

func testAccCloudManagedDatabasePostgresqlConnectionPoolConfig(name, pool string) string {
	return testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), "") + `
resource "ovh_cloud_managed_database_database" "database" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  service_id   = ovh_cloud_managed_database.db.id
  name         = "tfacc"
}

resource "ovh_cloud_managed_database_postgresql_user" "user" {
  service_name = ovh_cloud_managed_database.db.service_name
  service_id   = ovh_cloud_managed_database.db.id
  name         = "tfacc"
}

resource "ovh_cloud_managed_database_postgresql_connection_pool" "pool" {
  service_name = ovh_cloud_managed_database.db.service_name
  service_id   = ovh_cloud_managed_database.db.id
  name         = "tfacc"
  database_id  = ovh_cloud_managed_database_database.database.id
` + pool + `
}
`
}

func TestAccCloudManagedDatabasePostgresqlConnectionPool_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlConnectionPoolConfig(name, `
  mode    = "TRANSACTION"
  size    = 10
  user_id = ovh_cloud_managed_database_postgresql_user.user.id`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlDatabaseResource, "name", "tfacc"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlDatabaseResource, "default", "false"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "mode", "TRANSACTION"),
					resource.TestCheckResourceAttrPair(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "user_id", testAccCloudManagedDatabasePostgresqlUserResource, "id"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "port"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "uri"),
				),
			},
			{
				Config: testAccCloudManagedDatabasePostgresqlConnectionPoolConfig(name, `
  mode = "SESSION"
  size = 20`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "mode", "SESSION"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "size", "20"),
					resource.TestCheckNoResourceAttr(testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "user_id"),
				),
			},
			{
				ResourceName:      testAccCloudManagedDatabasePostgresqlConnectionPoolResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudManagedImportStateIdFunc(testAccCloudManagedDatabasePostgresqlConnectionPoolResource),
			},
			{
				ResourceName:      testAccCloudManagedDatabasePostgresqlDatabaseResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudManagedImportStateIdFunc(testAccCloudManagedDatabasePostgresqlDatabaseResource),
			},
		},
	})
}
