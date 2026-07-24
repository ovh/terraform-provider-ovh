package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageFileShareModel represents the Terraform model for the file storage share resource
type CloudStorageFileShareModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Protocol    ovhtypes.TfStringValue `tfsdk:"protocol"`
	ShareType   ovhtypes.TfStringValue `tfsdk:"share_type"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`
	ShareNetworkId   ovhtypes.TfStringValue `tfsdk:"share_network_id"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`
	Size types.Int64            `tfsdk:"size"`

	// Optional — mutable
	Description ovhtypes.TfStringValue `tfsdk:"description"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types
type CloudStorageFileShareAPIResponse struct {
	Id             string                                `json:"id"`
	Checksum       string                                `json:"checksum"`
	CreatedAt      string                                `json:"createdAt"`
	UpdatedAt      string                                `json:"updatedAt"`
	ResourceStatus string                                `json:"resourceStatus"`
	CurrentState   *CloudStorageFileShareAPICurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudStorageFileShareAPITargetSpec   `json:"targetSpec,omitempty"`
}

type CloudStorageFileShareAPICurrentState struct {
	Name            string                                   `json:"name,omitempty"`
	Description     string                                   `json:"description,omitempty"`
	Size            int64                                    `json:"size,omitempty"`
	Protocol        string                                   `json:"protocol,omitempty"`
	ShareType       string                                   `json:"shareType,omitempty"`
	Location        *CloudStorageFileShareAPILocation        `json:"location,omitempty"`
	ShareNetwork    *CloudStorageFileShareAPIShareNetworkRef `json:"shareNetwork,omitempty"`
	ExportLocations []CloudStorageFileShareAPIExportLocation `json:"exportLocations,omitempty"`
	Capabilities    []CloudStorageFileShareAPICapability     `json:"capabilities,omitempty"`
}

type CloudStorageFileShareAPILocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudStorageFileShareAPIExportLocation struct {
	Path      string `json:"path,omitempty"`
	Preferred bool   `json:"preferred"`
}

type CloudStorageFileShareAPICapability struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason,omitempty"`
}

type CloudStorageFileShareAPIShareNetworkRef struct {
	Id string `json:"id"`
}

type CloudStorageFileShareAPITargetSpec struct {
	Name         string                                   `json:"name,omitempty"`
	Description  string                                   `json:"description,omitempty"`
	Size         int64                                    `json:"size,omitempty"`
	Protocol     string                                   `json:"protocol,omitempty"`
	ShareType    string                                   `json:"shareType,omitempty"`
	Location     *CloudStorageFileShareAPILocation        `json:"location,omitempty"`
	ShareNetwork *CloudStorageFileShareAPIShareNetworkRef `json:"shareNetwork,omitempty"`
}

type CloudStorageFileShareAPIUpdateTargetSpec struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`
	Size        int64  `json:"size,omitempty"`
}

// Create payload
type CloudStorageFileShareCreatePayload struct {
	TargetSpec *CloudStorageFileShareAPITargetSpec `json:"targetSpec"`
}

// Update payload
type CloudStorageFileShareUpdatePayload struct {
	Checksum   string                                    `json:"checksum"`
	TargetSpec *CloudStorageFileShareAPIUpdateTargetSpec `json:"targetSpec"`
}

// fileShareLocationAttrTypes returns the attr types for the root-level location
// object exposed by the file share data sources.
func fileShareLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageFileShareModel) ToCreate(ctx context.Context) *CloudStorageFileShareCreatePayload {
	target := &CloudStorageFileShareAPITargetSpec{
		Name:      m.Name.ValueString(),
		Size:      m.Size.ValueInt64(),
		Protocol:  m.Protocol.ValueString(),
		ShareType: m.ShareType.ValueString(),
		Location:  &CloudStorageFileShareAPILocation{Region: m.Region.ValueString()},
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		target.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	// shareNetwork is required: always attach the reference.
	target.ShareNetwork = &CloudStorageFileShareAPIShareNetworkRef{Id: m.ShareNetworkId.ValueString()}

	return &CloudStorageFileShareCreatePayload{TargetSpec: target}
}

// ToUpdate converts the Terraform model to the API update payload
func (m *CloudStorageFileShareModel) ToUpdate(ctx context.Context, checksum string) *CloudStorageFileShareUpdatePayload {
	target := &CloudStorageFileShareAPIUpdateTargetSpec{
		Name: m.Name.ValueString(),
		Size: m.Size.ValueInt64(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	return &CloudStorageFileShareUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// FileShareCurrentStateAttrTypes returns the attribute types for the current_state object
func FileShareCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"size":        types.Int64Type,
		"protocol":    ovhtypes.TfStringType{},
		"share_type":  ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		}},
		"share_network_id": ovhtypes.TfStringType{},
		"export_locations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"path":      ovhtypes.TfStringType{},
			"preferred": types.BoolType,
		}}},
		"capabilities": types.ListType{ElemType: types.ObjectType{AttrTypes: FileShareCurrentStateCapabilityAttrTypes()}},
	}
}

// FileShareCurrentStateCapabilityAttrTypes returns the attribute types for current_state capabilities
func FileShareCurrentStateCapabilityAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":    ovhtypes.TfStringType{},
		"enabled": types.BoolType,
		"reason":  ovhtypes.TfStringType{},
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudStorageFileShareModel) MergeWith(ctx context.Context, response *CloudStorageFileShareAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildFileShareCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(FileShareCurrentStateAttrTypes())
	}

	// Set root-level fields from targetSpec
	if response.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Name)}
		m.Size = types.Int64Value(response.TargetSpec.Size)

		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}

		if response.TargetSpec.Protocol != "" {
			m.Protocol = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Protocol)}
		}
		if response.TargetSpec.ShareType != "" {
			m.ShareType = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.ShareType)}
		}
		if response.TargetSpec.ShareNetwork != nil && response.TargetSpec.ShareNetwork.Id != "" {
			m.ShareNetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.ShareNetwork.Id)}
		}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
		}
	}
}

// buildFileShareCurrentStateObject constructs the current_state object from the API response
func buildFileShareCurrentStateObject(ctx context.Context, state *CloudStorageFileShareAPICurrentState) types.Object {
	region := ""
	az := ""
	if state.Location != nil {
		region = state.Location.Region
		az = state.Location.AvailabilityZone
	}

	locObj, _ := types.ObjectValue(
		map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(az)},
		},
	)

	shareNetworkId := ""
	if state.ShareNetwork != nil {
		shareNetworkId = state.ShareNetwork.Id
	}

	// Build export_locations list
	exportLocationType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"path":      ovhtypes.TfStringType{},
		"preferred": types.BoolType,
	}}
	var exportLocValues []attr.Value
	for _, el := range state.ExportLocations {
		obj, _ := types.ObjectValue(
			map[string]attr.Type{
				"path":      ovhtypes.TfStringType{},
				"preferred": types.BoolType,
			},
			map[string]attr.Value{
				"path":      ovhtypes.TfStringValue{StringValue: types.StringValue(el.Path)},
				"preferred": types.BoolValue(el.Preferred),
			},
		)
		exportLocValues = append(exportLocValues, obj)
	}
	var exportLocList types.List
	if len(exportLocValues) > 0 {
		exportLocList, _ = types.ListValue(exportLocationType, exportLocValues)
	} else {
		exportLocList = types.ListValueMust(exportLocationType, []attr.Value{})
	}

	// Build capabilities list for current_state
	capabilityType := types.ObjectType{AttrTypes: FileShareCurrentStateCapabilityAttrTypes()}
	var capabilityValues []attr.Value
	for _, c := range state.Capabilities {
		obj, _ := types.ObjectValue(
			FileShareCurrentStateCapabilityAttrTypes(),
			map[string]attr.Value{
				"name":    ovhtypes.TfStringValue{StringValue: types.StringValue(c.Name)},
				"enabled": types.BoolValue(c.Enabled),
				"reason":  ovhtypes.TfStringValue{StringValue: types.StringValue(c.Reason)},
			},
		)
		capabilityValues = append(capabilityValues, obj)
	}
	var capabilityList types.List
	if len(capabilityValues) > 0 {
		capabilityList, _ = types.ListValue(capabilityType, capabilityValues)
	} else {
		capabilityList = types.ListValueMust(capabilityType, []attr.Value{})
	}

	obj, _ := types.ObjectValue(
		FileShareCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":             ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":      ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"size":             types.Int64Value(state.Size),
			"protocol":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Protocol)},
			"share_type":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.ShareType)},
			"location":         locObj,
			"share_network_id": ovhtypes.TfStringValue{StringValue: types.StringValue(shareNetworkId)},
			"export_locations": exportLocList,
			"capabilities":     capabilityList,
		},
	)

	return obj
}
