package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccHostingPrivateDatabaseUserImportBasic = `
resource "ovh_hosting_privatedatabase_user" "user" {
    service_name  = "%s"
    password      = "%s"
    user_name     = "%s"
}
`

func TestAccHostingPrivateDatabaseUser_importBasic(t *testing.T) {
	resourceName := "ovh_hosting_privatedatabase_user.user"
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	password := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseUserImportBasic, serviceName, password, userName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUser(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccHostingPrivateDatabaseUserImportId(resourceName),
			},
		},
	})
}

func testAccHostingPrivateDatabaseUserImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ds, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_hosting_privatedatabase_user not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			ds.Primary.Attributes["service_name"],
			ds.Primary.Attributes["user_name"],
		), nil
	}
}
