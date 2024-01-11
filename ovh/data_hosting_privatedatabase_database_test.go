package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceHostingPrivateDatabaseDatabaseConfig = `
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

resource "ovh_hosting_privatedatabase_database" "database" {
    service_name  = ovh_hosting_privatedatabase.database.service_name
    database_name = "%s"
}

data "ovh_hosting_privatedatabase_database" "database" {
	service_name  = ovh_hosting_privatedatabase.database.service_name
	database_name = ovh_hosting_privatedatabase_database.database.database_name
}
`

func TestAccDataSourceHostingPrivateDatabaseDatabase_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	displayName := acctest.RandomWithPrefix(test_prefix)
	dc := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	engine := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")

	config := fmt.Sprintf(
		testAccDataSourceHostingPrivateDatabaseDatabaseConfig,
		desc,
		displayName,
		dc,
		engine,
		databaseName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_database.database",
						"backup_time",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_database.database",
						"quota_used",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_database.database",
						"creation_date",
					),
				),
			},
		},
	})
}
