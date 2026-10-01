package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageFileShareSnapshots_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(test_prefix)
	snapshotName := acctest.RandomWithPrefix(test_prefix)

	config := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
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
  description  = "test file share snapshot"
  share_id     = ovh_cloud_storage_file_share.share.id
}

data "ovh_cloud_storage_file_share_snapshots" "snapshots" {
  service_name = "%s"

  depends_on = [ovh_cloud_storage_file_share_snapshot.snapshot]
}
`, serviceName, networkName, region, serviceName, shareName, region, serviceName, snapshotName, serviceName)

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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_snapshots.snapshots", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_snapshots.snapshots", "share_snapshots.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_snapshots.snapshots", "share_snapshots.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_snapshots.snapshots", "share_snapshots.0.share_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_snapshots.snapshots", "share_snapshots.0.current_state.size"),
				),
			},
		},
	})
}
