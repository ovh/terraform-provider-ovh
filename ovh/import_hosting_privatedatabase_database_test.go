package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccHostingPrivateDatabaseDatabaseImportBasic = `
resource "ovh_hosting_privatedatabase_database" "database" {
    service_name  = "%s"
    database_name = "%s"
}
`

func TestAccHostingPrivateDatabaseDatabase_importBasic(t *testing.T) {
	resourceName := "ovh_hosting_privatedatabase_database.database"
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseDatabaseImportBasic, serviceName, databaseName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostingPrivateDatabaseDatabaseImportId(resourceName),
			},
		},
	})
}

func testAccHostingPrivateDatabaseDatabaseImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ds, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_hosting_privatedatabase_database not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			ds.Primary.Attributes["service_name"],
			ds.Primary.Attributes["database_name"],
		), nil
	}
}
