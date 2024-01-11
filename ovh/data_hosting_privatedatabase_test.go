package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceHostingPrivateDatabaseConfig = `
data "ovh_order_cart" "mycart" {
	ovh_subsidiary = "fr"
	description    = "%s"
  }
	
  data "ovh_order_cart_product_plan" "database" {
	cart_id        = data.ovh_order_cart.mycart.id
	price_capacity = "renew"
	product        = "privateSQL"
	plan_code      = "private-sql-512-instance"
  }
	
  resource "ovh_hosting_privatedatabase" "database" {
	ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
	display_name   = "%s"
  
	plan {
	  duration     = "P1M"
	  plan_code    = data.ovh_order_cart_product_plan.database.plan_code
	  pricing_mode = data.ovh_order_cart_product_plan.database.selected_price[0].pricing_mode
  
	  configuration {
		label = "dc"
		value = "%s"
	  }
  
	  configuration {
		label = "engine"
		value = "%s"
	  }
	}
  }

data "ovh_hosting_privatedatabase" "database" {
  service_name = ovh_hosting_privatedatabase.database.service_name
}
`

func TestAccDataSourceHostingPrivateDatabase_basic(t *testing.T) {

	desc := acctest.RandomWithPrefix(test_prefix)
	displayName := acctest.RandomWithPrefix(test_prefix)
	dc := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	engine := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")

	config := fmt.Sprintf(
		testAccDataSourceHostingPrivateDatabaseConfig,
		desc,
		displayName,
		dc,
		engine,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"cpu",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"datacenter",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"offer",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"type",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"version_number",
					),
				),
			},
		},
	})
}
