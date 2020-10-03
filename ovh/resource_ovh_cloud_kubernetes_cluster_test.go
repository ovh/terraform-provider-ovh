package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccPublicCloudKubernetesClusterConfig = fmt.Sprintf(`
resource "ovh_cloud_kubernetes_cluster" "cluster" {
	project_id  = "%s"
  	name = "my cluster for acceptance tests"
	region = "GRA5"
	version = "1.15"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccPublicCloudKubernetesCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicCloudKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudKubernetesClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicCloudKubernetesClusterExists("ovh_cloud_kubernetes_cluster.cluster", t),
				),
			},
		},
	})
}

func testAccCheckPublicCloudKubernetesClusterExists(n string, t *testing.T) resource.TestCheckFunc {
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

		return cloudKubernetesClusterExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
	}
}

func testAccCheckPublicCloudKubernetesClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_kubernetes_cluster" {
			continue
		}

		err := cloudKubernetesClusterExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
		if err == nil {
			return fmt.Errorf("cloud > Kubernetes Cluster still exists")
		}

	}
	return nil
}
