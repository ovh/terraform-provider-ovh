package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCheckOvhIpLoadbalancingUdpFrontendConfig_basic = `
resource "ovh_iploadbalancing_udp_frontend" "testfrontend" {
	service_name = "%s"
	display_name = "aaa"
	port         = "102"
	zone         = "all"
}
`

const testAccCheckOvhIpLoadbalancingUdpFrontendConfig_update = `
resource "ovh_iploadbalancing_udp_frontend" "testfrontend" {
   service_name = "%s"
   display_name = "bbb"
   port         = "103,104"
   zone         = "all"
   disabled     = true
}
`

func TestAccIpLoadbalancingUdpFrontend_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFrontendConfig_basic, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "display_name", "aaa"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "disabled", "false"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "port", "102"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "zone", "all"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFrontendConfig_update, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "display_name", "bbb"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "disabled", "true"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "port", "103,104"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_frontend.testfrontend", "zone", "all"),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingUdpFrontend_importBasic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFrontendConfig_basic, iplb),
			},
			{
				ResourceName:                         "ovh_iploadbalancing_udp_frontend.testfrontend",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "frontend_id",
				ImportStateIdFunc:                    testAccIpLoadbalancinUdpFrontend_import("ovh_iploadbalancing_udp_frontend.testfrontend"),
			},
		},
	})
}

func testAccIpLoadbalancinUdpFrontend_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIpLoadbalancingUdpFrontend, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_iploadbalancing_udp_frontend not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testIpLoadbalancingUdpFrontend.Primary.Attributes["service_name"],
			testIpLoadbalancingUdpFrontend.Primary.Attributes["frontend_id"],
		), nil
	}
}
