package ovh

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func init() {
	resource.AddTestSweepers("ovh_me_identity_group_membership", &resource.Sweeper{
		Name:         "ovh_me_identity_group_membership",
		Dependencies: []string{"ovh_me_identity_user"},
		F:            testSweepMeIdentityGroupMembership,
	})
}

func testSweepMeIdentityGroupMembership(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	var groups []string
	if err := client.Get("/me/identity/group", &groups); err != nil {
		return fmt.Errorf("Error calling /me/identity/group:\n\t %q", err)
	}

	for _, groupName := range groups {
		if !strings.HasPrefix(groupName, test_prefix) {
			continue
		}

		var users []string
		endpoint := fmt.Sprintf("/me/identity/group/%s/user", url.PathEscape(groupName))
		if err := client.Get(endpoint, &users); err != nil {
			log.Printf("[WARN] Could not list users for group %s: %s", groupName, err)
			continue
		}

		for _, login := range users {
			log.Printf("[INFO] Removing user %s from group %s (sweeper)", login, groupName)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				deleteEndpoint := fmt.Sprintf("/me/identity/group/%s/user/%s",
					url.PathEscape(groupName), url.PathEscape(login))
				if err := client.Delete(deleteEndpoint, nil); err != nil {
					return resource.RetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccMeIdentityGroupMembership_basic(t *testing.T) {
	login := acctest.RandomWithPrefix(test_prefix)
	groupName := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityGroupMembershipConfig_basic, groupName, login, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_me_identity_group_membership.membership_1", "login", login),
					resource.TestCheckResourceAttr(
						"ovh_me_identity_group_membership.membership_1", "group", groupName),
				),
			},
		},
	})
}

func TestAccMeIdentityGroupMembership_forceNew(t *testing.T) {
	login := acctest.RandomWithPrefix(test_prefix)
	groupName1 := acctest.RandomWithPrefix(test_prefix)
	groupName2 := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityGroupMembershipConfig_twoGroups, groupName1, groupName2, login, password, "group_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_group_membership.membership_1", "group", groupName1),
				),
			},
			{
				// Changing the group must replace the membership
				Config: fmt.Sprintf(testAccMeIdentityGroupMembershipConfig_twoGroups, groupName1, groupName2, login, password, "group_2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_me_identity_group_membership.membership_1", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_group_membership.membership_1", "group", groupName2),
					resource.TestCheckResourceAttr("ovh_me_identity_group_membership.membership_1", "id", fmt.Sprintf("%s/%s", login, groupName2)),
				),
			},
			{
				// The user refreshes its memberships: only the new group remains
				Config: fmt.Sprintf(testAccMeIdentityGroupMembershipConfig_twoGroups, groupName1, groupName2, login, password, "group_2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_me_identity_user.user_1", "groups.#", "1"),
					resource.TestCheckTypeSetElemAttr("ovh_me_identity_user.user_1", "groups.*", groupName2),
				),
			},
		},
	})
}

func TestAccMeIdentityGroupMembership_importBasic(t *testing.T) {
	login := acctest.RandomWithPrefix(test_prefix)
	groupName := acctest.RandomWithPrefix(test_prefix)
	password := base64.StdEncoding.EncodeToString([]byte(acctest.RandomWithPrefix(test_prefix)))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMeIdentityGroupMembershipConfig_basic, groupName, login, password),
			},
			{
				ResourceName:      "ovh_me_identity_group_membership.membership_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%s", login, groupName),
			},
		},
	})
}

const testAccMeIdentityGroupMembershipConfig_basic = `
resource "ovh_me_identity_group" "group_1" {
	description = "Test group for membership resource"
	name        = "%s"
	role        = "NONE"
}

resource "ovh_me_identity_user" "user_1" {
	description = "Test user for membership resource"
	email       = "tf_acceptance_tests@example.com"
	group       = "DEFAULT"
	login       = "%s"
	password    = "%s"
}

resource "ovh_me_identity_group_membership" "membership_1" {
	login = ovh_me_identity_user.user_1.login
	group = ovh_me_identity_group.group_1.name
}
`

const testAccMeIdentityGroupMembershipConfig_twoGroups = `
resource "ovh_me_identity_group" "group_1" {
    description = "Test group 1 for membership resource"
    name        = "%s"
    role        = "NONE"
}

resource "ovh_me_identity_group" "group_2" {
    description = "Test group 2 for membership resource"
    name        = "%s"
    role        = "NONE"
}

resource "ovh_me_identity_user" "user_1" {
    description = "Test user for membership resource"
    email       = "tf_acceptance_tests@example.com"
    group       = "DEFAULT"
    login       = "%s"
    password    = "%s"
}

resource "ovh_me_identity_group_membership" "membership_1" {
    login = ovh_me_identity_user.user_1.login
    group = ovh_me_identity_group.%s.name
}
`
