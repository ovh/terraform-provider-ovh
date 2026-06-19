package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudNetworkPrivateVrackSubnet_basic(t *testing.T) {
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
  project_id   = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id
  name         = "%s"
  cidr         = "10.0.0.0/24"
  dhcp_enabled = true
  location = {
    region = "%s"
  }
}

data "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack_subnet.test.project_id
  network_id   = ovh_cloud_network_private_vrack_subnet.test.network_id
  id           = ovh_cloud_network_private_vrack_subnet.test.id
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
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "location.region", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "network_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "current_state.cidr"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack_subnet.test", "current_state.location.region"),
				),
			},
		},
	})
}
