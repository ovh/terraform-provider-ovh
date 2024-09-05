package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCheckOvhIpLoadbalancingUdpFarmServerConfig_basic = `
resource "ovh_iploadbalancing_udp_farm" "testfarm" {
	service_name = "%s"
	display_name = "aaa"
	port         = 102
	zone         = "all"
}

resource "ovh_iploadbalancing_udp_farm_server" "testserver" {
	service_name = "%s"
	farm_id      = "${ovh_iploadbalancing_udp_farm.testfarm.farm_id}"
	display_name = "mybackend1"
	address      = "10.0.0.11"
	status       = "active"
	port         = 80
  }
`

const testAccCheckOvhIpLoadbalancingUdpFarmServerConfig_update = `
resource "ovh_iploadbalancing_udp_farm" "testfarm" {
	service_name = "%s"
	display_name = "aaa"
	port         = 102
	zone         = "all"
}

resource "ovh_iploadbalancing_udp_farm_server" "testserver" {
	service_name = "%s"
	farm_id      = "${ovh_iploadbalancing_udp_farm.testfarm.farm_id}"
	display_name = "mybackend2"
	address      = "10.0.0.11"
	status       = "inactive"
	port         = 81
  }
`

func TestAccIpLoadbalancingUdpFarmServer_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmServerConfig_basic, iplb, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "display_name", "mybackend1"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "port", "80"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "status", "active"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmServerConfig_update, iplb, iplb),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "display_name", "mybackend2"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "port", "81"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_udp_farm_server.testserver", "status", "inactive"),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingUdpFarmServer_importBasic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIpLoadbalancing(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingUdpFarmServerConfig_basic, iplb, iplb),
			},
			{
				ResourceName:                         "ovh_iploadbalancing_udp_farm_server.testserver",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "server_id",
				ImportStateIdFunc:                    testAccIpLoadbalancinUdpFarmServer_import("ovh_iploadbalancing_udp_farm_server.testserver"),
			},
		},
	})
}

func testAccIpLoadbalancinUdpFarmServer_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIpLoadbalancingUdpFarmServer, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_iploadbalancing_udp_farm_server not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testIpLoadbalancingUdpFarmServer.Primary.Attributes["service_name"],
			testIpLoadbalancingUdpFarmServer.Primary.Attributes["farm_id"],
			testIpLoadbalancingUdpFarmServer.Primary.Attributes["server_id"],
		), nil
	}
}
