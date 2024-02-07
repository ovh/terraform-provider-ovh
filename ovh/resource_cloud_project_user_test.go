package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectUserConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudProjectUserWithRoleConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
 role_name    = "administrator"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudProjectUserWithRolesConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
 role_names   = ["administrator", "compute_operator"]
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccCloudProjectUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "description", "my user for acceptance tests"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user", t),
				),
			},
		},
	})
}

func TestAccCloudProjectUser_withRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectUserWithRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "description", "my user for acceptance tests"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "roles.0.name", "administrator"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user", t),
				),
			},
		},
	})
}

var updatedConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 role_names   = ["compute_operator"]
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccCloudProjectUser_withRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectUserWithRolesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "description", "my user for acceptance tests"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "roles.#", "2"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user", t),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user", "roles.#", "1"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user", t),
				),
			},
		},
	})
}

func testAccCheckCloudProjectUserOpenRC(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["openstack_rc.OS_AUTH_URL"] == "" {
			return fmt.Errorf("No openstack_rc.OS_AUTH_URL is set")
		}

		if rs.Primary.Attributes["openstack_rc.OS_TENANT_ID"] == "" {
			return fmt.Errorf("No openstack_rc.OS_TENANT_ID is set")
		}

		if rs.Primary.Attributes["openstack_rc.OS_TENANT_NAME"] == "" {
			return fmt.Errorf("No openstack_rc.OS_TENANT_NAME is set")
		}

		if rs.Primary.Attributes["openstack_rc.OS_USERNAME"] == "" {
			return fmt.Errorf("No openstack_rc.OS_USERNAME is set")
		}

		return nil
	}
}
