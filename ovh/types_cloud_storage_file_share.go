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
	NetworkId   ovhtypes.TfStringValue `tfsdk:"network_id"`
	SubnetId    ovhtypes.TfStringValue `tfsdk:"subnet_id"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`

	// Required — mutable
	Name ovhtypes.TfStringValue `tfsdk:"name"`
	Size types.Int64            `tfsdk:"size"`

	// Optional — mutable
	Description ovhtypes.TfStringValue `tfsdk:"description"`
	AccessRules types.List             `tfsdk:"access_rules"`

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
	Name            string                                      `json:"name,omitempty"`
	Description     string                                      `json:"description,omitempty"`
	Size            int64                                       `json:"size,omitempty"`
	Protocol        string                                      `json:"protocol,omitempty"`
	ShareType       string                                      `json:"shareType,omitempty"`
	NetworkId       string                                      `json:"networkId,omitempty"`
	SubnetId        string                                      `json:"subnetId,omitempty"`
	Location        *CloudStorageFileShareAPILocation           `json:"location,omitempty"`
	ExportLocations []CloudStorageFileShareAPIExportLocation    `json:"exportLocations,omitempty"`
	AccessRules     []CloudStorageFileShareAPICurrentAccessRule `json:"accessRules,omitempty"`
}

type CloudStorageFileShareAPILocation struct {
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudStorageFileShareAPIExportLocation struct {
	Path      string `json:"path,omitempty"`
	Preferred bool   `json:"preferred"`
}

type CloudStorageFileShareAPICurrentAccessRule struct {
	Id          string `json:"id,omitempty"`
	AccessTo    string `json:"accessTo,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
	State       string `json:"state,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type CloudStorageFileShareAPITargetSpec struct {
	Name        string                               `json:"name,omitempty"`
	Description string                               `json:"description,omitempty"`
	Size        int64                                `json:"size,omitempty"`
	Protocol    string                               `json:"protocol,omitempty"`
	ShareType   string                               `json:"shareType,omitempty"`
	NetworkId   string                               `json:"networkId,omitempty"`
	SubnetId    string                               `json:"subnetId,omitempty"`
	Location    *CloudStorageFileShareAPILocation    `json:"location,omitempty"`
	AccessRules []CloudStorageFileShareAPIAccessRule `json:"accessRules,omitempty"`
}

type CloudStorageFileShareAPIAccessRule struct {
	AccessTo    string `json:"accessTo"`
	AccessLevel string `json:"accessLevel"`
}

type CloudStorageFileShareAPIUpdateTargetSpec struct {
	Name        string                               `json:"name,omitempty"`
	Description string                               `json:"description"`
	Size        int64                                `json:"size,omitempty"`
	AccessRules []CloudStorageFileShareAPIAccessRule `json:"accessRules,omitempty"`
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

// AccessRule Terraform model for user-facing access_rules attribute
type CloudStorageFileShareAccessRuleModel struct {
	AccessTo    ovhtypes.TfStringValue `tfsdk:"access_to"`
	AccessLevel ovhtypes.TfStringValue `tfsdk:"access_level"`
}

// AccessRuleAttrTypes returns the attribute types for access_rules list elements
func FileShareAccessRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_to":    ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudStorageFileShareModel) ToCreate(ctx context.Context) *CloudStorageFileShareCreatePayload {
	target := &CloudStorageFileShareAPITargetSpec{
		Name:      m.Name.ValueString(),
		Size:      m.Size.ValueInt64(),
		Protocol:  m.Protocol.ValueString(),
		ShareType: m.ShareType.ValueString(),
		NetworkId: m.NetworkId.ValueString(),
		SubnetId:  m.SubnetId.ValueString(),
		Location:  &CloudStorageFileShareAPILocation{Region: m.Region.ValueString()},
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		target.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	target.AccessRules = extractAccessRules(ctx, m.AccessRules)

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

	target.AccessRules = extractAccessRules(ctx, m.AccessRules)

	return &CloudStorageFileShareUpdatePayload{Checksum: checksum, TargetSpec: target}
}

// extractAccessRules converts the Terraform list to API access rules
func extractAccessRules(ctx context.Context, accessRulesList types.List) []CloudStorageFileShareAPIAccessRule {
	if accessRulesList.IsNull() || accessRulesList.IsUnknown() {
		return nil
	}

	var rules []CloudStorageFileShareAccessRuleModel
	accessRulesList.ElementsAs(ctx, &rules, false)

	apiRules := make([]CloudStorageFileShareAPIAccessRule, len(rules))
	for i, r := range rules {
		apiRules[i] = CloudStorageFileShareAPIAccessRule{
			AccessTo:    r.AccessTo.ValueString(),
			AccessLevel: r.AccessLevel.ValueString(),
		}
	}

	return apiRules
}

// FileShareCurrentStateAttrTypes returns the attribute types for the current_state object
func FileShareCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"size":        types.Int64Type,
		"protocol":    ovhtypes.TfStringType{},
		"share_type":  ovhtypes.TfStringType{},
		"network_id":  ovhtypes.TfStringType{},
		"subnet_id":   ovhtypes.TfStringType{},
		"location": types.ObjectType{AttrTypes: map[string]attr.Type{
			"region":            ovhtypes.TfStringType{},
			"availability_zone": ovhtypes.TfStringType{},
		}},
		"export_locations": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"path":      ovhtypes.TfStringType{},
			"preferred": types.BoolType,
		}}},
		"access_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: FileShareCurrentStateAccessRuleAttrTypes()}},
	}
}

// FileShareCurrentStateAccessRuleAttrTypes returns the attribute types for current_state access_rules
func FileShareCurrentStateAccessRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           ovhtypes.TfStringType{},
		"access_to":    ovhtypes.TfStringType{},
		"access_level": ovhtypes.TfStringType{},
		"state":        ovhtypes.TfStringType{},
		"created_at":   ovhtypes.TfStringType{},
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
		if response.TargetSpec.NetworkId != "" {
			m.NetworkId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.NetworkId)}
		}
		if response.TargetSpec.SubnetId != "" {
			m.SubnetId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.SubnetId)}
		}
		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
		}

		// Set access_rules from targetSpec
		if response.TargetSpec.AccessRules != nil {
			ruleValues := make([]attr.Value, len(response.TargetSpec.AccessRules))
			for i, r := range response.TargetSpec.AccessRules {
				obj, _ := types.ObjectValue(
					FileShareAccessRuleAttrTypes(),
					map[string]attr.Value{
						"access_to":    ovhtypes.TfStringValue{StringValue: types.StringValue(r.AccessTo)},
						"access_level": ovhtypes.TfStringValue{StringValue: types.StringValue(r.AccessLevel)},
					},
				)
				ruleValues[i] = obj
			}
			m.AccessRules, _ = types.ListValue(types.ObjectType{AttrTypes: FileShareAccessRuleAttrTypes()}, ruleValues)
		} else if !m.AccessRules.IsNull() {
			m.AccessRules, _ = types.ListValue(types.ObjectType{AttrTypes: FileShareAccessRuleAttrTypes()}, []attr.Value{})
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

	// Build access_rules list for current_state
	accessRuleType := types.ObjectType{AttrTypes: FileShareCurrentStateAccessRuleAttrTypes()}
	var accessRuleValues []attr.Value
	for _, ar := range state.AccessRules {
		obj, _ := types.ObjectValue(
			FileShareCurrentStateAccessRuleAttrTypes(),
			map[string]attr.Value{
				"id":           ovhtypes.TfStringValue{StringValue: types.StringValue(ar.Id)},
				"access_to":    ovhtypes.TfStringValue{StringValue: types.StringValue(ar.AccessTo)},
				"access_level": ovhtypes.TfStringValue{StringValue: types.StringValue(ar.AccessLevel)},
				"state":        ovhtypes.TfStringValue{StringValue: types.StringValue(ar.State)},
				"created_at":   ovhtypes.TfStringValue{StringValue: types.StringValue(ar.CreatedAt)},
			},
		)
		accessRuleValues = append(accessRuleValues, obj)
	}
	var accessRuleList types.List
	if len(accessRuleValues) > 0 {
		accessRuleList, _ = types.ListValue(accessRuleType, accessRuleValues)
	} else {
		accessRuleList = types.ListValueMust(accessRuleType, []attr.Value{})
	}

	obj, _ := types.ObjectValue(
		FileShareCurrentStateAttrTypes(),
		map[string]attr.Value{
			"name":             ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
			"description":      ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"size":             types.Int64Value(state.Size),
			"protocol":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Protocol)},
			"share_type":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.ShareType)},
			"network_id":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.NetworkId)},
			"subnet_id":        ovhtypes.TfStringValue{StringValue: types.StringValue(state.SubnetId)},
			"location":         locObj,
			"export_locations": exportLocList,
			"access_rules":     accessRuleList,
		},
	)

	return obj
}
