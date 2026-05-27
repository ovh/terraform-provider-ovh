package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccStorageEfsShareConfig_basic = `
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}
`

func TestAccStorageEfsShare_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckStorageEfs(t); testAccCheckStorageEfsExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStorageEfsShareConfig_basic,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "created_at"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "name", "share"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "description", "My share"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "size", "100"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "status", "available"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "snapshot_id", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share_updated"
  description  = "My share updated"
  protocol     = "NFS"
  size         = 120
}
`, serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "created_at"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "name", "share_updated"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "description", "My share updated"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "size", "120"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "status", "available"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "snapshot_id", ""),
				),
			},
		},
	})

}

func TestAccStorageEfsShare_import(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")
	config := fmt.Sprintf(testAccStorageEfsShareConfig_basic, serviceName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckStorageEfs(t); testAccCheckStorageEfsExists(t) },
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share.share", "created_at"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "name", "share"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "description", "My share"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "size", "100"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "status", "available"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share.share", "snapshot_id", ""),
				),
			},
			{
				ResourceName:            "ovh_storage_efs_share.share",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ImportStateIdFunc:       testAccStorageEfsShareImportId("ovh_storage_efs_share.share"),
			},
		},
	})
}

func testAccStorageEfsShareImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
