package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMeIdentityGroupsDataSource_basic(t *testing.T) {
	desc := "Identity group created by Terraform Acc."
	role1 := "NONE"
	role2 := "ADMIN"
	groupName1 := acctest.RandomWithPrefix(test_prefix)
	groupName2 := acctest.RandomWithPrefix(test_prefix)

	preSetup := fmt.Sprintf(
		testAccMeIdentityGroupsDatasourceConfig_preSetup,
		desc,
		groupName1,
		role1,
		desc,
		groupName2,
		role2,
	)
	config := fmt.Sprintf(
		testAccMeIdentityGroupsDatasourceConfig_keys,
		desc,
		groupName1,
		role1,
		desc,
		groupName2,
		role2,
	)

	checks := checkIdentityGroupResourceAttr("ovh_me_identity_group.group_1", groupName1, desc, role1)
	checks = append(checks, checkIdentityGroupResourceAttr("ovh_me_identity_group.group_2", groupName2, desc, role2)...)

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

func checkIdentityGroupResourceAttr(name, g_name, desc, role string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "name", g_name),
		resource.TestCheckResourceAttr(name, "description", desc),
		resource.TestCheckResourceAttr(name, "role", role),
	}
}

const testAccMeIdentityGroupsDatasourceConfig_preSetup = `
resource "ovh_me_identity_group" "group_1" {
	description = "%s"
	name        = "%s"
	role        = "%s"
}

resource "ovh_me_identity_group" "group_2" {
	description = "%s"
	name        = "%s"
	role        = "%s"
}
`

const testAccMeIdentityGroupsDatasourceConfig_keys = `
resource "ovh_me_identity_group" "group_1" {
	description = "%s"
	name        = "%s"
	role        = "%s"
}

resource "ovh_me_identity_group" "group_2" {
	description = "%s"
	name        = "%s"
	role        = "%s"
}

data "ovh_me_identity_groups" "groups" {}

output "keys_present" {
	value = tostring(contains(data.ovh_me_identity_groups.groups.groups, ovh_me_identity_group.group_1.name) && contains(data.ovh_me_identity_groups.groups.groups, ovh_me_identity_group.group_2.name))
}
`
