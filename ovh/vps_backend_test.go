package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

const vpsTestId = "vps-1a2b3c4d.vps.ovh.net"

// vpsTestAPI is a fake API v1 serving one VPS, plus GET /v2/location. Each
// route answers the JSON value stored under its path, or 404 when there is
// none.
type vpsTestAPI struct {
	mu        sync.Mutex
	responses map[string]any
	requests  []string
}

func newVPSTestAPI() *vpsTestAPI {
	base := "/1.0/vps/" + vpsTestId
	return &vpsTestAPI{responses: map[string]any{
		base: map[string]any{
			"name":        vpsTestId,
			"displayName": "web-01",
			"state":       "running",
			"netbootMode": "local",
			"vcore":       4,
			"memoryLimit": 8192,
			"model":       map[string]any{"name": "vps-2025-model3", "disk": 160, "memory": 8192, "vcore": 4},
		},
		base + "/images/current": map[string]any{"id": "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d", "name": "Debian 12"},
		base + "/ips":            []string{"51.75.1.2", "2001:41d0::1"},
		base + "/datacenter":     map[string]any{"name": "GRA", "longName": "GRA1", "country": "fr"},
		"/v2/location":           vpsTestLocations,
	}}
}

func (a *vpsTestAPI) set(path string, response any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if response == nil {
		delete(a.responses, path)
		return
	}
	a.responses[path] = response
}

func (a *vpsTestAPI) count(request string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := 0
	for _, r := range a.requests {
		if r == request {
			n++
		}
	}
	return n
}

// backend starts the fake API and returns a v1 backend calling it.
func (a *vpsTestAPI) backend(t *testing.T) vpsBackend {
	t.Helper()
	// The client also reads these, and refuses two ways to authenticate.
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1.0/auth/time" {
			fmt.Fprint(w, time.Now().Unix())
			return
		}
		a.mu.Lock()
		a.requests = append(a.requests, r.Method+" "+r.URL.Path)
		response, found := a.responses[r.URL.Path]
		a.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet || !found {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"class": "Client::NotFound", "message": "The requested object does not exist"}`)
			return
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encoding the response of %s: %s", r.URL.Path, err)
		}
	}))
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return newVPSBackend(&Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())})
}

func vpsTestGet(t *testing.T, backend vpsBackend) (*VPSAPIServer, diag.Diagnostics) {
	t.Helper()
	srv, diags := backend.Get(context.Background(), vpsTestId)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if srv == nil {
		t.Fatal("got a nil VPS without error")
	}
	return srv, diags
}

func TestVPSV1BackendGet(t *testing.T) {
	api := newVPSTestAPI()
	srv, diags := vpsTestGet(t, api.backend(t))

	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
	if srv.Id != vpsTestId {
		t.Errorf("Id = %q, want %q", srv.Id, vpsTestId)
	}
	if srv.TargetSpec.DisplayName != "web-01" {
		t.Errorf("TargetSpec.DisplayName = %q, want %q", srv.TargetSpec.DisplayName, "web-01")
	}
	if srv.TargetSpec.PowerState == nil || *srv.TargetSpec.PowerState != "ACTIVE" {
		t.Errorf("TargetSpec.PowerState = %v, want ACTIVE", srv.TargetSpec.PowerState)
	}
	if from := srv.TargetSpec.Installation.From; from.ImageName == nil || *from.ImageName != "Debian 12" {
		t.Errorf("TargetSpec.Installation.From.ImageName = %v, want Debian 12", from.ImageName)
	}

	current := srv.CurrentState
	if current == nil {
		t.Fatal("CurrentState is nil")
	}
	if current.DisplayName != "web-01" {
		t.Errorf("CurrentState.DisplayName = %q, want %q", current.DisplayName, "web-01")
	}
	if current.PowerState != "ACTIVE" {
		t.Errorf("CurrentState.PowerState = %q, want ACTIVE", current.PowerState)
	}
	if current.Location.Region == nil || *current.Location.Region != "eu-west-gra" {
		t.Errorf("CurrentState.Location.Region = %v, want eu-west-gra", current.Location.Region)
	}
	if want := (VPSAPIHardware{Vcpus: 4, Ram: 8192, Disk: 160}); current.Model != want {
		t.Errorf("CurrentState.Model = %+v, want %+v", current.Model, want)
	}
	wantIPs := []VPSAPIIPAddress{{Ip: "51.75.1.2", Version: 4}, {Ip: "2001:41d0::1", Version: 6}}
	if len(current.IpAddresses) != len(wantIPs) {
		t.Fatalf("CurrentState.IpAddresses = %+v, want %+v", current.IpAddresses, wantIPs)
	}
	for i, want := range wantIPs {
		if current.IpAddresses[i] != want {
			t.Errorf("CurrentState.IpAddresses[%d] = %+v, want %+v", i, current.IpAddresses[i], want)
		}
	}
	from := current.Installation.From
	if from.ImageName == nil || *from.ImageName != "Debian 12" {
		t.Errorf("CurrentState.Installation.From.ImageName = %v, want Debian 12", from.ImageName)
	}
	if from.ImageId == nil || *from.ImageId != "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d" {
		t.Errorf("CurrentState.Installation.From.ImageId = %v, want 3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d", from.ImageId)
	}
}

func TestVPSV1BackendGetUnknownDatacenter(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/vps/"+vpsTestId+"/datacenter", map[string]any{"name": "XYZ"})
	srv, diags := vpsTestGet(t, api.backend(t))

	if srv.CurrentState.Location.Region != nil {
		t.Errorf("CurrentState.Location.Region = %q, want nil", *srv.CurrentState.Location.Region)
	}
	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1: %v", len(warnings), diags)
	}
	if text := warnings[0].Summary() + " " + warnings[0].Detail(); !strings.Contains(text, "XYZ") {
		t.Errorf("warning %q does not cite the datacenter code XYZ", text)
	}
}

func TestVPSV1BackendGetTransitional(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/vps/"+vpsTestId, map[string]any{
		"name":        vpsTestId,
		"displayName": "web-01",
		"state":       "installing",
		"vcore":       4,
		"memoryLimit": 8192,
		"model":       map[string]any{"disk": 160},
	})
	srv, _ := vpsTestGet(t, api.backend(t))

	if srv.TargetSpec.PowerState != nil {
		t.Errorf("TargetSpec.PowerState = %q, want nil", *srv.TargetSpec.PowerState)
	}
	if srv.CurrentState.PowerState != "UNKNOWN" {
		t.Errorf("CurrentState.PowerState = %q, want UNKNOWN", srv.CurrentState.PowerState)
	}
}

func TestVPSV1BackendGetLocationsOnce(t *testing.T) {
	api := newVPSTestAPI()
	backend := api.backend(t)
	vpsTestGet(t, backend)
	vpsTestGet(t, backend)

	if n := api.count("GET /v2/location"); n != 1 {
		t.Errorf("GET /v2/location called %d times, want 1", n)
	}
}

func TestVPSV1BackendGetNotFound(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/vps/"+vpsTestId, nil)
	srv, diags := api.backend(t).Get(context.Background(), vpsTestId)

	if !diags.HasError() {
		t.Fatalf("expected an error, got %+v", srv)
	}
	if srv != nil {
		t.Errorf("got %+v along with the error, want nil", srv)
	}
	want := "VPS " + vpsTestId + " not found"
	if text := diags.Errors()[0].Summary() + " " + diags.Errors()[0].Detail(); !strings.Contains(text, want) {
		t.Errorf("error %q does not contain %q", text, want)
	}
}

func TestVPSV1BackendGetLocationsError(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/v2/location", nil)
	backend := api.backend(t)

	srv, diags := backend.Get(context.Background(), vpsTestId)
	if !diags.HasError() {
		t.Fatalf("expected an error when GET /v2/location fails, got %+v", srv)
	}

	// The failure is not kept: the next read calls GET /v2/location again.
	api.set("/v2/location", vpsTestLocations)
	srv, _ = vpsTestGet(t, backend)
	if srv.CurrentState.Location.Region == nil || *srv.CurrentState.Location.Region != "eu-west-gra" {
		t.Errorf("CurrentState.Location.Region = %v, want eu-west-gra", srv.CurrentState.Location.Region)
	}
}

func TestVPSV1BackendGetUnknownImage(t *testing.T) {
	// API v1 cannot trace the image back when the first image of the chain no
	// longer exists (private image deleted after a restoration).
	api := newVPSTestAPI()
	api.set("/1.0/vps/"+vpsTestId+"/images/current", nil)
	srv, diags := vpsTestGet(t, api.backend(t))

	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
	for name, from := range map[string]VPSAPIInstallationFrom{
		"TargetSpec":   srv.TargetSpec.Installation.From,
		"CurrentState": srv.CurrentState.Installation.From,
	} {
		if from.ImageId != nil {
			t.Errorf("%s.Installation.From.ImageId = %q, want nil", name, *from.ImageId)
		}
		if from.ImageName != nil {
			t.Errorf("%s.Installation.From.ImageName = %q, want nil", name, *from.ImageName)
		}
	}
}

func TestVPSV1BackendGetEmptyImage(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/vps/"+vpsTestId+"/images/current", map[string]any{"id": "", "name": ""})
	srv, _ := vpsTestGet(t, api.backend(t))

	from := srv.CurrentState.Installation.From
	if from.ImageId != nil || from.ImageName != nil {
		t.Errorf("CurrentState.Installation.From = {%v, %v}, want nil fields", from.ImageId, from.ImageName)
	}
}

func TestVPSV1BackendGetImageError(t *testing.T) {
	api := newVPSTestAPI()
	backend := api.backend(t)
	// Any error other than a 404 stays an error.
	api.set("/1.0/vps/"+vpsTestId+"/images/current", "not an image")

	if srv, diags := backend.Get(context.Background(), vpsTestId); !diags.HasError() {
		t.Fatalf("expected an error, got %+v", srv)
	}
}
