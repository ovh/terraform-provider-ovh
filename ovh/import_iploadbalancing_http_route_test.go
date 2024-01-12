package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIpLoadbalancingHttpRoute_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpRouteConfig_basic, os.Getenv("OVH_IPLB_SERVICE_TEST"), "testroute", "0", "302", "https://test.url", "redirect"),
			},
			{
				ResourceName:      "ovh_iploadbalancing_http_route.testroute",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingHttpRouteImportId("ovh_iploadbalancing_http_route.testroute"),
			},
		},
	})
}

func testAccIpLoadbalancingHttpRouteImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testroute, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_http_route not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testroute.Primary.Attributes["service_name"],
			testroute.Primary.Attributes["id"],
		), nil
	}
}
