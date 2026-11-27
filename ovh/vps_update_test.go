package ovh

import (
	"context"
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
// (whose type is internal to it), to an empty one, holding installation_pending
// when pending.
func vpsTestNewPrivate(t *testing.T, field any, pending bool) {
	t.Helper()
	v := reflect.ValueOf(field).Elem()
	v.Set(reflect.New(v.Type().Elem()))
	if !pending {
		return
	}
	private := v.Interface().(interface {
		SetKey(context.Context, string, []byte) diag.Diagnostics
	})
	if diags := private.SetKey(context.Background(), vpsPrivateInstallationPending, []byte("true")); diags.HasError() {
		t.Fatalf("setting the private state: %v", diags)
	}
}

// vpsTestPending tells whether a private state holds installation_pending.
func vpsTestPending(t *testing.T, private interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
}) bool {
	t.Helper()
	value, diags := private.GetKey(context.Background(), vpsPrivateInstallationPending)
	if diags.HasError() {
		t.Fatalf("reading the private state: %v", diags)
	}
	return value != nil
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

	resp := vpsTestUpdatePending(t, api, state, state, true)
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
func vpsTestModifyPlanUpdate(t *testing.T, api *vpsTestAPI, state tfsdk.State, config, plan tfsdk.Plan, pending bool) *resource.ModifyPlanResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	req := resource.ModifyPlanRequest{
		Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw},
		Plan:   plan,
		State:  state,
	}
	vpsTestNewPrivate(t, &req.Private, pending)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), false)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), false)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), false)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), false)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(state), true)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), false)
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
		resp := vpsTestUpdatePending(t, api, state, plan, true)
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
		resp := vpsTestModifyPlanUpdate(t, api, state, config, tfsdk.Plan(plan), true)
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
