package ovh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeInstallationTemplateDataSource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplateDatasourceConfig_Basic,
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
						"data.ovh_me_installation_template.template",
						"template_name",
						templateName,
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_template.template",
						"partition_scheme.0.partition.0.type",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_template.template",
						"partition_scheme.0.partition.0.order",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_template.template",
						"partition_scheme.0.partition.0.mountpoint",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_template.template",
						"family",
					),
				),
			},
			{
				Config:      testAccMeInstallationTemplateDatasourceConfig_404,
				ExpectError: regexp.MustCompile("Your query returned no results. Please change your search criteria"),
			},
		},
	})
}

const testAccMeInstallationTemplateDatasourceConfig_404 = `
data "ovh_me_installation_template" "template" {
	template_name = "42"
  }
`
const testAccMeInstallationTemplateDatasourceConfig_Basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "debian12_64"
  template_name      = "%s"
}

data "ovh_me_installation_template" "template" {
  template_name = ovh_me_installation_template.template.template_name
  depends_on    = [ovh_me_installation_template.template]
}
`
