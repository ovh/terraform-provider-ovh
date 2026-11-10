package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// runStorageClassValidators runs every validator declared on a StringAttribute
// against the given value and returns whether any error diagnostic was raised.
func runStorageClassValidators(t *testing.T, attr schema.StringAttribute, value string) bool {
	t.Helper()

	req := validator.StringRequest{
		ConfigValue: types.StringValue(value),
	}

	for _, v := range attr.Validators {
		resp := &validator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			return true
		}
	}

	return false
}

func mustStringAttribute(t *testing.T, attr schema.Attribute, name string) schema.StringAttribute {
	t.Helper()
	sa, ok := attr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a StringAttribute (got %T)", name, attr)
	}
	return sa
}

func mustListNested(t *testing.T, attr schema.Attribute, name string) schema.ListNestedAttribute {
	t.Helper()
	la, ok := attr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a ListNestedAttribute (got %T)", name, attr)
	}
	return la
}

func mustSingleNested(t *testing.T, attr schema.Attribute, name string) schema.SingleNestedAttribute {
	t.Helper()
	sa, ok := attr.(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a SingleNestedAttribute (got %T)", name, attr)
	}
	return sa
}

// assertStorageClassLifecycleEnum asserts that a lifecycle storage_class validator
// accepts the full StorageClassLifecycleEnum (regression test for issue #1385)
// and rejects HIGH_PERF, which the API refuses as a transition target.
func assertStorageClassLifecycleEnum(t *testing.T, storageClass schema.StringAttribute, label string) {
	t.Helper()

	for _, value := range []string{"DEEP_ARCHIVE", "GLACIER_IR", "STANDARD", "STANDARD_IA"} {
		if runStorageClassValidators(t, storageClass, value) {
			t.Errorf("%s storage_class rejected %q but it should be accepted", label, value)
		}
	}

	for _, value := range []string{"HIGH_PERF", "NOPE"} {
		if !runStorageClassValidators(t, storageClass, value) {
			t.Errorf("%s storage_class accepted %q but it should be rejected", label, value)
		}
	}
}

func TestStorageClassLifecycleTransitionValidator(t *testing.T) {
	resources := map[string]resource.Resource{
		"ovh_cloud_project_storage_object_bucket_lifecycle_configuration": NewCloudProjectStorageLifecycleConfigurationResource(),
		"ovh_cloud_storage_object_bucket_lifecycle":                       NewCloudStorageObjectBucketLifecycleResource(),
	}

	for name, r := range resources {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		rules := mustListNested(t, resp.Schema.Attributes["rules"], "rules")
		for _, list := range []string{"transitions", "noncurrent_version_transitions"} {
			transitions := mustListNested(t, rules.NestedObject.Attributes[list], list)
			storageClass := mustStringAttribute(t, transitions.NestedObject.Attributes["storage_class"], "storage_class")
			assertStorageClassLifecycleEnum(t, storageClass, name+" "+list)
		}
	}
}

// TestStorageClassReplicationDestinationValidator asserts that the replication
// destination storage_class validator accepts the full StorageClassReplicationEnum
// (related drift from issue #1385) and still rejects garbage.
func TestStorageClassReplicationDestinationValidator(t *testing.T) {
	ctx := context.Background()

	r := NewCloudProjectStorageResource()
	resp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resp)

	replication := mustSingleNested(t, resp.Schema.Attributes["replication"], "replication")
	rules := mustListNested(t, replication.Attributes["rules"], "rules")
	destination := mustSingleNested(t, rules.NestedObject.Attributes["destination"], "destination")
	storageClass := mustStringAttribute(t, destination.Attributes["storage_class"], "storage_class")

	accepted := []string{
		"DEEP_ARCHIVE",
		"GLACIER",
		"GLACIER_IR",
		"HIGH_PERF",
		"INTELLIGENT_TIERING",
		"ONEZONE_IA",
		"STANDARD",
		"STANDARD_IA",
	}
	for _, value := range accepted {
		if runStorageClassValidators(t, storageClass, value) {
			t.Errorf("replication destination storage_class rejected %q but it should be accepted", value)
		}
	}

	if !runStorageClassValidators(t, storageClass, "NOPE") {
		t.Errorf("replication destination storage_class accepted bogus value %q but it should be rejected", "NOPE")
	}
}
