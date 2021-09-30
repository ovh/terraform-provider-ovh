package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_http_route_rule", &resource.Sweeper{
		Name: "ovh_iploadbalancing_http_route_rule",
		F:    testSweepIploadbalancingHttpRouteRule,
	})
}

func testSweepIploadbalancingHttpRouteRule(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")
	if iplb == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE_TEST is not set. No iploadbalancing_vrack_network to sweep")
		return nil
	}

	routes := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/route", iplb), &routes); err != nil {
		return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/route:\n\t %q", iplb, err)
	}

	if len(routes) == 0 {
		log.Print("[DEBUG] No http route to sweep")
		return nil
	}

	for _, f := range routes {
		route := &IPLoadbalancingHttpRoute{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d", iplb, f), &route); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/route/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*route.DisplayName, test_prefix) {
			continue
		}

		rules := make([]int64, 0)
		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d/rule", iplb, f), &rules); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/route/%d/rule:\n\t %q", iplb, f, err)
		}

		if len(rules) == 0 {
			log.Printf("[DEBUG] No rule to sweep on http route %s/http/route/%d", iplb, f)
			return nil
		}

		for _, s := range rules {
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d/rule/%d", iplb, f, s), nil); err != nil {
					return resource.RetryableError(err)
				}
				// Successful delete
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccIPLoadbalancingHttpRouteRuleBasicCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	displayName := acctest.RandomWithPrefix(test_prefix)
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
		PreCheck:     func() { testAccCheckIpLoadbalancingHttpRouteRulePreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingHttpRouteRuleDestroy,
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

func testAccCheckIpLoadbalancingHttpRouteRulePreCheck(t *testing.T) {
	testAccPreCheckIpLoadbalancing(t)
	testAccCheckIpLoadbalancingExists(t)
}

func testAccCheckIPLoadbalancingHttpRouteRuleDestroy(state *terraform.State) error {
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route_rule" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s/http/route/%s/rule/%s",
			os.Getenv("OVH_IPLB_SERVICE_TEST"),
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
