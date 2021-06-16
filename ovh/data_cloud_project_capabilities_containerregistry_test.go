package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectCapabilitiesContainerRegistryDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectCapabilitiesContainerRegistryDatasourceConfig_Basic,
		serviceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_capabilities_containerregistry.cap",
						"result.#",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_capabilities_containerregistry.cap",
						"result.0.region_name", "GRA"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_capabilities_containerregistry.cap",
						"result.0.plans.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_capabilities_containerregistry.cap",
						"result.0.plans.0.code",
					),
				),
			},
		},
	})
}

const testAccCloudProjectCapabilitiesContainerRegistryDatasourceConfig_Basic = `
data "ovh_cloud_project_capabilities_containerregistry" "cap" {
  service_name = "%s"
}
`
