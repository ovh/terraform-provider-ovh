package ovh

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	policies := []IamPolicy{}
	if err := client.Get("/v2/iam/policy", &policies); err != nil {
		return fmt.Errorf("Error calling /v2/iam/policy:\n\t %q", err)
	}

	if len(policies) == 0 {
		log.Print("[DEBUG] No iam policy to sweep")
		return nil
	}

	for _, pol := range policies {
		if !strings.HasPrefix(pol.Name, test_prefix) {
			continue
		}

		// skip seeping readonly attributes
		if pol.ReadOnly {
			continue
		}

		log.Printf("[DEBUG] IAM policy found %s: %s", pol.Name, pol.Id)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting iam policy %s: %s", pol.Name, pol.Id)
			if err := client.Delete(fmt.Sprintf("/v2/iam/policy/%s", pol.Id), nil); err != nil {
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
	exceptAction := "vps:apiovh:reinstall"
	config := fmt.Sprintf(testAccIamPolicyConfig, userName, userName, name, desc, res, exceptAction)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIamPolicyResourceAttr("ovh_iam_policy.policy1", name, desc, res, "*", exceptAction, "")...,
				),
			},
		},
	})
}

func TestAccIamPolicy_deny(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	denyAction := "vps:apiovh:reinstall"
	config := fmt.Sprintf(testAccIamPolicyDenyConfig, userName, userName, name, desc, res, denyAction)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIamPolicyResourceAttr("ovh_iam_policy.policy1", name, desc, res, "", "", denyAction)...,
				),
			},
		},
	})
}

func TestAccIamPolicy_withConditions(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with conditions created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	config := fmt.Sprintf(testAccIamPolicyConditionsConfig, userName, userName, name, desc, res)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "name", name),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "description", desc),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.operator", "OR"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.values.resource.Tag(environment)", "production"),
				),
			},
		},
	})
}

func TestAccIamPolicy_withOrConditions(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with OR conditions created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	config := fmt.Sprintf(testAccIamPolicyOrConditionsConfig, userName, userName, name, desc, res)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "name", name),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "description", desc),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.operator", "OR"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.values.resource.Tag(environment)", "production"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.1.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.1.values.resource.Tag(team)", "platform"),
				),
			},
		},
	})
}

func TestAccIamPolicy_withExpiration(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with expiration date created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	expiration := "2025-12-31T23:59:59Z"
	config := fmt.Sprintf(testAccIamPolicyWithExpirationConfig, userName, userName, name, desc, res, expiration)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "name", name),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "description", desc),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "expired_at", expiration),
				),
			},
		},
	})
}

func TestAccIamPolicy_withExpirationAndConditions(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with expiration and conditions created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	expiration := "2025-12-31T23:59:59Z"
	config := fmt.Sprintf(testAccIamPolicyWithExpirationAndConditionsConfig, userName, userName, name, desc, res, expiration)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "name", name),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "description", desc),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "expired_at", expiration),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.values.resource.Tag(Environment)", "development"),
				),
			},
		},
	})
}

func TestAccIamPolicy_withMaxDepthConditions(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with maximum depth conditions (3 levels) created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	config := fmt.Sprintf(testAccIamPolicyMaxDepthConditionsConfig, userName, userName, name, desc, res)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "name", name),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "description", desc),
					// Level 1: AND operator
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.operator", "AND"),
					// Level 2: First OR condition
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.operator", "OR"),
					// Level 3: First MATCH under first OR
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.condition.0.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.condition.0.values.resource.Tag(environment)", "production"),
					// Level 3: Second MATCH under first OR
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.condition.1.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.0.condition.1.values.resource.Tag(environment)", "staging"),
					// Level 2: Second condition (MATCH for team)
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.1.operator", "MATCH"),
					resource.TestCheckResourceAttr("ovh_iam_policy.policy1", "conditions.0.condition.1.values.resource.Tag(team)", "platform"),
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
	except 	    = ["%s"]
}
`

const testAccIamPolicyDenyConfig = `
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
	deny 	    = ["%s"]
}
`

const testAccIamPolicyConditionsConfig = `
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
	allow       = ["vps:apiovh:*"]
	
	conditions {
		operator = "OR"
		
		condition {
			operator = "MATCH"
			values = {
				"resource.Tag(environment)" = "production"
				"resource.Tag(team)"        = "platform"
			}
		}
		
		condition {
			operator = "MATCH"
			values = {
				"date(Europe/Paris).WeekDay" = "monday"
			}
		}
	}
}
`

const testAccIamPolicyOrConditionsConfig = `
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
	allow       = ["vps:apiovh:*"]
	
	conditions {
		operator = "OR"

		condition {
			operator = "MATCH"
			values = {
				"resource.Tag(environment)" = "production"
			}
		}

		condition {
			operator = "MATCH"
			values = {
				"resource.Tag(team)" = "platform"
			}
		}
	}
}
`

const testAccIamPolicyWithExpirationConfig = `
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
	allow       = ["vps:apiovh:*"]
	expired_at  = "%s"
}
`

const testAccIamPolicyWithExpirationAndConditionsConfig = `
resource "ovh_me_identity_user" "test_user" {
	login = "%s"
	email = "%s@terraform.test"
	password = "qwe123!@#"
}

resource "ovh_iam_policy" "policy1" {
  name        = "%s"
  description = "%s"

  identities = [ovh_me_identity_user.test_user.urn]
  resources  = ["%s"]

  allow = ["dnsZone:apiovh:get"]

  expired_at = "%s"

  conditions {
    operator = "MATCH"
    values = {
      "resource.Tag(Environment)" = "development"
    }
  }
}
`

const testAccIamPolicyMaxDepthConditionsConfig = `
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
	allow       = ["vps:apiovh:*"]
	
	conditions {
		operator = "AND"
		
		condition {
			operator = "OR"
			
			condition {
				operator = "MATCH"
				values = {
					"resource.Tag(environment)" = "production"
				}
			}
			
			condition {
				operator = "MATCH"
				values = {
					"resource.Tag(environment)" = "staging"
				}
			}
		}
		
		condition {
			operator = "MATCH"
			values = {
				"resource.Tag(team)" = "platform"
			}
		}
	}
}
`
