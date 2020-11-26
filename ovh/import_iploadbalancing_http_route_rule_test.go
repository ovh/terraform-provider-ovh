package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIpLoadbalancingHttpRouteRule_importBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpRouteRuleConfig_basic, os.Getenv("OVH_IPLB_SERVICE"), "Test rule", "header", "is", "false", "example.com", "Host"),
			},
			{
				ResourceName:      "ovh_iploadbalancing_http_route_rule.testrule",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingHttpRouteRuleImportId("ovh_iploadbalancing_http_route_rule.testrule"),
			},
		},
	})
}

func testAccIpLoadbalancingHttpRouteRuleImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testRouteRule, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_route_rule not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testRouteRule.Primary.Attributes["service_name"],
			testRouteRule.Primary.Attributes["route_id"],
			testRouteRule.Primary.Attributes["id"],
		), nil
	}
}
