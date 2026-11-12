package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

func TestCloudStorageObjectBucketReplicationAndLifecycleRead(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}

	for kind, newResource := range map[string]func(*Config) resource.Resource{
		"replication": func(c *Config) resource.Resource { return &cloudStorageObjectBucketReplicationResource{config: c} },
		"lifecycle":   func(c *Config) resource.Resource { return &cloudStorageObjectBucketLifecycleResource{config: c} },
	} {
		for name, tt := range map[string]struct {
			status      int
			body        string
			wantKept    bool
			wantWarning bool
			wantError   bool
		}{
			"region in maintenance": {status: http.StatusConflict, body: `{"class":"RegionInMaintenance","message":"region GRA is currently in maintenance"}`, wantKept: true, wantWarning: true},
			"object storage down":   {status: http.StatusBadGateway, body: `{"message":"object storage is unavailable"}`, wantKept: true, wantWarning: true},
			"internal error":        {status: http.StatusInternalServerError, body: `{"message":"internal error"}`, wantKept: true, wantWarning: true},
			"not found":             {status: http.StatusNotFound, body: `{"message":"not found"}`},
			"no rules":              {status: http.StatusOK, body: `{"rules":[]}`},
			"bad request":           {status: http.StatusBadRequest, body: `{"message":"bad request"}`, wantKept: true, wantError: true},
		} {
			t.Run(kind+"/"+name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if req.URL.Path == "/auth/time" {
						fmt.Fprint(w, time.Now().Unix())
						return
					}
					if req.Method != http.MethodGet || !strings.HasSuffix(req.URL.Path, "/storage/object/bucket/GRA_bucket/"+kind) {
						t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
					}
					w.WriteHeader(tt.status)
					fmt.Fprint(w, tt.body)
				}))
				defer server.Close()
				client, err := ovh.NewClient(server.URL, "key", "secret", "consumer")
				if err != nil {
					t.Fatal(err)
				}
				r := newResource(&Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())})

				var schemaResp resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
				objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
				values := map[string]tftypes.Value{}
				for attribute, attributeType := range objectType.AttributeTypes {
					values[attribute] = tftypes.NewValue(attributeType, nil)
				}
				values["service_name"] = tftypes.NewValue(tftypes.String, "project")
				values["bucket_id"] = tftypes.NewValue(tftypes.String, "GRA_bucket")
				values["id"] = tftypes.NewValue(tftypes.String, "project/GRA_bucket")
				prior := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, values)}

				resp := resource.ReadResponse{State: prior}
				r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

				if got := resp.Diagnostics.HasError(); got != tt.wantError {
					t.Errorf("HasError = %v, want %v: %v", got, tt.wantError, resp.Diagnostics)
				}
				warnings := resp.Diagnostics.Warnings()
				if got := len(warnings) > 0; got != tt.wantWarning {
					t.Errorf("warning = %v, want %v: %v", got, tt.wantWarning, resp.Diagnostics)
				}
				if tt.wantWarning {
					wantDetail := fmt.Sprintf("answered %d: ", tt.status)
					if detail := warnings[0].Detail(); !strings.Contains(detail, wantDetail) || !strings.Contains(detail, "/"+kind) {
						t.Errorf("warning detail %q lacks the endpoint or %q", detail, wantDetail)
					}
					if tt.status == http.StatusConflict && !strings.Contains(warnings[0].Detail(), "region GRA is currently in maintenance") {
						t.Errorf("warning detail %q lacks the API message", warnings[0].Detail())
					}
				}
				if kept := resp.State.Raw.Equal(prior.Raw); kept != tt.wantKept {
					t.Errorf("state kept = %v, want %v", kept, tt.wantKept)
				}
				if !tt.wantKept && !resp.State.Raw.IsNull() {
					t.Errorf("state not removed: %v", resp.State.Raw)
				}
			})
		}
	}
}
