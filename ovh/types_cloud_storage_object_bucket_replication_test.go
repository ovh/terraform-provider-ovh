package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func storageObjectBucketReplicationModelFrom(t *testing.T, rules []storageObjectBucketReplicationRuleModel) CloudStorageObjectBucketReplicationModel {
	t.Helper()

	list, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: storageObjectBucketReplicationRuleAttrTypes()}, rules)
	if diags.HasError() {
		t.Fatalf("build rules: %v", diags)
	}
	return CloudStorageObjectBucketReplicationModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
		Rules:       list,
	}
}

func storageObjectBucketReplicationRulesOf(t *testing.T, m CloudStorageObjectBucketReplicationModel) []storageObjectBucketReplicationRuleModel {
	t.Helper()

	rules, diags := m.rules(context.Background())
	if diags.HasError() {
		t.Fatalf("read rules: %v", diags)
	}
	return rules
}

func storageObjectBucketReplicationRule(id, region string, filter *storageObjectBucketReplicationFilterModel) storageObjectBucketReplicationRuleModel {
	return storageObjectBucketReplicationRuleModel{
		Id:                      ovhtypes.NewTfStringValue(id),
		Status:                  ovhtypes.NewTfStringValue("ENABLED"),
		Priority:                types.Int64Value(1),
		DeleteMarkerReplication: ovhtypes.NewTfStringValue("DISABLED"),
		Destination: storageObjectBucketReplicationDestinationModel{
			Name:         ovhtypes.NewTfStringValue("replica"),
			Region:       ovhtypes.NewTfStringValue(region),
			StorageClass: ovhtypes.TfStringValue{StringValue: types.StringNull()},
		},
		Filter: filter,
	}
}

func TestCloudStorageObjectBucketReplicationToPut(t *testing.T) {
	tags, _ := types.MapValueFrom(context.Background(), ovhtypes.TfStringType{}, map[string]ovhtypes.TfStringValue{"env": ovhtypes.NewTfStringValue("prod")})
	withFilter := storageObjectBucketReplicationRule("logs", "GRA", &storageObjectBucketReplicationFilterModel{
		Prefix: ovhtypes.NewTfStringValue("logs/"),
		Tags:   tags,
	})
	withFilter.Destination.StorageClass = ovhtypes.NewTfStringValue("STANDARD_IA")
	model := storageObjectBucketReplicationModelFrom(t, []storageObjectBucketReplicationRuleModel{
		withFilter,
		storageObjectBucketReplicationRule("all", "SBG", nil),
	})

	payload, diags := model.ToPut(context.Background())
	if diags.HasError() {
		t.Fatalf("ToPut: %v", diags)
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := `{"rules":[` +
		`{"id":"logs","status":"ENABLED","priority":1,"deleteMarkerReplication":"DISABLED","destination":{"name":"replica","region":"GRA","storageClass":"STANDARD_IA"},"filter":{"prefix":"logs/","tags":{"env":"prod"}}},` +
		`{"id":"all","status":"ENABLED","priority":1,"deleteMarkerReplication":"DISABLED","destination":{"name":"replica","region":"SBG"}}` +
		`]}`
	if string(raw) != want {
		t.Fatalf("payload mismatch\n got: %s\nwant: %s", raw, want)
	}
}

func TestCloudStorageObjectBucketReplicationToPut_EmptyRulesSentAsArray(t *testing.T) {
	model := CloudStorageObjectBucketReplicationModel{Rules: types.ListNull(types.ObjectType{AttrTypes: storageObjectBucketReplicationRuleAttrTypes()})}

	payload, diags := model.ToPut(context.Background())
	if diags.HasError() {
		t.Fatalf("ToPut: %v", diags)
	}
	raw, _ := json.Marshal(payload)
	if string(raw) != `{"rules":[]}` {
		t.Fatalf("got %s, want {\"rules\":[]}", raw)
	}
}

func TestCloudStorageObjectBucketReplicationMergeWith_KeepsPriorOrder(t *testing.T) {
	model := storageObjectBucketReplicationModelFrom(t, []storageObjectBucketReplicationRuleModel{
		storageObjectBucketReplicationRule("zeta", "GRA", nil),
		storageObjectBucketReplicationRule("alpha", "GRA", nil),
	})

	response := CloudStorageObjectBucketReplicationAPIRules{Rules: []CloudStorageObjectBucketReplicationAPIRule{
		{Id: "alpha", Status: "ENABLED", DeleteMarkerReplication: "DISABLED", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", Region: "GRA"}, Filter: &CloudStorageObjectBucketReplicationAPIFilter{}},
		{Id: "new", Status: "DISABLED", DeleteMarkerReplication: "ENABLED", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", Region: "GRA"}},
		{Id: "zeta", Status: "ENABLED", DeleteMarkerReplication: "DISABLED", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", Region: "GRA"}},
	}}

	if diags := model.MergeWith(context.Background(), &response); diags.HasError() {
		t.Fatalf("MergeWith: %v", diags)
	}

	rules := storageObjectBucketReplicationRulesOf(t, model)
	got := []string{}
	for _, r := range rules {
		got = append(got, r.Id.ValueString())
	}
	if len(got) != 3 || got[0] != "zeta" || got[1] != "alpha" || got[2] != "new" {
		t.Fatalf("rule order = %v, want [zeta alpha new]", got)
	}
	if rules[1].Filter != nil {
		t.Fatalf("an empty API filter must stay null when the prior had none, got %+v", rules[1].Filter)
	}
	if model.Id.ValueString() != "GRA_source" {
		t.Fatalf("id = %q, want the bucket id", model.Id.ValueString())
	}
}

func TestCloudStorageObjectBucketReplicationMergeWith_EmptyValues(t *testing.T) {
	model := storageObjectBucketReplicationModelFrom(t, []storageObjectBucketReplicationRuleModel{
		storageObjectBucketReplicationRule("all", "GRA", &storageObjectBucketReplicationFilterModel{
			Prefix: ovhtypes.NewTfStringValue(""),
			Tags:   types.MapNull(ovhtypes.TfStringType{}),
		}),
	})

	response := CloudStorageObjectBucketReplicationAPIRules{Rules: []CloudStorageObjectBucketReplicationAPIRule{
		{Id: "all", Status: "ENABLED", Priority: 1, DeleteMarkerReplication: "DISABLED", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica"}, Filter: &CloudStorageObjectBucketReplicationAPIFilter{}},
	}}

	if diags := model.MergeWith(context.Background(), &response); diags.HasError() {
		t.Fatalf("MergeWith: %v", diags)
	}

	rule := storageObjectBucketReplicationRulesOf(t, model)[0]
	if rule.Filter == nil || rule.Filter.Prefix.IsNull() || rule.Filter.Prefix.ValueString() != "" {
		t.Fatalf("filter.prefix = \"\" must be kept, got %+v", rule.Filter)
	}
	if rule.Destination.Region.ValueString() != "GRA" {
		t.Fatalf("an empty API destination region must keep the prior one, got %v", rule.Destination.Region)
	}
	if !rule.Destination.StorageClass.IsNull() {
		t.Fatalf("an empty API storage class must be null, got %v", rule.Destination.StorageClass)
	}
}

func TestCloudStorageObjectBucketReplicationMergeWith_Import(t *testing.T) {
	model := CloudStorageObjectBucketReplicationModel{
		BucketId: ovhtypes.NewTfStringValue("GRA_source"),
		Rules:    types.ListNull(types.ObjectType{AttrTypes: storageObjectBucketReplicationRuleAttrTypes()}),
	}

	response := CloudStorageObjectBucketReplicationAPIRules{Rules: []CloudStorageObjectBucketReplicationAPIRule{
		{Id: "all", Status: "ENABLED", DeleteMarkerReplication: "DISABLED", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", StorageClass: "STANDARD"}, Filter: &CloudStorageObjectBucketReplicationAPIFilter{Tags: map[string]string{"env": "prod"}}},
	}}

	if diags := model.MergeWith(context.Background(), &response); diags.HasError() {
		t.Fatalf("MergeWith: %v", diags)
	}

	rule := storageObjectBucketReplicationRulesOf(t, model)[0]
	if !rule.Destination.Region.IsNull() {
		t.Fatalf("an empty API destination region without prior must be null, got %v", rule.Destination.Region)
	}
	if rule.Destination.StorageClass.ValueString() != "STANDARD" || rule.Priority.ValueInt64() != 0 {
		t.Fatalf("unexpected destination or priority: %+v", rule)
	}
	if rule.Filter == nil || !rule.Filter.Prefix.IsNull() || len(rule.Filter.Tags.Elements()) != 1 {
		t.Fatalf("filter must carry the tag and a null prefix, got %+v", rule.Filter)
	}
}

func TestCloudStorageObjectBucketReplicationResourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()

	var schemaResp resource.SchemaResponse
	NewCloudStorageObjectBucketReplicationResource().Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	model := storageObjectBucketReplicationModelFrom(t, []storageObjectBucketReplicationRuleModel{
		storageObjectBucketReplicationRule("all", "GRA", &storageObjectBucketReplicationFilterModel{
			Prefix: ovhtypes.NewTfStringValue("logs/"),
			Tags:   types.MapNull(ovhtypes.TfStringType{}),
		}),
	})
	model.Id = model.BucketId

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
}

func TestCloudStorageObjectBucketReplicationJobResourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()

	var schemaResp resource.SchemaResponse
	NewCloudStorageObjectBucketReplicationJobResource().Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	model := CloudStorageObjectBucketReplicationJobModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
	}
	model.MergeWith(&CloudStorageObjectBucketReplicationJobAPIResponse{
		Id:        "3f0c9c1e-6d4b-4c55-9f0e-0d3a1d0b8f11",
		Status:    "ACTIVE",
		Total:     10,
		Completed: 4,
		Failed:    1,
		Stuck:     0,
		CreatedAt: "2026-11-11T10:00:00Z",
	})

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
}

func TestBuildStorageObjectBucketReplicationRule_Region(t *testing.T) {
	cases := []struct {
		name      string
		prior     *storageObjectBucketReplicationRuleModel
		apiRegion string
		want      string
	}{
		{"case-only difference keeps prior", ptrStorageObjectBucketReplicationRule("all", "sbg"), "SBG", "sbg"},
		{"real drift takes API", ptrStorageObjectBucketReplicationRule("all", "GRA"), "SBG", "SBG"},
		{"empty API keeps prior", ptrStorageObjectBucketReplicationRule("all", "sbg"), "", "sbg"},
		{"no prior takes API", nil, "SBG", "SBG"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			api := CloudStorageObjectBucketReplicationAPIRule{Id: "all", Destination: CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", Region: tc.apiRegion}}
			got := buildStorageObjectBucketReplicationRule(api, tc.prior).Destination.Region
			if got.IsNull() || got.ValueString() != tc.want {
				t.Fatalf("region = %v, want %q", got, tc.want)
			}
		})
	}
}

func ptrStorageObjectBucketReplicationRule(id, region string) *storageObjectBucketReplicationRuleModel {
	r := storageObjectBucketReplicationRule(id, region, nil)
	return &r
}
