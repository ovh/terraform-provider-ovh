package ovh

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

type TestAccIpLoadbalancingHttpFarmBackendProbeResponse struct {
	Match    string `json:"match"`
	Port     int    `json:"port"`
	Interval int    `json:"interval"`
	Negate   bool   `json:"negate"`
	Pattern  string `json:"pattern"`
	ForceSsl bool   `json:"forceSsl"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	Type     string `json:"type"`
}

type TestAccIpLoadbalancingHttpFarmResponse struct {
	Zone           string                                             `json:"zone"`
	VrackNetworkId int                                                `json:"vrackNetworkId"`
	Port           int                                                `json:"port"`
	Stickiness     string                                             `json:"stickiness"`
	FarmId         int                                                `json:"farmId"`
	Balance        string                                             `json:"balance"`
	Probe          TestAccIpLoadbalancingHttpFarmBackendProbeResponse `json:"probe"`
	DisplayName    string                                             `json:"displayName"`
}

func (r *TestAccIpLoadbalancingHttpFarmResponse) Equals(c *TestAccIpLoadbalancingHttpFarmResponse) bool {
	r.FarmId = 0
	if reflect.DeepEqual(r, c) {
		return true
	}
	return false
}

func testAccIpLoadbalancingHttpFarmTestStep(name, zone string, port, probePort, probeInterval int, probeType string) resource.TestStep {
	expected := &TestAccIpLoadbalancingHttpFarmResponse{
		Zone:        zone,
		Port:        port,
		DisplayName: name,
		Probe: TestAccIpLoadbalancingHttpFarmBackendProbeResponse{
			Port:     probePort,
			Interval: probeInterval,
			Type:     probeType,
		},
	}

	config := fmt.Sprintf(`
	resource "ovh_iploadbalancing_http_farm" "testfarm" {
		service_name = "%s"
		display_name = "%s"
		port = %d
		zone = "%s"
	  
		probe {
		  port = %d
		  interval = %d
		  type = "%s"
		}	  
	}
	`, os.Getenv("OVH_IPLB_SERVICE"), name, port, zone, probePort, probeInterval, probeType)

	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			testAccCheckIpLoadbalancingHttpFarmMatches(expected),
		),
	}
}

func TestAccIpLoadbalancingHttpFarmBasicCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpLoadbalancingHttpFarmDestroy,
		Steps: []resource.TestStep{
			testAccIpLoadbalancingHttpFarmTestStep("test-farm-v1", "all", 8080, 8888, 35, "http"),
			testAccIpLoadbalancingHttpFarmTestStep("test-farm-v2", "all", 8080, 9999, 60, "http"),
		},
	})
}

func testAccCheckIpLoadbalancingHttpFarmMatches(expected *TestAccIpLoadbalancingHttpFarmResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		name := "ovh_iploadbalancing_http_farm.testfarm"
		resource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.ID)
		response := &TestAccIpLoadbalancingHttpFarmResponse{}
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

func testAccCheckIpLoadbalancingHttpFarmDestroy(state *terraform.State) error {
	leftovers := false
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_farm" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			leftovers = true
		}
	}
	if leftovers {
		return fmt.Errorf("IpLoadbalancing farm still exists")
	}
	return nil
}
