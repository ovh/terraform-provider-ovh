package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudNetworkPrivateVrackNamePrefix = "tf-test-net-v2-"

func testAccPreCheckCloudNetworkPrivateVrack(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudNetworkPrivateVrackImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudNetworkPrivateVrack_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, networkName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateVrackImportStateIdFunc("ovh_cloud_network_private_vrack.test"),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrack_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, networkName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "name", updatedName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
				),
			},
		},
	})
}
