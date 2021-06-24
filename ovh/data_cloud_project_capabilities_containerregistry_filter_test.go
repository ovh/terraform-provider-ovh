package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectCapabilitiesContainerRegistryFilterDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	planName := "SMALL"
	region := "GRA"

	config := fmt.Sprintf(
		testAccCloudProjectCapabilitiesContainerRegistryFilterDatasourceConfig_Basic,
		serviceName,
		region,
		planName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_capabilities_containerregistry_filter.cap",
						"region",
						"GRA",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_capabilities_containerregistry_filter.cap",
						"name",
						planName,
					),
				),
			},
		},
	})
}

const testAccCloudProjectCapabilitiesContainerRegistryFilterDatasourceConfig_Basic = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "cap" {
  service_name = "%s"
  region       = "%s"
  plan_name    = "%s"
}
`
