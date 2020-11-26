package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: "data ovh_dedicated_servers servers {}",
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_dedicated_servers.servers",
					"result.#",
				),
			},
		},
	})
}
