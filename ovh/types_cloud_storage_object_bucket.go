package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageObjectBucketModel represents the Terraform model for the S3 bucket resource
type CloudStorageObjectBucketModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional
	Encryption  types.Object           `tfsdk:"encryption"`
	Versioning  types.Object           `tfsdk:"versioning"`
	ObjectLock  types.Object           `tfsdk:"object_lock"`
	Tags        types.Map              `tfsdk:"tags"`
	OwnerUserId ovhtypes.TfStringValue `tfsdk:"owner_user_id"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API Response types
type CloudBucketAPIResponse struct {
	Id             string                      `json:"id"`
	Checksum       string                      `json:"checksum"`
	CreatedAt      string                      `json:"createdAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	ResourceStatus string                      `json:"resourceStatus"`
	CurrentState   *CloudBucketAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudBucketAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudBucketAPICurrentState struct {
	Name       string                    `json:"name,omitempty"`
	Location   *CloudBucketAPILocation   `json:"location,omitempty"`
	Encryption *CloudBucketAPIEncryption `json:"encryption,omitempty"`
	Versioning *CloudBucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock *CloudBucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags       map[string]string         `json:"tags,omitempty"`
}

type CloudBucketAPILocation struct {
	Region string `json:"region"`
}

type CloudBucketAPIEncryption struct {
	Algorithm string `json:"algorithm"`
}

type CloudBucketAPIVersioning struct {
	Status string `json:"status"`
}

type CloudBucketAPIObjectLock struct {
	Mode           string `json:"mode"`
	RetentionDays  int64  `json:"retentionDays"`
	RetentionYears *int64 `json:"retentionYears,omitempty"`
}

type CloudBucketAPITargetSpec struct {
	Name        string                    `json:"name"`
	Location    *CloudBucketAPILocation   `json:"location,omitempty"`
	Encryption  *CloudBucketAPIEncryption `json:"encryption,omitempty"`
	Versioning  *CloudBucketAPIVersioning `json:"versioning,omitempty"`
	ObjectLock  *CloudBucketAPIObjectLock `json:"objectLock,omitempty"`
	Tags        map[string]string         `json:"tags,omitempty"`
	OwnerUserId string                    `json:"ownerUserId,omitempty"`
}

// Create payload
type CloudBucketCreatePayload struct {
	TargetSpec *CloudBucketAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudBucketUpdatePayload struct {
	Checksum   string                    `json:"checksum"`
	TargetSpec *CloudBucketAPITargetSpec `json:"targetSpec"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageObjectBucketModel) ToCreate(ctx context.Context) *CloudBucketCreatePayload {
	targetSpec := m.buildTargetSpec(ctx)

	return &CloudBucketCreatePayload{
		TargetSpec: targetSpec,
	}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageObjectBucketModel) ToUpdate(ctx context.Context, checksum string) *CloudBucketUpdatePayload {
	targetSpec := m.buildTargetSpec(ctx)
	targetSpec.Location = nil

	return &CloudBucketUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

func (m *CloudStorageObjectBucketModel) buildTargetSpec(ctx context.Context) *CloudBucketAPITargetSpec {
	targetSpec := &CloudBucketAPITargetSpec{
		Name: m.Name.ValueString(),
		Location: &CloudBucketAPILocation{
			Region: m.Region.ValueString(),
		},
	}

	if !m.OwnerUserId.IsNull() && !m.OwnerUserId.IsUnknown() {
		targetSpec.OwnerUserId = m.OwnerUserId.ValueString()
	}

	// Encryption
	if !m.Encryption.IsNull() && !m.Encryption.IsUnknown() {
		attrs := m.Encryption.Attributes()
		if v, ok := attrs["algorithm"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			targetSpec.Encryption = &CloudBucketAPIEncryption{
				Algorithm: v.ValueString(),
			}
		}
	}

	// Versioning
	if !m.Versioning.IsNull() && !m.Versioning.IsUnknown() {
		attrs := m.Versioning.Attributes()
		if v, ok := attrs["status"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			targetSpec.Versioning = &CloudBucketAPIVersioning{
				Status: v.ValueString(),
			}
		}
	}

	// ObjectLock
	if !m.ObjectLock.IsNull() && !m.ObjectLock.IsUnknown() {
		attrs := m.ObjectLock.Attributes()
		objLock := &CloudBucketAPIObjectLock{}

		if v, ok := attrs["mode"].(ovhtypes.TfStringValue); ok && !v.IsNull() && !v.IsUnknown() {
			objLock.Mode = v.ValueString()
		}
		if v, ok := attrs["retention_days"].(basetypes.Int64Value); ok && !v.IsNull() && !v.IsUnknown() {
			objLock.RetentionDays = v.ValueInt64()
		}
		if v, ok := attrs["retention_years"].(basetypes.Int64Value); ok && !v.IsNull() && !v.IsUnknown() {
			val := v.ValueInt64()
			objLock.RetentionYears = &val
		}

		targetSpec.ObjectLock = objLock
	}

	// Tags
	if !m.Tags.IsNull() && !m.Tags.IsUnknown() {
		tags := make(map[string]string)
		for k, v := range m.Tags.Elements() {
			if strVal, ok := v.(types.String); ok {
				tags[k] = strVal.ValueString()
			}
		}
		targetSpec.Tags = tags
	}

	return targetSpec
}

// BucketCurrentStateAttrTypes returns the attribute types for the current_state object
func BucketCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"region": ovhtypes.TfStringType{},
			},
		},
		"encryption": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"algorithm": ovhtypes.TfStringType{},
			},
		},
		"versioning": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"status": ovhtypes.TfStringType{},
			},
		},
		"object_lock": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"mode":            ovhtypes.TfStringType{},
				"retention_days":  types.Int64Type,
				"retention_years": types.Int64Type,
			},
		},
		"tags": types.MapType{
			ElemType: types.StringType,
		},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageObjectBucketModel) MergeWith(ctx context.Context, response *CloudBucketAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildBucketCurrentStateObject(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(BucketCurrentStateAttrTypes())
	}

	// Update region from targetSpec if available
	if response.TargetSpec != nil && response.TargetSpec.Location != nil {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
	}
}

func buildBucketCurrentStateObject(state *CloudBucketAPICurrentState) basetypes.ObjectValue {
	// Build location object
	var locationObj basetypes.ObjectValue
	if state.Location != nil {
		locationObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"region": ovhtypes.TfStringType{},
			},
			map[string]attr.Value{
				"region": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Location.Region)},
			},
		)
	} else {
		locationObj = types.ObjectNull(map[string]attr.Type{
			"region": ovhtypes.TfStringType{},
		})
	}

	// Build encryption object
	encryptionAttrTypes := map[string]attr.Type{
		"algorithm": ovhtypes.TfStringType{},
	}
	var encryptionObj basetypes.ObjectValue
	if state.Encryption != nil {
		encryptionObj, _ = types.ObjectValue(
			encryptionAttrTypes,
			map[string]attr.Value{
				"algorithm": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Encryption.Algorithm)},
			},
		)
	} else {
		encryptionObj = types.ObjectNull(encryptionAttrTypes)
	}

	// Build versioning object
	versioningAttrTypes := map[string]attr.Type{
		"status": ovhtypes.TfStringType{},
	}
	var versioningObj basetypes.ObjectValue
	if state.Versioning != nil {
		versioningObj, _ = types.ObjectValue(
			versioningAttrTypes,
			map[string]attr.Value{
				"status": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Versioning.Status)},
			},
		)
	} else {
		versioningObj = types.ObjectNull(versioningAttrTypes)
	}

	// Build object_lock object
	objectLockAttrTypes := map[string]attr.Type{
		"mode":            ovhtypes.TfStringType{},
		"retention_days":  types.Int64Type,
		"retention_years": types.Int64Type,
	}
	var objectLockObj basetypes.ObjectValue
	if state.ObjectLock != nil {
		var retentionYearsVal attr.Value
		if state.ObjectLock.RetentionYears != nil {
			retentionYearsVal = types.Int64Value(*state.ObjectLock.RetentionYears)
		} else {
			retentionYearsVal = types.Int64Null()
		}

		objectLockObj, _ = types.ObjectValue(
			objectLockAttrTypes,
			map[string]attr.Value{
				"mode":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.ObjectLock.Mode)},
				"retention_days":  types.Int64Value(state.ObjectLock.RetentionDays),
				"retention_years": retentionYearsVal,
			},
		)
	} else {
		objectLockObj = types.ObjectNull(objectLockAttrTypes)
	}

	// Build tags map
	var tagsVal basetypes.MapValue
	if state.Tags != nil {
		tagVals := make(map[string]attr.Value, len(state.Tags))
		for k, v := range state.Tags {
			tagVals[k] = types.StringValue(v)
		}
		tagsVal, _ = types.MapValue(types.StringType, tagVals)
	} else {
		tagsVal = types.MapNull(types.StringType)
	}

	currentStateObj, _ := types.ObjectValue(
		BucketCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"location":    locationObj,
			"encryption":  encryptionObj,
			"versioning":  versioningObj,
			"object_lock": objectLockObj,
			"tags":        tagsVal,
		},
	)

	return currentStateObj
}
