package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The attr type maps are shared between the resource and both data sources, so a
// field added to one side only produces a runtime "Missing Object Attribute Value"
// panic rather than a compile error.
func TestCloudStorageFileShareSchemaMatchesAttrTypes(t *testing.T) {
	ctx := context.Background()

	wantCurrentState := types.ObjectType{AttrTypes: FileShareCurrentStateAttrTypes()}
	wantEncryption := types.ObjectType{AttrTypes: FileShareEncryptionAttrTypes()}

	var resourceResp resource.SchemaResponse
	(&cloudStorageFileShareResource{}).Schema(ctx, resource.SchemaRequest{}, &resourceResp)
	if resourceResp.Diagnostics.HasError() {
		t.Fatalf("resource schema diagnostics: %v", resourceResp.Diagnostics)
	}

	if got := resourceResp.Schema.Attributes["current_state"].GetType(); !got.Equal(wantCurrentState) {
		t.Errorf("resource current_state type mismatch:\n got: %v\nwant: %v", got, wantCurrentState)
	}
	if got := resourceResp.Schema.Attributes["encryption"].GetType(); !got.Equal(wantEncryption) {
		t.Errorf("resource encryption type mismatch:\n got: %v\nwant: %v", got, wantEncryption)
	}

	var dataSourceResp datasource.SchemaResponse
	(&cloudStorageFileShareDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &dataSourceResp)
	if dataSourceResp.Diagnostics.HasError() {
		t.Fatalf("data source schema diagnostics: %v", dataSourceResp.Diagnostics)
	}

	if got := dataSourceResp.Schema.Attributes["current_state"].GetType(); !got.Equal(wantCurrentState) {
		t.Errorf("data source current_state type mismatch:\n got: %v\nwant: %v", got, wantCurrentState)
	}
	if got := dataSourceResp.Schema.Attributes["encryption"].GetType(); !got.Equal(wantEncryption) {
		t.Errorf("data source encryption type mismatch:\n got: %v\nwant: %v", got, wantEncryption)
	}

	var listResp datasource.SchemaResponse
	(&cloudStorageFileSharesDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &listResp)
	if listResp.Diagnostics.HasError() {
		t.Fatalf("list data source schema diagnostics: %v", listResp.Diagnostics)
	}

	wantList := types.ListType{ElemType: types.ObjectType{AttrTypes: fileShareListItemAttrTypes()}}
	if got := listResp.Schema.Attributes["file_shares"].GetType(); !got.Equal(wantList) {
		t.Errorf("list data source file_shares type mismatch:\n got: %v\nwant: %v", got, wantList)
	}
}

func TestCloudStorageFileShareToCreateIncludesEncryption(t *testing.T) {
	model := CloudStorageFileShareModel{
		Encryption: types.ObjectValueMust(
			FileShareEncryptionAttrTypes(),
			map[string]attr.Value{"enabled": types.BoolValue(true)},
		),
	}

	payload := model.ToCreate(context.Background())
	if payload.TargetSpec.Encryption == nil {
		t.Fatal("expected encryption in the create target spec")
	}
	if !payload.TargetSpec.Encryption.Enabled {
		t.Error("expected encryption.enabled to be true")
	}
}

func TestCloudStorageFileShareToCreateOmitsUnsetEncryption(t *testing.T) {
	model := CloudStorageFileShareModel{
		Encryption: types.ObjectNull(FileShareEncryptionAttrTypes()),
	}

	payload := model.ToCreate(context.Background())
	if payload.TargetSpec.Encryption != nil {
		t.Error("expected no encryption in the create target spec when unset")
	}
}
