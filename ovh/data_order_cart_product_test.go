package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccOrderCartProductBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product" "plans" {
 cart_id = data.ovh_order_cart.mycart.id
 product = "%s"
}
`

func TestAccDataSourceOrderCartIpLoadbalancing_basic(t *testing.T) {
	testAccDataSourceOrderCartProduct_basic(t, "ipLoadbalancing")
}

func TestAccDataSourceOrderCartCloud_basic(t *testing.T) {
	testAccDataSourceOrderCartProduct_basic(t, "cloud")
}

func testAccDataSourceOrderCartProduct_basic(t *testing.T, product string) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartProductBasic,
		desc,
		product,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_order_cart_product.plans", "result.#"),
				),
			},
		},
	})
}
