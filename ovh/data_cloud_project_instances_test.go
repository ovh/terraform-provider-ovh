package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCloudProjecInstancesConfig_basic = `
data "ovh_cloud_project_instances" "instances" {
   service_name = "%s"
   region = "%s"
}
`

func TestAccDataSourceCloudProjecInstances_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	config := fmt.Sprintf(
		testAccDataSourceCloudProjecInstancesConfig_basic,
		serviceName,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_instances.instances",
						"instances.0.flavor_id",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_instances.instances",
						"instances.0.flavor_name",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_instances.instances",
						"instances.0.id",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_instances.instances",
						"instances.0.image_id",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_instances.instances",
						"instances.0.ssh_key",
					),
				),
			},
		},
	})
}
