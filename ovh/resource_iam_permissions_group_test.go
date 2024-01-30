package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iam_permissions_group", &resource.Sweeper{
		Name: "ovh_iam_policy",
		F:    testSweepIamPermissionsGroup,
	})
}

func testSweepIamPermissionsGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	permissionsGroups := []IamPermissionsGroup{}
	if err := client.Get("/v2/iam/permissionsGroup", &permissionsGroups); err != nil {
		return fmt.Errorf("Error calling /v2/iam/permissionsGroup:\n\t %q", err)
	}

	if len(permissionsGroups) == 0 {
		log.Print("[DEBUG] No iam policy to sweep")
		return nil
	}

	for _, permGrp := range permissionsGroups {
		if !strings.HasPrefix(permGrp.Name, test_prefix) {
			continue
		}

		// skip sweeping permissions groups owned by ovh
		if permGrp.Owner == "ovh" {
			continue
		}

		log.Printf("[DEBUG] IAM policy found %s: %s", permGrp.Name, permGrp.Id)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting iam policy %s: %s", permGrp.Name, permGrp.Id)
			if err := client.Delete(fmt.Sprintf("/v2/iam/permissionsGroup/%s", url.QueryEscape(permGrp.Urn)), nil); err != nil {
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

func TestAccIamPermissionsGroup_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM permissions group created by Terraform Acc"
	allowAction := "account:apiovh:iam/policy/*"
	exceptAction := "account:apiovh:iam/policy/delete"
	denyAction := "account:apiovh:iam/policy/create"
	config := fmt.Sprintf(testAccIamPermissionsGroupConfig, name, desc, allowAction, exceptAction, denyAction)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIamPermissionsGroupResourceAttr("ovh_iam_permissions_group.permissions", name, desc, allowAction, exceptAction, denyAction)...,
				),
			}, {
				ResourceName:      "ovh_iam_permissions_group.permissions",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIamPermissionsGroupConfig = `
resource "ovh_iam_permissions_group" "permissions" {
	name        = "%s"
	description = "%s"
	allow       = ["%s"]
	except      = ["%s"]
	deny        = ["%s"]
}
`
