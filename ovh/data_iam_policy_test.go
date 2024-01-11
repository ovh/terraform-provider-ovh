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
