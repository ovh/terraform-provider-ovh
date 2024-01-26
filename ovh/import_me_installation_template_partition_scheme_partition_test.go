package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeInstallationTemplatePartitionSchemePartition_importBasic(t *testing.T) {
	installationTemplate := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeInstallationTemplatePartitionSchemePartitionResourceConfig_basic, installationTemplate),
			},
			{
				ResourceName:      "ovh_me_installation_template_partition_scheme_partition.root",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/myscheme//", installationTemplate),
			},
		},
	})
}
