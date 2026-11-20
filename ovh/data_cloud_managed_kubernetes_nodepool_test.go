package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesNodepoolDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")

	config := testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a") + `
data "ovh_cloud_managed_kubernetes_nodepool" "test" {
  service_name          = ovh_cloud_managed_kubernetes_nodepool.test.service_name
  managed_kubernetes_id = ovh_cloud_managed_kubernetes_nodepool.test.managed_kubernetes_id
  id                    = ovh_cloud_managed_kubernetes_nodepool.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes_nodepool.test", "name", name),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes_nodepool.test", "desired_nodes", "1"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes_nodepool.test", "node_config.flavor", "b3-8"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_kubernetes_nodepool.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_nodepool.test", "current_state.status.current_nodes"),
				),
			},
		},
	})
}
