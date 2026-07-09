package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageFileShare_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(test_prefix)

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

data "ovh_cloud_storage_file_share" "share" {
  service_name = "%s"
  id           = ovh_cloud_storage_file_share.share.id
}
`, serviceName, networkName, region, serviceName, shareName, region, serviceName)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_file_share.share", "id",
						"ovh_cloud_storage_file_share.share", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "name", shareName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "size", "150"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "share_type", "STANDARD_1AZ"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share.share", "location.region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share.share", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share.share", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share.share", "share_network_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share.share", "current_state.protocol"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share.share", "current_state.share_network_id"),
				),
			},
		},
	})
}
