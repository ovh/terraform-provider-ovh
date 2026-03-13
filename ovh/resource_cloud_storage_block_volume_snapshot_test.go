package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudStorageBlockVolumeSnapshotConfig = `
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
`

const testAccCloudStorageBlockVolumeSnapshotUpdatedConfig = `
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
`

func TestAccCloudStorageBlockVolumeSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)
	description := "test snapshot description"

	config := fmt.Sprintf(
		testAccCloudStorageBlockVolumeSnapshotConfig,
		serviceName, volumeName, region,
		serviceName, snapshotName, description, region,
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "description", description),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "volume_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_block_volume_snapshot.snapshot", "current_state.volume_id"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_block_volume_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf(
						"%s/%s",
						state.RootModule().Resources["ovh_cloud_storage_block_volume_snapshot.snapshot"].Primary.Attributes["service_name"],
						state.RootModule().Resources["ovh_cloud_storage_block_volume_snapshot.snapshot"].Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func TestAccCloudStorageBlockVolumeSnapshot_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)
	snapshotNameUpdated := acctest.RandomWithPrefix(test_prefix)
	description := "test snapshot description"
	descriptionUpdated := "updated snapshot description"

	config := fmt.Sprintf(
		testAccCloudStorageBlockVolumeSnapshotConfig,
		serviceName, volumeName, region,
		serviceName, snapshotName, description, region,
	)

	configUpdated := fmt.Sprintf(
		testAccCloudStorageBlockVolumeSnapshotUpdatedConfig,
		serviceName, volumeName, region,
		serviceName, snapshotNameUpdated, descriptionUpdated, region,
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "description", description),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "name", snapshotNameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_storage_block_volume_snapshot.snapshot", "description", descriptionUpdated),
				),
			},
		},
	})
}
