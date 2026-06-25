package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudNetworkPrivateVrackSubnets_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)
	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "%s"
  location = {
	region = "%s"
  }
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  dhcp_enabled = true
  location = {
    region = "%s"
  }
}

data "ovh_cloud_network_private_vrack_subnets" "test" {
  service_name = ovh_cloud_network_private_vrack_subnet.test.service_name
  network_id   = ovh_cloud_network_private_vrack_subnet.test.network_id
}
`, serviceName, networkName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnets.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnets.test", "network_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnets.test", "subnets.#"),
				),
			},
		},
	})
}
