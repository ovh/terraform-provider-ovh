package ovh

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudManagedKubernetesModel is shared by the ovh_cloud_managed_kubernetes
// resource and data source. TargetSpec objects are flattened to root attributes.
type CloudManagedKubernetesModel struct {
	// Required — immutable
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Plan        ovhtypes.TfStringValue `tfsdk:"plan"`

	// Optional — immutable
	IpAllocationPolicy types.Object `tfsdk:"ip_allocation_policy"`

	// Required — partially mutable
	Cni     types.Object `tfsdk:"cni"`
	Network types.Object `tfsdk:"network"`

	// Optional — mutable
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Version        ovhtypes.TfStringValue `tfsdk:"version"`
	ApiServer      types.Object           `tfsdk:"api_server"`
	Authentication types.Object           `tfsdk:"authentication"`
	AutoUpgrade    types.Object           `tfsdk:"auto_upgrade"`

	// Computed envelope
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// ---------- Nested objects ----------
// Each struct is both the API DTO (json) and the Terraform object (tfsdk):
// a nil pointer is an absent JSON key and a null Terraform object.

type MKSAuthorizedIpRule struct {
	Cidr ovhtypes.TfStringValue `tfsdk:"cidr" json:"cidr"`
}

type MKSAuthorizedIps struct {
	Rules []MKSAuthorizedIpRule `tfsdk:"rules" json:"rules"`
}

type MKSApiServerAccess struct {
	AuthorizedIps *MKSAuthorizedIps `tfsdk:"authorized_ips" json:"authorizedIps,omitempty"`
}

type MKSAdmissionPlugin struct {
	Enabled ovhtypes.TfBoolValue `tfsdk:"enabled" json:"enabled"`
}

type MKSApiServer struct {
	Access           *MKSApiServerAccess           `tfsdk:"access" json:"access,omitempty"`
	AdmissionPlugins map[string]MKSAdmissionPlugin `tfsdk:"admission_plugins" json:"admissionPlugins,omitempty"`
}

type MKSOidc struct {
	CaContent      ovhtypes.TfStringValue   `tfsdk:"ca_content" json:"caContent"`
	ClientId       ovhtypes.TfStringValue   `tfsdk:"client_id" json:"clientId"`
	GroupsClaim    []ovhtypes.TfStringValue `tfsdk:"groups_claim" json:"groupsClaim,omitempty"`
	GroupsPrefix   ovhtypes.TfStringValue   `tfsdk:"groups_prefix" json:"groupsPrefix"`
	IssuerUrl      ovhtypes.TfStringValue   `tfsdk:"issuer_url" json:"issuerUrl"`
	RequiredClaim  []ovhtypes.TfStringValue `tfsdk:"required_claim" json:"requiredClaim,omitempty"`
	SigningAlgs    []ovhtypes.TfStringValue `tfsdk:"signing_algs" json:"signingAlgs,omitempty"`
	UsernameClaim  ovhtypes.TfStringValue   `tfsdk:"username_claim" json:"usernameClaim"`
	UsernamePrefix ovhtypes.TfStringValue   `tfsdk:"username_prefix" json:"usernamePrefix"`
}

type MKSAuthentication struct {
	Type ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
	Oidc *MKSOidc               `tfsdk:"oidc" json:"oidc,omitempty"`
}

type MKSAutoUpgradeTrack struct {
	Enabled ovhtypes.TfBoolValue `tfsdk:"enabled" json:"enabled"`
}

type MKSAutoUpgrade struct {
	Patch *MKSAutoUpgradeTrack `tfsdk:"patch" json:"patch,omitempty"`
}

type MKSKubeProxyIptables struct {
	MinSyncPeriod ovhtypes.TfISO8601DurationValue `tfsdk:"min_sync_period" json:"minSyncPeriod"`
	SyncPeriod    ovhtypes.TfISO8601DurationValue `tfsdk:"sync_period" json:"syncPeriod"`
}

type MKSKubeProxyIpvs struct {
	MinSyncPeriod ovhtypes.TfISO8601DurationValue `tfsdk:"min_sync_period" json:"minSyncPeriod"`
	Scheduler     ovhtypes.TfStringValue          `tfsdk:"scheduler" json:"scheduler"`
	SyncPeriod    ovhtypes.TfISO8601DurationValue `tfsdk:"sync_period" json:"syncPeriod"`
	TcpFinTimeout ovhtypes.TfISO8601DurationValue `tfsdk:"tcp_fin_timeout" json:"tcpFinTimeout"`
	TcpTimeout    ovhtypes.TfISO8601DurationValue `tfsdk:"tcp_timeout" json:"tcpTimeout"`
	UdpTimeout    ovhtypes.TfISO8601DurationValue `tfsdk:"udp_timeout" json:"udpTimeout"`
}

type MKSKubeProxy struct {
	Mode     ovhtypes.TfStringValue `tfsdk:"mode" json:"mode"`
	Iptables *MKSKubeProxyIptables  `tfsdk:"iptables" json:"iptables,omitempty"`
	Ipvs     *MKSKubeProxyIpvs      `tfsdk:"ipvs" json:"ipvs,omitempty"`
}

type MKSCniCanal struct {
	KubeProxy *MKSKubeProxy `tfsdk:"kube_proxy" json:"kubeProxy,omitempty"`
}

type MKSCiliumResourceList struct {
	Cpu    ovhtypes.TfStringValue `tfsdk:"cpu" json:"cpu"`
	Memory ovhtypes.TfStringValue `tfsdk:"memory" json:"memory"`
}

type MKSCiliumResourceRequirements struct {
	Limits   *MKSCiliumResourceList `tfsdk:"limits" json:"limits,omitempty"`
	Requests *MKSCiliumResourceList `tfsdk:"requests" json:"requests,omitempty"`
}

type MKSCiliumHubbleUi struct {
	Enabled           ovhtypes.TfBoolValue           `tfsdk:"enabled" json:"enabled"`
	BackendResources  *MKSCiliumResourceRequirements `tfsdk:"backend_resources" json:"backendResources,omitempty"`
	FrontendResources *MKSCiliumResourceRequirements `tfsdk:"frontend_resources" json:"frontendResources,omitempty"`
}

type MKSCiliumHubbleRelay struct {
	Enabled ovhtypes.TfBoolValue `tfsdk:"enabled" json:"enabled"`
}

type MKSCiliumHubble struct {
	Enabled ovhtypes.TfBoolValue  `tfsdk:"enabled" json:"enabled"`
	Relay   *MKSCiliumHubbleRelay `tfsdk:"relay" json:"relay,omitempty"`
	Ui      *MKSCiliumHubbleUi    `tfsdk:"ui" json:"ui,omitempty"`
}

type MKSCiliumClusterMeshApiServer struct {
	NodePort    ovhtypes.TfInt64Value  `tfsdk:"node_port" json:"nodePort"`
	ServiceType ovhtypes.TfStringValue `tfsdk:"service_type" json:"serviceType"`
}

type MKSCiliumClusterMesh struct {
	Enabled   ovhtypes.TfBoolValue           `tfsdk:"enabled" json:"enabled"`
	ApiServer *MKSCiliumClusterMeshApiServer `tfsdk:"api_server" json:"apiServer,omitempty"`
}

type MKSCniCilium struct {
	ClusterId   ovhtypes.TfInt64Value `tfsdk:"cluster_id" json:"clusterId"`
	ClusterMesh *MKSCiliumClusterMesh `tfsdk:"cluster_mesh" json:"clusterMesh,omitempty"`
	Hubble      *MKSCiliumHubble      `tfsdk:"hubble" json:"hubble,omitempty"`
}

type MKSCni struct {
	Type   ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
	Canal  *MKSCniCanal           `tfsdk:"canal" json:"canal,omitempty"`
	Cilium *MKSCniCilium          `tfsdk:"cilium" json:"cilium,omitempty"`
}

type MKSIpAllocationPolicy struct {
	Pods     ovhtypes.TfStringValue `tfsdk:"pods" json:"pods"`
	Services ovhtypes.TfStringValue `tfsdk:"services" json:"services"`
}

type MKSUserGateway struct {
	Ip ovhtypes.TfStringValue `tfsdk:"ip" json:"ip"`
}

type MKSDefaultEgress struct {
	Type        ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
	UserGateway *MKSUserGateway        `tfsdk:"user_gateway" json:"userGateway,omitempty"`
}

type MKSNetworkVrack struct {
	PrivateNetworkId     ovhtypes.TfStringValue `tfsdk:"private_network_id" json:"privateNetworkId"`
	SubnetId             ovhtypes.TfStringValue `tfsdk:"subnet_id" json:"subnetId"`
	LoadBalancerSubnetId ovhtypes.TfStringValue `tfsdk:"load_balancer_subnet_id" json:"loadBalancerSubnetId"`
	DefaultEgress        *MKSDefaultEgress      `tfsdk:"default_egress" json:"defaultEgress,omitempty"`
}

type MKSNetwork struct {
	Type  ovhtypes.TfStringValue `tfsdk:"type" json:"type"`
	Vrack *MKSNetworkVrack       `tfsdk:"vrack" json:"vrack,omitempty"`
}

// ---------- API envelope ----------

type CloudManagedKubernetesAPITargetSpec struct {
	Name               ovhtypes.TfStringValue `json:"name"`
	Version            ovhtypes.TfStringValue `json:"version"`
	Region             ovhtypes.TfStringValue `json:"region"`
	Plan               ovhtypes.TfStringValue `json:"plan"`
	ApiServer          *MKSApiServer          `json:"apiServer,omitempty"`
	Authentication     *MKSAuthentication     `json:"authentication,omitempty"`
	AutoUpgrade        *MKSAutoUpgrade        `json:"autoUpgrade,omitempty"`
	Cni                *MKSCni                `json:"cni,omitempty"`
	IpAllocationPolicy *MKSIpAllocationPolicy `json:"ipAllocationPolicy,omitempty"`
	Network            *MKSNetwork            `json:"network,omitempty"`
}

// Update DTOs carry mutable fields only: the CNI/network variants, the
// kube-proxy mode and the network identifiers are immutable and must not be sent.
type CloudManagedKubernetesAPIUpdateTargetSpec struct {
	Name           ovhtypes.TfStringValue `json:"name"`
	Version        ovhtypes.TfStringValue `json:"version"`
	ApiServer      *MKSApiServer          `json:"apiServer,omitempty"`
	Authentication *MKSAuthentication     `json:"authentication,omitempty"`
	AutoUpgrade    *MKSAutoUpgrade        `json:"autoUpgrade,omitempty"`
	Cni            *MKSCniUpdate          `json:"cni,omitempty"`
	Network        *MKSNetworkUpdate      `json:"network,omitempty"`
}

type MKSKubeProxyUpdate struct {
	Iptables *MKSKubeProxyIptables `json:"iptables,omitempty"`
	Ipvs     *MKSKubeProxyIpvs     `json:"ipvs,omitempty"`
}

type MKSCniCanalUpdate struct {
	KubeProxy *MKSKubeProxyUpdate `json:"kubeProxy,omitempty"`
}

type MKSCniUpdate struct {
	Canal  *MKSCniCanalUpdate `json:"canal,omitempty"`
	Cilium *MKSCniCilium      `json:"cilium,omitempty"`
}

type MKSNetworkVrackUpdate struct {
	DefaultEgress *MKSDefaultEgress `json:"defaultEgress,omitempty"`
}

type MKSNetworkUpdate struct {
	Vrack *MKSNetworkVrackUpdate `json:"vrack,omitempty"`
}

type MKSApiServerPublicEndpoint struct {
	Url ovhtypes.TfStringValue `tfsdk:"url" json:"url"`
}

type MKSCurrentStateApiServerAccess struct {
	AuthorizedIps  *MKSAuthorizedIps           `tfsdk:"authorized_ips" json:"authorizedIps,omitempty"`
	PublicEndpoint *MKSApiServerPublicEndpoint `tfsdk:"public_endpoint" json:"publicEndpoint,omitempty"`
}

type MKSCurrentStateApiServer struct {
	Access           *MKSCurrentStateApiServerAccess `tfsdk:"access" json:"access,omitempty"`
	AdmissionPlugins map[string]MKSAdmissionPlugin   `tfsdk:"admission_plugins" json:"admissionPlugins,omitempty"`
}

type CloudManagedKubernetesAPICurrentState struct {
	Name               ovhtypes.TfStringValue    `tfsdk:"name" json:"name"`
	Region             ovhtypes.TfStringValue    `tfsdk:"region" json:"region"`
	Plan               ovhtypes.TfStringValue    `tfsdk:"plan" json:"plan"`
	Version            ovhtypes.TfStringValue    `tfsdk:"version" json:"version"`
	DeployedVersion    ovhtypes.TfStringValue    `tfsdk:"deployed_version" json:"deployedVersion"`
	ApiServer          *MKSCurrentStateApiServer `tfsdk:"api_server" json:"apiServer,omitempty"`
	Authentication     *MKSAuthentication        `tfsdk:"authentication" json:"authentication,omitempty"`
	AutoUpgrade        *MKSAutoUpgrade           `tfsdk:"auto_upgrade" json:"autoUpgrade,omitempty"`
	Cni                *MKSCni                   `tfsdk:"cni" json:"cni,omitempty"`
	IpAllocationPolicy *MKSIpAllocationPolicy    `tfsdk:"ip_allocation_policy" json:"ipAllocationPolicy,omitempty"`
	Network            *MKSNetwork               `tfsdk:"network" json:"network,omitempty"`
}

type CloudManagedKubernetesAPIResponse struct {
	Id             string                                 `json:"id"`
	Checksum       string                                 `json:"checksum"`
	CreatedAt      string                                 `json:"createdAt"`
	UpdatedAt      string                                 `json:"updatedAt"`
	ResourceStatus string                                 `json:"resourceStatus"`
	TargetSpec     *CloudManagedKubernetesAPITargetSpec   `json:"targetSpec,omitempty"`
	CurrentState   *CloudManagedKubernetesAPICurrentState `json:"currentState,omitempty"`
	CurrentTasks   []CloudResourceTask                    `json:"currentTasks,omitempty"`
}

type CloudManagedKubernetesCreatePayload struct {
	TargetSpec CloudManagedKubernetesAPITargetSpec `json:"targetSpec"`
}

type CloudManagedKubernetesUpdatePayload struct {
	Checksum   string                                    `json:"checksum"`
	TargetSpec CloudManagedKubernetesAPIUpdateTargetSpec `json:"targetSpec"`
}

// ---------- Conversions ----------

func (m *CloudManagedKubernetesModel) ToCreate(ctx context.Context, diags *diag.Diagnostics) *CloudManagedKubernetesCreatePayload {
	return &CloudManagedKubernetesCreatePayload{TargetSpec: CloudManagedKubernetesAPITargetSpec{
		Name:               m.Name,
		Version:            m.Version,
		Region:             m.Region,
		Plan:               m.Plan,
		ApiServer:          objectAs[MKSApiServer](ctx, m.ApiServer, diags),
		Authentication:     objectAs[MKSAuthentication](ctx, m.Authentication, diags),
		AutoUpgrade:        objectAs[MKSAutoUpgrade](ctx, m.AutoUpgrade, diags),
		Cni:                objectAs[MKSCni](ctx, m.Cni, diags),
		IpAllocationPolicy: objectAs[MKSIpAllocationPolicy](ctx, m.IpAllocationPolicy, diags),
		Network:            objectAs[MKSNetwork](ctx, m.Network, diags),
	}}
}

func (m *CloudManagedKubernetesModel) ToUpdate(ctx context.Context, checksum string, diags *diag.Diagnostics) *CloudManagedKubernetesUpdatePayload {
	ts := CloudManagedKubernetesAPIUpdateTargetSpec{
		Name:           m.Name,
		Version:        m.Version,
		ApiServer:      objectAs[MKSApiServer](ctx, m.ApiServer, diags),
		Authentication: objectAs[MKSAuthentication](ctx, m.Authentication, diags),
		AutoUpgrade:    objectAs[MKSAutoUpgrade](ctx, m.AutoUpgrade, diags),
	}

	if cni := objectAs[MKSCni](ctx, m.Cni, diags); cni != nil {
		ts.Cni = &MKSCniUpdate{Cilium: cni.Cilium}
		if cni.Canal != nil && cni.Canal.KubeProxy != nil {
			kp := cni.Canal.KubeProxy
			ts.Cni.Canal = &MKSCniCanalUpdate{KubeProxy: &MKSKubeProxyUpdate{Iptables: kp.Iptables, Ipvs: kp.Ipvs}}
		}
	}
	if network := objectAs[MKSNetwork](ctx, m.Network, diags); network != nil && network.Vrack != nil {
		ts.Network = &MKSNetworkUpdate{Vrack: &MKSNetworkVrackUpdate{DefaultEgress: network.Vrack.DefaultEgress}}
	}

	return &CloudManagedKubernetesUpdatePayload{Checksum: checksum, TargetSpec: ts}
}

// MergeWith copies the API response into the model.
func (m *CloudManagedKubernetesModel) MergeWith(ctx context.Context, response *CloudManagedKubernetesAPIResponse, diags *diag.Diagnostics) {
	attrTypes := managedKubernetesAttrTypes()

	m.Id = tfStr(response.Id)
	m.Checksum = tfStr(response.Checksum)
	m.CreatedAt = tfStr(response.CreatedAt)
	m.UpdatedAt = tfStr(response.UpdatedAt)
	m.ResourceStatus = tfStr(response.ResourceStatus)
	m.CurrentState = objectFrom(ctx, attrTypes["current_state"], response.CurrentState, diags)

	ts := response.TargetSpec
	if ts == nil {
		ts = &CloudManagedKubernetesAPITargetSpec{}
	}
	// The API returns every admission plugin: keep the declared ones, current_state has them all.
	if prior := objectAs[MKSApiServer](ctx, m.ApiServer, diags); prior != nil && prior.AdmissionPlugins != nil && ts.ApiServer != nil {
		ts.ApiServer.AdmissionPlugins = declaredAdmissionPlugins(prior.AdmissionPlugins, ts.ApiServer.AdmissionPlugins)
	}
	m.Name = ts.Name
	m.Region = ts.Region
	m.Plan = ts.Plan
	// The API stores the minor only: keep a declared full version of the same minor.
	if m.Version.IsNull() || m.Version.IsUnknown() || minorVersion(m.Version.ValueString()) != ts.Version.ValueString() {
		m.Version = ts.Version
	}
	m.ApiServer = objectFrom(ctx, attrTypes["api_server"], ts.ApiServer, diags)
	m.Authentication = objectFrom(ctx, attrTypes["authentication"], ts.Authentication, diags)
	m.AutoUpgrade = objectFrom(ctx, attrTypes["auto_upgrade"], ts.AutoUpgrade, diags)
	m.Cni = objectFrom(ctx, attrTypes["cni"], ts.Cni, diags)
	m.IpAllocationPolicy = objectFrom(ctx, attrTypes["ip_allocation_policy"], ts.IpAllocationPolicy, diags)
	m.Network = objectFrom(ctx, attrTypes["network"], ts.Network, diags)
}

// minorVersion returns the major.minor part of a Kubernetes version, e.g. 1.34 for 1.34.1-2.
func minorVersion(v string) string {
	if parts := strings.SplitN(v, ".", 3); len(parts) == 3 {
		return parts[0] + "." + parts[1]
	}
	return v
}

func declaredAdmissionPlugins(prior, api map[string]MKSAdmissionPlugin) map[string]MKSAdmissionPlugin {
	declared := make(map[string]MKSAdmissionPlugin, len(prior))
	for k := range prior {
		if v, ok := api[k]; ok {
			declared[k] = v
		}
	}
	return declared
}
