package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectVolume_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "volume"  {
						region_name = "%s"
						service_name = "%s"
						description = "test"
						name = "test"
						size = 15
						type = "classic"
					}
				`,
					regionName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume.volume", "volume_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "type", "classic"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "description", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "name", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "size", "15"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_volume" "volume"  {
						region_name = "%s"
						service_name = "%s"
						description = "test_updated"
						name = "test_updated"
						size = 20
						type = "classic"
					}
				`,
					regionName,
					serviceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "region_name", regionName),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_volume.volume", "volume_id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "type", "classic"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "description", "test_updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "name", "test_updated"),
					resource.TestCheckResourceAttr("ovh_cloud_project_volume.volume", "size", "20"),
				),
			},
		},
	})
}
