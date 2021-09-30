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
	resource.AddTestSweepers("ovh_iploadbalancing_http_route", &resource.Sweeper{
		Name: "ovh_iploadbalancing_http_route",
		Dependencies: []string{
			"ovh_iploadbalancing_http_route_rule",
		},
		F: testSweepIploadbalancingHttpRoute,
	})
}

func testSweepIploadbalancingHttpRoute(region string) error {
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
		return fmt.Errorf("Error calling /ipLoadbalancing/%s/http/route:\n\t %q", iplb, err)
	}

	if len(routes) == 0 {
		log.Print("[DEBUG] No http route to sweep")
		return nil
	}

	for _, f := range routes {
		route := &IPLoadbalancingHttpRoute{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d", iplb, f), &route); err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s/http/route/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*route.DisplayName, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/http/route/%d", iplb, f), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccIPLoadbalancingHttpRouteBasicCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)
	weight := "0"
	actionStatus := "302"
	actionTarget := "https://$${host}$${path}$${arguments}"
	actionType := "redirect"

	config := fmt.Sprintf(
		testAccCheckOvhIpLoadbalancingHttpRouteConfig_basic,
		serviceName,
		name,
		weight,
		actionStatus,
		actionTarget,
		actionType,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCheckIpLoadbalancingHttpRoutePreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingHttpRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "display_name", name),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "weight", weight),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.0.status", actionStatus),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.0.target", strings.Replace(actionTarget, "$$", "$", -1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.0.type", actionType),
				),
			},
		},
	})
}

func testAccCheckIpLoadbalancingHttpRoutePreCheck(t *testing.T) {
	testAccPreCheckIpLoadbalancing(t)
	testAccCheckIpLoadbalancingExists(t)
}

func testAccCheckIPLoadbalancingHttpRouteDestroy(state *terraform.State) error {
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", os.Getenv("OVH_IPLB_SERVICE_TEST"), resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			return fmt.Errorf("IpLoadbalancing route still exists")
		}
	}
	return nil
}

const testAccCheckOvhIpLoadbalancingHttpRouteConfig_basic = `
resource "ovh_iploadbalancing_http_route" "testroute" {
	service_name = "%s"
	display_name = "%s"
	weight = %s

	action {
	  status = %s
	  target = "%s"
	  type = "%s"
	}
}
`
