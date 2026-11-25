package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"go.uber.org/ratelimit"
)

const displayNameTestServiceName = "ns1234567.ip-1-2-3.us"

func newDisplayNameTestServer(t *testing.T, iamDisplayName string) (*httptest.Server, *int) {
	t.Helper()
	putServicesCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case req.URL.Path == "/auth/time":
			fmt.Fprint(w, time.Now().Unix())
		case req.Method == http.MethodPut && req.URL.Path == "/dedicated/server/"+displayNameTestServiceName:
			fmt.Fprint(w, `null`)
		case req.Method == http.MethodGet && req.URL.Path == "/dedicated/server/"+displayNameTestServiceName+"/serviceInfos":
			fmt.Fprint(w, `{"serviceId":42}`)
		case req.Method == http.MethodPut && req.URL.Path == "/services/42":
			putServicesCalls++
			fmt.Fprint(w, `null`)
		case req.Method == http.MethodGet && req.URL.Path == "/dedicated/server/"+displayNameTestServiceName:
			fmt.Fprintf(w, `{"name":%q,"iam":{"displayName":%q}}`, displayNameTestServiceName, iamDisplayName)
		default:
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server, &putServicesCalls
}

func newDisplayNameTestConfig(t *testing.T, serverURL string) *Config {
	t.Helper()
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}
	client, err := ovh.NewClient(serverURL, "key", "secret", "consumer")
	if err != nil {
		t.Fatal(err)
	}
	return &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
}

func shortenDisplayNamePropagationTimeout(t *testing.T) {
	t.Helper()
	previous := displayNamePropagationTimeout
	displayNamePropagationTimeout = time.Second
	t.Cleanup(func() { displayNamePropagationTimeout = previous })
}

func TestDedicatedServerUpdateWarnsWhenDisplayNameIsNotPropagated(t *testing.T) {
	shortenDisplayNamePropagationTimeout(t)

	for name, tt := range map[string]struct {
		iamDisplayName string
		wantWarnings   int
	}{
		"propagated":     {iamDisplayName: "new-name", wantWarnings: 0},
		"not propagated": {iamDisplayName: "old-name", wantWarnings: 1},
	} {
		t.Run(name, func(t *testing.T) {
			server, putServicesCalls := newDisplayNameTestServer(t, tt.iamDisplayName)
			r := &dedicatedServerResource{config: newDisplayNameTestConfig(t, server.URL)}

			plan := &DedicatedServerModel{
				ServiceName: ovhtypes.NewTfStringValue(displayNameTestServiceName),
				DisplayName: ovhtypes.NewTfStringValue("new-name"),
			}
			var diags diag.Diagnostics
			r.updateDedicatedServerResource(context.Background(), nil, plan, &DedicatedServerModel{}, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if got := diags.WarningsCount(); got != tt.wantWarnings {
				t.Errorf("warnings = %d, want %d: %v", got, tt.wantWarnings, diags)
			}
			if *putServicesCalls != 1 {
				t.Errorf("PUT /services calls = %d, want 1", *putServicesCalls)
			}
		})
	}
}

func TestDedicatedServerUpdateSDKWarnsWhenDisplayNameIsNotPropagated(t *testing.T) {
	shortenDisplayNamePropagationTimeout(t)
	server, _ := newDisplayNameTestServer(t, "old-name")
	config := newDisplayNameTestConfig(t, server.URL)

	d := schema.TestResourceDataRaw(t, resourceDedicatedServerUpdate().Schema, map[string]interface{}{
		"service_name": displayNameTestServiceName,
		"display_name": "new-name",
	})

	diags := resourceDedicatedServerUpdateCreateOrUpdate(context.Background(), d, config)

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if len(diags) == 0 {
		t.Errorf("want a warning about display name propagation, got none")
	}
}
