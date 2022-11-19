package ovh

import (
	"fmt"
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
}

type Customization struct {
	APIServer *APIServer `json:"apiServer,omitempty"`
}

type APIServer struct {
	AdmissionPlugins *AdmissionPlugins `json:"admissionPlugins,omitempty"`
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
	opts.Customization = loadCustomization(d.Get("customization"))
	return opts
}

func loadCustomization(i interface{}) *Customization {
	if i == nil {
		return nil
	}

	customizationOutput := Customization{
		APIServer: &APIServer{
			AdmissionPlugins: &AdmissionPlugins{},
		},
	}

	customizationSet := i.(*schema.Set).List()
	for _, customization := range customizationSet {
		apiServerSet := customization.(map[string]interface{})["apiserver"].(*schema.Set).List()
		for _, apiServer := range apiServerSet {
			admissionPluginsSet := apiServer.(map[string]interface{})["admissionplugins"].(*schema.Set).List()
			for _, admissionPlugins := range admissionPluginsSet {

				stringArray := admissionPlugins.(map[string]interface{})["enabled"].([]interface{})
				enabled := []string{}
				for _, s := range stringArray {
					enabled = append(enabled, s.(string))
				}
				customizationOutput.APIServer.AdmissionPlugins.Enabled = &enabled

				stringArray = admissionPlugins.(map[string]interface{})["disabled"].([]interface{})
				disabled := []string{}
				for _, s := range stringArray {
					disabled = append(disabled, s.(string))
				}
				customizationOutput.APIServer.AdmissionPlugins.Disabled = &disabled

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
	obj["customization"] = []map[string]interface{}{
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
	APIServer *APIServer `json:"apiServer"`
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
