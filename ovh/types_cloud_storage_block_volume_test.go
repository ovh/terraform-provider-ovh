package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// blockVolumeAPIResponse builds a minimal API response for MergeWith tests.
// targetEncryption and currentEncryption may be nil; currentState is omitted
// entirely when withCurrentState is false.
func blockVolumeAPIResponse(withCurrentState bool, targetEncryption, currentEncryption *CloudStorageBlockVolumeEncryption) *CloudStorageBlockVolumeAPIResponse {
	response := &CloudStorageBlockVolumeAPIResponse{
		Id:             "volume-id",
		Checksum:       "checksum-123",
		CreatedAt:      "2026-01-01T00:00:00Z",
		UpdatedAt:      "2026-01-01T00:00:00Z",
		ResourceStatus: "READY",
		TargetSpec: &CloudStorageBlockVolumeTarget{
			Location:   &CloudStorageBlockVolumeLocation{Region: "GRA9"},
			Name:       "test-volume",
			Size:       20,
			VolumeType: "CLASSIC",
			Encryption: targetEncryption,
		},
	}

	if withCurrentState {
		response.CurrentState = &CloudStorageBlockVolumeCurrentState{
			Location:   &CloudStorageBlockVolumeLocation{Region: "GRA9"},
			Name:       "test-volume",
			Size:       20,
			VolumeType: "CLASSIC",
			Status:     "available",
			Encryption: currentEncryption,
		}
	}

	return response
}

func encryptionObjectValue(t *testing.T, enabled bool) types.Object {
	t.Helper()
	return types.ObjectValueMust(
		BlockVolumeEncryptionAttrTypes(),
		map[string]attr.Value{
			"enabled": types.BoolValue(enabled),
		},
	)
}

// TestCloudStorageBlockVolumeMergeWith_EncryptionBackfillFromCurrentState
// covers the defensive backfill: when the API GET omits targetSpec.encryption,
// the root encryption attribute must be filled from currentState.encryption so
// that UseStateForUnknown keeps plans stable for volumes created without an
// encryption block.
func TestCloudStorageBlockVolumeMergeWith_EncryptionBackfillFromCurrentState(t *testing.T) {
	ctx := context.Background()

	t.Run("targetSpec encryption nil, currentState encryption set", func(t *testing.T) {
		model := CloudStorageBlockVolumeModel{
			Encryption: types.ObjectNull(BlockVolumeEncryptionAttrTypes()),
		}
		model.MergeWith(ctx, blockVolumeAPIResponse(true, nil, &CloudStorageBlockVolumeEncryption{Enabled: false}))

		want := encryptionObjectValue(t, false)
		if !model.Encryption.Equal(want) {
			t.Fatalf("expected encryption to be backfilled from currentState (%s), got %s", want, model.Encryption)
		}
	})

	t.Run("targetSpec encryption set wins over currentState", func(t *testing.T) {
		model := CloudStorageBlockVolumeModel{
			Encryption: types.ObjectUnknown(BlockVolumeEncryptionAttrTypes()),
		}
		model.MergeWith(ctx, blockVolumeAPIResponse(true,
			&CloudStorageBlockVolumeEncryption{Enabled: true},
			&CloudStorageBlockVolumeEncryption{Enabled: false},
		))

		want := encryptionObjectValue(t, true)
		if !model.Encryption.Equal(want) {
			t.Fatalf("expected encryption to come from targetSpec (%s), got %s", want, model.Encryption)
		}
	})

	t.Run("both targetSpec and currentState encryption absent leaves value as-is", func(t *testing.T) {
		prior := encryptionObjectValue(t, true)
		model := CloudStorageBlockVolumeModel{
			Encryption: prior,
		}
		model.MergeWith(ctx, blockVolumeAPIResponse(true, nil, nil))

		if !model.Encryption.Equal(prior) {
			t.Fatalf("expected encryption to be left untouched (%s), got %s", prior, model.Encryption)
		}
	})

	t.Run("currentState absent entirely leaves value as-is", func(t *testing.T) {
		prior := types.ObjectNull(BlockVolumeEncryptionAttrTypes())
		model := CloudStorageBlockVolumeModel{
			Encryption: prior,
		}
		model.MergeWith(ctx, blockVolumeAPIResponse(false, nil, nil))

		if !model.Encryption.Equal(prior) {
			t.Fatalf("expected encryption to be left untouched (%s), got %s", prior, model.Encryption)
		}
	})

	t.Run("both absent with unknown value falls back to null", func(t *testing.T) {
		model := CloudStorageBlockVolumeModel{
			Encryption: types.ObjectUnknown(BlockVolumeEncryptionAttrTypes()),
		}
		model.MergeWith(ctx, blockVolumeAPIResponse(true, nil, nil))

		if !model.Encryption.IsNull() {
			t.Fatalf("expected unknown encryption to be resolved to null, got %s", model.Encryption)
		}
	})
}

func TestCloudStorageBlockVolumeSchema_EncryptionEnabledRequiresReplace(t *testing.T) {
	ctx := context.Background()
	r := &cloudStorageBlockVolumeResource{}
	var resp resource.SchemaResponse

	r.Schema(ctx, resource.SchemaRequest{}, &resp)

	encryptionAttr, ok := resp.Schema.Attributes["encryption"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected encryption attribute to be SingleNestedAttribute, got %T", resp.Schema.Attributes["encryption"])
	}

	// UseStateForUnknown MUST come before RequiresReplace: it restores the
	// state value when the framework marks the unconfigured computed attribute
	// unknown, so RequiresReplace no longer sees a spurious diff.
	useStateDesc := objectplanmodifier.UseStateForUnknown().Description(ctx)
	replaceDesc := objectplanmodifier.RequiresReplace().Description(ctx)

	if len(encryptionAttr.PlanModifiers) != 2 {
		t.Fatalf("expected encryption to have exactly 2 plan modifiers (UseStateForUnknown, RequiresReplace), got %d", len(encryptionAttr.PlanModifiers))
	}
	if got := encryptionAttr.PlanModifiers[0].Description(ctx); got != useStateDesc {
		t.Fatalf("expected encryption plan modifier #0 to be UseStateForUnknown, got %q", got)
	}
	if got := encryptionAttr.PlanModifiers[1].Description(ctx); got != replaceDesc {
		t.Fatalf("expected encryption plan modifier #1 to be RequiresReplace, got %q", got)
	}

	enabledAttr, ok := encryptionAttr.Attributes["enabled"].(schema.BoolAttribute)
	if !ok {
		t.Fatalf("expected encryption.enabled attribute to be BoolAttribute, got %T", encryptionAttr.Attributes["enabled"])
	}

	boolUseStateDesc := boolplanmodifier.UseStateForUnknown().Description(ctx)
	boolReplaceDesc := boolplanmodifier.RequiresReplace().Description(ctx)

	if len(enabledAttr.PlanModifiers) != 2 {
		t.Fatalf("expected encryption.enabled to have exactly 2 plan modifiers (UseStateForUnknown, RequiresReplace), got %d", len(enabledAttr.PlanModifiers))
	}
	if got := enabledAttr.PlanModifiers[0].Description(ctx); got != boolUseStateDesc {
		t.Fatalf("expected encryption.enabled plan modifier #0 to be UseStateForUnknown, got %q", got)
	}
	if got := enabledAttr.PlanModifiers[1].Description(ctx); got != boolReplaceDesc {
		t.Fatalf("expected encryption.enabled plan modifier #1 to be RequiresReplace, got %q", got)
	}
}

// blockStorageResourceSchemas returns the schemas of the three block storage
// resources, keyed by a short name for test reporting.
func blockStorageResourceSchemas(ctx context.Context, t *testing.T) map[string]schema.Schema {
	t.Helper()

	schemas := map[string]schema.Schema{}

	var volumeResp resource.SchemaResponse
	(&cloudStorageBlockVolumeResource{}).Schema(ctx, resource.SchemaRequest{}, &volumeResp)
	schemas["volume"] = volumeResp.Schema

	var backupResp resource.SchemaResponse
	(&cloudStorageBlockVolumeBackupResource{}).Schema(ctx, resource.SchemaRequest{}, &backupResp)
	schemas["backup"] = backupResp.Schema

	var snapshotResp resource.SchemaResponse
	(&cloudStorageBlockVolumeSnapshotResource{}).Schema(ctx, resource.SchemaRequest{}, &snapshotResp)
	schemas["snapshot"] = snapshotResp.Schema

	return schemas
}

// stringAttrHasUseStateForUnknown checks that a string attribute carries the
// UseStateForUnknown plan modifier.
func stringAttrHasUseStateForUnknown(ctx context.Context, t *testing.T, s schema.Schema, name string) bool {
	t.Helper()

	a, ok := s.Attributes[name].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected %q attribute to be StringAttribute, got %T", name, s.Attributes[name])
	}

	useStateDesc := stringplanmodifier.UseStateForUnknown().Description(ctx)
	for _, pm := range a.PlanModifiers {
		if pm.Description(ctx) == useStateDesc {
			return true
		}
	}
	return false
}

// TestCloudStorageBlockSchemas_IDAndCreatedAtUseStateForUnknown ensures the
// computed id and created_at attributes keep their known state value during
// updates. Without UseStateForUnknown on id, an in-place volume update marks
// id as unknown, which cascades a spurious replacement to backups/snapshots
// referencing it through their RequiresReplace volume_id attribute.
func TestCloudStorageBlockSchemas_IDAndCreatedAtUseStateForUnknown(t *testing.T) {
	ctx := context.Background()

	for resName, s := range blockStorageResourceSchemas(ctx, t) {
		for _, attrName := range []string{"id", "created_at"} {
			if !stringAttrHasUseStateForUnknown(ctx, t, s, attrName) {
				t.Errorf("%s: expected %q to have UseStateForUnknown plan modifier", resName, attrName)
			}
		}
	}
}

// TestCloudStorageBlockSchemas_MutableComputedAttrsDoNotPinState ensures the
// attributes that legitimately change on update do NOT carry UseStateForUnknown.
func TestCloudStorageBlockSchemas_MutableComputedAttrsDoNotPinState(t *testing.T) {
	ctx := context.Background()

	for resName, s := range blockStorageResourceSchemas(ctx, t) {
		for _, attrName := range []string{"checksum", "updated_at", "resource_status"} {
			if stringAttrHasUseStateForUnknown(ctx, t, s, attrName) {
				t.Errorf("%s: %q legitimately changes on update and must not have UseStateForUnknown", resName, attrName)
			}
		}
	}
}

func runStringValidators(ctx context.Context, validators []validator.String, p path.Path, value types.String) *validator.StringResponse {
	req := validator.StringRequest{
		Path:           p,
		PathExpression: path.MatchRoot(p.String()),
		ConfigValue:    value,
	}
	resp := &validator.StringResponse{}
	for _, v := range validators {
		v.ValidateString(ctx, req, resp)
	}
	return resp
}

func TestCloudStorageBlockVolumeSchema_VolumeTypeValidator(t *testing.T) {
	ctx := context.Background()
	s := blockStorageResourceSchemas(ctx, t)["volume"]

	a, ok := s.Attributes["volume_type"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected volume_type to be StringAttribute, got %T", s.Attributes["volume_type"])
	}
	if len(a.Validators) == 0 {
		t.Fatal("expected volume_type to have a OneOf validator")
	}

	for _, valid := range []string{"CLASSIC", "HIGH_SPEED", "HIGH_SPEED_GEN2"} {
		resp := runStringValidators(ctx, a.Validators, path.Root("volume_type"), types.StringValue(valid))
		if resp.Diagnostics.HasError() {
			t.Errorf("expected %q to be a valid volume_type, got: %v", valid, resp.Diagnostics)
		}
	}

	resp := runStringValidators(ctx, a.Validators, path.Root("volume_type"), types.StringValue("NOT_A_VOLUME_TYPE"))
	if !resp.Diagnostics.HasError() {
		t.Error("expected an invalid volume_type to be rejected at plan time")
	}
}

func TestCloudStorageBlockVolumeSchema_SizeValidator(t *testing.T) {
	ctx := context.Background()
	s := blockStorageResourceSchemas(ctx, t)["volume"]

	a, ok := s.Attributes["size"].(schema.Int64Attribute)
	if !ok {
		t.Fatalf("expected size to be Int64Attribute, got %T", s.Attributes["size"])
	}
	if len(a.Validators) == 0 {
		t.Fatal("expected size to have an AtLeast validator")
	}

	for _, tc := range []struct {
		value     int64
		expectErr bool
	}{
		{0, true},
		{-1, true},
		{1, false},
		{100, false},
	} {
		req := validator.Int64Request{
			Path:           path.Root("size"),
			PathExpression: path.MatchRoot("size"),
			ConfigValue:    types.Int64Value(tc.value),
		}
		resp := &validator.Int64Response{}
		for _, v := range a.Validators {
			v.ValidateInt64(ctx, req, resp)
		}
		if resp.Diagnostics.HasError() != tc.expectErr {
			t.Errorf("size=%d: expected error=%t, got diagnostics: %v", tc.value, tc.expectErr, resp.Diagnostics)
		}
	}
}

func TestCloudStorageBlockSchemas_NameValidator(t *testing.T) {
	ctx := context.Background()

	for resName, s := range blockStorageResourceSchemas(ctx, t) {
		a, ok := s.Attributes["name"].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s: expected name to be StringAttribute, got %T", resName, s.Attributes["name"])
		}
		if len(a.Validators) == 0 {
			t.Fatalf("%s: expected name to have a LengthAtLeast validator", resName)
		}

		resp := runStringValidators(ctx, a.Validators, path.Root("name"), types.StringValue(""))
		if !resp.Diagnostics.HasError() {
			t.Errorf("%s: expected an empty name to be rejected at plan time", resName)
		}

		resp = runStringValidators(ctx, a.Validators, path.Root("name"), types.StringValue("valid-name"))
		if resp.Diagnostics.HasError() {
			t.Errorf("%s: expected a non-empty name to be accepted, got: %v", resName, resp.Diagnostics)
		}
	}
}

func TestCloudStorageBlockVolumeSchema_CreateFromExactlyOneOf(t *testing.T) {
	ctx := context.Background()
	s := blockStorageResourceSchemas(ctx, t)["volume"]

	createFromAttr, ok := s.Attributes["create_from"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected create_from to be SingleNestedAttribute, got %T", s.Attributes["create_from"])
	}

	backupIDAttr, ok := createFromAttr.Attributes["backup_id"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected create_from.backup_id to be StringAttribute, got %T", createFromAttr.Attributes["backup_id"])
	}

	// ExactlyOneOf needs a full config to run, so only assert its presence.
	if len(backupIDAttr.Validators) == 0 {
		t.Fatal("expected create_from.backup_id to have an ExactlyOneOf validator covering backup_id/snapshot_id/image_id")
	}
}
