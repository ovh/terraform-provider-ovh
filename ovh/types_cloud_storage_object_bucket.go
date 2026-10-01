package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageObjectBucketModel represents the Terraform model for the S3 bucket resource
type CloudStorageObjectBucketModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — mutable
	OwnerUserId ovhtypes.TfStringValue `tfsdk:"owner_user_id"`
	Tags        types.Map              `tfsdk:"tags"`
	Encryption  types.Object           `tfsdk:"encryption"`
	Versioning  types.Object           `tfsdk:"versioning"`
	ObjectLock  types.Object           `tfsdk:"object_lock"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudStorageObjectBucketAPIResponse struct {
	Id             string                                   `json:"id"`
	Checksum       string                                   `json:"checksum"`
	CreatedAt      string                                   `json:"createdAt"`
	UpdatedAt      string                                   `json:"updatedAt"`
	ResourceStatus string                                   `json:"resourceStatus"`
	TargetSpec     *CloudStorageObjectBucketAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudStorageObjectBucketAPICurrentState `json:"currentState,omitempty"`
	CurrentTasks   []CloudResourceTask                      `json:"currentTasks,omitempty"`
}

type CloudStorageObjectBucketAPILocation struct {
	Region string `json:"region,omitempty"`
}

type CloudStorageObjectBucketAPIEncryption struct {
	Algorithm string `json:"algorithm"`
}

type CloudStorageObjectBucketAPIVersioning struct {
	Status string `json:"status"`
}

type CloudStorageObjectBucketAPIObjectLock struct {
	Mode          string `json:"mode"`
	RetentionDays int64  `json:"retentionDays"`
	// retentionYears is read-only in the API contract: never sent, only read back.
	RetentionYears *int64 `json:"retentionYears,omitempty"`
}

type CloudStorageObjectBucketAPITargetSpec struct {
	Name        string                                 `json:"name,omitempty"`
	Location    *CloudStorageObjectBucketAPILocation   `json:"location,omitempty"`
	OwnerUserId string                                 `json:"ownerUserId,omitempty"`
	Encryption  *CloudStorageObjectBucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudStorageObjectBucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudStorageObjectBucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string                      `json:"tags,omitempty"`
}

type CloudStorageObjectBucketAPIUpdateTargetSpec struct {
	OwnerUserId string                                 `json:"ownerUserId,omitempty"`
	Encryption  *CloudStorageObjectBucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudStorageObjectBucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudStorageObjectBucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string                      `json:"tags,omitempty"`
}

type CloudStorageObjectBucketAPICurrentState struct {
	Name        string                                 `json:"name,omitempty"`
	Location    *CloudStorageObjectBucketAPILocation   `json:"location,omitempty"`
	Encryption  *CloudStorageObjectBucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudStorageObjectBucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudStorageObjectBucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string                      `json:"tags,omitempty"`
	VirtualHost string                                 `json:"virtualHost,omitempty"`
	// Pointers: the API only returns these on a single bucket GET, never on LIST.
	ObjectsCount *int64 `json:"objectsCount,omitempty"`
	ObjectsSize  *int64 `json:"objectsSize,omitempty"`
}

// Create payload
type CloudStorageObjectBucketCreatePayload struct {
	TargetSpec *CloudStorageObjectBucketAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudStorageObjectBucketUpdatePayload struct {
	Checksum   string                                       `json:"checksum"`
	TargetSpec *CloudStorageObjectBucketAPIUpdateTargetSpec `json:"targetSpec"`
}

// StorageObjectBucketLocationAttrTypes returns the attribute types for the location object
func StorageObjectBucketLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region": ovhtypes.TfStringType{},
	}
}

// StorageObjectBucketEncryptionAttrTypes returns the attribute types for the encryption object
func StorageObjectBucketEncryptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"algorithm": ovhtypes.TfStringType{},
	}
}

// StorageObjectBucketVersioningAttrTypes returns the attribute types for the versioning object
func StorageObjectBucketVersioningAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status": ovhtypes.TfStringType{},
	}
}

// StorageObjectBucketObjectLockAttrTypes returns the attribute types for the object_lock object
func StorageObjectBucketObjectLockAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":            ovhtypes.TfStringType{},
		"retention_days":  types.Int64Type,
		"retention_years": types.Int64Type,
	}
}

// retention_years is read-only, so it lives only under current_state.object_lock.
func StorageObjectBucketObjectLockSpecAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":           ovhtypes.TfStringType{},
		"retention_days": types.Int64Type,
	}
}

// StorageObjectBucketCurrentStateAttrTypes returns the attribute types for the current_state object
func StorageObjectBucketCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"location":      types.ObjectType{AttrTypes: StorageObjectBucketLocationAttrTypes()},
		"encryption":    types.ObjectType{AttrTypes: StorageObjectBucketEncryptionAttrTypes()},
		"versioning":    types.ObjectType{AttrTypes: StorageObjectBucketVersioningAttrTypes()},
		"object_lock":   types.ObjectType{AttrTypes: StorageObjectBucketObjectLockAttrTypes()},
		"tags":          types.MapType{ElemType: ovhtypes.TfStringType{}},
		"virtual_host":  ovhtypes.TfStringType{},
		"objects_count": types.Int64Type,
		"objects_size":  types.Int64Type,
	}
}

func nullableInt64(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

func buildStorageObjectBucketLocationObject(location *CloudStorageObjectBucketAPILocation) types.Object {
	if location == nil {
		return types.ObjectNull(StorageObjectBucketLocationAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageObjectBucketLocationAttrTypes(),
		map[string]attr.Value{
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(location.Region)},
		},
	)
	return obj
}

func buildStorageObjectBucketEncryptionObject(encryption *CloudStorageObjectBucketAPIEncryption) types.Object {
	if encryption == nil {
		return types.ObjectNull(StorageObjectBucketEncryptionAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageObjectBucketEncryptionAttrTypes(),
		map[string]attr.Value{
			"algorithm": ovhtypes.TfStringValue{StringValue: types.StringValue(encryption.Algorithm)},
		},
	)
	return obj
}

func buildStorageObjectBucketVersioningObject(versioning *CloudStorageObjectBucketAPIVersioning) types.Object {
	if versioning == nil {
		return types.ObjectNull(StorageObjectBucketVersioningAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageObjectBucketVersioningAttrTypes(),
		map[string]attr.Value{
			"status": ovhtypes.TfStringValue{StringValue: types.StringValue(versioning.Status)},
		},
	)
	return obj
}

func buildStorageObjectBucketObjectLockObject(objectLock *CloudStorageObjectBucketAPIObjectLock) types.Object {
	if objectLock == nil {
		return types.ObjectNull(StorageObjectBucketObjectLockAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageObjectBucketObjectLockAttrTypes(),
		map[string]attr.Value{
			"mode":            ovhtypes.TfStringValue{StringValue: types.StringValue(objectLock.Mode)},
			"retention_days":  types.Int64Value(objectLock.RetentionDays),
			"retention_years": nullableInt64(objectLock.RetentionYears),
		},
	)
	return obj
}

func buildStorageObjectBucketObjectLockSpecObject(objectLock *CloudStorageObjectBucketAPIObjectLock) types.Object {
	if objectLock == nil {
		return types.ObjectNull(StorageObjectBucketObjectLockSpecAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageObjectBucketObjectLockSpecAttrTypes(),
		map[string]attr.Value{
			"mode":           ovhtypes.TfStringValue{StringValue: types.StringValue(objectLock.Mode)},
			"retention_days": types.Int64Value(objectLock.RetentionDays),
		},
	)
	return obj
}

func buildStorageObjectBucketTagsMap(tags map[string]string) types.Map {
	if tags == nil {
		return types.MapNull(ovhtypes.TfStringType{})
	}

	elems := make(map[string]attr.Value, len(tags))
	for k, v := range tags {
		elems[k] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
	}
	m, _ := types.MapValue(ovhtypes.TfStringType{}, elems)
	return m
}

func storageObjectBucketObjectString(attrs map[string]attr.Value, key string) string {
	if v, ok := attrs[key]; ok {
		if sv, ok := v.(ovhtypes.TfStringValue); ok && !sv.IsNull() && !sv.IsUnknown() {
			return sv.ValueString()
		}
	}
	return ""
}

func storageObjectBucketObjectInt64(attrs map[string]attr.Value, key string) int64 {
	if v, ok := attrs[key]; ok {
		if iv, ok := v.(types.Int64); ok && !iv.IsNull() && !iv.IsUnknown() {
			return iv.ValueInt64()
		}
	}
	return 0
}

func storageObjectBucketEncryptionToAPI(o types.Object) *CloudStorageObjectBucketAPIEncryption {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	return &CloudStorageObjectBucketAPIEncryption{Algorithm: storageObjectBucketObjectString(o.Attributes(), "algorithm")}
}

func storageObjectBucketVersioningToAPI(o types.Object) *CloudStorageObjectBucketAPIVersioning {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	return &CloudStorageObjectBucketAPIVersioning{Status: storageObjectBucketObjectString(o.Attributes(), "status")}
}

func storageObjectBucketObjectLockToAPI(o types.Object) *CloudStorageObjectBucketAPIObjectLock {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	attrs := o.Attributes()
	return &CloudStorageObjectBucketAPIObjectLock{
		Mode:          storageObjectBucketObjectString(attrs, "mode"),
		RetentionDays: storageObjectBucketObjectInt64(attrs, "retention_days"),
	}
}

func storageObjectBucketTagsToAPI(m types.Map) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	tags := make(map[string]string, len(m.Elements()))
	for k, v := range m.Elements() {
		if sv, ok := v.(ovhtypes.TfStringValue); ok && !sv.IsNull() && !sv.IsUnknown() {
			tags[k] = sv.ValueString()
		}
	}
	return tags
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageObjectBucketModel) ToCreate() *CloudStorageObjectBucketCreatePayload {
	target := &CloudStorageObjectBucketAPITargetSpec{
		Name:       m.Name.ValueString(),
		Location:   &CloudStorageObjectBucketAPILocation{Region: m.Region.ValueString()},
		Encryption: storageObjectBucketEncryptionToAPI(m.Encryption),
		Versioning: storageObjectBucketVersioningToAPI(m.Versioning),
		ObjectLock: storageObjectBucketObjectLockToAPI(m.ObjectLock),
		Tags:       storageObjectBucketTagsToAPI(m.Tags),
	}

	if !m.OwnerUserId.IsNull() && !m.OwnerUserId.IsUnknown() {
		target.OwnerUserId = m.OwnerUserId.ValueString()
	}

	return &CloudStorageObjectBucketCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageObjectBucketModel) ToUpdate(checksum string) *CloudStorageObjectBucketUpdatePayload {
	target := &CloudStorageObjectBucketAPIUpdateTargetSpec{
		Encryption: storageObjectBucketEncryptionToAPI(m.Encryption),
		Versioning: storageObjectBucketVersioningToAPI(m.Versioning),
		ObjectLock: storageObjectBucketObjectLockToAPI(m.ObjectLock),
		Tags:       storageObjectBucketTagsToAPI(m.Tags),
	}

	if !m.OwnerUserId.IsNull() && !m.OwnerUserId.IsUnknown() {
		target.OwnerUserId = m.OwnerUserId.ValueString()
	}

	return &CloudStorageObjectBucketUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// buildStorageObjectBucketCurrentStateObject constructs the current_state object from the API response
func buildStorageObjectBucketCurrentStateObject(ctx context.Context, state *CloudStorageObjectBucketAPICurrentState) types.Object {
	obj, _ := types.ObjectValue(
		StorageObjectBucketCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"location":      buildStorageObjectBucketLocationObject(state.Location),
			"encryption":    buildStorageObjectBucketEncryptionObject(state.Encryption),
			"versioning":    buildStorageObjectBucketVersioningObject(state.Versioning),
			"object_lock":   buildStorageObjectBucketObjectLockObject(state.ObjectLock),
			"tags":          buildStorageObjectBucketTagsMap(state.Tags),
			"virtual_host":  nullableTfString(state.VirtualHost),
			"objects_count": nullableInt64(state.ObjectsCount),
			"objects_size":  nullableInt64(state.ObjectsSize),
		},
	)
	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageObjectBucketModel) MergeWith(ctx context.Context, response *CloudStorageObjectBucketAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildStorageObjectBucketCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(StorageObjectBucketCurrentStateAttrTypes())
	}

	if response.TargetSpec == nil {
		return
	}

	// Immutable and config-owned: the API normalizes them (region upper-cased), so writing
	// them back breaks plan consistency. Only fill when unset (import).
	if (m.Name.IsNull() || m.Name.IsUnknown()) && response.TargetSpec.Name != "" {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
	}
	if (m.Region.IsNull() || m.Region.IsUnknown()) && response.TargetSpec.Location != nil && response.TargetSpec.Location.Region != "" {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
	}

	m.OwnerUserId = nullableTfString(response.TargetSpec.OwnerUserId)
	m.Encryption = buildStorageObjectBucketEncryptionObject(response.TargetSpec.Encryption)
	m.Versioning = buildStorageObjectBucketVersioningObject(response.TargetSpec.Versioning)
	m.ObjectLock = buildStorageObjectBucketObjectLockSpecObject(response.TargetSpec.ObjectLock)

	// An empty tags map is dropped by the API (omitempty), so keep a configured
	// empty map instead of flipping it to null and breaking plan consistency.
	if response.TargetSpec.Tags == nil && !m.Tags.IsNull() && !m.Tags.IsUnknown() && len(m.Tags.Elements()) == 0 {
		return
	}
	m.Tags = buildStorageObjectBucketTagsMap(response.TargetSpec.Tags)
}
