package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccCloudStorageFileShareSnapshotConfig builds the HCL for a vRack network/subnet,
// a share-network backed by them, a share and a snapshot of that share.
func testAccCloudStorageFileShareSnapshotConfig(serviceName, region, vrackNetName, vrackSubnetName, networkName, shareName, snapshotName, description string) string {
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

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = "%s"
  name         = "%s"
  description  = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
}
`, serviceName, networkName, region, serviceName, shareName, region, serviceName, snapshotName, description)
}

func TestAccCloudStorageFileShareSnapshot_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	description := "test file share snapshot"

	config := testAccCloudStorageFileShareSnapshotConfig(
		serviceName, region, vrackNetName, vrackSubnetName,
		networkName, shareName, snapshotName, description,
	)

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
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "name", snapshotName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "description", description),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.share_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.size"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_file_share_snapshot.snapshot", "current_state.location.region", region),
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

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	snapshotNameUpdated := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareSnapshotNamePrefix)
	description := "test file share snapshot"
	descriptionUpdated := "updated file share snapshot"

	config := testAccCloudStorageFileShareSnapshotConfig(
		serviceName, region, vrackNetName, vrackSubnetName,
		networkName, shareName, snapshotName, description,
	)

	configUpdated := testAccCloudStorageFileShareSnapshotConfig(
		serviceName, region, vrackNetName, vrackSubnetName,
		networkName, shareName, snapshotNameUpdated, descriptionUpdated,
	)

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
