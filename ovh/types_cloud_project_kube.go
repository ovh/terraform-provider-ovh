package ovh

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type CloudProjectKubeUpdatePolicyOpts struct {
	UpdatePolicy string `json:"updatePolicy"`
}

// CloudProjectKubePutOpts update cluster options
type CloudProjectKubePutOpts struct {
	Name *string `json:"name,omitempty"`
}

type privateNetworkConfiguration struct {
	DefaultVrackGateway            string `json:"defaultVrackGateway"`
	PrivateNetworkRoutingAsDefault bool   `json:"privateNetworkRoutingAsDefault"`
}

type CloudProjectKubeUpdateLoadBalancersSubnetIdOpts struct {
	LoadBalancersSubnetId string `json:"loadBalancersSubnetId"`
}

type CloudProjectKubeCreateOpts struct {
	Name                        *string                             `json:"name,omitempty"`
	PrivateNetworkId            *string                             `json:"privateNetworkId,omitempty"`
	PrivateNetworkConfiguration *privateNetworkConfiguration        `json:"privateNetworkConfiguration,omitempty"`
	Region                      string                              `json:"region"`
	Version                     *string                             `json:"version,omitempty"`
	Plan                        *string                             `json:"plan,omitempty"`
	UpdatePolicy                *string                             `json:"updatePolicy,omitempty"`
	Customization               *Customization                      `json:"customization,omitempty"`
	KubeProxyMode               *string                             `json:"kubeProxyMode,omitempty"`
	LoadBalancersSubnetId       *string                             `json:"loadBalancersSubnetId,omitempty"`
	NodesSubnetId               *string                             `json:"nodesSubnetId,omitempty"`
	IPAllocationPolicy          *CloudProjectKubeIPAllocationPolicy `json:"ipAllocationPolicy,omitempty"`
}

type Customization struct {
	APIServer *APIServer                           `json:"apiServer,omitempty"`
	KubeProxy *kubeProxyCustomization              `json:"kubeProxy,omitempty"`
	Cilium    *CloudProjectKubeCiliumCustomization `json:"cilium,omitempty"`
}

type APIServer struct {
	AdmissionPlugins *AdmissionPlugins `json:"admissionPlugins,omitempty"`
}

type kubeProxyCustomization struct {
	IPTables *kubeProxyCustomizationIPTables `json:"iptables,omitempty"`
	IPVS     *kubeProxyCustomizationIPVS     `json:"ipvs,omitempty"`
}

type kubeProxyCustomizationIPTables struct {
	MinSyncPeriod *string `json:"minSyncPeriod,omitempty"`
	SyncPeriod    *string `json:"syncPeriod,omitempty"`
}

type kubeProxyCustomizationIPVS struct {
	MinSyncPeriod *string `json:"minSyncPeriod,omitempty"`
	Scheduler     *string `json:"scheduler,omitempty"`
	SyncPeriod    *string `json:"syncPeriod,omitempty"`
	TCPFinTimeout *string `json:"tcpFinTimeout,omitempty"`
	TCPTimeout    *string `json:"tcpTimeout,omitempty"`
	UDPTimeout    *string `json:"udpTimeout,omitempty"`
}

type CloudProjectKubeCiliumCustomization struct {
	ClusterId   *int                                            `json:"clusterId,omitempty"`
	ClusterMesh *CloudProjectKubeCiliumCustomizationClusterMesh `json:"clusterMesh,omitempty"`
	Hubble      *CloudProjectKubeCiliumCustomizationHubble      `json:"hubble,omitempty"`
}

type CloudProjectKubeCiliumCustomizationClusterMesh struct {
	Enabled   *bool                                                    `json:"enabled,omitempty"`
	ApiServer *CloudProjectKubeCiliumCustomizationClusterMeshApiServer `json:"apiServer,omitempty"`
}

type CloudProjectKubeCiliumCustomizationClusterMeshApiServer struct {
	NodePort    *int    `json:"nodePort,omitempty"`
	ServiceType *string `json:"serviceType,omitempty"`
}
type CloudProjectKubeCiliumCustomizationHubble struct {
	Enabled *bool                                           `json:"enabled,omitempty"`
	Relay   *CloudProjectKubeCiliumCustomizationHubbleRelay `json:"relay,omitempty"`
	UI      *CloudProjectKubeCiliumCustomizationHubbleUI    `json:"ui,omitempty"`
}

type CloudProjectKubeCiliumCustomizationHubbleRelay struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type CloudProjectKubeCiliumCustomizationHubbleUI struct {
	Enabled           *bool                                         `json:"enabled,omitempty"`
	BackendResources  *CloudProjectKubeCiliumCustomizationResources `json:"backendResources,omitempty"`
	FrontendResources *CloudProjectKubeCiliumCustomizationResources `json:"frontendResources,omitempty"`
}

type CloudProjectKubeCiliumCustomizationResources struct {
	Limits   *CloudProjectKubeCiliumCustomizationResourcesType `json:"limits,omitempty"`
	Requests *CloudProjectKubeCiliumCustomizationResourcesType `json:"requests,omitempty"`
}

type CloudProjectKubeCiliumCustomizationResourcesType struct {
	CPU    *string `json:"cpu,omitempty"`
	Memory *string `json:"memory,omitempty"`
}
type AdmissionPlugins struct {
	Enabled  *[]string `json:"enabled,omitempty"`
	Disabled *[]string `json:"disabled,omitempty"`
}

type CloudProjectKubeIPAllocationPolicy struct {
	PodsIpv4Cidr     *string `json:"podsIpv4Cidr,omitempty"`
	ServicesIpv4Cidr *string `json:"servicesIpv4Cidr,omitempty"`
}

func (opts *CloudProjectKubeCreateOpts) FromResource(d *schema.ResourceData) {
	opts.Region = d.Get(kubeRegionKey).(string)
	opts.Version = helpers.GetNilStringPointerFromData(d, kubeVersionKey)
	opts.Plan = helpers.GetNilStringPointerFromData(d, kubeClusterPlanKey)
	opts.Name = helpers.GetNilStringPointerFromData(d, kubeNameKey)
	opts.UpdatePolicy = helpers.GetNilStringPointerFromData(d, kubeClusterUpdatePolicyKey)
	opts.LoadBalancersSubnetId = helpers.GetNilStringPointerFromData(d, kubeClusterLoadBalancersSubnetIdKey)
	opts.NodesSubnetId = helpers.GetNilStringPointerFromData(d, kubeClusterNodesSubnetIdKey)
	opts.PrivateNetworkId = helpers.GetNilStringPointerFromData(d, kubeClusterPrivateNetworkIDKey)
	opts.PrivateNetworkConfiguration = loadPrivateNetworkConfiguration(d.Get(kubeClusterPrivateNetworkConfigurationKey))
	opts.KubeProxyMode = helpers.GetNilStringPointerFromData(d, kubeClusterProxyModeKey)

	opts.Customization = &Customization{
		APIServer: nil,
		KubeProxy: loadKubeProxyCustomization(d.Get(kubeClusterCustomizationKubeProxyKey)),
		Cilium:    loadCiliumCustomizationFromResource(d.Get(kubeClusterCustomizationCiliumKey)),
	}

	opts.IPAllocationPolicy = loadIPAllocationPolicyFromResource(d.Get(kubeClusterIpAllocationPolicyKey))

	// load the filled api server customization
	// both the new and the deprecated syntax are supported, but they are mutual exclusive
	if userIsUsingDeprecatedCustomizationSyntax(d) {
		log.Printf("[DEBUG] Using DEPRECATED syntax for api server customization")
		opts.Customization.APIServer = loadDeprecatedApiServerCustomization(d.Get(kubeClusterCustomization))
	} else {
		log.Printf("[DEBUG] Using new syntax for api server customization")
		opts.Customization.APIServer = loadApiServerCustomization(d.Get(kubeClusterCustomizationApiServerKey))
	}
}

func userIsUsingDeprecatedCustomizationSyntax(d *schema.ResourceData) bool {
	funcTypeSetNotNilAndNotEmpty := func(d *schema.ResourceData, key string) bool {
		return d.Get(key) != nil && len(d.Get(key).(*schema.Set).List()) > 0
	}

	return funcTypeSetNotNilAndNotEmpty(d, kubeClusterCustomization)
}

// loadApiServerCustomization reads the api server customization
func loadApiServerCustomization(apiServerAdmissionPlugins interface{}) *APIServer {
	if apiServerAdmissionPlugins == nil {
		return nil
	}

	apiServerOutput := &APIServer{
		AdmissionPlugins: &AdmissionPlugins{},
	}

	// Customization
	customizationSet := apiServerAdmissionPlugins.(*schema.Set).List()
	if len(customizationSet) > 0 {
		customization := customizationSet[0].(map[string]interface{})
		admissionPluginsSet := customization[kubeClusterCustomizationAdmissionPluginsKey].(*schema.Set).List()
		admissionPlugins := admissionPluginsSet[0].(map[string]interface{})

		readApiServerAdmissionPlugins(admissionPlugins, apiServerOutput)

		log.Printf("[DEBUG] Enabled admission plugins from new syntax: %v", apiServerOutput.AdmissionPlugins.Enabled)
		log.Printf("[DEBUG] Disabled admission plugins from new syntax: %v", apiServerOutput.AdmissionPlugins.Disabled)
	}

	return apiServerOutput
}

func loadIPAllocationPolicyFromResource(i interface{}) *CloudProjectKubeIPAllocationPolicy {
	if i == nil {
		return nil
	}

	ipAllocationPolicySet := i.(*schema.Set).List()
	if len(ipAllocationPolicySet) == 0 {
		return nil
	}

	// Due to this bug https://github.com/hashicorp/terraform-plugin-sdk/pull/1042
	// when updating the 'ip_allocation_policy' object there are two objects, one is empty, so take the non-empty one
	ipAllocationPolicyObject := ipAllocationPolicySet[0].(map[string]interface{}) // by default take the first one
	for _, to := range ipAllocationPolicySet {
		empty := true

		object := to.(map[string]interface{})
		if object[kubeClusterPodsIpv4CidrKey].(string) != "" || object[kubeClusterServicesIpv4CidrKey].(string) != "" {
			empty = false
		}

		if empty {
			continue
		}

		// We found the not empty object
		ipAllocationPolicyObject = object
	}

	return &CloudProjectKubeIPAllocationPolicy{
		PodsIpv4Cidr:     helpers.GetNilStringPointerFromData(ipAllocationPolicyObject, kubeClusterPodsIpv4CidrKey),
		ServicesIpv4Cidr: helpers.GetNilStringPointerFromData(ipAllocationPolicyObject, kubeClusterServicesIpv4CidrKey),
	}
}

func readApiServerAdmissionPlugins(admissionPlugins map[string]interface{}, apiServerOutput *APIServer) {
	// Enabled admission plugins
	{
		stringArray := admissionPlugins[kubeClusterCustomizationEnabledKey].([]interface{})
		enabled := make([]string, 0, len(stringArray))
		for _, s := range stringArray {
			enabled = append(enabled, s.(string))
		}
		apiServerOutput.AdmissionPlugins.Enabled = &enabled
	}

	// Disabled admission plugins
	{
		stringArray := admissionPlugins[kubeClusterCustomizationDisabledKey].([]interface{})
		disabled := make([]string, 0, len(stringArray))
		for _, s := range stringArray {
			disabled = append(disabled, s.(string))
		}
		apiServerOutput.AdmissionPlugins.Disabled = &disabled
	}
}

// loadDeprecatedApiServerCustomization reads the deprecated api server customization
// Deprecated, should be removed in the future
func loadDeprecatedApiServerCustomization(deprecatedApiServerCustomizationInterface interface{}) *APIServer {
	if deprecatedApiServerCustomizationInterface == nil {
		return nil
	}

	apiServerOutput := &APIServer{
		AdmissionPlugins: &AdmissionPlugins{},
	}

	oldCustomizationSet := deprecatedApiServerCustomizationInterface.(*schema.Set).List()
	if len(oldCustomizationSet) > 0 {
		oldApiServerCustomization := oldCustomizationSet[0].(map[string]interface{})
		oldApiServerCustomizationSet := oldApiServerCustomization[kubeClusterCustomizationApiServerNestedKey].(*schema.Set).List()

		if len(oldApiServerCustomizationSet) > 0 {
			oldApiServerCustomizationAdmissionPlugins := oldApiServerCustomizationSet[0].(map[string]interface{})
			oldApiServerCustomizationAdmissionPluginsSet := oldApiServerCustomizationAdmissionPlugins[kubeClusterCustomizationAdmissionPluginsKey].(*schema.Set).List()
			admissionPlugins := oldApiServerCustomizationAdmissionPluginsSet[0].(map[string]interface{})

			readApiServerAdmissionPlugins(admissionPlugins, apiServerOutput)
		}
	}

	log.Printf("[DEBUG] Enabled admission plugins from DEPRECATED syntax: %v", apiServerOutput.AdmissionPlugins.Enabled)
	log.Printf("[DEBUG] Disabled admission plugins from DEPRECATED syntax: %v", apiServerOutput.AdmissionPlugins.Disabled)

	return apiServerOutput
}

// loadKubeProxyCustomization reads the kube proxy customization
func loadKubeProxyCustomization(kubeProxyCustomizationInterface interface{}) *kubeProxyCustomization {
	if kubeProxyCustomizationInterface == nil {
		return nil
	}

	kubeProxyOutput := &kubeProxyCustomization{
		IPTables: &kubeProxyCustomizationIPTables{},
		IPVS:     &kubeProxyCustomizationIPVS{},
	}

	kubeProxySet := kubeProxyCustomizationInterface.(*schema.Set).List()
	if len(kubeProxySet) > 0 {
		kubeProxy := kubeProxySet[0].(map[string]interface{})

		// Nested IPTables customization
		{
			ipTablesSet := kubeProxy[kubeClusterCustomizationIptablesKey].(*schema.Set).List()
			if len(ipTablesSet) > 0 {
				ipTables := ipTablesSet[0].(map[string]interface{})
				kubeProxyOutput.IPTables.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipTables, kubeClusterCustomizationMinSyncPeriodKey)
				kubeProxyOutput.IPTables.SyncPeriod = helpers.GetNilStringPointerFromData(ipTables, kubeClusterCustomizationSyncPeriodKey)
			}
		}

		// Nested IPVS customization
		{
			ipvsSet := kubeProxy[kubeClusterCustomizationIpvsKey].(*schema.Set).List()
			if len(ipvsSet) > 0 {
				ipvs := ipvsSet[0].(map[string]interface{})
				kubeProxyOutput.IPVS.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationMinSyncPeriodKey)
				kubeProxyOutput.IPVS.Scheduler = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationSchedulerKey)
				kubeProxyOutput.IPVS.SyncPeriod = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationSyncPeriodKey)
				kubeProxyOutput.IPVS.TCPFinTimeout = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationTcpFinTimeoutKey)
				kubeProxyOutput.IPVS.TCPTimeout = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationTcpTimeoutKey)
				kubeProxyOutput.IPVS.UDPTimeout = helpers.GetNilStringPointerFromData(ipvs, kubeClusterCustomizationUdpTimeoutKey)
			}
		}
	}

	return kubeProxyOutput
}

func loadPrivateNetworkConfiguration(i interface{}) *privateNetworkConfiguration {
	if i == nil {
		return nil
	}
	pncOutput := privateNetworkConfiguration{}

	pncSet := i.(*schema.Set).List()
	for _, pnc := range pncSet {
		mapping := pnc.(map[string]interface{})
		pncOutput.DefaultVrackGateway = mapping[kubeClusterDefaultVrackGatewayKey].(string)
		pncOutput.PrivateNetworkRoutingAsDefault = mapping[kubeClusterPrivateNetworkRoutingAsDefault].(bool)
	}
	return &pncOutput
}

func (opts *CloudProjectKubeCreateOpts) String() string {
	var str string
	if opts.Name != nil {
		str = *opts.Name
	}

	str += fmt.Sprintf(" (%s)", opts.Region)

	if opts.Version != nil {
		str += fmt.Sprintf(": %s", *opts.Version)
	}

	return str
}

type CloudProjectKubeResponse struct {
	ControlPlaneIsUpToDate bool                               `json:"controlPlaneIsUpToDate"`
	Id                     string                             `json:"id"`
	IsUpToDate             bool                               `json:"isUpToDate"`
	LoadBalancersSubnetId  string                             `json:"loadBalancersSubnetId"`
	Name                   string                             `json:"name"`
	NextUpgradeVersions    []string                           `json:"nextUpgradeVersions"`
	NodesSubnetId          string                             `json:"nodesSubnetId"`
	NodesUrl               string                             `json:"nodesUrl"`
	PrivateNetworkId       string                             `json:"privateNetworkId"`
	Region                 string                             `json:"region"`
	Status                 string                             `json:"status"`
	UpdatePolicy           string                             `json:"updatePolicy"`
	Url                    string                             `json:"url"`
	Version                string                             `json:"version"`
	Plan                   string                             `json:"plan"`
	Customization          Customization                      `json:"customization"`
	KubeProxyMode          string                             `json:"kubeProxyMode"`
	IPAllocationPolicy     CloudProjectKubeIPAllocationPolicy `json:"ipAllocationPolicy"`
}

func (v *CloudProjectKubeResponse) ToMap(d *schema.ResourceData) map[string]interface{} {
	obj := make(map[string]interface{})
	obj[kubeClusterControlPlaneIsUpToDateKey] = v.ControlPlaneIsUpToDate
	obj[kubeNodeIdKey] = v.Id
	obj[kubeClusterIsUpToDateKey] = v.IsUpToDate
	obj[kubeClusterLoadBalancersSubnetIdKey] = v.LoadBalancersSubnetId
	obj[kubeNameKey] = v.Name
	obj[kubeClusterNextUpgradeVersionsKey] = v.NextUpgradeVersions
	obj[kubeClusterNodesSubnetIdKey] = v.NodesSubnetId
	obj[kubeClusterNodesUrlKey] = v.NodesUrl
	obj[kubeClusterPrivateNetworkIDKey] = v.PrivateNetworkId
	obj[kubeRegionKey] = v.Region
	obj[kubeStatusKey] = v.Status
	obj[kubeClusterUpdatePolicyKey] = v.UpdatePolicy
	obj[kubeClusterUrlKey] = v.Url
	obj[kubeClusterPlanKey] = v.Plan
	loadKubeIPAllocationPolicy(obj, v)
	versionPatch, err := version.NewVersion(v.Version)
	if err != nil {
		// if fail, return to the previous implementation
		obj[kubeVersionKey] = v.Version[:strings.LastIndex(v.Version, ".")]
	} else {
		// versionPatch.String() return a true semantic version (0.0.0)
		obj[kubeVersionKey] = v.Version[:strings.LastIndex(versionPatch.String(), ".")]
	}
	obj[kubeClusterProxyModeKey] = v.KubeProxyMode

	if v.Customization.APIServer != nil {
		if userIsUsingDeprecatedCustomizationSyntax(d) {
			loadDeprecatedApiServerCustomizationToMap(obj, v)
		} else {
			loadApiServerCustomizationToMap(obj, v)
		}
	}

	if v.Customization.KubeProxy != nil {
		loadKubeProxyCustomizationToMap(obj, v)
	}

	if v.Customization.Cilium != nil {
		loadKubeCustomizationCilium(obj, v)
	}

	return obj
}

func loadKubeProxyCustomizationToMap(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	obj[kubeClusterCustomizationKubeProxyKey] = []map[string]interface{}{{}}

	if v.Customization.KubeProxy.IPTables != nil {
		data := make(map[string]interface{})
		if vv := v.Customization.KubeProxy.IPTables.MinSyncPeriod; vv != nil && *vv != "" {
			data[kubeClusterCustomizationMinSyncPeriodKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPTables.SyncPeriod; vv != nil && *vv != "" {
			data[kubeClusterCustomizationSyncPeriodKey] = vv
		}

		if len(data) > 0 {
			obj[kubeClusterCustomizationKubeProxyKey].([]map[string]interface{})[0][kubeClusterCustomizationIptablesKey] = []map[string]interface{}{data}
		}
	}

	if v.Customization.KubeProxy.IPVS != nil {
		data := make(map[string]interface{})
		if vv := v.Customization.KubeProxy.IPVS.MinSyncPeriod; vv != nil && *vv != "" {
			data[kubeClusterCustomizationMinSyncPeriodKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.Scheduler; vv != nil && *vv != "" {
			data[kubeClusterCustomizationSchedulerKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.SyncPeriod; vv != nil && *vv != "" {
			data[kubeClusterCustomizationSyncPeriodKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.TCPFinTimeout; vv != nil && *vv != "" {
			data[kubeClusterCustomizationTcpFinTimeoutKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.TCPTimeout; vv != nil && *vv != "" {
			data[kubeClusterCustomizationTcpTimeoutKey] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.UDPTimeout; vv != nil && *vv != "" {
			data[kubeClusterCustomizationUdpTimeoutKey] = vv
		}

		if len(data) > 0 {
			obj[kubeClusterCustomizationKubeProxyKey].([]map[string]interface{})[0][kubeClusterCustomizationIpvsKey] = []map[string]interface{}{data}
		}
	}

	// Delete entire customization_kube_proxy if empty
	if len(obj[kubeClusterCustomizationKubeProxyKey].([]map[string]interface{})[0]) == 0 {
		delete(obj, kubeClusterCustomizationKubeProxyKey)
	}
}

// Deprecated: use loadApiServerCustomizationToMap instead
func loadDeprecatedApiServerCustomizationToMap(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	obj[kubeClusterCustomization] = []map[string]interface{}{
		{
			kubeClusterCustomizationApiServerNestedKey: []map[string]interface{}{
				{
					kubeClusterCustomizationAdmissionPluginsKey: []map[string]interface{}{
						{
							kubeClusterCustomizationEnabledKey:  v.Customization.APIServer.AdmissionPlugins.Enabled,
							kubeClusterCustomizationDisabledKey: v.Customization.APIServer.AdmissionPlugins.Disabled,
						},
					},
				},
			},
		},
	}
}

func loadApiServerCustomizationToMap(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	obj[kubeClusterCustomizationApiServerKey] = []map[string]interface{}{
		{
			kubeClusterCustomizationAdmissionPluginsKey: []map[string]interface{}{
				{
					kubeClusterCustomizationEnabledKey:  v.Customization.APIServer.AdmissionPlugins.Enabled,
					kubeClusterCustomizationDisabledKey: v.Customization.APIServer.AdmissionPlugins.Disabled,
				},
			},
		},
	}
}

func loadKubeIPAllocationPolicy(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	obj[kubeClusterIpAllocationPolicyKey] = []map[string]interface{}{
		{
			kubeClusterPodsIpv4CidrKey:     v.IPAllocationPolicy.PodsIpv4Cidr,
			kubeClusterServicesIpv4CidrKey: v.IPAllocationPolicy.ServicesIpv4Cidr,
		},
	}
}

func loadKubeCustomizationCilium(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	if v.Customization.Cilium == nil {
		return
	}

	ciliumMap := map[string]interface{}{}

	if v.Customization.Cilium.ClusterId != nil {
		ciliumMap[kubeClusterCiliumClusterID] = *v.Customization.Cilium.ClusterId
	}

	if v.Customization.Cilium.ClusterMesh != nil {
		clusterMeshMap := map[string]interface{}{}

		if v.Customization.Cilium.ClusterMesh.Enabled != nil {
			clusterMeshMap[kubeClusterCustomizationEnabledKey] = *v.Customization.Cilium.ClusterMesh.Enabled
		}

		if v.Customization.Cilium.ClusterMesh.ApiServer != nil {
			apiServerMap := map[string]interface{}{}
			if v.Customization.Cilium.ClusterMesh.ApiServer.NodePort != nil {
				apiServerMap[kubeClusterCiliumNodePortKey] = *v.Customization.Cilium.ClusterMesh.ApiServer.NodePort
			}
			if v.Customization.Cilium.ClusterMesh.ApiServer.ServiceType != nil {
				apiServerMap[kubeClusterCiliumServiceTypeKey] = *v.Customization.Cilium.ClusterMesh.ApiServer.ServiceType
			}
			if len(apiServerMap) > 0 {
				clusterMeshMap[kubeClusterCiliumApiServerKey] = []map[string]interface{}{apiServerMap}
			}
		}

		ciliumMap[kubeClusterCiliumClusterMeshKey] = []map[string]interface{}{clusterMeshMap}
	}

	if v.Customization.Cilium.Hubble != nil {
		hubbleMap := map[string]interface{}{}

		if v.Customization.Cilium.Hubble.Enabled != nil {
			hubbleMap[kubeClusterCustomizationEnabledKey] = *v.Customization.Cilium.Hubble.Enabled
		}

		if v.Customization.Cilium.Hubble.Relay != nil {
			relayMap := map[string]interface{}{}
			if v.Customization.Cilium.Hubble.Relay.Enabled != nil {
				relayMap[kubeClusterCustomizationEnabledKey] = *v.Customization.Cilium.Hubble.Relay.Enabled
			}
			hubbleMap[kubeClusterCiliumRelayKey] = []map[string]interface{}{relayMap}
		}

		if v.Customization.Cilium.Hubble.UI != nil {
			uiMap := map[string]interface{}{}

			if v.Customization.Cilium.Hubble.UI.Enabled != nil {
				uiMap[kubeClusterCustomizationEnabledKey] = *v.Customization.Cilium.Hubble.UI.Enabled
			}

			if v.Customization.Cilium.Hubble.UI.BackendResources != nil {
				uiMap[kubeClusterCiliumBackendResources] = []map[string]interface{}{loadCiliumResourcesToMap(v.Customization.Cilium.Hubble.UI.BackendResources)}
			}

			if v.Customization.Cilium.Hubble.UI.FrontendResources != nil {
				uiMap[kubeClusterCiliumFrontendResources] = []map[string]interface{}{loadCiliumResourcesToMap(v.Customization.Cilium.Hubble.UI.FrontendResources)}
			}

			hubbleMap[kubeClusterCiliumUiKey] = []map[string]interface{}{uiMap}
		}

		ciliumMap[kubeClusterCiliumHubbleKey] = []map[string]interface{}{hubbleMap}
	}

	obj[kubeClusterCustomizationCiliumKey] = []map[string]interface{}{ciliumMap}
}

func loadCiliumResourcesToMap(r *CloudProjectKubeCiliumCustomizationResources) map[string]interface{} {
	resourcesMap := map[string]interface{}{}

	if r.Limits != nil {
		limitsMap := map[string]interface{}{}
		if r.Limits.CPU != nil {
			limitsMap[kubeClusterCiliumCpuKey] = *r.Limits.CPU
		}
		if r.Limits.Memory != nil {
			limitsMap[kubeClusterCiliumMemoryKey] = *r.Limits.Memory
		}
		resourcesMap[kubeClusterCiliumLimitsKey] = []map[string]interface{}{limitsMap}
	}

	if r.Requests != nil {
		requestsMap := map[string]interface{}{}
		if r.Requests.CPU != nil {
			requestsMap[kubeClusterCiliumCpuKey] = *r.Requests.CPU
		}
		if r.Requests.Memory != nil {
			requestsMap[kubeClusterCiliumMemoryKey] = *r.Requests.Memory
		}
		resourcesMap[kubeClusterCiliumRequestsKey] = []map[string]interface{}{requestsMap}
	}

	return resourcesMap
}

// loadCiliumCustomizationFromResource parses the customization_cilium block from Terraform
// resource data into a CloudProjectKubeCiliumCustomization struct.
func loadCiliumCustomizationFromResource(i interface{}) *CloudProjectKubeCiliumCustomization {
	if i == nil {
		return nil
	}

	ciliumSet := i.(*schema.Set).List()
	if len(ciliumSet) == 0 {
		return nil
	}

	cilium := ciliumSet[0].(map[string]interface{})
	output := &CloudProjectKubeCiliumCustomization{}

	if v, ok := cilium[kubeClusterCiliumClusterID].(int); ok && v != 0 {
		output.ClusterId = &v
	}

	// cluster_mesh
	if clusterMeshRaw, ok := cilium[kubeClusterCiliumClusterMeshKey]; ok {
		clusterMeshSet := clusterMeshRaw.(*schema.Set).List()
		if len(clusterMeshSet) > 0 {
			clusterMeshMap := clusterMeshSet[0].(map[string]interface{})
			output.ClusterMesh = &CloudProjectKubeCiliumCustomizationClusterMesh{}

			if v, ok := clusterMeshMap[kubeClusterCustomizationEnabledKey].(bool); ok {
				output.ClusterMesh.Enabled = &v
			}

			if apiServerRaw, ok := clusterMeshMap[kubeClusterCiliumApiServerKey]; ok {
				apiServerSet := apiServerRaw.(*schema.Set).List()
				if len(apiServerSet) > 0 {
					apiServerMap := apiServerSet[0].(map[string]interface{})
					output.ClusterMesh.ApiServer = &CloudProjectKubeCiliumCustomizationClusterMeshApiServer{}

					if v, ok := apiServerMap[kubeClusterCiliumNodePortKey].(int); ok && v != 0 {
						output.ClusterMesh.ApiServer.NodePort = &v
					}
					if v, ok := apiServerMap[kubeClusterCiliumServiceTypeKey].(string); ok && v != "" {
						output.ClusterMesh.ApiServer.ServiceType = &v
					}
				}
			}
		}
	}

	// hubble
	if hubbleRaw, ok := cilium[kubeClusterCiliumHubbleKey]; ok {
		hubbleSet := hubbleRaw.(*schema.Set).List()
		if len(hubbleSet) > 0 {
			hubbleMap := hubbleSet[0].(map[string]interface{})
			output.Hubble = &CloudProjectKubeCiliumCustomizationHubble{}

			if v, ok := hubbleMap[kubeClusterCustomizationEnabledKey].(bool); ok {
				output.Hubble.Enabled = &v
			}

			// relay
			if relayRaw, ok := hubbleMap[kubeClusterCiliumRelayKey]; ok {
				relaySet := relayRaw.(*schema.Set).List()
				if len(relaySet) > 0 {
					relayMap := relaySet[0].(map[string]interface{})
					output.Hubble.Relay = &CloudProjectKubeCiliumCustomizationHubbleRelay{}
					if v, ok := relayMap[kubeClusterCustomizationEnabledKey].(bool); ok {
						output.Hubble.Relay.Enabled = &v
					}
				}
			}

			// ui
			if uiRaw, ok := hubbleMap[kubeClusterCiliumUiKey]; ok {
				uiSet := uiRaw.(*schema.Set).List()
				if len(uiSet) > 0 {
					uiMap := uiSet[0].(map[string]interface{})
					output.Hubble.UI = &CloudProjectKubeCiliumCustomizationHubbleUI{}

					if v, ok := uiMap[kubeClusterCustomizationEnabledKey].(bool); ok {
						output.Hubble.UI.Enabled = &v
					}

					if backendRaw, ok := uiMap[kubeClusterCiliumBackendResources]; ok {
						output.Hubble.UI.BackendResources = loadCiliumResourcesFromMap(backendRaw.(*schema.Set).List())
					}

					if frontendRaw, ok := uiMap[kubeClusterCiliumFrontendResources]; ok {
						output.Hubble.UI.FrontendResources = loadCiliumResourcesFromMap(frontendRaw.(*schema.Set).List())
					}
				}
			}
		}
	}

	return output
}

func loadCiliumResourcesFromMap(setList []interface{}) *CloudProjectKubeCiliumCustomizationResources {
	if len(setList) == 0 {
		return nil
	}

	resourcesMap := setList[0].(map[string]interface{})
	output := &CloudProjectKubeCiliumCustomizationResources{}

	if limitsRaw, ok := resourcesMap[kubeClusterCiliumLimitsKey]; ok {
		limitsList := limitsRaw.(*schema.Set).List()
		if len(limitsList) > 0 {
			limitsMap := limitsList[0].(map[string]interface{})
			output.Limits = &CloudProjectKubeCiliumCustomizationResourcesType{}
			if v, ok := limitsMap[kubeClusterCiliumCpuKey].(string); ok && v != "" {
				output.Limits.CPU = &v
			}
			if v, ok := limitsMap[kubeClusterCiliumMemoryKey].(string); ok && v != "" {
				output.Limits.Memory = &v
			}
		}
	}

	if requestsRaw, ok := resourcesMap[kubeClusterCiliumRequestsKey]; ok {
		requestsList := requestsRaw.(*schema.Set).List()
		if len(requestsList) > 0 {
			requestsMap := requestsList[0].(map[string]interface{})
			output.Requests = &CloudProjectKubeCiliumCustomizationResourcesType{}
			if v, ok := requestsMap[kubeClusterCiliumCpuKey].(string); ok && v != "" {
				output.Requests.CPU = &v
			}
			if v, ok := requestsMap[kubeClusterCiliumMemoryKey].(string); ok && v != "" {
				output.Requests.Memory = &v
			}
		}
	}

	return output
}

func (v *CloudProjectKubeResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", v.Name, v.Id, v.Status)
}

type CloudProjectKubeKubeConfigResponse struct {
	Content string `json:"content"`
}

type CloudProjectKubeUpdateOpts struct {
	Strategy string `json:"strategy"`
}

type CloudProjectKubeResetOpts struct {
	PrivateNetworkId *string `json:"privateNetworkId,omitempty"`
}

type CloudProjectKubeUpdatePNCOpts struct {
	DefaultVrackGateway            string `json:"defaultVrackGateway"`
	PrivateNetworkRoutingAsDefault bool   `json:"privateNetworkRoutingAsDefault"`
}

type CloudProjectKubeUpdateCustomizationOpts struct {
	APIServer *APIServer                           `json:"apiServer,omitempty"`
	KubeProxy *kubeProxyCustomization              `json:"kubeProxy,omitempty"`
	Cilium    *CloudProjectKubeCiliumCustomization `json:"cilium,omitempty"`
}

type CloudProjectKubeNodeResponse struct {
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	DeployedAt string `json:"deployedAt"`
	Flavor     string `json:"flavor"`
	Id         string `json:"id"`
	InstanceId string `json:"instanceId"`
	IsUpToDate bool   `json:"isUpToDate"`
	Name       string `json:"name"`
	NodePoolId string `json:"nodePoolId"`
	ProjectId  string `json:"projectId"`
	Status     string `json:"status"`
	Version    string `json:"version"`
}

func (v CloudProjectKubeNodeResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj[kubeCreatedAtKey] = v.CreatedAt
	obj[kubeNodeDeployedAtKey] = v.DeployedAt
	obj[kubeFlavorKey] = v.Flavor
	obj[kubeNodeIdKey] = v.Id
	obj[kubeNodeInstanceIdKey] = v.InstanceId
	obj[kubeClusterIsUpToDateKey] = v.IsUpToDate
	obj[kubeNameKey] = v.Name
	obj[kubeNodePoolIdKey] = v.NodePoolId
	obj[kubeProjectIdKey] = v.ProjectId
	obj[kubeStatusKey] = v.Status
	obj[kubeUpdatedAtKey] = v.UpdatedAt
	obj[kubeVersionKey] = v.Version
	return obj
}
