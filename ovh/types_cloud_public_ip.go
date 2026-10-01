package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// Shared API types for the public IP family (floating, extNet, additional, aggregate)

// CloudPublicIPLocation is the location of a public IP.
type CloudPublicIPLocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// CloudPublicIPAssociatedResource describes the resource a public IP is currently attached to.
type CloudPublicIPAssociatedResource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// CloudFloatingIPNetwork references the external network of a floating IP.
type CloudFloatingIPNetwork struct {
	ID string `json:"id"`
}

// CloudPublicIPTaskError is an error that occurred on an asynchronous public IP task.
type CloudPublicIPTaskError struct {
	Message string `json:"message"`
}

// CloudPublicIPResourceTask represents an asynchronous operation on a public IP.
type CloudPublicIPResourceTask struct {
	ID     string                   `json:"id"`
	Link   string                   `json:"link"`
	Type   string                   `json:"type"`
	Status string                   `json:"status,omitempty"`
	Errors []CloudPublicIPTaskError `json:"errors,omitempty"`
}

// Floating IP API types

// CloudFloatingIPCurrentState is the current state of a floating public IP.
type CloudFloatingIPCurrentState struct {
	ID                 string                           `json:"id"`
	IP                 string                           `json:"ip,omitempty"`
	Status             string                           `json:"status,omitempty"`
	Description        string                           `json:"description,omitempty"`
	Network            *CloudFloatingIPNetwork          `json:"network,omitempty"`
	AssociatedResource *CloudPublicIPAssociatedResource `json:"associatedResource,omitempty"`
	Location           *CloudPublicIPLocation           `json:"location,omitempty"`
}

// CloudFloatingIPTargetSpec is the desired specification for a floating public IP.
type CloudFloatingIPTargetSpec struct {
	Description string                 `json:"description,omitempty"`
	Location    *CloudPublicIPLocation `json:"location,omitempty"`
}

// CloudFloatingIPAPIResponse is the read envelope for a floating public IP.
type CloudFloatingIPAPIResponse struct {
	ID             string                       `json:"id"`
	Checksum       string                       `json:"checksum"`
	CreatedAt      string                       `json:"createdAt"`
	UpdatedAt      string                       `json:"updatedAt"`
	ResourceStatus string                       `json:"resourceStatus"`
	CurrentTasks   []CloudPublicIPResourceTask  `json:"currentTasks,omitempty"`
	CurrentState   *CloudFloatingIPCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudFloatingIPTargetSpec   `json:"targetSpec,omitempty"`
}

// CloudFloatingIPCreatePayload is the POST body to create a floating public IP.
type CloudFloatingIPCreatePayload struct {
	TargetSpec *CloudFloatingIPTargetSpec `json:"targetSpec"`
}

// CloudFloatingIPUpdateTargetSpec is the update spec of a floating public IP.
// Note: location is immutable, description is the only mutable field.
type CloudFloatingIPUpdateTargetSpec struct {
	Description string `json:"description"`
}

// CloudFloatingIPUpdatePayload is the PUT body to update a floating public IP.
type CloudFloatingIPUpdatePayload struct {
	Checksum   string                           `json:"checksum"`
	TargetSpec *CloudFloatingIPUpdateTargetSpec `json:"targetSpec"`
}

// ExtNet IP API types (read-only, no targetSpec)

// CloudExtNetIPCurrentState is the current state of an external-network public IP.
type CloudExtNetIPCurrentState struct {
	ID                 string                           `json:"id"`
	IP                 string                           `json:"ip,omitempty"`
	AssociatedResource *CloudPublicIPAssociatedResource `json:"associatedResource,omitempty"`
	Location           *CloudPublicIPLocation           `json:"location,omitempty"`
}

// CloudExtNetIPAPIResponse is the read envelope for an external-network public IP.
type CloudExtNetIPAPIResponse struct {
	ID             string                      `json:"id"`
	Checksum       string                      `json:"checksum"`
	CreatedAt      string                      `json:"createdAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	ResourceStatus string                      `json:"resourceStatus"`
	CurrentTasks   []CloudPublicIPResourceTask `json:"currentTasks,omitempty"`
	CurrentState   *CloudExtNetIPCurrentState  `json:"currentState,omitempty"`
}

// Additional IP API types (read-only, no targetSpec, no createdAt/updatedAt)

// CloudAdditionalIPCurrentState is the current state of an additional public IP (IP alias).
type CloudAdditionalIPCurrentState struct {
	ID                 string                           `json:"id"`
	IP                 string                           `json:"ip,omitempty"`
	IPBlock            string                           `json:"ipBlock,omitempty"`
	AssociatedResource *CloudPublicIPAssociatedResource `json:"associatedResource,omitempty"`
	Location           *CloudPublicIPLocation           `json:"location,omitempty"`
}

// CloudAdditionalIPAPIResponse is the read envelope for an additional public IP (IP alias).
type CloudAdditionalIPAPIResponse struct {
	ID             string                         `json:"id"`
	Checksum       string                         `json:"checksum"`
	ResourceStatus string                         `json:"resourceStatus"`
	CurrentTasks   []CloudPublicIPResourceTask    `json:"currentTasks,omitempty"`
	CurrentState   *CloudAdditionalIPCurrentState `json:"currentState,omitempty"`
}

// Aggregate API types

// CloudPublicIPSummary is the aggregated public IP entry returned by the
// list-all endpoint. It carries only the IP address and its type.
type CloudPublicIPSummary struct {
	IP   string `json:"ip"`
	Type string `json:"type"`
}

// Terraform attribute types

// cloudPublicIPLocationAttrTypes returns the attribute types for the nested location object.
func cloudPublicIPLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":            ovhtypes.TfStringType{},
		"availability_zone": ovhtypes.TfStringType{},
	}
}

// cloudPublicIPAssociatedResourceAttrTypes returns the attribute types for the
// nested associated_resource object.
func cloudPublicIPAssociatedResourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   ovhtypes.TfStringType{},
		"type": ovhtypes.TfStringType{},
	}
}

// cloudFloatingIPNetworkAttrTypes returns the attribute types for the nested network object.
func cloudFloatingIPNetworkAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

// cloudPublicIPTaskErrorAttrTypes returns the attribute types for a single
// error entry of a current_tasks entry.
func cloudPublicIPTaskErrorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"message": ovhtypes.TfStringType{},
	}
}

// cloudPublicIPResourceTaskAttrTypes returns the attribute types for a single current_tasks entry.
func cloudPublicIPResourceTaskAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"errors": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: cloudPublicIPTaskErrorAttrTypes(),
			},
		},
		"id":     ovhtypes.TfStringType{},
		"link":   ovhtypes.TfStringType{},
		"status": ovhtypes.TfStringType{},
		"type":   ovhtypes.TfStringType{},
	}
}

// CloudFloatingIPCurrentStateAttrTypes returns the attribute types for the
// current_state object of a floating public IP.
func CloudFloatingIPCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"ip":          ovhtypes.TfStringType{},
		"status":      ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"network": types.ObjectType{
			AttrTypes: cloudFloatingIPNetworkAttrTypes(),
		},
		"associated_resource": types.ObjectType{
			AttrTypes: cloudPublicIPAssociatedResourceAttrTypes(),
		},
		"location": types.ObjectType{
			AttrTypes: cloudPublicIPLocationAttrTypes(),
		},
	}
}

// CloudExtNetIPCurrentStateAttrTypes returns the attribute types for the
// current_state object of an external-network public IP.
func CloudExtNetIPCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
		"ip": ovhtypes.TfStringType{},
		"associated_resource": types.ObjectType{
			AttrTypes: cloudPublicIPAssociatedResourceAttrTypes(),
		},
		"location": types.ObjectType{
			AttrTypes: cloudPublicIPLocationAttrTypes(),
		},
	}
}

// CloudAdditionalIPCurrentStateAttrTypes returns the attribute types for the
// current_state object of an additional public IP.
func CloudAdditionalIPCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       ovhtypes.TfStringType{},
		"ip":       ovhtypes.TfStringType{},
		"ip_block": ovhtypes.TfStringType{},
		"associated_resource": types.ObjectType{
			AttrTypes: cloudPublicIPAssociatedResourceAttrTypes(),
		},
		"location": types.ObjectType{
			AttrTypes: cloudPublicIPLocationAttrTypes(),
		},
	}
}

// Terraform object builders

// buildCloudPublicIPLocationObject constructs the nested location object from the API location.
func buildCloudPublicIPLocationObject(location *CloudPublicIPLocation) basetypes.ObjectValue {
	if location == nil {
		return types.ObjectNull(cloudPublicIPLocationAttrTypes())
	}

	locationObj, _ := types.ObjectValue(
		cloudPublicIPLocationAttrTypes(),
		map[string]attr.Value{
			"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(location.Region)},
			"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(location.AvailabilityZone)},
		},
	)

	return locationObj
}

// buildCloudPublicIPAssociatedResourceObject constructs the nested associated_resource
// object from the API associated resource.
func buildCloudPublicIPAssociatedResourceObject(associatedResource *CloudPublicIPAssociatedResource) basetypes.ObjectValue {
	if associatedResource == nil {
		return types.ObjectNull(cloudPublicIPAssociatedResourceAttrTypes())
	}

	associatedResourceObj, _ := types.ObjectValue(
		cloudPublicIPAssociatedResourceAttrTypes(),
		map[string]attr.Value{
			"id":   ovhtypes.TfStringValue{StringValue: types.StringValue(associatedResource.ID)},
			"type": ovhtypes.TfStringValue{StringValue: types.StringValue(associatedResource.Type)},
		},
	)

	return associatedResourceObj
}

// buildCloudFloatingIPNetworkObject constructs the nested network object from the API network.
func buildCloudFloatingIPNetworkObject(network *CloudFloatingIPNetwork) basetypes.ObjectValue {
	if network == nil {
		return types.ObjectNull(cloudFloatingIPNetworkAttrTypes())
	}

	networkObj, _ := types.ObjectValue(
		cloudFloatingIPNetworkAttrTypes(),
		map[string]attr.Value{
			"id": ovhtypes.TfStringValue{StringValue: types.StringValue(network.ID)},
		},
	)

	return networkObj
}

// buildCloudPublicIPTaskErrorsList constructs the errors list of a current_tasks entry
// from the API task errors.
func buildCloudPublicIPTaskErrorsList(taskErrors []CloudPublicIPTaskError) basetypes.ListValue {
	errorObjType := types.ObjectType{AttrTypes: cloudPublicIPTaskErrorAttrTypes()}

	if taskErrors == nil {
		return types.ListNull(errorObjType)
	}

	errorObjs := make([]attr.Value, len(taskErrors))
	for i, taskError := range taskErrors {
		errorObj, _ := types.ObjectValue(
			cloudPublicIPTaskErrorAttrTypes(),
			map[string]attr.Value{
				"message": ovhtypes.TfStringValue{StringValue: types.StringValue(taskError.Message)},
			},
		)
		errorObjs[i] = errorObj
	}

	errorsVal, _ := types.ListValue(errorObjType, errorObjs)

	return errorsVal
}

// buildCloudPublicIPCurrentTasksList constructs the current_tasks list from the API tasks.
func buildCloudPublicIPCurrentTasksList(tasks []CloudPublicIPResourceTask) basetypes.ListValue {
	taskObjType := types.ObjectType{AttrTypes: cloudPublicIPResourceTaskAttrTypes()}

	if tasks == nil {
		return types.ListNull(taskObjType)
	}

	taskObjs := make([]attr.Value, len(tasks))
	for i, task := range tasks {
		taskObj, _ := types.ObjectValue(
			cloudPublicIPResourceTaskAttrTypes(),
			map[string]attr.Value{
				"errors": buildCloudPublicIPTaskErrorsList(task.Errors),
				"id":     ovhtypes.TfStringValue{StringValue: types.StringValue(task.ID)},
				"link":   ovhtypes.TfStringValue{StringValue: types.StringValue(task.Link)},
				"status": ovhtypes.TfStringValue{StringValue: types.StringValue(task.Status)},
				"type":   ovhtypes.TfStringValue{StringValue: types.StringValue(task.Type)},
			},
		)
		taskObjs[i] = taskObj
	}

	tasksVal, _ := types.ListValue(taskObjType, taskObjs)

	return tasksVal
}

// buildCloudFloatingIPCurrentStateObject constructs the current_state object
// of a floating public IP from the API response.
func buildCloudFloatingIPCurrentStateObject(ctx context.Context, state *CloudFloatingIPCurrentState) basetypes.ObjectValue {
	currentStateObj, _ := types.ObjectValue(
		CloudFloatingIPCurrentStateAttrTypes(),
		map[string]attr.Value{
			"id":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.ID)},
			"ip":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.IP)},
			"status":              ovhtypes.TfStringValue{StringValue: types.StringValue(state.Status)},
			"description":         ovhtypes.TfStringValue{StringValue: types.StringValue(state.Description)},
			"network":             buildCloudFloatingIPNetworkObject(state.Network),
			"associated_resource": buildCloudPublicIPAssociatedResourceObject(state.AssociatedResource),
			"location":            buildCloudPublicIPLocationObject(state.Location),
		},
	)

	return currentStateObj
}

// buildCloudExtNetIPCurrentStateObject constructs the current_state object
// of an external-network public IP from the API response.
func buildCloudExtNetIPCurrentStateObject(ctx context.Context, state *CloudExtNetIPCurrentState) basetypes.ObjectValue {
	currentStateObj, _ := types.ObjectValue(
		CloudExtNetIPCurrentStateAttrTypes(),
		map[string]attr.Value{
			"id":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.ID)},
			"ip":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.IP)},
			"associated_resource": buildCloudPublicIPAssociatedResourceObject(state.AssociatedResource),
			"location":            buildCloudPublicIPLocationObject(state.Location),
		},
	)

	return currentStateObj
}

// buildCloudAdditionalIPCurrentStateObject constructs the current_state object
// of an additional public IP from the API response.
func buildCloudAdditionalIPCurrentStateObject(ctx context.Context, state *CloudAdditionalIPCurrentState) basetypes.ObjectValue {
	currentStateObj, _ := types.ObjectValue(
		CloudAdditionalIPCurrentStateAttrTypes(),
		map[string]attr.Value{
			"id":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.ID)},
			"ip":                  ovhtypes.TfStringValue{StringValue: types.StringValue(state.IP)},
			"ip_block":            ovhtypes.TfStringValue{StringValue: types.StringValue(state.IPBlock)},
			"associated_resource": buildCloudPublicIPAssociatedResourceObject(state.AssociatedResource),
			"location":            buildCloudPublicIPLocationObject(state.Location),
		},
	)

	return currentStateObj
}

// Terraform list-item attribute types for the plural datasources

// CloudFloatingIPListItemAttrTypes returns the attribute types for a single
// floating IP element of the plural datasource list.
func CloudFloatingIPListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"description": ovhtypes.TfStringType{},
		"location": types.ObjectType{
			AttrTypes: cloudPublicIPLocationAttrTypes(),
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: CloudFloatingIPCurrentStateAttrTypes(),
		},
		"current_tasks": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: cloudPublicIPResourceTaskAttrTypes(),
			},
		},
	}
}

// CloudExtNetIPListItemAttrTypes returns the attribute types for a single
// external-network IP element of the plural datasource list.
func CloudExtNetIPListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: CloudExtNetIPCurrentStateAttrTypes(),
		},
		"current_tasks": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: cloudPublicIPResourceTaskAttrTypes(),
			},
		},
	}
}

// CloudAdditionalIPListItemAttrTypes returns the attribute types for a single
// additional IP element of the plural datasource list.
func CloudAdditionalIPListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"checksum":        ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: CloudAdditionalIPCurrentStateAttrTypes(),
		},
		"current_tasks": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: cloudPublicIPResourceTaskAttrTypes(),
			},
		},
	}
}

// Terraform list-item builders for the plural datasources

// buildCloudFloatingIPListItemObject builds a single floating IP element of
// the plural datasource list from an API response, reusing the shared builders
// for nested objects.
func buildCloudFloatingIPListItemObject(ctx context.Context, response *CloudFloatingIPAPIResponse) basetypes.ObjectValue {
	descriptionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	locationVal := types.ObjectNull(cloudPublicIPLocationAttrTypes())

	if response.TargetSpec != nil {
		descriptionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		locationVal = buildCloudPublicIPLocationObject(response.TargetSpec.Location)
	}

	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildCloudFloatingIPCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(CloudFloatingIPCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		CloudFloatingIPListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.ID)},
			"description":     descriptionVal,
			"location":        locationVal,
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
			"current_tasks":   buildCloudPublicIPCurrentTasksList(response.CurrentTasks),
		},
	)

	return obj
}

// buildCloudExtNetIPListItemObject builds a single external-network IP element
// of the plural datasource list from an API response, reusing the shared builders
// for nested objects.
func buildCloudExtNetIPListItemObject(ctx context.Context, response *CloudExtNetIPAPIResponse) basetypes.ObjectValue {
	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildCloudExtNetIPCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(CloudExtNetIPCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		CloudExtNetIPListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.ID)},
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":      ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
			"current_tasks":   buildCloudPublicIPCurrentTasksList(response.CurrentTasks),
		},
	)

	return obj
}

// buildCloudAdditionalIPListItemObject builds a single additional IP element
// of the plural datasource list from an API response, reusing the shared builders
// for nested objects.
func buildCloudAdditionalIPListItemObject(ctx context.Context, response *CloudAdditionalIPAPIResponse) basetypes.ObjectValue {
	var currentStateVal attr.Value
	if response.CurrentState != nil {
		currentStateVal = buildCloudAdditionalIPCurrentStateObject(ctx, response.CurrentState)
	} else {
		currentStateVal = types.ObjectNull(CloudAdditionalIPCurrentStateAttrTypes())
	}

	obj, _ := types.ObjectValue(
		CloudAdditionalIPListItemAttrTypes(),
		map[string]attr.Value{
			"id":              ovhtypes.TfStringValue{StringValue: types.StringValue(response.ID)},
			"checksum":        ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"resource_status": ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":   currentStateVal,
			"current_tasks":   buildCloudPublicIPCurrentTasksList(response.CurrentTasks),
		},
	)

	return obj
}

// Terraform model for the floating IP resource

// CloudFloatingIPModel represents the Terraform model for the floating IP resource.
type CloudFloatingIPModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`

	// Optional — immutable
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`

	// Optional — mutable
	Description ovhtypes.TfStringValue `tfsdk:"description"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
	CurrentTasks   types.List             `tfsdk:"current_tasks"`
}

// ToCreate converts the Terraform model to the API create payload
func (m *CloudFloatingIPModel) ToCreate() *CloudFloatingIPCreatePayload {
	targetSpec := &CloudFloatingIPTargetSpec{
		Location: &CloudPublicIPLocation{
			Region: m.Region.ValueString(),
		},
	}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	// Handle optional availability zone
	if !m.AvailabilityZone.IsNull() && !m.AvailabilityZone.IsUnknown() {
		targetSpec.Location.AvailabilityZone = m.AvailabilityZone.ValueString()
	}

	return &CloudFloatingIPCreatePayload{TargetSpec: targetSpec}
}

// ToUpdate converts the Terraform model to the API update payload
// Note: location is immutable, description is the only mutable field
func (m *CloudFloatingIPModel) ToUpdate(checksum string) *CloudFloatingIPUpdatePayload {
	targetSpec := &CloudFloatingIPUpdateTargetSpec{}

	// Handle optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		targetSpec.Description = m.Description.ValueString()
	}

	return &CloudFloatingIPUpdatePayload{
		Checksum:   checksum,
		TargetSpec: targetSpec,
	}
}

// MergeWith merges API response data into the Terraform model
func (m *CloudFloatingIPModel) MergeWith(ctx context.Context, response *CloudFloatingIPAPIResponse) {
	// Never overwrite the tracked IP with an empty response ID
	if response.ID != "" {
		m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ID)}
	}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	// Build current_state from API currentState
	if response.CurrentState != nil {
		m.CurrentState = buildCloudFloatingIPCurrentStateObject(ctx, response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(CloudFloatingIPCurrentStateAttrTypes())
	}

	m.CurrentTasks = buildCloudPublicIPCurrentTasksList(response.CurrentTasks)

	// Set flattened root-level fields from targetSpec
	if response.TargetSpec != nil {
		// Keep description null if user didn't set it and API returns empty
		if response.TargetSpec.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
			m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Description)}
		}

		if response.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.Region)}
			if response.TargetSpec.Location.AvailabilityZone != "" {
				m.AvailabilityZone = ovhtypes.TfStringValue{StringValue: types.StringValue(response.TargetSpec.Location.AvailabilityZone)}
			}
		}
	}
}
