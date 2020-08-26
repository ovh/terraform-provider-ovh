package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccIpLoadbalancingHttpFarmServerConfig = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}
resource "ovh_iploadbalancing_http_farm" "testfarm" {
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

resource "ovh_iploadbalancing_http_farm_server" "testfarmserver" {
  service_name           = data.ovh_iploadbalancing.iplb.id
  farm_id                = ovh_iploadbalancing_http_farm.testfarm.id
  display_name           = "%s"
  probe                  = "true"
  proxy_protocol_version = "v1"
  status                 = "active"
  address                = "%s"
}
`
)

func TestAccIpLoadbalancingHttpFarmServer_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingHttpFarmServerConfig_basic,
			},
			{
				ResourceName:      "ovh_iploadbalancing_http_farm_server.testfarmserver",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingHttpFarmServerImportId("ovh_iploadbalancing_http_farm_server.testfarmserver"),
			},
		},
	})
}

func testAccIpLoadbalancingHttpFarmServerImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testFarmServer, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_http_farm not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testFarmServer.Primary.Attributes["service_name"],
			testFarmServer.Primary.Attributes["farm_id"],
			testFarmServer.Primary.Attributes["id"],
		), nil
	}
}

// an OVH IPv4 is required for servers
// ping.ovh.net ip is used for test purposes
var httpServerAddress = "198.27.92.1"
var testAccIpLoadbalancingHttpFarmServerConfig_basic = fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig,
	os.Getenv("OVH_IPLB_SERVICE"), "testfarm", 12345, "all", "testserver", httpServerAddress)
