package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccPublicCloudKubernetesNodeConfig = fmt.Sprintf(`
data "ovh_cloud_kubernetes_cluster" "cluster" {
	project_id  = "%s"
	name = "%s"
}
resource "ovh_cloud_kubernetes_node" "node" {
	project_id  = "%s"
	cluster_id = data.ovh_cloud_kubernetes_cluster.cluster.id
  	name = "acceptance-node"
	flavor = "b2-7"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"), os.Getenv("OVH_KUBERNETES_CLUSTER_NAME"), os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccPublicCloudKubernetesNode_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicCloudKubernetesNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudKubernetesNodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicCloudKubernetesNodeExists("ovh_cloud_kubernetes_node.node", t),
				),
			},
		},
	})
}

func testAccCheckPublicCloudKubernetesNodeExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("No Project ID is set")
		}

		return cloudKubernetesNodeExists(rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_id"], rs.Primary.ID, config.OVHClient)
	}
}

func testAccCheckPublicCloudKubernetesNodeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_kubernetes_cluster" {
			continue
		}

		err := cloudKubernetesNodeExists(rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_id"], rs.Primary.ID, config.OVHClient)
		if err == nil {
			return fmt.Errorf("cloud > Kubernetes Node still exists")
		}

	}
	return nil
}
