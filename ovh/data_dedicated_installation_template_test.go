package ovh

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedInstallationTemplateDataSource_basic(t *testing.T) {
	testAccInstallationTemplateDatasourceConfig_404 := `data "ovh_dedicated_installation_template" "notemplate" {
		template_name = "42"
	}`
	testAccInstallationTemplateDatasourceConfig_Basic := `data "ovh_dedicated_installation_template" "template" {
		template_name= "debian12_64"
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstallationTemplateDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_installation_template.template",
						"template_name",
						"debian12_64",
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
