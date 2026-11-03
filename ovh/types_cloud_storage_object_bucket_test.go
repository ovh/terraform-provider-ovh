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

func TestBuildStorageObjectBucketObjectLock(t *testing.T) {
	days, years := int64(30), int64(1)
	cases := map[string]struct {
		lock     *CloudStorageObjectBucketAPIObjectLock
		wantNull bool
		mode     string
		days     *int64
		years    *int64
	}{
		"absent":           {lock: nil, wantNull: true},
		"disabled":         {lock: &CloudStorageObjectBucketAPIObjectLock{Status: "DISABLED"}, wantNull: true},
		"enabled, no rule": {lock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED"}},
		"days rule": {
			lock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED", Mode: "GOVERNANCE", RetentionDays: &days},
			mode: "GOVERNANCE", days: &days,
		},
		"years rule": {
			lock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED", Mode: "GOVERNANCE", RetentionYears: &years},
			mode: "GOVERNANCE", years: &years,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			spec := buildStorageObjectBucketObjectLockSpecObject(tc.lock)
			current := buildStorageObjectBucketObjectLockObject(tc.lock)
			if spec.IsNull() != tc.wantNull || current.IsNull() != tc.wantNull {
				t.Fatalf("null = (spec %t, current %t), want %t", spec.IsNull(), current.IsNull(), tc.wantNull)
			}
			if tc.wantNull {
				return
			}

			attrs := current.Attributes()
			if got := nullableTfString(tc.mode); !attrs["mode"].Equal(got) || !spec.Attributes()["mode"].Equal(got) {
				t.Fatalf("mode = %v, want %v", attrs["mode"], got)
			}
			if got := nullableInt64(tc.days); !attrs["retention_days"].Equal(got) || !spec.Attributes()["retention_days"].Equal(got) {
				t.Fatalf("retention_days = %v, want %v", attrs["retention_days"], got)
			}
			if got := nullableInt64(tc.years); !attrs["retention_years"].Equal(got) {
				t.Fatalf("retention_years = %v, want %v", attrs["retention_years"], got)
			}
		})
	}
}

func TestStorageObjectBucketObjectLockToAPI(t *testing.T) {
	cases := map[string]struct {
		mode attr.Value
		days attr.Value
		want string
	}{
		"no rule": {
			mode: ovhtypes.NewTfStringNull(),
			days: types.Int64Null(),
			want: `{"status":"ENABLED"}`,
		},
		"days rule": {
			mode: ovhtypes.NewTfStringValue("GOVERNANCE"),
			days: types.Int64Value(30),
			want: `{"status":"ENABLED","mode":"GOVERNANCE","retentionDays":30}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			lock := types.ObjectValueMust(StorageObjectBucketObjectLockSpecAttrTypes(), map[string]attr.Value{
				"mode":           tc.mode,
				"retention_days": tc.days,
			})
			model := CloudStorageObjectBucketModel{ObjectLock: lock, Tags: types.MapNull(ovhtypes.TfStringType{})}

			createSpec, err := json.Marshal(model.ToCreate().TargetSpec.ObjectLock)
			if err != nil {
				t.Fatalf("marshal create object lock: %v", err)
			}
			if string(createSpec) != tc.want {
				t.Fatalf("create objectLock = %s, want %s", createSpec, tc.want)
			}
			if got := storageObjectBucketUpdateTargetSpec(t, model.ToUpdate("checksum-123"))["objectLock"]; string(got) != tc.want {
				t.Fatalf("update objectLock = %s, want %s", got, tc.want)
			}
		})
	}
}
