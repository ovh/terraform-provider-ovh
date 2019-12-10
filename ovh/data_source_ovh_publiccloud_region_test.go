package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPublicCloudRegionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPublicCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudRegionDatasourceConfig,
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

var testAccPublicCloudRegionDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id = "%s"
}

data "ovh_cloud_region" "region_attr" {
  count = 3
  project_id = data.ovh_cloud_regions.regions.project_id
  name = element(sort(data.ovh_cloud_regions.regions.names), count.index)
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))
