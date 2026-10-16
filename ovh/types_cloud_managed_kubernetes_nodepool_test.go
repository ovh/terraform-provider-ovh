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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

const testMKSNodepoolResponse = `{
  "id": "np1", "checksum": "abc", "createdAt": "2026-01-01T00:00:00Z", "updatedAt": "2026-01-01T00:00:00Z",
  "resourceStatus": "READY",
  "targetSpec": {
    "name": "pool", "availabilityZone": "eu-west-par-a", "billingType": "HOURLY", "desiredNodes": 3,
    "nodeConfig": {"flavor": "b3-8", "attachFloatingIp": {"enabled": false}},
    "antiAffinity": {"enabled": true, "policy": "HARD"},
    "autoscaler": {"enabled": true, "minNodes": 1, "maxNodes": 5, "scaleDown": {"utilizationThreshold": 0.5}},
    "template": {
      "metadata": {"labels": {"team": "a"}, "finalizers": ["f1"]},
      "spec": {"taints": [{"effect": "NoSchedule", "key": "k", "value": "v"}], "unschedulable": false}
    }
  },
  "currentState": {
    "name": "pool", "desiredNodes": 3, "nodeConfig": {"flavor": "b3-8"},
    "rollout": {"strategy": "ROLLING_UPDATE"}, "status": {"currentNodes": 3, "readyNodes": 2, "upToDateNodes": 3}
  }
}`

func TestCloudManagedKubernetesNodepool_mergeAndUpdate(t *testing.T) {
	ctx := context.Background()

	var rs resource.SchemaResponse
	NewCloudManagedKubernetesNodepoolResource().Schema(ctx, resource.SchemaRequest{}, &rs)
	if d := rs.Schema.ValidateImplementation(ctx); d.HasError() {
		t.Fatalf("invalid resource schema: %v", d)
	}
	var ds datasource.SchemaResponse
	NewCloudManagedKubernetesNodepoolDataSource().(*cloudManagedKubernetesNodepoolDataSource).Schema(ctx, datasource.SchemaRequest{}, &ds)
	if d := ds.Schema.ValidateImplementation(ctx); d.HasError() {
		t.Fatalf("invalid data source schema: %v", d)
	}

	var response CloudManagedKubernetesNodepoolAPIResponse
	if err := json.Unmarshal([]byte(testMKSNodepoolResponse), &response); err != nil {
		t.Fatal(err)
	}

	var diags diag.Diagnostics
	data := CloudManagedKubernetesNodepoolModel{ServiceName: tfStr("p1"), ManagedKubernetesId: tfStr("c1")}
	data.MergeWith(ctx, &response, &diags)
	if diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}

	// The API omits empty annotations; state must keep {} to match a configured empty map.
	if ann := data.Template.Attributes()["metadata"].(types.Object).Attributes()["annotations"]; ann.IsNull() {
		t.Errorf("omitted annotations must merge as an empty map, got null")
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
	var planned CloudManagedKubernetesNodepoolModel
	if d := state.Get(ctx, &planned); d.HasError() {
		t.Fatalf("state get: %v", d)
	}
	payload, err := json.Marshal(planned.ToUpdate(ctx, "abc", &diags))
	if err != nil || diags.HasError() {
		t.Fatalf("update payload: %v %v", err, diags)
	}
	body := string(payload)
	for _, forbidden := range []string{`"flavor"`, `"billingType"`, `"availabilityZone"`, `"antiAffinity"`, `"name"`} {
		if strings.Contains(body, forbidden) {
			t.Errorf("update payload must not carry %s: %s", forbidden, body)
		}
	}
	for _, expected := range []string{`"checksum":"abc"`, `"desiredNodes":3`, `"utilizationThreshold":0.5`, `"labels":{"team":"a"}`, `"taints":[{"effect":"NoSchedule","key":"k","value":"v"}]`} {
		if !strings.Contains(body, expected) {
			t.Errorf("update payload must carry %s: %s", expected, body)
		}
	}

	createBody, err := json.Marshal(planned.ToCreate(ctx, &diags))
	if err != nil || diags.HasError() {
		t.Fatalf("create payload: %v %v", err, diags)
	}
	if !strings.Contains(string(createBody), `"nodeConfig":{"flavor":"b3-8"`) {
		t.Errorf("create payload must carry the node config: %s", createBody)
	}

	// An omitted desired_nodes is unknown at plan time and must not be sent.
	planned.DesiredNodes = ovhtypes.TfInt64Value{Int64Value: types.Int64Unknown()}
	if body, _ := json.Marshal(planned.ToCreate(ctx, &diags)); strings.Contains(string(body), "desiredNodes") {
		t.Errorf("create payload must omit an unset desiredNodes: %s", body)
	}
}

// A merge keeps every API value, including what the config omitted.
func TestCloudManagedKubernetesNodepool_mergeKeepsApiValues(t *testing.T) {
	ctx := context.Background()
	var response CloudManagedKubernetesNodepoolAPIResponse
	if err := json.Unmarshal([]byte(testMKSNodepoolResponse), &response); err != nil {
		t.Fatal(err)
	}

	var diags diag.Diagnostics
	attrTypes := managedKubernetesNodepoolAttrTypes()
	data := CloudManagedKubernetesNodepoolModel{
		ServiceName:         tfStr("p1"),
		ManagedKubernetesId: tfStr("c1"),
		AvailabilityZone:    ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		NodeConfig:          objectFrom(ctx, attrTypes["node_config"], &MKSNodeConfig{Flavor: tfStr("b3-8")}, &diags),
		AntiAffinity:        types.ObjectUnknown(attrTypes["anti_affinity"]),
		Autoscaler:          types.ObjectUnknown(attrTypes["autoscaler"]),
		Template: objectFrom(ctx, attrTypes["template"], &MKSNodeTemplate{Metadata: &MKSNodeTemplateMetadata{
			Labels: map[string]ovhtypes.TfStringValue{"team": tfStr("a")},
		}}, &diags),
	}
	data.MergeWith(ctx, &response, &diags)
	if diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}

	for name, v := range map[string]attr.Value{
		"availability_zone":              data.AvailabilityZone,
		"anti_affinity":                  data.AntiAffinity,
		"autoscaler":                     data.Autoscaler,
		"node_config.attach_floating_ip": data.NodeConfig.Attributes()["attach_floating_ip"],
		"template.spec":                  data.Template.Attributes()["spec"],
		"template.metadata.labels":       data.Template.Attributes()["metadata"].(types.Object).Attributes()["labels"],
	} {
		if v.IsNull() || v.IsUnknown() {
			t.Errorf("%s must come from the API, got %s", name, v)
		}
	}
	if data.CurrentState.IsNull() || data.DesiredNodes.ValueInt64() != 3 {
		t.Error("current_state and computed attributes must come from the API")
	}
}
