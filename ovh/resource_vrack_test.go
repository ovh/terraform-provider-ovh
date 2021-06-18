package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccVrackBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product_plan" "vrack" {
 cart_id        = data.ovh_order_cart.mycart.id
 price_capacity = "renew"
 product        = "vrack"
 plan_code      = "vrack"
}

resource "ovh_vrack" "vrack" {
 ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
 name          = "%s"
 payment_mean  = "fidelity"
 description   = "%s"

 plan {
   duration     = data.ovh_order_cart_product_plan.vrack.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.vrack.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.vrack.selected_price.0.pricing_mode
 }
}
`

func TestAccResourceVrack_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVrackBasic,
		desc,
		name,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderVrack(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vrack.vrack", "name", name),
					resource.TestCheckResourceAttr(
						"ovh_vrack.vrack", "description", desc),
					resource.TestCheckResourceAttrSet(
						"ovh_vrack.vrack", "service_name"),
				),
			},
		},
	})
}
