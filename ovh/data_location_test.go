package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "ovh_location" "paname" {
						name = "eu-west-par"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_location.paname", "code", "par"),
					resource.TestCheckResourceAttr("data.ovh_location.paname", "type", "REGION-3-AZ"),
					resource.TestCheckResourceAttr("data.ovh_location.paname", "city_name", "Paris"),
					resource.TestCheckResourceAttr("data.ovh_location.paname", "country_code", "FR"),
					resource.TestCheckResourceAttrSet("data.ovh_location.paname", "availability_zones.#"),
				),
			},
		},
	})
}
