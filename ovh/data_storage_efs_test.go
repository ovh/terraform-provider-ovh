package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStorageEfs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckStorageEfs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "ovh_storage_efs" "efs" {
  service_name = "%s"
}`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "created_at"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "region"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "performance_level"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "product"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "quota"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs.efs", "iam.urn"),
					resource.TestCheckResourceAttr("data.ovh_storage_efs.efs", "status", "running"),
				),
			},
		},
	})
}
