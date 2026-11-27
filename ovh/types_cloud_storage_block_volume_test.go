package ovh

import (
	"context"
	"encoding/json"
	"slices"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func blockVolumeUpdateTargetSpecKeys(t *testing.T, payload *CloudStorageBlockVolumeUpdatePayload) (map[string]json.RawMessage, string) {
	t.Helper()

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal update payload: %v", err)
	}

	var body struct {
		Checksum   string                     `json:"checksum"`
		TargetSpec map[string]json.RawMessage `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal update payload: %v", err)
	}

	return body.TargetSpec, body.Checksum
}

func TestCloudStorageBlockVolumeToUpdate_SendsOnlyMutableFields(t *testing.T) {
	model := CloudStorageBlockVolumeModel{
		Name:             ovhtypes.NewTfStringValue("test-volume"),
		Size:             types.Int64Value(20),
		Region:           ovhtypes.NewTfStringValue("GRA9"),
		AvailabilityZone: ovhtypes.NewTfStringValue("eu-west-gra-a"),
		VolumeType:       ovhtypes.NewTfStringValue("CLASSIC"),
		Encryption: types.ObjectValueMust(
			BlockVolumeEncryptionAttrTypes(),
			map[string]attr.Value{
				"enabled": types.BoolValue(true),
				"kms":     types.ObjectNull(BlockVolumeEncryptionKMSAttrTypes()),
			},
		),
		CreateFrom: types.ObjectValueMust(
			CreateFromAttrTypes(),
			map[string]attr.Value{
				"backup_id":   ovhtypes.NewTfStringValue("backup-1"),
				"snapshot_id": ovhtypes.TfStringValue{StringValue: types.StringNull()},
				"image_id":    ovhtypes.TfStringValue{StringValue: types.StringNull()},
			},
		),
	}

	spec, checksum := blockVolumeUpdateTargetSpecKeys(t, model.ToUpdate("checksum-123"))

	if checksum != "checksum-123" {
		t.Fatalf("unexpected checksum: got %q", checksum)
	}

	got := make([]string, 0, len(spec))
	for k := range spec {
		got = append(got, k)
	}
	sort.Strings(got)

	want := []string{"name", "size", "volumeType"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected targetSpec keys: got %v, want %v", got, want)
	}

	if string(spec["name"]) != `"test-volume"` {
		t.Fatalf("unexpected targetSpec.name: got %s", spec["name"])
	}
	if string(spec["size"]) != "20" {
		t.Fatalf("unexpected targetSpec.size: got %s", spec["size"])
	}
	if string(spec["volumeType"]) != `"CLASSIC"` {
		t.Fatalf("unexpected targetSpec.volumeType: got %s", spec["volumeType"])
	}
}

func TestCloudStorageBlockVolumeToUpdate_AlwaysCarriesNameAndSize(t *testing.T) {
	model := CloudStorageBlockVolumeModel{
		Name:       ovhtypes.NewTfStringValue(""),
		Size:       types.Int64Value(0),
		VolumeType: ovhtypes.NewTfStringValue("HIGH_SPEED"),
	}

	spec, _ := blockVolumeUpdateTargetSpecKeys(t, model.ToUpdate("checksum-123"))

	for _, key := range []string{"name", "size"} {
		if _, ok := spec[key]; !ok {
			t.Fatalf("targetSpec.%s must always be on the wire: a real PUT clears an absent field", key)
		}
	}
}

func TestCloudStorageBlockVolumeToUpdate_OmitsEmptyVolumeType(t *testing.T) {
	model := CloudStorageBlockVolumeModel{
		Name:       ovhtypes.NewTfStringValue("test-volume"),
		Size:       types.Int64Value(20),
		VolumeType: ovhtypes.TfStringValue{StringValue: types.StringNull()},
		CurrentState: types.ObjectValueMust(
			BlockVolumeCurrentStateAttrTypes(),
			map[string]attr.Value{
				"location": types.ObjectValueMust(
					map[string]attr.Type{"region": ovhtypes.TfStringType{}, "availability_zone": ovhtypes.TfStringType{}},
					map[string]attr.Value{
						"region":            ovhtypes.NewTfStringValue("GRA9"),
						"availability_zone": ovhtypes.NewTfStringValue("eu-west-gra-a"),
					},
				),
				"name":               ovhtypes.NewTfStringValue("test-volume"),
				"size":               types.Int64Value(20),
				"volume_type":        ovhtypes.NewTfStringValue("CLASSIC"),
				"bootable":           types.BoolValue(false),
				"status":             ovhtypes.NewTfStringValue("available"),
				"encryption":         types.ObjectNull(BlockVolumeEncryptionAttrTypes()),
				"attached_instances": types.ListNull(types.ObjectType{AttrTypes: BlockVolumeAttachedInstanceAttrTypes()}),
			},
		),
	}

	spec, _ := blockVolumeUpdateTargetSpecKeys(t, model.ToUpdate("checksum-123"))

	if v, ok := spec["volumeType"]; ok {
		t.Fatalf("empty volume_type must not reach the PUT (not a VolumeTypeEnum member), got %s", v)
	}
}

func TestCloudStorageBlockVolumeSchema_EncryptionEnabledRequiresReplace(t *testing.T) {
	r := &cloudStorageBlockVolumeResource{}
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	encryptionAttr, ok := resp.Schema.Attributes["encryption"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected encryption attribute to be SingleNestedAttribute, got %T", resp.Schema.Attributes["encryption"])
	}

	if len(encryptionAttr.PlanModifiers) == 0 {
		t.Fatal("expected encryption object attribute to have replace plan modifier")
	}

	enabledAttr, ok := encryptionAttr.Attributes["enabled"].(schema.BoolAttribute)
	if !ok {
		t.Fatalf("expected encryption.enabled attribute to be BoolAttribute, got %T", encryptionAttr.Attributes["enabled"])
	}

	if len(enabledAttr.PlanModifiers) == 0 {
		t.Fatal("expected encryption.enabled to have RequiresReplace plan modifier")
	}
}

func blockVolumeKMSObject(domainID, serviceKeyID string) types.Object {
	return types.ObjectValueMust(
		BlockVolumeEncryptionKMSAttrTypes(),
		map[string]attr.Value{
			"domain_id":      ovhtypes.NewTfStringValue(domainID),
			"service_key_id": ovhtypes.NewTfStringValue(serviceKeyID),
		},
	)
}

// Runs the modifiers in schema order, the way the framework chains them on one attribute.
func planBlockVolumeKMS(t *testing.T, configValue, planValue, stateValue types.Object) (types.Object, bool) {
	t.Helper()

	r := &cloudStorageBlockVolumeResource{}
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

	kmsAttr := schemaResp.Schema.Attributes["encryption"].(schema.SingleNestedAttribute).Attributes["kms"].(schema.SingleNestedAttribute)
	if !kmsAttr.Optional || !kmsAttr.Computed {
		t.Fatalf("encryption.kms must be Optional+Computed, got Optional=%t Computed=%t", kmsAttr.Optional, kmsAttr.Computed)
	}

	existingResource := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, map[string]tftypes.Value{})
	requiresReplace := false
	for _, modifier := range kmsAttr.PlanModifiers {
		req := planmodifier.ObjectRequest{
			Path:        path.Root("encryption").AtName("kms"),
			ConfigValue: configValue,
			PlanValue:   planValue,
			StateValue:  stateValue,
			Plan:        tfsdk.Plan{Raw: existingResource},
			State:       tfsdk.State{Raw: existingResource},
		}
		resp := &planmodifier.ObjectResponse{PlanValue: planValue}
		modifier.PlanModifyObject(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("plan modifier failed: %v", resp.Diagnostics)
		}
		planValue = resp.PlanValue
		requiresReplace = requiresReplace || resp.RequiresReplace
	}

	return planValue, requiresReplace
}

func TestCloudStorageBlockVolumeSchema_EncryptionKMS(t *testing.T) {
	keyInState := blockVolumeKMSObject("11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222")
	otherKey := blockVolumeKMSObject("33333333-3333-3333-3333-333333333333", "44444444-4444-4444-4444-444444444444")
	kmsNull := types.ObjectNull(BlockVolumeEncryptionKMSAttrTypes())
	kmsUnknown := types.ObjectUnknown(BlockVolumeEncryptionKMSAttrTypes())

	tests := []struct {
		name                string
		configValue         types.Object
		planValue           types.Object
		stateValue          types.Object
		wantPlanValue       types.Object
		wantRequiresReplace bool
	}{
		{
			name:          "key filled by the API from a source backup or snapshot, omitted in config",
			configValue:   kmsNull,
			planValue:     kmsUnknown,
			stateValue:    keyInState,
			wantPlanValue: keyInState,
		},
		{
			name:          "no key anywhere",
			configValue:   kmsNull,
			planValue:     kmsUnknown,
			stateValue:    kmsNull,
			wantPlanValue: kmsNull,
		},
		{
			name:          "configured key unchanged",
			configValue:   keyInState,
			planValue:     keyInState,
			stateValue:    keyInState,
			wantPlanValue: keyInState,
		},
		{
			name:                "configured key changed",
			configValue:         otherKey,
			planValue:           otherKey,
			stateValue:          keyInState,
			wantPlanValue:       otherKey,
			wantRequiresReplace: true,
		},
		{
			name:                "key added to a volume without one",
			configValue:         otherKey,
			planValue:           otherKey,
			stateValue:          kmsNull,
			wantPlanValue:       otherKey,
			wantRequiresReplace: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotPlanValue, gotRequiresReplace := planBlockVolumeKMS(t, tc.configValue, tc.planValue, tc.stateValue)

			if !gotPlanValue.Equal(tc.wantPlanValue) {
				t.Fatalf("plan value: got %s, want %s", gotPlanValue, tc.wantPlanValue)
			}
			if gotRequiresReplace != tc.wantRequiresReplace {
				t.Fatalf("requires replace: got %t, want %t", gotRequiresReplace, tc.wantRequiresReplace)
			}
		})
	}
}

func TestCloudStorageBlockVolumeMergeWith_EncryptionKMS(t *testing.T) {
	planWithKMSOmitted := types.ObjectValueMust(
		BlockVolumeEncryptionAttrTypes(),
		map[string]attr.Value{
			"enabled": types.BoolValue(true),
			"kms":     types.ObjectUnknown(BlockVolumeEncryptionKMSAttrTypes()),
		},
	)
	inheritedKey := &CloudStorageBlockVolumeEncryptionKMS{
		DomainID:     "11111111-1111-1111-1111-111111111111",
		ServiceKeyID: "22222222-2222-2222-2222-222222222222",
	}

	tests := []struct {
		name                 string
		targetSpecEncryption *CloudStorageBlockVolumeEncryption
		wantKMS              types.Object
	}{
		{
			name:                 "API returns the inherited key",
			targetSpecEncryption: &CloudStorageBlockVolumeEncryption{Enabled: true, Kms: inheritedKey},
			wantKMS:              blockVolumeKMSObject(inheritedKey.DomainID, inheritedKey.ServiceKeyID),
		},
		{
			name:                 "API returns encryption without a key",
			targetSpecEncryption: &CloudStorageBlockVolumeEncryption{Enabled: true},
			wantKMS:              types.ObjectNull(BlockVolumeEncryptionKMSAttrTypes()),
		},
		{
			name:                 "API returns no encryption",
			targetSpecEncryption: nil,
			wantKMS:              types.ObjectNull(BlockVolumeEncryptionKMSAttrTypes()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			model := CloudStorageBlockVolumeModel{Encryption: planWithKMSOmitted}

			model.MergeWith(context.Background(), &CloudStorageBlockVolumeAPIResponse{
				TargetSpec: &CloudStorageBlockVolumeTarget{
					Location:   &CloudStorageBlockVolumeLocation{Region: "GRA9"},
					Encryption: tc.targetSpecEncryption,
				},
			})

			attrs := model.Encryption.Attributes()
			if !attrs["kms"].Equal(tc.wantKMS) {
				t.Fatalf("encryption.kms: got %s, want %s", attrs["kms"], tc.wantKMS)
			}
			if !attrs["enabled"].Equal(types.BoolValue(true)) {
				t.Fatalf("encryption.enabled: got %s, want true", attrs["enabled"])
			}
		})
	}
}
