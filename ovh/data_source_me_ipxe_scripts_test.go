package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIpxeScriptsDataSource_basic(t *testing.T) {
	scriptName := acctest.RandomWithPrefix(test_prefix)
	presetup := fmt.Sprintf(
		testAccMeIpxeScriptsDatasourceConfig_presetup,
		scriptName,
	)
	config := fmt.Sprintf(
		testAccMeIpxeScriptsDatasourceConfig_Basic,
		scriptName,
		scriptName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: presetup,
				Check: resource.TestCheckResourceAttr(
					"ovh_me_ipxe_script.script",
					"name",
					scriptName,
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_ipxe_scripts.scripts",
						"result.#",
					),
					resource.TestCheckOutput(
						"check",
						"true",
					),
				),
			},
		},
	})
}

const testAccMeIpxeScriptsDatasourceConfig_presetup = `
resource "ovh_me_ipxe_script" "script" {
  name        = "%s"
  script      = "test"
}
`

const testAccMeIpxeScriptsDatasourceConfig_Basic = `
resource "ovh_me_ipxe_script" "script" {
  name        = "%s"
  script      = "test"
}

data "ovh_me_ipxe_scripts" "scripts" {}

output check {
  value = tostring(contains(data.ovh_me_ipxe_scripts.scripts.result, "%s"))
}
`
