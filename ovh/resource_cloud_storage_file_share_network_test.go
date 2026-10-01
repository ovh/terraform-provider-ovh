package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceCloudStorageFileShareVrackSubnetNamePrefix = "tf-test-vracksubnet-v2-"

// testAccVrackNetworkSubnetConfig returns the HCL for a vRack private network and
// subnet used as the backing network for the file storage share-network tests.
// The share-network references:
//   - network_id = ovh_cloud_network_private_vrack.vrack_net.id
//   - subnet_id  = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
func testAccVrackNetworkSubnetConfig(serviceName, region, netName, subnetName string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "vrack_net" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "vrack_subnet" {
  service_name = ovh_cloud_network_private_vrack.vrack_net.service_name
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  region       = ovh_cloud_network_private_vrack.vrack_net.region
}
`, serviceName, netName, region, subnetName)
}

func TestAccCloudStorageFileShareNetwork_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)

	config := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  name         = "%s"
  description  = "test share network"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}
`, serviceName, networkName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckVRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_network.network", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_network.network", "name", networkName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_network.network", "description", "test share network"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_file_share_network.network", "network_id", "ovh_cloud_network_private_vrack.vrack_net", "id"),
					resource.TestCheckResourceAttrPair("ovh_cloud_storage_file_share_network.network", "subnet_id", "ovh_cloud_network_private_vrack_subnet.vrack_subnet", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_network.network", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "resource_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "current_state.network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_network.network", "current_state.subnet_id"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_file_share_network.network",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageFileShareNetworkImportStateIdFunc("ovh_cloud_storage_file_share_network.network"),
			},
		},
	})
}

const testAccResourceCloudStorageFileShareNetworkNamePrefix = "tf-test-sharenet-v2-"

func testAccCloudStorageFileShareNetworkImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
