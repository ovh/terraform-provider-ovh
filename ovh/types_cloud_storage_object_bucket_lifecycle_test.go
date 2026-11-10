package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func storageObjectBucketLifecycleModelFrom(t *testing.T, rules []storageObjectBucketLifecycleRuleModel) CloudStorageObjectBucketLifecycleModel {
	t.Helper()

	list, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: storageObjectBucketLifecycleRuleAttrTypes()}, rules)
	if diags.HasError() {
		t.Fatalf("build rules: %v", diags)
	}
	return CloudStorageObjectBucketLifecycleModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_bucket"),
		Rules:       list,
	}
}

func storageObjectBucketLifecycleRulesOf(t *testing.T, m CloudStorageObjectBucketLifecycleModel) []storageObjectBucketLifecycleRuleModel {
	t.Helper()

	rules, diags := m.rules(context.Background())
	if diags.HasError() {
		t.Fatalf("read rules: %v", diags)
	}
	return rules
}

func TestCloudStorageObjectBucketLifecycleToPut(t *testing.T) {
	model := storageObjectBucketLifecycleModelFrom(t, []storageObjectBucketLifecycleRuleModel{
		{
			Id:     ovhtypes.NewTfStringValue("expire-logs"),
			Status: ovhtypes.NewTfStringValue("ENABLED"),
			Filter: &storageObjectBucketLifecycleFilterModel{
				Prefix:                ovhtypes.NewTfStringValue("logs/"),
				Tags:                  types.MapNull(ovhtypes.TfStringType{}),
				ObjectSizeGreaterThan: types.Int64Null(),
				ObjectSizeLessThan:    types.Int64Null(),
			},
			Expiration: &storageObjectBucketLifecycleExpirationModel{
				Date:                      ovhtypes.TfStringValue{StringValue: types.StringNull()},
				Days:                      types.Int64Value(30),
				ExpiredObjectDeleteMarker: types.BoolNull(),
			},
			Transitions: []storageObjectBucketLifecycleTransitionModel{
				{
					Date:         ovhtypes.TfStringValue{StringValue: types.StringNull()},
					Days:         types.Int64Value(7),
					StorageClass: ovhtypes.NewTfStringValue("STANDARD_IA"),
				},
			},
		},
	})

	payload, diags := model.ToPut(context.Background())
	if diags.HasError() {
		t.Fatalf("ToPut: %v", diags)
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := `{"rules":[{"id":"expire-logs","status":"ENABLED","filter":{"prefix":"logs/"},"expiration":{"days":30},"transitions":[{"days":7,"storageClass":"STANDARD_IA"}]}]}`
	if string(raw) != want {
		t.Fatalf("payload mismatch\n got: %s\nwant: %s", raw, want)
	}
}

func TestCloudStorageObjectBucketLifecycleToPut_EmptyRulesSentAsArray(t *testing.T) {
	model := CloudStorageObjectBucketLifecycleModel{Rules: types.ListNull(types.ObjectType{AttrTypes: storageObjectBucketLifecycleRuleAttrTypes()})}

	payload, diags := model.ToPut(context.Background())
	if diags.HasError() {
		t.Fatalf("ToPut: %v", diags)
	}
	raw, _ := json.Marshal(payload)
	if string(raw) != `{"rules":[]}` {
		t.Fatalf("got %s, want {\"rules\":[]}", raw)
	}
}

func TestCloudStorageObjectBucketLifecycleMergeWith_KeepsPriorOrder(t *testing.T) {
	model := storageObjectBucketLifecycleModelFrom(t, []storageObjectBucketLifecycleRuleModel{
		{Id: ovhtypes.NewTfStringValue("zeta"), Status: ovhtypes.NewTfStringValue("ENABLED")},
		{Id: ovhtypes.NewTfStringValue("alpha"), Status: ovhtypes.NewTfStringValue("ENABLED")},
	})

	response := CloudStorageObjectBucketLifecycleAPIRules{Rules: []CloudStorageObjectBucketLifecycleAPIRule{
		{Id: "alpha", Status: "ENABLED", Filter: &CloudStorageObjectBucketLifecycleAPIFilter{}},
		{Id: "new", Status: "DISABLED"},
		{Id: "zeta", Status: "ENABLED"},
	}}

	if diags := model.MergeWith(context.Background(), &response); diags.HasError() {
		t.Fatalf("MergeWith: %v", diags)
	}

	rules := storageObjectBucketLifecycleRulesOf(t, model)
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
	if model.Id.ValueString() != "GRA_bucket" {
		t.Fatalf("id = %q, want the bucket id", model.Id.ValueString())
	}
}

func TestCloudStorageObjectBucketLifecycleMergeWith_KeepsSpelledOutZeroValues(t *testing.T) {
	model := storageObjectBucketLifecycleModelFrom(t, []storageObjectBucketLifecycleRuleModel{
		{
			Id:     ovhtypes.NewTfStringValue("all"),
			Status: ovhtypes.NewTfStringValue("ENABLED"),
			Filter: &storageObjectBucketLifecycleFilterModel{
				Prefix:                ovhtypes.NewTfStringValue(""),
				Tags:                  types.MapNull(ovhtypes.TfStringType{}),
				ObjectSizeGreaterThan: types.Int64Null(),
				ObjectSizeLessThan:    types.Int64Null(),
			},
			Expiration: &storageObjectBucketLifecycleExpirationModel{
				Date:                      ovhtypes.TfStringValue{StringValue: types.StringNull()},
				Days:                      types.Int64Value(10),
				ExpiredObjectDeleteMarker: types.BoolValue(false),
			},
			Transitions: []storageObjectBucketLifecycleTransitionModel{
				{Date: ovhtypes.TfStringValue{StringValue: types.StringNull()}, Days: types.Int64Value(30), StorageClass: ovhtypes.NewTfStringValue("STANDARD_IA")},
				{Date: ovhtypes.TfStringValue{StringValue: types.StringNull()}, Days: types.Int64Value(90), StorageClass: ovhtypes.NewTfStringValue("DEEP_ARCHIVE")},
			},
		},
	})

	response := CloudStorageObjectBucketLifecycleAPIRules{Rules: []CloudStorageObjectBucketLifecycleAPIRule{
		{
			Id:         "all",
			Status:     "ENABLED",
			Filter:     &CloudStorageObjectBucketLifecycleAPIFilter{},
			Expiration: &CloudStorageObjectBucketLifecycleAPIExpiration{Days: 10},
			Transitions: []CloudStorageObjectBucketLifecycleAPITransition{
				{Days: 90, StorageClass: "DEEP_ARCHIVE"},
				{Days: 30, StorageClass: "STANDARD_IA"},
			},
		},
	}}

	if diags := model.MergeWith(context.Background(), &response); diags.HasError() {
		t.Fatalf("MergeWith: %v", diags)
	}

	rule := storageObjectBucketLifecycleRulesOf(t, model)[0]
	if rule.Filter == nil || rule.Filter.Prefix.IsNull() || rule.Filter.Prefix.ValueString() != "" {
		t.Fatalf("filter.prefix = \"\" must be kept, got %+v", rule.Filter)
	}
	if rule.Expiration.ExpiredObjectDeleteMarker.IsNull() || rule.Expiration.ExpiredObjectDeleteMarker.ValueBool() {
		t.Fatalf("expired_object_delete_marker = false must be kept, got %v", rule.Expiration.ExpiredObjectDeleteMarker)
	}
	if rule.Transitions[0].StorageClass.ValueString() != "STANDARD_IA" || rule.Transitions[1].StorageClass.ValueString() != "DEEP_ARCHIVE" {
		t.Fatalf("transitions must keep the prior order, got %+v", rule.Transitions)
	}
}
