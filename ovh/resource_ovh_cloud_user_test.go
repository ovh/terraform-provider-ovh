package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccCloudUserConfig = fmt.Sprintf(`
resource "ovh_cloud_user" "user" {
	project_id  = "%s"
  description = "my user for acceptance tests"
}
`, os.Getenv("OVH_PUBLIC_CLOUD"))

func TestAccCloudUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCloud(t); testAccCheckCloudExists(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudUserExists("ovh_cloud_user.user", t),
					testAccCheckCloudUserOpenRC("ovh_cloud_user.user", t),
				),
			},
		},
	})
}

func testAccCheckCloudUserExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("No Project ID is set")
		}

		return cloudUserExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
	}
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

func testAccCheckCloudUserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_user" {
			continue
		}

		err := cloudUserExists(rs.Primary.Attributes["project_id"], rs.Primary.ID, config.OVHClient)
		if err == nil {
			return fmt.Errorf("VRack > Public Cloud User still exists")
		}

	}
	return nil
}
