package ovh

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func storageObjectBucketUpdateTargetSpec(t *testing.T, payload *CloudStorageObjectBucketUpdatePayload) map[string]json.RawMessage {
	t.Helper()

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal update payload: %v", err)
	}

	var body struct {
		TargetSpec map[string]json.RawMessage `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal update payload: %v", err)
	}

	return body.TargetSpec
}

func TestCloudStorageObjectBucketToUpdate_TagsAlwaysSent(t *testing.T) {
	cases := map[string]struct {
		tags types.Map
		want string
	}{
		"null map clears tags": {
			tags: types.MapNull(ovhtypes.TfStringType{}),
			want: `{}`,
		},
		"empty map clears tags": {
			tags: types.MapValueMust(ovhtypes.TfStringType{}, map[string]attr.Value{}),
			want: `{}`,
		},
		"values are sent": {
			tags: types.MapValueMust(ovhtypes.TfStringType{}, map[string]attr.Value{
				"env": ovhtypes.NewTfStringValue("test"),
			}),
			want: `{"env":"test"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			model := CloudStorageObjectBucketModel{Tags: tc.tags}
			spec := storageObjectBucketUpdateTargetSpec(t, model.ToUpdate("checksum-123"))

			got, ok := spec["tags"]
			if !ok {
				t.Fatalf("tags key missing from update target spec: %v", spec)
			}
			if string(got) != tc.want {
				t.Fatalf("tags = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestCloudStorageObjectBucketToCreate_NullTagsOmitted(t *testing.T) {
	model := CloudStorageObjectBucketModel{
		Name:   ovhtypes.NewTfStringValue("bucket"),
		Region: ovhtypes.NewTfStringValue("GRA"),
		Tags:   types.MapNull(ovhtypes.TfStringType{}),
	}

	raw, err := json.Marshal(model.ToCreate())
	if err != nil {
		t.Fatalf("marshal create payload: %v", err)
	}

	var body struct {
		TargetSpec map[string]json.RawMessage `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal create payload: %v", err)
	}
	if _, ok := body.TargetSpec["tags"]; ok {
		t.Fatalf("tags should be omitted from create target spec when null: %v", body.TargetSpec)
	}
}
