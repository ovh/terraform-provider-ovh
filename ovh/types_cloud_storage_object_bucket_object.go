package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketObjectAPILock struct {
	Mode        string `json:"mode"`
	RetainUntil string `json:"retainUntil"`
}

type CloudStorageObjectBucketObjectAPIRestoreStatus struct {
	InProgress bool   `json:"inProgress"`
	ExpiresAt  string `json:"expiresAt,omitempty"`
}

type CloudStorageObjectBucketObjectAPIResponse struct {
	Key               string                                          `json:"key"`
	Size              int64                                           `json:"size"`
	Etag              string                                          `json:"etag,omitempty"`
	LastModifiedAt    string                                          `json:"lastModifiedAt,omitempty"`
	StorageClass      string                                          `json:"storageClass,omitempty"`
	VersionId         string                                          `json:"versionId,omitempty"`
	IsLatest          bool                                            `json:"isLatest,omitempty"`
	IsDeleteMarker    bool                                            `json:"isDeleteMarker,omitempty"`
	IsCommonPrefix    bool                                            `json:"isCommonPrefix,omitempty"`
	Lock              *CloudStorageObjectBucketObjectAPILock          `json:"lock,omitempty"`
	LegalHold         string                                          `json:"legalHold,omitempty"`
	RestoreStatus     *CloudStorageObjectBucketObjectAPIRestoreStatus `json:"restoreStatus,omitempty"`
	ReplicationStatus string                                          `json:"replicationStatus,omitempty"`
}

func storageObjectBucketObjectLockAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":         ovhtypes.TfStringType{},
		"retain_until": ovhtypes.TfStringType{},
	}
}

func storageObjectBucketObjectRestoreStatusAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"in_progress": types.BoolType,
		"expires_at":  ovhtypes.TfStringType{},
	}
}

func storageObjectBucketObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":                ovhtypes.TfStringType{},
		"size":               types.Int64Type,
		"etag":               ovhtypes.TfStringType{},
		"last_modified_at":   ovhtypes.TfStringType{},
		"storage_class":      ovhtypes.TfStringType{},
		"version_id":         ovhtypes.TfStringType{},
		"is_latest":          types.BoolType,
		"is_delete_marker":   types.BoolType,
		"is_common_prefix":   types.BoolType,
		"legal_hold":         ovhtypes.TfStringType{},
		"lock":               types.ObjectType{AttrTypes: storageObjectBucketObjectLockAttrTypes()},
		"restore_status":     types.ObjectType{AttrTypes: storageObjectBucketObjectRestoreStatusAttrTypes()},
		"replication_status": ovhtypes.TfStringType{},
	}
}

func storageObjectBucketObjectDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Object key (path within the bucket)",
		},
		"size": schema.Int64Attribute{
			Computed:    true,
			Description: "Object size in bytes",
		},
		"etag": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Entity tag (hash of the object content), null for a delete marker",
		},
		"last_modified_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last modification timestamp, null for a common prefix",
		},
		"storage_class": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Storage class of the object (DEEP_ARCHIVE, GLACIER_IR, HIGH_PERF, STANDARD, STANDARD_IA), null for a delete marker or a common prefix",
		},
		"version_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Version identifier",
		},
		"is_latest": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this is the latest version of the object",
		},
		"is_delete_marker": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this version is a delete marker",
		},
		"is_common_prefix": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this entry represents a common prefix (virtual directory)",
		},
		"legal_hold": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Object legal hold status (ON, OFF)",
		},
		"lock": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Object lock retention configuration",
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Lock retention mode (COMPLIANCE, GOVERNANCE)",
				},
				"retain_until": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Date until which the object is locked",
				},
			},
		},
		"restore_status": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Restore status for archived objects",
			Attributes: map[string]schema.Attribute{
				"in_progress": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether a restore operation is currently in progress",
				},
				"expires_at": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Expiration date of the restored copy",
				},
			},
		},
		"replication_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Replication status of the object (COMPLETED, FAILED, NONE, PENDING, REPLICA)",
		},
	}
}

func buildStorageObjectBucketObjectLock(lock *CloudStorageObjectBucketObjectAPILock) types.Object {
	if lock == nil {
		return types.ObjectNull(storageObjectBucketObjectLockAttrTypes())
	}
	obj, _ := types.ObjectValue(storageObjectBucketObjectLockAttrTypes(), map[string]attr.Value{
		"mode":         ovhtypes.TfStringValue{StringValue: types.StringValue(lock.Mode)},
		"retain_until": ovhtypes.TfStringValue{StringValue: types.StringValue(lock.RetainUntil)},
	})
	return obj
}

func buildStorageObjectBucketObjectRestoreStatus(status *CloudStorageObjectBucketObjectAPIRestoreStatus) types.Object {
	if status == nil {
		return types.ObjectNull(storageObjectBucketObjectRestoreStatusAttrTypes())
	}
	obj, _ := types.ObjectValue(storageObjectBucketObjectRestoreStatusAttrTypes(), map[string]attr.Value{
		"in_progress": types.BoolValue(status.InProgress),
		"expires_at":  nullableTfString(status.ExpiresAt),
	})
	return obj
}

func buildStorageObjectBucketObjectValue(object *CloudStorageObjectBucketObjectAPIResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(storageObjectBucketObjectAttrTypes(), map[string]attr.Value{
		"key":                ovhtypes.TfStringValue{StringValue: types.StringValue(object.Key)},
		"size":               types.Int64Value(object.Size),
		"etag":               nullableTfString(object.Etag),
		"last_modified_at":   nullableTfString(object.LastModifiedAt),
		"storage_class":      nullableTfString(object.StorageClass),
		"version_id":         nullableTfString(object.VersionId),
		"is_latest":          types.BoolValue(object.IsLatest),
		"is_delete_marker":   types.BoolValue(object.IsDeleteMarker),
		"is_common_prefix":   types.BoolValue(object.IsCommonPrefix),
		"legal_hold":         nullableTfString(object.LegalHold),
		"lock":               buildStorageObjectBucketObjectLock(object.Lock),
		"restore_status":     buildStorageObjectBucketObjectRestoreStatus(object.RestoreStatus),
		"replication_status": nullableTfString(object.ReplicationStatus),
	})
}
