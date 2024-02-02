package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

var testAccVrackIpConfig = `
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
 description    = data.ovh_order_cart.mycart.description
 name           = data.ovh_order_cart.mycart.description
 ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary

 plan {
   duration     = data.ovh_order_cart_product_plan.vrack.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.vrack.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.vrack.selected_price.0.pricing_mode
 }
}

data "ovh_order_cart_product_plan" "ipblock" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ip"
  plan_code      = "ip-v4-s30-ripe"
}

resource "ovh_ip_service" "ipblock" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description    = data.ovh_order_cart.mycart.description

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

resource "ovh_vrack_ip" "vrackblock" {
  service_name = ovh_vrack.vrack.service_name
  block        = ovh_ip_service.ipblock.ip
}
`

func init() {
	resource.AddTestSweepers("ovh_vrack_ip", &resource.Sweeper{
		Name: "ovh_vrack_ip",
		F:    testSweepVrackIp,
	})
}

func testSweepVrackIp(region string) error {
	config, err := sharedConfigForRegion(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	res := []string{}
	endpoint := "/vrack"

	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_vrack to sweep")
		return nil
	}

	for _, vrackId := range res {
		log.Printf("[DEBUG] Will read vrack %s", vrackId)

		vrack := &Vrack{}
		endpoint := fmt.Sprintf(
			"/vrack/%s",
			url.PathEscape(vrackId),
		)
		if err := config.OVHClient.Get(endpoint, vrack); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(*vrack.Description, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Will read vrack ip attach for %s", vrackId)

		ips := []string{}
		endpoint = fmt.Sprintf(
			"/vrack/%s/ip",
			url.PathEscape(vrackId),
		)

		for _, ip := range ips {
			ipService := &IpService{}

			endpoint := fmt.Sprintf(
				"/ip/service/%s",
				url.PathEscape(ip),
			)
			if err := config.OVHClient.Get(endpoint, ipService); err != nil {
				return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
			}

			if !strings.HasPrefix(*ipService.Description, test_prefix) {
				continue
			}

			endpoint = fmt.Sprintf("/vrack/%s/ip/%s",
				url.PathEscape(vrackId),
				url.PathEscape(ip),
			)

			vrackblock := &VrackIp{}

			if err := config.OVHClient.Get(endpoint, vrackblock); err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					continue
				}
				return err
			}

			task := &VrackTask{}

			if err := config.OVHClient.Delete(endpoint, task); err != nil {
				return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, ip, err)
			}

			if err := waitForVrackTask(task, config.OVHClient); err != nil {
				return fmt.Errorf("Error waiting for vrack (%s) to detach cloud project (%s): %s", vrackId, ip, err)
			}

		}
	}

	return nil
}

func TestAccVrackIp_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccVrackIpConfig,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackIpPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_vrack_ip.vrackblock", "service_name"),
					resource.TestCheckResourceAttrSet("ovh_vrack_ip.vrackblock", "block"),
				),
			},
		},
	})
}

func testAccCheckVrackIpPreCheck(t *testing.T) {
	testAccPreCheckOrderVrack(t)
	testAccPreCheckOrderIpService(t)
}
