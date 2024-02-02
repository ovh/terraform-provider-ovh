package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccDataCloudProjectUserConfig_basic = `
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

data "ovh_cloud_project_user" "user" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
 depends_on   = [ovh_cloud_project_user.user]
}
`

func TestAccDataCloudProjectUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccDataCloudProjectUserConfig_basic, serviceName)
	resourceName := "data.ovh_cloud_project_user.user"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "my user for acceptance tests"),
				),
			},
		},
	})
}
