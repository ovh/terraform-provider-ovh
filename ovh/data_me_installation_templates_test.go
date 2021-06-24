package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeInstallationTemplatesDataSource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	presetup := fmt.Sprintf(
		testAccMeInstallationTemplatesDatasourceConfig_presetup,
		templateName,
	)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatesDatasourceConfig_Basic,
		templateName,
		templateName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: presetup,
				Check: resource.TestCheckResourceAttr(
					"ovh_me_installation_template.template",
					"template_name",
					templateName,
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_templates.templates",
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

const testAccMeInstallationTemplatesDatasourceConfig_presetup = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}
`

const testAccMeInstallationTemplatesDatasourceConfig_Basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}

data "ovh_me_installation_templates" "templates" {}

output check {
  value = tostring(contains(data.ovh_me_installation_templates.templates.result, "%s"))
}
`
