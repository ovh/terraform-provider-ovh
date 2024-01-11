package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMeIdentityGroupDataSource_basic(t *testing.T) {
	groupName := acctest.RandomWithPrefix(test_prefix)
	desc := "Identity group created by Terraform Acc."
	role := "NONE"

	config := fmt.Sprintf(testAccMeIdentityGroupDatasourceConfig, desc, groupName, role)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIdentityGroupResourceAttr("data.ovh_me_identity_group.group_1", groupName, desc, role)...,
				),
			},
		},
	})
}

const testAccMeIdentityGroupDatasourceConfig = `
resource "ovh_me_identity_group" "group_1" {
	description = "%s"
	name        = "%s"
  	role        = "%s"
}

data "ovh_me_identity_group" "group_1" {
  name       = ovh_me_identity_group.group_1.name
  depends_on = [ovh_me_identity_group.group_1]
}
`
