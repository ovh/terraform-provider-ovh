package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudS3BucketModel represents the Terraform model for the S3 bucket resource
type CloudS3BucketModel struct {
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
type CloudS3BucketAPIResponse struct {
	Id             string                        `json:"id"`
	Checksum       string                        `json:"checksum"`
	CreatedAt      string                        `json:"createdAt"`
	UpdatedAt      string                        `json:"updatedAt"`
	ResourceStatus string                        `json:"resourceStatus"`
	TargetSpec     *CloudS3BucketAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudS3BucketAPICurrentState `json:"currentState,omitempty"`
	CurrentTasks   []CloudResourceTask           `json:"currentTasks,omitempty"`
}

type CloudS3BucketAPILocation struct {
	Region string `json:"region,omitempty"`
}

type CloudS3BucketAPIEncryption struct {
	Algorithm string `json:"algorithm"`
}

type CloudS3BucketAPIVersioning struct {
	Status string `json:"status"`
}

type CloudS3BucketAPIObjectLock struct {
	Mode          string `json:"mode"`
	RetentionDays int64  `json:"retentionDays"`
	// retentionYears is read-only in the API contract: never sent, only read back.
	RetentionYears *int64 `json:"retentionYears,omitempty"`
}

type CloudS3BucketAPITargetSpec struct {
	Name        string                      `json:"name,omitempty"`
	Location    *CloudS3BucketAPILocation   `json:"location,omitempty"`
	OwnerUserId string                      `json:"ownerUserId,omitempty"`
	Encryption  *CloudS3BucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudS3BucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudS3BucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string           `json:"tags,omitempty"`
}

type CloudS3BucketAPIUpdateTargetSpec struct {
	OwnerUserId string                      `json:"ownerUserId,omitempty"`
	Encryption  *CloudS3BucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudS3BucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudS3BucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string           `json:"tags,omitempty"`
}

type CloudS3BucketAPICurrentState struct {
	Name        string                      `json:"name,omitempty"`
	Location    *CloudS3BucketAPILocation   `json:"location,omitempty"`
	Encryption  *CloudS3BucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudS3BucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudS3BucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string           `json:"tags,omitempty"`
	VirtualHost string                      `json:"virtualHost,omitempty"`
	// Pointers: the API only returns these on a single bucket GET, never on LIST.
	ObjectsCount *int64 `json:"objectsCount,omitempty"`
	ObjectsSize  *int64 `json:"objectsSize,omitempty"`
}

// Create payload
type CloudS3BucketCreatePayload struct {
	TargetSpec *CloudS3BucketAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudS3BucketUpdatePayload struct {
	Checksum   string                            `json:"checksum"`
	TargetSpec *CloudS3BucketAPIUpdateTargetSpec `json:"targetSpec"`
}

// S3BucketLocationAttrTypes returns the attribute types for the location object
func S3BucketLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region": ovhtypes.TfStringType{},
	}
}

// S3BucketEncryptionAttrTypes returns the attribute types for the encryption object
func S3BucketEncryptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"algorithm": ovhtypes.TfStringType{},
	}
}

// S3BucketVersioningAttrTypes returns the attribute types for the versioning object
func S3BucketVersioningAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status": ovhtypes.TfStringType{},
	}
}

// S3BucketObjectLockAttrTypes returns the attribute types for the object_lock object
func S3BucketObjectLockAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":            ovhtypes.TfStringType{},
		"retention_days":  types.Int64Type,
		"retention_years": types.Int64Type,
	}
}

// retention_years is read-only, so it lives only under current_state.object_lock.
func S3BucketObjectLockSpecAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":           ovhtypes.TfStringType{},
		"retention_days": types.Int64Type,
	}
}

// S3BucketCurrentStateAttrTypes returns the attribute types for the current_state object
func S3BucketCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"location":      types.ObjectType{AttrTypes: S3BucketLocationAttrTypes()},
		"encryption":    types.ObjectType{AttrTypes: S3BucketEncryptionAttrTypes()},
		"versioning":    types.ObjectType{AttrTypes: S3BucketVersioningAttrTypes()},
		"object_lock":   types.ObjectType{AttrTypes: S3BucketObjectLockAttrTypes()},
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

func buildS3BucketLocationObject(location *CloudS3BucketAPILocation) types.Object {
	if location == nil {
		return types.ObjectNull(S3BucketLocationAttrTypes())
	}

	obj, _ := types.ObjectValue(
		S3BucketLocationAttrTypes(),
		map[string]attr.Value{
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(location.Region)},
		},
	)
	return obj
}

func buildS3BucketEncryptionObject(encryption *CloudS3BucketAPIEncryption) types.Object {
	if encryption == nil {
		return types.ObjectNull(S3BucketEncryptionAttrTypes())
	}

	obj, _ := types.ObjectValue(
		S3BucketEncryptionAttrTypes(),
		map[string]attr.Value{
			"algorithm": ovhtypes.TfStringValue{StringValue: types.StringValue(encryption.Algorithm)},
		},
	)
	return obj
}

func buildS3BucketVersioningObject(versioning *CloudS3BucketAPIVersioning) types.Object {
	if versioning == nil {
		return types.ObjectNull(S3BucketVersioningAttrTypes())
	}

	obj, _ := types.ObjectValue(
		S3BucketVersioningAttrTypes(),
		map[string]attr.Value{
			"status": ovhtypes.TfStringValue{StringValue: types.StringValue(versioning.Status)},
		},
	)
	return obj
}

func buildS3BucketObjectLockObject(objectLock *CloudS3BucketAPIObjectLock) types.Object {
	if objectLock == nil {
		return types.ObjectNull(S3BucketObjectLockAttrTypes())
	}

	obj, _ := types.ObjectValue(
		S3BucketObjectLockAttrTypes(),
		map[string]attr.Value{
			"mode":            ovhtypes.TfStringValue{StringValue: types.StringValue(objectLock.Mode)},
			"retention_days":  types.Int64Value(objectLock.RetentionDays),
			"retention_years": nullableInt64(objectLock.RetentionYears),
		},
	)
	return obj
}

func buildS3BucketObjectLockSpecObject(objectLock *CloudS3BucketAPIObjectLock) types.Object {
	if objectLock == nil {
		return types.ObjectNull(S3BucketObjectLockSpecAttrTypes())
	}

	obj, _ := types.ObjectValue(
		S3BucketObjectLockSpecAttrTypes(),
		map[string]attr.Value{
			"mode":           ovhtypes.TfStringValue{StringValue: types.StringValue(objectLock.Mode)},
			"retention_days": types.Int64Value(objectLock.RetentionDays),
		},
	)
	return obj
}

func buildS3BucketTagsMap(tags map[string]string) types.Map {
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

func s3BucketObjectString(attrs map[string]attr.Value, key string) string {
	if v, ok := attrs[key]; ok {
		if sv, ok := v.(ovhtypes.TfStringValue); ok && !sv.IsNull() && !sv.IsUnknown() {
			return sv.ValueString()
		}
	}
	return ""
}

func s3BucketObjectInt64(attrs map[string]attr.Value, key string) int64 {
	if v, ok := attrs[key]; ok {
		if iv, ok := v.(types.Int64); ok && !iv.IsNull() && !iv.IsUnknown() {
			return iv.ValueInt64()
		}
	}
	return 0
}

func s3BucketEncryptionToAPI(o types.Object) *CloudS3BucketAPIEncryption {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	return &CloudS3BucketAPIEncryption{Algorithm: s3BucketObjectString(o.Attributes(), "algorithm")}
}

func s3BucketVersioningToAPI(o types.Object) *CloudS3BucketAPIVersioning {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	return &CloudS3BucketAPIVersioning{Status: s3BucketObjectString(o.Attributes(), "status")}
}

func s3BucketObjectLockToAPI(o types.Object) *CloudS3BucketAPIObjectLock {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	attrs := o.Attributes()
	return &CloudS3BucketAPIObjectLock{
		Mode:          s3BucketObjectString(attrs, "mode"),
		RetentionDays: s3BucketObjectInt64(attrs, "retention_days"),
	}
}

func s3BucketTagsToAPI(m types.Map) map[string]string {
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
func (m *CloudS3BucketModel) ToCreate() *CloudS3BucketCreatePayload {
	target := &CloudS3BucketAPITargetSpec{
		Name:       m.Name.ValueString(),
		Location:   &CloudS3BucketAPILocation{Region: m.Region.ValueString()},
		Encryption: s3BucketEncryptionToAPI(m.Encryption),
		Versioning: s3BucketVersioningToAPI(m.Versioning),
		ObjectLock: s3BucketObjectLockToAPI(m.ObjectLock),
		Tags:       s3BucketTagsToAPI(m.Tags),
	}

	if !m.OwnerUserId.IsNull() && !m.OwnerUserId.IsUnknown() {
		target.OwnerUserId = m.OwnerUserId.ValueString()
	}

	return &CloudS3BucketCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudS3BucketModel) ToUpdate(checksum string) *CloudS3BucketUpdatePayload {
	target := &CloudS3BucketAPIUpdateTargetSpec{
		Encryption: s3BucketEncryptionToAPI(m.Encryption),
		Versioning: s3BucketVersioningToAPI(m.Versioning),
		ObjectLock: s3BucketObjectLockToAPI(m.ObjectLock),
		Tags:       s3BucketTagsToAPI(m.Tags),
	}

	if !m.OwnerUserId.IsNull() && !m.OwnerUserId.IsUnknown() {
		target.OwnerUserId = m.OwnerUserId.ValueString()
	}

	return &CloudS3BucketUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// buildS3BucketCurrentStateObject constructs the current_state object from the API response
func buildS3BucketCurrentStateObject(ctx context.Context, state *CloudS3BucketAPICurrentState) types.Object {
	obj, _ := types.ObjectValue(
		S3BucketCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"location":      buildS3BucketLocationObject(state.Location),
			"encryption":    buildS3BucketEncryptionObject(state.Encryption),
			"versioning":    buildS3BucketVersioningObject(state.Versioning),
			"object_lock":   buildS3BucketObjectLockObject(state.ObjectLock),
			"tags":          buildS3BucketTagsMap(state.Tags),
			"virtual_host":  nullableTfString(state.VirtualHost),
			"objects_count": nullableInt64(state.ObjectsCount),
			"objects_size":  nullableInt64(state.ObjectsSize),
		},
	)
	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudS3BucketModel) MergeWith(ctx context.Context, response *CloudS3BucketAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildS3BucketCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(S3BucketCurrentStateAttrTypes())
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
	m.Encryption = buildS3BucketEncryptionObject(response.TargetSpec.Encryption)
	m.Versioning = buildS3BucketVersioningObject(response.TargetSpec.Versioning)
	m.ObjectLock = buildS3BucketObjectLockSpecObject(response.TargetSpec.ObjectLock)

	// An empty tags map is dropped by the API (omitempty), so keep a configured
	// empty map instead of flipping it to null and breaking plan consistency.
	if response.TargetSpec.Tags == nil && !m.Tags.IsNull() && !m.Tags.IsUnknown() && len(m.Tags.Elements()) == 0 {
		return
	}
	m.Tags = buildS3BucketTagsMap(response.TargetSpec.Tags)
}
