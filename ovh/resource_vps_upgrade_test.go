package ovh

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// vpsTestV0OrderedState is the state of a VPS ordered by the former ovh_vps
// resource (schema version 0), as Terraform stores it: every attribute of the
// former schema, from ovh/resource_vps_gen.go and the order attributes of
// ovh/order_resource_gen.go, before this version.
const vpsTestV0OrderedState = `{
  "id": "vps-1a2b3c4d.vps.ovh.net",
  "service_name": "vps-1a2b3c4d.vps.ovh.net",
  "cluster": "",
  "display_name": "web-01",
  "do_not_send_password": true,
  "iam": {
    "display_name": "web-01",
    "id": "5f0b6a1e-3c2d-4e5f-8a9b-0c1d2e3f4a5b",
    "tags": {"ovh:env": "prod"},
    "urn": "urn:v1:eu:resource:vps:vps-1a2b3c4d.vps.ovh.net"
  },
  "image_id": null,
  "keymap": null,
  "memory_limit": 8192,
  "model": {
    "available_options": ["additionalDisk", "automatedBackup", "snapshot"],
    "datacenter": ["gra", "sbg"],
    "disk": 75,
    "maximum_additionnal_ip": 16,
    "memory": 8192,
    "name": "vps-2025-model1",
    "offer": "VPS vps-2025-model1",
    "vcore": 4,
    "version": "2025v1"
  },
  "monitoring_ip_blocks": ["92.222.184.0/24"],
  "name": "vps-1a2b3c4d.vps.ovh.net",
  "netboot_mode": "local",
  "offer_type": "ssd",
  "public_ssh_key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMigrated user@laptop",
  "sla_monitoring": false,
  "state": "running",
  "vcore": 4,
  "zone": "Region OpenStack: os-gra9",
  "order": {
    "date": "2025-03-04T10:00:00+01:00",
    "details": [
      {"description": "VPS vps-2025-model1", "detail_type": "DURATION", "domain": "*001", "order_detail_id": 123456, "quantity": "1"}
    ],
    "expiration_date": "2025-03-18T10:00:00+01:00",
    "order_id": 42
  },
  "ovh_subsidiary": "FR",
  "plan": [
    {
      "configuration": [
        {"label": "vps_datacenter", "value": "GRA"},
        {"label": "vps_os", "value": "Debian 12"}
      ],
      "duration": "P1M",
      "item_id": null,
      "plan_code": "vps-2025-model1",
      "pricing_mode": "default",
      "quantity": null
    }
  ],
  "plan_option": [
    {
      "configuration": [{"label": "size", "value": "100"}],
      "duration": "P1M",
      "plan_code": "option-disk-100",
      "pricing_mode": "default",
      "quantity": 1
    }
  ]
}`

// vpsTestV0ImportedState is the state of a VPS the former ovh_vps resource
// imported: no order, no plan, no password setting.
const vpsTestV0ImportedState = `{
  "id": "vps-1a2b3c4d.vps.ovh.net",
  "service_name": "vps-1a2b3c4d.vps.ovh.net",
  "cluster": "",
  "display_name": "web-01",
  "do_not_send_password": null,
  "iam": null,
  "image_id": null,
  "keymap": null,
  "memory_limit": 8192,
  "model": null,
  "monitoring_ip_blocks": [],
  "name": "vps-1a2b3c4d.vps.ovh.net",
  "netboot_mode": "local",
  "offer_type": "ssd",
  "public_ssh_key": null,
  "sla_monitoring": false,
  "state": "running",
  "vcore": 4,
  "zone": "Region OpenStack: os-gra9",
  "order": null,
  "ovh_subsidiary": null,
  "plan": [],
  "plan_option": null
}`

// vpsTestUpgrade upgrades a state of schema version 0 through the protocol
// server of the provider, as Terraform does, and returns the upgraded state.
func vpsTestUpgrade(t *testing.T, rawState string) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	server, err := providerserver.NewProtocol6WithError(&OvhProvider{})()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
		TypeName: "ovh_vps",
		Version:  0,
		RawState: &tfprotov6.RawState{JSON: []byte(rawState)},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		t.Errorf("upgrade: %s: %s: %s", d.Severity, d.Summary, d.Detail)
	}
	if t.Failed() {
		t.FailNow()
	}
	s := vpsTestSchema()
	raw, err := resp.UpgradedState.Unmarshal(s.Type().TerraformType(ctx))
	if err != nil {
		t.Fatal(err)
	}
	return tfsdk.State{Schema: s, Raw: raw}
}

func TestUpgradeVPSStateV0Ordered(t *testing.T) {
	state := vpsTestUpgrade(t, vpsTestV0OrderedState)
	vpsTestCheckValues(t, state.Raw, map[string]string{
		"id":                                    vpsTestId,
		"plan_code":                             "vps-2025-model1",
		"display_name":                          "web-01",
		"power_state":                           "<null>",
		"location":                              "<null>",
		"installation.from":                     "<null>",
		"installation.ssh_key":                  "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMigrated user@laptop",
		"installation.post_installation_script": "<null>",
		"installation.send_password_by_email":   "false",
		"options":                               "<null>",
		"order.duration":                        "P1M",
		"order.pricing_mode":                    "default",
		"order.order_id":                        "42",
		"order.options":                         "<object>",
		"timeouts":                              "<null>",
		"current_state":                         "<null>",
	})
	var options []string
	if diags := state.GetAttribute(context.Background(), path.Root("order").AtName("options"), &options); diags.HasError() {
		t.Fatalf("reading order.options: %v", diags)
	}
	if want := []string{"option-disk-100"}; !slices.Equal(options, want) {
		t.Errorf("order.options = %q, want %q", options, want)
	}
}

// An empty plan_option gives a null order.options.
func TestUpgradeVPSStateV0NoPlanOption(t *testing.T) {
	var v0 map[string]any
	if err := json.Unmarshal([]byte(vpsTestV0OrderedState), &v0); err != nil {
		t.Fatal(err)
	}
	for name, planOption := range map[string]string{"empty": `[]`, "null": `null`} {
		t.Run(name, func(t *testing.T) {
			v0["plan_option"] = json.RawMessage(planOption)
			raw, err := json.Marshal(v0)
			if err != nil {
				t.Fatal(err)
			}
			vpsTestCheckValues(t, vpsTestUpgrade(t, string(raw)).Raw, map[string]string{
				"order.order_id": "42",
				"order.options":  "<null>",
			})
		})
	}
}

func TestUpgradeVPSStateV0Imported(t *testing.T) {
	for name, plan := range map[string]string{"empty plan": `[]`, "null plan": `null`} {
		t.Run(name, func(t *testing.T) {
			var v0 map[string]any
			if err := json.Unmarshal([]byte(vpsTestV0ImportedState), &v0); err != nil {
				t.Fatal(err)
			}
			v0["plan"] = json.RawMessage(plan)
			raw, err := json.Marshal(v0)
			if err != nil {
				t.Fatal(err)
			}

			state := vpsTestUpgrade(t, string(raw))
			vpsTestCheckValues(t, state.Raw, map[string]string{
				"id":            vpsTestId,
				"plan_code":     "<null>",
				"display_name":  "web-01",
				"power_state":   "<null>",
				"location":      "<null>",
				"installation":  "<null>",
				"options":       "<null>",
				"order":         "<null>",
				"timeouts":      "<null>",
				"current_state": "<null>",
			})
		})
	}
}

// The id of states written before the former resource had one comes from
// service_name.
func TestUpgradeVPSStateV0IdFromServiceName(t *testing.T) {
	var v0 map[string]any
	if err := json.Unmarshal([]byte(vpsTestV0ImportedState), &v0); err != nil {
		t.Fatal(err)
	}
	v0["id"] = nil
	raw, err := json.Marshal(v0)
	if err != nil {
		t.Fatal(err)
	}
	vpsTestCheckValues(t, vpsTestUpgrade(t, string(raw)).Raw, map[string]string{"id": vpsTestId})
}

// The fixtures hold exactly the attributes of the former schema, and the frozen
// copy decodes them without ignoring any attribute.
func TestVPSResourceSchemaV0DecodesOldState(t *testing.T) {
	ctx := context.Background()
	v0 := vpsResourceSchemaV0()
	typ := v0.Type().TerraformType(ctx)
	var names []string
	for name := range v0.Attributes {
		names = append(names, name)
	}
	slices.Sort(names)

	for name, fixture := range map[string]string{"ordered": vpsTestV0OrderedState, "imported": vpsTestV0ImportedState} {
		t.Run(name, func(t *testing.T) {
			var attributes map[string]any
			if err := json.Unmarshal([]byte(fixture), &attributes); err != nil {
				t.Fatal(err)
			}
			var keys []string
			for key := range attributes {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			if !slices.Equal(keys, names) {
				t.Errorf("fixture attributes %v, want the attributes of the former schema %v", keys, names)
			}

			if _, err := (tfprotov6.RawState{JSON: []byte(fixture)}).Unmarshal(typ); err != nil {
				t.Errorf("decoding the fixture with the former schema: %s", err)
			}
		})
	}
}

// vpsTestProposedNew is the proposed new state Terraform sends to the plan of
// the provider: the configuration, with the prior values of the computed
// attributes the configuration leaves null.
func vpsTestProposedNew(t *testing.T, attributes map[string]schema.Attribute, prior, config tftypes.Value) tftypes.Value {
	t.Helper()
	var priorValues, configValues map[string]tftypes.Value
	if !prior.IsNull() {
		if err := prior.As(&priorValues); err != nil {
			t.Fatal(err)
		}
	}
	if err := config.As(&configValues); err != nil {
		t.Fatal(err)
	}
	proposed := make(map[string]tftypes.Value, len(configValues))
	for name, c := range configValues {
		a := attributes[name]
		p, ok := priorValues[name]
		if !ok {
			p = tftypes.NewValue(c.Type(), nil)
		}
		nested, isNested := a.(schema.SingleNestedAttribute)
		switch {
		case a.IsComputed() && !a.IsOptional():
			proposed[name] = p
		case c.IsNull() && a.IsComputed():
			proposed[name] = p
		case c.IsNull() || !isNested:
			proposed[name] = c
		default:
			proposed[name] = vpsTestProposedNew(t, nested.Attributes, p, c)
		}
	}
	return tftypes.NewValue(config.Type(), proposed)
}

// vpsTestPlanThroughServer runs the plan of the provider through its protocol
// server, as Terraform does for an existing resource.
func vpsTestPlanThroughServer(t *testing.T, prior tfsdk.State, config tfsdk.Config) *tfprotov6.PlanResourceChangeResponse {
	t.Helper()
	ctx := context.Background()
	typ := prior.Schema.Type().TerraformType(ctx)
	dynamic := func(v tftypes.Value) *tfprotov6.DynamicValue {
		d, err := tfprotov6.NewDynamicValue(typ, v)
		if err != nil {
			t.Fatal(err)
		}
		return &d
	}
	server, err := providerserver.NewProtocol6WithError(&OvhProvider{})()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "ovh_vps",
		PriorState:       dynamic(prior.Raw),
		ProposedNewState: dynamic(vpsTestProposedNew(t, vpsResourceAttributes(), prior.Raw, config.Raw)),
		Config:           dynamic(config.Raw),
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

// A state migrated from the former ovh_vps, then read, plans neither a
// replacement nor an update with a configuration equivalent to the former one.
func TestVPSResourceMigrationPlansNoChange(t *testing.T) {
	cases := map[string]struct {
		v0     string
		config map[string]attr.Value
		// billingCalls is the number of calls to /services by the Read.
		billingCalls int
	}{
		"ordered": {
			v0: vpsTestV0OrderedState,
			config: map[string]attr.Value{
				"plan_code":                           tfStr("vps-2025-model1"),
				"location.region":                     tfStr("eu-west-gra"),
				"installation.from.image_name":        tfStr("Debian 12"),
				"installation.ssh_key":                tfStr("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMigrated user@laptop"),
				"installation.send_password_by_email": types.BoolValue(false),
			},
		},
		"imported": {
			v0: vpsTestV0ImportedState,
			config: map[string]attr.Value{
				"plan_code":                    tfStr("vps-2025-model3"),
				"location.region":              tfStr("eu-west-gra"),
				"installation.from.image_name": tfStr("Debian 12"),
				"display_name":                 tfStr("web-01"),
			},
			billingCalls: 2,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			upgraded := vpsTestUpgrade(t, c.v0)

			api := newVPSTestAPI()
			vpsTestAddService(api)
			read := vpsTestRead(t, api, upgraded)
			if read.Diagnostics.HasError() || read.Diagnostics.WarningsCount() > 0 {
				t.Fatalf("unexpected read diagnostics: %v", read.Diagnostics)
			}
			if got := vpsTestBillingCalls(api); got != c.billingCalls {
				t.Errorf("Read called /services %d times, want %d", got, c.billingCalls)
			}

			plan := vpsTestPlan(t, c.config)
			config := tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}
			resp := vpsTestPlanThroughServer(t, read.State, config)
			for _, d := range resp.Diagnostics {
				t.Errorf("plan: %s: %s: %s", d.Severity, d.Summary, d.Detail)
			}
			if len(resp.RequiresReplace) > 0 {
				t.Errorf("plan replaces the VPS for %v", resp.RequiresReplace)
			}
			planned, err := resp.PlannedState.Unmarshal(read.State.Raw.Type())
			if err != nil {
				t.Fatal(err)
			}
			if vpsUpdatePlanned(planned, read.State.Raw, config.Raw, false) {
				t.Errorf("the shared rule sees an update in the plan")
			}
			if diffs, err := planned.Diff(read.State.Raw); err != nil {
				t.Fatal(err)
			} else if len(diffs) > 0 {
				for _, d := range diffs {
					t.Errorf("plan changes %s: %v → %v", d.Path, d.Value1, d.Value2)
				}
			}
		})
	}
}
