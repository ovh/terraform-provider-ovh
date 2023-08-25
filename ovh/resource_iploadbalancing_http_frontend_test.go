package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_http_frontend", &resource.Sweeper{
		Name: "ovh_iploadbalancing_http_frontend",
		F:    testSweepIploadbalancingHttpFrontend,
	})
}

func testSweepIploadbalancingHttpFrontend(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")
	if iplb == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE_TEST is not set. No iploadbalancing_vrack_network to sweep")
		return nil
	}

	frontends := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/frontend", iplb), &frontends); err != nil {
		return fmt.Errorf("Error calling /ipLoadbalancing/%s/http/frontend:\n\t %q", iplb, err)
	}

	if len(frontends) == 0 {
		log.Print("[DEBUG] No frontend to sweep")
		return nil
	}

	for _, f := range frontends {
		frontend := &IpLoadbalancingHttpFrontend{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%d", iplb, f), &frontend); err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s/http/frontend/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(frontend.DisplayName, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/http/frontend/%d", iplb, f), nil); err != nil {
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

func TestAccIpLoadbalancingHttpFrontend_basic(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")
	ipfo := os.Getenv("OVH_IPLB_IPFO_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpFrontendConfig_basic, iplb, test_prefix, ipfo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "display_name", test_prefix),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "port", "22280"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "disabled", "true"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "allowed_source.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpFrontendConfig_update, iplb, test_prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "display_name", test_prefix),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "ssl", "false"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "port", "22280,22443"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "disabled", "false"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "allowed_source.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpFrontendConfig_basic, iplb, test_prefix, ipfo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "display_name", test_prefix),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "ssl", "true"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "port", "22280"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "disabled", "true"),
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "allowed_source.#", "0"),
				),
			},

			{
				ResourceName:        TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME,
				ImportState:         true,
				ImportStateIdPrefix: iplb + "/",
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccIpLoadbalancingHttpFrontend_withfarm(t *testing.T) {
	iplb := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhIpLoadbalancingHttpFrontendConfig_withfarm, iplb, test_prefix, test_prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "display_name", test_prefix),
					resource.TestCheckResourceAttrSet(
						TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME, "default_farm_id"),
				),
			},
		},
	})
}

const testAccCheckOvhIpLoadbalancingHttpFrontendConfig_basic = `
resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
   service_name = "%s"
   display_name = "%s"
   zone     = "all"
   port     = "22280"
   disabled = true
   ssl      = true
   dedicated_ipfo = %s
}
`
const testAccCheckOvhIpLoadbalancingHttpFrontendConfig_update = `
resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
   service_name   = "%s"
   display_name   = "%s"
   zone           = "all"
   port           = "22280,22443"
   allowed_source = ["8.8.8.8/32"]
   http_header    = ["X-Ip-Header %%ci", "X-Port-Header %%cp"]
}
`

const testAccCheckOvhIpLoadbalancingHttpFrontendConfig_withfarm = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}

resource "ovh_iploadbalancing_http_farm" "farm" {
   service_name = data.ovh_iploadbalancing.iplb.service_name
   display_name = "%s"
   zone         = "all"
   port         = 22280
}

resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
   service_name    = data.ovh_iploadbalancing.iplb.service_name
   display_name    = "%s"
   zone            = "all"
   port            = "22280,22443"
   default_farm_id = ovh_iploadbalancing_http_farm.farm.id
}
`
const TEST_ACC_IPLOADBALANCING_HTTP_FRONTEND_RES_NAME = "ovh_iploadbalancing_http_frontend.testfrontend"
