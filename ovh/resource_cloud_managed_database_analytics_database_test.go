package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

const (
	testAccCloudManagedDatabasePostgresqlDatabaseResource = "ovh_cloud_managed_database_database.database"
	managedDatabaseDatabaseTestPath                       = managedDatabasePostgresqlTestItemPath + "/database"
	managedDatabaseDatabaseTestAnswer                     = `{"id": "d1", "name": "app", "default": false, "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "updatedAt": "2026-11-12T10:00:00Z"}`
)

func newManagedDatabaseDatabaseTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*cloudManagedDatabaseDatabaseResource, func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabaseDatabaseResource().(*cloudManagedDatabaseDatabaseResource)
	r.config = config
	return r, bodies
}

func managedDatabaseDatabaseTestModel() CloudManagedDatabaseDatabaseModel {
	return CloudManagedDatabaseDatabaseModel{
		ServiceName: types.StringValue("pj"),
		Engine:      types.StringValue("postgresql"),
		ServiceId:   types.StringValue(managedDatabasePostgresqlTestId),
		Name:        types.StringValue("app"),
		Id:          types.StringValue("d1"),
		Default:     types.BoolValue(false),
		Checksum:    types.StringValue("c1"),
		CreatedAt:   types.StringValue("2026-11-12T10:00:00Z"),
		UpdatedAt:   types.StringValue("2026-11-12T10:00:00Z"),
	}
}

func TestCloudManagedDatabasePostgresqlDatabaseCreate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseDatabaseTestPath: {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, managedDatabaseDatabaseTestAnswer}},
	})

	plan := managedDatabaseDatabaseTestModel()
	plan.Id, plan.Default, plan.Checksum = types.StringUnknown(), types.BoolUnknown(), types.StringUnknown()
	plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	posts := bodies("POST " + managedDatabaseDatabaseTestPath)
	if len(posts) != 2 {
		t.Fatalf("expected the POST to be retried once after the 409, got %d POST", len(posts))
	}
	assertManagedDatabaseJSON(t, json.RawMessage(posts[1]), `{"name": "app"}`)
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
}

// The forward path can answer an error to a POST that made the database: it is then taken from the list.
func TestCloudManagedDatabasePostgresqlDatabaseCreateFailedButMade(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseDatabaseTestPath: {{http.StatusNotFound, `{"class": "Client::NotFound::BadRequest", "message": "Bad Request"}`}},
		"GET " + managedDatabaseDatabaseTestPath:  {{http.StatusOK, `[` + managedDatabaseDatabaseTestAnswer + `]`}},
	})

	plan := managedDatabaseDatabaseTestModel()
	plan.Id, plan.Default, plan.Checksum = types.StringUnknown(), types.BoolUnknown(), types.StringUnknown()
	plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var got CloudManagedDatabaseDatabaseModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != "d1" {
		t.Errorf("expected the database d1 of the list, got %s", got.Id)
	}
}

// The database can be on any page of the list. The fake API of newManagedDatabaseTestConfig can't send
// the header of the next page.
func TestCloudManagedDatabasePostgresqlDatabaseCreateFailedButMadeOnLaterPage(t *testing.T) {
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().Unix())
	})
	mux.HandleFunc(managedDatabaseDatabaseTestPath, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"class": "Client::NotFound::BadRequest", "message": "Bad Request"}`)
		case r.Header.Get("X-Pagination-Cursor") == "":
			w.Header().Set("X-Pagination-Cursor-Next", "2")
			fmt.Fprint(w, `[]`)
		default:
			fmt.Fprint(w, `[`+managedDatabaseDatabaseTestAnswer+`]`)
		}
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	r := NewCloudManagedDatabaseDatabaseResource().(*cloudManagedDatabaseDatabaseResource)
	r.config = &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}

	plan := managedDatabaseDatabaseTestModel()
	plan.Id, plan.Default, plan.Checksum = types.StringUnknown(), types.BoolUnknown(), types.StringUnknown()
	plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var got CloudManagedDatabaseDatabaseModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Id.ValueString() != "d1" {
		t.Errorf("expected the database d1 of the second page, got %s", got.Id)
	}
}

// A database that was there before the POST isn't taken over: the error says how to import it.
func TestCloudManagedDatabasePostgresqlDatabaseCreateAlreadyExists(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseDatabaseTestPath: {{http.StatusConflict, `{"class": "Client::Conflict::DatabaseAlready", "message": "database app already exists"}`}},
		"GET " + managedDatabaseDatabaseTestPath:  {{http.StatusOK, `[` + managedDatabaseDatabaseTestAnswer + `]`}},
	})

	plan := managedDatabaseDatabaseTestModel()
	plan.Id, plan.Default, plan.Checksum = types.StringUnknown(), types.BoolUnknown(), types.StringUnknown()
	plan.CreatedAt, plan.UpdatedAt = types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if detail, id := resp.Diagnostics.Errors()[0].Detail(), "pj/postgresql/"+managedDatabasePostgresqlTestId+"/d1"; !strings.Contains(detail, id) {
		t.Errorf("expected the error to give the import ID %s, got %s", id, detail)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected no state, got %v", resp.State.Raw)
	}
}

// On apiv2's native path, a database stays until its workflow drops it, and a pool using it until its own
// workflow removes it: the DELETE of a pool's database right after the pool's would be refused.
func TestCloudManagedDatabasePostgresqlDatabaseDelete(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseDatabaseTestPath + "/d1": {{http.StatusOK, managedDatabaseDatabaseTestAnswer}},
		"GET " + managedDatabaseDatabaseTestPath + "/d1":    {{http.StatusOK, managedDatabaseDatabaseTestAnswer}, {http.StatusNotFound, `{"message": "database not found"}`}},
	})

	state := managedDatabaseDatabaseTestModel()
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got := len(bodies("GET " + managedDatabaseDatabaseTestPath + "/d1")); got != 2 {
		t.Errorf("expected the delete to wait for the 404, got %d GET", got)
	}
}

func TestCloudManagedDatabasePostgresqlDatabaseDeleteDefault(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{})
	state := managedDatabaseDatabaseTestModel()
	state.Default = types.BoolValue(true)
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() || resp.Diagnostics.WarningsCount() != 1 {
		t.Errorf("expected only a warning, got %v", resp.Diagnostics)
	}
	if requests := bodies(""); len(requests) != 0 {
		t.Errorf("expected no request for the default database, got %v", requests)
	}
}

func TestCloudManagedDatabasePostgresqlDatabaseReadNotFound(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseDatabaseTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabaseDatabaseTestPath + "/d1": {{http.StatusNotFound, `{"message": "database not found"}`}},
	})

	state := managedDatabaseDatabaseTestModel()
	resp := fwresource.ReadResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Read(ctx, fwresource.ReadRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected the database to be removed from the state, got %v", resp.State.Raw)
	}
}
