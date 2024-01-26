package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: "data ovh_vpss servers {}",
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_vpss.servers",
					"result.#",
				),
			},
		},
	})
}
