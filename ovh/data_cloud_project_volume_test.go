package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	volumeId := os.Getenv("OVH_CLOUD_PROJECT_VOLUME_ID_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_volume" "volume" {
						service_name = "%s"
						region_name  = "%s"
						volume_id           = "%s"
					}
				`,
					serviceName,
					regionName,
					volumeId,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_volume.volume", "region_name", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_volume.volume", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_volume.volume", "volume_id"),
				),
			},
		},
	})
}
