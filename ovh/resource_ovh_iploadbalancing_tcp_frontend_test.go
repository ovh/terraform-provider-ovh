package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_tcp_frontend", &resource.Sweeper{
		Name: "ovh_iploadbalancing_tcp_frontend",
		F:    testSweepIploadbalancingTcpFrontend,
	})
}

func testSweepIploadbalancingTcpFrontend(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	iplb := os.Getenv("OVH_IPLB_SERVICE")
	if iplb == "" {
		return fmt.Errorf("OVH_IPLB_SERVICE must be set")
	}

	frontends := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/frontend", iplb), &frontends); err != nil {
		return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/frontend:\n\t %q", iplb, err)
	}

	if len(frontends) == 0 {
		log.Print("[DEBUG] No frontend to sweep")
		return nil
	}

	for _, f := range frontends {
		frontend := &IpLoadbalancingTcpFrontend{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/frontend/%d", iplb, f), &frontend); err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/frontend/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(frontend.DisplayName, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/tcp/frontend/%d", iplb, f), nil); err != nil {
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

func TestAccOvhIpLoadbalancingTcpFrontend_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingTcpFrontendConfig_basic, iplb, test_prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "display_name", test_prefix),
					resource.TestCheckNoResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "default_farm_id"),
					resource.TestCheckNoResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "default_ssl_id"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "disabled", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingTcpFrontendConfig_update, iplb, test_prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "display_name", test_prefix),
					resource.TestCheckNoResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "default_farm_id"),
					resource.TestCheckNoResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "default_ssl_id"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "ssl", "false"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "disabled", "false"),
				),
			},
		},
	})
}

func TestAccOvhIpLoadbalancingTcpFrontend_withfarm(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingTcpFrontendConfig_withfarm, iplb, test_prefix, test_prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "display_name", test_prefix),
					resource.TestCheckResourceAttrSet(
						"ovh_iploadbalancing_tcp_frontend.testfrontend", "default_farm_id"),
				),
			},
		},
	})
}

const testAccCheckOvhIpLoadbalancingTcpFrontendConfig_basic = `
resource "ovh_iploadbalancing_tcp_frontend" "testfrontend" {
   service_name = "%s"
   display_name = "%s"
   zone = "all"
   port = "22280,22443"
   disabled = true
   ssl = true
}
`
const testAccCheckOvhIpLoadbalancingTcpFrontendConfig_update = `
resource "ovh_iploadbalancing_tcp_frontend" "testfrontend" {
   service_name = "%s"
   display_name = "%s"
   zone = "all"
   port = "22280,22443"
}
`

const testAccCheckOvhIpLoadbalancingTcpFrontendConfig_withfarm = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}

resource "ovh_iploadbalancing_tcp_farm" "farm" {
   service_name = "${data.ovh_iploadbalancing.iplb.service_name}"
   display_name = "%s"
   zone = "all"
   port = 22280
}

resource "ovh_iploadbalancing_tcp_frontend" "testfrontend" {
   service_name = "${data.ovh_iploadbalancing.iplb.service_name}"
   display_name = "%s"
   zone = "all"
   port = "22280,22443"
   default_farm_id = "${ovh_iploadbalancing_tcp_farm.farm.id}"
}
`
