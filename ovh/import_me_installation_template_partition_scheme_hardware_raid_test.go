package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeInstallationTemplatePartitionSchemeHardwareRaid_importBasic(t *testing.T) {
	installationTemplate := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeInstallationTemplatePartitionSchemeHardwareRaidResourceConfig_basic, installationTemplate),
			},
			{
				ResourceName:      "ovh_me_installation_template_partition_scheme_hardware_raid.group1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/myscheme/group1", installationTemplate),
			},
		},
	})
}
