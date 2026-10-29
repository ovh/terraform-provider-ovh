package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageSwiftContainerModel represents the Terraform model for the Swift container resource
type CloudStorageSwiftContainerModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — immutable
	Archive types.Bool `tfsdk:"archive"`

	// Optional — mutable
	Type ovhtypes.TfStringValue `tfsdk:"type"`
	Cors types.Set              `tfsdk:"cors"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudStorageSwiftContainerAPIResponse struct {
	Id             string                                     `json:"id"`
	Checksum       string                                     `json:"checksum"`
	CreatedAt      string                                     `json:"createdAt"`
	UpdatedAt      string                                     `json:"updatedAt"`
	ResourceStatus string                                     `json:"resourceStatus"`
	TargetSpec     *CloudStorageSwiftContainerAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudStorageSwiftContainerAPICurrentState `json:"currentState,omitempty"`
	CurrentTasks   []CloudResourceTask                        `json:"currentTasks,omitempty"`
}

type CloudStorageSwiftContainerAPILocation struct {
	Region string `json:"region,omitempty"`
}

type CloudStorageSwiftContainerAPITargetSpec struct {
	Name     string                                 `json:"name,omitempty"`
	Location *CloudStorageSwiftContainerAPILocation `json:"location,omitempty"`
	Archive  *bool                                  `json:"archive,omitempty"`
	Type     string                                 `json:"type,omitempty"`
	Cors     []string                               `json:"cors,omitempty"`
}

type CloudStorageSwiftContainerAPIUpdateTargetSpec struct {
	Type string `json:"type,omitempty"`
	// Never omitted: null and an empty list both clear the origins, so an explicit list keeps the intent visible.
	Cors []string `json:"cors"`
}

type CloudStorageSwiftContainerAPICurrentState struct {
	Name     string                                 `json:"name,omitempty"`
	Location *CloudStorageSwiftContainerAPILocation `json:"location,omitempty"`
	Archive  bool                                   `json:"archive"`
	Type     string                                 `json:"type,omitempty"`
	Cors     []string                               `json:"cors,omitempty"`
	// Pointers: the API fills these on a single container GET, a LIST may leave them out.
	ObjectsCount *int64 `json:"objectsCount,omitempty"`
	ObjectsSize  *int64 `json:"objectsSize,omitempty"`
	StaticUrl    string `json:"staticUrl,omitempty"`
}

// CloudStorageSwiftContainerAPIObject is one object of a container, returned by the synchronous object routes
type CloudStorageSwiftContainerAPIObject struct {
	Key            string `json:"key"`
	Size           int64  `json:"size"`
	ContentType    string `json:"contentType"`
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
	Etag           string `json:"etag"`
	RetrievalState string `json:"retrievalState,omitempty"`
	RetrievalDelay *int64 `json:"retrievalDelay,omitempty"`
}

// Create payload
type CloudStorageSwiftContainerCreatePayload struct {
	TargetSpec *CloudStorageSwiftContainerAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudStorageSwiftContainerUpdatePayload struct {
	Checksum   string                                         `json:"checksum"`
	TargetSpec *CloudStorageSwiftContainerAPIUpdateTargetSpec `json:"targetSpec"`
}

// StorageSwiftContainerLocationAttrTypes returns the attribute types for the location object
func StorageSwiftContainerLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region": ovhtypes.TfStringType{},
	}
}

// StorageSwiftContainerCurrentStateAttrTypes returns the attribute types for the current_state object
func StorageSwiftContainerCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"location":      types.ObjectType{AttrTypes: StorageSwiftContainerLocationAttrTypes()},
		"archive":       types.BoolType,
		"type":          ovhtypes.TfStringType{},
		"cors":          types.SetType{ElemType: ovhtypes.TfStringType{}},
		"objects_count": types.Int64Type,
		"objects_size":  types.Int64Type,
		"static_url":    ovhtypes.TfStringType{},
	}
}

// StorageSwiftContainerObjectAttrTypes returns the attribute types for one object of the objects data source
func StorageSwiftContainerObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":              ovhtypes.TfStringType{},
		"size":             types.Int64Type,
		"content_type":     ovhtypes.TfStringType{},
		"last_modified_at": ovhtypes.TfStringType{},
		"etag":             ovhtypes.TfStringType{},
		"retrieval_state":  ovhtypes.TfStringType{},
		"retrieval_delay":  types.Int64Type,
	}
}

func buildStorageSwiftContainerLocationObject(location *CloudStorageSwiftContainerAPILocation) types.Object {
	if location == nil {
		return types.ObjectNull(StorageSwiftContainerLocationAttrTypes())
	}

	obj, _ := types.ObjectValue(
		StorageSwiftContainerLocationAttrTypes(),
		map[string]attr.Value{
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(location.Region)},
		},
	)
	return obj
}

// buildStorageSwiftContainerCorsSet never returns null: a null set would not match a configured empty one.
func buildStorageSwiftContainerCorsSet(cors []string) types.Set {
	elems := make([]attr.Value, 0, len(cors))
	for _, origin := range cors {
		elems = append(elems, ovhtypes.TfStringValue{StringValue: types.StringValue(origin)})
	}
	set, _ := types.SetValue(ovhtypes.TfStringType{}, elems)
	return set
}

func storageSwiftContainerCorsToAPI(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	cors := make([]string, 0, len(set.Elements()))
	for _, v := range set.Elements() {
		if sv, ok := v.(ovhtypes.TfStringValue); ok && !sv.IsNull() && !sv.IsUnknown() {
			cors = append(cors, sv.ValueString())
		}
	}
	return cors
}

func buildStorageSwiftContainerObjectValue(object *CloudStorageSwiftContainerAPIObject) types.Object {
	obj, _ := types.ObjectValue(
		StorageSwiftContainerObjectAttrTypes(),
		map[string]attr.Value{
			"key":              ovhtypes.TfStringValue{StringValue: types.StringValue(object.Key)},
			"size":             types.Int64Value(object.Size),
			"content_type":     ovhtypes.TfStringValue{StringValue: types.StringValue(object.ContentType)},
			"last_modified_at": nullableTfString(object.LastModifiedAt),
			"etag":             ovhtypes.TfStringValue{StringValue: types.StringValue(object.Etag)},
			"retrieval_state":  nullableTfString(object.RetrievalState),
			"retrieval_delay":  nullableInt64(object.RetrievalDelay),
		},
	)
	return obj
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageSwiftContainerModel) ToCreate() *CloudStorageSwiftContainerCreatePayload {
	target := &CloudStorageSwiftContainerAPITargetSpec{
		Name:     m.Name.ValueString(),
		Location: &CloudStorageSwiftContainerAPILocation{Region: m.Region.ValueString()},
		Cors:     storageSwiftContainerCorsToAPI(m.Cors),
	}

	if !m.Archive.IsNull() && !m.Archive.IsUnknown() {
		archive := m.Archive.ValueBool()
		target.Archive = &archive
	}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		target.Type = m.Type.ValueString()
	}

	return &CloudStorageSwiftContainerCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageSwiftContainerModel) ToUpdate(checksum string) *CloudStorageSwiftContainerUpdatePayload {
	cors := storageSwiftContainerCorsToAPI(m.Cors)
	if cors == nil {
		cors = []string{}
	}
	target := &CloudStorageSwiftContainerAPIUpdateTargetSpec{Cors: cors}

	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		target.Type = m.Type.ValueString()
	}

	return &CloudStorageSwiftContainerUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// buildStorageSwiftContainerCurrentStateObject constructs the current_state object from the API response
func buildStorageSwiftContainerCurrentStateObject(ctx context.Context, state *CloudStorageSwiftContainerAPICurrentState) types.Object {
	obj, _ := types.ObjectValue(
		StorageSwiftContainerCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"location":      buildStorageSwiftContainerLocationObject(state.Location),
			"archive":       types.BoolValue(state.Archive),
			"type":          nullableTfString(state.Type),
			"cors":          buildStorageSwiftContainerCorsSet(state.Cors),
			"objects_count": nullableInt64(state.ObjectsCount),
			"objects_size":  nullableInt64(state.ObjectsSize),
			"static_url":    nullableTfString(state.StaticUrl),
		},
	)
	return obj
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageSwiftContainerModel) MergeWith(ctx context.Context, response *CloudStorageSwiftContainerAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildStorageSwiftContainerCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(StorageSwiftContainerCurrentStateAttrTypes())
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

	// The API stores the defaults it applied (archive false, type PRIVATE) in the target
	// spec; until it does, the observed state is the only source of the effective value.
	archive := false
	if response.TargetSpec.Archive != nil {
		archive = *response.TargetSpec.Archive
	} else if response.CurrentState != nil {
		archive = response.CurrentState.Archive
	}
	m.Archive = types.BoolValue(archive)

	containerType := response.TargetSpec.Type
	if containerType == "" && response.CurrentState != nil {
		containerType = response.CurrentState.Type
	}
	if containerType != "" {
		m.Type = ovhtypes.TfStringValue{StringValue: types.StringValue(containerType)}
	} else if m.Type.IsUnknown() {
		m.Type = ovhtypes.NewTfStringNull()
	}

	m.Cors = buildStorageSwiftContainerCorsSet(response.TargetSpec.Cors)
}
