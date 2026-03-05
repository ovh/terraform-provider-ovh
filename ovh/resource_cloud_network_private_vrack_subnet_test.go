package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix = "tf-test-subnet-v2-"

func testAccPreCheckCloudNetworkPrivateVrackSubnet(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
}

func testAccCloudNetworkPrivateVrackSubnetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["network_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}

func TestAccCloudNetworkPrivateVrackSubnet_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  region       = "%s"
  dhcp_enabled = true
}
`, serviceName, networkName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrackSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.cidr"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateVrackSubnetImportStateIdFunc("ovh_cloud_network_private_vrack_subnet.test"),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrackSubnet_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  region       = "%s"
  description  = "initial description"
}
`, serviceName, networkName, region, subnetName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  region       = "%s"
  description  = "updated description"
}
`, serviceName, networkName, region, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrackSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "initial description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "updated description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
		},
	})
}
