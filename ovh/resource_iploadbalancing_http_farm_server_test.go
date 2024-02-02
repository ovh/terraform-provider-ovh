package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAccIpLoadbalancingHttpFarmServerConfig_templ = `
data ovh_iploadbalancing iplb {
  service_name = "%s"
}

resource ovh_iploadbalancing_http_farm testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port = 8080
  zone = "all"
  probe {
    port     = 8080
    interval = 30
    type     = "http"
  }
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step0 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name = data.ovh_iploadbalancing.iplb.id
  farm_id      = ovh_iploadbalancing_http_farm.testacc.id
  address      = "10.0.0.11"
  status       = "active"
  display_name = "testBackendA"
  port         = 80
  weight       = 3
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step1 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendA"
  port = 8080
  weight = 3
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step2 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  weight = 2
  probe = true
  backup = true
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step3 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "inactive"
  display_name = "testBackendB"
  port = 80
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step4 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v2"
  ssl = true
  weight = 2
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step5 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v1"
  ssl    = true
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step6 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  ssl = true
  backup = true
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step7 = `
%s
resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  ssl = true
  backup = true
  on_marked_down = "shutdown-sessions"
  proxy_protocol_version = "v1"
}
`
	TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME = "ovh_iploadbalancing_http_farm_server.testacc"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_http_farm_server", &resource.Sweeper{
		Name: "ovh_iploadbalancing_http_farm_server",
		F:    testSweepIploadbalancingHttpFarmServer,
	})
}

func testSweepIploadbalancingHttpFarmServer(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")
	if iplb == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE_TEST is not set. No iploadbalancing_vrack_network to sweep")
		return nil
	}

	farms := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/farm", iplb), &farms); err != nil {
		return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/farm:\n\t %q", iplb, err)
	}

	if len(farms) == 0 {
		log.Print("[DEBUG] No http farm to sweep")
		return nil
	}

	for _, f := range farms {
		farm := &IpLoadbalancingFarm{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%d", iplb, f), &farm); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/farm/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*farm.DisplayName, test_prefix) {
			continue
		}

		servers := make([]int64, 0)
		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%d/server", iplb, f), &servers); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/http/farm/%d/server:\n\t %q", iplb, f, err)
		}

		if len(servers) == 0 {
			log.Printf("[DEBUG] No server to sweep on http farm %s/http/farm/%d", iplb, f)
			return nil
		}

		for _, s := range servers {
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%d/server/%d", iplb, f, s), nil); err != nil {
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

func TestAccIpLoadbalancingHttpFarmServerBasic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	prefix := fmt.Sprintf(
		testAccIpLoadbalancingHttpFarmServerConfig_templ,
		serviceName,
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step0, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "display_name", "testBackendA"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "80"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "3"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step1, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "display_name", "testBackendA"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "3"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step2, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "display_name", "testBackendB"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "2"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "probe", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step3, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "inactive"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "80"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "probe", "false"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "backup", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step4, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "2"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "proxy_protocol_version", "v2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step6, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "1"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step5, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "1"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "backup", "false"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "proxy_protocol_version", "v1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step7, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "address", "10.0.0.11"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "status", "active"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "port", "8080"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "weight", "1"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "backup", "true"),
					resource.TestCheckResourceAttr(TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME, "on_marked_down", "shutdown-sessions"),
				),
			},
			{
				ResourceName:      TEST_ACC_IPLOADBALANCING_HTTP_FARM_SRV_RES_NAME,
				ImportState:       true,
				ImportStateIdFunc: getImportStateId,
				ImportStateVerify: true,
			},
		},
	})
}

func getImportStateId(state *terraform.State) (string, error) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")

	if len(state.Modules) != 1 {
		return "", fmt.Errorf("Unexpected modules length %d", len(state.Modules))
	}
	var mod = state.Modules[0]
	var farmId = mod.Resources["ovh_iploadbalancing_http_farm.testacc"].Primary.ID
	var serverId = mod.Resources["ovh_iploadbalancing_http_farm_server.testacc"].Primary.ID

	var result = serviceName + "/" + farmId + "/" + serverId
	log.Printf("[DEBUG] ID to import %s", result)
	return result, nil

}
