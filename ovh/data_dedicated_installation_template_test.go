package ovh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstallationTemplateDataSource_basic(t *testing.T) {
	templateName := "debian12_64"
	config := fmt.Sprintf(
		testAccInstallationTemplateDatasourceConfig_Basic,
		templateName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_installation_template.template",
						"template_name",
						templateName,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_installation_template.template",
						"family",
						"linux",
					),
				),
			},
			{
				Config:      testAccInstallationTemplateDatasourceConfig_404,
				ExpectError: regexp.MustCompile("Your query returned no results. Please change your search criteria"),
			},
		},
	})
}

const testAccInstallationTemplateDatasourceConfig_404 = `
data "ovh_dedicated_installation_template" "notemplate" {
	template_name = "42"
  }
`
const testAccInstallationTemplateDatasourceConfig_Basic = `
data "ovh_dedicated_installation_template" "template" {
    template_name      = "%s"
}
`
