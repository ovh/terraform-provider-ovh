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

func TestAccCloudRegionsDataSource_withProjectIdEnvVar(t *testing.T) {
	os.Setenv("OVH_PROJECT_ID", os.Getenv("OVH_PUBLIC_CLOUD"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionsDatasourceConfig_withProjectIdEnvVar,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

var testAccCloudRegionsDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id = "%s"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccCloudRegionsDatasourceConfig_withProjectIdEnvVar = `
data "ovh_cloud_regions" "regions" {}
`

var testAccCloudRegionsDatasourceConfig_withNetwork = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id      = "%s"
  has_services_up = ["network"]
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))
