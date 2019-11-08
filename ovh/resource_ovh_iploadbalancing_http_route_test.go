package ovh

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIPLoadbalancingRouteHTTPBasicCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	name := "test-route-redirect-https"
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
		PreCheck:     func() { testAccCheckIpLoadbalancingRouteHTTPPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingRouteHTTPDestroy,
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
						"ovh_iploadbalancing_http_route.testroute", "action.859787636.status", actionStatus),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.859787636.target", strings.Replace(actionTarget, "$$", "$", -1)),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_http_route.testroute", "action.859787636.type", actionType),
				),
			},
		},
	})
}

func testAccCheckIpLoadbalancingRouteHTTPPreCheck(t *testing.T) {
	testAccPreCheckIpLoadbalancing(t)
	testAccCheckIpLoadbalancingExists(t)
}

func testAccCheckIPLoadbalancingRouteHTTPDestroy(state *terraform.State) error {
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.ID)
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
