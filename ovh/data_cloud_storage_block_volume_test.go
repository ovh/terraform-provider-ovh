package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudStorageBlockVolumeConfig = `
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}

data "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  id           = ovh_cloud_storage_block_volume.volume.id
}
`

func TestAccDataSourceCloudStorageBlockVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccDataSourceCloudStorageBlockVolumeConfig,
		serviceName, volumeName, region,
		serviceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_storage_block_volume.volume", "id",
						"ovh_cloud_storage_block_volume.volume", "id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume.volume", "name", volumeName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume.volume", "size", "10"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume.volume", "volume_type", "CLASSIC"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volume.volume", "location.region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volume.volume", "resource_status"),
				),
			},
		},
	})
}
