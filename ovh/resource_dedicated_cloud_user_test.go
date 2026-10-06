package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCheckOvhDedicatedCloudUserConfig_basic = `
resource "ovh_dedicated_cloud_user" "testuser" {
	service_name = "%s"
	name         = "%s"
	email        = "john.doe@example.com"
	first_name   = "John"
	last_name    = "Doe"
}
`

const testAccCheckOvhDedicatedCloudUserConfig_update = `
resource "ovh_dedicated_cloud_user" "testuser" {
	service_name = "%s"
	name         = "%s"
	email        = "jane.doe@example.com"
	first_name   = "Jane"
	last_name    = "Doe"
}
`

func TestAccDedicatedCloudUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	// The dedicatedCloud user "name" is validated server-side as a shortname:
	// hyphens (and thus acctest.RandomWithPrefix) are rejected with a 400.
	userName := "tfacctest" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCloudUser(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDedicatedCloudUserConfig_basic, serviceName, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user.testuser", "name", userName),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user.testuser", "email", "john.doe@example.com"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user.testuser", "first_name", "John"),
					resource.TestCheckResourceAttrSet(
						"ovh_dedicated_cloud_user.testuser", "user_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_dedicated_cloud_user.testuser", "login"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckOvhDedicatedCloudUserConfig_update, serviceName, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user.testuser", "email", "jane.doe@example.com"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user.testuser", "first_name", "Jane"),
				),
			},
		},
	})
}

func TestAccDedicatedCloudUser_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	// The dedicatedCloud user "name" is validated server-side as a shortname:
	// hyphens (and thus acctest.RandomWithPrefix) are rejected with a 400.
	userName := "tfacctest" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCloudUser(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDedicatedCloudUserConfig_basic, serviceName, userName),
			},
			{
				ResourceName:                         "ovh_dedicated_cloud_user.testuser",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateVerifyIgnore:              []string{"password", "right", "network_role", "vm_network_role", "can_add_ressource", "expiration_date"},
				ImportStateIdFunc:                    testAccDedicatedCloudUser_import("ovh_dedicated_cloud_user.testuser"),
			},
		},
	})
}

func testAccDedicatedCloudUser_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_dedicated_cloud_user not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["user_id"],
		), nil
	}
}
