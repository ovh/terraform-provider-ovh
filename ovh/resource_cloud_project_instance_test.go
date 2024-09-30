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
				image_id = "%s"
			}
			flavor {
				flavor_id = "%s"
			}
			name = "%s"
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
		os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_INSTANCE_NAME_TEST"),
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
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "flavor_name"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST")),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST")),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", os.Getenv("OVH_CLOUD_PROJECT_INSTANCE_NAME_TEST")),
				),
			},
		},
	})
}
