package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// helpers

func tfStringMap(ctx context.Context, kv map[string]string) ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] {
	elems := make(map[string]attr.Value, len(kv))
	for k, v := range kv {
		elems[k] = ovhtypes.NewTfStringValue(v)
	}
	m, _ := ovhtypes.NewTfMapNestedValue[ovhtypes.TfStringValue](ctx, elems)
	return m
}

func nullTagMap(ctx context.Context) ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] {
	return ovhtypes.NewNullTfMapNestedValue[ovhtypes.TfStringValue](ctx)
}

func unknownTagMap(ctx context.Context) ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] {
	return ovhtypes.NewUnknownTfMapNestedValue[ovhtypes.TfStringValue](ctx)
}

// ToCreate

func TestCloudProjectRegionStorageModel_ToCreate_TagsUnknown(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: unknownTagMap(ctx),
	}

	result := m.ToCreate()

	if result.Tags != nil {
		t.Errorf("expected Tags to be nil in writable model when source is unknown, got %v", result.Tags)
	}
}

func TestCloudProjectRegionStorageModel_ToCreate_TagsNull(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: nullTagMap(ctx),
	}

	result := m.ToCreate()

	if result.Tags == nil {
		t.Fatal("expected Tags pointer to be non-nil when source is null (signals clearing)")
	}
	if !result.Tags.IsNull() {
		t.Errorf("expected Tags value to be null, got %v", result.Tags)
	}
}

func TestCloudProjectRegionStorageModel_ToCreate_TagsWithValues(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"env": "prod", "team": "platform"}),
	}

	result := m.ToCreate()

	if result.Tags == nil {
		t.Fatal("expected Tags pointer to be non-nil")
	}
	elems := result.Tags.Elements()
	if len(elems) != 2 {
		t.Fatalf("expected 2 tag elements, got %d", len(elems))
	}
	if v, ok := elems["env"]; !ok || v.(ovhtypes.TfStringValue).ValueString() != "prod" {
		t.Errorf("unexpected value for tag 'env': %v", v)
	}
}

// ToUpdate

func TestCloudProjectRegionStorageModel_ToUpdate_TagsUnknown(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: unknownTagMap(ctx),
	}

	result := m.ToUpdate()

	if result.Tags != nil {
		t.Errorf("expected Tags to be nil in writable model when source is unknown, got %v", result.Tags)
	}
}

func TestCloudProjectRegionStorageModel_ToUpdate_TagsNull(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: nullTagMap(ctx),
	}

	result := m.ToUpdate()

	if result.Tags == nil {
		t.Fatal("expected Tags pointer to be non-nil when source is null (signals clearing)")
	}
	if !result.Tags.IsNull() {
		t.Errorf("expected Tags value to be null, got %v", result.Tags)
	}
}

func TestCloudProjectRegionStorageModel_ToUpdate_TagsWithValues(t *testing.T) {
	ctx := context.Background()
	m := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"env": "staging"}),
	}

	result := m.ToUpdate()

	if result.Tags == nil {
		t.Fatal("expected Tags pointer to be non-nil")
	}
	elems := result.Tags.Elements()
	if len(elems) != 1 {
		t.Fatalf("expected 1 tag element, got %d", len(elems))
	}
}

// MergeWith

func TestCloudProjectRegionStorageModel_MergeWith_TagsMergedFromOther(t *testing.T) {
	ctx := context.Background()
	target := CloudProjectRegionStorageModel{
		Tags: unknownTagMap(ctx),
	}
	source := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"env": "prod"}),
	}

	target.MergeWith(&source)

	if target.Tags.IsUnknown() || target.Tags.IsNull() {
		t.Errorf("expected Tags to be populated after MergeWith, got %v", target.Tags)
	}
	elems := target.Tags.Elements()
	if _, ok := elems["env"]; !ok {
		t.Errorf("expected 'env' key in merged Tags")
	}
}

func TestCloudProjectRegionStorageModel_MergeWith_TagsNotOverwrittenWhenKnown(t *testing.T) {
	ctx := context.Background()
	target := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"env": "prod"}),
	}
	source := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"env": "staging", "other": "val"}),
	}

	target.MergeWith(&source)

	elems := target.Tags.Elements()
	if v, ok := elems["env"]; !ok || v.(ovhtypes.TfStringValue).ValueString() != "prod" {
		t.Errorf("expected existing Tags to not be overwritten by MergeWith, got %v", v)
	}
	if _, ok := elems["other"]; ok {
		t.Errorf("expected source-only key 'other' to not appear after MergeWith")
	}
}

func TestCloudProjectRegionStorageModel_MergeWith_NullTargetGetsSourceValue(t *testing.T) {
	ctx := context.Background()
	target := CloudProjectRegionStorageModel{
		Tags: nullTagMap(ctx),
	}
	source := CloudProjectRegionStorageModel{
		Tags: tfStringMap(ctx, map[string]string{"k": "v"}),
	}

	target.MergeWith(&source)

	if target.Tags.IsNull() {
		t.Errorf("expected null Tags to be replaced by source value after MergeWith")
	}
}

// JSON marshaling of writable model

func TestCloudProjectRegionStorageWritableModel_JSON_TagsOmittedWhenNilPointer(t *testing.T) {
	m := CloudProjectRegionStorageWritableModel{
		Tags: nil,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if _, found := parsed["tags"]; found {
		t.Errorf("expected 'tags' to be absent from JSON when pointer is nil, got: %s", data)
	}
}

func TestCloudProjectRegionStorageWritableModel_JSON_TagsNullWhenPointerToNullValue(t *testing.T) {
	ctx := context.Background()
	nullTags := nullTagMap(ctx)
	m := CloudProjectRegionStorageWritableModel{
		Tags: &nullTags,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	raw, found := parsed["tags"]
	if !found {
		t.Fatalf("expected 'tags' to be present in JSON, got: %s", data)
	}
	if string(raw) != "null" {
		t.Errorf("expected 'tags' to be null in JSON, got: %s", raw)
	}
}

func TestCloudProjectRegionStorageWritableModel_JSON_TagsWithValues(t *testing.T) {
	ctx := context.Background()
	tags := tfStringMap(ctx, map[string]string{"env": "prod"})
	m := CloudProjectRegionStorageWritableModel{
		Tags: &tags,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	raw, found := parsed["tags"]
	if !found {
		t.Fatalf("expected 'tags' to be present in JSON, got: %s", data)
	}

	var tagMap map[string]string
	if err := json.Unmarshal(raw, &tagMap); err != nil {
		t.Fatalf("unexpected unmarshal of tags: %v", err)
	}
	if tagMap["env"] != "prod" {
		t.Errorf("expected tags[env]=prod, got: %v", tagMap)
	}
}

// JSON unmarshaling of API response into model

func TestCloudProjectRegionStorageModel_UnmarshalJSON_TagsFromAPIResponse(t *testing.T) {
	apiResponse := `{
		"name": "my-bucket",
		"region": "GRA",
		"tags": {"environment": "production", "team": "platform"},
		"objectsCount": 0,
		"objectsSize": 0
	}`

	var m CloudProjectRegionStorageModel
	if err := json.Unmarshal([]byte(apiResponse), &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if m.Tags.IsNull() || m.Tags.IsUnknown() {
		t.Fatalf("expected Tags to be known and non-null, got null=%v unknown=%v", m.Tags.IsNull(), m.Tags.IsUnknown())
	}

	elems := m.Tags.Elements()
	if len(elems) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(elems))
	}
	if v, ok := elems["environment"]; !ok || v.(ovhtypes.TfStringValue).ValueString() != "production" {
		t.Errorf("unexpected value for tag 'environment': %v", v)
	}
	if v, ok := elems["team"]; !ok || v.(ovhtypes.TfStringValue).ValueString() != "platform" {
		t.Errorf("unexpected value for tag 'team': %v", v)
	}
}

func TestCloudProjectRegionStorageModel_UnmarshalJSON_TagsNullFromAPIResponse(t *testing.T) {
	apiResponse := `{
		"name": "my-bucket",
		"region": "GRA",
		"tags": null,
		"objectsCount": 0,
		"objectsSize": 0
	}`

	var m CloudProjectRegionStorageModel
	if err := json.Unmarshal([]byte(apiResponse), &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if !m.Tags.IsNull() {
		t.Errorf("expected Tags to be null when API returns null, got %v", m.Tags)
	}
}

func TestCloudProjectRegionStorageModel_UnmarshalJSON_TagsAbsentFromAPIResponse(t *testing.T) {
	apiResponse := `{
		"name": "my-bucket",
		"region": "GRA",
		"objectsCount": 0,
		"objectsSize": 0
	}`

	var m CloudProjectRegionStorageModel
	if err := json.Unmarshal([]byte(apiResponse), &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	// When the field is absent from JSON, TfMapNestedValue stays at zero value (null).
	if !m.Tags.IsNull() {
		t.Errorf("expected Tags to be null when absent from API response, got %v", m.Tags)
	}
}
