package ovh

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

const testMKSResponse = `{
  "id": "c1", "checksum": "abc", "createdAt": "2026-01-01T00:00:00Z", "updatedAt": "2026-01-01T00:00:00Z",
  "resourceStatus": "READY",
  "targetSpec": {
    "name": "my-cluster", "version": "1.34", "region": "GRA11", "plan": "FREE",
    "apiServer": {"access": {"authorizedIps": {"rules": [{"cidr": "0.0.0.0/0"}]}}, "admissionPlugins": {"AlwaysPullImages": {"enabled": true}}},
    "autoUpgrade": {"patch": {"enabled": true}},
    "cni": {"type": "CANAL", "canal": {"kubeProxy": {"mode": "IPVS", "ipvs": {"scheduler": "ROUND_ROBIN", "syncPeriod": "PT60S"}}}},
    "network": {"type": "VRACK", "vrack": {"privateNetworkId": "n1", "subnetId": "s1", "defaultEgress": {"type": "NODE_PUBLIC_INTERFACE"}}}
  },
  "currentState": {
    "name": "my-cluster", "deployedVersion": "1.34.1-2",
    "apiServer": {"access": {"publicEndpoint": {"url": "https://example"}}},
    "authentication": {"type": "OIDC", "oidc": {"clientId": "id", "issuerUrl": "https://issuer", "signingAlgs": ["RS256"]}}
  }
}`

func TestCloudManagedKubernetes_mergeAndUpdate(t *testing.T) {
	ctx := context.Background()

	var rs resource.SchemaResponse
	NewCloudManagedKubernetesResource().Schema(ctx, resource.SchemaRequest{}, &rs)
	if d := rs.Schema.ValidateImplementation(ctx); d.HasError() {
		t.Fatalf("invalid resource schema: %v", d)
	}
	var ds datasource.SchemaResponse
	NewCloudManagedKubernetesDataSource().(*cloudManagedKubernetesDataSource).Schema(ctx, datasource.SchemaRequest{}, &ds)
	if d := ds.Schema.ValidateImplementation(ctx); d.HasError() {
		t.Fatalf("invalid data source schema: %v", d)
	}

	var response CloudManagedKubernetesAPIResponse
	if err := json.Unmarshal([]byte(testMKSResponse), &response); err != nil {
		t.Fatal(err)
	}

	var diags diag.Diagnostics
	data := CloudManagedKubernetesModel{ServiceName: tfStr("p1"), Version: tfStr("1.34")}
	data.MergeWith(ctx, &response, &diags)
	if diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}
	if got := data.Version.ValueString(); got != "1.34" {
		t.Errorf("configured minor version must be kept, got %q", got)
	}

	// The API returns the minor only: a declared full version of that minor is kept, another minor is not.
	full := CloudManagedKubernetesModel{Version: tfStr("1.34.1-2")}
	full.MergeWith(ctx, &response, &diags)
	if got := full.Version.ValueString(); got != "1.34.1-2" {
		t.Errorf("configured full version must be kept, got %q", got)
	}
	other := CloudManagedKubernetesModel{Version: tfStr("1.33.4-1")}
	other.MergeWith(ctx, &response, &diags)
	if got := other.Version.ValueString(); got != "1.34" {
		t.Errorf("a full version of another minor must take the API value, got %q", got)
	}

	// Setting the model in both states checks it matches each schema's types.
	state := tfsdk.State{Schema: rs.Schema, Raw: tftypes.NewValue(rs.Schema.Type().TerraformType(ctx), nil)}
	if d := state.Set(ctx, &data); d.HasError() {
		t.Fatalf("resource state set: %v", d)
	}
	dsState := tfsdk.State{Schema: ds.Schema, Raw: tftypes.NewValue(ds.Schema.Type().TerraformType(ctx), nil)}
	if d := dsState.Set(ctx, &data); d.HasError() {
		t.Fatalf("data source state set: %v", d)
	}

	// Read back from state, as Update does with the plan.
	var planned CloudManagedKubernetesModel
	if d := state.Get(ctx, &planned); d.HasError() {
		t.Fatalf("state get: %v", d)
	}
	payload, err := json.Marshal(planned.ToUpdate(ctx, "abc", &diags))
	if err != nil || diags.HasError() {
		t.Fatalf("update payload: %v %v", err, diags)
	}
	body := string(payload)
	for _, forbidden := range []string{`"type":"CANAL"`, `"mode"`, `"privateNetworkId"`, `"region"`, `"plan"`} {
		if strings.Contains(body, forbidden) {
			t.Errorf("update payload must not carry %s: %s", forbidden, body)
		}
	}
	for _, expected := range []string{`"checksum":"abc"`, `"scheduler":"ROUND_ROBIN"`, `"defaultEgress":{"type":"NODE_PUBLIC_INTERFACE"`, `"version":"1.34"`} {
		if !strings.Contains(body, expected) {
			t.Errorf("update payload must carry %s: %s", expected, body)
		}
	}
}

func TestCloudManagedKubernetes_durationSemanticEquals(t *testing.T) {
	ctx := context.Background()
	d := func(s string) ovhtypes.TfISO8601DurationValue {
		return ovhtypes.TfISO8601DurationValue{TfStringValue: tfStr(s)}
	}
	if eq, _ := d("PT0S").StringSemanticEquals(ctx, d("P0D")); !eq {
		t.Error("PT0S and P0D must be semantically equal")
	}
	if eq, _ := d("PT0S").StringSemanticEquals(ctx, d("PT60S")); eq {
		t.Error("a non-zero duration must not equal zero")
	}
}

// at walks nested object attributes by name.
func at(v attr.Value, path ...string) attr.Value {
	for _, p := range path {
		v = v.(types.Object).Attributes()[p]
	}
	return v
}

// A merge keeps every API value, including what the config omitted, and merging the
// same response twice changes nothing: that is what keeps the next plan empty.
func TestCloudManagedKubernetes_mergeKeepsApiValues(t *testing.T) {
	ctx := context.Background()
	var response CloudManagedKubernetesAPIResponse
	if err := json.Unmarshal([]byte(testMKSResponse), &response); err != nil {
		t.Fatal(err)
	}

	var diags diag.Diagnostics
	attrTypes := managedKubernetesAttrTypes()
	data := CloudManagedKubernetesModel{
		ServiceName: tfStr("p1"), Region: tfStr("GRA11"), Plan: tfStr("FREE"),
		Cni:                objectFrom(ctx, attrTypes["cni"], &MKSCni{Type: tfStr("CANAL"), Canal: &MKSCniCanal{KubeProxy: &MKSKubeProxy{Mode: tfStr("IPVS")}}}, &diags),
		Network:            objectFrom(ctx, attrTypes["network"], &MKSNetwork{Type: tfStr("VRACK"), Vrack: &MKSNetworkVrack{PrivateNetworkId: tfStr("n1"), SubnetId: tfStr("s1")}}, &diags),
		ApiServer:          types.ObjectUnknown(attrTypes["api_server"]),
		Authentication:     types.ObjectUnknown(attrTypes["authentication"]),
		AutoUpgrade:        types.ObjectUnknown(attrTypes["auto_upgrade"]),
		IpAllocationPolicy: types.ObjectUnknown(attrTypes["ip_allocation_policy"]),
	}
	data.MergeWith(ctx, &response, &diags)
	if diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}

	for name, v := range map[string]attr.Value{
		"api_server":                   data.ApiServer,
		"auto_upgrade":                 data.AutoUpgrade,
		"cni.canal.kube_proxy.ipvs":    at(data.Cni, "canal", "kube_proxy", "ipvs"),
		"network.vrack.default_egress": at(data.Network, "vrack", "default_egress"),
	} {
		if v.IsNull() || v.IsUnknown() {
			t.Errorf("%s omitted from the config must come from the API, got %s", name, v)
		}
	}
	for name, v := range map[string]attr.Value{"authentication": data.Authentication, "cni.cilium": at(data.Cni, "cilium")} {
		if !v.IsNull() {
			t.Errorf("%s absent from the API must be null, got %s", name, v)
		}
	}
	if data.Name.ValueString() != "my-cluster" || data.CurrentState.IsNull() {
		t.Error("computed root attributes and current_state must come from the API")
	}

	first := data
	data.MergeWith(ctx, &response, &diags)
	for name, pair := range map[string][2]attr.Value{
		"cni": {first.Cni, data.Cni}, "network": {first.Network, data.Network}, "api_server": {first.ApiServer, data.ApiServer},
	} {
		if !pair[0].Equal(pair[1]) {
			t.Errorf("a second merge must not change %s: %s -> %s", name, pair[0], pair[1])
		}
	}
}

// Every optional attribute is also computed, so a config that omits it keeps the state
// value: this is what makes an import plan empty. Below a computed attribute nothing may be
// Required: Terraform core would then read a null config as a removal and plan a change.
func assertOptionalIsComputed(t *testing.T, prefix string, attrs map[string]schema.Attribute, belowComputed bool) {
	t.Helper()
	for name, a := range attrs {
		path := prefix + name
		if a.IsOptional() && !a.IsComputed() {
			t.Errorf("%s is Optional but not Computed", path)
		}
		if belowComputed && !a.IsComputed() {
			t.Errorf("%s is Required below a Computed attribute", path)
		}
		switch v := a.(type) {
		case schema.SingleNestedAttribute:
			assertOptionalIsComputed(t, path+".", v.Attributes, v.IsComputed())
		case schema.ListNestedAttribute:
			assertOptionalIsComputed(t, path+".", v.NestedObject.Attributes, v.IsComputed())
		case schema.MapNestedAttribute:
			assertOptionalIsComputed(t, path+".", v.NestedObject.Attributes, v.IsComputed())
		}
	}
}

func TestCloudManagedKubernetes_optionalIsComputed(t *testing.T) {
	assertOptionalIsComputed(t, "", managedKubernetesResourceAttributes(), false)
	assertOptionalIsComputed(t, "", managedKubernetesNodepoolResourceAttributes(), false)
}

func TestCloudManagedKubernetes_declaredAdmissionPlugins(t *testing.T) {
	ctx := context.Background()
	response := func() *CloudManagedKubernetesAPIResponse {
		return &CloudManagedKubernetesAPIResponse{TargetSpec: &CloudManagedKubernetesAPITargetSpec{ApiServer: &MKSApiServer{
			AdmissionPlugins: map[string]MKSAdmissionPlugin{
				"AlwaysPullImages": {Enabled: ovhtypes.NewTfBoolValue(true)},
				"NodeRestriction":  {Enabled: ovhtypes.NewTfBoolValue(false)},
			},
		}}}
	}
	plugins := func(m CloudManagedKubernetesModel) map[string]MKSAdmissionPlugin {
		var diags diag.Diagnostics
		return objectAs[MKSApiServer](ctx, m.ApiServer, &diags).AdmissionPlugins
	}

	var diags diag.Diagnostics
	// No prior value (import, data source): every plugin is read.
	var imported CloudManagedKubernetesModel
	imported.MergeWith(ctx, response(), &diags)
	if got := plugins(imported); len(got) != 2 {
		t.Errorf("without prior value every plugin must be read, got %v", got)
	}

	// Declared plugins only: the API's extra plugin is left to current_state.
	declared := CloudManagedKubernetesModel{ApiServer: objectFrom(ctx, managedKubernetesAttrTypes()["api_server"], &MKSApiServer{
		AdmissionPlugins: map[string]MKSAdmissionPlugin{"NodeRestriction": {Enabled: ovhtypes.NewTfBoolValue(true)}},
	}, &diags)}
	declared.MergeWith(ctx, response(), &diags)
	if diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}
	got := plugins(declared)
	if _, extra := got["AlwaysPullImages"]; extra || len(got) != 1 {
		t.Errorf("only declared plugins must be kept, got %v", got)
	}
	if got["NodeRestriction"].Enabled.ValueBool() {
		t.Errorf("the API value of a declared plugin must win to surface drift")
	}
}

func TestCloudManagedKubernetes_dataSourceSchemas(t *testing.T) {
	ctx := context.Background()
	for _, newDataSource := range []func() datasource.DataSource{
		NewCloudManagedKubernetesKubeconfigDataSource,
		NewCloudManagedKubernetesPlansDataSource,
		NewCloudManagedKubernetesRegionsDataSource,
		NewCloudManagedKubernetesVersionsDataSource,
	} {
		var resp datasource.SchemaResponse
		newDataSource().Schema(ctx, datasource.SchemaRequest{}, &resp)
		if d := resp.Schema.ValidateImplementation(ctx); d.HasError() {
			t.Errorf("invalid data source schema: %v", d)
		}
	}
}

func TestVariantValidator(t *testing.T) {
	ctx := context.Background()
	block := types.ObjectType{AttrTypes: map[string]attr.Type{}}
	attrTypes := map[string]attr.Type{"type": types.StringType, "canal": block, "cilium": block}
	cni := func(typ string, canal, cilium bool) types.Object {
		set := func(b bool) attr.Value {
			if b {
				return types.ObjectValueMust(block.AttrTypes, map[string]attr.Value{})
			}
			return types.ObjectNull(block.AttrTypes)
		}
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{"type": types.StringValue(typ), "canal": set(canal), "cilium": set(cilium)})
	}
	v := variantValidator{"type", []string{"canal", "cilium"}}
	for name, tc := range map[string]struct {
		value   types.Object
		wantErr bool
	}{
		"matching block":   {cni("CANAL", true, false), false},
		"no block":         {cni("CILIUM", false, false), false},
		"mismatched block": {cni("CANAL", false, true), true},
	} {
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{ConfigValue: tc.value}, &resp)
		if resp.Diagnostics.HasError() != tc.wantErr {
			t.Errorf("%s: got errors %v", name, resp.Diagnostics)
		}
	}
}
