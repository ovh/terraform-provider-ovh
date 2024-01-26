package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeIdentityGroup_importBasic(t *testing.T) {
	resourceName := "ovh_me_identity_group.group_1"
	groupName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityGroupConfig_import, groupName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMeIdentityGroupConfig_import = `
resource "ovh_me_identity_group" "group_1" {
	description = "tf acc import test"
  	name        = "%s"
  	role        = "ADMIN"
}
`
