package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccOrderCartProductPlanBasic = `
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "%s"
}

data "ovh_order_cart_product_plan" "plan" {
 cart_id        = data.ovh_order_cart.mycart.id
 price_capacity = "renew"
 product        = "%s"
 plan_code      = "%s"
}
`

func TestAccDataSourceOrderCartIpLoadbalancingPlan_basic(t *testing.T) {
	testAccDataSourceOrderCartProductPlan_basic(t, "ipLoadbalancing", "iplb-lb2")
}

func TestAccDataSourceOrderCartCloudPlan_basic(t *testing.T) {
	testAccDataSourceOrderCartProductPlan_basic(t, "cloud", "project")
}

func testAccDataSourceOrderCartProductPlan_basic(t *testing.T, product, planCode string) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccOrderCartProductPlanBasic,
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
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_plan.plan", "plan_code", planCode),
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_plan.plan", "product_type", "delivery"),
					resource.TestCheckResourceAttr(
						"data.ovh_order_cart_product_plan.plan",
						"selected_price.0.duration", "P1M"),
				),
			},
		},
	})
}
