package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStorageFileShare_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	config := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "%s"
  name             = "%s"
  size             = 150
  region           = "%s"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.network.id
}
`, serviceName, networkName, region, serviceName, shareName, region)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "name", shareName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "size", "150"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "share_type", "STANDARD_1AZ"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "resource_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "share_network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "current_state.protocol"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "current_state.share_network_id"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_file_share.share",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageFileShareImportStateIdFunc("ovh_cloud_storage_file_share.share"),
			},
		},
	})
}

func TestAccCloudStorageFileShare_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	vrackConfig := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName)

	config := vrackConfig + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "%s"
  name             = "%s"
  size             = 150
  region           = "%s"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.network.id
}
`, serviceName, networkName, region, serviceName, shareName, region)

	updatedConfig := vrackConfig + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  name         = "%s"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "%s"
  name             = "%s"
  size             = 200
  region           = "%s"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.network.id
  description      = "updated description"

  access_rules = [
	{
		access_to    = "10.0.0.0/24"
		access_level = "READ_WRITE"
	}
  ]
}
`, serviceName, networkName, region, serviceName, updatedName, region)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "name", shareName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "size", "150"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "size", "200"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "access_rules.0.access_to", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "access_rules.0.access_level", "READ_WRITE"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "checksum"),
				),
			},
		},
	})
}

const testAccResourceCloudStorageFileShareNamePrefix = "tf-test-fileshare-v2-"

func testAccCloudStorageFileShareImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
