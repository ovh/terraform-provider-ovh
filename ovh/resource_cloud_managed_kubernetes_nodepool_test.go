package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name string, desiredNodes int, label string) string {
	return testAccCloudManagedKubernetesConfig(serviceName, region, name, true) + fmt.Sprintf(`
resource "ovh_cloud_managed_kubernetes_nodepool" "test" {
  service_name          = ovh_cloud_managed_kubernetes.test.service_name
  managed_kubernetes_id = ovh_cloud_managed_kubernetes.test.id
  name                  = "%s"
  desired_nodes         = %d

  node_config = { flavor = "b3-8" }

  template = {
    metadata = { labels = { team = "%s" } }
  }
}
`, name, desiredNodes, label)
}

func testAccCloudManagedKubernetesNodepoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		a := rs.Primary.Attributes
		return a["service_name"] + "/" + a["managed_kubernetes_id"] + "/" + a["id"], nil
	}
}

func TestAccCloudManagedKubernetesNodepool_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")
	importIdFunc := testAccCloudManagedKubernetesNodepoolImportStateIdFunc("ovh_cloud_managed_kubernetes_nodepool.test")
	var importId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a"),
				Check: resource.ComposeTestCheckFunc(
					// Saved for the import step, which runs after the nodepool left the state.
					func(s *terraform.State) (err error) {
						importId, err = importIdFunc(s)
						return err
					},
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "desired_nodes", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "node_config.flavor", "b3-8"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "template.metadata.labels.team", "a"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrPair("ovh_cloud_managed_kubernetes_nodepool.test", "managed_kubernetes_id", "ovh_cloud_managed_kubernetes.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes_nodepool.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes_nodepool.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes_nodepool.test", "billing_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes_nodepool.test", "current_state.status.current_nodes"),
				),
			},
			{
				// Forgets the nodepool without destroying it, so the next step can import it into the state.
				Config: testAccCloudManagedKubernetesConfig(serviceName, region, name, true) + `
removed {
  from = ovh_cloud_managed_kubernetes_nodepool.test
  lifecycle { destroy = false }
}
`,
			},
			{
				Config:             testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a"),
				ResourceName:       "ovh_cloud_managed_kubernetes_nodepool.test",
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc:  func(*terraform.State) (string, error) { return importId, nil },
			},
			{
				// Apply no-op after the import: records the nodepool dependencies in the
				// state so the post-test destroy deletes it before the cluster.
				Config: testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a"),
			},
		},
	})
}

func TestAccCloudManagedKubernetesNodepool_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 1, "a"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "desired_nodes", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "template.metadata.labels.team", "a"),
				),
			},
			{
				Config: testAccCloudManagedKubernetesNodepoolConfig(serviceName, region, name, 2, "b"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "desired_nodes", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "template.metadata.labels.team", "b"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes_nodepool.test", "resource_status", "READY"),
				),
			},
		},
	})
}
