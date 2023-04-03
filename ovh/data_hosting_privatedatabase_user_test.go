package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceHostingPrivateDatabaseUserConfig = `
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

resource "ovh_hosting_privatedatabase_user" "description" {
    service_name  = ovh_hosting_privatedatabase.database.service_name
    password      = "%s"
    user_name     = "%s"
}

data "ovh_hosting_privatedatabase_user" "description" {
	service_name = ovh_hosting_privatedatabase.database.service_name
	user_name = ovh_hosting_privatedatabase_user.description.user_name
  }
`

func TestAccDataSourceHostingPrivateDatabaseUser_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	displayName := acctest.RandomWithPrefix(test_prefix)
	dc := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	engine := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	password := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")

	config := fmt.Sprintf(
		testAccDataSourceHostingPrivateDatabaseUserConfig,
		desc,
		displayName,
		dc,
		engine,
		password,
		userName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUser(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase_user.description",
						"creation_date",
					),
				),
			},
		},
	})
}
