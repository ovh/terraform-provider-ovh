package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeInstallationTemplatePartitionSchemeResource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatePartitionSchemeResourceConfig_Basic,
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
						"ovh_me_installation_template.template",
						"template_name",
						templateName,
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme.scheme",
						"name",
						"myscheme",
					),
				),
			},
		},
	})
}

func TestAccMeInstallationTemplatePartitionSchemeResource_priority_zero(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatePartitionSchemeResourceConfig_priority_zero,
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
						"ovh_me_installation_template.template",
						"template_name",
						templateName,
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme.scheme",
						"name",
						"myscheme",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme.scheme",
						"priority",
						"0",
					),
				),
			},
		},
	})
}

const testAccMeInstallationTemplatePartitionSchemeResourceConfig_Basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "debian12_64"
  template_name      = "%s"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name      = ovh_me_installation_template.template.template_name
  name               = "myscheme"
  priority           = 1
}
`

const testAccMeInstallationTemplatePartitionSchemeResourceConfig_priority_zero = `
resource "ovh_me_installation_template" "template" {
  base_template_name               = "debian12_64"
  template_name                    = "%s"
  remove_default_partition_schemes = true
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name      = ovh_me_installation_template.template.template_name
  name               = "myscheme"
  priority           = 0
}
`
