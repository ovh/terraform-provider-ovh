package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectContainerRegistriesDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistriesDatasourceConfig_Basic,
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
						"data.ovh_cloud_project_containerregistries.regs",
						"result.#",
					),
				),
			},
		},
	})
}

const testAccCloudProjectContainerRegistriesDatasourceConfig_Basic = `
data "ovh_cloud_project_containerregistries" "regs" {
  service_name = "%s"
}
`
