package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStorageBlockVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}
`, serviceName, volumeName, region)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "10"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "resource_status"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "current_state.volume_type"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "current_state.status"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_storage_block_volume.volume",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageBlockVolumeImportStateIdFunc("ovh_cloud_storage_block_volume.volume"),
			},
		},
	})
}

func TestAccCloudStorageBlockVolume_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	volumeName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudStorageBlockVolumeNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}
`, serviceName, volumeName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 20
  region       = "%s"
  volume_type  = "CLASSIC"
}
`, serviceName, updatedName, region)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "10"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume.volume", "size", "20"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume.volume", "checksum"),
				),
			},
		},
	})
}

const testAccResourceCloudStorageBlockVolumeNamePrefix = "tf-test-volume-v2-"

func testAccCloudStorageBlockVolumeImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
