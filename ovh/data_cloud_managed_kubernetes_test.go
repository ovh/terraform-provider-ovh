package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")

	config := testAccCloudManagedKubernetesConfig(serviceName, region, name, true) + `
data "ovh_cloud_managed_kubernetes" "test" {
  service_name = ovh_cloud_managed_kubernetes.test.service_name
  id           = ovh_cloud_managed_kubernetes.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes.test", "plan", "FREE"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_kubernetes.test", "version", "ovh_cloud_managed_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes.test", "current_state.api_server.access.public_endpoint.url"),
				),
			},
		},
	})
}
