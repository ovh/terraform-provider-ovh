package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStorageFileShareAcl_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	vrackNetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	vrackSubnetName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareVrackSubnetNamePrefix)
	networkName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNetworkNamePrefix)
	shareName := acctest.RandomWithPrefix(testAccResourceCloudStorageFileShareNamePrefix)

	config := testAccCloudStorageFileShareAclShareConfig(serviceName, region, vrackNetName, vrackSubnetName, networkName, shareName) + fmt.Sprintf(`
resource "ovh_cloud_storage_file_share_acl" "acl" {
  service_name = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_ONLY"
}

data "ovh_cloud_storage_file_share_acl" "acl" {
  service_name = "%s"
  share_id     = ovh_cloud_storage_file_share.share.id
  id           = ovh_cloud_storage_file_share_acl.acl.id

  depends_on = [ovh_cloud_storage_file_share_acl.acl]
}
`, serviceName, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckVRack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_acl.acl", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_file_share_acl.acl", "id",
						"ovh_cloud_storage_file_share_acl.acl", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_acl.acl", "access_to", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_file_share_acl.acl", "access_level", "READ_ONLY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_acl.acl", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_acl.acl", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_acl.acl", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_file_share_acl.acl", "current_state.state"),
				),
			},
		},
	})
}
