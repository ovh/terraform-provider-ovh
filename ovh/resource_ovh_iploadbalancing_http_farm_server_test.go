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

	iplb := os.Getenv("OVH_IPLB_SERVICE")
	if iplb == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE is not set. No iploadbalancing_vrack_network to sweep")
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
	prefix := fmt.Sprintf(
		testAccIpLoadbalancingHttpFarmServerConfig_templ,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step0, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step1, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step2, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendB"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step3, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "inactive"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step4, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "proxy_protocol_version", "v2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step5, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "proxy_protocol_version", "v1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step6, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "true"),
				),
			},
		},
	})
}
