package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "subnet"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "vlan"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "nat_ip"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_network.network", "farm_id.#"),
				),
			},
		},
	})
}

var testAccIpLoadbalancingVrackNetworkDatasourceConfig_basic = fmt.Sprintf(`
data ovh_iploadbalancing "iplb" {
  service_name = "%s"
}

data ovh_iploadbalancing_vrack_networks "networks" {
  service_name = data.ovh_iploadbalancing.iplb.service_name
}

data ovh_iploadbalancing_vrack_network "network" {
  service_name = data.ovh_iploadbalancing.iplb.service_name
  vrack_network_id  = data.ovh_iploadbalancing_vrack_networks.networks.result[0]
}
`, os.Getenv("OVH_IPLB_SERVICE"))
