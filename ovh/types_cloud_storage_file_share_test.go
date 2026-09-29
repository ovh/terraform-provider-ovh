package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
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
	wantCreateFrom := types.ObjectType{AttrTypes: FileShareCreateFromAttrTypes()}
	if got := resourceResp.Schema.Attributes["create_from"].GetType(); !got.Equal(wantCreateFrom) {
		t.Errorf("resource create_from type mismatch:\n got: %v\nwant: %v", got, wantCreateFrom)
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

func fileShareCreateFromSnapshot(snapshotId string) types.Object {
	return types.ObjectValueMust(
		FileShareCreateFromAttrTypes(),
		map[string]attr.Value{"snapshot_id": ovhtypes.NewTfStringValue(snapshotId)},
	)
}

func TestCloudStorageFileShareToCreateFromSnapshotOmitsSourceDerivedFields(t *testing.T) {
	model := CloudStorageFileShareModel{
		Name:           ovhtypes.NewTfStringValue("share-from-snapshot"),
		Protocol:       ovhtypes.NewTfStringValue("NFS"),
		Region:         ovhtypes.NewTfStringValue("region-a"),
		Size:           types.Int64Unknown(),
		ShareType:      ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		ShareNetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		Encryption:     types.ObjectUnknown(FileShareEncryptionAttrTypes()),
		CreateFrom:     fileShareCreateFromSnapshot("snapshot-id"),
	}

	targetSpec := model.ToCreate(context.Background()).TargetSpec
	if targetSpec.CreateFrom == nil || targetSpec.CreateFrom.SnapshotId != "snapshot-id" {
		t.Fatalf("expected createFrom.snapshotId in the create target spec, got %+v", targetSpec.CreateFrom)
	}
	if targetSpec.Size != 0 {
		t.Errorf("expected size to be omitted, got %d", targetSpec.Size)
	}
	if targetSpec.ShareType != "" {
		t.Errorf("expected shareType to be omitted, got %q", targetSpec.ShareType)
	}
	if targetSpec.ShareNetwork != nil {
		t.Errorf("expected shareNetwork to be omitted, got %+v", targetSpec.ShareNetwork)
	}
	if targetSpec.Encryption != nil {
		t.Errorf("expected encryption to be omitted, got %+v", targetSpec.Encryption)
	}
}

func TestCloudStorageFileShareToCreateWithoutSourceOmitsCreateFrom(t *testing.T) {
	model := CloudStorageFileShareModel{
		Size:           types.Int64Value(150),
		ShareType:      ovhtypes.NewTfStringValue("STANDARD_1AZ"),
		ShareNetworkId: ovhtypes.NewTfStringValue("share-network-id"),
		Encryption:     types.ObjectNull(FileShareEncryptionAttrTypes()),
		CreateFrom:     types.ObjectNull(FileShareCreateFromAttrTypes()),
	}

	targetSpec := model.ToCreate(context.Background()).TargetSpec
	if targetSpec.CreateFrom != nil {
		t.Errorf("expected no createFrom, got %+v", targetSpec.CreateFrom)
	}
	if targetSpec.Size != 150 || targetSpec.ShareType != "STANDARD_1AZ" {
		t.Errorf("expected size and shareType to be sent, got %d / %q", targetSpec.Size, targetSpec.ShareType)
	}
	if targetSpec.ShareNetwork == nil || targetSpec.ShareNetwork.Id != "share-network-id" {
		t.Errorf("expected shareNetwork to be sent, got %+v", targetSpec.ShareNetwork)
	}
}

func TestCloudStorageFileShareMergeWithTakesSourceValuesFromTargetSpec(t *testing.T) {
	model := CloudStorageFileShareModel{
		Size:           types.Int64Unknown(),
		ShareType:      ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		ShareNetworkId: ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
		Encryption:     types.ObjectUnknown(FileShareEncryptionAttrTypes()),
		CreateFrom:     fileShareCreateFromSnapshot("snapshot-id"),
	}

	model.MergeWith(context.Background(), &CloudStorageFileShareAPIResponse{
		TargetSpec: &CloudStorageFileShareAPITargetSpec{
			Size:         100,
			ShareType:    "STANDARD_1AZ",
			ShareNetwork: &CloudStorageFileShareAPIShareNetworkRef{Id: "source-share-network-id"},
			Encryption:   &CloudStorageFileShareAPIEncryption{Enabled: true},
			CreateFrom:   &CloudStorageFileShareAPICreateFrom{SnapshotId: "snapshot-id"},
		},
	})

	if model.Size.ValueInt64() != 100 {
		t.Errorf("expected size 100, got %d", model.Size.ValueInt64())
	}
	if model.ShareType.ValueString() != "STANDARD_1AZ" {
		t.Errorf("expected share_type from targetSpec, got %q", model.ShareType.ValueString())
	}
	if model.ShareNetworkId.ValueString() != "source-share-network-id" {
		t.Errorf("expected share_network_id from targetSpec, got %q", model.ShareNetworkId.ValueString())
	}
	if encryptionEnabled, isKnown := model.configuredEncryptionEnabled(); !isKnown || !encryptionEnabled {
		t.Errorf("expected encryption.enabled true, got %t (known: %t)", encryptionEnabled, isKnown)
	}
	if !model.CreateFrom.Equal(fileShareCreateFromSnapshot("snapshot-id")) {
		t.Errorf("expected create_from.snapshot_id to be kept, got %v", model.CreateFrom)
	}
}

func TestCloudStorageFileShareMergeWithImportsCreateFrom(t *testing.T) {
	model := CloudStorageFileShareModel{CreateFrom: types.ObjectNull(FileShareCreateFromAttrTypes())}

	model.MergeWith(context.Background(), &CloudStorageFileShareAPIResponse{
		TargetSpec: &CloudStorageFileShareAPITargetSpec{
			CreateFrom: &CloudStorageFileShareAPICreateFrom{SnapshotId: "snapshot-id"},
		},
	})

	if !model.CreateFrom.Equal(fileShareCreateFromSnapshot("snapshot-id")) {
		t.Errorf("expected create_from read from targetSpec, got %v", model.CreateFrom)
	}
}

func TestCloudStorageFileShareAttributesReplacedBySnapshotSource(t *testing.T) {
	encryptedSourceResponse := &CloudStorageFileShareAPIResponse{
		TargetSpec: &CloudStorageFileShareAPITargetSpec{
			ShareType:  "STANDARD_1AZ",
			Encryption: &CloudStorageFileShareAPIEncryption{Enabled: true},
		},
	}
	encryptionDisabled := types.ObjectValueMust(
		FileShareEncryptionAttrTypes(),
		map[string]attr.Value{"enabled": types.BoolValue(false)},
	)

	testCases := []struct {
		name                   string
		model                  CloudStorageFileShareModel
		wantReplacedAttributes int
	}{
		{
			name: "values omitted",
			model: CloudStorageFileShareModel{
				ShareType:  ovhtypes.TfStringValue{StringValue: types.StringUnknown()},
				Encryption: types.ObjectUnknown(FileShareEncryptionAttrTypes()),
				CreateFrom: fileShareCreateFromSnapshot("snapshot-id"),
			},
			wantReplacedAttributes: 0,
		},
		{
			name: "values matching the source",
			model: CloudStorageFileShareModel{
				ShareType: ovhtypes.NewTfStringValue("STANDARD_1AZ"),
				Encryption: types.ObjectValueMust(
					FileShareEncryptionAttrTypes(),
					map[string]attr.Value{"enabled": types.BoolValue(true)},
				),
				CreateFrom: fileShareCreateFromSnapshot("snapshot-id"),
			},
			wantReplacedAttributes: 0,
		},
		{
			name: "values differing from the source",
			model: CloudStorageFileShareModel{
				ShareType:  ovhtypes.NewTfStringValue("OTHER_TYPE"),
				Encryption: encryptionDisabled,
				CreateFrom: fileShareCreateFromSnapshot("snapshot-id"),
			},
			wantReplacedAttributes: 2,
		},
		{
			name: "no create_from",
			model: CloudStorageFileShareModel{
				ShareType:  ovhtypes.NewTfStringValue("OTHER_TYPE"),
				Encryption: encryptionDisabled,
				CreateFrom: types.ObjectNull(FileShareCreateFromAttrTypes()),
			},
			wantReplacedAttributes: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			replacedAttributes := testCase.model.attributesReplacedBySnapshotSource(encryptedSourceResponse)
			if len(replacedAttributes) != testCase.wantReplacedAttributes {
				t.Errorf("expected %d replaced attributes, got %v", testCase.wantReplacedAttributes, replacedAttributes)
			}
		})
	}
}
