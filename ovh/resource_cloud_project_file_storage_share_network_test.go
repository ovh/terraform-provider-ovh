package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectFileStorageShareNetwork_basic(t *testing.T) {
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
  name         = "test_share_network"
  description  = "Test share network"
  network_id   = "%s"
  subnet_id    = "%s"
}`, serviceName, regionName, networkId, subnetId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "region_name", regionName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share_network.sn", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_file_storage_share_network.sn", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "name", "test_share_network"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "description", "Test share network"),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "network_id", networkId),
					resource.TestCheckResourceAttr("ovh_cloud_project_file_storage_share_network.sn", "subnet_id", subnetId),
				),
			},
		},
	})
}
