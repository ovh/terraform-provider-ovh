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
	networkId := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_SUBNET_ID_TEST")

	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "%s"
  name         = "%s"
  size         = 100
  region       = "%s"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "%s"
  subnet_id    = "%s"
}
`, serviceName, shareName, region, networkId, subnetId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "name", shareName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "size", "100"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "share_type", "STANDARD_1AZ"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "network_id", networkId),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "subnet_id", subnetId),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "resource_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share.share", "current_state.protocol"),
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
	networkId := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_SUBNET_ID_TEST")

	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "%s"
  name         = "%s"
  size         = 100
  region       = "%s"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "%s"
  subnet_id    = "%s"
}
`, serviceName, shareName, region, networkId, subnetId)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_file_share" "share" {
  service_name = "%s"
  name         = "%s"
  size         = 200
  region       = "%s"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "%s"
  subnet_id    = "%s"
  description  = "updated description"

  access_rules {
    access_to    = "10.0.0.0/24"
    access_level = "READ_WRITE"
  }
}
`, serviceName, updatedName, region, networkId, subnetId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "name", shareName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share.share", "size", "100"),
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
