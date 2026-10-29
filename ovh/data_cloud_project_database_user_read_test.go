package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

func TestCloudProjectDatabaseUserRead(t *testing.T) {
	tests := []struct {
		engine     string
		resource   func() *schema.Resource
		dataSource func() *schema.Resource
	}{
		{"mysql", resourceCloudProjectDatabaseUser, dataSourceCloudProjectDatabaseUser},
		{"postgresql", resourceCloudProjectDatabasePostgresqlUser, dataSourceCloudProjectDatabasePostgresqlUser},
		{"mongodb", resourceCloudProjectDatabaseMongodbUser, dataSourceCloudProjectDatabaseMongodbUser},
		{"clickhouse", resourceCloudProjectDatabaseClickhouseUser, dataSourceCloudProjectDatabaseClickhouseUser},
		{"opensearch", resourceCloudProjectDatabaseOpensearchUser, dataSourceCloudProjectDatabaseOpensearchUser},
		{"valkey", resourceCloudProjectDatabaseValkeyUser, dataSourceCloudProjectDatabaseValkeyUser},
	}
	for _, tt := range tests {
		t.Run(tt.engine, func(t *testing.T) {
			for _, mode := range []string{"adopt built-in user", "data source by name", "data source by id"} {
				t.Run(mode, func(t *testing.T) {
					name := "avnadmin"
					if tt.engine == "valkey" {
						name = "default"
					}
					values := map[string]interface{}{
						"service_name": "project", "cluster_id": "cluster", "engine": tt.engine,
					}
					r := tt.dataSource()
					if mode == "adopt built-in user" {
						r = tt.resource()
					}
					if mode == "data source by id" {
						values["id"] = "user-id"
					} else {
						values["name"] = name
					}
					d := schema.TestResourceDataRaw(t, r.Schema, values)
					endpoint := fmt.Sprintf("/cloud/project/project/database/%s/cluster/user", tt.engine)
					var requests []string
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						if req.URL.Path == "/auth/time" {
							fmt.Fprint(w, time.Now().Unix())
							return
						}
						requests = append(requests, req.Method+" "+req.URL.Path)
						switch {
						case req.Method == http.MethodGet && req.URL.Path == endpoint:
							fmt.Fprint(w, `["user-id"]`)
						case req.Method == http.MethodGet && req.URL.Path == endpoint+"/user-id":
							fmt.Fprintf(w, `{"id":"user-id","username":%q,"status":"READY"}`, name)
						default:
							t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
							http.NotFound(w, req)
						}
					}))
					defer server.Close()
					client, err := ovh.NewClient(server.URL, "key", "secret", "consumer")
					if err != nil {
						t.Fatal(err)
					}
					config := &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
					read := tt.dataSource().ReadContext
					var diags diag.Diagnostics
					if mode == "adopt built-in user" {
						// Adoption passes resource data without an id attribute to the data source
						// reader, then updates the existing user using the discovered state ID.
						if d.Get("id") != nil {
							t.Fatal("resource must have no id attribute")
						}
						updated := false
						diags = postCloudProjectDatabaseUser(context.Background(), d, config, tt.engine, read, nil,
							func(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
								updated = true
								if data.Id() != "user-id" {
									t.Errorf("update received ID %q, want user-id", data.Id())
								}
								return nil
							}, nil)
						if !updated {
							t.Error("built-in user was not updated after adoption")
						}
					} else {
						diags = read(context.Background(), d, config)
					}
					if diags.HasError() {
						t.Fatalf("read failed: %v", diags)
					}
					if d.Id() != "user-id" || d.Get("name") != name || d.Get("status") != "READY" {
						t.Errorf("unexpected user state: id=%q name=%v status=%v", d.Id(), d.Get("name"), d.Get("status"))
					}
					wantRequests := []string{"GET " + endpoint + "/user-id"}
					if mode != "data source by id" {
						wantRequests = append([]string{"GET " + endpoint}, wantRequests...)
					}
					if !reflect.DeepEqual(requests, wantRequests) {
						t.Errorf("requests = %v, want %v", requests, wantRequests)
					}
				})
			}
		})
	}
}
