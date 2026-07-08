package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageFileShareNetwork_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(test_prefix)

	config := testAccVrackNetworkSubnetConfig(serviceName, region, vrackNetName, vrackSubnetName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  name         = "%s"
  description  = "test share network"
  network_id   = ovh_cloud_network_private_vrack.vrack_net.id
  subnet_id    = ovh_cloud_network_private_vrack_subnet.vrack_subnet.id
  region       = "%s"
}

data "ovh_cloud_storage_file_share_network" "network" {
  service_name = "%s"
  id           = ovh_cloud_storage_file_share_network.network.id
}
`, serviceName, networkName, region, serviceName)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_network.network", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_file_share_network.network", "id",
						"ovh_cloud_storage_file_share_network.network", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_network.network", "name", networkName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_file_share_network.network", "network_id",
						"ovh_cloud_network_private_vrack.vrack_net", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_file_share_network.network", "subnet_id",
						"ovh_cloud_network_private_vrack_subnet.vrack_subnet", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_network.network", "location.region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_network.network", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_network.network", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_network.network", "current_state.network_id"),
				),
			},
		},
	})
}
