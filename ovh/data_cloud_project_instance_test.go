package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjecInstance_basic(t *testing.T) {

	config := fmt.Sprintf(
		testAccDataSourceCloudProjectInstance,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_INSTANCE_ID_TEST"),
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
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "flavor_name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "flavor_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "image_id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "ssh_key"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_instance.test", "region"),
				),
			},
		},
	})
}

var testAccDataSourceCloudProjectInstance = `
data "ovh_cloud_project_instance" "test" {
	service_name = "%s"
	region = "%s"
	instance_id = "%s"
}
`
