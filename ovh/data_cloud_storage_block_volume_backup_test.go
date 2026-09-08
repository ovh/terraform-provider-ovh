package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudStorageBlockVolumeBackupConfig = `
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
  %s
}

data "ovh_cloud_storage_block_volume_backup" "backup" {
  service_name = "%s"
  id           = ovh_cloud_storage_block_volume_backup.backup.id
}
`

func TestAccDataSourceCloudStorageBlockVolumeBackup_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	backupName := acctest.RandomWithPrefix(test_prefix)
	description := "test backup description"
	// Optional: only set for 3AZ regions, where the API reports an availability zone.
	az := os.Getenv("OVH_CLOUD_PROJECT_AZ_TEST")

	azConfig := ""
	if az != "" {
		azConfig = fmt.Sprintf("availability_zone = %q", az)
	}

	config := fmt.Sprintf(
		testAccDataSourceCloudStorageBlockVolumeBackupConfig,
		serviceName, volumeName, region,
		serviceName, backupName, description, region, azConfig,
		serviceName,
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backup.backup", "service_name", serviceName),
		resource.TestCheckResourceAttrPair(
			"data.ovh_cloud_storage_block_volume_backup.backup", "id",
			"ovh_cloud_storage_block_volume_backup.backup", "id",
		),
		resource.TestCheckResourceAttrPair(
			"data.ovh_cloud_storage_block_volume_backup.backup", "volume_id",
			"ovh_cloud_storage_block_volume.volume", "id",
		),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backup.backup", "name", backupName),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backup.backup", "description", description),
		resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volume_backup.backup", "size"),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backup.backup", "location.region", region),
	}
	if az != "" {
		checks = append(checks, resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backup.backup", "location.availability_zone", az))
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
