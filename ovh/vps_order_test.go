package ovh

import (
	"context"
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// vpsTestCatalog is an extract of GET /order/catalog/public/vps?ovhSubsidiary=FR
// (2026-11-25): plan vps-2025-model1 as served, and a plan of another range.
var vpsTestCatalog = map[string]any{
	"catalogId": 1234,
	"locale":    map[string]any{"subsidiary": "FR", "currencyCode": "EUR"},
	"plans": []any{
		map[string]any{
			"planCode": "vps-2025-model1.LZ",
			"configurations": []any{
				map[string]any{"name": "vps_datacenter", "isCustom": false, "isMandatory": true, "values": []string{"EU-WEST-LZ-AMS"}},
				map[string]any{"name": "vps_os", "isCustom": false, "isMandatory": false, "values": []string{"Debian 13"}},
			},
		},
		map[string]any{
			"planCode": "vps-2025-model1",
			"configurations": []any{
				map[string]any{"name": "infrastructure", "isCustom": false, "isMandatory": false, "values": []string{"preproduction", "production"}},
				map[string]any{"name": "region", "isCustom": false, "isMandatory": false, "values": []string{"canada", "europe"}},
				map[string]any{"name": "vps_datacenter", "isCustom": false, "isMandatory": true, "values": vpsTestCatalogDatacenters},
				map[string]any{"name": "vps_install_rtm", "isCustom": false, "isMandatory": false, "values": []string{"if_available", "no"}},
				map[string]any{"name": "vps_os", "isCustom": false, "isMandatory": false, "values": vpsTestCatalogImages},
			},
		},
	},
}

// vpsTestCatalogImages is the vps_os configuration of plan vps-2025-model1
// (FR, 2026-11-25), shortened.
var vpsTestCatalogImages = []string{"AlmaLinux 9", "Debian 12", "Debian 13", "Ubuntu 24.04", "Windows Server 2025 Standard (Desktop)"}

func TestVPSOrderOfferFor(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/order/catalog/public/vps", vpsTestCatalog)
	config := api.config(t)

	offer, err := vpsOrderOfferFor(context.Background(), config.OVHClient, "FR", "vps-2025-model1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !slices.Equal(offer.Datacenters, vpsTestCatalogDatacenters) {
		t.Errorf("datacenters: got %v, want %v", offer.Datacenters, vpsTestCatalogDatacenters)
	}
	if !slices.Equal(offer.Images, vpsTestCatalogImages) {
		t.Errorf("images: got %v, want %v", offer.Images, vpsTestCatalogImages)
	}
}

func TestVPSOrderOfferForUnknownPlan(t *testing.T) {
	api := newVPSTestAPI()
	api.set("/1.0/order/catalog/public/vps", vpsTestCatalog)
	config := api.config(t)

	_, err := vpsOrderOfferFor(context.Background(), config.OVHClient, "FR", "vps-2025-model9")
	if err == nil {
		t.Fatal("got no error, want one for a plan absent from the catalog")
	}
	if !strings.Contains(err.Error(), "vps-2025-model9") {
		t.Errorf("error %q does not name the plan", err)
	}
}

func TestVPSOrderOfferForCatalogError(t *testing.T) {
	config := newVPSTestAPI().config(t) // no catalog: 404

	if _, err := vpsOrderOfferFor(context.Background(), config.OVHClient, "FR", "vps-2025-model1"); err == nil {
		t.Fatal("got no error, want the error of the catalog call")
	}
}

func TestVPSOrderModel(t *testing.T) {
	order, diags := vpsOrderModel(context.Background(), "vps-2025-model1", "GRA", "Debian 12", "P1M", "default", "FR")
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	if got := order.OvhSubsidiary.ValueString(); got != "FR" {
		t.Errorf("subsidiary: got %q, want FR", got)
	}
	if n := len(order.PlanOption.Elements()); n != 0 {
		t.Errorf("plan options: got %d, want none", n)
	}
	plans := order.Plan.Elements()
	if len(plans) != 1 {
		t.Fatalf("plans: got %d, want 1", len(plans))
	}

	got, err := json.Marshal(plans[0])
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	// itemId and domain are null in every PlanValue posted to a cart.
	want := `{"planCode":"vps-2025-model1","duration":"P1M","pricingMode":"default","quantity":1,"itemId":null,"domain":null,` +
		`"configuration":[{"label":"vps_datacenter","value":"GRA"},{"label":"vps_os","value":"Debian 12"}]}`
	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("plan:\n got  %s\n want %s", got, want)
	}
}

func TestVPSOrderImageForImageId(t *testing.T) {
	for _, tc := range []struct {
		images []string
		want   string
		ok     bool
	}{
		{[]string{"AlmaLinux 8 - cPanel", "Windows Server 2022", "AlmaLinux 9", "Debian 12"}, "AlmaLinux 9", true},
		{[]string{"Debian 12 - Plesk", "Debian 12"}, "Debian 12", true},
		{[]string{"Windows Server 2025 Standard (Desktop)", "Ubuntu 24.04"}, "Ubuntu 24.04", true},
		{[]string{"AlmaLinux 8 - cPanel", "Windows Server 2022"}, "", false},
		{nil, "", false},
	} {
		got, ok := vpsOrderImageForImageId(tc.images)
		if got != tc.want || ok != tc.ok {
			t.Errorf("%q: got %q, %t, want %q, %t", tc.images, got, ok, tc.want, tc.ok)
		}
	}
}
