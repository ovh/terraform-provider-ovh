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

// CloudKeyManagerSecretModel represents the Terraform model for the KMS secret resource
type CloudKeyManagerSecretModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	// Optional (immutable)
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`
	Name             ovhtypes.TfStringValue `tfsdk:"name"`
	SecretType       ovhtypes.TfStringValue `tfsdk:"secret_type"`

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
type CloudKeyManagerSecretAPIResponse struct {
	Id             string                             `json:"id"`
	Checksum       string                             `json:"checksum"`
	CreatedAt      string                             `json:"createdAt"`
	UpdatedAt      string                             `json:"updatedAt"`
	ResourceStatus string                             `json:"resourceStatus"`
	CurrentState   *CloudKeyManagerSecretCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudKeyManagerSecretTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudKeyManagerSecretCurrentState struct {
	Name               string                   `json:"name,omitempty"`
	SecretType         string                   `json:"secretType,omitempty"`
	Algorithm          string                   `json:"algorithm,omitempty"`
	BitLength          *int64                   `json:"bitLength,omitempty"`
	Mode               string                   `json:"mode,omitempty"`
	PayloadContentType string                   `json:"payloadContentType,omitempty"`
	Expiration         string                   `json:"expiration,omitempty"`
	SecretRef          string                   `json:"secretRef,omitempty"`
	Status             string                   `json:"status,omitempty"`
	Location           *CloudKeyManagerLocation `json:"location,omitempty"`
	Metadata           map[string]string        `json:"metadata,omitempty"`
}

type CloudKeyManagerLocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudKeyManagerSecretTargetSpec struct {
	Name               string                   `json:"name,omitempty"`
	SecretType         string                   `json:"secretType,omitempty"`
	Algorithm          string                   `json:"algorithm,omitempty"`
	BitLength          *int64                   `json:"bitLength,omitempty"`
	Mode               string                   `json:"mode,omitempty"`
	Payload            string                   `json:"payload,omitempty"`
	PayloadContentType string                   `json:"payloadContentType,omitempty"`
	Expiration         string                   `json:"expiration,omitempty"`
	Location           *CloudKeyManagerLocation `json:"location,omitempty"`
	Metadata           map[string]string        `json:"metadata,omitempty"`
}

// Create payload
type CloudKeyManagerSecretCreatePayload struct {
	TargetSpec *CloudKeyManagerSecretTargetSpec `json:"targetSpec"`
}

// Update payload
type CloudKeyManagerSecretUpdatePayload struct {
	Checksum   string                                 `json:"checksum"`
	TargetSpec *CloudKeyManagerSecretUpdateTargetSpec `json:"targetSpec"`
}

type CloudKeyManagerSecretUpdateTargetSpec struct {
	Metadata map[string]string `json:"metadata"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudKeyManagerSecretModel) ToCreate(ctx context.Context) *CloudKeyManagerSecretCreatePayload {
	targetSpec := &CloudKeyManagerSecretTargetSpec{
		Name:       m.Name.ValueString(),
		SecretType: m.SecretType.ValueString(),
		Location:   &CloudKeyManagerLocation{Region: m.Region.ValueString()},
	}

	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
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

	return &CloudKeyManagerSecretCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudKeyManagerSecretModel) ToUpdate(ctx context.Context) *CloudKeyManagerSecretUpdatePayload {
	metadata := map[string]string{}
	if !m.Metadata.IsNull() && !m.Metadata.IsUnknown() {
		m.Metadata.ElementsAs(ctx, &metadata, false)
	}
	return &CloudKeyManagerSecretUpdatePayload{
		Checksum:   m.Checksum.ValueString(),
		TargetSpec: &CloudKeyManagerSecretUpdateTargetSpec{Metadata: metadata},
	}
}

// keyManagerLocationAttrTypes returns the attribute types for the nested
// location object exposed by the current_state and the data sources.
func keyManagerLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// buildKeyManagerLocationObject converts an API location into the nested
// Terraform location object, returning a null object when the API omits it.
func buildKeyManagerLocationObject(location *CloudKeyManagerLocation) types.Object {
	if location == nil {
		return types.ObjectNull(keyManagerLocationAttrTypes())
	}

	obj, _ := types.ObjectValue(
		keyManagerLocationAttrTypes(),
		map[string]attr.Value{
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(location.Region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(location.AvailabilityZone)},
		},
	)
	return obj
}

// KeyManagerSecretCurrentStateAttrTypes returns the attribute types for the secret current_state object
func KeyManagerSecretCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                 ovhtypes.TfStringType{},
		"secret_type":          ovhtypes.TfStringType{},
		"algorithm":            ovhtypes.TfStringType{},
		"bit_length":           types.Int64Type,
		"mode":                 ovhtypes.TfStringType{},
		"payload_content_type": ovhtypes.TfStringType{},
		"expiration":           ovhtypes.TfStringType{},
		"secret_ref":           ovhtypes.TfStringType{},
		"status":               ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: keyManagerLocationAttrTypes(),
		},
		"metadata": types.MapType{ElemType: types.StringType},
	}
}

func buildKeyManagerSecretCurrentStateObject(ctx context.Context, state *CloudKeyManagerSecretCurrentState) types.Object {
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
		KeyManagerSecretCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":                 ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"secret_type":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.SecretType)},
			"algorithm":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.Algorithm)},
			"bit_length":           types.Int64Value(bitLength),
			"mode":                 ovhtypes.TfStringValue{StringValue: types.StringValue(state.Mode)},
			"payload_content_type": ovhtypes.TfStringValue{StringValue: types.StringValue(state.PayloadContentType)},
			"expiration":           ovhtypes.TfStringValue{StringValue: types.StringValue(state.Expiration)},
			"secret_ref":           ovhtypes.TfStringValue{StringValue: types.StringValue(state.SecretRef)},
			"status":               ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"location":             buildKeyManagerLocationObject(state.Location),
			"metadata":             metadataMap,
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform secret model
func (m *CloudKeyManagerSecretModel) MergeWith(ctx context.Context, response *CloudKeyManagerSecretAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeyManagerSecretCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeyManagerSecretCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.SecretType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.SecretType)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
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
		// The root metadata attribute is config-driven: the API never adds
		// server-side keys and updates replace it wholesale, so it is never
		// overwritten from the API response. The live server value remains
		// visible through current_state.metadata.
		if m.Metadata.IsUnknown() {
			m.Metadata = types.MapNull(types.StringType)
		}
	}
}

// ===========================
// Secret Payload types
// ===========================

// CloudKeyManagerSecretPayloadDataSourceModel is the Terraform model for reading a secret payload
type CloudKeyManagerSecretPayloadDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId    ovhtypes.TfStringValue `tfsdk:"secret_id"`
	Payload     ovhtypes.TfStringValue `tfsdk:"payload"`
}

// CloudKeyManagerSecretPayloadAPIResponse is the API response holding the secret payload
type CloudKeyManagerSecretPayloadAPIResponse struct {
	Payload string `json:"payload"`
}

// ===========================
// Container types
// ===========================

// CloudKeyManagerContainerModel represents the Terraform model for the KMS container resource
type CloudKeyManagerContainerModel struct {
	// Required
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	// Optional (immutable)
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`
	Name             ovhtypes.TfStringValue `tfsdk:"name"`
	Type             ovhtypes.TfStringValue `tfsdk:"type"`

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
type CloudKeyManagerContainerAPIResponse struct {
	Id             string                                `json:"id"`
	Checksum       string                                `json:"checksum"`
	CreatedAt      string                                `json:"createdAt"`
	UpdatedAt      string                                `json:"updatedAt"`
	ResourceStatus string                                `json:"resourceStatus"`
	CurrentState   *CloudKeyManagerContainerCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudKeyManagerContainerTargetSpec   `json:"targetSpec,omitempty"`
}

type CloudKeyManagerContainerCurrentState struct {
	Name         string                     `json:"name,omitempty"`
	Type         string                     `json:"type,omitempty"`
	SecretRefs   []CloudKeyManagerSecretRef `json:"secretRefs,omitempty"`
	ContainerRef string                     `json:"containerRef,omitempty"`
	Status       string                     `json:"status,omitempty"`
	Location     *CloudKeyManagerLocation   `json:"location,omitempty"`
}

type CloudKeyManagerContainerTargetSpec struct {
	Name       string                     `json:"name,omitempty"`
	Type       string                     `json:"type,omitempty"`
	SecretRefs []CloudKeyManagerSecretRef `json:"secretRefs,omitempty"`
	Location   *CloudKeyManagerLocation   `json:"location,omitempty"`
}

type CloudKeyManagerSecretRef struct {
	Name   string                   `json:"name,omitempty"`
	Secret *CloudKeyManagerSecretID `json:"secret,omitempty"`
}

type CloudKeyManagerSecretID struct {
	ID string `json:"id"`
}

// Create payload
type CloudKeyManagerContainerCreatePayload struct {
	TargetSpec *CloudKeyManagerContainerTargetSpec `json:"targetSpec"`
}

// Update payload
type CloudKeyManagerContainerUpdatePayload struct {
	Checksum   string                                    `json:"checksum"`
	TargetSpec *CloudKeyManagerContainerUpdateTargetSpec `json:"targetSpec"`
}

type CloudKeyManagerContainerUpdateTargetSpec struct {
	SecretRefs []CloudKeyManagerSecretRef `json:"secretRefs"`
}

// SecretRefModel is the Terraform model for a secret reference
type CloudKeyManagerSecretRefModel struct {
	Name     ovhtypes.TfStringValue `tfsdk:"name"`
	SecretId ovhtypes.TfStringValue `tfsdk:"secret_id"`
}

// SecretRefAttrTypes returns the attribute types for a secret_ref object
func KeyManagerSecretRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      ovhtypes.TfStringType{},
		"secret_id": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudKeyManagerContainerModel) ToCreate(ctx context.Context) *CloudKeyManagerContainerCreatePayload {
	targetSpec := &CloudKeyManagerContainerTargetSpec{
		Name:     m.Name.ValueString(),
		Type:     m.Type.ValueString(),
		Location: &CloudKeyManagerLocation{Region: m.Region.ValueString()},
	}

	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	if !m.SecretRefs.IsNull() && !m.SecretRefs.IsUnknown() {
		var refs []CloudKeyManagerSecretRefModel
		m.SecretRefs.ElementsAs(ctx, &refs, false)
		for _, ref := range refs {
			sr := CloudKeyManagerSecretRef{
				Name: ref.Name.ValueString(),
				Secret: &CloudKeyManagerSecretID{
					ID: ref.SecretId.ValueString(),
				},
			}
			targetSpec.SecretRefs = append(targetSpec.SecretRefs, sr)
		}
	}

	return &CloudKeyManagerContainerCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudKeyManagerContainerModel) ToUpdate(ctx context.Context) *CloudKeyManagerContainerUpdatePayload {
	var refs []CloudKeyManagerSecretRef
	if !m.SecretRefs.IsNull() && !m.SecretRefs.IsUnknown() {
		var refModels []CloudKeyManagerSecretRefModel
		m.SecretRefs.ElementsAs(ctx, &refModels, false)
		for _, ref := range refModels {
			refs = append(refs, CloudKeyManagerSecretRef{
				Name:   ref.Name.ValueString(),
				Secret: &CloudKeyManagerSecretID{ID: ref.SecretId.ValueString()},
			})
		}
	}
	if refs == nil {
		refs = []CloudKeyManagerSecretRef{}
	}
	return &CloudKeyManagerContainerUpdatePayload{
		Checksum: m.Checksum.ValueString(),
		TargetSpec: &CloudKeyManagerContainerUpdateTargetSpec{
			SecretRefs: refs,
		},
	}
}

// KeyManagerContainerCurrentStateAttrTypes returns the attribute types for the container current_state object
func KeyManagerContainerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"type":          ovhtypes.TfStringType{},
		"container_ref": ovhtypes.TfStringType{},
		"status":        ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: keyManagerLocationAttrTypes(),
		},
		"secret_refs": types.ListType{ElemType: types.ObjectType{
			AttrTypes: KeyManagerSecretRefAttrTypes(),
		}},
	}
}

func buildKeyManagerContainerCurrentStateObject(ctx context.Context, state *CloudKeyManagerContainerCurrentState) types.Object {
	// Build secret_refs list
	var secretRefValues []attr.Value
	for _, sr := range state.SecretRefs {
		secretID := ""
		if sr.Secret != nil {
			secretID = sr.Secret.ID
		}
		refObj, _ := types.ObjectValue(
			KeyManagerSecretRefAttrTypes(),
			map[string]attr.Value{
				"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(sr.Name)},
				"secret_id": ovhtypes.TfStringValue{StringValue: types.StringValue(secretID)},
			},
		)
		secretRefValues = append(secretRefValues, refObj)
	}

	secretRefsList, _ := types.ListValue(
		types.ObjectType{AttrTypes: KeyManagerSecretRefAttrTypes()},
		secretRefValues,
	)

	obj, _ := types.ObjectValue(
		KeyManagerContainerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"type":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Type)},
			"container_ref": ovhtypes.TfStringValue{StringValue: types.StringValue(state.ContainerRef)},
			"status":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"location":      buildKeyManagerLocationObject(state.Location),
			"secret_refs":   secretRefsList,
		},
	)

	return obj
}

// MergeWith merges API response data into the Terraform container model
func (m *CloudKeyManagerContainerModel) MergeWith(ctx context.Context, response *CloudKeyManagerContainerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeyManagerContainerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeyManagerContainerCurrentStateAttrTypes())
	}

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Type = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Type)}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
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
					KeyManagerSecretRefAttrTypes(),
					map[string]attr.Value{
						"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(sr.Name)},
						"secret_id": ovhtypes.TfStringValue{StringValue: types.StringValue(secretID)},
					},
				)
				refValues = append(refValues, refObj)
			}
			secretRefsList, _ := types.ListValue(
				types.ObjectType{AttrTypes: KeyManagerSecretRefAttrTypes()},
				refValues,
			)
			m.SecretRefs = secretRefsList
		} else if m.SecretRefs.IsNull() || m.SecretRefs.IsUnknown() {
			m.SecretRefs = types.ListNull(types.ObjectType{AttrTypes: KeyManagerSecretRefAttrTypes()})
		}
	}
}

// ===========================
// Data source models
// ===========================

// CloudKeyManagerSecretDataSourceModel is the Terraform model for the single
// secret data source. It mirrors the resource model but exposes the location
// as a nested object (the secret is fetched by ID, so the location is
// computed, not user input).
type CloudKeyManagerSecretDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId    ovhtypes.TfStringValue `tfsdk:"secret_id"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Location       types.Object           `tfsdk:"location"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	SecretType     ovhtypes.TfStringValue `tfsdk:"secret_type"`
	Metadata       types.Map              `tfsdk:"metadata"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// MergeWith merges API response data into the data source model
func (m *CloudKeyManagerSecretDataSourceModel) MergeWith(ctx context.Context, response *CloudKeyManagerSecretAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeyManagerSecretCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeyManagerSecretCurrentStateAttrTypes())
	}

	m.Location = types.ObjectNull(keyManagerLocationAttrTypes())

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.SecretType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.SecretType)}
		m.Location = buildKeyManagerLocationObject(response.TargetSpec.Location)
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

// CloudKeyManagerSecretsDataSourceModel is the Terraform model for the secrets list data source
type CloudKeyManagerSecretsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Secrets     types.List             `tfsdk:"secrets"`
}

// CloudKeyManagerContainerDataSourceModel is the Terraform model for the
// single container data source. It mirrors the resource model but exposes the
// location as a nested object (the container is fetched by ID, so the
// location is computed, not user input).
type CloudKeyManagerContainerDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId ovhtypes.TfStringValue `tfsdk:"container_id"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	Location       types.Object           `tfsdk:"location"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	ContainerType  ovhtypes.TfStringValue `tfsdk:"type"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// MergeWith merges API response data into the container data source model
func (m *CloudKeyManagerContainerDataSourceModel) MergeWith(ctx context.Context, response *CloudKeyManagerContainerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildKeyManagerContainerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(KeyManagerContainerCurrentStateAttrTypes())
	}

	m.Location = types.ObjectNull(keyManagerLocationAttrTypes())

	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.ContainerType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Type)}
		m.Location = buildKeyManagerLocationObject(response.TargetSpec.Location)
	}
}

// CloudKeyManagerContainersDataSourceModel is the Terraform model for the containers list data source
type CloudKeyManagerContainersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Containers  types.List             `tfsdk:"containers"`
}

// KeyManagerSecretListItemAttrTypes returns attr types for secret items in data source list
func KeyManagerSecretListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: keyManagerLocationAttrTypes(),
		},
		"name":          ovhtypes.TfStringType{},
		"secret_type":   ovhtypes.TfStringType{},
		"current_state": types.ObjectType{AttrTypes: KeyManagerSecretCurrentStateAttrTypes()},
	}
}

// KeyManagerContainerListItemAttrTypes returns attr types for container items in data source list
func KeyManagerContainerListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: keyManagerLocationAttrTypes(),
		},
		"name":          ovhtypes.TfStringType{},
		"type":          ovhtypes.TfStringType{},
		"current_state": types.ObjectType{AttrTypes: KeyManagerContainerCurrentStateAttrTypes()},
	}
}

// ===========================
// Secret Consumer types
// ===========================

// CloudKeyManagerSecretConsumerModel represents the Terraform model for a secret consumer
type CloudKeyManagerSecretConsumerModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId     ovhtypes.TfStringValue `tfsdk:"secret_id"`
	Service      ovhtypes.TfStringValue `tfsdk:"service"`
	ResourceType ovhtypes.TfStringValue `tfsdk:"resource_type"`
	ResourceId   ovhtypes.TfStringValue `tfsdk:"resource_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
}

// CloudKeyManagerSecretConsumerAPIResponse is the API response for a single consumer
type CloudKeyManagerSecretConsumerAPIResponse struct {
	Id           string `json:"id"`
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeyManagerSecretConsumerPayload is the request body for registering a consumer
type CloudKeyManagerSecretConsumerPayload struct {
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeyManagerSecretConsumersDataSourceModel is the Terraform model for listing consumers
type CloudKeyManagerSecretConsumersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId    ovhtypes.TfStringValue `tfsdk:"secret_id"`
	Consumers   types.List             `tfsdk:"consumers"`
}

// CloudKeyManagerSecretConsumerDataSourceModel is the Terraform model for reading a single secret consumer by its computed id
type CloudKeyManagerSecretConsumerDataSourceModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	SecretId     ovhtypes.TfStringValue `tfsdk:"secret_id"`
	ConsumerId   ovhtypes.TfStringValue `tfsdk:"consumer_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	Service      ovhtypes.TfStringValue `tfsdk:"service"`
	ResourceType ovhtypes.TfStringValue `tfsdk:"resource_type"`
	ResourceId   ovhtypes.TfStringValue `tfsdk:"resource_id"`
}

// KeyManagerSecretConsumerAttrTypes returns the attribute types for a consumer list item
func KeyManagerSecretConsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"service":       ovhtypes.TfStringType{},
		"resource_type": ovhtypes.TfStringType{},
		"resource_id":   ovhtypes.TfStringType{},
	}
}

// ===========================
// Container Consumer types
// ===========================

// CloudKeyManagerContainerConsumerModel represents the Terraform model for a container consumer
type CloudKeyManagerContainerConsumerModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId  ovhtypes.TfStringValue `tfsdk:"container_id"`
	Service      ovhtypes.TfStringValue `tfsdk:"service"`
	ResourceType ovhtypes.TfStringValue `tfsdk:"resource_type"`
	ResourceId   ovhtypes.TfStringValue `tfsdk:"resource_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
}

// CloudKeyManagerContainerConsumerAPIResponse is the API response for a single container consumer
type CloudKeyManagerContainerConsumerAPIResponse struct {
	Id           string `json:"id"`
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeyManagerContainerConsumerPayload is the request body for registering a container consumer
type CloudKeyManagerContainerConsumerPayload struct {
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

// CloudKeyManagerContainerConsumersDataSourceModel is the Terraform model for listing container consumers
type CloudKeyManagerContainerConsumersDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId ovhtypes.TfStringValue `tfsdk:"container_id"`
	Consumers   types.List             `tfsdk:"consumers"`
}

// CloudKeyManagerContainerConsumerDataSourceModel is the Terraform model for reading a single container consumer by its computed id
type CloudKeyManagerContainerConsumerDataSourceModel struct {
	ServiceName  ovhtypes.TfStringValue `tfsdk:"service_name"`
	ContainerId  ovhtypes.TfStringValue `tfsdk:"container_id"`
	ConsumerId   ovhtypes.TfStringValue `tfsdk:"consumer_id"`
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	Service      ovhtypes.TfStringValue `tfsdk:"service"`
	ResourceType ovhtypes.TfStringValue `tfsdk:"resource_type"`
	ResourceId   ovhtypes.TfStringValue `tfsdk:"resource_id"`
}

// KeyManagerContainerConsumerAttrTypes returns the attribute types for a container consumer list item
func KeyManagerContainerConsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"service":       ovhtypes.TfStringType{},
		"resource_type": ovhtypes.TfStringType{},
		"resource_id":   ovhtypes.TfStringType{},
	}
}
