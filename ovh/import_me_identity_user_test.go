package ovh

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIdentityUser_importBasic(t *testing.T) {
	resourceName := "ovh_me_identity_user.user_1"
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityUserConfig_import, login, password),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

const testAccMeIdentityUserConfig_import = `
resource "ovh_me_identity_user" "user_1" {
	description = "tf acc import test"
  	email       = "tf_import@example.com"
  	group       = "DEFAULT"
  	login       = "%s"
  	password    = "%s"
}
`
