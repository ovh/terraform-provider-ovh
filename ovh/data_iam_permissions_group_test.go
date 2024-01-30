package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamPermissionsGroupDataSource_basic(t *testing.T) {
	grpName := acctest.RandomWithPrefix(test_prefix)

	desc := "Permissions group created by Terraform Acc"
	allow := "account:apiovh:iam/policy/*"
	except := "account:apiovh:iam/policy/delete"
	deny := "account:apiovh:iam/policy/create"

	preSetup := fmt.Sprintf(
		testAccIamPermissionsGroupDatasourceConfig_preSetup,
		grpName,
		desc,
		allow,
		except,
		deny,
	)
	config := fmt.Sprintf(
		testAccIamPermissionsGroupDatasourceConfig_keys,
		grpName,
		desc,
		allow,
		except,
		deny,
	)

	checks := checkIamPermissionsGroupResourceAttr("ovh_iam_permissions_group.permissions", grpName, desc, allow, except, deny)
	dataCheck := checkIamPermissionsGroupResourceAttr("data.ovh_iam_permissions_group.permissions", grpName, desc, allow, except, deny)
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
					append(
						dataCheck,
						resource.TestCheckOutput("keys_present", "true"),
					)...,
				),
			},
		},
	})
}

func checkIamPermissionsGroupResourceAttr(name, permName, desc, allowAction, exceptAction, denyAction string) []resource.TestCheckFunc {
	// we are not checking identity urn because they are dynamic and depend on the test account NIC
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "name", permName),
		resource.TestCheckResourceAttr(name, "description", desc),
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

const testAccIamPermissionsGroupDatasourceConfig_preSetup = `
resource "ovh_iam_permissions_group" "permissions" {
	name        = "%s"
	description = "%s"
	allow       = ["%s"]
	except      = ["%s"]
	deny        = ["%s"]
}
`

const testAccIamPermissionsGroupDatasourceConfig_keys = `
resource "ovh_iam_permissions_group" "permissions" {
	name        = "%s"
	description = "%s"
	allow       = ["%s"]
	except      = ["%s"]
	deny        = ["%s"]
}

data "ovh_iam_permissions_group" "permissions" {
	urn = ovh_iam_permissions_group.permissions.urn
}


data "ovh_iam_permissions_groups" "groups" {}

output "keys_present" {
	value = tostring(
		contains(data.ovh_iam_permissions_groups.groups.urns, ovh_iam_permissions_group.permissions.urn)
	)
}
`
