package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudInstanceAutobackupModel is the Terraform model for the autobackup resource.
type CloudInstanceAutobackupModel struct {
	// Required — immutable (all trigger replace)
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	ImageName   ovhtypes.TfStringValue `tfsdk:"image_name"`
	Cron        ovhtypes.TfStringValue `tfsdk:"cron"`
	Rotation    ovhtypes.TfInt64Value  `tfsdk:"rotation"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	InstanceId  ovhtypes.TfStringValue `tfsdk:"instance_id"`

	// Optional — immutable
	Distant types.Object `tfsdk:"distant"`

	// Computed
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API response types (match the JSON envelope from the API).

type CloudInstanceAutobackupAPIResponse struct {
	Id             string                                  `json:"id"`
	Checksum       string                                  `json:"checksum"`
	ResourceStatus string                                  `json:"resourceStatus"`
	CreatedAt      string                                  `json:"createdAt"`
	UpdatedAt      string                                  `json:"updatedAt"`
	TargetSpec     *CloudInstanceAutobackupAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudInstanceAutobackupAPICurrentState `json:"currentState,omitempty"`
}

type CloudInstanceAutobackupAPITargetSpec struct {
	Name      string                              `json:"name"`
	ImageName string                              `json:"imageName"`
	Cron      string                              `json:"cron"`
	Rotation  int64                               `json:"rotation"`
	Location  *CloudInstanceAutobackupAPILocation `json:"location"`
	Instance  *CloudInstanceAutobackupAPIRef      `json:"instance"`
	Distant   *CloudInstanceAutobackupAPIDistant  `json:"distant,omitempty"`
}

type CloudInstanceAutobackupAPICurrentState struct {
	Name              string                                `json:"name"`
	ImageName         string                                `json:"imageName"`
	Cron              string                                `json:"cron"`
	Rotation          int64                                 `json:"rotation"`
	Location          *CloudInstanceAutobackupAPILocation   `json:"location"`
	Instance          *CloudInstanceAutobackupAPIRef        `json:"instance"`
	WorkflowName      string                                `json:"workflowName,omitempty"`
	NextExecutionTime string                                `json:"nextExecutionTime,omitempty"`
	Distant           *CloudInstanceAutobackupAPIDistant    `json:"distant,omitempty"`
	LastExecutions    []CloudInstanceAutobackupAPIExecution `json:"lastExecutions"`
}

type CloudInstanceAutobackupAPIExecution struct {
	Id           string  `json:"id"`
	StartedAt    *string `json:"startedAt"`
	UpdatedAt    *string `json:"updatedAt"`
	State        string  `json:"state"`
	ErrorMessage *string `json:"errorMessage"`
}

type CloudInstanceAutobackupAPILocation struct {
	Region string `json:"region"`
}

type CloudInstanceAutobackupAPIRef struct {
	Id string `json:"id"`
}

type CloudInstanceAutobackupAPIDistant struct {
	Region    string `json:"region"`
	ImageName string `json:"imageName"`
}

// CloudInstanceAutobackupCreatePayload wraps fields into the API targetSpec envelope.
type CloudInstanceAutobackupCreatePayload struct {
	TargetSpec *CloudInstanceAutobackupAPITargetSpec `json:"targetSpec"`
}

// cloudInstanceAutobackupDistantAttrTypes returns the attribute types for the distant nested object.
func cloudInstanceAutobackupDistantAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region":     ovhtypes.TfStringType{},
		"image_name": ovhtypes.TfStringType{},
	}
}

// ToCreate converts the Terraform model into the API create payload.
func (m *CloudInstanceAutobackupModel) ToCreate() *CloudInstanceAutobackupCreatePayload {
	spec := &CloudInstanceAutobackupAPITargetSpec{
		Name:      m.Name.ValueString(),
		ImageName: m.ImageName.ValueString(),
		Cron:      m.Cron.ValueString(),
		Rotation:  m.Rotation.ValueInt64(),
		Location:  &CloudInstanceAutobackupAPILocation{Region: m.Region.ValueString()},
		Instance:  &CloudInstanceAutobackupAPIRef{Id: m.InstanceId.ValueString()},
	}

	if !m.Distant.IsNull() && !m.Distant.IsUnknown() {
		distantAttrs := m.Distant.Attributes()
		region := distantAttrs["region"].(ovhtypes.TfStringValue).ValueString()
		imageName := distantAttrs["image_name"].(ovhtypes.TfStringValue).ValueString()
		spec.Distant = &CloudInstanceAutobackupAPIDistant{
			Region:    region,
			ImageName: imageName,
		}
	}

	return &CloudInstanceAutobackupCreatePayload{TargetSpec: spec}
}

// MergeWith merges API response data into the Terraform model.
func (m *CloudInstanceAutobackupModel) MergeWith(ctx context.Context, resp *CloudInstanceAutobackupAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.Checksum)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.ResourceStatus)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.UpdatedAt)}

	// Set root-level fields from targetSpec
	if resp.TargetSpec != nil {
		m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Name)}
		m.ImageName = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.ImageName)}
		m.Cron = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Cron)}
		m.Rotation = ovhtypes.TfInt64Value{Int64Value: types.Int64Value(resp.TargetSpec.Rotation)}
		if resp.TargetSpec.Location != nil {
			m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Location.Region)}
		}
		if resp.TargetSpec.Instance != nil {
			m.InstanceId = ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Instance.Id)}
		}
		// ServiceName is not returned by the API — kept from config
		if resp.TargetSpec.Distant != nil {
			distObj, _ := types.ObjectValue(cloudInstanceAutobackupDistantAttrTypes(), map[string]attr.Value{
				"region":     ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Distant.Region)},
				"image_name": ovhtypes.TfStringValue{StringValue: types.StringValue(resp.TargetSpec.Distant.ImageName)},
			})
			m.Distant = distObj
		} else {
			m.Distant = types.ObjectNull(cloudInstanceAutobackupDistantAttrTypes())
		}
	}

	// Build current_state object
	if resp.CurrentState != nil {
		m.CurrentState = buildAutobackupCurrentStateObject(resp.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(cloudInstanceAutobackupCurrentStateAttrTypes())
	}
}

// cloudInstanceAutobackupExecutionAttrTypes returns the attribute types for a single last_executions entry.
func cloudInstanceAutobackupExecutionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            ovhtypes.TfStringType{},
		"started_at":    ovhtypes.TfStringType{},
		"updated_at":    ovhtypes.TfStringType{},
		"state":         ovhtypes.TfStringType{},
		"error_message": ovhtypes.TfStringType{},
	}
}

// cloudInstanceAutobackupCurrentStateAttrTypes returns the attribute types map for building current_state objects.
func cloudInstanceAutobackupCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                ovhtypes.TfStringType{},
		"image_name":          ovhtypes.TfStringType{},
		"cron":                ovhtypes.TfStringType{},
		"rotation":            types.Int64Type,
		"region":              ovhtypes.TfStringType{},
		"instance_id":         ovhtypes.TfStringType{},
		"workflow_name":       ovhtypes.TfStringType{},
		"next_execution_time": ovhtypes.TfStringType{},
		"distant": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"region":     ovhtypes.TfStringType{},
				"image_name": ovhtypes.TfStringType{},
			},
		},
		"last_executions": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: cloudInstanceAutobackupExecutionAttrTypes(),
			},
		},
	}
}

// nullableTfString converts a *string from the API into the matching ovhtypes.TfStringValue,
// using a null TF value when the pointer is nil.
func nullableTfString(s *string) ovhtypes.TfStringValue {
	if s == nil {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return ovhtypes.TfStringValue{StringValue: types.StringValue(*s)}
}

func buildAutobackupCurrentStateObject(state *CloudInstanceAutobackupAPICurrentState) types.Object {
	region := ""
	if state.Location != nil {
		region = state.Location.Region
	}

	instanceId := ""
	if state.Instance != nil {
		instanceId = state.Instance.Id
	}

	attrs := map[string]attr.Value{
		"name":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Name)},
		"image_name":          ovhtypes.TfStringValue{StringValue: types.StringValue(state.ImageName)},
		"cron":                ovhtypes.TfStringValue{StringValue: types.StringValue(state.Cron)},
		"rotation":            types.Int64Value(state.Rotation),
		"region":              ovhtypes.TfStringValue{StringValue: types.StringValue(region)},
		"instance_id":         ovhtypes.TfStringValue{StringValue: types.StringValue(instanceId)},
		"workflow_name":       ovhtypes.TfStringValue{StringValue: types.StringValue(state.WorkflowName)},
		"next_execution_time": ovhtypes.TfStringValue{StringValue: types.StringValue(state.NextExecutionTime)},
	}

	if state.Distant != nil {
		distObj, _ := types.ObjectValue(
			cloudInstanceAutobackupDistantAttrTypes(),
			map[string]attr.Value{
				"region":     ovhtypes.TfStringValue{StringValue: types.StringValue(state.Distant.Region)},
				"image_name": ovhtypes.TfStringValue{StringValue: types.StringValue(state.Distant.ImageName)},
			},
		)
		attrs["distant"] = distObj
	} else {
		attrs["distant"] = types.ObjectNull(cloudInstanceAutobackupDistantAttrTypes())
	}

	// Build last_executions list (null when the API returned null/absent — typical of LIST or Mistral failure).
	executionElemType := types.ObjectType{AttrTypes: cloudInstanceAutobackupExecutionAttrTypes()}
	var executionsVal basetypes.ListValue
	if state.LastExecutions != nil {
		execObjs := make([]attr.Value, len(state.LastExecutions))
		for i, exec := range state.LastExecutions {
			execObj, _ := types.ObjectValue(
				cloudInstanceAutobackupExecutionAttrTypes(),
				map[string]attr.Value{
					"id":            ovhtypes.TfStringValue{StringValue: types.StringValue(exec.Id)},
					"started_at":    nullableTfString(exec.StartedAt),
					"updated_at":    nullableTfString(exec.UpdatedAt),
					"state":         ovhtypes.TfStringValue{StringValue: types.StringValue(exec.State)},
					"error_message": nullableTfString(exec.ErrorMessage),
				},
			)
			execObjs[i] = execObj
		}
		executionsVal, _ = types.ListValue(executionElemType, execObjs)
	} else {
		executionsVal = types.ListNull(executionElemType)
	}
	attrs["last_executions"] = executionsVal

	obj, _ := types.ObjectValue(cloudInstanceAutobackupCurrentStateAttrTypes(), attrs)
	return obj
}
