package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccOrderCartProductOptionsPlanBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product_options_plan" "plan" {
 cart_id           = data.ovh_order_cart.mycart.id
 price_capacity    = "renew"
 product           = "%s"
 plan_code         = "%s"
 options_plan_code = "%s"
}
`

func TestAccDataSourceOrderCartIpLoadbalancingOptionsPlan_basic(t *testing.T) {
	testAccDataSourceOrderCartProductOptionsPlan_basic(
		t,
		"ipLoadbalancing",
		"iplb-lb2",
		"iplb-zone-lb2-bhs",
	)
}

func TestAccDataSourceOrderCartCloudOptionsPlan_basic(t *testing.T) {
	testAccDataSourceOrderCartProductOptionsPlan_basic(
		t,
		"cloud",
		"project",
		"certification.hds",
	)
}

func testAccDataSourceOrderCartProductOptionsPlan_basic(t *testing.T, product, planCode, optionsCode string) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartProductOptionsPlanBasic,
		desc,
		product,
		planCode,
		optionsCode,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_options_plan.plan", "plan_code", optionsCode),
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_options_plan.plan", "product_type", "delivery"),
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_options_plan.plan",
						"selected_price.0.duration", "P1M"),
				),
			},
		},
	})
}
