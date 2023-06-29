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
	resource.AddTestSweepers("ovh_iam_policy", &resource.Sweeper{
		Name: "ovh_iam_policy",
		F:    testSweepIamPolicy,
	})
}

func testSweepIamPolicy(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	policyIds := []string{}
	if err := client.Get("/v2/iam/policy", &policyIds); err != nil {
		return fmt.Errorf("Error calling /v2/iam/policy:\n\t %q", err)
	}

	if len(policyIds) == 0 {
		log.Print("[DEBUG] No identity groups to sweep")
		return nil
	}

	for _, polId := range policyIds {
		var polDetails IamPolicy
		if err := client.Get(fmt.Sprintf("/v2/iam/policy/%s", polId), &polDetails); err != nil {
			return err
		}

		if !strings.HasPrefix(polDetails.Name, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] IAM policy found %s: %s", polDetails.Name, polId)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting iam policy %s: %s", polDetails.Name, polId)
			if err := client.Delete(fmt.Sprintf("/v2/iam/policy/%s", polId), nil); err != nil {
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

func TestAccIamPolicy_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	config := fmt.Sprintf(testAccIamPolicyConfig, userName, userName, name, desc, res)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIamPolicyResourceAttr("ovh_iam_policy.policy1", name, desc, res)...,
				),
			},
		},
	})
}

const testAccIamPolicyConfig = `
resource "ovh_me_identity_user" "test_user" {
	login = "%s"
	email = "%s@terraform.test"
	password = "qwe123!@#"
}

resource "ovh_iam_policy" "policy1" {
	name        = "%s"
	description = "%s"
	identities  = [ovh_me_identity_user.test_user.urn]
	resources   = ["%s"]
	allow 	    = ["*"]
}
`
