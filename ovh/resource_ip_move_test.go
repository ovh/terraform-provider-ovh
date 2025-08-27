package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccIpMoveConfig = `
data "ovh_order_cart" "mycart" {
	ovh_subsidiary = "fr"
	description    = "Test cart"
}

data "ovh_order_cart_product_plan" "ipblock" {
	cart_id        = data.ovh_order_cart.mycart.id
	price_capacity = "renew"
	product        = "ip"
	plan_code      = "ip-failover-ripe"
}

resource "ovh_ip_service" "ipblock" {
	ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
	description   = "Test IP"

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

resource "ovh_ip_move" "move" {
    ip = ovh_ip_service.ipblock.ip
    routed_to {
        service_name = "%s"
    }
}
`

func TestAccIpMove_basic(t *testing.T) {
	routedToServiceName := os.Getenv("OVH_IP_MOVE_SERVICE_NAME_TEST")

	moveConfig := fmt.Sprintf(testAccIpMoveConfig, routedToServiceName)
	parkConfig := fmt.Sprintf(testAccIpMoveConfig, "")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpMove(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: moveConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", routedToServiceName),
				),
			},
			{
				Config: parkConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", ""),
				),
			},
		},
	})
}

func TestAccIpMove_block(t *testing.T) {
	ipBlock := os.Getenv("OVH_IP_BLOCK_MOVE_TEST")
	routedToServiceName := os.Getenv("OVH_IP_MOVE_SERVICE_NAME_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpMove(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_ip_move" "move" {
						ip = "%s"
						routed_to {
							service_name = "%s"
						}
					}`,
					ipBlock, routedToServiceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", routedToServiceName),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "ovh_ip_move" "move" {
						ip = "%s"
						routed_to {
							service_name = ""
						}
					}`,
					ipBlock),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", ""),
				),
			},
		},
	})
}
