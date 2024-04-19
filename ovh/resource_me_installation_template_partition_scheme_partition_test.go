package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeInstallationTemplatePartitionSchemePartitionResource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatePartitionSchemePartitionResourceConfig_basic,
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
						"ovh_me_installation_template_partition_scheme_partition.root",
						"size",
						"400",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme_partition.root",
						"type",
						"primary",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme_partition.home",
						"size",
						"500",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme_partition.home",
						"type",
						"logical",
					),
				),
			},
		},
	})
}

const testAccMeInstallationTemplatePartitionSchemePartitionResourceConfig_basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "debian12_64"
  template_name      = "%s"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name      = ovh_me_installation_template.template.template_name
  name               = "myscheme"
  priority           = 1
}

resource "ovh_me_installation_template_partition_scheme_partition" "root" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  mountpoint    = "/"
  filesystem    = "ext4"
  size          = "400"
  order         = 1
  type          = "primary"
}

resource "ovh_me_installation_template_partition_scheme_partition" "home" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  mountpoint    = "/home"
  filesystem    = "ext4"
  size          = "500"
  order         = 2
  type          = "logical"
}
`
