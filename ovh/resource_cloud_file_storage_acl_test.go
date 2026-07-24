package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCloudFileStorageAclShareConfig(serviceName, region, vrackNetName, vrackSubnetName, networkName, shareName string) string {
	return testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
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
}

func TestAccCloudFileStorageAcl_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	shareConfig := testAccCloudFileStorageAclShareConfig(serviceName, region, vrackNetName, vrackSubnetName, networkName, shareName)

	config := shareConfig + fmt.Sprintf(`
resource "ovh_cloud_file_storage_acl" "acl" {
  service_name = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_ONLY"
}
`, serviceName)

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
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "access_to", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "access_level", "READ_ONLY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "share_id"),
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "current_state.state"),
				),
			},
			{
				ResourceName:      "ovh_cloud_file_storage_acl.acl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudFileStorageAclImportStateIdFunc("ovh_cloud_file_storage_acl.acl"),
			},
		},
	})
}

func TestAccCloudFileStorageAcl_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	shareConfig := testAccCloudFileStorageAclShareConfig(serviceName, region, vrackNetName, vrackSubnetName, networkName, shareName)

	config := shareConfig + fmt.Sprintf(`
resource "ovh_cloud_file_storage_acl" "acl" {
  service_name = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_ONLY"
}
`, serviceName)

	updatedConfig := shareConfig + fmt.Sprintf(`
resource "ovh_cloud_file_storage_acl" "acl" {
  service_name = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_WRITE"
}
`, serviceName)

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
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "access_level", "READ_ONLY"),
				),
			},
			{
				// accessLevel is mutable but the backend deletes+recreates the
				// underlying Manila rule, so the resource ID changes across
				// this update — the framework/provider must handle that.
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_file_storage_acl.acl", "access_level", "READ_WRITE"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_file_storage_acl.acl", "checksum"),
				),
			},
		},
	})
}

func testAccCloudFileStorageAclImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["share_id"], rs.Primary.Attributes["id"]), nil
	}
}
