package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudRegionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.0", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.0", "services.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.1", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.1", "services.#"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.2", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_region.region_attr.2", "services.#"),
				),
			},
		},
	})
}

var testAccCloudRegionDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id = "%s"
}

data "ovh_cloud_region" "region_attr" {
  count = 3
  project_id = data.ovh_cloud_regions.regions.project_id
  name = element(sort(data.ovh_cloud_regions.regions.names), count.index)
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))
