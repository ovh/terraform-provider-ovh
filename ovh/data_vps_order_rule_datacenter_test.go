package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSOrderRuleDatacenterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSOrderRuleDatacenterDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_order_rule_datacenter.test", "datacenters.#"),
				),
			},
		},
	})
}

const testAccVPSOrderRuleDatacenterDatasourceConfig = `
data "ovh_vps_order_rule_datacenter" "test" {
  ovh_subsidiary = "US"
  plan_code      = "vps-2025-model1"
}
`
