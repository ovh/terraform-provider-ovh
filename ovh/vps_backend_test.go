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

// vpsTestAPI is a fake API v1 serving VPS, plus GET /v2/location. Each route
// answers the JSON value stored under its path, a vpsTestStatus error, or 404
// when there is none.
type vpsTestAPI struct {
	mu        sync.Mutex
	responses map[string]any
	requests  []string
	// delays holds the time each path waits before answering.
	delays map[string]time.Duration
	// inFlight and maxInFlight count the GET /vps/{serviceName} being served.
	inFlight, maxInFlight int
}

// vpsTestStatus is an error answer of the fake API.
type vpsTestStatus struct {
	code    int
	message string
}

// vpsTestExpired is the answer of API v1 on the routes of an expired VPS.
var vpsTestExpired = vpsTestStatus{code: 460, message: "This service is expired"}

func newVPSTestAPI() *vpsTestAPI {
	api := &vpsTestAPI{
		responses: map[string]any{
			"/1.0/vps":     []string{vpsTestId},
			"/v2/location": vpsTestLocations,
		},
		delays: map[string]time.Duration{},
	}
	api.addVPS(vpsTestId, "web-01")
	return api
}

// addVPS serves a running VPS without option under id.
func (a *vpsTestAPI) addVPS(id, displayName string) {
	base := "/1.0/vps/" + id
	a.mu.Lock()
	defer a.mu.Unlock()
	a.responses[base] = map[string]any{
		"name":        id,
		"displayName": displayName,
		"state":       "running",
		"netbootMode": "local",
		"vcore":       4,
		"memoryLimit": 8192,
		"model":       map[string]any{"name": "vps-2025-model3", "disk": 160, "memory": 8192, "vcore": 4},
	}
	a.responses[base+"/images/current"] = map[string]any{"id": "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d", "name": "Debian 12"}
	a.responses[base+"/ips"] = []string{"51.75.1.2", "2001:41d0::1"}
	a.responses[base+"/datacenter"] = map[string]any{"name": "GRA", "longName": "GRA1", "country": "fr"}
	a.responses[base+"/option"] = []string{}
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
	return newVPSBackend(a.config(t))
}

// config starts the fake API and returns a provider configuration calling it.
func (a *vpsTestAPI) config(t *testing.T) *Config {
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
		vpsRoute := strings.HasPrefix(r.URL.Path, "/1.0/vps/") && strings.Count(r.URL.Path, "/") == 3
		a.mu.Lock()
		a.requests = append(a.requests, r.Method+" "+r.URL.Path)
		response, found := a.responses[r.URL.Path]
		delay := a.delays[r.URL.Path]
		if vpsRoute {
			a.inFlight++
			a.maxInFlight = max(a.maxInFlight, a.inFlight)
		}
		a.mu.Unlock()

		time.Sleep(delay)
		if vpsRoute {
			a.mu.Lock()
			a.inFlight--
			a.mu.Unlock()
		}

		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet || !found {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"class": "Client::NotFound", "message": "The requested object does not exist"}`)
			return
		}
		if status, ok := response.(vpsTestStatus); ok {
			w.WriteHeader(status.code)
			fmt.Fprintf(w, `{"message": %q}`, status.message)
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
	return &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
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

// vpsTestBase is the API v1 path of the test VPS.
const vpsTestBase = "/1.0/vps/" + vpsTestId

func TestVPSV1BackendGetOptions(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/option", []string{"automatedBackup", "snapshot"})
	api.set(vpsTestBase+"/automatedBackup", map[string]any{
		"state":               "enabled",
		"schedule":            "03:00:00",
		"rotation":            7,
		"serviceResourceName": vpsTestId + "-autobackup",
	})
	srv, diags := vpsTestGet(t, api.backend(t))

	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
	options := srv.CurrentState.Options
	autoBackup := options.AutoBackup
	if !autoBackup.Enabled || autoBackup.Cron == nil || *autoBackup.Cron != "0 3 * * *" || autoBackup.Rotation == nil || *autoBackup.Rotation != 7 {
		t.Errorf("CurrentState.Options.AutoBackup = {%t, %v, %v}, want {true, 0 3 * * *, 7}", autoBackup.Enabled, autoBackup.Cron, autoBackup.Rotation)
	}
	if !options.Snapshot.Enabled {
		t.Error("CurrentState.Options.Snapshot.Enabled = false, want true")
	}
	if options.AdditionalDisks.Enabled || options.AdditionalDisks.Size != nil {
		t.Errorf("CurrentState.Options.AdditionalDisks = {%t, %v}, want {false, nil}", options.AdditionalDisks.Enabled, options.AdditionalDisks.Size)
	}
	if target := srv.TargetSpec.Options.AutoBackup; target == nil || target.Cron != "0 3 * * *" {
		t.Errorf("TargetSpec.Options.AutoBackup = %+v, want {Cron: 0 3 * * *}", target)
	}
	if n := api.count("GET " + vpsTestBase + "/disks"); n != 0 {
		t.Errorf("GET /disks called %d times without the additionalDisk option", n)
	}
}

func TestVPSV1BackendGetNoOption(t *testing.T) {
	api := newVPSTestAPI()
	srv, _ := vpsTestGet(t, api.backend(t))

	options := srv.CurrentState.Options
	if options.AutoBackup.Enabled || options.AutoBackup.Cron != nil || options.AutoBackup.Rotation != nil {
		t.Errorf("CurrentState.Options.AutoBackup = %+v, want disabled with null fields", options.AutoBackup)
	}
	if options.Snapshot.Enabled {
		t.Error("CurrentState.Options.Snapshot.Enabled = true, want false")
	}
	if options.AdditionalDisks.Enabled || options.AdditionalDisks.Size != nil {
		t.Errorf("CurrentState.Options.AdditionalDisks = %+v, want disabled with null size", options.AdditionalDisks)
	}
	if srv.TargetSpec.Options.AutoBackup != nil {
		t.Errorf("TargetSpec.Options.AutoBackup = %+v, want nil", srv.TargetSpec.Options.AutoBackup)
	}
	for _, path := range []string{"/automatedBackup", "/disks"} {
		if n := api.count("GET " + vpsTestBase + path); n != 0 {
			t.Errorf("GET %s called %d times without the option", path, n)
		}
	}
}

func TestVPSV1BackendGetAutoBackupDisabled(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/option", []string{"automatedBackup"})
	api.set(vpsTestBase+"/automatedBackup", map[string]any{"state": "disabled", "schedule": "03:00:00", "rotation": 7})
	srv, _ := vpsTestGet(t, api.backend(t))

	autoBackup := srv.CurrentState.Options.AutoBackup
	if !autoBackup.Enabled || autoBackup.Cron != nil || autoBackup.Rotation == nil || *autoBackup.Rotation != 7 {
		t.Errorf("CurrentState.Options.AutoBackup = {%t, %v, %v}, want {true, nil, 7}", autoBackup.Enabled, autoBackup.Cron, autoBackup.Rotation)
	}
	if srv.TargetSpec.Options.AutoBackup != nil {
		t.Errorf("TargetSpec.Options.AutoBackup = %+v, want nil", srv.TargetSpec.Options.AutoBackup)
	}
}

func TestVPSV1BackendGetAutoBackupNoSchedule(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/option", []string{"automatedBackup"})
	api.set(vpsTestBase+"/automatedBackup", map[string]any{"state": "enabled", "schedule": nil, "rotation": 7})
	srv, diags := vpsTestGet(t, api.backend(t))

	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
	if cron := srv.CurrentState.Options.AutoBackup.Cron; cron != nil {
		t.Errorf("CurrentState.Options.AutoBackup.Cron = %q, want nil", *cron)
	}
	if srv.TargetSpec.Options.AutoBackup != nil {
		t.Errorf("TargetSpec.Options.AutoBackup = %+v, want nil", srv.TargetSpec.Options.AutoBackup)
	}
}

func TestVPSV1BackendGetAdditionalDisk(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/option", []string{"additionalDisk"})
	api.set(vpsTestBase+"/disks", []int64{746694019, 746694020})
	api.set(vpsTestBase+"/disks/746694019", map[string]any{"id": 746694019, "type": "primary", "size": 20, "state": "connected"})
	api.set(vpsTestBase+"/disks/746694020", map[string]any{"id": 746694020, "type": "additional", "size": 50, "state": "connected"})
	srv, _ := vpsTestGet(t, api.backend(t))

	disks := srv.CurrentState.Options.AdditionalDisks
	if !disks.Enabled || disks.Size == nil || *disks.Size != 50 {
		t.Errorf("CurrentState.Options.AdditionalDisks = {%t, %v}, want {true, 50}", disks.Enabled, disks.Size)
	}
}

func TestVPSV1BackendGetExpired(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/images/current", vpsTestExpired)
	srv, diags := api.backend(t).Get(context.Background(), vpsTestId)

	if !diags.HasError() {
		t.Fatalf("expected an error, got %+v", srv)
	}
	want := "VPS " + vpsTestId + " is expired"
	if text := diags.Errors()[0].Summary() + " " + diags.Errors()[0].Detail(); !strings.Contains(text, want) {
		t.Errorf("error %q does not contain %q", text, want)
	}
}

// vpsTestListAPI serves the VPS a, b and c, in this order.
func vpsTestListAPI() *vpsTestAPI {
	api := newVPSTestAPI()
	ids := []string{"vps-aaaaaaaa.vps.ovh.net", "vps-bbbbbbbb.vps.ovh.net", "vps-cccccccc.vps.ovh.net"}
	for _, id := range ids {
		api.addVPS(id, strings.TrimSuffix(id, ".vps.ovh.net"))
	}
	api.set("/1.0/vps", ids)
	return api
}

func vpsTestIds(servers []VPSAPIServer) []string {
	ids := make([]string, 0, len(servers))
	for _, srv := range servers {
		ids = append(ids, srv.Id)
	}
	return ids
}

func TestVPSV1BackendList(t *testing.T) {
	api := vpsTestListAPI()
	// b is deleted between GET /vps and its own GET.
	api.set("/1.0/vps/vps-bbbbbbbb.vps.ovh.net", nil)
	servers, diags := api.backend(t).List(context.Background())

	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
	want := "vps-aaaaaaaa.vps.ovh.net vps-cccccccc.vps.ovh.net"
	if got := strings.Join(vpsTestIds(servers), " "); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if len(servers) == 2 && servers[0].TargetSpec.DisplayName != "vps-aaaaaaaa" {
		t.Errorf("servers[0].TargetSpec.DisplayName = %q, want vps-aaaaaaaa", servers[0].TargetSpec.DisplayName)
	}
}

func TestVPSV1BackendListExpired(t *testing.T) {
	api := vpsTestListAPI()
	api.set("/1.0/vps/vps-bbbbbbbb.vps.ovh.net/option", vpsTestExpired)
	servers, diags := api.backend(t).List(context.Background())

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	want := "vps-aaaaaaaa.vps.ovh.net vps-cccccccc.vps.ovh.net"
	if got := strings.Join(vpsTestIds(servers), " "); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1: %v", len(warnings), diags)
	}
	if text := warnings[0].Summary() + " " + warnings[0].Detail(); !strings.Contains(text, "vps-bbbbbbbb.vps.ovh.net") {
		t.Errorf("warning %q does not cite the expired VPS", text)
	}
}

func TestVPSV1BackendListError(t *testing.T) {
	api := vpsTestListAPI()
	api.set("/1.0/vps/vps-bbbbbbbb.vps.ovh.net", vpsTestStatus{code: http.StatusInternalServerError, message: "Failed to fetch sub object"})
	servers, diags := api.backend(t).List(context.Background())

	if !diags.HasError() {
		t.Fatalf("expected an error, got %v", vpsTestIds(servers))
	}
	if servers != nil {
		t.Errorf("got %v along with the error, want nil", vpsTestIds(servers))
	}
}

func TestVPSV1BackendListKeepsOrder(t *testing.T) {
	api := vpsTestListAPI()
	// Answers arrive in the reverse order of GET /vps.
	api.delays["/1.0/vps/vps-aaaaaaaa.vps.ovh.net"] = 200 * time.Millisecond
	api.delays["/1.0/vps/vps-bbbbbbbb.vps.ovh.net"] = 100 * time.Millisecond
	servers, diags := api.backend(t).List(context.Background())

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	want := "vps-aaaaaaaa.vps.ovh.net vps-bbbbbbbb.vps.ovh.net vps-cccccccc.vps.ovh.net"
	if got := strings.Join(vpsTestIds(servers), " "); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if api.maxInFlight < 2 {
		t.Errorf("at most %d VPS read at once, want them read in parallel", api.maxInFlight)
	}
}

func TestVPSV1BackendListParallelLimit(t *testing.T) {
	api := newVPSTestAPI()
	ids := make([]string, 0, 20)
	for i := range 20 {
		id := fmt.Sprintf("vps-%08d.vps.ovh.net", i)
		api.addVPS(id, id)
		api.delays["/1.0/vps/"+id] = 50 * time.Millisecond
		ids = append(ids, id)
	}
	api.set("/1.0/vps", ids)
	servers, diags := api.backend(t).List(context.Background())

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if got := strings.Join(vpsTestIds(servers), " "); got != strings.Join(ids, " ") {
		t.Errorf("got %q, want %q", got, strings.Join(ids, " "))
	}
	if api.maxInFlight > vpsV1ListParallelism {
		t.Errorf("%d VPS read at once, want at most %d", api.maxInFlight, vpsV1ListParallelism)
	}
	if api.maxInFlight < 2 {
		t.Errorf("at most %d VPS read at once, want them read in parallel", api.maxInFlight)
	}
}

func TestVPSBackendLocationsSharedByProvider(t *testing.T) {
	api := newVPSTestAPI()
	config := api.config(t)
	// Each data source read builds its own backend.
	vpsTestGet(t, newVPSBackend(config))
	vpsTestGet(t, newVPSBackend(config))

	if n := api.count("GET /v2/location"); n != 1 {
		t.Errorf("GET /v2/location called %d times for one provider, want 1", n)
	}

	// Another provider configuration has its own cache.
	vpsTestGet(t, newVPSBackend(api.config(t)))
	if n := api.count("GET /v2/location"); n != 2 {
		t.Errorf("GET /v2/location called %d times for two providers, want 2", n)
	}
}
