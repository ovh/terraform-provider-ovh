package ovh

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestVpsResourceSchemaModelVersionValidation(t *testing.T) {
	ctx := context.Background()
	schema := VpsResourceSchema(ctx)
	model := schema.Attributes["model"].(resourceschema.SingleNestedAttribute)
	version := model.Attributes["version"].(resourceschema.StringAttribute)
	versionValidator := version.Validators[0]
	expectedDescription := `value must be one of: ["2013v1" "2014v1" "2015v1" "2017v1" "2017v2" "2017v3" "2018v1" "2018v2" "2019v1" "2025v1" "2027v1"]`

	if got := versionValidator.Description(ctx); got != expectedDescription {
		t.Fatalf("expected plan type validator description %q, got %q", expectedDescription, got)
	}

	for _, planType := range []string{
		"2013v1",
		"2014v1",
		"2015v1",
		"2017v1",
		"2017v2",
		"2017v3",
		"2018v1",
		"2018v2",
		"2019v1",
		"2025v1",
		"2027v1",
	} {
		t.Run(planType, func(t *testing.T) {
			var response validator.StringResponse

			versionValidator.ValidateString(ctx, validator.StringRequest{
				ConfigValue: types.StringValue(planType),
			}, &response)

			if response.Diagnostics.HasError() {
				t.Fatalf("expected plan type %q to be accepted, got diagnostics: %#v", planType, response.Diagnostics)
			}
		})
	}

	var response validator.StringResponse
	versionValidator.ValidateString(ctx, validator.StringRequest{
		ConfigValue: types.StringValue("unsupported-vps-plan"),
	}, &response)
	if !response.Diagnostics.HasError() {
		t.Fatal("expected an unsupported plan type to be rejected")
	}
}

// vpsTestServiceId is the service id of vpsTestId in the fake API.
const vpsTestServiceId = 12345

// vpsTestAddService serves GET /services?resourceName= and GET
// /services/{serviceId} for vpsTestId. Its billing differs from the defaults of
// the resource, and from the plan code of vpsTestFullState.
func vpsTestAddService(api *vpsTestAPI) {
	api.set("/1.0/services", []int{vpsTestServiceId})
	api.set(fmt.Sprintf("/1.0/services/%d", vpsTestServiceId), map[string]any{
		"billing": map[string]any{
			"plan":    map[string]any{"code": "vps-2025-model3"},
			"pricing": map[string]any{"duration": "P12M", "pricingMode": "upfront12"},
		},
	})
}

// vpsTestBillingCalls is the number of calls to the /services routes.
func vpsTestBillingCalls(api *vpsTestAPI) int {
	return api.count("GET /1.0/services") + api.count(fmt.Sprintf("GET /1.0/services/%d", vpsTestServiceId))
}

// vpsTestFullState is the state of vpsTestId as an apply leaves it, with
// values that differ from what API v1 answers, overridden by overrides (a nil
// value removes the attribute).
func vpsTestFullState(t *testing.T, overrides map[string]attr.Value) tfsdk.State {
	t.Helper()
	values := map[string]attr.Value{
		"id":           tfStr(vpsTestId),
		"plan_code":    tfStr("vps-2025-model1"),
		"display_name": tfStr("old-name"),
		"power_state":  tfStr("ACTIVE"),
		"location": types.ObjectValueMust(vpsLocationAttrTypes(), map[string]attr.Value{
			"region": tfStr("eu-west-sbg"),
		}),
		"installation": types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr("11111111-2222-4333-8444-555555555555"),
				"image_name": tfStr("Ubuntu 24.04"),
			}),
			"ssh_key":                  tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test"),
			"post_installation_script": tfStr("#!/bin/sh\ntrue\n"),
			"send_password_by_email":   types.BoolValue(false),
		}),
		"options": types.ObjectValueMust(vpsOptionsAttrTypes(), map[string]attr.Value{
			"auto_backup": types.ObjectValueMust(vpsAutoBackupAttrTypes(), map[string]attr.Value{
				"cron": tfStr("0 3 * * *"),
			}),
		}),
		"order": types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     tfStr("P1M"),
			"pricing_mode": tfStr("default"),
			"order_id":     types.Int64Value(42),
		}),
		"timeouts": types.ObjectValueMust(vpsTimeoutsAttrTypes(), map[string]attr.Value{
			"create": tfStr("2h"),
			"update": vpsTfStringPtr(nil),
			"delete": vpsTfStringPtr(nil),
		}),
	}
	for name, value := range overrides {
		if value == nil {
			delete(values, name)
			continue
		}
		values[name] = value
	}
	plan := vpsTestPlan(t, values)
	return tfsdk.State{Schema: plan.Schema, Raw: plan.Raw}
}

// vpsTestRead runs the Read of the resource on the state.
func vpsTestRead(t *testing.T, api *vpsTestAPI, state tfsdk.State) *resource.ReadResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: state.Schema, Raw: state.Raw.Copy()}}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)
	return resp
}

// vpsTestValue returns the value of a state at a dotted path, as a string:
// "<null>" for a null value or a value under a null object, "<unknown>", and
// "<object>" for a known object.
func vpsTestValue(t *testing.T, raw tftypes.Value, dotted string) string {
	t.Helper()
	value := raw
	for _, name := range strings.Split(dotted, ".") {
		if value.IsNull() {
			return "<null>"
		}
		var attributes map[string]tftypes.Value
		if err := value.As(&attributes); err != nil {
			t.Fatalf("%s: %s", dotted, err)
		}
		next, ok := attributes[name]
		if !ok {
			t.Fatalf("%s: no attribute %q", dotted, name)
		}
		value = next
	}

	switch {
	case value.IsNull():
		return "<null>"
	case !value.IsKnown():
		return "<unknown>"
	}
	switch {
	case value.Type().Is(tftypes.String):
		var s string
		_ = value.As(&s)
		return s
	case value.Type().Is(tftypes.Bool):
		var b bool
		_ = value.As(&b)
		return fmt.Sprint(b)
	case value.Type().Is(tftypes.Number):
		n := new(big.Float)
		_ = value.As(&n)
		return n.Text('f', -1)
	default:
		return "<object>"
	}
}

func vpsTestCheckValues(t *testing.T, raw tftypes.Value, want map[string]string) {
	t.Helper()
	for dotted, expected := range want {
		if got := vpsTestValue(t, raw, dotted); got != expected {
			t.Errorf("%s: got %q, want %q", dotted, got, expected)
		}
	}
}

func TestVPSResourceReadImport(t *testing.T) {
	ctx := context.Background()
	api := newVPSTestAPI()
	vpsTestAddService(api)
	r := &vpsResource{config: api.config(t)}

	s := vpsTestSchema()
	imported := &resource.ImportStateResponse{
		State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)},
	}
	r.ImportState(ctx, resource.ImportStateRequest{ID: vpsTestId}, imported)
	if imported.Diagnostics.HasError() {
		t.Fatalf("unexpected import error: %v", imported.Diagnostics)
	}
	vpsTestCheckValues(t, imported.State.Raw, map[string]string{
		"id":            vpsTestId,
		"plan_code":     "<null>",
		"location":      "<null>",
		"installation":  "<null>",
		"order":         "<null>",
		"current_state": "<null>",
	})

	resp := &resource.ReadResponse{State: tfsdk.State{Schema: s, Raw: imported.State.Raw.Copy()}}
	r.Read(ctx, resource.ReadRequest{State: imported.State}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":                                    vpsTestId,
		"plan_code":                             "vps-2025-model3",
		"location.region":                       "eu-west-gra",
		"display_name":                          "web-01",
		"power_state":                           "ACTIVE",
		"installation.from.image_id":            "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d",
		"installation.from.image_name":          "Debian 12",
		"installation.ssh_key":                  "<null>",
		"installation.post_installation_script": "<null>",
		"installation.send_password_by_email":   "true",
		"options.auto_backup":                   "<null>",
		"order.duration":                        "P12M",
		"order.pricing_mode":                    "upfront12",
		"order.order_id":                        "<null>",
		"timeouts":                              "<null>",
		"current_state.display_name":            "web-01",
		"current_state.power_state":             "ACTIVE",
		"current_state.location.region":         "eu-west-gra",
		"current_state.model.vcpus":             "4",
	})
	if got := api.count("GET /1.0/services"); got != 1 {
		t.Errorf("GET /services called %d times, want 1", got)
	}
	if got := api.count(fmt.Sprintf("GET /1.0/services/%d", vpsTestServiceId)); got != 1 {
		t.Errorf("GET /services/{serviceId} called %d times, want 1", got)
	}
}

func TestVPSResourceReadKeepsConfigured(t *testing.T) {
	api := newVPSTestAPI()
	resp := vpsTestRead(t, api, vpsTestFullState(t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":              vpsTestId,
		"plan_code":       "vps-2025-model1",
		"location.region": "eu-west-sbg",
		// Always read from API v1.
		"display_name":        "web-01",
		"power_state":         "ACTIVE",
		"options.auto_backup": "<null>",
		// Never read back.
		"installation.from.image_id":            "11111111-2222-4333-8444-555555555555",
		"installation.from.image_name":          "Ubuntu 24.04",
		"installation.ssh_key":                  "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test",
		"installation.post_installation_script": "#!/bin/sh\ntrue\n",
		"installation.send_password_by_email":   "false",
		"order.duration":                        "P1M",
		"order.pricing_mode":                    "default",
		"order.order_id":                        "42",
		"timeouts.create":                       "2h",
		// The current state is what API v1 answers.
		"current_state.display_name":                     "web-01",
		"current_state.installation.from.image_name":     "Debian 12",
		"current_state.location.region":                  "eu-west-gra",
		"current_state.options.auto_backup.enabled":      "false",
		"current_state.options.additional_disks.enabled": "false",
	})
}

func TestVPSResourceReadTransitionalKeepsPowerState(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase, map[string]any{
		"name":        vpsTestId,
		"displayName": "web-01",
		"state":       "installing",
		"vcore":       4,
		"memoryLimit": 8192,
		"model":       map[string]any{"name": "vps-2025-model3", "disk": 160},
	})
	resp := vpsTestRead(t, api, vpsTestFullState(t, map[string]attr.Value{"power_state": tfStr("SHUTOFF")}))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"power_state":               "SHUTOFF",
		"current_state.power_state": "UNKNOWN",
	})
}

func TestVPSResourceReadNotFound(t *testing.T) {
	api := newVPSTestAPI()
	resp := vpsTestRead(t, api, vpsTestFullState(t, map[string]attr.Value{"id": tfStr("vps-00000000.vps.ovh.net")}))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("a VPS that no longer exists must be removed from the state")
	}
}

func TestVPSResourceReadExpired(t *testing.T) {
	api := newVPSTestAPI()
	api.set(vpsTestBase, vpsTestExpired)
	state := vpsTestFullState(t, nil)
	resp := vpsTestRead(t, api, state)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for an expired VPS")
	}
	errs := resp.Diagnostics.Errors()
	if summary := errs[0].Summary(); !strings.Contains(summary, "expired") {
		t.Errorf("summary %q does not say the VPS is expired", summary)
	}
	for _, hint := range []string{"Renew", "terraform state rm"} {
		if detail := errs[0].Detail(); !strings.Contains(detail, hint) {
			t.Errorf("detail %q does not hold %q", detail, hint)
		}
	}
	if !resp.State.Raw.Equal(state.Raw) {
		t.Error("an expired VPS must stay in the state, unchanged")
	}
}

func TestVPSResourceReadNoBillingCallWhenKnown(t *testing.T) {
	api := newVPSTestAPI()
	vpsTestAddService(api)
	resp := vpsTestRead(t, api, vpsTestFullState(t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	if got := vpsTestBillingCalls(api); got != 0 {
		t.Errorf("/services called %d times, want 0", got)
	}
}

func TestVPSResourceReadFillsOnlyNullBilling(t *testing.T) {
	api := newVPSTestAPI()
	vpsTestAddService(api)
	resp := vpsTestRead(t, api, vpsTestFullState(t, map[string]attr.Value{"order": nil}))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"plan_code":          "vps-2025-model1",
		"order.duration":     "P12M",
		"order.pricing_mode": "upfront12",
		"order.order_id":     "<null>",
	})
}

func TestVPSResourceReadBillingError(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/services", vpsTestStatus{code: http.StatusInternalServerError, message: "Internal server error"})
	state := vpsTestFullState(t, map[string]attr.Value{"plan_code": nil})
	resp := vpsTestRead(t, api, state)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the billing of the VPS cannot be read")
	}
	if !resp.State.Raw.Equal(state.Raw) {
		t.Error("the state must stay unchanged")
	}
}

func TestServiceIdFromResourceNameEmpty(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/services", []int{})
	_, err := serviceIdFromResourceName(api.config(t).OVHClient, vpsTestId)
	if err == nil {
		t.Fatal("expected an error when no service has the resource name")
	}
	if !strings.Contains(err.Error(), vpsTestId) {
		t.Errorf("error %q does not name the resource", err)
	}
}
