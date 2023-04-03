package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceHostingPrivateDatabaseUserGrantConfig = `
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

resource "ovh_hosting_privatedatabase_user" "user" {
    service_name  = ovh_hosting_privatedatabase.database.service_name
    password      = "%s"
    user_name     = "%s"
}

resource "ovh_hosting_privatedatabase_user_grant" "grant" {
    service_name  = ovh_hosting_privatedatabase.database.service_name
    user_name     = ovh_hosting_privatedatabase_user.user.user_name
    database_name = ovh_hosting_privatedatabase_database.database.database_name
    grant         = "%s"
}

data "ovh_hosting_privatedatabase_user_grant" "description" {
	service_name = ovh_hosting_privatedatabase.database.service_name
	database_name = ovh_hosting_privatedatabase_user_grant.grant.database_name
	user_name = ovh_hosting_privatedatabase_user.user.user_name
}
`

func TestAccDataSourceHostingPrivateDatabaseUserGrant_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	displayName := acctest.RandomWithPrefix(test_prefix)
	dc := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	engine := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	password := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	grant := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_GRANT_TEST")

	config := fmt.Sprintf(
		testAccDataSourceHostingPrivateDatabaseUserGrantConfig,
		desc,
		displayName,
		dc,
		engine,
		databaseName,
		password,
		userName,
		grant,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUserGrant(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_user_grant.description",
						"creation_date",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_hosting_privatedatabase_user_grant.description",
						"grant",
						grant,
					),
				),
			},
		},
	})
}
