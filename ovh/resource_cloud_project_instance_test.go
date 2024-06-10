package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectInstance_basic(t *testing.T) {
	var testCreateLoadBalancerLogSubscription = fmt.Sprintf(`
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
			name = "haproxy"
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
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateLoadBalancerLogSubscription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
				),
			},
		},
	})
}
