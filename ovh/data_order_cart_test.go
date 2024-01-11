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
)

const testAccOrderCartBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}
`

const testAccOrderCartAssignBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
 assign         = true
}
`

func init() {
	resource.AddTestSweepers("ovh_order_cart", &resource.Sweeper{
		Name: "ovh_order_cart",
		F:    testSweepOrderCart,
	})
}

func testSweepOrderCart(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	carts := make([]string, 0)
	if err := config.OVHClient.Get("/order/cart", &carts); err != nil {
		return fmt.Errorf("Error calling GET /order/cart:\n\t %q", err)
	}

	if len(carts) == 0 {
		log.Print("[DEBUG] No carts to sweep")
		return nil
	}

	for _, cart := range carts {
		r := &OrderCart{}
		log.Printf("[DEBUG] Will get order cart: %v", cart)
		endpoint := fmt.Sprintf(
			"/order/cart/%s",
			url.PathEscape(cart),
		)

		if err := config.OVHClient.Get(endpoint, r); err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

		if r.Description != nil && !strings.HasPrefix(*r.Description, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := config.OVHClient.Delete(endpoint, nil); err != nil {
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

func TestAccDataSourceOrderCart_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartBasic,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart.mycart", "description", desc),
					resource.TestCheckResourceAttrSet(
						"data.ovh_order_cart.mycart", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceOrderCartAssign_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartAssignBasic,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart.mycart", "description", desc),
					resource.TestCheckResourceAttrSet(
						"data.ovh_order_cart.mycart", "id"),
				),
			},
		},
	})
}
