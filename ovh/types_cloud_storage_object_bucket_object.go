package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// --- API payloads shared between the object-lock resource and the object data sources ---

// CloudBucketObjectRetention mirrors publicCloud.storage.object.ObjectRetention.
type CloudBucketObjectRetention struct {
	Mode        string `json:"mode"`
	RetainUntil string `json:"retainUntil"`
}

// CloudBucketObjectRestoreStatus mirrors publicCloud.storage.object.ObjectRestoreStatus.
type CloudBucketObjectRestoreStatus struct {
	ExpireDate string `json:"expireDate,omitempty"`
	InProgress bool   `json:"inProgress"`
}

// CloudBucketObjectAPI mirrors publicCloud.storage.object.Object.
type CloudBucketObjectAPI struct {
	ETag              string                          `json:"eTag,omitempty"`
	IsCommonPrefix    bool                            `json:"isCommonPrefix,omitempty"`
	IsDeleteMarker    bool                            `json:"isDeleteMarker,omitempty"`
	IsLatest          bool                            `json:"isLatest,omitempty"`
	Key               string                          `json:"key"`
	LastModified      string                          `json:"lastModified,omitempty"`
	LegalHold         string                          `json:"legalHold,omitempty"`
	ReplicationStatus string                          `json:"replicationStatus,omitempty"`
	RestoreStatus     *CloudBucketObjectRestoreStatus `json:"restoreStatus,omitempty"`
	Retention         *CloudBucketObjectRetention     `json:"retention,omitempty"`
	Size              int64                           `json:"size,omitempty"`
	StorageClass      string                          `json:"storageClass,omitempty"`
	VersionID         string                          `json:"versionId,omitempty"`
}

// CloudBucketObjectUpdatePayload mirrors publicCloud.storage.object.ObjectUpdate.
type CloudBucketObjectUpdatePayload struct {
	LegalHold string                      `json:"legalHold,omitempty"`
	Retention *CloudBucketObjectRetention `json:"retention,omitempty"`
}

// CloudBucketObjectListResponse mirrors publicCloud.storage.object.ObjectListResponse.
type CloudBucketObjectListResponse struct {
	IsTruncated   bool                   `json:"isTruncated"`
	NextKeyMarker *string                `json:"nextKeyMarker,omitempty"`
	Objects       []CloudBucketObjectAPI `json:"objects"`
}

// CloudBucketObjectVersionListResponse mirrors publicCloud.storage.object.ObjectVersionListResponse.
type CloudBucketObjectVersionListResponse struct {
	IsTruncated         bool                   `json:"isTruncated"`
	NextKeyMarker       *string                `json:"nextKeyMarker,omitempty"`
	NextVersionIDMarker *string                `json:"nextVersionIdMarker,omitempty"`
	Objects             []CloudBucketObjectAPI `json:"objects"`
}

// --- Attribute type helpers ---

func cloudBucketObjectRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":              ovhtypes.TfStringType{},
		"retain_until_date": ovhtypes.TfStringType{},
	}
}

func cloudBucketObjectLegalHoldAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status": ovhtypes.TfStringType{},
	}
}

func cloudBucketObjectRestoreStatusAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"expire_date": ovhtypes.TfStringType{},
		"in_progress": types.BoolType,
	}
}

// cloudBucketObjectAttrTypes returns the map[string]attr.Type describing a single
// object as exposed by the data sources. It mirrors CloudBucketObjectAPI but in
// snake_case.
func cloudBucketObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"e_tag":              ovhtypes.TfStringType{},
		"is_common_prefix":   types.BoolType,
		"is_delete_marker":   types.BoolType,
		"is_latest":          types.BoolType,
		"key":                ovhtypes.TfStringType{},
		"last_modified":      ovhtypes.TfStringType{},
		"legal_hold":         ovhtypes.TfStringType{},
		"replication_status": ovhtypes.TfStringType{},
		"restore_status":     types.ObjectType{AttrTypes: cloudBucketObjectRestoreStatusAttrTypes()},
		"retention":          types.ObjectType{AttrTypes: cloudBucketObjectRetentionAttrTypes()},
		"size":               types.Int64Type,
		"storage_class":      ovhtypes.TfStringType{},
		"version_id":         ovhtypes.TfStringType{},
	}
}

// --- Conversion helpers ---

// apiObjectRetentionToTFObject converts a CloudBucketObjectRetention into a
// Terraform object matching cloudBucketObjectRetentionAttrTypes. Returns null when
// the input is nil.
func apiObjectRetentionToTFObject(r *CloudBucketObjectRetention) basetypes.ObjectValue {
	if r == nil {
		return types.ObjectNull(cloudBucketObjectRetentionAttrTypes())
	}
	obj, _ := types.ObjectValue(cloudBucketObjectRetentionAttrTypes(), map[string]attr.Value{
		"mode":              ovhtypes.TfStringValue{StringValue: types.StringValue(r.Mode)},
		"retain_until_date": ovhtypes.TfStringValue{StringValue: types.StringValue(r.RetainUntil)},
	})
	return obj
}

// apiLegalHoldStringToTFObject converts a legal-hold status string to a Terraform
// object. Empty string becomes null.
func apiLegalHoldStringToTFObject(status string) basetypes.ObjectValue {
	if status == "" {
		return types.ObjectNull(cloudBucketObjectLegalHoldAttrTypes())
	}
	obj, _ := types.ObjectValue(cloudBucketObjectLegalHoldAttrTypes(), map[string]attr.Value{
		"status": ovhtypes.TfStringValue{StringValue: types.StringValue(status)},
	})
	return obj
}

// apiRestoreStatusToTFObject converts a CloudBucketObjectRestoreStatus to a
// Terraform object. Null when input is nil.
func apiRestoreStatusToTFObject(r *CloudBucketObjectRestoreStatus) basetypes.ObjectValue {
	if r == nil {
		return types.ObjectNull(cloudBucketObjectRestoreStatusAttrTypes())
	}
	expireDateVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if r.ExpireDate != "" {
		expireDateVal = ovhtypes.TfStringValue{StringValue: types.StringValue(r.ExpireDate)}
	}
	obj, _ := types.ObjectValue(cloudBucketObjectRestoreStatusAttrTypes(), map[string]attr.Value{
		"expire_date": expireDateVal,
		"in_progress": types.BoolValue(r.InProgress),
	})
	return obj
}

// apiObjectToTFObject converts a CloudBucketObjectAPI to a Terraform object
// matching cloudBucketObjectAttrTypes.
func apiObjectToTFObject(o CloudBucketObjectAPI) basetypes.ObjectValue {
	eTagVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.ETag != "" {
		eTagVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.ETag)}
	}
	lastModifiedVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.LastModified != "" {
		lastModifiedVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.LastModified)}
	}
	legalHoldVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.LegalHold != "" {
		legalHoldVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.LegalHold)}
	}
	replicationStatusVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.ReplicationStatus != "" {
		replicationStatusVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.ReplicationStatus)}
	}
	storageClassVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.StorageClass != "" {
		storageClassVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.StorageClass)}
	}
	versionIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if o.VersionID != "" {
		versionIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(o.VersionID)}
	}

	obj, _ := types.ObjectValue(cloudBucketObjectAttrTypes(), map[string]attr.Value{
		"e_tag":              eTagVal,
		"is_common_prefix":   types.BoolValue(o.IsCommonPrefix),
		"is_delete_marker":   types.BoolValue(o.IsDeleteMarker),
		"is_latest":          types.BoolValue(o.IsLatest),
		"key":                ovhtypes.TfStringValue{StringValue: types.StringValue(o.Key)},
		"last_modified":      lastModifiedVal,
		"legal_hold":         legalHoldVal,
		"replication_status": replicationStatusVal,
		"restore_status":     apiRestoreStatusToTFObject(o.RestoreStatus),
		"retention":          apiObjectRetentionToTFObject(o.Retention),
		"size":               types.Int64Value(o.Size),
		"storage_class":      storageClassVal,
		"version_id":         versionIdVal,
	})
	return obj
}

// apiObjectListToTFList converts a slice of CloudBucketObjectAPI to a Terraform
// list of objects matching cloudBucketObjectAttrTypes.
func apiObjectListToTFList(list []CloudBucketObjectAPI) types.List {
	elemType := types.ObjectType{AttrTypes: cloudBucketObjectAttrTypes()}
	elems := make([]attr.Value, 0, len(list))
	for _, o := range list {
		elems = append(elems, apiObjectToTFObject(o))
	}
	out, _ := types.ListValue(elemType, elems)
	return out
}
