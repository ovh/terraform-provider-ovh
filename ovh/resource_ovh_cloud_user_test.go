package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccCloudUserConfig = fmt.Sprintf(`
resource "ovh_cloud_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccCloudUserWithRoleConfig = fmt.Sprintf(`
resource "ovh_cloud_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
 role_name    = "administrator"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccCloudUserWithRolesConfig = fmt.Sprintf(`
resource "ovh_cloud_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
 role_names   = ["administrator", "compute_operator"]
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

var testAccCloudUserDeprecatedConfig = fmt.Sprintf(`
resource "ovh_cloud_user" "user" {
  project_id  = "%s"
  description = "my user for acceptance tests"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccCloudUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "description", "my user for acceptance tests"),
					testAccCheckCloudUserOpenRC("ovh_cloud_user.user", t),
				),
			},
		},
	})
}

func TestAccCloudUserDeprecated_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserDeprecatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "description", "my user for acceptance tests"),
					testAccCheckCloudUserOpenRC("ovh_cloud_user.user", t),
				),
			},
		},
	})
}

func TestAccCloudUser_withRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserWithRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "description", "my user for acceptance tests"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "roles.0.name", "administrator"),
					testAccCheckCloudUserOpenRC("ovh_cloud_user.user", t),
				),
			},
		},
	})
}

func TestAccCloudUser_withRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserWithRolesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "description", "my user for acceptance tests"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_user.user", "roles.#", "2"),
					testAccCheckCloudUserOpenRC("ovh_cloud_user.user", t),
				),
			},
		},
	})
}

func testAccCheckCloudUserOpenRC(n string, t *testing.T) resource.TestCheckFunc {
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
