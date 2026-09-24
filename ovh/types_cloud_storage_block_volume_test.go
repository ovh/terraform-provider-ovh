package ovh

import (
	"context"
	"encoding/json"
	"slices"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
