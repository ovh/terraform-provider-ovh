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

const testAccIpLoadbalancingBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product_plan" "iplb" {
 cart_id        = data.ovh_order_cart.mycart.id
 price_capacity = "renew"
 product        = "ipLoadbalancing"
 plan_code      = "iplb-lb1"
}

data "ovh_order_cart_product_options_plan" "bhs" {
 cart_id           = data.ovh_order_cart_product_plan.iplb.cart_id
 price_capacity    = data.ovh_order_cart_product_plan.iplb.price_capacity
 product           = data.ovh_order_cart_product_plan.iplb.product
 plan_code         = data.ovh_order_cart_product_plan.iplb.plan_code
 options_plan_code = "iplb-zone-lb1-rbx"
}

resource "ovh_iploadbalancing" "iplb-lb1" {
 ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
 display_name   = "%s"

 plan {
   duration     = data.ovh_order_cart_product_plan.iplb.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.iplb.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.iplb.selected_price.0.pricing_mode
 }

 plan_option {
   duration     = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_options_plan.bhs.plan_code
   pricing_mode = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.pricing_mode
 }
}
`

const testAccIpLoadbalancingInternal = `
resource "ovh_iploadbalancing" "iplb-internal" {
 ovh_subsidiary = "fr"
 display_name   = "%s"

 plan {
   catalog_name = "iplb_private_beta"
   duration     = "P1M"
   plan_code    = "iplb-service-monitoring"
   pricing_mode = "default"
 }

 plan_option {
   catalog_name = "iplb_private_beta"
   duration     = "P1M"
   plan_code    = "iplb-zone-gra"
   pricing_mode = "default"
 }
}
`

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing", &resource.Sweeper{
		Name: "ovh_iploadbalancing",
		F:    testSweepIpLoadbalancing,
	})
}

func testSweepIpLoadbalancing(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceNames := make([]string, 0)
	if err := config.OVHClient.Get("/ipLoadbalancing", &serviceNames); err != nil {
		return fmt.Errorf("Error calling GET /ipLoadbalancing:\n\t %q", err)
	}

	if len(serviceNames) == 0 {
		log.Print("[DEBUG] No ipLoadbalancing to sweep")
		return nil
	}

	for _, serviceName := range serviceNames {
		r := &IpLoadbalancing{}
		log.Printf("[DEBUG] Will get ipLoadbalancing: %v", serviceName)
		endpoint := fmt.Sprintf(
			"/ipLoadbalancing/%s",
			url.PathEscape(serviceName),
		)

		if err := config.OVHClient.Get(endpoint, r); err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

		if r.DisplayName == nil || !strings.HasPrefix(*r.DisplayName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Will delete ipLoadbalancing: %v", serviceName)

		terminate := func() (string, error) {
			log.Printf("[DEBUG] Will terminate ipLoadbalancing %s", serviceName)
			endpoint := fmt.Sprintf(
				"/ipLoadbalancing/%s/terminate",
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
			log.Printf("[DEBUG] Will confirm termination of ipLoadbalancing %s", serviceName)
			endpoint := fmt.Sprintf(
				"/ipLoadbalancing/%s/confirmTermination",
				url.PathEscape(serviceName),
			)
			if err := config.OVHClient.Post(endpoint, &IpLoadbalancingConfirmTerminationOpts{Token: token}, nil); err != nil {
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

func TestAccResourceIpLoadbalancing_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccIpLoadbalancingBasic,
		desc,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderIpLoadbalancing(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_iploadbalancing.iplb-lb1", "ipv4"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing.iplb-lb1", "display_name", desc),
					resource.TestCheckResourceAttrSet(
						"ovh_iploadbalancing.iplb-lb1", "urn"),
				),
			},
		},
	})
}

func TestAccResourceIpLoadbalancing_internal(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccIpLoadbalancingInternal,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderIpLoadbalancing(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_iploadbalancing.iplb-internal", "ipv4"),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing.iplb-internal", "display_name", desc),
					resource.TestCheckResourceAttrSet(
						"ovh_iploadbalancing.iplb-internal", "urn"),
				),
			},
		},
	})
}
