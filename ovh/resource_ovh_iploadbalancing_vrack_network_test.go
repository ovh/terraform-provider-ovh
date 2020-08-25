package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_vrack_network", &resource.Sweeper{
		Name: "ovh_iploadbalancing_vrack_network",
		F:    testSweepIpLoadbalancingVrackNetwork,
	})
}

const (
	testAccIpLoadbalancingVrackNetworkVlan1001 = "1001"
	testAccIpLoadbalancingVrackNetworkVlan1002 = "1002"
	testAccIpLoadbalancingVrackNetworkSubnet   = "10.0.1.0/24"
	testAccIpLoadbalancingVrackNetworkNatIp    = "10.0.1.0/27"
	testAccIpLoadbalancingVrackNetworkConfig   = `
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
`
)

func testSweepIpLoadbalancingVrackNetwork(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE is not set. No iploadbalancing_vrack_network to sweep")
		return nil
	}

	log.Print("[DEBUG] Sweeping iploadbalancing_vrack_network")

	get_network_ids := func(vlanId string) ([]int64, error) {
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/vrack/network?vlan=%s",
			url.PathEscape(serviceName),
			url.PathEscape(vlanId),
		)

		result := make([]int64, 0)

		if err := client.Get(endpoint, &result); err != nil {
			if err.(*ovh.APIError).Code == 404 {
				return nil, nil
			}
			return nil, err
		}

		return result, nil
	}

	delete_farms := func(farmType string, networkId int64) error {
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/%s/farm?vrackNetworkId=%d",
			url.PathEscape(serviceName),
			url.PathEscape(farmType),
			networkId,
		)

		result := make([]int64, 0)

		if err := client.Get(endpoint, &result); err != nil {
			if err.(*ovh.APIError).Code == 404 {
				return nil
			}
			return err
		}

		for _, farmId := range result {
			endpoint := fmt.Sprintf(
				"/ipLoadbalancing/%s/%s/farm/%d",
				url.PathEscape(serviceName),
				url.PathEscape(farmType),
				farmId,
			)
			// delete the farm
			log.Printf("[DEBUG] Calling DELETE on %v", endpoint)
			if err := client.Delete(endpoint, nil); err != nil {
				return fmt.Errorf("Error calling DELETE %s:\n\t %q", endpoint, err)
			}
		}
		return nil
	}

	resultVlan1001, err := get_network_ids(testAccIpLoadbalancingVrackNetworkVlan1001)
	if err != nil {
		return err
	}

	resultVlan1002, err := get_network_ids(testAccIpLoadbalancingVrackNetworkVlan1002)
	if err != nil {
		return err
	}

	result := append(resultVlan1001, resultVlan1002...)
	for _, id := range result {
		// delete farms, then delete vrack network
		if err := delete_farms("http", id); err != nil {
			return err
		}

		if err := delete_farms("tcp", id); err != nil {
			return err
		}

		// delete the vrack network
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/vrack/network/%d",
			url.PathEscape(serviceName),
			id,
		)

		log.Printf("[DEBUG] Calling DELETE on %v", endpoint)
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
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "subnet", testAccIpLoadbalancingVrackNetworkSubnet),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "vlan", testAccIpLoadbalancingVrackNetworkVlan1001),
					resource.TestCheckResourceAttrSet("ovh_iploadbalancing_vrack_network.network", "id"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "nat_ip", testAccIpLoadbalancingVrackNetworkNatIp),
				),
			},
			{
				Config: testAccIpLoadbalancingVrackNetworkConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "subnet", testAccIpLoadbalancingVrackNetworkSubnet),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "vlan", testAccIpLoadbalancingVrackNetworkVlan1002),
					resource.TestCheckResourceAttrSet("ovh_iploadbalancing_vrack_network.network", "id"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_vrack_network.network", "nat_ip", testAccIpLoadbalancingVrackNetworkNatIp),
				),
			},
		},
	})
}

var testAccIpLoadbalancingVrackNetworkConfig_basic = fmt.Sprintf(testAccIpLoadbalancingVrackNetworkConfig,
	os.Getenv("OVH_IPLB_SERVICE"),
	os.Getenv("OVH_VRACK"),
	testAccIpLoadbalancingVrackNetworkSubnet,
	testAccIpLoadbalancingVrackNetworkVlan1001,
	testAccIpLoadbalancingVrackNetworkNatIp,
)

var testAccIpLoadbalancingVrackNetworkConfig_update = fmt.Sprintf(testAccIpLoadbalancingVrackNetworkConfig,
	os.Getenv("OVH_IPLB_SERVICE"),
	os.Getenv("OVH_VRACK"),
	testAccIpLoadbalancingVrackNetworkSubnet,
	testAccIpLoadbalancingVrackNetworkVlan1002,
	testAccIpLoadbalancingVrackNetworkNatIp,
)
