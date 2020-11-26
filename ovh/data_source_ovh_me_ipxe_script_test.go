package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIpxeScriptDataSource_basic(t *testing.T) {
	scriptName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(testAccMeIpxeScriptDatasourceConfig, scriptName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_me_ipxe_script.script", "name", scriptName),
					resource.TestCheckResourceAttr(
						"data.ovh_me_ipxe_script.script", "script", "test"),
				),
			},
		},
	})
}

const testAccMeIpxeScriptDatasourceConfig = `
resource "ovh_me_ipxe_script" "script" {
  name        = "%s"
  script      = "test"
}

data "ovh_me_ipxe_script" "script" {
  name = ovh_me_ipxe_script.script.name
  depends_on    = [ovh_me_ipxe_script.script]
}
`
