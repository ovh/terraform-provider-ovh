package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStorageEfsShareAccessPaths_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckStorageEfs(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}

data "ovh_storage_efs_share_access_paths" "access_paths" {
  service_name = "%s"
  share_id = ovh_storage_efs_share.share.id
}`,
					serviceName,
					serviceName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_storage_efs_share_access_paths.access_paths", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs_share_access_paths.access_paths", "share_id"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs_share_access_paths.access_paths", "access_paths.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs_share_access_paths.access_paths", "access_paths.0.path"),
					resource.TestCheckResourceAttrSet("data.ovh_storage_efs_share_access_paths.access_paths", "access_paths.0.preferred"),
				),
			},
		},
	})
}
