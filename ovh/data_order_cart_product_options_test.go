package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccOrderCartProductOptionsBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product_options" "options" {
 cart_id   = data.ovh_order_cart.mycart.id
 product   = "%s"
 plan_code = "%s"
}
`

func TestAccDataSourceOrderCartIpLoadbalancingOptions_basic(t *testing.T) {
	testAccDataSourceOrderCartProductOptions_basic(t, "ipLoadbalancing", "iplb-lb2")
}

func TestAccDataSourceOrderCartCloudOptions_basic(t *testing.T) {
	testAccDataSourceOrderCartProductOptions_basic(t, "cloud", "project")
}

func testAccDataSourceOrderCartProductOptions_basic(t *testing.T, product, planCode string) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartProductOptionsBasic,
		desc,
		product,
		planCode,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_order_cart_product_options.options", "result.#"),
				),
			},
		},
	})
}
