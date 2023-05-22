package ovh

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_me_identity_group", &resource.Sweeper{
		Name: "ovh_me_identity_group",
		F:    testSweepMeIdentityGroup,
	})
}

func testSweepMeIdentityGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	names := []string{}
	if err := client.Get("/me/identity/group", &names); err != nil {
		return fmt.Errorf("Error calling /me/identity/group:\n\t %q", err)
	}

	if len(names) == 0 {
		log.Print("[DEBUG] No identity groups to sweep")
		return nil
	}

	for _, keyName := range names {
		if !strings.HasPrefix(keyName, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Identity group found %v", keyName)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting identity group %v", keyName)
			if err := client.Delete(fmt.Sprintf("/me/identity/group/%s", keyName), nil); err != nil {
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

func TestAccMeIdentityGroup_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "Identity group created by Terraform Acc."
	role := "NONE"
	config := fmt.Sprintf(testAccMeIdentityGroupConfig, desc, name, role)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityGroupResourceAttr("ovh_me_identity_group.group_1", name, desc, role)...,
				),
			},
		},
	})
}

const testAccMeIdentityGroupConfig = `
resource "ovh_me_identity_group" "group_1" {
	description = "%s"
  	name        = "%s"
  	role        = "%s"
}
`
