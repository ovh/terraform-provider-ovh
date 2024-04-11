package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeInstallationTemplatePartitionSchemeHardwareRaidResource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatePartitionSchemeHardwareRaidResourceConfig_basic,
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
						"ovh_me_installation_template_partition_scheme_hardware_raid.group1",
						"name",
						"group1",
					),
					resource.TestCheckResourceAttr(
						"ovh_me_installation_template_partition_scheme_hardware_raid.group2",
						"name",
						"group2",
					),
				),
			},
		},
	})
}

const testAccMeInstallationTemplatePartitionSchemeHardwareRaidResourceConfig_basic = `
resource "ovh_me_installation_template" "template" {
  base_template_name = "debian12_64"
  template_name      = "%s"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name      = ovh_me_installation_template.template.template_name
  name               = "myscheme"
  priority           = 1
}

resource "ovh_me_installation_template_partition_scheme_hardware_raid" "group1" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  name          = "group1"
  disks         = ["[c1:d1,c1:d2,c1:d3]", "[c1:d10,c1:d20,c1:d30]"]
  mode          = "raid50"
  step          = 1
}

resource "ovh_me_installation_template_partition_scheme_hardware_raid" "group2" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  name          = "group2"
  disks         = ["c2:d11", "c2:d21"]
  mode          = "raid1"
  step          = 2
}
`
