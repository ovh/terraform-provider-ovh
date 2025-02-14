package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjectInstances_basic(t *testing.T) {
	config := fmt.Sprintf(`
			data "ovh_cloud_project_floatingips" "instances" {
				service_name = "%s"
				region       = "%s"
			}
		`,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"),
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
