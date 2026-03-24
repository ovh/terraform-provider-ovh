package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// ===========================
// Secret types
// ===========================

// CloudKeymanagerSecretModel represents the Terraform model for the KMS secret resource
type CloudKeymanagerSecretModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	SecretType  ovhtypes.TfStringValue `tfsdk:"secret_type"`

	// Optional
	Algorithm          ovhtypes.TfStringValue `tfsdk:"algorithm"`
	BitLength          types.Int64            `tfsdk:"bit_length"`
	Mode               ovhtypes.TfStringValue `tfsdk:"mode"`
	Payload            ovhtypes.TfStringValue `tfsdk:"payload"`
	PayloadContentType ovhtypes.TfStringValue `tfsdk:"payload_content_type"`
	Expiration         ovhtypes.TfStringValue `tfsdk:"expiration"`
	Metadata           types.Map              `tfsdk:"metadata"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types for Secret
type CloudKeymanagerSecretAPIResponse struct {
	Id             string                            `json:"id"`
	Checksum       string                            `json:"checksum"`
	CreatedAt      string                            `json:"createdAt"`
	UpdatedAt      string                            `json:"updatedAt"`
	ResourceStatus string                            `json:"resourceStatus"`
	CurrentState   *CloudKeymanagerSecretCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudKeymanagerSecretTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudKeymanagerSecretCurrentState struct {
	Name                   string                   `json:"name,omitempty"`
	SecretType             string                   `json:"secretType,omitempty"`
	Algorithm              string                   `json:"algorithm,omitempty"`
	BitLength              *int64                   `json:"bitLength,omitempty"`
	Mode                   string                   `json:"mode,omitempty"`
	PayloadContentType     string                   `json:"payloadContentType,omitempty"`
	Expiration             string                   `json:"expiration,omitempty"`
	SecretRef              string                   `json:"secretRef,omitempty"`
	Status                 string                   `json:"status,omitempty"`
	Location               *CloudKeymanagerLocation `json:"location,omitempty"`
	Metadata               map[string]string        `json:"metadata,omitempty"`
}

type CloudKeymanagerLocation struct {
	Region string `json:"region,omitempty"`
}

type CloudKeymanagerSecretTargetSpec struct {
	Name               string                   `json:"name,omitempty"`
	SecretType         string                   `json:"secretType,omitempty"`
	Algorithm          string                   `json:"algorithm,omitempty"`
	BitLength          *int64                   `json:"bitLength,omitempty"`
	Mode               string                   `json:"mode,omitempty"`
	Payload            string                   `json:"payload,omitempty"`
	PayloadContentType string                   `json:"payloadContentType,omitempty"`
	Expiration         string                   `json:"expiration,omitempty"`
	Location           *CloudKeymanagerLocation `json:"location,omitempty"`
	Metadata           map[string]string        `json:"metadata,omitempty"`
}

// Create payload
type CloudKeymanagerSecretCreatePayload struct {
	TargetSpec *CloudKeymanagerSecretTargetSpec `json:"targetSpec"`
}

// Update payload
type CloudKeymanagerSecretUpdatePayload struct {
	Checksum   string                                `json:"checksum"`
	TargetSpec *CloudKeymanagerSecretUpdateTargetSpec `json:"targetSpec"`
}

type CloudKeymanagerSecretUpdateTargetSpec struct {
	Metadata map[string]string `json:"metadata"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudKeymanagerSecretModel) ToCreate(ctx context.Context) *CloudKeymanagerSecretCreatePayload {
	targetSpec := &CloudKeymanagerSecretTargetSpec{
		Name:       m.Name.ValueString(),
		SecretType: m.SecretType.ValueString(),
		Location:   &CloudKeymanagerLocation{Region: m.Region.ValueString()},
	}

	if !m.Algorithm.IsNull() && !m.Algorithm.IsUnknown() {
		targetSpec.Algorithm = m.Algorithm.ValueString()
	}
	if !m.BitLength.IsNull() && !m.BitLength.IsUnknown() {
		v := m.BitLength.ValueInt64()
		targetSpec.BitLength = &v
	}
	if !m.Mode.IsNull() && !m.Mode.IsUnknown() {
		targetSpec.Mode = m.Mode.ValueString()
	}
	if !m.Payload.IsNull() && !m.Payload.IsUnknown() {
		targetSpec.Payload = m.Payload.ValueString()
	}
	if !m.PayloadContentType.IsNull() && !m.PayloadContentType.IsUnknown() {
		targetSpec.PayloadContentType = m.PayloadContentType.ValueString()
	}
	if !m.Expiration.IsNull() && !m.Expiration.IsUnknown() {
		targetSpec.Expiration = m.Expiration.ValueString()
	}
	if !m.Metadata.IsNull() && !m.Metadata.IsUnknown() {
		metadata := map[string]string{}
		m.Metadata.ElementsAs(ctx, &metadata, false)
		targetSpec.Metadata = metadata
	}

	return &CloudKeymanagerSecretCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudKeymanagerSecretModel) ToUpdate(ctx context.Context) *CloudKeymanagerSecretUpdatePayload {
	metadata := map[string]string{}
	if !m.Metadata.IsNull() && !m.Metadata.IsUnknown() {
		m.Metadata.ElementsAs(ctx, &metadata, false)
	}
	return &CloudKeymanagerSecretUpdatePayload{
		Checksum:   m.Checksum.ValueString(),
		TargetSpec: &CloudKeymanagerSecretUpdateTargetSpec{Metadata: metadata},
	}
}

// KeymanagerSecretCurrentStateAttrTypes returns the attribute types for the secret current_state object
func KeymanagerSecretCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                     ovhtypes.TfStringType{},
		"secret_type":              ovhtypes.TfStringType{},
		"algorithm":                ovhtypes.TfStringType{},
		"bit_length":               types.Int64Type,
		"mode":                     ovhtypes.TfStringType{},
		"payload_content_type":     ovhtypes.TfStringType{},
		"expiration":               ovhtypes.TfStringType{},
		"secret_ref":               ovhtypes.TfStringType{},
		"status":                   ovhtypes.TfStringType{},
		"region":                   ovhtypes.TfStringType{},
		"metadata":                 types.MapType{ElemType: types.StringType},
	}
}

func buildKeymanagerSecretCurrentStateObject(ctx context.Context, state *CloudKeymanagerSecretCurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	bitLength := int64(0)
	if state.BitLength != nil {
		bitLength = *state.BitLength
	}

	// Build metadata map
	metadataMap := types.MapNull(types.StringType)
	if len(state.Metadata) > 0 {
		elems := make(map[string]attr.Value)
		for k, v := range state.Metadata {
			elems[k] = types.StringValue(v)
		}
		metadataMap, _ = types.MapValue(types.StringType, elems)
	}

	obj, _ := types.ObjectValue(
		KeymanagerSecretCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                     ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"secret_type":              ovhtypes.TfStringValue{StringValue: types.StringValue(state.SecretType)},
			"algorithm":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Algorithm)},
			"bit_length":               types.Int64Value(bitLength),
			"mode":                     ovhtypes.TfStringValue{StringValue: types.StringValue(state.Mode)},
			"payload_content_type":     ovhtypes.TfStringValue{StringValue: types.StringValue(state.PayloadContentType)},
			"expiration":               ovhtypes.TfStringValue{StringValue: types.StringValue(state.Expiration)},
			"secret_ref":               ovhtypes.TfStringValue{StringValue: types.StringValue(state.SecretRef)},
			"status":                   ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"region":                   ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"metadata":                 metadataMap,
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform secret model
func (m *CloudKeymanagerSecretModel) MergeWith(ctx context.Context, response *CloudKeymanagerSecretAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeymanagerSecretCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeymanagerSecretCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.SecretType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.SecretType)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		if response.TargetSpec.Algorithm != "" {
			m.Algorithm = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Algorithm)}
		}
		if response.TargetSpec.BitLength != nil {
			m.BitLength = types.Int64Value(*response.TargetSpec.BitLength)
		}
		if response.TargetSpec.Mode != "" {
			m.Mode = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Mode)}
		}
		if response.TargetSpec.PayloadContentType != "" {
			m.PayloadContentType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.PayloadContentType)}
		}
		if response.TargetSpec.Expiration != "" {
			m.Expiration = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Expiration)}
		}
		if response.TargetSpec.Metadata != nil && len(response.TargetSpec.Metadata) > 0 {
			elems := make(map[string]attr.Value)
			for k, v := range response.TargetSpec.Metadata {
				elems[k] = types.StringValue(v)
			}
			m.Metadata, _ = types.MapValue(types.StringType, elems)
		} else if m.Metadata.IsNull() || m.Metadata.IsUnknown() {
			m.Metadata = types.MapNull(types.StringType)
		}
	}
}

// ===========================
// Container types
// ===========================

// CloudKeymanagerContainerModel represents the Terraform model for the KMS container resource
type CloudKeymanagerContainerModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Type        ovhtypes.TfStringValue `tfsdk:"type"`

	// Optional
	SecretRefs types.List `tfsdk:"secret_refs"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types for Container
type CloudKeymanagerContainerAPIResponse struct {
	Id             string                                `json:"id"`
	Checksum       string                                `json:"checksum"`
	CreatedAt      string                                `json:"createdAt"`
	UpdatedAt      string                                `json:"updatedAt"`
	ResourceStatus string                                `json:"resourceStatus"`
	CurrentState   *CloudKeymanagerContainerCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudKeymanagerContainerTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudKeymanagerContainerCurrentState struct {
	Name         string                       `json:"name,omitempty"`
	Type         string                       `json:"type,omitempty"`
	SecretRefs   []CloudKeymanagerSecretRef   `json:"secretRefs,omitempty"`
	ContainerRef string                       `json:"containerRef,omitempty"`
	Status       string                       `json:"status,omitempty"`
	Location     *CloudKeymanagerLocation     `json:"location,omitempty"`
}

type CloudKeymanagerContainerTargetSpec struct {
	Name       string                     `json:"name,omitempty"`
	Type       string                     `json:"type,omitempty"`
	SecretRefs []CloudKeymanagerSecretRef `json:"secretRefs,omitempty"`
	Location   *CloudKeymanagerLocation   `json:"location,omitempty"`
}

type CloudKeymanagerSecretRef struct {
	Name   string                    `json:"name,omitempty"`
	Secret *CloudKeymanagerSecretID  `json:"secret,omitempty"`
}

type CloudKeymanagerSecretID struct {
	ID string `json:"id"`
}

// Create payload
type CloudKeymanagerContainerCreatePayload struct {
	TargetSpec *CloudKeymanagerContainerTargetSpec `json:"targetSpec"`
}

// Update payload
type CloudKeymanagerContainerUpdatePayload struct {
	Checksum   string                                    `json:"checksum"`
	TargetSpec *CloudKeymanagerContainerUpdateTargetSpec `json:"targetSpec"`
}

type CloudKeymanagerContainerUpdateTargetSpec struct {
	SecretRefs []CloudKeymanagerSecretRef `json:"secretRefs"`
}

// SecretRefModel is the Terraform model for a secret reference
type CloudKeymanagerSecretRefModel struct {
	Name     ovhtypes.TfStringValue `tfsdk:"name"`
	SecretId ovhtypes.TfStringValue `tfsdk:"secret_id"`
}

// SecretRefAttrTypes returns the attribute types for a secret_ref object
func KeymanagerSecretRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      ovhtypes.TfStringType{},
		"secret_id": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudKeymanagerContainerModel) ToCreate(ctx context.Context) *CloudKeymanagerContainerCreatePayload {
	targetSpec := &CloudKeymanagerContainerTargetSpec{
		Name:     m.Name.ValueString(),
		Type:     m.Type.ValueString(),
		Location: &CloudKeymanagerLocation{Region: m.Region.ValueString()},
	}

	if !m.SecretRefs.IsNull() && !m.SecretRefs.IsUnknown() {
		var refs []CloudKeymanagerSecretRefModel
		m.SecretRefs.ElementsAs(ctx, &refs, false)
		for _, ref := range refs {
			sr := CloudKeymanagerSecretRef{
				Name: ref.Name.ValueString(),
				Secret: &CloudKeymanagerSecretID{
					ID: ref.SecretId.ValueString(),
				},
			}
			targetSpec.SecretRefs = append(targetSpec.SecretRefs, sr)
		}
	}

	return &CloudKeymanagerContainerCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudKeymanagerContainerModel) ToUpdate(ctx context.Context) *CloudKeymanagerContainerUpdatePayload {
	var refs []CloudKeymanagerSecretRef
	if !m.SecretRefs.IsNull() && !m.SecretRefs.IsUnknown() {
		var refModels []CloudKeymanagerSecretRefModel
		m.SecretRefs.ElementsAs(ctx, &refModels, false)
		for _, ref := range refModels {
			refs = append(refs, CloudKeymanagerSecretRef{
				Name:   ref.Name.ValueString(),
				Secret: &CloudKeymanagerSecretID{ID: ref.SecretId.ValueString()},
			})
		}
	}
	if refs == nil {
		refs = []CloudKeymanagerSecretRef{}
	}
	return &CloudKeymanagerContainerUpdatePayload{
		Checksum: m.Checksum.ValueString(),
		TargetSpec: &CloudKeymanagerContainerUpdateTargetSpec{
			SecretRefs: refs,
		},
	}
}

// KeymanagerContainerCurrentStateAttrTypes returns the attribute types for the container current_state object
func KeymanagerContainerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"type":          ovhtypes.TfStringType{},
		"container_ref": ovhtypes.TfStringType{},
		"status":        ovhtypes.TfStringType{},
		"region":        ovhtypes.TfStringType{},
		"secret_refs": types.ListType{ElemType: types.ObjectType{
			AttrTypes: KeymanagerSecretRefAttrTypes(),
		}},
	}
}

func buildKeymanagerContainerCurrentStateObject(ctx context.Context, state *CloudKeymanagerContainerCurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	// Build secret_refs list
	var secretRefValues []attr.Value
	for _, sr := range state.SecretRefs {
		secretID := ""
		if sr.Secret != nil {
			secretID = sr.Secret.ID
		}
		refObj, _ := types.ObjectValue(
			KeymanagerSecretRefAttrTypes(),
			map[string]attr.Value{
				"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(sr.Name)},
				"secret_id": ovhtypes.TfStringValue{StringValue: types.StringValue(secretID)},
			},
		)
		secretRefValues = append(secretRefValues, refObj)
	}

	secretRefsList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeymanagerSecretRefAttrTypes()},
		secretRefValues,
	)

	obj, _ := types.ObjectValue(
		KeymanagerContainerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"type":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Type)},
			"container_ref": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ContainerRef)},
			"status":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"region":        ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"secret_refs":   secretRefsList,
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform container model
func (m *CloudKeymanagerContainerModel) MergeWith(ctx context.Context, response *CloudKeymanagerContainerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeymanagerContainerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeymanagerContainerCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Type = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Type)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}

		// Build secret_refs from targetSpec
		if len(response.TargetSpec.SecretRefs) > 0 {
			var refValues []attr.Value
			for _, sr := range response.TargetSpec.SecretRefs {
				secretID := ""
				if sr.Secret != nil {
					secretID = sr.Secret.ID
				}
				refObj, _ := types.ObjectValue(
					KeymanagerSecretRefAttrTypes(),
					map[string]attr.Value{
						"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(sr.Name)},
						"secret_id": ovhtypes.TfStringValue{StringValue: types.StringValue(secretID)},
					},
				)
				refValues = append(refValues, refObj)
			}
			secretRefsList, _ := types.ListValue(
				types.ObjectType{AttrTypes: KeymanagerSecretRefAttrTypes()},
				refValues,
			)
			m.SecretRefs = secretRefsList
		} else if m.SecretRefs.IsNull() || m.SecretRefs.IsUnknown() {
			m.SecretRefs = types.ListNull(types.ObjectType{AttrTypes: KeymanagerSecretRefAttrTypes()})
		}
	}
}

// ===========================
// Data source models
// ===========================

// CloudKeymanagerSecretDataSourceModel is the Terraform model for the single secret data source
type CloudKeymanagerSecretDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId    ovhtypes.TfStringValue `tfsdk:"secret_id"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Region         ovhtypes.TfStringValue `tfsdk:"region"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	SecretType     ovhtypes.TfStringValue `tfsdk:"secret_type"`
	Metadata       types.Map              `tfsdk:"metadata"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// MergeWith merges API response data into the data source model
func (m *CloudKeymanagerSecretDataSourceModel) MergeWith(ctx context.Context, response *CloudKeymanagerSecretAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeymanagerSecretCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeymanagerSecretCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.SecretType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.SecretType)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
		if response.TargetSpec.Metadata != nil && len(response.TargetSpec.Metadata) > 0 {
			elems := make(map[string]attr.Value)
			for k, v := range response.TargetSpec.Metadata {
				elems[k] = types.StringValue(v)
			}
			m.Metadata, _ = types.MapValue(types.StringType, elems)
		} else {
			m.Metadata = types.MapNull(types.StringType)
		}
	}
}

// CloudKeymanagerSecretsDataSourceModel is the Terraform model for the secrets list data source
type CloudKeymanagerSecretsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Secrets     types.List             `tfsdk:"secrets"`
}

// CloudKeymanagerContainerDataSourceModel is the Terraform model for the single container data source
type CloudKeymanagerContainerDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId ovhtypes.TfStringValue `tfsdk:"container_id"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Region         ovhtypes.TfStringValue `tfsdk:"region"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	ContainerType  ovhtypes.TfStringValue `tfsdk:"type"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// MergeWith merges API response data into the container data source model
func (m *CloudKeymanagerContainerDataSourceModel) MergeWith(ctx context.Context, response *CloudKeymanagerContainerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeymanagerContainerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeymanagerContainerCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.ContainerType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Type)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
		}
	}
}

// CloudKeymanagerContainersDataSourceModel is the Terraform model for the containers list data source
type CloudKeymanagerContainersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Containers  types.List             `tfsdk:"containers"`
}

// KeymanagerSecretListItemAttrTypes returns attr types for secret items in data source list
func KeymanagerSecretListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"region":          ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"secret_type":     ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: KeymanagerSecretCurrentStateAttrTypes()},
	}
}

// KeymanagerContainerListItemAttrTypes returns attr types for container items in data source list
func KeymanagerContainerListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"region":          ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"type":            ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: KeymanagerContainerCurrentStateAttrTypes()},
	}
}

// ===========================
// Secret Consumer types
// ===========================

// CloudKeymanagerSecretConsumerModel represents the Terraform model for a secret consumer
type CloudKeymanagerSecretConsumerModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId     ovhtypes.TfStringValue `tfsdk:"secret_id"`
	Service      ovhtypes.TfStringValue `tfsdk:"service"`
	ResourceType ovhtypes.TfStringValue `tfsdk:"resource_type"`
	ResourceId   ovhtypes.TfStringValue `tfsdk:"resource_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
}

// CloudKeymanagerSecretConsumerAPIResponse is the API response for a single consumer
type CloudKeymanagerSecretConsumerAPIResponse struct {
	Id           string `json:"id"`
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeymanagerSecretConsumerPayload is the request body for registering a consumer
type CloudKeymanagerSecretConsumerPayload struct {
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeymanagerSecretConsumersDataSourceModel is the Terraform model for listing consumers
type CloudKeymanagerSecretConsumersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId    ovhtypes.TfStringValue `tfsdk:"secret_id"`
	Consumers   types.List             `tfsdk:"consumers"`
}

// KeymanagerSecretConsumerAttrTypes returns the attribute types for a consumer list item
func KeymanagerSecretConsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"service":       ovhtypes.TfStringType{},
		"resource_type": ovhtypes.TfStringType{},
		"resource_id":   ovhtypes.TfStringType{},
	}
}
