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
	resource.AddTestSweepers("ovh_iploadbalancing_tcp_route", &resource.Sweeper{
		Name: "ovh_iploadbalancing_tcp_route",
		Dependencies: []string{
			"ovh_iploadbalancing_tcp_route_rule",
		},
		F: testSweepIploadbalancingTcpRoute,
	})
}

func testSweepIploadbalancingTcpRoute(region string) error {
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
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/route", iplb), &routes); err != nil {
		return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/route:\n\t %q", iplb, err)
	}

	if len(routes) == 0 {
		log.Print("[DEBUG] No tcp route to sweep")
		return nil
	}

	for _, f := range routes {
		route := &IPLoadbalancingTcpRoute{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%d", iplb, f), &route); err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/route/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*route.DisplayName, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%d", iplb, f), nil); err != nil {
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

func TestAccIPLoadbalancingTcpRouteBasicCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	name := acctest.RandomWithPrefix(test_prefix)
	weight := "0"
	config := fmt.Sprintf(
		testAccCheckOvhIpLoadbalancingTcpRouteConfig_basic,
		serviceName,
		name,
		weight,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCheckIpLoadbalancingTcpRoutePreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingTcpRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_route.testroute", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_route.testroute", "display_name", name),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_route.testroute", "weight", weight),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_route.testroute", "action.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_route.testroute", "action.0.type", "reject"),
				),
			},
		},
	})
}

func testAccCheckIpLoadbalancingTcpRoutePreCheck(t *testing.T) {
	testAccPreCheckIpLoadbalancing(t)
	testAccCheckIpLoadbalancingExists(t)
}

func testAccCheckIPLoadbalancingTcpRouteDestroy(state *terraform.State) error {
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_tcp_route" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/route/%s", os.Getenv("OVH_IPLB_SERVICE_TEST"), resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			return fmt.Errorf("IpLoadbalancing route still exists")
		}
	}
	return nil
}

const testAccCheckOvhIpLoadbalancingTcpRouteConfig_basic = `
resource "ovh_iploadbalancing_tcp_route" "testroute" {
	service_name = "%s"
	display_name = "%s"
	weight = %s

	action {
	  type = "reject"
	}
}
`
