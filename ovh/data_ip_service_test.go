package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceIpService_basic = `
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

data "ovh_ip_service" "myip" {
 service_name  = ovh_ip_service.ipblock.service_name
}
`

func TestAccDataSourceIpService_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccDataSourceIpService_basic,
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
					resource.TestCheckResourceAttr(
						"data.ovh_ip_service.myip",
						"description",
						desc,
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_ip_service.myip",
						"ip",
					),
				),
			},
		},
	})
}
