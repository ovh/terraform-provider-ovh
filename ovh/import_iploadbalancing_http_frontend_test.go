package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccIpLoadbalancingHttpFrontendConfig = `
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

resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
  service_name    = data.ovh_iploadbalancing.iplb.id
  default_farm_id = ovh_iploadbalancing_http_farm.testfarm.id
  display_name    = "%s"
  zone            = "all"
  port            = 12345
}
`
)

func TestAccIpLoadbalancingHttpFrontend_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingHttpFrontendConfig_basic,
			},
			{
				ResourceName:      "ovh_iploadbalancing_http_frontend.testfrontend",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingHttpFrontendImportId("ovh_iploadbalancing_http_frontend.testfrontend"),
			},
		},
	})
}

func testAccIpLoadbalancingHttpFrontendImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testfrontend, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_http_frontend not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testfrontend.Primary.Attributes["service_name"],
			testfrontend.Primary.Attributes["id"],
		), nil
	}
}

var testAccIpLoadbalancingHttpFrontendConfig_basic = fmt.Sprintf(testAccIpLoadbalancingHttpFrontendConfig,
	os.Getenv("OVH_IPLB_SERVICE"), "testfarm", 12345, "all", "testfrontend")
