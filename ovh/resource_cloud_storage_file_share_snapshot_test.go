package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCloudStorageFileShareSnapshotConfig = `
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

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  region       = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
}
`

const testAccCloudStorageFileShareSnapshotUpdatedConfig = `
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

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  region       = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
}
`

func TestAccCloudStorageFileShareSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	networkId := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_SUBNET_ID_TEST")
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	description := "test file share snapshot"

	config := fmt.Sprintf(
		testAccCloudStorageFileShareSnapshotConfig,
		serviceName, shareName, region, networkId, subnetId,
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "description", description),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.share_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.share_proto"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_file_share_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf(
						"%s/%s",
						state.RootModule().Resources["ovh_cloud_storage_file_share_snapshot.snapshot"].Primary.Attributes["service_name"],
						state.RootModule().Resources["ovh_cloud_storage_file_share_snapshot.snapshot"].Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func TestAccCloudStorageFileShareSnapshot_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	networkId := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_SUBNET_ID_TEST")
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotNameUpdated := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	description := "test file share snapshot"
	descriptionUpdated := "updated file share snapshot"

	config := fmt.Sprintf(
		testAccCloudStorageFileShareSnapshotConfig,
		serviceName, shareName, region, networkId, subnetId,
		serviceName, snapshotName, description, region,
	)

	configUpdated := fmt.Sprintf(
		testAccCloudStorageFileShareSnapshotUpdatedConfig,
		serviceName, shareName, region, networkId, subnetId,
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
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "description", description),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "name", snapshotNameUpdated),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "description", descriptionUpdated),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "checksum"),
				),
			},
		},
	})
}

const testAccResourceCloudStorageFileShareSnapshotNamePrefix = "tf-test-fileshare-snap-v2-"
