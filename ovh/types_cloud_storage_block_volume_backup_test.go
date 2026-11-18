package ovh

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
