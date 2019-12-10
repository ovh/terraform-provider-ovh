package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPublicCloudRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPublicCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudRegionsDatasourceConfig,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

func TestAccPublicCloudRegionsDataSource_withNetworkUp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPublicCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudRegionsDatasourceConfig_withNetwork,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

func TestAccPublicCloudRegionsDataSource_withProjectIdEnvVar(t *testing.T) {
	os.Setenv("OVH_PROJECT_ID", os.Getenv("OVH_PUBLIC_CLOUD"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPublicCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicCloudRegionsDatasourceConfig_withProjectIdEnvVar,
				Check:  resource.TestCheckResourceAttrSet("data.ovh_cloud_regions.regions", "names.#"),
			},
		},
	})
}

var testAccPublicCloudRegionsDatasourceConfig = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id = "%s"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccPublicCloudRegionsDatasourceConfig_withProjectIdEnvVar = `
data "ovh_cloud_regions" "regions" {}
`

var testAccPublicCloudRegionsDatasourceConfig_withNetwork = fmt.Sprintf(`
data "ovh_cloud_regions" "regions" {
  project_id      = "%s"
  has_services_up = ["network"]
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))
