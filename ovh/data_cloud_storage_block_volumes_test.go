package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudStorageBlockVolumesConfig = `
resource "ovh_cloud_storage_block_volume" "volume" {
  service_name = "%s"
  name         = "%s"
  size         = 10
  region       = "%s"
  volume_type  = "CLASSIC"
}

data "ovh_cloud_storage_block_volumes" "volumes" {
  service_name = "%s"
  region       = "%s"

  depends_on = [ovh_cloud_storage_block_volume.volume]
}
`

func TestAccDataSourceCloudStorageBlockVolumes_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccDataSourceCloudStorageBlockVolumesConfig,
		serviceName, volumeName, region,
		serviceName, region,
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
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volumes.volumes", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volumes.volumes", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volumes.volumes", "volumes.0.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volumes.volumes", "volumes.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volumes.volumes", "volumes.0.size"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_block_volumes.volumes", "volumes.0.volume_type"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_block_volumes.volumes", "volumes.0.location.region", region),
				),
			},
		},
	})
}
