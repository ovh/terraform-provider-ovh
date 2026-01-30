package ovh

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeIdentityUserToken_basic(t *testing.T) {
	login := acctest.RandomWithPrefix("tf_user")
	tokenName := acctest.RandomWithPrefix("tf_token")
	tokenDesc := "Token created by Terraform acceptance test"

	email := "tf_acc_token@example.com"
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix("pass")))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMeIdentityUserTokenConfig_basic(login, email, password, tokenName, tokenDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_basic", "user_login", login),
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_basic", "name", tokenName),
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_basic", "description", tokenDesc),
					resource.TestCheckResourceAttrSet("ovh_me_identity_user_token.token_basic", "token"),
					resource.TestCheckResourceAttrSet("ovh_me_identity_user_token.token_basic", "creation"),
				),
			},
		},
	})
}

func TestAccMeIdentityUserToken_expiresIn(t *testing.T) {
	login := acctest.RandomWithPrefix("tf_user_in")
	tokenName := acctest.RandomWithPrefix("tf_token_in")
	tokenDesc := "Token with expires_in"

	email := "tf_acc_token_in@example.com"
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix("pass")))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMeIdentityUserTokenConfig_expiresIn(login, email, password, tokenName, tokenDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_expires_in", "expires_in", "3600"),
					resource.TestCheckResourceAttrSet("ovh_me_identity_user_token.token_expires_in", "expires_at"),
				),
			},
		},
	})
}

func TestAccMeIdentityUserToken_expiresAt(t *testing.T) {
	login := acctest.RandomWithPrefix("tf_user_at")
	tokenName := acctest.RandomWithPrefix("tf_token_at")
	tokenDesc := "Token with expires_at"

	email := "tf_acc_token_at@example.com"
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix("pass")))

	// Create a future timestamp
	tomorrow := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMeIdentityUserTokenConfig_expiresAt(login, email, password, tokenName, tokenDesc, tomorrow),
				Check: resource.ComposeTestCheckFunc(
					// We check if it's set. Exact match might fail due to normalization (Z vs offset)
					// unless we normalize it ourselves. For basic check, ensuring it's set is good.
					resource.TestCheckResourceAttrSet("ovh_me_identity_user_token.token_expires_at", "expires_at"),
				),
			},
		},
	})
}

func TestAccMeIdentityUserToken_update(t *testing.T) {
	login := acctest.RandomWithPrefix("tf_user_upd")
	tokenName := acctest.RandomWithPrefix("tf_token_upd")
	tokenDesc := "Token initial"
	tokenDescUpdated := "Token updated"

	email := "tf_acc_token_upd@example.com"
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix("pass")))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMeIdentityUserTokenConfig_update(login, email, password, tokenName, tokenDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_update", "description", tokenDesc),
				),
			},
			{
				Config: testAccMeIdentityUserTokenConfig_update(login, email, password, tokenName, tokenDescUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user_token.token_update", "description", tokenDescUpdated),
				),
			},
		},
	})
}

func testAccMeIdentityUserTokenConfig_basic(login, email, password, tokenName, tokenDesc string) string {
	return fmt.Sprintf(`
resource "ovh_me_identity_user" "user_1" {
	description = "User for token basic"
	email       = "%s"
	group       = "DEFAULT"
	login       = "%s"
	password    = "%s"
}

resource "ovh_me_identity_user_token" "token_basic" {
	user_login  = ovh_me_identity_user.user_1.login
	name        = "%s"
	description = "%s"
}
`, email, login, password, tokenName, tokenDesc)
}

func testAccMeIdentityUserTokenConfig_expiresIn(login, email, password, tokenName, tokenDesc string) string {
	return fmt.Sprintf(`
resource "ovh_me_identity_user" "user_2" {
	description = "User for token expires_in"
	email       = "%s"
	group       = "DEFAULT"
	login       = "%s"
	password    = "%s"
}

resource "ovh_me_identity_user_token" "token_expires_in" {
	user_login  = ovh_me_identity_user.user_2.login
	name        = "%s"
	description = "%s"
	expires_in  = 3600
}
`, email, login, password, tokenName, tokenDesc)
}

func testAccMeIdentityUserTokenConfig_expiresAt(login, email, password, tokenName, tokenDesc, expiresAt string) string {
	return fmt.Sprintf(`
resource "ovh_me_identity_user" "user_3" {
	description = "User for token expires_at"
	email       = "%s"
	group       = "DEFAULT"
	login       = "%s"
	password    = "%s"
}

resource "ovh_me_identity_user_token" "token_expires_at" {
	user_login  = ovh_me_identity_user.user_3.login
	name        = "%s"
	description = "%s"
	expires_at  = "%s"
}
`, email, login, password, tokenName, tokenDesc, expiresAt)
}

func testAccMeIdentityUserTokenConfig_update(login, email, password, tokenName, tokenDesc string) string {
	return fmt.Sprintf(`
resource "ovh_me_identity_user" "user_4" {
	description = "User for token update"
	email       = "%s"
	group       = "DEFAULT"
	login       = "%s"
	password    = "%s"
}

resource "ovh_me_identity_user_token" "token_update" {
	user_login  = ovh_me_identity_user.user_4.login
	name        = "%s"
	description = "%s"
}
`, email, login, password, tokenName, tokenDesc)
}
