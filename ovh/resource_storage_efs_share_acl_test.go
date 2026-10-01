package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccStorageEfsShareAclConfig_basic = `
resource "ovh_storage_efs_share" "share" {
  service_name = "%s"
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}

resource "ovh_storage_efs_share_acl" "acl" {
  service_name = "%s"
  share_id = ovh_storage_efs_share.share.id
  access_level = "rw"
  access_to = "10.0.0.1/32"
}
`

func TestAccStorageEfsShareAcl_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckStorageEfs(t); testAccCheckStorageEfsExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStorageEfsShareAclConfig_basic,
					serviceName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "created_at"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_level", "rw"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_to", "10.0.0.1/32"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_type", "ip"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "status", "active"),
				),
			},
		},
	})

}

func TestAccStorageEfsShareAcl_import(t *testing.T) {
	serviceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")
	config := fmt.Sprintf(testAccStorageEfsShareAclConfig_basic, serviceName, serviceName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckStorageEfs(t); testAccCheckStorageEfsExists(t) },
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "share_id"),
					resource.TestCheckResourceAttrSet("ovh_storage_efs_share_acl.acl", "created_at"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_level", "rw"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_to", "10.0.0.1/32"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "access_type", "ip"),
					resource.TestCheckResourceAttr("ovh_storage_efs_share_acl.acl", "status", "active"),
				),
			},
			{
				ResourceName:            "ovh_storage_efs_share_acl.acl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ImportStateIdFunc:       testAccStorageEfsShareAclImportId("ovh_storage_efs_share_acl.acl"),
			},
		},
	})
}

func testAccStorageEfsShareAclImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["share_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
