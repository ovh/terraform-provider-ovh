package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionsDatasourceConfig,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

func TestAccCloudRegionsDeprecatedDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionsDatasourceDeprecatedConfig,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

func TestAccCloudRegionsDataSource_withNetworkUp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionsDatasourceConfig_withNetwork,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

var testAccCloudRegionsDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  service_name = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudRegionsDatasourceDeprecatedConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id = "%s"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudRegionsDatasourceConfig_withNetwork = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  service_name    = "%s"
  has_services_up = ["network"]
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))
