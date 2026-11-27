package ovh

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// Images of the fake API: images/current answers Debian 12.
const (
	vpsTestDebianId = "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d"
	vpsTestUbuntuId = "11111111-2222-4333-8444-555555555555"
	vpsTestSshKey   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINew new"
)

// vpsTestUpdateAPI is newVPSTestAPI with fast polls, the action routes, the
// images Debian 12 and Ubuntu 24.04 available, and the automated backup
// subscribed and scheduled at 3:00.
func vpsTestUpdateAPI(t *testing.T) (*vpsTestAPI, *vpsTestActions) {
	t.Helper()
	vpsTestFastPolls(t)
	api := newVPSTestAPI()
	api.set(vpsTestBase+"/images/available", []string{vpsTestDebianId, vpsTestUbuntuId})
	api.set(vpsTestBase+"/images/available/"+vpsTestDebianId, map[string]any{"id": vpsTestDebianId, "name": "Debian 12"})
	api.set(vpsTestBase+"/images/available/"+vpsTestUbuntuId, map[string]any{"id": vpsTestUbuntuId, "name": "Ubuntu 24.04"})
	api.set(vpsTestBase+"/option", []string{"automatedBackup"})
	api.set(vpsTestBase+"/automatedBackup", map[string]any{"state": "enabled", "schedule": "03:00:00", "rotation": 7})
	return api, api.serveActions()
}

// vpsTestInstallation is the installation of the VPS of the fake API (Debian
// 12, password sent by email) with sshKey and no script.
func vpsTestInstallation(sshKey attr.Value) attr.Value {
	return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
		"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
			"image_id":   tfStr(vpsTestDebianId),
			"image_name": tfStr("Debian 12"),
		}),
		"ssh_key":                  sshKey,
		"post_installation_script": vpsTfStringPtr(nil),
		"send_password_by_email":   types.BoolValue(true),
	})
}

// vpsTestUpdateState is the state of the VPS of vpsTestUpdateAPI as an apply
// leaves it, with overrides.
func vpsTestUpdateState(t *testing.T, overrides map[string]attr.Value) tfsdk.State {
	t.Helper()
	values := map[string]attr.Value{
		"display_name": tfStr("web-01"),
		"installation": vpsTestInstallation(tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old")),
	}
	for name, value := range overrides {
		values[name] = value
	}
	return vpsTestFullState(t, values)
}

// vpsTestNewPrivate sets *field, a resource private state of the framework
// (whose type is internal to it), to a new one holding keys, each set to
// []byte("true").
func vpsTestNewPrivate(t *testing.T, field any, keys ...string) {
	t.Helper()
	v := reflect.ValueOf(field).Elem()
	v.Set(reflect.New(v.Type().Elem()))
	private := v.Interface().(interface {
		SetKey(context.Context, string, []byte) diag.Diagnostics
	})
	for _, key := range keys {
		if diags := private.SetKey(context.Background(), key, []byte("true")); diags.HasError() {
			t.Fatalf("setting the private state: %v", diags)
		}
	}
}

// vpsTestPrivateGetter is a resource private state of the framework.
type vpsTestPrivateGetter interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
}

// vpsTestHasKey tells whether a private state holds key.
func vpsTestHasKey(t *testing.T, private vpsTestPrivateGetter, key string) bool {
	t.Helper()
	value, diags := private.GetKey(context.Background(), key)
	if diags.HasError() {
		t.Fatalf("reading the private state: %v", diags)
	}
	return value != nil
}

// vpsTestPending tells whether a private state holds installation_pending.
func vpsTestPending(t *testing.T, private vpsTestPrivateGetter) bool {
	t.Helper()
	return vpsTestHasKey(t, private, vpsPrivateInstallationPending)
}

// vpsTestPrivateKeys returns the keys of the resource a private state holds.
func vpsTestPrivateKeys(t *testing.T, private vpsTestPrivateGetter) []string {
	t.Helper()
	var keys []string
	for _, key := range []string{vpsPrivateInstallationPending, vpsPrivateInstallationKnown} {
		if vpsTestHasKey(t, private, key) {
			keys = append(keys, key)
		}
	}
	return keys
}

// vpsTestActionBody returns the body of the only request to the action path
// (such as "POST /rebuild") among requests, decoded.
func vpsTestActionBody(t *testing.T, requests []string, action string) map[string]any {
	t.Helper()
	method, suffix, _ := strings.Cut(action, " ")
	prefix := method + " " + vpsTestBase + suffix + " "
	var bodies []string
	for _, request := range requests {
		if body, ok := strings.CutPrefix(request, prefix); ok {
			bodies = append(bodies, body)
		}
	}
	if len(bodies) != 1 {
		t.Fatalf("%s: got %d requests with a body in %q, want 1", action, len(bodies), requests)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(bodies[0]), &body); err != nil {
		t.Fatalf("%s: decoding %s: %s", action, bodies[0], err)
	}
	return body
}

// vpsTestActionNames returns the requests as "METHOD /suffix", without the path
// of the VPS nor the body.
func vpsTestActionNames(requests []string) []string {
	names := make([]string, 0, len(requests))
	for _, request := range requests {
		method, rest, _ := strings.Cut(request, " ")
		route, _, _ := strings.Cut(rest, " ")
		names = append(names, method+" "+strings.TrimPrefix(route, vpsTestBase))
	}
	return names
}

func vpsTestCheckBody(t *testing.T, action string, got map[string]any, want string) {
	t.Helper()
	var wantBody map[string]any
	if err := json.Unmarshal([]byte(want), &wantBody); err != nil {
		t.Fatalf("decoding %s: %s", want, err)
	}
	if !reflect.DeepEqual(got, wantBody) {
		gotJSON, _ := json.Marshal(got)
		t.Errorf("%s: body %s, want %s", action, gotJSON, want)
	}
}

func TestVPSResourceUpdateDisplayName(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	state := vpsTestUpdateState(t, nil)
	plan := vpsTestUpdateState(t, map[string]attr.Value{"display_name": tfStr("web-02")})

	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	requests := actions.list()
	if got, want := vpsTestActionNames(requests), []string{"PUT "}; !slices.Equal(got, want) {
		t.Errorf("got requests %q, want %q", got, want)
	}
	vpsTestCheckBody(t, "PUT", vpsTestActionBody(t, requests, "PUT "), `{"displayName":"web-02"}`)
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"display_name":               "web-02",
		"current_state.display_name": "web-02",
	})
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("the state holds unknown values: %v", resp.State.Raw)
	}
}

func TestVPSResourceUpdateOrder(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	state := vpsTestUpdateState(t, nil)
	plan := vpsTestUpdateState(t, map[string]attr.Value{
		"display_name": tfStr("web-02"),
		"installation": vpsTestInstallation(tfStr(vpsTestSshKey)),
		"options": types.ObjectValueMust(vpsOptionsAttrTypes(), map[string]attr.Value{
			"auto_backup": types.ObjectValueMust(vpsAutoBackupAttrTypes(), map[string]attr.Value{"cron": tfStr("0 2 * * *")}),
		}),
		"power_state": tfStr("SHUTOFF"),
	})

	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	requests := actions.list()
	want := []string{"PUT ", "POST /rebuild", "POST /automatedBackup/reschedule", "POST /stop"}
	if got := vpsTestActionNames(requests); !slices.Equal(got, want) {
		t.Errorf("got requests %q, want %q", got, want)
	}
	vpsTestCheckBody(t, "rebuild", vpsTestActionBody(t, requests, "POST /rebuild"),
		`{"imageId":"`+vpsTestDebianId+`","publicSshKey":"`+vpsTestSshKey+`","doNotSendPassword":false}`)
	vpsTestCheckBody(t, "reschedule", vpsTestActionBody(t, requests, "POST /automatedBackup/reschedule"), `{"schedule":"02:00:00"}`)
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"display_name":              "web-02",
		"installation.ssh_key":      vpsTestSshKey,
		"options.auto_backup.cron":  "0 2 * * *",
		"power_state":               "SHUTOFF",
		"current_state.power_state": "SHUTOFF",
	})
}

func TestVPSResourceUpdateRebuildFailsKeepsPriorInstallation(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	api.handle("POST "+vpsTestBase+"/rebuild", func(body []byte) any {
		actions.add("POST "+vpsTestBase+"/rebuild", body)
		return api.serveTask(7, "reinstallVm", "error")
	})
	state := vpsTestUpdateState(t, nil)
	plan := vpsTestUpdateState(t, map[string]attr.Value{
		"display_name": tfStr("web-02"),
		"installation": vpsTestInstallation(tfStr(vpsTestSshKey)),
		"options": types.ObjectValueMust(vpsOptionsAttrTypes(), map[string]attr.Value{
			"auto_backup": types.ObjectValueMust(vpsAutoBackupAttrTypes(), map[string]attr.Value{"cron": tfStr("0 2 * * *")}),
		}),
	})

	resp := vpsTestUpdate(t, api, state, plan)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error: the reinstallation failed")
	}
	if got := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(got, "reinstallVm") {
		t.Errorf("detail %q does not name the task reinstallVm", got)
	}
	if got := vpsTestActionNames(actions.list()); !slices.Equal(got, []string{"PUT ", "POST /rebuild"}) {
		t.Errorf("got requests %q, want the rename and the reinstallation only", got)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"display_name":             "web-02",
		"installation.ssh_key":     "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old",
		"options.auto_backup.cron": "0 3 * * *",
	})
}

func TestVPSResourceUpdateSendPasswordOnlyNoRebuild(t *testing.T) {
	api, _ := vpsTestUpdateAPI(t)
	installation := func(send bool) attr.Value {
		return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(vpsTestDebianId),
				"image_name": tfStr("Debian 12"),
			}),
			"ssh_key":                  tfStr(vpsTestSshKey),
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(send),
		})
	}
	state := vpsTestUpdateState(t, map[string]attr.Value{"installation": installation(true)})
	plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": installation(false)})

	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{"installation.send_password_by_email": "false"})
	if got := vpsTestRequests(api); got != 0 {
		t.Errorf("%d API calls, want 0", got)
	}
}

func TestVPSResourceUpdateInstallationPending(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	state := vpsTestUpdateState(t, nil)

	resp := vpsTestUpdatePending(t, api, state, state, vpsPrivateInstallationPending)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	if got := vpsTestActionNames(actions.list()); !slices.Equal(got, []string{"POST /rebuild"}) {
		t.Errorf("got requests %q, want the reinstallation only", got)
	}
	if vpsTestPending(t, resp.Private) {
		t.Error("installation_pending is still set after the reinstallation")
	}
}

// vpsTestModifyPlanUpdate runs ModifyPlan from the state to the plan of
// config: config, with the values the plan modifiers would set from the
// state.
func vpsTestModifyPlanUpdate(t *testing.T, api *vpsTestAPI, state tfsdk.State, config, plan tfsdk.Plan, keys ...string) *resource.ModifyPlanResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	req := resource.ModifyPlanRequest{
		Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw},
		Plan:   plan,
		State:  state,
	}
	vpsTestNewPrivate(t, &req.Private, keys...)
	resp := &resource.ModifyPlanResponse{Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw.Copy()}, Private: req.Private}
	r.ModifyPlan(context.Background(), req, resp)
	return resp
}

// vpsTestUpdateConfig is the configuration of the VPS of vpsTestUpdateState,
// with overrides (dotted paths).
func vpsTestUpdateConfig(t *testing.T, overrides map[string]attr.Value) tfsdk.Plan {
	t.Helper()
	values := map[string]attr.Value{
		"plan_code":                    tfStr("vps-2025-model1"),
		"location.region":              tfStr("eu-west-sbg"),
		"display_name":                 tfStr("web-01"),
		"installation.from.image_name": tfStr("Debian 12"),
		"installation.ssh_key":         tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"),
	}
	for name, value := range overrides {
		values[name] = value
	}
	return vpsTestPlan(t, values)
}

func TestVPSResourceModifyPlanUpdate(t *testing.T) {
	reinstallWarning := "This apply reinstalls VPS " + vpsTestId + ", which erases its disk"
	pendingWarning := "The installation requested when the VPS was ordered is not applied yet: this apply reinstalls the VPS, which erases its disk"
	warnings := func(diags diag.Diagnostics) []string {
		var out []string
		for _, d := range diags.Warnings() {
			out = append(out, d.Summary()+": "+d.Detail())
		}
		return out
	}
	hasWarning := func(diags diag.Diagnostics, want string) bool {
		return slices.ContainsFunc(warnings(diags), func(w string) bool { return strings.Contains(w, want) })
	}

	t.Run("image name not available", func(t *testing.T) {
		api, actions := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, nil)
		config := vpsTestUpdateConfig(t, map[string]attr.Value{"installation.from.image_name": tfStr("Windows")})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(vpsTestDebianId),
				"image_name": tfStr("Windows"),
			}),
			"ssh_key":                  tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"),
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(true),
		})})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan))
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error: Windows is not available")
		}
		d, ok := resp.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
		if !ok || d.Path().String() != "installation.from.image_name" {
			t.Errorf("error %v is not on installation.from.image_name", resp.Diagnostics.Errors()[0])
		}
		if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, "Debian 12, Ubuntu 24.04") {
			t.Errorf("detail %q does not list the available images", detail)
		}
		if got := actions.list(); len(got) != 0 {
			t.Errorf("got action requests %q, want none", got)
		}
	})

	t.Run("image id not available", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, nil)
		const unknownImage = "99999999-2222-4333-8444-555555555555"
		config := vpsTestUpdateConfig(t, map[string]attr.Value{
			"installation.from.image_name": vpsTfStringPtr(nil),
			"installation.from.image_id":   tfStr(unknownImage),
		})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(unknownImage),
				"image_name": tfStr("Debian 12"),
			}),
			"ssh_key":                  tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"),
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(true),
		})})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan))
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error: the image id is not available")
		}
		d, ok := resp.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
		if !ok || d.Path().String() != "installation.from.image_id" {
			t.Errorf("error %v is not on installation.from.image_id", resp.Diagnostics.Errors()[0])
		}
		if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, "Debian 12, Ubuntu 24.04") {
			t.Errorf("detail %q does not list the available images", detail)
		}
	})

	t.Run("reinstallation", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, nil)
		config := vpsTestUpdateConfig(t, map[string]attr.Value{"installation.ssh_key": tfStr(vpsTestSshKey)})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": vpsTestInstallation(tfStr(vpsTestSshKey))})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan))
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected plan error: %v", resp.Diagnostics)
		}
		if !hasWarning(resp.Diagnostics, reinstallWarning) {
			t.Errorf("warnings %q, want %q", warnings(resp.Diagnostics), reinstallWarning)
		}
	})

	// The image the configuration does not set is null in the plan: the
	// reinstallation is from the configured one.
	t.Run("reinstallation from the image name", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, map[string]attr.Value{"installation": types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(vpsTestDebianId),
				"image_name": tfStr("Debian 12"),
			}),
			"ssh_key":                  tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"),
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(true),
		})})
		config := vpsTestUpdateConfig(t, map[string]attr.Value{"installation.from.image_name": tfStr("Ubuntu 24.04")})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(vpsTestDebianId),
				"image_name": tfStr("Ubuntu 24.04"),
			}),
			"ssh_key":                  tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"),
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(true),
		})})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan))
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected plan error: %v", resp.Diagnostics)
		}
		vpsTestCheckValues(t, resp.Plan.Raw, map[string]string{
			"installation.from.image_name": "Ubuntu 24.04",
			"installation.from.image_id":   "<null>",
		})
	})

	t.Run("installation pending", func(t *testing.T) {
		api, actions := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, nil)
		config := vpsTestUpdateConfig(t, nil)
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(state), vpsPrivateInstallationPending)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected plan error: %v", resp.Diagnostics)
		}
		if !hasWarning(resp.Diagnostics, pendingWarning) {
			t.Errorf("warnings %q, want %q", warnings(resp.Diagnostics), pendingWarning)
		}
		if hasWarning(resp.Diagnostics, reinstallWarning) {
			t.Errorf("warnings %q: the reinstallation is told twice", warnings(resp.Diagnostics))
		}
		if got := vpsTestValue(t, resp.Plan.Raw, "current_state"); got != "<unknown>" {
			t.Errorf("current_state = %s, want <unknown>", got)
		}
		if !vpsTestPending(t, resp.Private) {
			t.Error("installation_pending is no longer set")
		}
		if got := actions.list(); len(got) != 0 {
			t.Errorf("got action requests %q, want none", got)
		}
	})

	t.Run("display name only", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := vpsTestUpdateState(t, nil)
		config := vpsTestUpdateConfig(t, map[string]attr.Value{"display_name": tfStr("web-02")})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"display_name": tfStr("web-02")})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan))
		if len(resp.Diagnostics) != 0 {
			t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
		}
		if got := vpsTestRequests(api); got != 0 {
			t.Errorf("%d API calls, want 0", got)
		}
	})
}

func TestVPSResourceValidateConfigSendPassword(t *testing.T) {
	for name, tc := range map[string]struct {
		values map[string]attr.Value
		ok     bool
	}{
		"false without key": {map[string]attr.Value{"installation.send_password_by_email": types.BoolValue(false)}, false},
		"false with key": {map[string]attr.Value{
			"installation.send_password_by_email": types.BoolValue(false),
			"installation.ssh_key":                tfStr(vpsTestSshKey),
		}, true},
		"false with unknown key": {map[string]attr.Value{
			"installation.send_password_by_email": types.BoolValue(false),
			"installation.ssh_key":                ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		}, true},
		"true without key": {map[string]attr.Value{"installation.send_password_by_email": types.BoolValue(true)}, true},
		"not set":          {nil, true},
	} {
		t.Run(name, func(t *testing.T) {
			config := vpsTestCreateConfig(t, tc.values)
			resp := &resource.ValidateConfigResponse{}
			(&vpsResource{}).ValidateConfig(context.Background(), resource.ValidateConfigRequest{Config: config}, resp)
			if resp.Diagnostics.HasError() == tc.ok {
				t.Fatalf("errors %v, want ok = %t", resp.Diagnostics, tc.ok)
			}
			if tc.ok {
				return
			}
			d, isPath := resp.Diagnostics.Errors()[0].(diag.DiagnosticWithPath)
			if !isPath || d.Path().String() != "installation.send_password_by_email" {
				t.Errorf("error %v is not on installation.send_password_by_email", resp.Diagnostics.Errors()[0])
			}
			want := "Setting send_password_by_email to false needs installation.ssh_key: without an SSH key, the password is the only way to log in"
			if got := resp.Diagnostics.Errors()[0].Detail(); got != want {
				t.Errorf("detail %q, want %q", got, want)
			}
		})
	}
}

// A VPS whose order is not delivered yet has no id: nothing can be applied to
// it, and the plan checks nothing.
func TestVPSResourceUpdateNotDelivered(t *testing.T) {
	ordered := func(t *testing.T, overrides map[string]attr.Value) tfsdk.State {
		values := map[string]attr.Value{
			"id": vpsTfStringPtr(nil),
			"order": types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
				"duration":     tfStr("P1M"),
				"pricing_mode": tfStr("default"),
				"order_id":     types.Int64Value(vpsTestOrderId),
				"options":      types.ListNull(types.StringType),
			}),
		}
		for name, value := range overrides {
			values[name] = value
		}
		return vpsTestUpdateState(t, values)
	}

	t.Run("update", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := ordered(t, nil)
		plan := ordered(t, map[string]attr.Value{"display_name": tfStr("web-02")})
		resp := vpsTestUpdatePending(t, api, state, plan, vpsPrivateInstallationPending)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error: the VPS is not delivered yet")
		}
		if got := resp.Diagnostics.Errors()[0].Summary(); got != "The VPS is not delivered yet" {
			t.Errorf("summary %q", got)
		}
		if got := vpsTestRequests(api); got != 0 {
			t.Errorf("%d API calls, want 0", got)
		}
		if !resp.State.Raw.Equal(state.Raw) {
			t.Error("the state must stay unchanged")
		}
	})

	t.Run("plan", func(t *testing.T) {
		api, _ := vpsTestUpdateAPI(t)
		state := ordered(t, nil)
		config := vpsTestUpdateConfig(t, map[string]attr.Value{"installation.ssh_key": tfStr(vpsTestSshKey)})
		plan := ordered(t, map[string]attr.Value{"installation": vpsTestInstallation(tfStr(vpsTestSshKey))})
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), vpsPrivateInstallationPending)
		if len(resp.Diagnostics) != 0 {
			t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
		}
		if got := vpsTestRequests(api); got != 0 {
			t.Errorf("%d API calls, want 0", got)
		}
	})
}

// A reinstallation starts the VPS: the power state of the plan is applied
// after it, even when it does not change.
func TestVPSResourceUpdateRebuildKeepsPowerState(t *testing.T) {
	for _, tc := range []struct {
		name, powerState, state, netboot string
		// netbootAfter is the netboot mode the reinstallation leaves.
		netbootAfter string
		want         []string
	}{
		{"shutoff", "SHUTOFF", "stopped", "local", "local", []string{"POST /rebuild", "POST /stop"}},
		{"rescue, netboot kept", "RESCUE", "rescued", "rescue", "rescue", []string{"POST /rebuild", "POST /reboot"}},
		{"rescue, netboot reset", "RESCUE", "rescued", "rescue", "local", []string{"POST /rebuild", "PUT ", "POST /reboot"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api, actions := vpsTestUpdateAPI(t)
			api.change(vpsTestBase, func(vps map[string]any) {
				vps["state"], vps["netbootMode"] = tc.state, tc.netboot
			})
			api.handle("POST "+vpsTestBase+"/rebuild", func(body []byte) any {
				actions.add("POST "+vpsTestBase+"/rebuild", body)
				api.change(vpsTestBase, func(vps map[string]any) {
					vps["state"], vps["netbootMode"] = "running", tc.netbootAfter
				})
				return api.serveTask(7, "reinstallVm", "done")
			})
			state := vpsTestUpdateState(t, map[string]attr.Value{"power_state": tfStr(tc.powerState)})
			plan := vpsTestUpdateState(t, map[string]attr.Value{
				"power_state":  tfStr(tc.powerState),
				"installation": vpsTestInstallation(tfStr(vpsTestSshKey)),
			})

			resp := vpsTestUpdate(t, api, state, plan)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected update error: %v", resp.Diagnostics)
			}
			requests := actions.list()
			if got := vpsTestActionNames(requests); !slices.Equal(got, tc.want) {
				t.Errorf("got requests %q, want %q", got, tc.want)
			}
			if slices.Contains(tc.want, "PUT ") {
				vpsTestCheckBody(t, "PUT", vpsTestActionBody(t, requests, "PUT "), `{"netbootMode":"rescue"}`)
			}
			vpsTestCheckValues(t, resp.State.Raw, map[string]string{
				"power_state":               tc.powerState,
				"current_state.power_state": tc.powerState,
				"installation.ssh_key":      vpsTestSshKey,
			})
		})
	}
}

// order.options only applies to the order: a change is recorded in the state.
func TestVPSResourceUpdateOrderOptionsOnlyRecorded(t *testing.T) {
	api, _ := vpsTestUpdateAPI(t)
	order := func(options types.List) attr.Value {
		return types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     tfStr("P1M"),
			"pricing_mode": tfStr("default"),
			"order_id":     types.Int64Value(42),
			"options":      options,
		})
	}
	state := vpsTestUpdateState(t, map[string]attr.Value{"order": order(types.ListNull(types.StringType))})
	plan := vpsTestUpdateState(t, map[string]attr.Value{"order": order(types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("option-snapshot-2025-model1"),
	}))})

	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.Equal(plan.Raw) {
		t.Errorf("state = %v, want the plan %v", resp.State.Raw, plan.Raw)
	}
	if got := vpsTestRequests(api); got != 0 {
		t.Errorf("%d API calls, want 0", got)
	}
}

// vpsTestInstallationWith is vpsTestInstallation with the key and the script
// given, null when empty.
func vpsTestInstallationWith(sshKey, script string) attr.Value {
	value := func(s string) attr.Value {
		if s == "" {
			return vpsTfStringPtr(nil)
		}
		return tfStr(s)
	}
	return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
		"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
			"image_id":   tfStr(vpsTestDebianId),
			"image_name": tfStr("Debian 12"),
		}),
		"ssh_key":                  value(sshKey),
		"post_installation_script": value(script),
		"send_password_by_email":   types.BoolValue(true),
	})
}

// A VPS Terraform installed (installation_known) has no key or script when the
// state has none: setting one reinstalls the VPS. Without the marker, as after
// an import or the migration of the state, the state does not know: setting
// one is only recorded (TestVPSResourceUpdateSshKeyAfterImport).
func TestVPSResourceInstallationKnown(t *testing.T) {
	reinstallWarning := "This apply reinstalls VPS " + vpsTestId + ", which erases its disk"
	for _, attribute := range []struct{ name, sshKey, script string }{
		{"ssh_key", vpsTestSshKey, ""},
		{"post_installation_script", "", "#!/bin/sh\ntrue\n"},
	} {
		t.Run(attribute.name+" added", func(t *testing.T) {
			api, actions := vpsTestUpdateAPI(t)
			state := vpsTestUpdateState(t, map[string]attr.Value{"installation": vpsTestInstallationWith("", "")})
			plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": vpsTestInstallationWith(attribute.sshKey, attribute.script)})
			config := vpsTestUpdateConfig(t, map[string]attr.Value{
				"installation.ssh_key":                  vpsTestInstallationWith(attribute.sshKey, attribute.script).(types.Object).Attributes()["ssh_key"],
				"installation.post_installation_script": vpsTestInstallationWith(attribute.sshKey, attribute.script).(types.Object).Attributes()["post_installation_script"],
			})

			planResp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), vpsPrivateInstallationKnown)
			if planResp.Diagnostics.HasError() {
				t.Fatalf("unexpected plan error: %v", planResp.Diagnostics)
			}
			if !slices.ContainsFunc(planResp.Diagnostics.Warnings(), func(d diag.Diagnostic) bool { return strings.Contains(d.Detail(), reinstallWarning) }) {
				t.Errorf("warnings %v, want %q", planResp.Diagnostics.Warnings(), reinstallWarning)
			}

			resp := vpsTestUpdatePending(t, api, state, plan, vpsPrivateInstallationKnown)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected update error: %v", resp.Diagnostics)
			}
			if got := vpsTestActionNames(actions.list()); !slices.Equal(got, []string{"POST /rebuild"}) {
				t.Errorf("got requests %q, want the reinstallation", got)
			}
			if !vpsTestHasKey(t, resp.Private, vpsPrivateInstallationKnown) {
				t.Error("installation_known is no longer set after the reinstallation")
			}
		})
	}

	t.Run("key removed then added back", func(t *testing.T) {
		api, actions := vpsTestUpdateAPI(t)
		withKey := vpsTestUpdateState(t, map[string]attr.Value{"installation": vpsTestInstallationWith(vpsTestSshKey, "")})
		withoutKey := vpsTestUpdateState(t, map[string]attr.Value{"installation": vpsTestInstallationWith("", "")})

		removed := vpsTestUpdatePending(t, api, withKey, withoutKey, vpsPrivateInstallationKnown)
		if removed.Diagnostics.HasError() {
			t.Fatalf("unexpected update error: %v", removed.Diagnostics)
		}
		added := vpsTestUpdatePending(t, api, tfsdk.State{Schema: removed.State.Schema, Raw: removed.State.Raw},
			withKey, vpsTestPrivateKeys(t, removed.Private)...)
		if added.Diagnostics.HasError() {
			t.Fatalf("unexpected update error: %v", added.Diagnostics)
		}
		requests := actions.list()
		if got := vpsTestActionNames(requests); !slices.Equal(got, []string{"POST /rebuild", "POST /rebuild"}) {
			t.Errorf("got requests %q, want a reinstallation without the key, then one with it", got)
		}
		if len(requests) == 2 && (strings.Contains(requests[0], "publicSshKey") || !strings.Contains(requests[1], vpsTestSshKey)) {
			t.Errorf("got requests %q, want the key in the second reinstallation only", requests)
		}
	})

	// Created from image_name, the VPS has a null image_id in the state:
	// setting image_id instead reinstalls it.
	t.Run("image_id set instead of image_name", func(t *testing.T) {
		api, actions := vpsTestUpdateAPI(t)
		from := func(imageId attr.Value) attr.Value {
			return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
				"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
					"image_id":   imageId,
					"image_name": tfStr("Debian 12"),
				}),
				"ssh_key":                  vpsTfStringPtr(nil),
				"post_installation_script": vpsTfStringPtr(nil),
				"send_password_by_email":   types.BoolValue(true),
			})
		}
		state := vpsTestUpdateState(t, map[string]attr.Value{"installation": from(vpsTfStringPtr(nil))})
		plan := vpsTestUpdateState(t, map[string]attr.Value{"installation": from(tfStr(vpsTestUbuntuId))})
		resp := vpsTestUpdatePending(t, api, state, plan, vpsPrivateInstallationKnown)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update error: %v", resp.Diagnostics)
		}
		requests := actions.list()
		if got := vpsTestActionNames(requests); !slices.Equal(got, []string{"POST /rebuild"}) {
			t.Fatalf("got requests %q, want the reinstallation", got)
		}
		vpsTestCheckBody(t, "rebuild", vpsTestActionBody(t, requests, "POST /rebuild"), `{"imageId":"`+vpsTestUbuntuId+`"}`)
	})

	// The migration of the state of the former ovh_vps sets no marker: a key
	// the state has none of is only recorded.
	t.Run("key added after the migration", func(t *testing.T) {
		api, actions := vpsTestUpdateAPI(t)
		vpsTestAddService(api)
		read := vpsTestRead(t, api, vpsTestUpgrade(t, vpsTestV0ImportedState))
		if read.Diagnostics.HasError() {
			t.Fatalf("unexpected read error: %v", read.Diagnostics)
		}
		state := tfsdk.State{Schema: read.State.Schema, Raw: read.State.Raw}
		var plan tfsdk.State = tfsdk.State{Schema: state.Schema, Raw: state.Raw.Copy()}
		if diags := plan.SetAttribute(context.Background(), path.Root("installation").AtName("ssh_key"), tfStr(vpsTestSshKey)); diags.HasError() {
			t.Fatalf("setting the key: %v", diags)
		}
		resp := vpsTestUpdate(t, api, state, plan)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update error: %v", resp.Diagnostics)
		}
		if got := actions.list(); len(got) != 0 {
			t.Errorf("got action requests %q, want none", got)
		}
		vpsTestCheckValues(t, resp.State.Raw, map[string]string{"installation.ssh_key": vpsTestSshKey})
	})
}

// The acceptance test orders a VPS without key, then adds one: the key
// reinstalls the VPS.
func TestVPSResourceCreateThenAddSshKey(t *testing.T) {
	api, _, actions := vpsTestCreateSettingsAPI(t)
	created := vpsTestCreate(t, api, vpsTestCreatePlan(t, nil))
	if created.Diagnostics.HasError() {
		t.Fatalf("unexpected create error: %v", created.Diagnostics)
	}
	keys := vpsTestPrivateKeys(t, created.Private)
	if !slices.Contains(keys, vpsPrivateInstallationKnown) {
		t.Fatalf("private keys after the creation %q, want installation_known", keys)
	}

	state := tfsdk.State{Schema: created.State.Schema, Raw: created.State.Raw}
	plan := tfsdk.State{Schema: state.Schema, Raw: state.Raw.Copy()}
	if diags := plan.SetAttribute(context.Background(), path.Root("installation").AtName("ssh_key"), tfStr(vpsTestSshKey)); diags.HasError() {
		t.Fatalf("setting the key: %v", diags)
	}
	config := vpsTestCreateConfig(t, map[string]attr.Value{"installation.ssh_key": tfStr(vpsTestSshKey)})

	planResp := vpsTestModifyPlanUpdate(t, api, state, tfsdk.Plan(config), tfsdk.Plan(plan), keys...)
	if planResp.Diagnostics.HasError() {
		t.Fatalf("unexpected plan error: %v", planResp.Diagnostics)
	}
	if planResp.Diagnostics.WarningsCount() == 0 {
		t.Error("the plan does not warn of the reinstallation")
	}
	resp := vpsTestUpdatePending(t, api, state, tfsdk.State{Schema: planResp.Plan.Schema, Raw: planResp.Plan.Raw}, keys...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	if got := vpsTestActionNames(actions.list()); !slices.Equal(got, []string{"POST /rebuild"}) {
		t.Errorf("got requests %q, want the reinstallation", got)
	}
}

// vpsTestCheckPlansConsistent checks that the plan made at apply time keeps
// every value the plan made at plan time knows, as Terraform requires.
func vpsTestCheckPlansConsistent(t *testing.T, atPlan, atApply tftypes.Value) {
	t.Helper()
	diffs, err := atPlan.Diff(atApply)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range diffs {
		if d.Value1 != nil && d.Value1.IsFullyKnown() {
			t.Errorf("inconsistent final plan: %s planned %v, then %v", d.Path, d.Value1, d.Value2)
		}
	}
}

// A value of the installation unknown at plan time may reinstall the VPS: the
// plan warns, and plans unknown the image the configuration does not set, so
// that the plan made at apply time is consistent with it.
func TestVPSResourceModifyPlanUnknownInstallation(t *testing.T) {
	unknown := ovhtypes.TfStringValue{StringValue: types.StringUnknown()}
	mayWarning := "This apply may reinstall VPS " + vpsTestId + ", which erases its disk: installation depends on values known only at apply time"
	installation := func(imageName, sshKey, script attr.Value) attr.Value {
		return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr(vpsTestDebianId),
				"image_name": imageName,
			}),
			"ssh_key":                  sshKey,
			"post_installation_script": script,
			"send_password_by_email":   types.BoolValue(true),
		})
	}
	oldKey, oldScript := tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOld old"), tfStr("#!/bin/sh\nold\n")
	debian := tfStr("Debian 12")

	for _, tc := range []struct {
		name string
		// attribute is the dotted path of the value unknown at plan time.
		attribute       string
		known           attr.Value
		atPlan, atApply attr.Value
	}{
		{"ssh_key", "installation.ssh_key", tfStr(vpsTestSshKey),
			installation(debian, unknown, oldScript), installation(debian, tfStr(vpsTestSshKey), oldScript)},
		{"post_installation_script", "installation.post_installation_script", tfStr("#!/bin/sh\nnew\n"),
			installation(debian, oldKey, unknown), installation(debian, oldKey, tfStr("#!/bin/sh\nnew\n"))},
		{"image_name", "installation.from.image_name", tfStr("Ubuntu 24.04"),
			installation(unknown, oldKey, oldScript), installation(tfStr("Ubuntu 24.04"), oldKey, oldScript)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api, actions := vpsTestUpdateAPI(t)
			state := vpsTestUpdateState(t, map[string]attr.Value{"installation": installation(debian, oldKey, oldScript)})
			configAt := func(value attr.Value) tfsdk.Plan {
				return vpsTestUpdateConfig(t, map[string]attr.Value{
					"installation.ssh_key":                  oldKey,
					"installation.post_installation_script": oldScript,
					tc.attribute:                            value,
				})
			}

			atPlan := vpsTestModifyPlanUpdate(t, api, state, configAt(unknown), tfsdk.Plan(vpsTestUpdateState(t, map[string]attr.Value{"installation": tc.atPlan})))
			if atPlan.Diagnostics.HasError() {
				t.Fatalf("unexpected plan error: %v", atPlan.Diagnostics)
			}
			if !slices.ContainsFunc(atPlan.Diagnostics.Warnings(), func(d diag.Diagnostic) bool { return d.Detail() == mayWarning }) {
				t.Errorf("warnings %v, want %q", atPlan.Diagnostics.Warnings(), mayWarning)
			}
			if got := vpsTestValue(t, atPlan.Plan.Raw, "installation.from.image_id"); got != "<unknown>" {
				t.Errorf("installation.from.image_id planned %s, want <unknown>", got)
			}

			atApply := vpsTestModifyPlanUpdate(t, api, state, configAt(tc.known), tfsdk.Plan(vpsTestUpdateState(t, map[string]attr.Value{"installation": tc.atApply})))
			if atApply.Diagnostics.HasError() {
				t.Fatalf("unexpected apply plan error: %v", atApply.Diagnostics)
			}
			vpsTestCheckPlansConsistent(t, atPlan.Plan.Raw, atApply.Plan.Raw)
			if got := actions.list(); len(got) != 0 {
				t.Errorf("got action requests %q, want none", got)
			}
		})
	}
}

// A reinstallation API v1 accepted runs on, even when its wait ends: the state
// takes the new installation, and the next plan does not reinstall again.
func TestVPSResourceUpdateRebuildStillRunning(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	api.handle("POST "+vpsTestBase+"/rebuild", func(body []byte) any {
		actions.add("POST "+vpsTestBase+"/rebuild", body)
		return api.serveTask(7, "reinstallVm", "doing")
	})
	timeouts := types.ObjectValueMust(vpsTimeoutsAttrTypes(), map[string]attr.Value{
		"create": vpsTfStringPtr(nil),
		"delete": vpsTfStringPtr(nil),
		"update": tfStr("200ms"),
	})
	state := vpsTestUpdateState(t, map[string]attr.Value{"timeouts": timeouts})
	plan := vpsTestUpdateState(t, map[string]attr.Value{
		"timeouts":     timeouts,
		"installation": vpsTestInstallation(tfStr(vpsTestSshKey)),
		"power_state":  tfStr("SHUTOFF"),
	})

	resp := vpsTestUpdatePending(t, api, state, plan, vpsPrivateInstallationPending)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error: the reinstallation does not end in time")
	}
	if d := resp.Diagnostics.Errors()[0]; !strings.Contains(d.Summary(), "The reinstallation of VPS "+vpsTestId+" is still running") {
		t.Errorf("error %q: %q, want the reinstallation still running", d.Summary(), d.Detail())
	}
	if got := vpsTestActionNames(actions.list()); !slices.Equal(got, []string{"POST /rebuild"}) {
		t.Errorf("got requests %q, want the reinstallation only", got)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"installation.ssh_key": vpsTestSshKey,
		"power_state":          "ACTIVE",
	})
	if vpsTestPending(t, resp.Private) {
		t.Error("installation_pending is still set: the next plan would reinstall the VPS again")
	}
	if !vpsTestHasKey(t, resp.Private, vpsPrivateInstallationKnown) {
		t.Error("installation_known is not set after the reinstallation")
	}
}

// API v1 answers 409 while another task runs on the VPS: the rename is tried
// again.
func TestVPSResourceUpdateDisplayNameRetriesConflict(t *testing.T) {
	api, actions := vpsTestUpdateAPI(t)
	calls := 0
	api.handle("PUT "+vpsTestBase, func(body []byte) any {
		calls++
		if calls == 1 {
			return vpsTestStatus{code: http.StatusConflict, message: "A task is already running"}
		}
		actions.add("PUT "+vpsTestBase, body)
		api.change(vpsTestBase, func(vps map[string]any) { vps["displayName"] = "web-02" })
		return nil
	})
	state := vpsTestUpdateState(t, nil)
	plan := vpsTestUpdateState(t, map[string]attr.Value{"display_name": tfStr("web-02")})

	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	if calls != 2 {
		t.Errorf("PUT called %d times, want 2", calls)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{"current_state.display_name": "web-02"})
}
