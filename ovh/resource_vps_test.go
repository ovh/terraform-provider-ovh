package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsBasic = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "myvps" {
  display_name = "%s"
  netboot_mode = "rescue"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]
}
`

func TestAccResourceVps_basic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVpsBasic,
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderVPS(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "netboot_mode", "rescue"),
					resource.TestCheckResourceAttr(
						"ovh_vps.myvps", "display_name", displayName),
				),
			},
		},
	})
}
