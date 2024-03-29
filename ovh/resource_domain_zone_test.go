package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccDomainZoneBasic = `
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
}

data "ovh_order_cart_product_plan" "zone" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "dns"
  plan_code      = "zone"
}

resource "ovh_domain_zone" "zone" {
 ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary

 plan {
   duration     = data.ovh_order_cart_product_plan.zone.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.zone.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.zone.selected_price.0.pricing_mode

   configuration {
     label = "zone"
     value = "%s"
   }

   configuration {
     label = "template"
     value = "minimized"
   }
 }
}
`

func init() {
	resource.AddTestSweepers("ovh_domain_zone", &resource.Sweeper{
		Name: "ovh_domain_zone",
		F:    testSweepDomainZone,
	})
}

func testSweepDomainZone(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	zoneNames := make([]string, 0)
	if err := config.OVHClient.Get("/domain/zone", &zoneNames); err != nil {
		return fmt.Errorf("Error calling GET /domain/zone:\n\t %q", err)
	}

	if len(zoneNames) == 0 {
		log.Print("[DEBUG] No domainZone to sweep")
		return nil
	}

	for _, zoneName := range zoneNames {
		if !strings.HasPrefix(zoneName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Will delete domainZone: %v", zoneName)

		terminate := func() (string, error) {
			log.Printf("[DEBUG] Will terminate domainZone %s", zoneName)
			endpoint := fmt.Sprintf(
				"/domain/zone/%s/terminate",
				url.PathEscape(zoneName),
			)
			if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
					return "", nil
				}
				return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
			}
			return zoneName, nil
		}

		confirmTerminate := func(token string) error {
			log.Printf("[DEBUG] Will confirm termination of domainZone %s", zoneName)
			endpoint := fmt.Sprintf(
				"/domain/zone/%s/confirmTermination",
				url.PathEscape(zoneName),
			)
			if err := config.OVHClient.Post(endpoint, &DomainZoneConfirmTerminationOpts{Token: token}, nil); err != nil {
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

func TestAccResourceDomainZone_basic(t *testing.T) {
	domain := os.Getenv("OVH_TESTACC_ORDER_DOMAIN")
	prefix := acctest.RandomWithPrefix(test_prefix)
	name := fmt.Sprintf("%s.%s", prefix, domain)
	config := fmt.Sprintf(
		testAccDomainZoneBasic,
		name,
	)

	t.Logf("[INFO] Will order test zone: %v", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderDomainZone(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_domain_zone.zone", "name", name),
					resource.TestCheckResourceAttrSet(
						"ovh_domain_zone.zone", "urn"),
				),
			},
		},
	})
}
