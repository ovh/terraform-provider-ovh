package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIpLoadbalancingVrackNetworksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingVrackNetworksDatasourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_networks.networks", "result.#"),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingVrackNetworksDataSource_withFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingVrackNetworksDatasourceConfig_withFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_vrack_networks.networks", "result.#"),
				),
			},
		},
	})
}

var testAccIpLoadbalancingVrackNetworksDatasourceConfig_basic = fmt.Sprintf(`
data ovh_iploadbalancing "iplb" {
  service_name = "%s"
}

data ovh_iploadbalancing_vrack_networks "networks" {
  service_name = data.ovh_iploadbalancing.iplb.service_name
  subnet = "10.0.0.0/24"
}
`, os.Getenv("OVH_IPLB_SERVICE"))

var testAccIpLoadbalancingVrackNetworksDatasourceConfig_withFilters = fmt.Sprintf(`
data ovh_iploadbalancing "iplb" {
  service_name = "%s"
}

data ovh_iploadbalancing_vrack_networks "networks" {
  service_name = data.ovh_iploadbalancing.iplb.service_name
  subnet       = "10.0.0.0/24"
  vlan_id      = 0
}
`, os.Getenv("OVH_IPLB_SERVICE"))
