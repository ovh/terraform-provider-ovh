package ovh

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

type TestAccIPLoadbalancingRouteHTTPActionResponse struct {
	Target string `json:"target,omitempty"`
	Status int    `json:"status,omitempty"`
	Type   string `json:"type"`
}

type TestAccIPLoadbalancingRouteHTTPResponse struct {
	Weight      int                                           `json:"weight"`
	Action      TestAccIPLoadbalancingRouteHTTPActionResponse `json:"action"`
	RouteID     int                                           `json:"routeId"`
	DisplayName string                                        `json:"displayName"`
	FrontendID  int                                           `json:"frontendId"`
}

func (r *TestAccIPLoadbalancingRouteHTTPResponse) Equals(c *TestAccIPLoadbalancingRouteHTTPResponse) bool {
	r.RouteID = 0
	if reflect.DeepEqual(r, c) {
		return true
	}
	return false
}

func testAccIPLoadbalancingRouteHTTPTestStep(name string, weight int, actionStatus int, actionTarget string, actionType string) resource.TestStep {
	expected := &TestAccIPLoadbalancingRouteHTTPResponse{
		Weight:      weight,
		DisplayName: name,
		Action: TestAccIPLoadbalancingRouteHTTPActionResponse{
			Target: actionTarget,
			Status: actionStatus,
			Type:   actionType,
		},
	}

	config := fmt.Sprintf(`
	resource "ovh_iploadbalancing_http_route" "testroute" {
		service_name = "%s"
		display_name = "%s"
		weight = %d

		action {
		  status = %d
		  target = "%s"
		  type = "%s"
		}
	}
	`, os.Getenv("OVH_IPLB_SERVICE"), name, weight, actionStatus, actionTarget, actionType)

	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			testAccCheckIPLoadbalancingRouteHTTPMatches(expected),
		),
	}
}

func TestAccIPLoadbalancingRouteHTTPBasicCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPLoadbalancingRouteHTTPDestroy,
		Steps: []resource.TestStep{
			testAccIPLoadbalancingRouteHTTPTestStep("test-route-redirect-https", 0, 302, "https://${host}${path}${arguments}", "redirect"),
		},
	})
}

func testAccCheckIPLoadbalancingRouteHTTPMatches(expected *TestAccIPLoadbalancingRouteHTTPResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		name := "ovh_iploadbalancing_http_route.testroute"
		resource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.ID)
		response := &TestAccIPLoadbalancingRouteHTTPResponse{}
		err := config.OVHClient.Get(endpoint, response)
		if err != nil {
			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		}
		if !response.Equals(expected) {
			return fmt.Errorf("%s %s state differs from expected", name, resource.Primary.ID)
		}
		return nil
	}
}

func testAccCheckIPLoadbalancingRouteHTTPDestroy(state *terraform.State) error {
	leftovers := false
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_route" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/route/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			leftovers = true
		}
	}
	if leftovers {
		return fmt.Errorf("IpLoadbalancing route still exists")
	}
	return nil
}
