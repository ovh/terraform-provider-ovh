package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccStorageEfsConfig_basic = `
data "ovh_me" "myaccount" {}

resource "ovh_storage_efs" "efs" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
  name           = "MyEFS"

  plan = [
    {
      plan_code = "enterprise-file-storage-premium-1tb"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region"
          value = "%s"
        },
        {
          label = "network"
          value = "vrack"
        }
      ]
    }
  ]
}
`

func TestAccResourceStorageEfs_basic(t *testing.T) {
	region := os.Getenv("OVH_STORAGE_EFS_REGION_TEST")

	config := fmt.Sprintf(
		testAccStorageEfsConfig_basic,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrderStorageEfs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "created_at"),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "iam.urn"),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "id"),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "name"),
					resource.TestCheckResourceAttr(
						"ovh_storage_efs.efs", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "performance_level"),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "product"),
					resource.TestCheckResourceAttrSet(
						"ovh_storage_efs.efs", "quota"),
					resource.TestCheckResourceAttr(
						"ovh_storage_efs.efs", "status", "running"),
				),
			},
		},
	})
}
