package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccIpServiceBasic = `
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
  description    = "%s"
}

data "ovh_order_cart_product_plan" "ipblock" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ip"
  plan_code      = "ip-v4-s30-ripe"
}

resource "ovh_ip_service" "ipblock" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description   = "%s"

 plan {
   duration     = data.ovh_order_cart_product_plan.ipblock.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.ipblock.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.ipblock.selected_price.0.pricing_mode

   configuration {
     label = "country"
     value = "FR"
   }
 }
}
`

func init() {
	resource.AddTestSweepers("ovh_ip_service", &resource.Sweeper{
		Name:         "ovh_ip_service",
		Dependencies: []string{"ovh_vrack_ip"},
		F:            testSweepIp,
	})
}

func testSweepIp(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceNames := make([]string, 0)
	if err := config.OVHClient.Get("/ip/service", &serviceNames); err != nil {
		return fmt.Errorf("Error calling GET /ip/service:\n\t %q", err)
	}

	if len(serviceNames) == 0 {
		log.Print("[DEBUG] No ip to sweep")
		return nil
	}

	for _, serviceName := range serviceNames {
		r := &IpService{}
		log.Printf("[DEBUG] Will get ip: %v", serviceName)
		endpoint := fmt.Sprintf(
			"/ip/service/%s",
			url.PathEscape(serviceName),
		)

		if err := config.OVHClient.Get(endpoint, r); err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

		if r.Description == nil || !strings.HasPrefix(*r.Description, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Will delete ip: %v", serviceName)

		terminate := func() (string, error) {
			log.Printf("[DEBUG] Will terminate ip %s", serviceName)
			endpoint := fmt.Sprintf(
				"/ip/service/%s/terminate",
				url.PathEscape(serviceName),
			)
			if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
					return "", nil
				}
				return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
			}
			return serviceName, nil
		}

		confirmTerminate := func(token string) error {
			log.Printf("[DEBUG] Will confirm termination of ip %s", serviceName)
			endpoint := fmt.Sprintf(
				"/ip/service/%s/confirmTermination",
				url.PathEscape(serviceName),
			)
			if err := config.OVHClient.Post(endpoint, &IpServiceConfirmTerminationOpts{Token: token}, nil); err != nil {
				return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
			}
			return nil
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := orderDeleteFromResource(nil, config, terminate, confirmTerminate); err != nil {
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

func TestAccResourceIpService_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccIpServiceBasic,
		desc,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderIpService(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_ip_service.ipblock", "ip"),
					resource.TestCheckResourceAttr(
						"ovh_ip_service.ipblock", "description", desc),
				),
			},
		},
	})
}
