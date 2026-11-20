package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func TestCloudStorageBlockVolumeBackupMergeWith_CopiesCurrentStateStatus(t *testing.T) {
	var response CloudStorageBlockVolumeBackupAPIResponse
	if err := json.Unmarshal([]byte(`{"id":"backup-id","currentState":{"name":"weekly","status":"ERROR_RESTORING"}}`), &response); err != nil {
		t.Fatalf("unmarshal backup response: %v", err)
	}

	var model CloudStorageBlockVolumeBackupModel
	model.MergeWith(context.Background(), &response)

	statusValue, ok := model.CurrentState.Attributes()["status"].(ovhtypes.TfStringValue)
	if !ok {
		t.Fatalf("expected current_state.status to be TfStringValue, got %T", model.CurrentState.Attributes()["status"])
	}
	if statusValue.ValueString() != "ERROR_RESTORING" {
		t.Fatalf("expected current_state.status ERROR_RESTORING, got %q", statusValue.ValueString())
	}
}

func mergeBackupResponseJSON(t *testing.T, responseJSON string) CloudStorageBlockVolumeBackupModel {
	t.Helper()

	var response CloudStorageBlockVolumeBackupAPIResponse
	if err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
		t.Fatalf("unmarshal backup response: %v", err)
	}

	var model CloudStorageBlockVolumeBackupModel
	model.MergeWith(context.Background(), &response)
	return model
}

func TestCloudStorageBlockVolumeBackupMergeWith_CopiesVolumeTypeAndCustomerKeyEncryption(t *testing.T) {
	model := mergeBackupResponseJSON(t, `{"id":"backup-id","currentState":{"name":"weekly","volumeType":"HIGH_SPEED_GEN2","encryption":{"enabled":true,"kms":{"domainId":"domain-1","serviceKeyId":"key-1"}}}}`)
	currentState := model.CurrentState.Attributes()

	volumeType, ok := currentState["volume_type"].(ovhtypes.TfStringValue)
	if !ok {
		t.Fatalf("expected current_state.volume_type to be TfStringValue, got %T", currentState["volume_type"])
	}
	if volumeType.ValueString() != "HIGH_SPEED_GEN2" {
		t.Fatalf("expected current_state.volume_type HIGH_SPEED_GEN2, got %q", volumeType.ValueString())
	}

	encryption, ok := currentState["encryption"].(types.Object)
	if !ok || encryption.IsNull() {
		t.Fatalf("expected current_state.encryption to be a non-null object, got %v", currentState["encryption"])
	}
	if enabled := encryption.Attributes()["enabled"]; !enabled.Equal(types.BoolValue(true)) {
		t.Fatalf("expected current_state.encryption.enabled true, got %v", enabled)
	}
	keyManagement, ok := encryption.Attributes()["kms"].(types.Object)
	if !ok || keyManagement.IsNull() {
		t.Fatalf("expected current_state.encryption.kms to be a non-null object, got %v", encryption.Attributes()["kms"])
	}
	if domainID := keyManagement.Attributes()["domain_id"]; !domainID.Equal(ovhtypes.NewTfStringValue("domain-1")) {
		t.Fatalf("expected kms.domain_id domain-1, got %v", domainID)
	}
	if serviceKeyID := keyManagement.Attributes()["service_key_id"]; !serviceKeyID.Equal(ovhtypes.NewTfStringValue("key-1")) {
		t.Fatalf("expected kms.service_key_id key-1, got %v", serviceKeyID)
	}
}

func TestCloudStorageBlockVolumeBackupMergeWith_PlatformEncryptionHasNullKms(t *testing.T) {
	model := mergeBackupResponseJSON(t, `{"id":"backup-id","currentState":{"volumeType":"CLASSIC","encryption":{"enabled":true}}}`)

	encryption, ok := model.CurrentState.Attributes()["encryption"].(types.Object)
	if !ok || encryption.IsNull() {
		t.Fatalf("expected current_state.encryption to be a non-null object, got %v", model.CurrentState.Attributes()["encryption"])
	}
	if keyManagement, ok := encryption.Attributes()["kms"].(types.Object); !ok || !keyManagement.IsNull() {
		t.Fatalf("expected current_state.encryption.kms to be null, got %v", encryption.Attributes()["kms"])
	}
}

func TestCloudStorageBlockVolumeBackupMergeWith_NullVolumeTypeAndEncryptionStayNull(t *testing.T) {
	model := mergeBackupResponseJSON(t, `{"id":"backup-id","currentState":{"name":"weekly","volumeType":null,"encryption":null}}`)
	currentState := model.CurrentState.Attributes()

	if volumeType, ok := currentState["volume_type"].(ovhtypes.TfStringValue); !ok || !volumeType.IsNull() {
		t.Fatalf("expected current_state.volume_type to be null, got %v", currentState["volume_type"])
	}
	if encryption, ok := currentState["encryption"].(types.Object); !ok || !encryption.IsNull() {
		t.Fatalf("expected current_state.encryption to be null, got %v", currentState["encryption"])
	}
}

func TestCloudStorageBlockVolumeBackupSchema_CurrentStateMatchesAttrTypes(t *testing.T) {
	backupResource := &cloudStorageBlockVolumeBackupResource{}
	var resp resource.SchemaResponse

	backupResource.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	currentStateAttr, ok := resp.Schema.Attributes["current_state"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected current_state attribute to be SingleNestedAttribute, got %T", resp.Schema.Attributes["current_state"])
	}

	currentStateAttrTypes := BackupCurrentStateAttrTypes()
	if len(currentStateAttr.Attributes) != len(currentStateAttrTypes) {
		t.Fatalf("schema has %d current_state attributes, attr types have %d", len(currentStateAttr.Attributes), len(currentStateAttrTypes))
	}
	for attributeName := range currentStateAttr.Attributes {
		if _, found := currentStateAttrTypes[attributeName]; !found {
			t.Errorf("current_state.%s is in the schema but not in BackupCurrentStateAttrTypes", attributeName)
		}
	}
}
