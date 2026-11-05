package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

var storageObjectBucketNoMode = ovhtypes.TfStringValue{StringValue: types.StringNull()}

func storageObjectBucketObjectLockValue(t *testing.T, status string, mode ovhtypes.TfStringValue, retentionDays, retentionYears types.Int64) types.Object {
	t.Helper()

	return types.ObjectValueMust(StorageObjectBucketObjectLockAttrTypes(), map[string]attr.Value{
		"status":          ovhtypes.NewTfStringValue(status),
		"mode":            mode,
		"retention_days":  retentionDays,
		"retention_years": retentionYears,
	})
}

func TestCloudStorageObjectBucketToCreate_ObjectLock(t *testing.T) {
	cases := map[string]struct {
		objectLock types.Object
		want       string
		omitted    bool
	}{
		"null object_lock is omitted": {
			objectLock: types.ObjectNull(StorageObjectBucketObjectLockAttrTypes()),
			omitted:    true,
		},
		"disabled sends status only": {
			objectLock: storageObjectBucketObjectLockValue(t, "DISABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null()),
			want:       `{"status":"DISABLED"}`,
		},
		"enabled without rule sends status only": {
			objectLock: storageObjectBucketObjectLockValue(t, "ENABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null()),
			want:       `{"status":"ENABLED"}`,
		},
		"enabled with days rule sends mode and retentionDays": {
			objectLock: storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("GOVERNANCE"), types.Int64Value(30), types.Int64Null()),
			want:       `{"status":"ENABLED","mode":"GOVERNANCE","retentionDays":30}`,
		},
		"enabled with years rule sends mode and retentionYears": {
			objectLock: storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("GOVERNANCE"), types.Int64Null(), types.Int64Value(1)),
			want:       `{"status":"ENABLED","mode":"GOVERNANCE","retentionYears":1}`,
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

			updated := storageObjectBucketUpdateTargetSpec(t, model.ToUpdate("checksum-123"))["objectLock"]
			if string(updated) != tc.want {
				t.Fatalf("update objectLock = %s, want %s", updated, tc.want)
			}
		})
	}
}

func TestCloudStorageObjectBucketMergeWith_ObjectLock(t *testing.T) {
	ctx := context.Background()
	disabled := &CloudStorageObjectBucketAPIObjectLock{Status: "DISABLED"}

	t.Run("unconfigured object_lock stays null when current state reports DISABLED", func(t *testing.T) {
		model := CloudStorageObjectBucketModel{ObjectLock: types.ObjectNull(StorageObjectBucketObjectLockAttrTypes())}
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

	t.Run("unconfigured object_lock stays null when the target spec is a bare DISABLED", func(t *testing.T) {
		model := CloudStorageObjectBucketModel{ObjectLock: types.ObjectNull(StorageObjectBucketObjectLockAttrTypes())}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec: &CloudStorageObjectBucketAPITargetSpec{Name: "bucket", ObjectLock: disabled},
		})

		if !model.ObjectLock.IsNull() {
			t.Fatalf("object_lock = %v, want null", model.ObjectLock)
		}
	})

	t.Run("configured DISABLED object_lock is kept", func(t *testing.T) {
		model := CloudStorageObjectBucketModel{
			ObjectLock: storageObjectBucketObjectLockValue(t, "DISABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null()),
		}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec: &CloudStorageObjectBucketAPITargetSpec{Name: "bucket", ObjectLock: disabled},
		})

		if got := storageObjectBucketObjectLockStatus(model.ObjectLock); got != "DISABLED" {
			t.Fatalf("object_lock.status = %q, want DISABLED", got)
		}
	})

	t.Run("enabled without rule maps mode and retention period to null", func(t *testing.T) {
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
		if !attrs["mode"].IsNull() || !attrs["retention_days"].IsNull() || !attrs["retention_years"].IsNull() {
			t.Fatalf("object_lock rule should be null without a default retention rule: %v", attrs)
		}
	})

	t.Run("enabled with days rule maps mode and retention_days", func(t *testing.T) {
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
		if !attrs["retention_years"].IsNull() {
			t.Fatalf("object_lock.retention_years = %v, want null", attrs["retention_years"])
		}
	})

	t.Run("enabled with years rule round-trips retention_years only", func(t *testing.T) {
		years := int64(1)
		model := CloudStorageObjectBucketModel{}
		model.MergeWith(ctx, &CloudStorageObjectBucketAPIResponse{
			TargetSpec: &CloudStorageObjectBucketAPITargetSpec{
				Name:       "bucket",
				ObjectLock: &CloudStorageObjectBucketAPIObjectLock{Status: "ENABLED", Mode: "GOVERNANCE", RetentionYears: &years},
			},
		})

		attrs := model.ObjectLock.Attributes()
		if got := attrs["retention_years"].(types.Int64).ValueInt64(); got != 1 {
			t.Fatalf("object_lock.retention_years = %d, want 1", got)
		}
		if !attrs["retention_days"].IsNull() {
			t.Fatalf("object_lock.retention_days = %v, want null", attrs["retention_days"])
		}
		got := storageObjectBucketObjectLockToAPI(model.ObjectLock)
		if got.RetentionDays != nil || got.RetentionYears == nil || *got.RetentionYears != 1 {
			t.Fatalf("round trip objectLock = %+v, want retentionYears 1 only", got)
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

func TestStorageObjectBucketObjectLockStatusChanged(t *testing.T) {
	ctx := context.Background()
	null := types.ObjectNull(StorageObjectBucketObjectLockAttrTypes())
	disabled := storageObjectBucketObjectLockValue(t, "DISABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null())
	enabled := storageObjectBucketObjectLockValue(t, "ENABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null())
	rule30 := storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("GOVERNANCE"), types.Int64Value(30), types.Int64Null())
	rule60 := storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("COMPLIANCE"), types.Int64Value(60), types.Int64Null())
	rule1y := storageObjectBucketObjectLockValue(t, "ENABLED", ovhtypes.NewTfStringValue("GOVERNANCE"), types.Int64Null(), types.Int64Value(1))
	unknownStatus := types.ObjectValueMust(StorageObjectBucketObjectLockAttrTypes(), map[string]attr.Value{
		"status":          ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		"mode":            storageObjectBucketNoMode,
		"retention_days":  types.Int64Null(),
		"retention_years": types.Int64Null(),
	})
	rule30NoStatus := types.ObjectValueMust(StorageObjectBucketObjectLockAttrTypes(), map[string]attr.Value{
		"status":          ovhtypes.TfStringValue{StringValue: types.StringNull()},
		"mode":            ovhtypes.NewTfStringValue("GOVERNANCE"),
		"retention_days":  types.Int64Value(30),
		"retention_years": types.Int64Null(),
	})

	cases := map[string]struct {
		state, plan types.Object
		want        bool
	}{
		"null to DISABLED is in place":        {state: null, plan: disabled, want: false},
		"DISABLED to null is in place":        {state: disabled, plan: null, want: false},
		"retention rule change is in place":   {state: rule30, plan: rule60, want: false},
		"days to years rule is in place":      {state: rule30, plan: rule1y, want: false},
		"retention rule removal is in place":  {state: rule30, plan: enabled, want: false},
		"null to ENABLED replaces":            {state: null, plan: enabled, want: true},
		"DISABLED to ENABLED replaces":        {state: disabled, plan: rule30, want: true},
		"ENABLED to null replaces":            {state: rule30, plan: null, want: true},
		"ENABLED to DISABLED replaces":        {state: enabled, plan: disabled, want: true},
		"unknown planned object is in place":  {state: disabled, plan: types.ObjectUnknown(StorageObjectBucketObjectLockAttrTypes()), want: false},
		"unknown planned status is in place":  {state: enabled, plan: unknownStatus, want: false},
		"2.24.0 state to ENABLED is in place": {state: rule30NoStatus, plan: rule30, want: false},
		"2.24.0 state to DISABLED replaces":   {state: rule30NoStatus, plan: disabled, want: true},
		"2.24.0 state to null replaces":       {state: rule30NoStatus, plan: null, want: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := &objectplanmodifier.RequiresReplaceIfFuncResponse{}
			storageObjectBucketObjectLockStatusChanged(ctx, planmodifier.ObjectRequest{StateValue: tc.state, PlanValue: tc.plan}, resp)
			if resp.RequiresReplace != tc.want {
				t.Fatalf("RequiresReplace = %t, want %t", resp.RequiresReplace, tc.want)
			}
		})
	}
}

func TestStorageObjectBucketObjectLockRetentionRule(t *testing.T) {
	ctx := context.Background()
	governance := ovhtypes.NewTfStringValue("GOVERNANCE")

	cases := map[string]struct {
		lock    types.Object
		wantErr bool
	}{
		"enabled without rule":     {lock: storageObjectBucketObjectLockValue(t, "ENABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Null())},
		"mode with days":           {lock: storageObjectBucketObjectLockValue(t, "ENABLED", governance, types.Int64Value(30), types.Int64Null())},
		"mode with years":          {lock: storageObjectBucketObjectLockValue(t, "ENABLED", governance, types.Int64Null(), types.Int64Value(1))},
		"mode with days and years": {lock: storageObjectBucketObjectLockValue(t, "ENABLED", governance, types.Int64Value(30), types.Int64Value(1)), wantErr: true},
		"mode without period":      {lock: storageObjectBucketObjectLockValue(t, "ENABLED", governance, types.Int64Null(), types.Int64Null()), wantErr: true},
		"days without mode":        {lock: storageObjectBucketObjectLockValue(t, "ENABLED", storageObjectBucketNoMode, types.Int64Value(30), types.Int64Null()), wantErr: true},
		"years without mode":       {lock: storageObjectBucketObjectLockValue(t, "ENABLED", storageObjectBucketNoMode, types.Int64Null(), types.Int64Value(1)), wantErr: true},
		"years rule with DISABLED": {lock: storageObjectBucketObjectLockValue(t, "DISABLED", governance, types.Int64Null(), types.Int64Value(1)), wantErr: true},
		"days rule without status": {lock: types.ObjectValueMust(StorageObjectBucketObjectLockAttrTypes(), map[string]attr.Value{
			"status":          ovhtypes.TfStringValue{StringValue: types.StringNull()},
			"mode":            governance,
			"retention_days":  types.Int64Value(30),
			"retention_years": types.Int64Null(),
		})},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := &validator.ObjectResponse{}
			storageObjectBucketObjectLockRetentionRule{}.ValidateObject(ctx, validator.ObjectRequest{
				Path:        path.Root("object_lock"),
				ConfigValue: tc.lock,
			}, resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Fatalf("HasError = %t, want %t: %v", resp.Diagnostics.HasError(), tc.wantErr, resp.Diagnostics)
			}
		})
	}
}
