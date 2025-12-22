package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpLoadbalancingsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data ovh_iploadbalancing_nat_ips iplbs {
					}`,
				Check: resource.TestCheckResourceAttrSet("data.ovh_iploadbalancings.iplbs", "iploadbalancings.#"),
			},
		},
	})
}
