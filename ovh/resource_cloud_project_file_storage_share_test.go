package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectFileStorageShare_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_REGION_TEST")
	networkId := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_SUBNET_ID_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_REGION_TEST")
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_NETWORK_ID_TEST")
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_SUBNET_ID_TEST")
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "ovh_cloud_project_file_storage_share" "share" {
  service_name = "%s"
  region_name  = "%s"
  name         = "test_share"
  description  = "Test file storage share"
  size         = 500
  type         = "standard-1az"
  network_id   = "%s"
  subnet_id    = "%s"
}`, serviceName, regionName, networkId, subnetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "region_name", regionName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share.share", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "name", "test_share"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "description", "Test file storage share"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "size", "500"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "type", "standard-1az"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "network_id", networkId),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "subnet_id", subnetId),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "status", "available"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "ovh_cloud_project_file_storage_share" "share" {
  service_name = "%s"
  region_name  = "%s"
  name         = "test_share_updated"
  description  = "Test file storage share updated"
  size         = 600
  type         = "standard-1az"
  network_id   = "%s"
  subnet_id    = "%s"
}`, serviceName, regionName, networkId, subnetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "region_name", regionName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share.share", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share.share", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "name", "test_share_updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "description", "Test file storage share updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "size", "600"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "type", "standard-1az"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "network_id", networkId),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "subnet_id", subnetId),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "status", "available"),
				),
			},
		},
	})
}

func TestAccCloudProjectFileStorageShare_withShareNetworkId(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_REGION_TEST")
	networkId := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_NETWORK_ID_TEST")
	subnetId := os.Getenv("OVH_CLOUD_PROJECT_FILE_STORAGE_SUBNET_ID_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_REGION_TEST")
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_NETWORK_ID_TEST")
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FILE_STORAGE_SUBNET_ID_TEST")
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "ovh_cloud_project_file_storage_share_network" "sn" {
  service_name = "%s"
  region_name  = "%s"
  name         = "test_share_network_for_share"
  network_id   = "%s"
  subnet_id    = "%s"
}

resource "ovh_cloud_project_file_storage_share" "share" {
  service_name     = "%s"
  region_name      = "%s"
  name             = "test_share_with_sn"
  description      = "Test with explicit share network"
  size             = 500
  type             = "standard-1az"
  share_network_id = ovh_cloud_project_file_storage_share_network.sn.id
}`, serviceName, regionName, networkId, subnetId, serviceName, regionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share.share", "id"),
					resource.TestCheckResourceAttrPair(
						"ovh_cloud_project_file_storage_share.share", "share_network_id",
						"ovh_cloud_project_file_storage_share_network.sn", "id",
					),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share.share", "status", "available"),
				),
			},
		},
	})
}
