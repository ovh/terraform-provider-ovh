package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamPolicyDataSource_basic(t *testing.T) {
	usrLogin := acctest.RandomWithPrefix(test_prefix)
	grpName := acctest.RandomWithPrefix(test_prefix)

	desc := "Iam policy created by Terraform Acc"
	policyName1 := acctest.RandomWithPrefix(test_prefix)
	policyName2 := acctest.RandomWithPrefix(test_prefix)
	resource1 := "urn:v1:eu:resource:vrack:*"
	resource2 := "urn:v1:eu:resource:vps:*"
	allow1 := "*"
	except1 := "vrack:apiovh:dedicatedServer/detach"
	deny2 := "*"

	preSetup := fmt.Sprintf(
		testAccIamPolicyDatasourceConfig_preSetup,
		usrLogin,
		usrLogin,
		grpName,
		policyName1,
		desc,
		resource1,
		allow1,
		except1,
		policyName2,
		desc,
		resource2,
		deny2,
	)
	config := fmt.Sprintf(
		testAccIamPolicyDatasourceConfig_keys,
		usrLogin,
		usrLogin,
		grpName,
		policyName1,
		desc,
		resource1,
		allow1,
		except1,
		policyName2,
		desc,
		resource2,
		deny2,
	)

	checks := checkIamPolicyResourceAttr("ovh_iam_policy.policy_1", policyName1, desc, resource1, allow1, except1, "")
	checks = append(checks, checkIamPolicyResourceAttr("ovh_iam_policy.policy_2", policyName2, desc, resource2, "", "", deny2)...)

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

func TestAccIamPolicyDataSource_withConditionsAndExpiration(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	desc := "IAM policy with conditions and expiration created by Terraform Acc"
	userName := acctest.RandomWithPrefix(test_prefix)
	res := "urn:v1:eu:resource:vps:*"
	expiration := "2025-12-31T23:59:59Z"
	config := fmt.Sprintf(testAccIamPolicyDataSourceConfig, userName, userName, name, desc, res, expiration)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "name", name),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "description", desc),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "expired_at", expiration),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.#", "1"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.operator", "OR"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.#", "2"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.0.operator", "MATCH"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.0.values.%", "2"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.0.values.resource.Tag(environment)", "production"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.0.values.resource.Tag(team)", "platform"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.1.operator", "MATCH"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.1.values.%", "1"),
					resource.TestCheckResourceAttr("data.ovh_iam_policy.policy", "conditions.0.condition.1.values.date(Europe/Paris).WeekDay", "monday"),
				),
			},
		},
	})
}

func checkIamPolicyResourceAttr(name, polName, desc, resourceURN, allowAction, exceptAction, denyAction string) []resource.TestCheckFunc {
	// we are not checking identity urn because they are dynamic and depend on the test account NIC
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "name", polName),
		resource.TestCheckResourceAttr(name, "description", desc),
		resource.TestCheckTypeSetElemAttr(name, "resources.*", resourceURN),
	}
	if allowAction != "" {
		checks = append(checks, resource.TestCheckTypeSetElemAttr(name, "allow.*", allowAction))
	}
	if exceptAction != "" {
		checks = append(checks, resource.TestCheckTypeSetElemAttr(name, "except.*", exceptAction))
	}
	if denyAction != "" {
		checks = append(checks, resource.TestCheckTypeSetElemAttr(name, "deny.*", denyAction))
	}
	return checks
}

const testAccIamPolicyDatasourceConfig_preSetup = `
resource "ovh_me_identity_user" "user_1" {
	login = "%s"
	email = "%s@terraform.test"
	password = "qwerty123!@#"
}

resource "ovh_me_identity_group" "group_1" {
	name = "%s"
}

resource "ovh_iam_policy" "policy_1" {
	name        = "%s"
	description = "%s"
	identities  = [ovh_me_identity_user.user_1.urn]
	resources   = ["%s"]
	allow       = ["%s"]
	except       = ["%s"]
}

resource "ovh_iam_policy" "policy_2" {
	name        = "%s"
	description = "%s"
	identities  = [ovh_me_identity_group.group_1.urn]
	resources   = ["%s"]
	deny       = ["%s"]
}
`

const testAccIamPolicyDatasourceConfig_keys = `
resource "ovh_me_identity_user" "user_1" {
	login = "%s"
	email = "%s@terraform.test"
	password = "qwerty123!@#"
}

resource "ovh_me_identity_group" "group_1" {
	name = "%s"
}

resource "ovh_iam_policy" "policy_1" {
	name        = "%s"
	description = "%s"
	identities  = [ovh_me_identity_user.user_1.urn]
	resources   = ["%s"]
	allow       = ["%s"]
	except       = ["%s"]
}

resource "ovh_iam_policy" "policy_2" {
	name        = "%s"
	description = "%s"
	identities  = [ovh_me_identity_group.group_1.urn]
	resources   = ["%s"]
	deny       = ["%s"]
}

data "ovh_iam_policies" "policies" {}

output "keys_present" {
	value = tostring(
		contains(data.ovh_iam_policies.policies.policies, ovh_iam_policy.policy_1.id) &&
		contains(data.ovh_iam_policies.policies.policies, ovh_iam_policy.policy_2.id)
	)
}
`

const testAccIamPolicyDataSourceConfig = `
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

data "ovh_iam_policy" "policy" {
	id = ovh_iam_policy.policy1.id
}
`
