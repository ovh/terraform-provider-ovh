package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectRegionsDatasourceConfig,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_project_regions.regions", "names.#"),
			},
		},
	})
}

func TestAccCloudProjectRegionsDataSource_withNetworkUp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectRegionsDatasourceConfig_withNetwork,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_project_regions.regions", "names.#"),
			},
		},
	})
}

var testAccCloudProjectRegionsDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_project_regions" "regions" {
  service_name = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudProjectRegionsDatasourceConfig_withNetwork = fmt.Sprintf(`
data "ovh_cloud_project_regions" "regions" {
  service_name    = "%s"
  has_services_up = ["network"]
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))
