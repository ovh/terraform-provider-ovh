package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccHostingPrivateDatabaseWhitelistConfig = `
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

resource "ovh_hosting_privatedatabase_whitelist" "whitelist" {
    service_name = ovh_hosting_privatedatabase.database.service_name
    ip           = "%s"
    name         = "%s"
    sftp         = "%s"
    service      = "%s"
}

data "ovh_hosting_privatedatabase_whitelist" "description" {
	service_name = ovh_hosting_privatedatabase.database.service_name
	ip = ovh_hosting_privatedatabase_whitelist.whitelist.ip
}
`

func TestAccDataSourceHostingPrivateDatabaseWhitelist_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	displayName := acctest.RandomWithPrefix(test_prefix)
	dc := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	engine := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	ip := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST")
	name := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_NAME_TEST")
	sftp := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SFTP_TEST")
	service := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SERVICE_TEST")

	config := fmt.Sprintf(
		testAccHostingPrivateDatabaseWhitelistConfig,
		desc,
		displayName,
		dc,
		engine,
		ip,
		name,
		sftp,
		service,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseWhitelist(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_whitelist.description",
						"creation_date",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_whitelist.description",
						"last_update",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_whitelist.description",
						"name",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_whitelist.description",
						"service",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_whitelist.description",
						"status",
					),
				),
			},
		},
	})
}
