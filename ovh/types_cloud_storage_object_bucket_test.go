package ovh

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestCloudStorageObjectBucketMergeWith_ObjectLock(t *testing.T) {
	cases := map[string]struct {
		objectLock string
		wantNull   bool
	}{
		// What the API returns for a bucket created without Object Lock, e.g. through
		// ovh_cloud_project_storage: no mode, no retention.
		"disabled lock is null": {
			objectLock: `{"status": "DISABLED"}`,
			wantNull:   true,
		},
		"absent lock is null": {
			objectLock: `null`,
			wantNull:   true,
		},
		"armed lock is kept": {
			objectLock: `{"status": "ENABLED", "mode": "COMPLIANCE", "retentionDays": 30}`,
			wantNull:   false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			raw := fmt.Sprintf(`{
				"id": "GRA_bucket",
				"resourceStatus": "READY",
				"targetSpec": {"name": "bucket", "location": {"region": "GRA"}, "objectLock": %s},
				"currentState": {"name": "bucket", "location": {"region": "GRA"}, "objectLock": %s}
			}`, tc.objectLock, tc.objectLock)

			var response CloudStorageObjectBucketAPIResponse
			if err := json.Unmarshal([]byte(raw), &response); err != nil {
				t.Fatalf("unmarshal response: %v", err)
			}

			model := CloudStorageObjectBucketModel{}
			model.MergeWith(context.Background(), &response)

			if model.ObjectLock.IsNull() != tc.wantNull {
				t.Fatalf("object_lock = %v, want null: %v", model.ObjectLock, tc.wantNull)
			}
			current := model.CurrentState.Attributes()["object_lock"].(types.Object)
			if current.IsNull() != tc.wantNull {
				t.Fatalf("current_state.object_lock = %v, want null: %v", current, tc.wantNull)
			}
			if !tc.wantNull {
				if got := model.ObjectLock.Attributes()["mode"].(ovhtypes.TfStringValue).ValueString(); got != "COMPLIANCE" {
					t.Fatalf("object_lock.mode = %q, want COMPLIANCE", got)
				}
			}
		})
	}
}

func TestCloudStorageObjectBucketToCreate_ObjectLockStatusSent(t *testing.T) {
	model := CloudStorageObjectBucketModel{
		Name:   ovhtypes.NewTfStringValue("bucket"),
		Region: ovhtypes.NewTfStringValue("GRA"),
		ObjectLock: types.ObjectValueMust(StorageObjectBucketObjectLockSpecAttrTypes(), map[string]attr.Value{
			"mode":           ovhtypes.NewTfStringValue("COMPLIANCE"),
			"retention_days": types.Int64Value(30),
		}),
	}

	raw, err := json.Marshal(model.ToCreate())
	if err != nil {
		t.Fatalf("marshal create payload: %v", err)
	}

	var body struct {
		TargetSpec struct {
			ObjectLock map[string]json.RawMessage `json:"objectLock"`
		} `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal create payload: %v", err)
	}
	// status is mandatory in the API contract when arming the lock.
	if got := string(body.TargetSpec.ObjectLock["status"]); got != `"ENABLED"` {
		t.Fatalf("objectLock.status = %s, want \"ENABLED\": %v", got, body.TargetSpec.ObjectLock)
	}
}
