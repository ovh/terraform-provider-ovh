package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeInstallationTemplate_importBasic(t *testing.T) {
	installationTemplate := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeInstallationTemplateResourceConfig_Basic, installationTemplate),
			},
			{
				ResourceName:      "ovh_me_installation_template.template",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("centos7_64/%s", installationTemplate),
			},
		},
	})
}
