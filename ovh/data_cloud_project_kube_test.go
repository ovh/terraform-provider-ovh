package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectKubeDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeDatasourceConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)
	matchVersion := regexp.MustCompile(`^` + version + `(\..*)?$`)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "name", name),
					resource.TestMatchResourceAttr("data.ovh_cloud_project_kube.cluster", "version", matchVersion),

					// Check kubeconfig is present and not empty
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig"),

					// Check kubeconfig_attributes are present
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.host"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_certificate"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_key"),
				),
			},
		},
	})
}

func TestAccCloudProjectKubeDataSource_kubeProxy(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeDatasourceKubeProxyConfig,
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
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "kube_proxy_mode", "ipvs"),

					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.iptables.0.sync_period", "PT30S"),

					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.min_sync_period", "PT30S"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.sync_period", "PT30S"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_fin_timeout", "PT30S"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.tcp_timeout", "PT30S"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_kube.cluster", "customization_kube_proxy.0.ipvs.0.udp_timeout", "PT30S"),

					// Check kubeconfig is present and not empty
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig"),

					// Check kubeconfig_attributes are present
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.host"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_certificate"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_kube.cluster", "kubeconfig_attributes.0.client_key"),
				),
			},
		},
	})
}

var testAccCloudProjectKubeDatasourceConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
    name          = "%s"
	region        = "%s"
	version = "%s"
}

data "ovh_cloud_project_kube" "cluster" {
  service_name = ovh_cloud_project_kube.cluster.service_name
  kube_id = ovh_cloud_project_kube.cluster.id
}
`

var testAccCloudProjectKubeDatasourceKubeProxyConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
	
	kube_proxy_mode = "ipvs"
	customization_kube_proxy {
		iptables {
      		min_sync_period = "PT30S"
			sync_period = "PT30S"
    	}
    	
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

data "ovh_cloud_project_kube" "cluster" {
  service_name = ovh_cloud_project_kube.cluster.service_name
  kube_id = ovh_cloud_project_kube.cluster.id
}
`
