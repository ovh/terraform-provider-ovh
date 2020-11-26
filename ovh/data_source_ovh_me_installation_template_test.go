package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		},
	})
}

const testAccMeInstallationTemplateDatasourceConfig_Basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}

data "ovh_me_installation_template" "template" {
  template_name = ovh_me_installation_template.template.template_name
  depends_on    = [ovh_me_installation_template.template]
}
`
