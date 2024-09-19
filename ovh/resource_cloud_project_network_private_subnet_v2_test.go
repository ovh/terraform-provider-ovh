package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccCloudProjectNetworkPrivateSubnetV2Config_basic = `
%s

resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name = ovh_cloud_project_network_private.network.service_name

  # whatever region, for test purpose
  network_id        = element(tolist(ovh_cloud_project_network_private.network.regions_attributes), 0).openstackid
  region            = element(tolist(ovh_cloud_project_network_private.network.regions), 0)
  name              = "my_new_subnet"
  cidr              = "192.168.169.0/24"
  dns_nameservers   = ["1.1.1.1"]
  host_routes     = [
    {
      destination = "192.168.169.0/24", 
      nexthop = "192.168.169.254"
    }
  ]
  allocation_pools   = [
    {
      start = "192.168.169.100", 
      end = "192.168.169.200"
    }
  ]
  dhcp              = true
  gateway_ip        = "192.168.169.253"
  enable_gateway_ip = true
}
`

func TestAccCloudProjectNetworkPrivateSubnetV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckcCloudProjectNetworkPrivateSubnetPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateSubnetConfig(testAccCloudProjectNetworkPrivateSubnetV2Config_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "network_id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_network_private_subnet_v2.subnet", "region"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "gateway_ip", "192.168.169.253"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pools.0.start", "192.168.169.100"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "allocation_pools.0.end", "192.168.169.200"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_routes.0.destination", "192.168.169.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "host_routes.0.nexthop", "192.168.169.254"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "dhcp", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "enable_gateway_ip", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_network_private_subnet_v2.subnet", "cidr", "192.168.169.0/24"),
				),
			},
		},
	})
}
