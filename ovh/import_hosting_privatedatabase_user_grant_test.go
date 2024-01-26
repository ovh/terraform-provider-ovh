package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccHostingPrivateDatabaseUserGrantImportBasic = `
resource "ovh_hosting_privatedatabase_user_grant" "grant" {
    service_name  = "%s"
    user_name     = "%s"
    database_name = "%s"
    grant         = "%s"
}
`

func TestAccHostingPrivateDatabaseUserGrant_importBasic(t *testing.T) {
	resourceName := "ovh_hosting_privatedatabase_user_grant.grant"
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	grantName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_GRANT_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseUserGrantImportBasic, serviceName, userName, databaseName, grantName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUserGrant(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccHostingPrivateDatabaseUserGrantImportId(resourceName),
			},
		},
	})
}

func testAccHostingPrivateDatabaseUserGrantImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ds, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_hosting_privatedatabase_user_grant not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			ds.Primary.Attributes["service_name"],
			ds.Primary.Attributes["user_name"],
			ds.Primary.Attributes["database_name"],
			ds.Primary.Attributes["grant"],
		), nil
	}
}
