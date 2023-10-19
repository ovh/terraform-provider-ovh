package ovh

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_project_kube", &resource.Sweeper{
		Name: "ovh_cloud_project_kube",
		Dependencies: []string{
			"ovh_cloud_project_kube_nodepool",
		},
		F: testSweepCloudProjectKube,
	})
}

func testSweepCloudProjectKube(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No kube to sweep")
		return nil
	}

	kubeIds := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/cloud/project/%s/kube", serviceName), &kubeIds); err != nil {
		return fmt.Errorf("Error calling GET /cloud/project/%s/kube:\n\t %q", serviceName, err)
	}
	for _, kubeId := range kubeIds {
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
		res := &CloudProjectKubeResponse{}
		log.Printf("[DEBUG] read kube %s from project: %s", kubeId, serviceName)
		if err := client.Get(endpoint, res); err != nil {
			return err
		}
		if !strings.HasPrefix(res.Name, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}

var testAccCloudProjectKubeUpdatePolicyConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	update_policy = "%s"
}
`

var testAccCloudProjectKubeUpdateLoadBalancersSubnetIdConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	load_balancers_subnet_id = "%s"
}
`

var testAccCloudProjectKubeNodesSubnetIdCreateConfigDefaultValues = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	nodes_subnet_id = "%s"
}
`

var testAccCloudProjectKubeConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	version = "%s"
}
`

var testAccCloudProjectKubeEmptyVersionConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	update_policy = "ALWAYS_UPDATE"
}
`

var testAccCloudProjectKubeVRackConfig = `
resource "ovh_vrack_cloudproject" "attach" {
	service_name = "{{ .VrackID }}"
	project_id   = "{{ .ServiceName }}"
}

resource "ovh_cloud_project_network_private" "network" {
	service_name = "{{ .ServiceName }}"
	vlan_id    = 0
	name       = "terraform_testacc_private_net"
	regions    = ["{{ .Region }}"]
	depends_on = [ovh_vrack_cloudproject.attach]
}

resource "ovh_cloud_project_network_private_subnet" "networksubnet" {
  service_name = ovh_cloud_project_network_private.network.service_name
  network_id   = ovh_cloud_project_network_private.network.id

  # whatever region, for test purpose
  region     = "{{ .Region }}"
  start      = "192.168.168.100"
  end        = "192.168.168.200"
  network    = "192.168.168.0/24"
  dhcp       = true
  no_gateway = false

  depends_on   = [ovh_cloud_project_network_private.network]
}

resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "{{ .ServiceName }}"
	name          = "{{ .Name }}"
	region        = "{{ .Region }}"

	private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]

	private_network_configuration {
		default_vrack_gateway              = "{{ .DefaultVrackGateway }}"
		private_network_routing_as_default = {{ .PrivateNetworkRoutingAsDefault }}
	}

	depends_on = [
		ovh_cloud_project_network_private.network
	]
}
`

var testAccCloudProjectKubeVRackWithSubnetsConfig = `
resource "ovh_vrack_cloudproject" "attach" {
	service_name = "{{ .VrackID }}"
	project_id   = "{{ .ServiceName }}"
}

resource "ovh_cloud_project_network_private" "network" {
	service_name = "{{ .ServiceName }}"
	vlan_id    = 0
	name       = "terraform_testacc_private_net"
	regions    = ["{{ .Region }}"]
	depends_on = [ovh_vrack_cloudproject.attach]
}

resource "ovh_cloud_project_network_private_subnet" "networksubnetForNodes" {
	service_name = ovh_cloud_project_network_private.network.service_name
	network_id   = ovh_cloud_project_network_private.network.id

	# whatever region, for test purpose
	region     = "{{ .Region }}"
	start      = "192.168.168.64"
	end        = "192.168.168.127"
	network    = "192.168.168.64/26"
	dhcp       = true
	no_gateway = false

	depends_on   = [ovh_cloud_project_network_private.network]
}

resource "ovh_cloud_project_network_private_subnet" "networksubnetForLoadbalancers1" {
	service_name = ovh_cloud_project_network_private.network.service_name
	network_id   = ovh_cloud_project_network_private.network.id

	# whatever region, for test purpose
	region     = "{{ .Region }}"
	start      = "192.168.168.128"
	end        = "192.168.168.191"
	network    = "192.168.168.128/26"
	dhcp       = true
	no_gateway = false

	depends_on   = [ovh_cloud_project_network_private.network]
}

resource "ovh_cloud_project_network_private_subnet" "networksubnetForLoadbalancers2" {
    service_name = ovh_cloud_project_network_private.network.service_name
    network_id   = ovh_cloud_project_network_private.network.id

    # whatever region, for test purpose
    region     = "{{ .Region }}"
    start      = "192.168.168.192"
    end        = "192.168.168.255"
    network    = "192.168.168.192/26"
    dhcp       = true
    no_gateway = false

    depends_on   = [ovh_cloud_project_network_private.network]
}

resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "{{ .ServiceName }}"
	name          = "{{ .Name }}"
	region        = "{{ .Region }}"

	private_network_configuration {
		default_vrack_gateway              = "{{ .DefaultVrackGateway }}"
		private_network_routing_as_default = {{ .PrivateNetworkRoutingAsDefault }}
	}

	nodes_subnet_id = ovh_cloud_project_network_private.networksubnetForNodes
	load_balancers_subnet_id = ovh_cloud_project_network_private.{{ .LoadBalancersSubnetName }}

	depends_on = [
		ovh_cloud_project_network_private.network,
		ovh_cloud_project_network_private_subnet.{{ .LoadBalancersSubnetName }}
	]
}
`

var testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsCreateConfigDefaultValues = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}
`

var testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsUpdateConfigEnabledAndDisabled = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	customization_apiserver {
		admissionplugins {
			enabled = ["NodeRestriction"]
			disabled = ["AlwaysPullImages"]
		}
	}
}
`

var testAccCloudProjectKubeDeprecatedCustomizationApiServerAdmissionPluginsUpdateConfigEnabledAndDisabled = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	customization {
		apiserver {
			admissionplugins {
				enabled = ["NodeRestriction"]
				disabled = ["AlwaysPullImages"]
			}
		}
	}
}
`

type configData struct {
	Region                         string
	Regions                        string
	ServiceName                    string
	VrackID                        string
	Name                           string
	DefaultVrackGateway            string
	PrivateNetworkRoutingAsDefault bool
	LoadBalancersSubnetName        string
}

func TestAccCloudProjectKubeCustomizationApiServerAdmissionPlugins(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	createConfig := fmt.Sprintf(
		testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsCreateConfigDefaultValues,
		serviceName,
		name,
		region,
	)

	updatedConfigEnabledAndDisabled := fmt.Sprintf(
		testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsUpdateConfigEnabledAndDisabled,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// no apiserver customization, should contain default values from API
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.disabled.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.0", "AlwaysPullImages"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.1", "NodeRestriction"),

					// Conflicts with the old schema
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.#", "0"),
				),
			},
			{
				Config: updatedConfigEnabledAndDisabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.0", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.disabled.0", "AlwaysPullImages"),

					// Conflicts with the old schema
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.#", "0"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeNodesSubnetId(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)
	nodesSubnetId := "158c9998-18da-4bef-a8db-4891b1736574"

	createConfig := fmt.Sprintf(
		testAccCloudProjectKubeNodesSubnetIdCreateConfigDefaultValues,
		serviceName,
		name,
		region,
		nodesSubnetId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// no apiserver customization, should contain default values from API
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "nodes_subnet_id", nodesSubnetId),

					// Conflicts with the old schema
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.#", "0"),
				),
			},
		},
	})
}

// TestAccCloudProjectKubeDeprecatedCustomizationApiServerAdmissionPlugins aims to test that
// values are the same between customization_apiserver.admissionplugins and customization.apiserver.admissionplugins.
// This is deprecated and will be removed in the future.
func TestAccCloudProjectKubeDeprecatedCustomizationApiServerAdmissionPlugins(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	createConfig := fmt.Sprintf(
		testAccCloudProjectKubeDeprecatedCustomizationApiServerAdmissionPluginsUpdateConfigEnabledAndDisabled,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),

					// Deprecated configuration
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.0", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.0", "AlwaysPullImages"),

					// Conflicts with the new schema
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.#", "0"),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_kube_proxy_iptables(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "iptables"
	customization_kube_proxy {
		iptables {
        	min_sync_period = "PT0S"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfigWithDifferentTime := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "iptables"
	customization_kube_proxy {
		iptables {
        	min_sync_period = "P0D"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfigNewArgument := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "iptables"
	customization_kube_proxy {
		iptables {
        	min_sync_period = "PT30S"
			sync_period = "PT30S"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// no kube proxy mode specified, should contain default values from API
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", ".customization_kube_proxy.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.#", "0"),
				),
			},
			{
				Config: updatedConfigWithDifferentTime,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.#", "0"),
				),
			},
			{
				Config: updatedConfigNewArgument,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.#", "0"),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_kube_proxy_ipvs(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	kube_proxy_mode = "ipvs"
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "ipvs"
	customization_kube_proxy {
		ipvs {
        	min_sync_period = "PT0S"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfigWithDifferentTime := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "ipvs"
	customization_kube_proxy {
		ipvs {
        	min_sync_period = "P0D"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfigAllArguments := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
	service_name    = "%s"
	name            = "%s"
	region          = "%s"
	
	kube_proxy_mode = "ipvs"
	customization_kube_proxy {
		ipvs {
        	min_sync_period = "PT30S"
			sync_period = "PT30S"
			scheduler = "rr"
			tcp_fin_timeout = "PT30S"
			tcp_timeout = "PT30S"
			udp_timeout = "PT30S"
		}
    }
}
`,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// no kube proxy mode specified, should contain default values from API
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", ".customization_kube_proxy.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.#", "0"),
				),
			},
			{
				Config: updatedConfigWithDifferentTime,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", ""),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.#", "0"),
				),
			},
			{
				Config: updatedConfigAllArguments,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.#", "0"),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_customization_full_deprecated(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	erroredConfigKubeProxyMode := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"
  kube_proxy_mode = "foo"
}
`,
		serviceName,
		name,
		region,
	)

	erroredConfigInvalidRFC3339Duration := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"

  customization_kube_proxy {
    iptables {
      min_sync_period = "foo"
      sync_period     = "foo"
    }
    ipvs {
      min_sync_period = "foo"
      scheduler       = "rr"
      sync_period     = "foo"
      tcp_fin_timeout = "foo"
      tcp_timeout     = "foo"
      udp_timeout     = "foo"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	erroredConfigInvalidScheduler := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"

  customization_kube_proxy {
    ipvs {
      scheduler       = "foo"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	config := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"
  kube_proxy_mode = "iptables"

  customization {
    apiserver {
      admissionplugins {
        enabled  = ["NodeRestriction"]
        disabled = ["AlwaysPullImages"]
      }
    }
  }

  customization_kube_proxy {
    iptables {
      min_sync_period = "PT0S"
      sync_period     = "PT0S"
    }
    ipvs {
      min_sync_period = "PT0S"
      scheduler       = "rr"
      sync_period     = "PT0S"
      tcp_fin_timeout = "PT0S"
      tcp_timeout     = "PT0S"
      udp_timeout     = "PT0S"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"
  kube_proxy_mode = "iptables"

  customization {
    apiserver {
	  admissionplugins {
	    enabled  = ["AlwaysPullImages", "NodeRestriction"]
	    disabled = []
	  }
    }
  }

  customization_kube_proxy {
    iptables {
      min_sync_period = "PT30S"
      sync_period     = "PT30S"
    }
    ipvs {
      min_sync_period = "PT30S"
      scheduler       = "rr"
      sync_period     = "PT30S"
      tcp_fin_timeout = "PT30S"
      tcp_timeout     = "PT30S"
      udp_timeout     = "PT30S"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      erroredConfigKubeProxyMode,
				ExpectError: regexp.MustCompile(`is not among valid values`),
			},
			{
				Config:      erroredConfigInvalidRFC3339Duration,
				ExpectError: regexp.MustCompile(`does not match RFC3339 duration`),
			},
			{
				Config:      erroredConfigInvalidScheduler,
				ExpectError: regexp.MustCompile(`is not among valid values`),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),

					// customization_kube_proxy - ipvs
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT0S"),

					// customization_kube_proxy - iptables
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT0S"),

					// customization - apiserver
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.0", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.0", "AlwaysPullImages"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),

					// customization_kube_proxy - ipvs
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT30S"),

					// customization_kube_proxy - iptables
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT30S"),

					// customization - apiserver
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.0", "AlwaysPullImages"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.1", "NodeRestriction"),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_customization_full(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"
  kube_proxy_mode = "iptables"

  customization_apiserver {
    admissionplugins {
      enabled  = ["NodeRestriction"]
      disabled = ["AlwaysPullImages"]
    }
  }

  customization_kube_proxy {
    iptables {
      min_sync_period = "PT0S"
      sync_period     = "PT0S"
    }
    ipvs {
      min_sync_period = "PT0S"
      scheduler       = "rr"
      sync_period     = "PT0S"
      tcp_fin_timeout = "PT0S"
      tcp_timeout     = "PT0S"
      udp_timeout     = "PT0S"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_project_kube" "cluster" {
  service_name    = "%s"
  name            = "%s"
  region          = "%s"
  kube_proxy_mode = "iptables"

  customization_apiserver {
	admissionplugins {
	  enabled  = ["AlwaysPullImages", "NodeRestriction"]
	  disabled = []
	}
  }

  customization_kube_proxy {
    iptables {
      min_sync_period = "PT30S"
      sync_period     = "PT30S"
    }
    ipvs {
      min_sync_period = "PT30S"
      scheduler       = "rr"
      sync_period     = "PT30S"
      tcp_fin_timeout = "PT30S"
      tcp_timeout     = "PT30S"
      udp_timeout     = "PT30S"
    }
  }
}
`,
		serviceName,
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),

					// customization_kube_proxy - ipvs
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT0S"),

					// customization_kube_proxy - iptables
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT0S"),

					// customization - apiserver
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.0", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.disabled.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.disabled.0", "AlwaysPullImages"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),

					// customization_kube_proxy - ipvs
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT30S"),

					// customization_kube_proxy - iptables
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT30S"),

					// customization - apiserver
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.disabled.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.0", "AlwaysPullImages"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization_apiserver.0.admissionplugins.0.enabled.1", "NodeRestriction"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeVRack(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	vrackID := os.Getenv("OVH_VRACK_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	tmpl := template.Must(template.New("config").Parse(testAccCloudProjectKubeVRackConfig))

	var config bytes.Buffer
	var configUpdated bytes.Buffer
	configData1 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         region,
		DefaultVrackGateway:            "",
		PrivateNetworkRoutingAsDefault: false,
	}
	configData2 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         region,
		DefaultVrackGateway:            "10.4.0.1",
		PrivateNetworkRoutingAsDefault: true,
	}

	err := tmpl.Execute(&config, &configData1)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&configUpdated, &configData2)
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
			testAccPreCheckKubernetesVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData1.DefaultVrackGateway),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData1.PrivateNetworkRoutingAsDefault)),
				),
			},
			{
				Config: configUpdated.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData2.DefaultVrackGateway),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData2.PrivateNetworkRoutingAsDefault)),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	resourceName := "ovh_cloud_project_kube.cluster"
	config := fmt.Sprintf(
		testAccCloudProjectKubeConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "version", version),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.host"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_certificate"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), state.RootModule().Resources[resourceName].Primary.ID), nil
				},
				ImportStateVerifyIgnore: []string{"kubeconfig"}, // we must ignore kubeconfig because certificate is not the same
			},
		},
	})
}

// TestAccCloudProjectKubeEmptyVersion_basic
// create a public cluster
// check some properties
// update cluster name
// check some properties && cluster updated name
func TestAccCloudProjectKubeEmptyVersion_basic(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccCloudProjectKubeEmptyVersionConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	updatedName := acctest.RandomWithPrefix(test_prefix)
	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeEmptyVersionConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		updatedName,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, updatedName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeUpdatePolicy_basic(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccCloudProjectKubeUpdatePolicyConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		"ALWAYS_UPDATE",
	)

	updatedName := acctest.RandomWithPrefix(test_prefix)
	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeUpdatePolicyConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		updatedName,
		region,
		"NEVER_UPDATE",
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "update_policy", "ALWAYS_UPDATE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "update_policy", "NEVER_UPDATE"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeUpdateVersion_basic(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	updatedName := acctest.RandomWithPrefix(test_prefix)

	version1 := os.Getenv("OVH_CLOUD_PROJECT_KUBE_PREV_VERSION_TEST")
	version2 := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectKubeConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version1,
	)

	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		updatedName,
		region,
		version2,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterVersionKey, version1),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterVersionKey, version2),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeUpdateLoadBalancersSubnetId_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	vrackID := os.Getenv("OVH_VRACK_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	tmpl := template.Must(template.New("config").Parse(testAccCloudProjectKubeVRackWithSubnetsConfig))

	var config bytes.Buffer
	var configUpdated bytes.Buffer
	configData1 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         region,
		DefaultVrackGateway:            "10.4.0.1",
		PrivateNetworkRoutingAsDefault: true,
		LoadBalancersSubnetName:        "networksubnetForLoadbalancers1",
	}
	configData2 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         region,
		DefaultVrackGateway:            "10.4.0.1",
		PrivateNetworkRoutingAsDefault: true,
		// Update load_balancers_subnet_id to networksubnetForLoadbalancers2
		LoadBalancersSubnetName: "networksubnetForLoadbalancers2",
	}

	err := tmpl.Execute(&config, &configData1)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&configUpdated, &configData2)
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
			testAccPreCheckKubernetesVRack(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData1.DefaultVrackGateway),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData1.PrivateNetworkRoutingAsDefault)),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "load_balancers_subnet_id", "ovh_cloud_project_network_private_subnet.networksubnetForLoadbalancers1"),
				),
			},
			{
				Config: configUpdated.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData2.DefaultVrackGateway),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData2.PrivateNetworkRoutingAsDefault)),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "load_balancers_subnet_id", "ovh_cloud_project_network_private_subnet.networksubnetForLoadbalancers2"),
				),
			},
		},
	})
}

func TestCustomIPVSIPTablesSchemaSetFunc(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expectedHex string
	}{
		{
			name: "Input with P0D value",
			input: map[string]interface{}{
				"key1": "P0D",
				"key2": "value2",
			},
			expectedHex: fmt.Sprintf("%#v", map[string]interface{}{
				"key1": "PT0S",
				"key2": "value2",
			}),
		},
		{
			name: "Input without P0D value",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectedHex: fmt.Sprintf("%#v", map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHash := schema.HashString(tt.expectedHex)
			if got := CustomIPVSIPTablesSchemaSetFunc()(tt.input); got != expectedHash {
				t.Errorf("CustomIPVSIPTablesSchemaSetFunc() = %v, want %v", got, expectedHash)
			}
		})
	}
}

func TestCustomApiServerAdmissionPluginsSchemaSetFunc(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expectedHex string
	}{
		{
			name: "No plugins",
			input: map[string]interface{}{
				"enabled":  []interface{}{},
				"disabled": []interface{}{},
			},
			expectedHex: fmt.Sprintf("%#v", map[string]interface{}{
				"enabled":  []interface{}{},
				"disabled": []interface{}{},
			}),
		},
		{
			name: "Should reorder plugins",
			input: map[string]interface{}{
				"enabled":  []interface{}{"foo", "bar"},
				"disabled": []interface{}{"bar", "foo"},
			},
			expectedHex: fmt.Sprintf("%#v", map[string]interface{}{
				"enabled":  []interface{}{"bar", "foo"},
				"disabled": []interface{}{"bar", "foo"},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHash := schema.HashString(tt.expectedHex)
			if got := CustomApiServerAdmissionPluginsSchemaSetFunc()(tt.input); got != expectedHash {
				t.Errorf("CustomApiServerAdmissionPluginsSchemaSetFunc() = %v, want %v", got, expectedHash)
			}
		})
	}
}
