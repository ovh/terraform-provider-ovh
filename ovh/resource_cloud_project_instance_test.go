package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectInstance_basic(t *testing.T) {
	var testCreateInstance = fmt.Sprintf(`
			resource "ovh_cloud_project_instance" "instance" {
				service_name = "%s"
				region = "%s"
				billing_period = "hourly"
				boot_from {
					image_id = "46d79419-c5bb-4b78-90a5-75ea3d6f5767"
				}
				flavor {
					flavor_id = "21774833-28aa-442f-86b2-156343e3217f"
				}
				name = "TestInstance"
				ssh_key {
					name = "%s"
				}
				network {
					public = true
				}
			}
		`,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_name", "b3-8"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", "21774833-28aa-442f-86b2-156343e3217f"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", "46d79419-c5bb-4b78-90a5-75ea3d6f5767"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", "TestInstance"),
				),
			},
		},
	})
}
