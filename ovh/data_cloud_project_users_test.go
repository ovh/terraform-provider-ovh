package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataCloudProjectUsersConfig_basic = `
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

data "ovh_cloud_project_users" "project_users" {
 service_name = ovh_cloud_project_user.user.service_name
}

output "returns_users" {
  value = length(data.ovh_cloud_project_users.project_users.users) >= 1
}

output "user_found" {
  value = contains([for user in data.ovh_cloud_project_users.project_users.users: user.user_id], ovh_cloud_project_user.user.id)
}
`

func TestAccDataCloudProjectUsers_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccDataCloudProjectUsersConfig_basic, serviceName)
	resourceName := "data.ovh_cloud_project_users.project_users"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
					resource.TestCheckOutput("returns_users", "true"),
					resource.TestCheckOutput("user_found", "true"),
				),
			},
		},
	})
}
