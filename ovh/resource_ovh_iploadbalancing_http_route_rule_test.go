package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIPLoadbalancingRouteHTTPRuleBasicCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	displayName := "Test rule"
	field := "header"
	match := "is"
	negate := "false"
	pattern := "example.com"
	subField := "Host"

	config := fmt.Sprintf(
		testAccCheckOvhIpLoadbalancingHttpRouteRuleConfig_basic,
		serviceName,
		displayName,
		field,
		match,
		negate,
		pattern,
		subField,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCheckIpLoadbalancingRouteHTTPRulePreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingRouteHTTPRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "field", field),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "match", match),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "negate", negate),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "pattern", pattern),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route_rule.testrule", "sub_field", subField),
				),
			},
		},
	})
}

func testAccCheckIpLoadbalancingRouteHTTPRulePreCheck(t *testing.T) {
	testAccPreCheckIpLoadbalancing(t)
	testAccCheckIpLoadbalancingExists(t)
}

func testAccCheckIPLoadbalancingRouteHTTPRuleDestroy(state *terraform.State) error {
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route_rule" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/http/route/%s/rule/%s",
			os.Getenv("OVH_IPLB_SERVICE"),
			resource.Primary.Attributes["route_id"],
			resource.Primary.ID,
		)

		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			return fmt.Errorf("IpLoadbalancing http route rule still exists")
		}
	}
	return nil
}

const testAccCheckOvhIpLoadbalancingHttpRouteRuleConfig_basic = `
resource "ovh_iploadbalancing_http_route" "testroute" {
	service_name = "%s"
	display_name = "%s"
	weight = 0

	action {
		status = 302
		target = "http://example.com"
		type = "redirect"
	}
}

resource "ovh_iploadbalancing_http_route_rule" "testrule" {
	service_name = "${ovh_iploadbalancing_http_route.testroute.service_name}"
	route_id  = "${ovh_iploadbalancing_http_route.testroute.id}"
	display_name = "${ovh_iploadbalancing_http_route.testroute.display_name}"
	field = "%s"
	match = "%s"
	negate = %s
	pattern = "%s"
	sub_field = "%s"
}
`
