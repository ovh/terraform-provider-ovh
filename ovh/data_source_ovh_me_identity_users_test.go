package ovh

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIdentityUsersDataSource_basic(t *testing.T) {
	desc := "Identity user created by Terraform Acc."
	email1 := "tf_acceptance_tests_1@example.com"
	email2 := "tf_acceptance_tests_2@example.com"
	group := "DEFAULT"
	login1 := acctest.RandomWithPrefix(test_prefix)
	login2 := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	preSetup := fmt.Sprintf(
		testAccMeIdentityUsersDatasourceConfig_preSetup,
		desc,
		email1,
		group,
		login1,
		password,
		desc,
		email2,
		group,
		login2,
		password,
	)
	config := fmt.Sprintf(
		testAccMeIdentityUsersDatasourceConfig_keys,
		desc,
		email1,
		group,
		login1,
		password,
		desc,
		email2,
		group,
		login2,
		password,
	)

	checks := checkIdentityUserResourceAttr("ovh_me_identity_user.user_1", desc, email1, group, login1)
	checks = append(checks, checkIdentityUserResourceAttr("ovh_me_identity_user.user_1", desc, email1, group, login1)...)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: preSetup,
				Check:  resource.ComposeTestCheckFunc(checks...),
			}, {
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"keys_present", "true"),
				),
			},
		},
	})
}

func checkIdentityUserResourceAttr(name, desc, email, group, login string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "description", desc),
		resource.TestCheckResourceAttr(name, "email", email),
		resource.TestCheckResourceAttr(name, "group", group),
		resource.TestCheckResourceAttr(name, "login", login),
	}
}

const testAccMeIdentityUsersDatasourceConfig_preSetup = `
resource "ovh_me_identity_user" "user_1" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}

resource "ovh_me_identity_user" "user_2" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}
`

const testAccMeIdentityUsersDatasourceConfig_keys = `
resource "ovh_me_identity_user" "user_1" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}

resource "ovh_me_identity_user" "user_2" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}

data "ovh_me_identity_users" "users" {}

output "keys_present" {
    value = tostring(contains(data.ovh_me_identity_users.users.users, ovh_me_identity_user.user_1.login) && contains(data.ovh_me_identity_users.users.users, ovh_me_identity_user.user_2.login))
}
`
