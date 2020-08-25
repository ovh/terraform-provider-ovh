package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIpLoadbalancingVrackNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingVrackNetworkDatasourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing_vrack_network.network", "subnet", testAccIpLoadbalancingVrackNetworkSubnet),
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing_vrack_network.network", "vlan", testAccIpLoadbalancingVrackNetworkVlan1001),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "id"),
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing_vrack_network.network", "nat_ip", testAccIpLoadbalancingVrackNetworkNatIp),
				),
			},
		},
	})
}

var testAccIpLoadbalancingVrackNetworkDatasourceConfig_basic = fmt.Sprintf(`
data ovh_iploadbalancing "iplb" {
  service_name = "%s"
}

resource "ovh_vrack_iploadbalancing" "viplb" {
  service_name     = "%s"
  ip_loadbalancing = data.ovh_iploadbalancing.iplb.service_name
}

resource ovh_iploadbalancing_vrack_network "network" {
  service_name = ovh_vrack_iploadbalancing.viplb.ip_loadbalancing
  subnet       = "%s"
  vlan         = %s
  nat_ip       = "%s"
  display_name = "terraform_testacc"
}

data ovh_iploadbalancing_vrack_network "network" {
  service_name = data.ovh_iploadbalancing.iplb.service_name
  vrack_network_id  = ovh_iploadbalancing_vrack_network.network.vrack_network_id
}
`,
	os.Getenv("OVH_IPLB_SERVICE"),
	os.Getenv("OVH_VRACK"),
	testAccIpLoadbalancingVrackNetworkSubnet,
	testAccIpLoadbalancingVrackNetworkVlan1001,
	testAccIpLoadbalancingVrackNetworkNatIp,
)
