package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCheckOvhIpLoadbalancingUdpFarmConfig_basic = `
resource "ovh_iploadbalancing_udp_farm" "testfarm" {
	service_name = "%s"
	display_name = "aaa"
	port         = 102
	zone         = "all"
}
`

const testAccCheckOvhIpLoadbalancingUdpFarmConfig_update = `
resource "ovh_iploadbalancing_udp_farm" "testfarm" {
   service_name = "%s"
   display_name = "bbb"
   port         = 103
   zone         = "all"
}
`

func TestAccIpLoadbalancingUdpFarm_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmConfig_basic, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "display_name", "aaa"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "port", "102"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "zone", "all"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmConfig_update, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "display_name", "bbb"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "port", "103"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm.testfarm", "zone", "all"),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingUdpFarm_importBasic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmConfig_basic, iplb),
			},
			{
				ResourceName:                         "ovh_iploadbalancing_udp_farm.testfarm",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "farm_id",
				ImportStateIdFunc:                    testAccIpLoadbalancinUdpFarm_import("ovh_iploadbalancing_udp_farm.testfarm"),
			},
		},
	})
}

func testAccIpLoadbalancinUdpFarm_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIpLoadbalancingUdpFarm, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_iploadbalancing_udp_farm not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testIpLoadbalancingUdpFarm.Primary.Attributes["service_name"],
			testIpLoadbalancingUdpFarm.Primary.Attributes["farm_id"],
		), nil
	}
}
