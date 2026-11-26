package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudManagedKubernetesNodepoolModel is shared by the
// ovh_cloud_managed_kubernetes_nodepool resource and data source.
type CloudManagedKubernetesNodepoolModel struct {
	// Required — immutable
	ServiceName         ovhtypes.TfStringValue `tfsdk:"service_name"`
	ManagedKubernetesId ovhtypes.TfStringValue `tfsdk:"managed_kubernetes_id"`
	NodeConfig          types.Object           `tfsdk:"node_config"`

	// Optional — immutable
	Name             ovhtypes.TfStringValue `tfsdk:"name"`
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone"`
	BillingType      ovhtypes.TfStringValue `tfsdk:"billing_type"`
	AntiAffinity     types.Object           `tfsdk:"anti_affinity"`

	// Optional — mutable
	DesiredNodes ovhtypes.TfInt64Value `tfsdk:"desired_nodes"`
	Autoscaler   types.Object          `tfsdk:"autoscaler"`
	Template     types.Object          `tfsdk:"template"`

	// Computed envelope
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// ---------- Nested objects (API DTO and Terraform object) ----------

type MKSAttachFloatingIp struct {
	Enabled ovhtypes.TfBoolValue `tfsdk:"enabled" json:"enabled"`
}

type MKSNodeConfig struct {
	Flavor           ovhtypes.TfStringValue `tfsdk:"flavor" json:"flavor"`
	AttachFloatingIp *MKSAttachFloatingIp   `tfsdk:"attach_floating_ip" json:"attachFloatingIp,omitempty"`
}

type MKSAntiAffinity struct {
	Enabled ovhtypes.TfBoolValue   `tfsdk:"enabled" json:"enabled"`
	Policy  ovhtypes.TfStringValue `tfsdk:"policy" json:"policy"`
}

type MKSAutoscalerScaleDown struct {
	UnneededTimeSeconds  ovhtypes.TfInt64Value  `tfsdk:"unneeded_time_seconds" json:"unneededTimeSeconds"`
	UnreadyTimeSeconds   ovhtypes.TfInt64Value  `tfsdk:"unready_time_seconds" json:"unreadyTimeSeconds"`
	UtilizationThreshold ovhtypes.TfNumberValue `tfsdk:"utilization_threshold" json:"utilizationThreshold"`
}

type MKSAutoscaler struct {
	Enabled   ovhtypes.TfBoolValue    `tfsdk:"enabled" json:"enabled"`
	MinNodes  ovhtypes.TfInt64Value   `tfsdk:"min_nodes" json:"minNodes"`
	MaxNodes  ovhtypes.TfInt64Value   `tfsdk:"max_nodes" json:"maxNodes"`
	ScaleDown *MKSAutoscalerScaleDown `tfsdk:"scale_down" json:"scaleDown,omitempty"`
}

type MKSTaint struct {
	Effect ovhtypes.TfStringValue `tfsdk:"effect" json:"effect"`
	Key    ovhtypes.TfStringValue `tfsdk:"key" json:"key"`
	Value  ovhtypes.TfStringValue `tfsdk:"value" json:"value"`
}

// Collections have no omitempty so an empty list or map clears the API value.
type MKSNodeTemplateMetadata struct {
	Annotations map[string]ovhtypes.TfStringValue `tfsdk:"annotations" json:"annotations"`
	Finalizers  []ovhtypes.TfStringValue          `tfsdk:"finalizers" json:"finalizers"`
	Labels      map[string]ovhtypes.TfStringValue `tfsdk:"labels" json:"labels"`
}

type MKSNodeTemplateSpec struct {
	Taints        []MKSTaint           `tfsdk:"taints" json:"taints"`
	Unschedulable ovhtypes.TfBoolValue `tfsdk:"unschedulable" json:"unschedulable"`
}

type MKSNodeTemplate struct {
	Metadata *MKSNodeTemplateMetadata `tfsdk:"metadata" json:"metadata,omitempty"`
	Spec     *MKSNodeTemplateSpec     `tfsdk:"spec" json:"spec,omitempty"`
}

// normalizeEmpty maps omitted collections to empty ones, since the API drops empty values that may be configured as {} or [].
func (t *MKSNodeTemplate) normalizeEmpty() {
	if t == nil {
		return
	}
	if md := t.Metadata; md != nil {
		if md.Annotations == nil {
			md.Annotations = map[string]ovhtypes.TfStringValue{}
		}
		if md.Labels == nil {
			md.Labels = map[string]ovhtypes.TfStringValue{}
		}
		if md.Finalizers == nil {
			md.Finalizers = []ovhtypes.TfStringValue{}
		}
	}
	if t.Spec != nil && t.Spec.Taints == nil {
		t.Spec.Taints = []MKSTaint{}
	}
}

type MKSNodepoolRollout struct {
	Strategy ovhtypes.TfStringValue `tfsdk:"strategy" json:"strategy"`
}

type MKSNodepoolStatus struct {
	CurrentNodes  ovhtypes.TfInt64Value `tfsdk:"current_nodes" json:"currentNodes"`
	ReadyNodes    ovhtypes.TfInt64Value `tfsdk:"ready_nodes" json:"readyNodes"`
	UpToDateNodes ovhtypes.TfInt64Value `tfsdk:"up_to_date_nodes" json:"upToDateNodes"`
}

// ---------- API envelope ----------

type CloudManagedKubernetesNodepoolAPITargetSpec struct {
	Name             ovhtypes.TfStringValue `json:"name"`
	AvailabilityZone ovhtypes.TfStringValue `json:"availabilityZone"`
	BillingType      ovhtypes.TfStringValue `json:"billingType"`
	DesiredNodes     *int64                 `json:"desiredNodes,omitempty"`
	NodeConfig       *MKSNodeConfig         `json:"nodeConfig,omitempty"`
	AntiAffinity     *MKSAntiAffinity       `json:"antiAffinity,omitempty"`
	Autoscaler       *MKSAutoscaler         `json:"autoscaler,omitempty"`
	Template         *MKSNodeTemplate       `json:"template,omitempty"`
}

// The update spec carries mutable fields only.
type CloudManagedKubernetesNodepoolAPIUpdateTargetSpec struct {
	DesiredNodes *int64           `json:"desiredNodes,omitempty"`
	Autoscaler   *MKSAutoscaler   `json:"autoscaler,omitempty"`
	Template     *MKSNodeTemplate `json:"template,omitempty"`
}

type CloudManagedKubernetesNodepoolAPICurrentState struct {
	Name             ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone" json:"availabilityZone"`
	BillingType      ovhtypes.TfStringValue `tfsdk:"billing_type" json:"billingType"`
	DesiredNodes     ovhtypes.TfInt64Value  `tfsdk:"desired_nodes" json:"desiredNodes"`
	NodeConfig       *MKSNodeConfig         `tfsdk:"node_config" json:"nodeConfig,omitempty"`
	AntiAffinity     *MKSAntiAffinity       `tfsdk:"anti_affinity" json:"antiAffinity,omitempty"`
	Autoscaler       *MKSAutoscaler         `tfsdk:"autoscaler" json:"autoscaler,omitempty"`
	Template         *MKSNodeTemplate       `tfsdk:"template" json:"template,omitempty"`
	Rollout          *MKSNodepoolRollout    `tfsdk:"rollout" json:"rollout,omitempty"`
	Status           *MKSNodepoolStatus     `tfsdk:"status" json:"status,omitempty"`
}

type CloudManagedKubernetesNodepoolAPIResponse struct {
	Id             string                                         `json:"id"`
	Checksum       string                                         `json:"checksum"`
	CreatedAt      string                                         `json:"createdAt"`
	UpdatedAt      string                                         `json:"updatedAt"`
	ResourceStatus string                                         `json:"resourceStatus"`
	TargetSpec     *CloudManagedKubernetesNodepoolAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudManagedKubernetesNodepoolAPICurrentState `json:"currentState,omitempty"`
	CurrentTasks   []CloudResourceTask                            `json:"currentTasks,omitempty"`
}

type CloudManagedKubernetesNodepoolCreatePayload struct {
	TargetSpec CloudManagedKubernetesNodepoolAPITargetSpec `json:"targetSpec"`
}

type CloudManagedKubernetesNodepoolUpdatePayload struct {
	Checksum   string                                            `json:"checksum"`
	TargetSpec CloudManagedKubernetesNodepoolAPIUpdateTargetSpec `json:"targetSpec"`
}

// ---------- Conversions ----------

func (m *CloudManagedKubernetesNodepoolModel) ToCreate(ctx context.Context, diags *diag.Diagnostics) *CloudManagedKubernetesNodepoolCreatePayload {
	return &CloudManagedKubernetesNodepoolCreatePayload{TargetSpec: CloudManagedKubernetesNodepoolAPITargetSpec{
		Name:             m.Name,
		AvailabilityZone: m.AvailabilityZone,
		BillingType:      m.BillingType,
		DesiredNodes:     knownInt64Pointer(m.DesiredNodes),
		NodeConfig:       objectAs[MKSNodeConfig](ctx, m.NodeConfig, diags),
		AntiAffinity:     objectAs[MKSAntiAffinity](ctx, m.AntiAffinity, diags),
		Autoscaler:       objectAs[MKSAutoscaler](ctx, m.Autoscaler, diags),
		Template:         objectAs[MKSNodeTemplate](ctx, m.Template, diags),
	}}
}

func (m *CloudManagedKubernetesNodepoolModel) ToUpdate(ctx context.Context, checksum string, diags *diag.Diagnostics) *CloudManagedKubernetesNodepoolUpdatePayload {
	return &CloudManagedKubernetesNodepoolUpdatePayload{Checksum: checksum, TargetSpec: CloudManagedKubernetesNodepoolAPIUpdateTargetSpec{
		DesiredNodes: knownInt64Pointer(m.DesiredNodes),
		Autoscaler:   objectAs[MKSAutoscaler](ctx, m.Autoscaler, diags),
		Template:     objectAs[MKSNodeTemplate](ctx, m.Template, diags),
	}}
}

// MergeWith copies the API response into the model.
// service_name and managed_kubernetes_id come from the route and are left untouched.
func (m *CloudManagedKubernetesNodepoolModel) MergeWith(ctx context.Context, response *CloudManagedKubernetesNodepoolAPIResponse, diags *diag.Diagnostics) {
	attrTypes := managedKubernetesNodepoolAttrTypes()

	m.Id = tfStr(response.Id)
	m.Checksum = tfStr(response.Checksum)
	m.CreatedAt = tfStr(response.CreatedAt)
	m.UpdatedAt = tfStr(response.UpdatedAt)
	m.ResourceStatus = tfStr(response.ResourceStatus)
	m.CurrentState = objectFrom(ctx, attrTypes["current_state"], response.CurrentState, diags)

	ts := response.TargetSpec
	if ts == nil {
		ts = &CloudManagedKubernetesNodepoolAPITargetSpec{}
	}
	m.Name = ts.Name
	m.BillingType = ts.BillingType
	m.DesiredNodes = ovhtypes.TfInt64Value{Int64Value: types.Int64PointerValue(ts.DesiredNodes)}
	m.AvailabilityZone = ts.AvailabilityZone
	m.NodeConfig = objectFrom(ctx, attrTypes["node_config"], ts.NodeConfig, diags)
	m.AntiAffinity = objectFrom(ctx, attrTypes["anti_affinity"], ts.AntiAffinity, diags)
	m.Autoscaler = objectFrom(ctx, attrTypes["autoscaler"], ts.Autoscaler, diags)
	ts.Template.normalizeEmpty()
	m.Template = objectFrom(ctx, attrTypes["template"], ts.Template, diags)
}

// knownInt64Pointer omits an unset value from the payload so the API applies its default.
func knownInt64Pointer(v ovhtypes.TfInt64Value) *int64 {
	if v.IsUnknown() {
		return nil
	}
	return v.ValueInt64Pointer()
}
