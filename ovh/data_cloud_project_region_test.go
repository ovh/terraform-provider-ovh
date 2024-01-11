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
				Config: testAccCloudProjectRegionDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.0", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.0", "services.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.1", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.1", "services.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.2", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_region.region_attr.2", "services.#"),
				),
			},
		},
	})
}

var testAccCloudProjectRegionDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_project_regions" "regions" {
  service_name = "%s"
}

data "ovh_cloud_project_region" "region_attr" {
  count        = 3
  service_name = data.ovh_cloud_project_regions.regions.service_name
  name         = element(sort(data.ovh_cloud_project_regions.regions.names), count.index)
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))
