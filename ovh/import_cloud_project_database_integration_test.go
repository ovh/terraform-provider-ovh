package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudProjectDatabaseIntegration_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	versionPsql := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_POSTGRESQL_VERSION_TEST")
	versionOs := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_OPENSEARCH_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	descriptionPsql := acctest.RandomWithPrefix(test_prefix)
	descriptionOs := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseIntegrationConfig,
		serviceName,
		descriptionPsql,
		versionPsql,
		region,
		flavor,
		serviceName,
		descriptionOs,
		versionOs,
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
				ResourceName:      "ovh_cloud_project_database_integration.integration",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseIntegrationImportId("ovh_cloud_project_database_integration.integration"),
			},
		},
	})
}

func testAccCloudProjectDatabaseIntegrationImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testIntegration, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_integration not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			testIntegration.Primary.Attributes["service_name"],
			testIntegration.Primary.Attributes["engine"],
			testIntegration.Primary.Attributes["cluster_id"],
			testIntegration.Primary.Attributes["id"],
		), nil
	}
}
