package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudStorageBlockVolumeSnapshotsConfig = `
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume_snapshot" "snapshot" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  region       = "%s"
  volume_id    = ovh_cloud_storage_block_volume.volume.id
}

data "ovh_cloud_storage_block_volume_snapshots" "snapshots" {
  service_name = "%s"
  region       = "%s"
  volume_id    = ovh_cloud_storage_block_volume.volume.id

  depends_on = [ovh_cloud_storage_block_volume_snapshot.snapshot]
}
`

func TestAccDataSourceCloudStorageBlockVolumeSnapshots_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)
	description := "test snapshot description"

	config := fmt.Sprintf(
		testAccDataSourceCloudStorageBlockVolumeSnapshotsConfig,
		serviceName, volumeName, region,
		serviceName, snapshotName, description, region,
		serviceName, region,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_snapshots.snapshots", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_snapshots.snapshots", "region", region),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_block_volume_snapshots.snapshots", "volume_id",
						"ovh_cloud_storage_block_volume.volume", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_block_volume_snapshots.snapshots", "snapshots.0.id",
						"ovh_cloud_storage_block_volume_snapshot.snapshot", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_snapshots.snapshots", "snapshots.0.name", snapshotName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volume_snapshots.snapshots", "snapshots.0.size"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume_snapshots.snapshots", "snapshots.0.location.region", region),
				),
			},
		},
	})
}
