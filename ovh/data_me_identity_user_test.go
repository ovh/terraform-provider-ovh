package ovh

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeIdentityUserDataSource_basic(t *testing.T) {
	desc := "Identity user created by Terraform Acc."
	email := "tf_acceptance_tests@example.com"
	group := "DEFAULT"
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))
	config := fmt.Sprintf(testAccMeIdentityUserDatasourceConfig, desc, email, group, login, password)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityUserResourceAttr("data.ovh_me_identity_user.user_1", desc, email, group, login)...,
				),
			},
		},
	})
}

const testAccMeIdentityUserDatasourceConfig = `
resource "ovh_me_identity_user" "user_1" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}

data "ovh_me_identity_user" "user_1" {
  user       = ovh_me_identity_user.user_1.login
  depends_on = [ovh_me_identity_user.user_1]
}
`
