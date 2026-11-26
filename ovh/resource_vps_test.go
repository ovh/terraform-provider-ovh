package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	acc "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

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

// vpsTestOrderId and vpsTestCartId identify the order of the fake API.
const (
	vpsTestOrderId = 99
	vpsTestCartId  = "cart-1"
	vpsTestItemId  = 7
)

// vpsTestFastPolls shortens the waits of the resource for the test.
func vpsTestFastPolls(t *testing.T) {
	t.Helper()
	delay, interval := vpsPollDelay, vpsPollInterval
	vpsPollDelay, vpsPollInterval = time.Millisecond, time.Millisecond
	t.Cleanup(func() { vpsPollDelay, vpsPollInterval = delay, interval })
}

// vpsTestCart holds the bodies of the cart requests of the fake API.
type vpsTestCart struct {
	mu             sync.Mutex
	item           map[string]any
	configurations []map[string]any
}

// vpsTestOrderAPI is newVPSTestAPI with the routes of an order of
// vps-2025-model1 paid by the default payment method and delivered as
// vpsTestId.
func vpsTestOrderAPI(t *testing.T) (*vpsTestAPI, *vpsTestCart) {
	t.Helper()
	vpsTestFastPolls(t)
	api := newVPSTestAPI()
	cart := &vpsTestCart{}
	decode := func(body []byte) map[string]any {
		var m map[string]any
		if err := json.Unmarshal(body, &m); err != nil {
			t.Errorf("decoding %s: %s", body, err)
		}
		return m
	}
	orderBase := fmt.Sprintf("/1.0/me/order/%d", vpsTestOrderId)
	cartBase := "/1.0/order/cart/" + vpsTestCartId

	api.set("/1.0/me", map[string]any{"ovhSubsidiary": "FR"})
	api.set("/1.0/order/catalog/public/vps", vpsTestCatalog)
	api.handle("POST /1.0/order/cart", func([]byte) any { return map[string]any{"cartId": vpsTestCartId} })
	api.handle("POST "+cartBase+"/assign", func([]byte) any { return nil })
	api.handle("POST "+cartBase+"/vps", func(body []byte) any {
		cart.mu.Lock()
		defer cart.mu.Unlock()
		cart.item = decode(body)
		return map[string]any{"cartId": vpsTestCartId, "itemId": vpsTestItemId}
	})
	api.handle(fmt.Sprintf("POST %s/item/%d/configuration", cartBase, vpsTestItemId), func(body []byte) any {
		cart.mu.Lock()
		defer cart.mu.Unlock()
		cart.configurations = append(cart.configurations, decode(body))
		return map[string]any{"id": len(cart.configurations)}
	})
	api.set("/1.0/me/payment/method", []int64{1})
	api.handle("POST "+cartBase+"/checkout", func([]byte) any {
		return map[string]any{"orderId": vpsTestOrderId, "prices": map[string]any{"withoutTax": map[string]any{"value": 5.99}}}
	})
	api.handle("POST "+orderBase+"/pay", func([]byte) any { return nil })
	api.set(orderBase+"/status", "delivered")
	api.set(orderBase+"/details", []int64{1})
	api.set(orderBase+"/details/1/extension", map[string]any{"order": map[string]any{"plan": map[string]any{"code": "vps-2025-model1"}}})
	api.set(orderBase+"/details/1", map[string]any{"domain": vpsTestId, "orderDetailId": 1})
	return api, cart
}

// vpsTestPaidOrders is the number of orders paid through the fake API.
func vpsTestPaidOrders(api *vpsTestAPI) int {
	return api.count(fmt.Sprintf("POST /1.0/me/order/%d/pay", vpsTestOrderId))
}

// vpsTestCreateConfig is the configuration of a VPS to create, with overrides
// (dotted paths).
func vpsTestCreateConfig(t *testing.T, overrides map[string]attr.Value) tfsdk.Config {
	t.Helper()
	values := map[string]attr.Value{
		"plan_code":                    tfStr("vps-2025-model1"),
		"location.region":              tfStr("eu-west-gra"),
		"installation.from.image_name": tfStr("Debian 12"),
	}
	maps.Copy(values, overrides)
	plan := vpsTestPlan(t, values)
	return tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}
}

// vpsTestCreatePlan is the plan of vpsTestCreateConfig as the framework makes
// it: defaults set, the computed attributes the configuration leaves out
// unknown.
func vpsTestCreatePlan(t *testing.T, overrides map[string]attr.Value) tfsdk.Plan {
	t.Helper()
	unknown := ovhtypes.TfStringValue{StringValue: types.StringUnknown()}
	values := map[string]attr.Value{
		"id":                                  unknown,
		"plan_code":                           tfStr("vps-2025-model1"),
		"display_name":                        unknown,
		"power_state":                         tfStr("ACTIVE"),
		"location.region":                     tfStr("eu-west-gra"),
		"installation.from.image_name":        tfStr("Debian 12"),
		"installation.from.image_id":          unknown,
		"installation.send_password_by_email": types.BoolValue(true),
		"options":                             types.ObjectUnknown(vpsOptionsAttrTypes()),
		"order":                               types.ObjectUnknown(vpsOrderAttrTypes()),
		"current_state":                       types.ObjectUnknown(vpsCurrentStateAttrTypes()),
	}
	maps.Copy(values, overrides)
	return vpsTestPlan(t, values)
}

// vpsTestCreate runs the Create of the resource on the plan.
func vpsTestCreate(t *testing.T, api *vpsTestAPI, plan tfsdk.Plan) *resource.CreateResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema, Raw: tftypes.NewValue(plan.Raw.Type(), nil)}}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)
	return resp
}

func TestVPSResourceCreate(t *testing.T) {
	api, cart := vpsTestOrderAPI(t)
	// Delivered, the VPS is installing for two reads: Create waits for it.
	gets := 0
	api.handle("GET "+vpsTestBase, func([]byte) any {
		gets++
		state := "running"
		if gets <= 2 {
			state = "installing"
		}
		return map[string]any{
			"name": vpsTestId, "displayName": vpsTestId, "state": state, "vcore": 4, "memoryLimit": 8192,
			"model": map[string]any{"name": "vps-2025-model1", "disk": 160},
		}
	})

	resp := vpsTestCreate(t, api, vpsTestCreatePlan(t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected create error: %v", resp.Diagnostics)
	}

	if got, want := cart.item["planCode"], "vps-2025-model1"; got != want {
		t.Errorf("cart item planCode = %v, want %v", got, want)
	}
	if got, want := cart.item["duration"], "P1M"; got != want {
		t.Errorf("cart item duration = %v, want %v", got, want)
	}
	if got, want := cart.item["pricingMode"], "default"; got != want {
		t.Errorf("cart item pricingMode = %v, want %v", got, want)
	}
	configurations := map[string]any{}
	for _, c := range cart.configurations {
		configurations[fmt.Sprint(c["label"])] = c["value"]
	}
	wantConfigurations := map[string]any{"vps_datacenter": "GRA", "vps_os": "Debian 12"}
	if !reflect.DeepEqual(configurations, wantConfigurations) {
		t.Errorf("cart configurations = %v, want %v", configurations, wantConfigurations)
	}
	if got := vpsTestPaidOrders(api); got != 1 {
		t.Errorf("%d orders paid, want 1", got)
	}

	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("the state holds unknown values: %v", resp.State.Raw)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":                                  vpsTestId,
		"plan_code":                           "vps-2025-model1",
		"location.region":                     "eu-west-gra",
		"display_name":                        vpsTestId,
		"power_state":                         "ACTIVE",
		"installation.from.image_name":        "Debian 12",
		"installation.from.image_id":          "<null>",
		"installation.send_password_by_email": "true",
		"order.duration":                      "P1M",
		"order.pricing_mode":                  "default",
		"order.order_id":                      fmt.Sprint(vpsTestOrderId),
		"current_state.power_state":           "ACTIVE",
		"current_state.installation.from.image_name": "Debian 12",
	})
	if got := vpsTestBillingCalls(api); got != 0 {
		t.Errorf("/services called %d times, want 0", got)
	}
}

func TestVPSResourceCreateOrderSettings(t *testing.T) {
	api, cart := vpsTestOrderAPI(t)
	plan := vpsTestCreatePlan(t, map[string]attr.Value{
		"order": types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     tfStr("P12M"),
			"pricing_mode": tfStr("upfront12"),
			"order_id":     types.Int64Unknown(),
		}),
	})
	resp := vpsTestCreate(t, api, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected create error: %v", resp.Diagnostics)
	}
	if got, want := cart.item["duration"], "P12M"; got != want {
		t.Errorf("cart item duration = %v, want %v", got, want)
	}
	if got, want := cart.item["pricingMode"], "upfront12"; got != want {
		t.Errorf("cart item pricingMode = %v, want %v", got, want)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"order.duration":     "P12M",
		"order.pricing_mode": "upfront12",
		"order.order_id":     fmt.Sprint(vpsTestOrderId),
	})
}

func TestVPSResourceCreateRegionNotInCatalog(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	resp := vpsTestCreate(t, api, vpsTestCreatePlan(t, map[string]attr.Value{"location.region": tfStr("eu-west-xyz")}))
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for a region the plan cannot be ordered in")
	}
	detail := resp.Diagnostics.Errors()[0].Detail()
	for _, want := range []string{"eu-west-xyz", "eu-west-gra", "eu-west-rbx"} {
		if !strings.Contains(detail, want) {
			t.Errorf("detail %q does not list %q", detail, want)
		}
	}
	if got := api.count("POST /1.0/order/cart"); got != 0 {
		t.Errorf("POST /order/cart called %d times, want 0", got)
	}
	if !resp.State.Raw.IsNull() {
		t.Error("nothing must be saved in the state")
	}
}

func TestVPSResourceCreateImageNotInCatalog(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	resp := vpsTestCreate(t, api, vpsTestCreatePlan(t, map[string]attr.Value{"installation.from.image_name": tfStr("Debian 11")}))
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for an image the plan cannot be ordered with")
	}
	detail := resp.Diagnostics.Errors()[0].Detail()
	for _, want := range []string{"Debian 11", "Debian 12", "Ubuntu 24.04"} {
		if !strings.Contains(detail, want) {
			t.Errorf("detail %q does not list %q", detail, want)
		}
	}
	if got := api.count("POST /1.0/order/cart"); got != 0 {
		t.Errorf("POST /order/cart called %d times, want 0", got)
	}
}

func TestVPSResourceCreateDeliveryTimeout(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	api.set(fmt.Sprintf("/1.0/me/order/%d/status", vpsTestOrderId), "delivering")
	plan := vpsTestCreatePlan(t, map[string]attr.Value{
		"timeouts": types.ObjectValueMust(vpsTimeoutsAttrTypes(), map[string]attr.Value{
			"create": tfStr("100ms"),
			"update": vpsTfStringPtr(nil),
			"delete": vpsTfStringPtr(nil),
		}),
	})
	resp := vpsTestCreate(t, api, plan)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the order is not delivered in time")
	}
	if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, "terraform untaint") {
		t.Errorf("detail %q does not tell to untaint the resource", detail)
	}
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("the state holds unknown values: %v", resp.State.Raw)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":             "<null>",
		"order.order_id": fmt.Sprint(vpsTestOrderId),
	})

	// Delivered later: the next Read waits for it again, and fills the state.
	api.set(fmt.Sprintf("/1.0/me/order/%d/status", vpsTestOrderId), "delivered")
	read := vpsTestRead(t, api, resp.State)
	if read.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", read.Diagnostics)
	}
	vpsTestCheckValues(t, read.State.Raw, map[string]string{
		"id":                         vpsTestId,
		"order.order_id":             fmt.Sprint(vpsTestOrderId),
		"current_state.display_name": "web-01",
	})
	if got := vpsTestPaidOrders(api); got != 1 {
		t.Errorf("%d orders paid, want 1", got)
	}
}

// vpsTestOrderStatusPath is the status route of the order of the fake API.
var vpsTestOrderStatusPath = fmt.Sprintf("/1.0/me/order/%d/status", vpsTestOrderId)

// vpsTestOrderedState is the state Create leaves when the order of the fake
// API is not delivered in time: its order id, no id.
func vpsTestOrderedState(t *testing.T) tfsdk.State {
	t.Helper()
	return vpsTestFullState(t, map[string]attr.Value{
		"id": nil,
		"order": types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     tfStr("P1M"),
			"pricing_mode": tfStr("default"),
			"order_id":     types.Int64Value(vpsTestOrderId),
		}),
	})
}

// Read never waits: a plan does not block on an order being delivered.
func TestVPSResourceReadDeliveryStillPending(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	api.set(vpsTestOrderStatusPath, "delivering")
	state := vpsTestOrderedState(t)
	resp := vpsTestRead(t, api, state)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	warnings := resp.Diagnostics.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1: %v", len(warnings), resp.Diagnostics)
	}
	if text := warnings[0].Summary() + warnings[0].Detail(); !strings.Contains(text, fmt.Sprint(vpsTestOrderId)) {
		t.Errorf("warning %q does not name the order", text)
	}
	if !resp.State.Raw.Equal(state.Raw) {
		t.Error("the state must stay unchanged")
	}
	if got := api.count("GET " + vpsTestOrderStatusPath); got != 1 {
		t.Errorf("GET /me/order/{orderId}/status called %d times, want 1", got)
	}
}

// Delivered but still installing: Read sets the id and reads the VPS as it is,
// without waiting for the installation.
func TestVPSResourceReadDeliveredOrder(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	api.set(vpsTestBase, map[string]any{
		"name": vpsTestId, "displayName": vpsTestId, "state": "installing", "vcore": 4, "memoryLimit": 8192,
		"model": map[string]any{"name": "vps-2025-model1", "disk": 160},
	})
	resp := vpsTestRead(t, api, vpsTestOrderedState(t))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":                        vpsTestId,
		"order.order_id":            fmt.Sprint(vpsTestOrderId),
		"power_state":               "ACTIVE",
		"current_state.power_state": "UNKNOWN",
	})
	if got := api.count("GET " + vpsTestOrderStatusPath); got != 1 {
		t.Errorf("GET /me/order/{orderId}/status called %d times, want 1", got)
	}
	if got := api.count("GET " + vpsTestBase); got != 1 {
		t.Errorf("GET /vps/{serviceName} called %d times, want 1", got)
	}
}

// Delivered, but API v1 does not serve the VPS yet: the id is kept, so that the
// next apply does not order a new VPS.
func TestVPSResourceReadDeliveredOrderNotReadableYet(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	api.set(vpsTestBase, nil)
	resp := vpsTestRead(t, api, vpsTestOrderedState(t))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read error: %v", resp.Diagnostics)
	}
	if len(resp.Diagnostics.Warnings()) != 1 {
		t.Errorf("want 1 warning: %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("the VPS of a delivered order must stay in the state")
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{
		"id":             vpsTestId,
		"order.order_id": fmt.Sprint(vpsTestOrderId),
	})
}

func TestVPSResourceModifyPlanRefusesPostOrderSettings(t *testing.T) {
	cases := []struct {
		attribute string
		value     attr.Value
		refused   bool
	}{
		{"installation.from.image_id", tfStr("11111111-2222-4333-8444-555555555555"), true},
		{"installation.ssh_key", tfStr("ssh-ed25519 AAAA"), true},
		{"installation.post_installation_script", tfStr("#!/bin/sh"), true},
		{"installation.send_password_by_email", types.BoolValue(false), true},
		{"installation.send_password_by_email", types.BoolValue(true), false},
		{"power_state", tfStr("SHUTOFF"), true},
		{"power_state", tfStr("ACTIVE"), false},
		{"display_name", tfStr("web-01"), true},
		{"options.auto_backup.cron", tfStr("0 3 * * *"), true},
		{"order.duration", tfStr("P12M"), false},
	}
	ctx := context.Background()
	r := &vpsResource{}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s=%v", c.attribute, c.value), func(t *testing.T) {
			config := vpsTestCreateConfig(t, map[string]attr.Value{c.attribute: c.value})
			plan := tfsdk.Plan{Schema: config.Schema, Raw: config.Raw}
			null := tftypes.NewValue(config.Raw.Type(), nil)

			// Creation.
			resp := &resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(ctx, resource.ModifyPlanRequest{
				Config: config,
				Plan:   plan,
				State:  tfsdk.State{Schema: config.Schema, Raw: null},
			}, resp)
			if got := resp.Diagnostics.HasError(); got != c.refused {
				t.Fatalf("refused = %t, want %t: %v", got, c.refused, resp.Diagnostics)
			}
			if c.refused {
				d := resp.Diagnostics.Errors()[0]
				// Setting it after the creation is refused too, or only
				// recorded: the detail does not tell to set it then.
				if !strings.HasPrefix(d.Detail(), "Setting "+c.attribute) || !strings.HasSuffix(d.Detail(), " when creating a VPS is not supported yet") {
					t.Errorf("detail %q, want \"Setting %s… when creating a VPS is not supported yet\"", d.Detail(), c.attribute)
				}
				if strings.Contains(d.Detail(), "PCINT") {
					t.Errorf("detail %q names a ticket", d.Detail())
				}
			}

			// The VPS exists: nothing to refuse here.
			state := vpsTestFullState(t, nil)
			resp = &resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan, State: state}, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("unexpected error when the VPS exists: %v", resp.Diagnostics)
			}
		})
	}
}

// vpsTestUpdate runs the Update of the resource from the state to the plan,
// the configuration being the plan.
func vpsTestUpdate(t *testing.T, api *vpsTestAPI, state tfsdk.State, plan tfsdk.State) *resource.UpdateResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: state.Schema, Raw: state.Raw.Copy()}}
	r.Update(context.Background(), resource.UpdateRequest{
		Plan:   tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw},
		Config: tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw},
		State:  state,
	}, resp)
	return resp
}

// vpsTestRequests is the number of requests the fake API got.
func vpsTestRequests(api *vpsTestAPI) int {
	api.mu.Lock()
	defer api.mu.Unlock()
	return len(api.requests)
}

func TestVPSResourceUpdateTimeouts(t *testing.T) {
	api := newVPSTestAPI()
	state := vpsTestFullState(t, nil)
	plan := vpsTestFullState(t, map[string]attr.Value{
		"timeouts": types.ObjectValueMust(vpsTimeoutsAttrTypes(), map[string]attr.Value{
			"create": tfStr("3h"),
			"update": vpsTfStringPtr(nil),
			"delete": vpsTfStringPtr(nil),
		}),
	})
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

func TestVPSResourceUpdateSshKeyAfterImport(t *testing.T) {
	api := newVPSTestAPI()
	installation := func(sshKey attr.Value) attr.Value {
		return types.ObjectValueMust(vpsResourceInstallationAttrTypes(), map[string]attr.Value{
			"from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
				"image_id":   tfStr("3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d"),
				"image_name": tfStr("Debian 12"),
			}),
			"ssh_key":                  sshKey,
			"post_installation_script": vpsTfStringPtr(nil),
			"send_password_by_email":   types.BoolValue(true),
		})
	}
	state := vpsTestFullState(t, map[string]attr.Value{"installation": installation(vpsTfStringPtr(nil))})
	plan := vpsTestFullState(t, map[string]attr.Value{"installation": installation(tfStr("ssh-ed25519 AAAA"))})
	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{"installation.ssh_key": "ssh-ed25519 AAAA"})
	if got := vpsTestRequests(api); got != 0 {
		t.Errorf("%d API calls, want 0", got)
	}
}

func TestVPSResourceUpdateKeepsComputedFromState(t *testing.T) {
	api := newVPSTestAPI()
	state := vpsTestFullState(t, nil)
	plan := vpsTestFullState(t, map[string]attr.Value{
		"order": types.ObjectValueMust(vpsOrderAttrTypes(), map[string]attr.Value{
			"duration":     tfStr("P12M"),
			"pricing_mode": tfStr("default"),
			"order_id":     types.Int64Unknown(),
		}),
	})
	resp := vpsTestUpdate(t, api, state, plan)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected update error: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("the state holds unknown values: %v", resp.State.Raw)
	}
	vpsTestCheckValues(t, resp.State.Raw, map[string]string{"order.duration": "P12M", "order.order_id": "42"})
}

func TestVPSResourceUpdateDisplayName(t *testing.T) {
	api := newVPSTestAPI()
	state := vpsTestFullState(t, nil)
	plan := vpsTestFullState(t, map[string]attr.Value{"display_name": tfStr("web-02")})
	resp := vpsTestUpdate(t, api, state, plan)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error: updating a VPS is not supported yet")
	}
	if summary := resp.Diagnostics.Errors()[0].Summary(); summary != "Updating a VPS is not supported yet" {
		t.Errorf("summary = %q", summary)
	}
	if !resp.State.Raw.Equal(state.Raw) {
		t.Error("the state must stay unchanged")
	}
	if got := vpsTestRequests(api); got != 0 {
		t.Errorf("%d API calls, want 0", got)
	}
}

// vpsTestTerminationAPI is newVPSTestAPI with the termination routes of
// vpsTestId: POST /terminate sends the termination email.
func vpsTestTerminationAPI(t *testing.T) (*vpsTestAPI, *[]string) {
	t.Helper()
	api := newVPSTestAPI()
	var tokens []string
	api.set("/1.0/me/notification/email/history", []int64{1})
	api.set("/1.0/me/notification/email/history/1", map[string]any{"id": 1, "body": "Welcome"})
	api.handle("POST "+vpsTestBase+"/terminate", func([]byte) any {
		api.set("/1.0/me/notification/email/history", []int64{1, 2})
		api.set("/1.0/me/notification/email/history/2", map[string]any{
			"id":   2,
			"body": "Termination of " + vpsTestId + ": https://www.ovh.com/manager/#/billing/confirmTerminate?id=abc123&token=tok456 to confirm",
		})
		return map[string]any{}
	})
	api.handle("POST "+vpsTestBase+"/confirmTermination", func(body []byte) any {
		var opts ConfirmTerminationOpts
		_ = json.Unmarshal(body, &opts)
		tokens = append(tokens, opts.Token)
		return "terminated"
	})
	return api, &tokens
}

func vpsTestDelete(t *testing.T, api *vpsTestAPI, state tfsdk.State) *resource.DeleteResponse {
	t.Helper()
	r := &vpsResource{config: api.config(t)}
	resp := &resource.DeleteResponse{State: tfsdk.State{Schema: state.Schema, Raw: state.Raw.Copy()}}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)
	return resp
}

func TestVPSResourceDelete(t *testing.T) {
	api, tokens := vpsTestTerminationAPI(t)
	resp := vpsTestDelete(t, api, vpsTestFullState(t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected delete error: %v", resp.Diagnostics)
	}
	if !reflect.DeepEqual(*tokens, []string{"tok456"}) {
		t.Errorf("confirmation tokens = %v, want [tok456]", *tokens)
	}
}

func TestVPSResourceDeleteExpired(t *testing.T) {
	api, tokens := vpsTestTerminationAPI(t)
	api.handle("POST "+vpsTestBase+"/terminate", func([]byte) any { return vpsTestExpired })
	resp := vpsTestDelete(t, api, vpsTestFullState(t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected delete error: %v", resp.Diagnostics)
	}
	if len(*tokens) != 0 {
		t.Errorf("termination confirmed with %v, want no confirmation", *tokens)
	}
}

func TestVPSResourceDeleteNotDelivered(t *testing.T) {
	api, _ := vpsTestTerminationAPI(t)
	resp := vpsTestDelete(t, api, vpsTestFullState(t, map[string]attr.Value{"id": nil}))
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error: the VPS of the order is not delivered")
	}
	if got := api.count("POST " + vpsTestBase + "/terminate"); got != 0 {
		t.Errorf("POST /terminate called %d times, want 0", got)
	}
}

// A value unknown at plan time, such as a variable from another resource, is
// refused by Create, before the order.
func TestVPSResourceCreateRefusesPostOrderSettings(t *testing.T) {
	api, _ := vpsTestOrderAPI(t)
	plan := vpsTestCreatePlan(t, map[string]attr.Value{"installation.ssh_key": tfStr("ssh-ed25519 AAAA")})
	resp := vpsTestCreate(t, api, plan)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected installation.ssh_key to be refused")
	}
	if got := api.count("POST /1.0/order/cart"); got != 0 {
		t.Errorf("POST /order/cart called %d times, want 0", got)
	}
}

const testAccVPSResourceConfig = `
resource "ovh_vps" "vps" {
  plan_code = %q

  location = {
    region = %q
  }

  installation = {
    from = {
      image_name = %q
    }
  }
}
`

// TestAccVPSResource_order orders a VPS, which is billed, then terminates it.
func TestAccVPSResource_order(t *testing.T) {
	planCode := os.Getenv("OVH_VPS_PLAN_CODE")
	region := os.Getenv("OVH_VPS_REGION")
	imageName := os.Getenv("OVH_VPS_IMAGE_NAME")

	acc.Test(t, acc.TestCase{
		PreCheck: func() {
			testAccPreCheckOrderVPS(t)
			checkEnvOrSkip(t, "OVH_VPS_PLAN_CODE")
			checkEnvOrSkip(t, "OVH_VPS_REGION")
			checkEnvOrSkip(t, "OVH_VPS_IMAGE_NAME")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []acc.TestStep{
			{
				Config: fmt.Sprintf(testAccVPSResourceConfig, planCode, region, imageName),
				Check: acc.ComposeTestCheckFunc(
					acc.TestCheckResourceAttrSet("ovh_vps.vps", "id"),
					acc.TestCheckResourceAttrSet("ovh_vps.vps", "order.order_id"),
					acc.TestCheckResourceAttr("ovh_vps.vps", "plan_code", planCode),
					acc.TestCheckResourceAttr("ovh_vps.vps", "current_state.location.region", region),
					acc.TestCheckResourceAttr("ovh_vps.vps", "installation.from.image_name", imageName),
				),
			},
		},
	})
}

// TestAccVPSResource_import imports the existing VPS OVH_VPS. The import step
// does not persist its state (no import state persistence): the test ends
// with an empty state, so it never terminates the VPS. The configuration is
// only validated: an import does not plan it.
func TestAccVPSResource_import(t *testing.T) {
	vps := os.Getenv("OVH_VPS")

	acc.Test(t, acc.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			checkEnvOrSkip(t, "OVH_VPS")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []acc.TestStep{
			{
				Config:        fmt.Sprintf(testAccVPSResourceConfig, "vps-2025-model1", "eu-west-gra", "Debian 12"),
				ResourceName:  "ovh_vps.vps",
				ImportState:   true,
				ImportStateId: vps,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("imported %d resources, want 1", len(states))
					}
					attributes := states[0].Attributes
					if attributes["id"] != vps {
						return fmt.Errorf("id = %q, want %q", attributes["id"], vps)
					}
					for _, name := range []string{"plan_code", "order.duration", "order.pricing_mode", "location.region", "installation.from.image_name", "current_state.power_state"} {
						if attributes[name] == "" {
							return fmt.Errorf("%s is not set after the import", name)
						}
					}
					return nil
				},
			},
		},
	})
}
