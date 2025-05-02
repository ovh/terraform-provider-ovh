package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectRegionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_cloud_project_region" "region" {
						service_name = "%s"
						name         = "EU-WEST-PAR"
					}
				`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region", "type"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region", "services.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region", "availability_zones.#"),
				),
			},
		},
	})
}
