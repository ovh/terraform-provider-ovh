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

var testAccCloudProjectUserWithRotateConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user_rotate" {
 service_name = "%s"
 description  = "my user for acceptance tests with rotation"
 password_reset =  "2025-04-29"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

var testAccCloudProjectUserWithRotateUpdatedConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user_rotate" {
 service_name = "%s"
 description  = "my user for acceptance tests with rotation"
 password_reset =  "2025-04-30"
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccCloudProjectUser_withRotate(t *testing.T) {
	var oldPassword, newPassword string
	var userId string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectUserWithRotateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user_rotate", "description", "my user for acceptance tests with rotation"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user_rotate", t),
					testAccCheckCloudProjectUserPassword("ovh_cloud_project_user.user_rotate", &oldPassword),
					testAccCheckCloudProjectUserId("ovh_cloud_project_user.user_rotate", &userId),
				),
			},
			{
				Config: testAccCloudProjectUserWithRotateUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_user.user_rotate", "description", "my user for acceptance tests with rotation"),
					testAccCheckCloudProjectUserOpenRC("ovh_cloud_project_user.user_rotate", t),
					testAccCheckCloudProjectUserPassword("ovh_cloud_project_user.user_rotate", &newPassword), // Store the second password
					testAccCheckCloudProjectUserId("ovh_cloud_project_user.user_rotate", &userId),            // Same user ID (not recreated)
					testAccCheckCloudProjectUserDifferentPasswords(&oldPassword, &newPassword),               // Compare first and second passwords
				),
			},
		},
	})
}

func testAccCheckCloudProjectUserPassword(n string, password *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["password"] == "" {
			return fmt.Errorf("No password is set")
		}

		*password = rs.Primary.Attributes["password"]
		return nil
	}
}

// testAccCheckCloudProjectUserDifferentPasswords creates a test check function that verifies
// if the password has been changed during rotation by comparing old and new password values
func testAccCheckCloudProjectUserDifferentPasswords(oldPassword *string, newPassword *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Verify that we have actually captured the old password
		if *oldPassword == "" {
			return fmt.Errorf("Old password was not captured")
		}

		// Verify that we have actually captured the new password
		if *newPassword == "" {
			return fmt.Errorf("New password was not captured")
		}

		// Check that the password has actually changed
		if *oldPassword == *newPassword {
			return fmt.Errorf("Password did not change after rotation (both are '%s')", *oldPassword)
		}

		return nil
	}
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

func testAccCheckCloudProjectUserId(n string, userId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		*userId = rs.Primary.ID
		return nil
	}
}
