package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudStorageBlockVolumeBackupConfig = `
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  region       = "%s"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
`

const testAccCloudStorageBlockVolumeBackupUpdatedConfig = `
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  region       = "%s"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}
`

func TestAccCloudStorageBlockVolumeBackup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	backupName := acctest.RandomWithPrefix(test_prefix)
	description := "test backup description"

	config := fmt.Sprintf(
		testAccCloudStorageBlockVolumeBackupConfig,
		serviceName, volumeName, region,
		serviceName, backupName, description, region,
	)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "name", backupName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "description", description),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "volume_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_backup.backup", "current_state.volume_id"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_block_volume_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf(
						"%s/%s",
						state.RootModule().Resources["ovh_cloud_storage_block_volume_backup.backup"].Primary.Attributes["service_name"],
						state.RootModule().Resources["ovh_cloud_storage_block_volume_backup.backup"].Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func TestAccCloudStorageBlockVolumeBackup_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	backupName := acctest.RandomWithPrefix(test_prefix)
	backupNameUpdated := acctest.RandomWithPrefix(test_prefix)
	description := "test backup description"
	descriptionUpdated := "updated backup description"

	config := fmt.Sprintf(
		testAccCloudStorageBlockVolumeBackupConfig,
		serviceName, volumeName, region,
		serviceName, backupName, description, region,
	)

	configUpdated := fmt.Sprintf(
		testAccCloudStorageBlockVolumeBackupUpdatedConfig,
		serviceName, volumeName, region,
		serviceName, backupNameUpdated, descriptionUpdated, region,
	)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "name", backupName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "description", description),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "name", backupNameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_backup.backup", "description", descriptionUpdated),
				),
			},
		},
	})
}
