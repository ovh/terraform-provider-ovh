package ovh

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
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
	Name                        *string                      `json:"name,omitempty"`
	PrivateNetworkId            *string                      `json:"privateNetworkId,omitempty"`
	PrivateNetworkConfiguration *privateNetworkConfiguration `json:"privateNetworkConfiguration,omitempty"`
	Region                      string                       `json:"region"`
	Version                     *string                      `json:"version,omitempty"`
	UpdatePolicy                *string                      `json:"updatePolicy,omitempty"`
	Customization               *Customization               `json:"customization,omitempty"`
	KubeProxyMode               *string                      `json:"kubeProxyMode,omitempty"`
	LoadBalancersSubnetId       *string                      `json:"loadBalancersSubnetId,omitempty"`
}

type Customization struct {
	APIServer *APIServer              `json:"apiServer,omitempty"`
	KubeProxy *kubeProxyCustomization `json:"kubeProxy,omitempty"`
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

type AdmissionPlugins struct {
	Enabled  *[]string `json:"enabled,omitempty"`
	Disabled *[]string `json:"disabled,omitempty"`
}

func (opts *CloudProjectKubeCreateOpts) FromResource(d *schema.ResourceData) {
	opts.Region = d.Get("region").(string)
	opts.Version = helpers.GetNilStringPointerFromData(d, "version")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	opts.UpdatePolicy = helpers.GetNilStringPointerFromData(d, "update_policy")
	opts.LoadBalancersSubnetId = helpers.GetNilStringPointerFromData(d, "load_balancers_subnet_id")
	opts.PrivateNetworkId = helpers.GetNilStringPointerFromData(d, "private_network_id")
	opts.PrivateNetworkConfiguration = loadPrivateNetworkConfiguration(d.Get("private_network_configuration"))
	opts.KubeProxyMode = helpers.GetNilStringPointerFromData(d, kubeClusterProxyModeKey)

	opts.Customization = &Customization{
		APIServer: nil,
		KubeProxy: loadKubeProxyCustomization(d.Get(kubeClusterCustomizationKubeProxyKey)),
	}

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
		admissionPluginsSet := customization["admissionplugins"].(*schema.Set).List()
		admissionPlugins := admissionPluginsSet[0].(map[string]interface{})

		readApiServerAdmissionPlugins(admissionPlugins, apiServerOutput)

		log.Printf("[DEBUG] Enabled admission plugins from new syntax: %v", apiServerOutput.AdmissionPlugins.Enabled)
		log.Printf("[DEBUG] Disabled admission plugins from new syntax: %v", apiServerOutput.AdmissionPlugins.Disabled)
	}

	return apiServerOutput
}

func readApiServerAdmissionPlugins(admissionPlugins map[string]interface{}, apiServerOutput *APIServer) {
	// Enabled admission plugins
	{
		stringArray := admissionPlugins["enabled"].([]interface{})
		enabled := make([]string, 0, len(stringArray))
		for _, s := range stringArray {
			enabled = append(enabled, s.(string))
		}
		apiServerOutput.AdmissionPlugins.Enabled = &enabled
	}

	// Disabled admission plugins
	{
		stringArray := admissionPlugins["disabled"].([]interface{})
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
		oldApiServerCustomizationSet := oldApiServerCustomization["apiserver"].(*schema.Set).List()

		if len(oldApiServerCustomizationSet) > 0 {
			oldApiServerCustomizationAdmissionPlugins := oldApiServerCustomizationSet[0].(map[string]interface{})
			oldApiServerCustomizationAdmissionPluginsSet := oldApiServerCustomizationAdmissionPlugins["admissionplugins"].(*schema.Set).List()
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
			ipTablesSet := kubeProxy["iptables"].(*schema.Set).List()
			if len(ipTablesSet) > 0 {
				ipTables := ipTablesSet[0].(map[string]interface{})
				kubeProxyOutput.IPTables.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipTables, "min_sync_period")
				kubeProxyOutput.IPTables.SyncPeriod = helpers.GetNilStringPointerFromData(ipTables, "sync_period")
			}
		}

		// Nested IPVS customization
		{
			ipvsSet := kubeProxy["ipvs"].(*schema.Set).List()
			if len(ipvsSet) > 0 {
				ipvs := ipvsSet[0].(map[string]interface{})
				kubeProxyOutput.IPVS.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipvs, "min_sync_period")
				kubeProxyOutput.IPVS.Scheduler = helpers.GetNilStringPointerFromData(ipvs, "scheduler")
				kubeProxyOutput.IPVS.SyncPeriod = helpers.GetNilStringPointerFromData(ipvs, "sync_period")
				kubeProxyOutput.IPVS.TCPFinTimeout = helpers.GetNilStringPointerFromData(ipvs, "tcp_fin_timeout")
				kubeProxyOutput.IPVS.TCPTimeout = helpers.GetNilStringPointerFromData(ipvs, "tcp_timeout")
				kubeProxyOutput.IPVS.UDPTimeout = helpers.GetNilStringPointerFromData(ipvs, "udp_timeout")
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
		pncOutput.DefaultVrackGateway = mapping["default_vrack_gateway"].(string)
		pncOutput.PrivateNetworkRoutingAsDefault = mapping["private_network_routing_as_default"].(bool)
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
	ControlPlaneIsUpToDate bool          `json:"controlPlaneIsUpToDate"`
	Id                     string        `json:"id"`
	IsUpToDate             bool          `json:"isUpToDate"`
	LoadBalancersSubnetId  string        `json:"loadBalancersSubnetId"`
	Name                   string        `json:"name"`
	NextUpgradeVersions    []string      `json:"nextUpgradeVersions"`
	NodesUrl               string        `json:"nodesUrl"`
	PrivateNetworkId       string        `json:"privateNetworkId"`
	Region                 string        `json:"region"`
	Status                 string        `json:"status"`
	UpdatePolicy           string        `json:"updatePolicy"`
	Url                    string        `json:"url"`
	Version                string        `json:"version"`
	Customization          Customization `json:"customization"`
	KubeProxyMode          string        `json:"kubeProxyMode"`
}

func (v *CloudProjectKubeResponse) ToMap(d *schema.ResourceData) map[string]interface{} {
	obj := make(map[string]interface{})
	obj["control_plane_is_up_to_date"] = v.ControlPlaneIsUpToDate
	obj["id"] = v.Id
	obj["is_up_to_date"] = v.IsUpToDate
	obj["load_balancers_subnet_id"] = v.LoadBalancersSubnetId
	obj["name"] = v.Name
	obj["next_upgrade_versions"] = v.NextUpgradeVersions
	obj["nodes_url"] = v.NodesUrl
	obj["private_network_id"] = v.PrivateNetworkId
	obj["region"] = v.Region
	obj["status"] = v.Status
	obj["update_policy"] = v.UpdatePolicy
	obj["url"] = v.Url
	obj["version"] = v.Version[:strings.LastIndex(v.Version, ".")]
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

	return obj
}

func loadKubeProxyCustomizationToMap(obj map[string]interface{}, v *CloudProjectKubeResponse) {
	obj[kubeClusterCustomizationKubeProxyKey] = []map[string]interface{}{{}}

	if v.Customization.KubeProxy.IPTables != nil {
		data := make(map[string]interface{})
		if vv := v.Customization.KubeProxy.IPTables.MinSyncPeriod; vv != nil && *vv != "" {
			data["min_sync_period"] = vv
		}

		if vv := v.Customization.KubeProxy.IPTables.SyncPeriod; vv != nil && *vv != "" {
			data["sync_period"] = vv
		}

		if len(data) > 0 {
			obj[kubeClusterCustomizationKubeProxyKey].([]map[string]interface{})[0]["iptables"] = []map[string]interface{}{data}
		}
	}

	if v.Customization.KubeProxy.IPVS != nil {
		data := make(map[string]interface{})
		if vv := v.Customization.KubeProxy.IPVS.MinSyncPeriod; vv != nil && *vv != "" {
			data["min_sync_period"] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.Scheduler; vv != nil && *vv != "" {
			data["scheduler"] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.SyncPeriod; vv != nil && *vv != "" {
			data["sync_period"] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.TCPFinTimeout; vv != nil && *vv != "" {
			data["tcp_fin_timeout"] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.TCPTimeout; vv != nil && *vv != "" {
			data["tcp_timeout"] = vv
		}

		if vv := v.Customization.KubeProxy.IPVS.UDPTimeout; vv != nil && *vv != "" {
			data["udp_timeout"] = vv
		}

		if len(data) > 0 {
			obj[kubeClusterCustomizationKubeProxyKey].([]map[string]interface{})[0]["ipvs"] = []map[string]interface{}{data}
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
			"apiserver": []map[string]interface{}{
				{
					"admissionplugins": []map[string]interface{}{
						{
							"enabled":  v.Customization.APIServer.AdmissionPlugins.Enabled,
							"disabled": v.Customization.APIServer.AdmissionPlugins.Disabled,
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
			"admissionplugins": []map[string]interface{}{
				{
					"enabled":  v.Customization.APIServer.AdmissionPlugins.Enabled,
					"disabled": v.Customization.APIServer.AdmissionPlugins.Disabled,
				},
			},
		},
	}
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
	APIServer *APIServer              `json:"apiServer,omitempty"`
	KubeProxy *kubeProxyCustomization `json:"kubeProxy,omitempty"`
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
	obj["created_at"] = v.CreatedAt
	obj["deployed_at"] = v.DeployedAt
	obj["flavor"] = v.Flavor
	obj["id"] = v.Id
	obj["instance_id"] = v.InstanceId
	obj["is_up_to_date"] = v.IsUpToDate
	obj["name"] = v.Name
	obj["node_pool_id"] = v.NodePoolId
	obj["project_id"] = v.ProjectId
	obj["status"] = v.Status
	obj["updated_at"] = v.UpdatedAt
	obj["version"] = v.Version
	return obj
}
