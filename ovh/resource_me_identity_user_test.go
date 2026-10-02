package ovh

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	// Collect group memberships once, as the API only lists users per group
	var groups []string
	if err := client.Get("/me/identity/group", &groups); err != nil {
		return fmt.Errorf("Error calling /me/identity/group:\n\t %q", err)
	}
	userGroups := map[string][]string{}
	for _, groupName := range groups {
		var users []string
		if err := client.Get(fmt.Sprintf("/me/identity/group/%s/user", url.PathEscape(groupName)), &users); err != nil {
			log.Printf("[WARN] Could not list users for group %s: %s", groupName, err)
			continue
		}
		for _, u := range users {
			if strings.HasPrefix(u, test_prefix) {
				userGroups[u] = append(userGroups[u], groupName)
			}
		}
	}

	for _, keyName := range names {
		if !strings.HasPrefix(keyName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Identity user found %v", keyName)
		for _, groupName := range userGroups[keyName] {
			log.Printf("[INFO] Removing user %s from group %s", keyName, groupName)
			if err := client.Delete(fmt.Sprintf("/me/identity/group/%s/user/%s", url.PathEscape(groupName), url.PathEscape(keyName)), nil); err != nil {
				log.Printf("[WARN] Could not remove user %s from group %s: %s", keyName, groupName, err)
			}
		}

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

func TestAccMeIdentityUser_withGroups(t *testing.T) {
	desc := "Identity user created by Terraform Acc."
	email := "tf_acceptance_tests@example.com"
	group := "DEFAULT"
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))
	groupName1 := acctest.RandomWithPrefix(test_prefix)
	groupName2 := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Create user with one additional group
				Config: fmt.Sprintf(testAccMeIdentityUserWithOneGroupConfig, groupName1, groupName2, desc, email, group, login, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "description", desc),
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "email", email),
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "group", group),
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "login", login),
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "groups.#", "1"),
				),
			},
			{
				// Step 2: Update to two additional groups
				Config: fmt.Sprintf(testAccMeIdentityUserWithTwoGroupsConfig, groupName1, groupName2, desc, email, group, login, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "groups.#", "2"),
				),
			},
			{
				// Step 3: Remove all additional groups
				Config: fmt.Sprintf(testAccMeIdentityUserWithNoGroupsConfig, groupName1, groupName2, desc, email, group, login, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccMeIdentityUser_mainGroupInGroups(t *testing.T) {
	login := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccMeIdentityUserMainGroupInGroupsConfig, login, password),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`groups must not contain the main group "DEFAULT"`),
			},
		},
	})
}

const testAccMeIdentityUserMainGroupInGroupsConfig = `
resource "ovh_me_identity_user" "user_1" {
	email    = "tf_acceptance_tests@example.com"
	group    = "DEFAULT"
	login    = "%s"
	password = "%s"
	groups   = ["DEFAULT"]
}
`

const testAccMeIdentityUserConfig = `
resource "ovh_me_identity_user" "user_1" {
	description = "%s"
  	email       = "%s"
  	group       = "%s"
  	login       = "%s"
  	password    = "%s"
}
`

const testAccMeIdentityUserWithOneGroupConfig = `
resource "ovh_me_identity_group" "group_1" {
	description = "test group 1"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_group" "group_2" {
	description = "test group 2"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_user" "user_1" {
	description = "%s"
	email       = "%s"
	group       = "%s"
	login       = "%s"
	password    = "%s"
	groups      = [ovh_me_identity_group.group_1.name]
}
`

const testAccMeIdentityUserWithTwoGroupsConfig = `
resource "ovh_me_identity_group" "group_1" {
	description = "test group 1"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_group" "group_2" {
	description = "test group 2"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_user" "user_1" {
	description = "%s"
	email       = "%s"
	group       = "%s"
	login       = "%s"
	password    = "%s"
	groups      = [ovh_me_identity_group.group_1.name, ovh_me_identity_group.group_2.name]
}
`

const testAccMeIdentityUserWithNoGroupsConfig = `
resource "ovh_me_identity_group" "group_1" {
	description = "test group 1"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_group" "group_2" {
	description = "test group 2"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_user" "user_1" {
	description = "%s"
	email       = "%s"
	group       = "%s"
	login       = "%s"
	password    = "%s"
	groups      = []
}
`
