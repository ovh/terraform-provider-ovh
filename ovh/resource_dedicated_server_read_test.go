package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

func TestDedicatedServerReadRemovesDeletedServer(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}

	for name, tt := range map[string]struct {
		status    int
		wantKept  bool
		wantError bool
	}{
		"not found":      {status: http.StatusNotFound},
		"internal error": {status: http.StatusInternalServerError, wantKept: true, wantError: true},
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path == "/auth/time" {
					fmt.Fprint(w, time.Now().Unix())
					return
				}
				if req.Method != http.MethodGet || req.URL.Path != "/dedicated/server/ns1234567.ip-1-2-3.eu" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
				}
				w.WriteHeader(tt.status)
				fmt.Fprint(w, `{"message":"The requested object (serviceName = ns1234567.ip-1-2-3.eu) does not exist"}`)
			}))
			defer server.Close()

			client, err := ovh.NewClient(server.URL, "key", "secret", "consumer")
			if err != nil {
				t.Fatal(err)
			}
			r := &dedicatedServerResource{config: &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}}

			var schemaResp resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
			objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
			values := map[string]tftypes.Value{}
			for attribute, attributeType := range objectType.AttributeTypes {
				values[attribute] = tftypes.NewValue(attributeType, nil)
			}
			values["service_name"] = tftypes.NewValue(tftypes.String, "ns1234567.ip-1-2-3.eu")
			values["id"] = tftypes.NewValue(tftypes.String, "ns1234567.ip-1-2-3.eu")
			prior := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, values)}

			resp := resource.ReadResponse{State: prior}
			r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

			if got := resp.Diagnostics.HasError(); got != tt.wantError {
				t.Errorf("HasError = %v, want %v: %v", got, tt.wantError, resp.Diagnostics)
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
