package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudStorageBlockVolumeBackupsConfig = `
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

data "ovh_cloud_storage_block_volume_backups" "backups" {
  service_name = "%s"
  region       = "%s"
  volume_id    = ovh_cloud_storage_block_volume.volume.id

  depends_on = [ovh_cloud_storage_block_volume_backup.backup]
}
`

func TestAccDataSourceCloudStorageBlockVolumeBackups_basic(t *testing.T) {
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
		testAccDataSourceCloudStorageBlockVolumeBackupsConfig,
		serviceName, volumeName, region,
		serviceName, backupName, description, region, azConfig,
		serviceName, region,
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backups.backups", "service_name", serviceName),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backups.backups", "region", region),
		resource.TestCheckResourceAttrPair(
			"data.ovh_cloud_storage_block_volume_backups.backups", "volume_id",
			"ovh_cloud_storage_block_volume.volume", "id",
		),
		resource.TestCheckResourceAttrPair(
			"data.ovh_cloud_storage_block_volume_backups.backups", "backups.0.id",
			"ovh_cloud_storage_block_volume_backup.backup", "id",
		),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backups.backups", "backups.0.name", backupName),
		resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volume_backups.backups", "backups.0.size"),
		resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backups.backups", "backups.0.location.region", region),
	}
	if az != "" {
		checks = append(checks, resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_backups.backups", "backups.0.location.availability_zone", az))
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
