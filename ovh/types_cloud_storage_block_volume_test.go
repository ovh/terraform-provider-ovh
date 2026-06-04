package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func TestCloudStorageBlockVolumeToUpdate_DoesNotIncludeEncryption(t *testing.T) {
	model := CloudStorageBlockVolumeModel{
		Name:       ovhtypes.NewTfStringValue("test-volume"),
		Size:       types.Int64Value(20),
		VolumeType: ovhtypes.NewTfStringValue("CLASSIC"),
		Encryption: types.ObjectValueMust(
			BlockVolumeEncryptionAttrTypes(),
			map[string]attr.Value{
				"enabled": types.BoolValue(true),
			},
		),
	}

	payload := model.ToUpdate("checksum-123")
	if payload == nil {
		t.Fatal("expected update payload to be non-nil")
	}

	if payload.Checksum != "checksum-123" {
		t.Fatalf("unexpected checksum: got %q", payload.Checksum)
	}

	if payload.TargetSpec == nil {
		t.Fatal("expected targetSpec to be non-nil")
	}

	if payload.TargetSpec.Name != "test-volume" {
		t.Fatalf("unexpected targetSpec.name: got %q", payload.TargetSpec.Name)
	}

	if payload.TargetSpec.Size != 20 {
		t.Fatalf("unexpected targetSpec.size: got %d", payload.TargetSpec.Size)
	}

	if payload.TargetSpec.VolumeType != "CLASSIC" {
		t.Fatalf("unexpected targetSpec.volumeType: got %q", payload.TargetSpec.VolumeType)
	}

	if payload.TargetSpec.Encryption != nil {
		t.Fatalf("expected targetSpec.encryption to be nil in update payload, got %+v", payload.TargetSpec.Encryption)
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
