package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccIpLoadbalancingTcpFarmConfig = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}

resource "ovh_iploadbalancing_tcp_farm" "testfarm" {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port             = "%d"
  zone             = "%s"
  balance 		   = "roundrobin"
  probe {
        interval = 30
        type = "oco"
  }
}
`
)

func TestAccIpLoadbalancingTcpFarm_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingTcpFarmConfig_basic,
			},
			{
				ResourceName:      "ovh_iploadbalancing_tcp_farm.testfarm",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingTcpFarmImportId("ovh_iploadbalancing_tcp_farm.testfarm"),
			},
		},
	})
}

func testAccIpLoadbalancingTcpFarmImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testfarm, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_tcp_farm not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testfarm.Primary.Attributes["service_name"],
			testfarm.Primary.Attributes["id"],
		), nil
	}
}

var testAccIpLoadbalancingTcpFarmConfig_basic = fmt.Sprintf(testAccIpLoadbalancingTcpFarmConfig,
	os.Getenv("OVH_IPLB_SERVICE"), "testfarm", 12345, "all")
