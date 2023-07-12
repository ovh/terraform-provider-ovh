package ovh

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_me_identity_user", &resource.Sweeper{
		Name: "ovh_me_identity_user",
		F:    testSweepMeIdentityUser,
	})
}

func testSweepMeIdentityUser(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	names := []string{}
	if err := client.Get("/me/identity/user", &names); err != nil {
		return fmt.Errorf("Error calling /me/identity/user:\n\t %q", err)
	}

	if len(names) == 0 {
		log.Print("[DEBUG] No identity users to sweep")
		return nil
	}

	for _, keyName := range names {
		if !strings.HasPrefix(keyName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Identity user found %v", keyName)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting identity user %v", keyName)
			if err := client.Delete(fmt.Sprintf("/me/identity/user/%s", keyName), nil); err != nil {
				return resource.RetryableError(err)
			}

			// Successful delete
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func TestAccMeIdentityUser_basic(t *testing.T) {
	desc := "Identity user created by Terraform Acc."
	email := "tf_acceptance_tests@example.com"
	group := "DEFAULT"
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))
	config := fmt.Sprintf(testAccMeIdentityUserConfig, desc, email, group, login, password)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityUserResourceAttr("ovh_me_identity_user.user_1", desc, email, group, login)...,
				),
			},
		},
	})
}

func TestAccMeIdentityUser_update(t *testing.T) {
	desc := "Identity user created by Terraform Acc."
	email := "tf_acceptance_tests@example.com"
	group := "DEFAULT"
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	newEmail := "tf_acceptance_tests_new@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityUserConfig, desc, email, group, login, password),
				Check: resource.ComposeTestCheckFunc(
					checkIdentityUserResourceAttr("ovh_me_identity_user.user_1", desc, email, group, login)...,
				),
			},
			{
				Config: fmt.Sprintf(testAccMeIdentityUserConfig, desc, newEmail, group, login, password),
				Check: resource.ComposeTestCheckFunc(
					checkIdentityUserResourceAttr("ovh_me_identity_user.user_1", desc, newEmail, group, login)...,
				),
			}},
	})
}

const testAccMeIdentityUserConfig = `
resource "ovh_me_identity_user" "user_1" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}
`
