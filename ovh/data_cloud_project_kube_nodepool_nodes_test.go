package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectKubeNodePoolNodesDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectKubeNodePoolNodesDataSourceConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.#"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.created_at"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.instance_id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.is_up_to_date"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.node_pool_id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.status"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.updated_at"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.version"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.flavor", "b2-7"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube_nodepool_nodes.nodePoolNodesDataSource", "nodes.0.project_id", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				),
			},
		},
	})
}

var testAccCloudProjectKubeNodePoolNodesDataSourceConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}

resource "ovh_cloud_project_kube_nodepool" "pool" {
	service_name  = ovh_cloud_project_kube.cluster.service_name
	kube_id       = ovh_cloud_project_kube.cluster.id
	name          = ovh_cloud_project_kube.cluster.name
	flavor_name   = "b2-7"
	desired_nodes = 1
	min_nodes     = 0
	max_nodes     = 2

	depends_on = [
		ovh_cloud_project_kube.cluster
	]
}

data "ovh_cloud_project_kube_nodepool_nodes" "nodePoolNodesDataSource" {
  service_name  = ovh_cloud_project_kube.cluster.service_name
  kube_id       = ovh_cloud_project_kube.cluster.id
  name       	= ovh_cloud_project_kube_nodepool.pool.name

  depends_on = [
    ovh_cloud_project_kube_nodepool.pool
  ]
}
`
