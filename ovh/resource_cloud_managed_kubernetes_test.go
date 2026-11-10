package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudManagedKubernetesConfig(serviceName, region, name string, autoUpgrade bool) string {
	return fmt.Sprintf(`
resource "ovh_cloud_managed_kubernetes" "test" {
  service_name = "%s"
  region       = "%s"
  plan         = "FREE"
  name         = "%s"

  cni     = { type = "CANAL" }
  network = { type = "PUBLIC" }

  auto_upgrade = {
    patch = { enabled = %t }
  }
}
`, serviceName, region, name, autoUpgrade)
}

func testAccCloudManagedKubernetesImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["service_name"] + "/" + rs.Primary.Attributes["id"], nil
	}
}

func TestAccCloudManagedKubernetes_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")
	importIdFunc := testAccCloudManagedKubernetesImportStateIdFunc("ovh_cloud_managed_kubernetes.test")
	var importId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedKubernetesConfig(serviceName, region, name, true),
				Check: resource.ComposeTestCheckFunc(
					// Saved for the import step, which runs after the cluster left the state.
					func(s *terraform.State) (err error) {
						importId, err = importIdFunc(s)
						return err
					},
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "plan", "FREE"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "cni.type", "CANAL"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "network.type", "PUBLIC"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes.test", "current_state.deployed_version"),
					resource.TestCheckResourceAttrSet("ovh_cloud_managed_kubernetes.test", "current_state.api_server.access.public_endpoint.url"),
				),
			},
			{
				// Forgets the cluster without destroying it, so the next step can import it into the state.
				Config: `
removed {
  from = ovh_cloud_managed_kubernetes.test
  lifecycle { destroy = false }
}
`,
			},
			{
				Config:             testAccCloudManagedKubernetesConfig(serviceName, region, name, true),
				ResourceName:       "ovh_cloud_managed_kubernetes.test",
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc:  func(*terraform.State) (string, error) { return importId, nil },
			},
			{
				// The imported state must plan no change against the original config.
				Config:   testAccCloudManagedKubernetesConfig(serviceName, region, name, true),
				PlanOnly: true,
			},
		},
	})
}

// testAccCloudManagedKubernetesConfigRestricted adds an API server allowlist to the base config.
func testAccCloudManagedKubernetesConfigRestricted(serviceName, region, name string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_managed_kubernetes" "test" {
  service_name = "%s"
  region       = "%s"
  plan         = "FREE"
  name         = "%s"

  cni     = { type = "CANAL" }
  network = { type = "PUBLIC" }

  api_server = {
    access = { authorized_ips = { rules = [{ cidr = "203.0.113.0/24" }] } }
  }
}
`, serviceName, region, name)
}

// auto_upgrade.patch.enabled is not covered: the API forces it back to true on the FREE plan.
func TestAccCloudManagedKubernetes_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	name := acctest.RandomWithPrefix("test-mks")
	nameUpdated := name + "-upd"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedKubernetes(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedKubernetesConfig(serviceName, region, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "api_server.access.authorized_ips.rules.0.cidr", "0.0.0.0/0"),
				),
			},
			{
				Config: testAccCloudManagedKubernetesConfigRestricted(serviceName, region, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "name", nameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "api_server.access.authorized_ips.rules.0.cidr", "203.0.113.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_managed_kubernetes.test", "resource_status", "READY"),
				),
			},
		},
	})
}
