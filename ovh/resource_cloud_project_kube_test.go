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
}
`

var testAccCloudProjectKubeVRackConfig = `
resource "ovh_cloud_project_network_private" "network1" {
  service_name = "{{ .ServiceName }}"
  name         = "dhcp-default-gateway"
  regions      = {{ .Regions }}
  vlan_id      = "{{ .Vlanid }}"
}

resource "ovh_cloud_project_network_private_subnet" "network1subnetSBG5" {
  service_name = "{{ .ServiceName }}"
  network_id   = ovh_cloud_project_network_private.network1.id
  start        = "10.6.0.2" #first ip is burn for gateway ip
  end          = "10.6.255.254"
  network      = "10.6.0.0/16"
  dhcp         = true
  region       = "{{ .Region }}"
  no_gateway   = false
  depends_on   = [ovh_cloud_project_network_private.network1]
}

resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "{{ .ServiceName }}"
    name          = "{{ .Name }}"
	region        = "{{ .Region }}"

    private_network_id = ovh_cloud_project_network_private.network1.regions_attributes[index(ovh_cloud_project_network_private.network1.regions_attributes.*.region, "{{ .Region }}")].openstackid

    private_network_configuration {
      default_vrack_gateway              = "{{ .DefaultVrackGateway }}"
      private_network_routing_as_default = {{ .PrivateNetworkRoutingAsDefault }}
    }

    depends_on = [
      ovh_cloud_project_network_private.network1
    ]
}
`

type configData struct {
	Region                         string
	Regions                        string
	Vlanid                         string
	ServiceName                    string
	Name                           string
	DefaultVrackGateway            string
	PrivateNetworkRoutingAsDefault bool
}

func TestAccCloudProjectKubeVRack(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	tmp, err := template.New("config").Parse(testAccCloudProjectKubeVRackConfig)
	if err != nil {
		panic(err)
	}

	var config bytes.Buffer
	var configUpdated bytes.Buffer
	configData1 := configData{
		Region:                         region,
		Regions:                        `["` + region + `"]`,
		Vlanid:                         "6",
		ServiceName:                    serviceName,
		Name:                           name,
		DefaultVrackGateway:            "",
		PrivateNetworkRoutingAsDefault: false,
	}
	configData2 := configData{
		Region:                         region,
		Regions:                        `["` + region + `"]`,
		Vlanid:                         "6",
		ServiceName:                    serviceName,
		Name:                           name,
		DefaultVrackGateway:            "10.4.0.1",
		PrivateNetworkRoutingAsDefault: true,
	}

	err = tmp.Execute(&config, &configData1)
	if err != nil {
		panic(err)
	}

	err = tmp.Execute(&configUpdated, &configData2)
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData1.DefaultVrackGateway),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster",
						"private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData1.PrivateNetworkRoutingAsDefault)),
				),
			},
			{
				Config: configUpdated.String(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "private_network_configuration.0.default_vrack_gateway", configData2.DefaultVrackGateway),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster",
						"private_network_configuration.0.private_network_routing_as_default", strconv.FormatBool(configData2.PrivateNetworkRoutingAsDefault)),
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
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "version", version),
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
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, updatedName),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
				),
			},
		},
	})
}
