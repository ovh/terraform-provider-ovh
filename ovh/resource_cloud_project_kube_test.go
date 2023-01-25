package ovh

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

var testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsCreateConfig = `
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

var testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsUpdateConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	customization {
		apiserver {
			admissionplugins {
				enabled = ["AlwaysPullImages","NodeRestriction"]
			}
		}
	}
}
`

var testAccCloudProjectKubeKubeProxyIPVS = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"

	kube_proxy_mode = "ipvs"
	customization {
		kube_proxy {
			ipvs {
				min_sync_period = "PT0S"
				sync_period = "PT30S"
				scheduler = "rr"
				tcp_fin_timeout = "PT0S"
				tcp_timeout = "PT0S"
				udp_timeout = "PT0S"
			}

			iptables {
				sync_period = "PT30S"
				min_sync_period = "PT0S"
			}
		}
	}
}
`

var testAccCloudProjectKubeUpdatedKubeProxyIPVS = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"

	kube_proxy_mode = "ipvs"
	customization {
		kube_proxy {
			ipvs {
				min_sync_period = "PT30S"
				sync_period = "PT10S"
				scheduler = "rr"
				tcp_fin_timeout = "PT30S"
				tcp_timeout = "PT30S"
				udp_timeout = "PT30S"
			}

			iptables {
				sync_period = "PT30S"
				min_sync_period = "PT0S"
			}
		}
	}
}
`

var testAccCloudProjectKubeKubeProxyIPTables = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"

	kube_proxy_mode = "iptables"
	customization {
		kube_proxy {
			ipvs {
				min_sync_period = "PT0S"
				sync_period = "PT30S"
				scheduler = "rr"
				tcp_fin_timeout = "PT0S"
				tcp_timeout = "PT0S"
				udp_timeout = "PT0S"
			}

			iptables {
				sync_period = "PT30S"
				min_sync_period = "PT0S"
			}
		}
	}
}
`

var testAccCloudProjectKubeUpdatedKubeProxyIPTables = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"

	kube_proxy_mode = "iptables"
	customization {
		kube_proxy {
			ipvs {
				min_sync_period = "PT0S"
				sync_period = "PT30S"
				scheduler = "rr"
				tcp_fin_timeout = "PT0S"
				tcp_timeout = "PT0S"
				udp_timeout = "PT0S"
			}

			iptables {
				sync_period = "PT60S"
				min_sync_period = "PT10S"
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
}

func TestAccCloudProjectKubeCustomizationApiServerAdmissionPlugins(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsCreateConfig,
		serviceName,
		name,
		region,
	)
	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeCustomizationApiServerAdmissionPluginsUpdateConfig,
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
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.0", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.0", "AlwaysPullImages"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.0", "AlwaysPullImages"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.enabled.1", "NodeRestriction"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.apiserver.0.admissionplugins.0.disabled.#", "0"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeVRack(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	vrackID := os.Getenv("OVH_VRACK_SERVICE_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	tmpl, err := template.New("config").Parse(testAccCloudProjectKubeVRackConfig)
	if err != nil {
		panic(err)
	}

	var config bytes.Buffer
	var configUpdated bytes.Buffer
	configData1 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         "GRA5",
		DefaultVrackGateway:            "",
		PrivateNetworkRoutingAsDefault: false,
	}
	configData2 := configData{
		ServiceName:                    serviceName,
		VrackID:                        vrackID,
		Name:                           name,
		Region:                         "GRA5",
		DefaultVrackGateway:            "10.4.0.1",
		PrivateNetworkRoutingAsDefault: true,
	}

	err = tmpl.Execute(&config, &configData1)
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
				),
			},
		},
	})
}

func TestAccCloudProjectKube_kube_proxy_ipvs(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeKubeProxyIPVS,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	updatableConfig := fmt.Sprintf(
		testAccCloudProjectKubeUpdatedKubeProxyIPVS,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
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
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.tcp_timeout", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.udp_timeout", "PT0S"),
				),
			},
			{
				Config: updatableConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.sync_period", "PT10S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.tcp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.ipvs.0.udp_timeout", "PT30S"),
				),
			},
		},
	})
}

func TestAccCloudProjectKube_kube_proxy_iptables(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeKubeProxyIPTables,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeUpdatedKubeProxyIPTables,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
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
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.iptables.0.min_sync_period", "PT0S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.iptables.0.sync_period", "PT30S"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "kube_proxy_mode", "iptables"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.iptables.0.min_sync_period", "PT10S"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "customization.0.kube_proxy.0.iptables.0.sync_period", "PT60S"),
				),
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

	version1 := "1.22"
	version2 := "1.23"

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
