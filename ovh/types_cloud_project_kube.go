package ovh

import (
	"encoding/json"
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

type CloudProjectKubeCreateOpts struct {
	Name                        *string                      `json:"name,omitempty"`
	PrivateNetworkId            *string                      `json:"privateNetworkId,omitempty"`
	PrivateNetworkConfiguration *privateNetworkConfiguration `json:"privateNetworkConfiguration,omitempty"`
	Region                      string                       `json:"region"`
	Version                     *string                      `json:"version,omitempty"`
	UpdatePolicy                *string                      `json:"updatePolicy,omitempty"`
	Customization               *Customization               `json:"customization,omitempty"`
	KubeProxyMode               *string                      `json:"kubeProxyMode"`
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

func (opts *CloudProjectKubeCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeCreateOpts {
	opts.Region = d.Get("region").(string)
	opts.Version = helpers.GetNilStringPointerFromData(d, "version")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	opts.UpdatePolicy = helpers.GetNilStringPointerFromData(d, "update_policy")
	opts.PrivateNetworkId = helpers.GetNilStringPointerFromData(d, "private_network_id")
	opts.PrivateNetworkConfiguration = loadPrivateNetworkConfiguration(d.Get("private_network_configuration"))
	opts.Customization = loadCustomization(d.Get(kubeClusterCustomizationApiServerKey), d.Get(kubeClusterCustomizationKubeProxyKey))
	opts.KubeProxyMode = helpers.GetNilStringPointerFromData(d, kubeClusterProxyModeKey)
	return opts
}

func loadCustomization(apiServerAdmissionPlugins interface{}, kubeProxyCustomizationInterface interface{}) *Customization {
	if apiServerAdmissionPlugins == nil && kubeProxyCustomizationInterface == nil {
		return nil
	}

	customizationOutput := Customization{
		APIServer: &APIServer{
			AdmissionPlugins: &AdmissionPlugins{},
		},
		KubeProxy: &kubeProxyCustomization{
			IPTables: &kubeProxyCustomizationIPTables{},
			IPVS:     &kubeProxyCustomizationIPVS{},
		},
	}

	// Customization
	customizationSet := apiServerAdmissionPlugins.(*schema.Set).List()
	if len(customizationSet) > 0 {
		customization := customizationSet[0].(map[string]interface{})
		admissionPluginsSet := customization["admissionplugins"].(*schema.Set).List()
		admissionPlugins := admissionPluginsSet[0].(map[string]interface{})

		// Enabled admission plugins
		{
			stringArray := admissionPlugins["enabled"].([]interface{})
			enabled := []string{}
			for _, s := range stringArray {
				enabled = append(enabled, s.(string))
			}
			customizationOutput.APIServer.AdmissionPlugins.Enabled = &enabled
		}

		// Disabled admission plugins
		{
			stringArray := admissionPlugins["disabled"].([]interface{})
			disabled := []string{}
			for _, s := range stringArray {
				disabled = append(disabled, s.(string))
			}
			customizationOutput.APIServer.AdmissionPlugins.Disabled = &disabled
		}
	}

	// Nested KubeProxy customization
	kubeProxySet := kubeProxyCustomizationInterface.(*schema.Set).List()
	if len(kubeProxySet) > 0 {
		kubeProxy := kubeProxySet[0].(map[string]interface{})

		// Nested IPTables customization
		{
			ipTablesSet := kubeProxy["iptables"].(*schema.Set).List()
			if len(ipTablesSet) > 0 {
				ipTables := ipTablesSet[0].(map[string]interface{})
				customizationOutput.KubeProxy.IPTables.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipTables, "min_sync_period")
				customizationOutput.KubeProxy.IPTables.SyncPeriod = helpers.GetNilStringPointerFromData(ipTables, "sync_period")
			}
		}

		// Nested IPVS customization
		{
			ipvsSet := kubeProxy["ipvs"].(*schema.Set).List()
			if len(ipvsSet) > 0 {
				ipvs := ipvsSet[0].(map[string]interface{})
				customizationOutput.KubeProxy.IPVS.MinSyncPeriod = helpers.GetNilStringPointerFromData(ipvs, "min_sync_period")
				customizationOutput.KubeProxy.IPVS.Scheduler = helpers.GetNilStringPointerFromData(ipvs, "scheduler")
				customizationOutput.KubeProxy.IPVS.SyncPeriod = helpers.GetNilStringPointerFromData(ipvs, "sync_period")
				customizationOutput.KubeProxy.IPVS.TCPFinTimeout = helpers.GetNilStringPointerFromData(ipvs, "tcp_fin_timeout")
				customizationOutput.KubeProxy.IPVS.TCPTimeout = helpers.GetNilStringPointerFromData(ipvs, "tcp_timeout")
				customizationOutput.KubeProxy.IPVS.UDPTimeout = helpers.GetNilStringPointerFromData(ipvs, "udp_timeout")
			}
		}
	}

	return &customizationOutput
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

func (s *CloudProjectKubeCreateOpts) String() string {
	return fmt.Sprintf("%s(%s): %s", *s.Name, s.Region, *s.Version)
}

type CloudProjectKubeResponse struct {
	ControlPlaneIsUpToDate bool          `json:"controlPlaneIsUpToDate"`
	Id                     string        `json:"id"`
	IsUpToDate             bool          `json:"isUpToDate"`
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

func (v CloudProjectKubeResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["control_plane_is_up_to_date"] = v.ControlPlaneIsUpToDate
	obj["id"] = v.Id
	obj["is_up_to_date"] = v.IsUpToDate
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

	obj["customization_apiserver"] = []map[string]interface{}{
		{
			"admissionplugins": []map[string]interface{}{
				{
					"enabled":  v.Customization.APIServer.AdmissionPlugins.Enabled,
					"disabled": v.Customization.APIServer.AdmissionPlugins.Disabled,
				},
			},
		},
	}

	obj["customization_kube_proxy"] = []map[string]interface{}{
		{
			"iptables": []map[string]interface{}{
				{
					"min_sync_period": v.Customization.KubeProxy.IPTables.MinSyncPeriod,
					"sync_period":     v.Customization.KubeProxy.IPTables.SyncPeriod,
				},
			},
			"ipvs": []map[string]interface{}{
				{
					"min_sync_period": v.Customization.KubeProxy.IPVS.MinSyncPeriod,
					"scheduler":       v.Customization.KubeProxy.IPVS.Scheduler,
					"sync_period":     v.Customization.KubeProxy.IPVS.SyncPeriod,
					"tcp_fin_timeout": v.Customization.KubeProxy.IPVS.TCPFinTimeout,
					"tcp_timeout":     v.Customization.KubeProxy.IPVS.TCPTimeout,
					"udp_timeout":     v.Customization.KubeProxy.IPVS.UDPTimeout,
				},
			},
		},
	}

	if len(obj["customization_kube_proxy"].([]map[string]interface{})) > 0 {

		if len(obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})) > 0 {
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})[0]["min_sync_period"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})[0], "min_sync_period")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})[0]["sync_period"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})[0], "sync_period")
			}
			if len(obj["customization_kube_proxy"].([]map[string]interface{})[0]["iptables"].([]map[string]interface{})[0]) == 0 {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0], "iptables")
			}
		}

		if len(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})) > 0 {
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["min_sync_period"].(*string) == nil ||
				*obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["min_sync_period"].(*string) == "" {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "min_sync_period")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["scheduler"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "scheduler")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["sync_period"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "sync_period")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["tcp_fin_timeout"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "tcp_fin_timeout")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["tcp_timeout"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "tcp_timeout")
			}
			if obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]["udp_timeout"].(*string) == nil {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0], "udp_timeout")
			}
			if len(obj["customization_kube_proxy"].([]map[string]interface{})[0]["ipvs"].([]map[string]interface{})[0]) == 0 {
				delete(obj["customization_kube_proxy"].([]map[string]interface{})[0], "ipvs")
			}
		}

		if len(obj["customization_kube_proxy"].([]map[string]interface{})[0]) == 0 {
			delete(obj, "customization_kube_proxy")
		}
	}

	i, _ := json.MarshalIndent(v, "", "  ")
	o, _ := json.MarshalIndent(obj, "", "  ")

	log.Printf("##### %s\n\n%s\n\n%s\n\n", i, o)

	return obj
}

func (s *CloudProjectKubeResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Name, s.Id, s.Status)
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
