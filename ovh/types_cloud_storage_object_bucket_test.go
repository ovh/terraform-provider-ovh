package ovh

import (
	"context"
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

func storageObjectBucketObjectLockValue(t *testing.T, status string, mode ovhtypes.TfStringValue, retentionDays types.Int64) types.Object {
	t.Helper()

	return types.ObjectValueMust(StorageObjectBucketObjectLockSpecAttrTypes(), map[string]attr.Value{
		"status":         ovhtypes.NewTfStringValue(status),
		"mode":           mode,
		"retention_days": retentionDays,
	})
}

func TestCloudStorageObjectBucketToCreate_ObjectLock(t *testing.T) {
	cases := map[string]struct {
		objectLock types.Object
		want       string
		omitted    bool
	}{
		"null object_lock is omitted": {
			objectLock: types.ObjectNull(StorageObjectBucketObjectLockSpecAttrTypes()),
			omitted:    true,
		},
		"disabled sends status only": {
			objectLock: storageObjectBucketObjectLockValue(t, "DISABLED", ovhtypes.TfStringValue{StringValue: types.StringNull()}, types.Int64Null()),
			want:       `{"status":"DISABLED"}`,
		},
		"enabled without rule sends status only": {
			objectLock: storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.TfStringValue{StringValue: types.StringNull()}, types.Int64Null()),
			want:       `{"status":"ENABLED"}`,
		},
		"enabled with rule sends mode and retentionDays": {
			objectLock: storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("GOVERNANCE"), types.Int64Value(30)),
			want:       `{"status":"ENABLED","mode":"GOVERNANCE","retentionDays":30}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			model := CloudStorageObjectBucketModel{
				Name:       ovhtypes.NewTfStringValue("bucket"),
				Region:     ovhtypes.NewTfStringValue("GRA"),
				ObjectLock: tc.objectLock,
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

			got, ok := body.TargetSpec["objectLock"]
			if tc.omitted {
				if ok {
					t.Fatalf("objectLock should be omitted: %s", got)
				}
				return
			}
			if !ok {
				t.Fatalf("objectLock missing from create target spec: %v", body.TargetSpec)
			}
			if string(got) != tc.want {
				t.Fatalf("objectLock = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestCloudStorageObjectBucketMergeWith_ObjectLock(t *testing.T) {
	ctx := context.Background()
	disabled := &CloudStorageObjectBucketAPIObjectLock{Status: "DISABLED"}

	t.Run("unconfigured object_lock stays null when current state reports DISABLED", func(t *testing.T) {
		model := CloudStorageObjectBucketModel{ObjectLock: types.ObjectNull(StorageObjectBucketObjectLockSpecAttrTypes())}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec:   &CloudStorageObjectBucketAPITargetSpec{Name: "bucket"},
			CurrentState: &CloudStorageObjectBucketAPICurrentState{Name: "bucket", ObjectLock: disabled},
		})

		if !model.ObjectLock.IsNull() {
			t.Fatalf("object_lock = %v, want null", model.ObjectLock)
		}
		current := model.CurrentState.Attributes()["object_lock"].(types.Object).Attributes()
		if got := current["status"].(ovhtypes.TfStringValue).ValueString(); got != "DISABLED" {
			t.Fatalf("current_state.object_lock.status = %q, want DISABLED", got)
		}
	})

	t.Run("enabled without rule maps mode and retention_days to null", func(t *testing.T) {
		model := CloudStorageObjectBucketModel{}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec: &CloudStorageObjectBucketAPITargetSpec{
				Name:       "bucket",
				ObjectLock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED"},
			},
		})

		attrs := model.ObjectLock.Attributes()
		if got := attrs["status"].(ovhtypes.TfStringValue).ValueString(); got != "ENABLED" {
			t.Fatalf("object_lock.status = %q, want ENABLED", got)
		}
		if !attrs["mode"].IsNull() || !attrs["retention_days"].IsNull() {
			t.Fatalf("object_lock rule should be null without a default retention rule: %v", attrs)
		}
	})

	t.Run("enabled with rule maps mode and retention_days", func(t *testing.T) {
		days := int64(30)
		model := CloudStorageObjectBucketModel{}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec: &CloudStorageObjectBucketAPITargetSpec{
				Name:       "bucket",
				ObjectLock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED", Mode: "GOVERNANCE", RetentionDays: &days},
			},
		})

		attrs := model.ObjectLock.Attributes()
		if got := attrs["mode"].(ovhtypes.TfStringValue).ValueString(); got != "GOVERNANCE" {
			t.Fatalf("object_lock.mode = %q, want GOVERNANCE", got)
		}
		if got := attrs["retention_days"].(types.Int64).ValueInt64(); got != 30 {
			t.Fatalf("object_lock.retention_days = %d, want 30", got)
		}
	})
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
