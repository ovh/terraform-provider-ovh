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
	testAccIpLoadbalancingTcpFarmServerConfig_templ = `
data ovh_iploadbalancing iplb {
  service_name = "%s"
}

resource ovh_iploadbalancing_tcp_farm testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port = 8080
  zone = "all"
  probe {
    port     = 8080
    interval = 30
    type     = "tcp"
  }
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step0 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name = data.ovh_iploadbalancing.iplb.id
  farm_id      = ovh_iploadbalancing_tcp_farm.testacc.id
  address      = "10.0.0.11"
  status       = "active"
  display_name = "testBackendA"
  port         = 80
  weight       = 3
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step1 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendA"
  port = 8080
  weight = 3
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step2 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  weight = 2
  probe = true
  backup = true
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step3 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "inactive"
  display_name = "testBackendB"
  port = 80
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step4 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v2"
  ssl = true
  weight = 2
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step5 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v1"
  ssl    = true
  backup = false
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step6 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  ssl = true
  backup = true
}
`
	testAccIpLoadbalancingTcpFarmServerConfig_step7 = `
%s

resource ovh_iploadbalancing_tcp_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_tcp_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  ssl = true
  backup = true
  on_marked_down = "shutdown-sessions"
}
`
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_tcp_farm_server", &resource.Sweeper{
		Name: "ovh_iploadbalancing_tcp_farm_server",
		F:    testSweepIploadbalancingTcpFarmServer,
	})
}

func testSweepIploadbalancingTcpFarmServer(region string) error {
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
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm", iplb), &farms); err != nil {
		return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/tcp/farm:\n\t %q", iplb, err)
	}

	if len(farms) == 0 {
		log.Print("[DEBUG] No tcp farm to sweep")
		return nil
	}

	for _, f := range farms {
		farm := &IpLoadbalancingFarm{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d", iplb, f), &farm); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/tcp/farm/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*farm.DisplayName, test_prefix) {
			continue
		}

		servers := make([]int64, 0)
		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server", iplb, f), &servers); err != nil {
			return fmt.Errorf("Error calling GET /ipLoadbalancing/%s/tcp/farm/%d/server:\n\t %q", iplb, f, err)
		}

		if len(servers) == 0 {
			log.Printf("[DEBUG] No server to sweep on tcp farm %s/tcp/farm/%d", iplb, f)
			return nil
		}

		for _, s := range servers {
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server/%d", iplb, f, s), nil); err != nil {
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

func TestAccIpLoadbalancingTcpFarmServerBasic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	prefix := fmt.Sprintf(
		testAccIpLoadbalancingTcpFarmServerConfig_templ,
		os.Getenv("OVH_IPLB_SERVICE_TEST"),
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step0, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step1, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step2, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "display_name", "testBackendB"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "probe", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step3, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "inactive"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "probe", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "backup", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step4, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "proxy_protocol_version", "v2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step5, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "backup", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "proxy_protocol_version", "v1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step6, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingTcpFarmServerConfig_step7, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "backup", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm_server.testacc", "on_marked_down", "shutdown-sessions"),
				),
			},
		},
	})
}
