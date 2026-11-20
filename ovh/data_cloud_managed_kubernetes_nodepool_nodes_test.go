package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudManagedKubernetesNodepoolNodesDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")

	config := testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a") + `
data "ovh_cloud_managed_kubernetes_nodepool_nodes" "test" {
  service_name          = ovh_cloud_managed_kubernetes_nodepool.test.service_name
  managed_kubernetes_id = ovh_cloud_managed_kubernetes_nodepool.test.managed_kubernetes_id
  nodepool_id           = ovh_cloud_managed_kubernetes_nodepool.test.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_kubernetes_nodepool_nodes.test", "nodes.#", "ovh_cloud_managed_kubernetes_nodepool.test", "current_state.status.current_nodes"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_kubernetes_nodepool_nodes.test", "nodes.0.id"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_kubernetes_nodepool_nodes.test", "nodes.0.current_state.flavor", "ovh_cloud_managed_kubernetes_nodepool.test", "node_config.flavor"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_kubernetes_nodepool_nodes.test", "nodes.0.current_state.nodepool_id", "ovh_cloud_managed_kubernetes_nodepool.test", "id"),
				),
			},
		},
	})
}
