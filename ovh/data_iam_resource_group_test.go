package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIamResourceGroupDataSource_basic(t *testing.T) {
	resourceGroupName1 := acctest.RandomWithPrefix(test_prefix)
	resourceGroupName2 := acctest.RandomWithPrefix(test_prefix)

	resource1 := "urn:v1:eu:resource:vrack:" + os.Getenv("OVH_VRACK_SERVICE_TEST")
	resource2 := "urn:v1:eu:resource:vps:" + os.Getenv("OVH_VPS")

	preSetup := fmt.Sprintf(
		testAccIamResourceGroupDatasourceConfigInit,
		resourceGroupName1,
		resource1,
		resourceGroupName2,
		resource1,
		resource2,
	)

	dataConfig := fmt.Sprintf(
		testAccIamResourceGroupDatasourceConfigData,
		resourceGroupName1,
		resource1,
		resourceGroupName2,
		resource1,
		resource2,
	)

	config := fmt.Sprintf(
		testAccIamResourceGroupDatasourceConfigList,
		resourceGroupName1,
		resource1,
		resourceGroupName2,
		resource1,
		resource2,
	)

	checks := append(
		checkIamResourceGroupResourceAttr("ovh_iam_resource_group.resource_group_1", resourceGroupName1, resource1),
		checkIamResourceGroupResourceAttr("ovh_iam_resource_group.resource_group_2", resourceGroupName2, resource1, resource2)...,
	)

	checksData := append(
		checkIamResourceGroupResourceAttr("data.ovh_iam_resource_group.resource_group_1", resourceGroupName1, resource1),
		checkIamResourceGroupResourceAttr("data.ovh_iam_resource_group.resource_group_2", resourceGroupName2, resource1, resource2)...,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIamResourceGroup(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: preSetup,
				Check:  resource.ComposeTestCheckFunc(checks...),
			}, {
				Config: dataConfig,
				Check:  resource.ComposeTestCheckFunc(checksData...),
			}, {
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("keys_present", "true"),
				),
			},
		},
	})
}

func checkIamResourceGroupResourceAttr(name, grpName string, resourceURNs ...string) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "name", grpName),
	}
	for _, urn := range resourceURNs {
		checks = append(checks,
			resource.TestCheckTypeSetElemAttr(name, "resources.*", urn),
		)
	}
	return checks
}

const testAccIamResourceGroupDatasourceConfigInit = `
resource "ovh_iam_resource_group" "resource_group_1" {
	name        = "%s"
	resources   = ["%s"]
}

resource "ovh_iam_resource_group" "resource_group_2" {
	name        = "%s"
	resources   = ["%s", "%s"]
}
`

const testAccIamResourceGroupDatasourceConfigData = `
resource "ovh_iam_resource_group" "resource_group_1" {
	name        = "%s"
	resources   = ["%s"]
}

resource "ovh_iam_resource_group" "resource_group_2" {
	name        = "%s"
	resources   = ["%s", "%s"]
}

data "ovh_iam_resource_group" "resource_group_1" {
	id = ovh_iam_resource_group.resource_group_1.id
}

data "ovh_iam_resource_group" "resource_group_2" {
	id = ovh_iam_resource_group.resource_group_2.id
}
`

const testAccIamResourceGroupDatasourceConfigList = `
resource "ovh_iam_resource_group" "resource_group_1" {
	name        = "%s"
	resources   = ["%s"]
}

resource "ovh_iam_resource_group" "resource_group_2" {
	name        = "%s"
	resources   = ["%s", "%s"]
}

data "ovh_iam_resource_group" "data_resource_group_1" {
	id = ovh_iam_resource_group.resource_group_1.id
}


data "ovh_iam_resource_groups" "resource_groups" {}

output "keys_present" {
	value = tostring(
		contains(data.ovh_iam_resource_groups.resource_groups.resource_groups, ovh_iam_resource_group.resource_group_1.id) &&
		contains(data.ovh_iam_resource_groups.resource_groups.resource_groups, ovh_iam_resource_group.resource_group_2.id)
	)
}
`
