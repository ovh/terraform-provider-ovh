package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_vrack_network", &resource.Sweeper{
		Name: "ovh_iploadbalancing_vrack_network",
		F:    testSweepIpLoadbalancingVrackNetwork,
	})
}

const (
	ipLoadbalancingVrackNetworkVlan   = "1001"
	ipLoadbalancingVrackNetworkSubnet = "10.0.1.0/24"
	ipLoadbalancingVrackNetworkNatIp  = "10.0.1.0/27"
)

func testSweepIpLoadbalancingVrackNetwork(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	if serviceName == "" {
		return fmt.Errorf("OVH_IPLB_SERVICE env var is required")
	}

	endpoint := fmt.Sprintf(
		"/ipLoadbalancing/%s/vrack/network?subnet=%s&vlan=%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipLoadbalancingVrackNetworkSubnet),
		url.PathEscape(ipLoadbalancingVrackNetworkVlan),
	)

	result := make([]int64, 0)

	if err := client.Get(endpoint, result); err != nil {
		if err.(*ovh.APIError).Code == 404 {
			return nil
		}
		return err
	}

	for _, id := range result {
		// delete farms, then delete vrack network
		endpoint = fmt.Sprintf(
			"/ipLoadbalancing/%s/tcp/farm?vrackNetworkId=%d",
			url.PathEscape(serviceName),
			id,
		)

		farms := make([]int64, 0)
		if err := client.Get(endpoint, farms); err != nil && !(err.(*ovh.APIError).Code == 404) {
			return err
		}
		for _, farmId := range farms {
			endpoint = fmt.Sprintf(
				"/ipLoadbalancing/%s/tcp/farm/%d",
				url.PathEscape(serviceName),
				farmId,
			)
			// delete the farm
			if err := client.Delete(endpoint, nil); err != nil {
				return fmt.Errorf("Error calling DELETE %s:\n\t %q", endpoint, err)
			}
		}

		// delete the vrack network
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/vrack/network/%d",
			url.PathEscape(serviceName),
			id,
		)
		if err := client.Delete(endpoint, nil); err != nil {
			return fmt.Errorf("Error calling DELETE %s:\n\t %q", endpoint, err)
		}
	}

	return nil
}

func TestAccIpLoadbalancingVrackNetwork_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackIpLoadbalancingPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingVrackNetworkConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "subnet", ipLoadbalancingVrackNetworkSubnet),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "vlan", ipLoadbalancingVrackNetworkVlan),
					resource.TestCheckResourceAttrSet("ovh_iploadbalancing_vrack_network.network", "id"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "nat_ip", ipLoadbalancingVrackNetworkNatIp),
				),
			},
		},
	})
}

var testAccIpLoadbalancingVrackNetworkConfig_basic = fmt.Sprintf(`
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

resource "ovh_iploadbalancing_tcp_farm" "testfarm" {
  service_name     = data.ovh_iploadbalancing.iplb.service_name
  display_name     = "terraform_testacc"
  port             = 80
  vrack_network_id = ovh_iploadbalancing_vrack_network.network.vrack_network_id
  zone             = tolist(data.ovh_iploadbalancing.iplb.zone)[0]
}
`,
	os.Getenv("OVH_IPLB_SERVICE"),
	os.Getenv("OVH_VRACK"),
	ipLoadbalancingVrackNetworkSubnet,
	ipLoadbalancingVrackNetworkVlan,
	ipLoadbalancingVrackNetworkNatIp,
)
