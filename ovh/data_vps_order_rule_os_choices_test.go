package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSOrderRuleOSChoicesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSOrderRuleOSChoicesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_order_rule_os_choices.test", "choices.#"),
				),
			},
		},
	})
}

const testAccVPSOrderRuleOSChoicesDatasourceConfig = `
data "ovh_vps_order_rule_os_choices" "test" {
  datacenter = "GRA"
  os         = "linux"
}
`
