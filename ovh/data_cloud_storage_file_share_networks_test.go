package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageFileShareNetworks_basic(t *testing.T) {
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

data "ovh_cloud_storage_file_share_networks" "networks" {
  service_name = "%s"
  region       = "%s"

  depends_on = [ovh_cloud_storage_file_share_network.network]
}
`, serviceName, networkName, region, serviceName, region)

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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_networks.networks", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_networks.networks", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_networks.networks", "share_networks.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_networks.networks", "share_networks.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_networks.networks", "share_networks.0.network_id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_networks.networks", "share_networks.0.location.region", region),
				),
			},
		},
	})
}
