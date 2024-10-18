package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectVolumes_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_volumes" "volume" {
						service_name = "%s"
						region_name  = "%s"
					}
				`,
					serviceName,
					regionName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_volumes.volume", "region_name", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_volumes.volume", "volumes.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_volumes.volume", "volumes.0.id"),
				),
			},
		},
	})
}
