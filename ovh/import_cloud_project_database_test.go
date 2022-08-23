package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudProjectDatabase_importBasic(t *testing.T) {
	description := acctest.RandomWithPrefix(test_prefix)
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		description,
		engine,
		version,
		region,
		flavor,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_project_database.db",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseImportId("ovh_cloud_project_database.db"),
			},
		},
	})
}

func testAccCloudProjectDatabaseImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testDatabase, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testDatabase.Primary.Attributes["service_name"],
			testDatabase.Primary.Attributes["engine"],
			testDatabase.Primary.Attributes["id"],
		), nil
	}
}
