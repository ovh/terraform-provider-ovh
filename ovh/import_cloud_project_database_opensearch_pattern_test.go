package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseOpensearchPattern_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_OPENSEARCH_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	maxIndexCount := 2
	pattern := "logs_*"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseOpensearchPatternConfig,
		serviceName,
		description,
		version,
		region,
		flavor,
		maxIndexCount,
		pattern,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "ovh_cloud_project_database_opensearch_pattern.pattern",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectDatabaseOpensearchPatternImportId("ovh_cloud_project_database_opensearch_pattern.pattern"),
			},
		},
	})
}

func testAccCloudProjectDatabaseOpensearchPatternImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testOpensearchPattern, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_opensearch_pattern not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testOpensearchPattern.Primary.Attributes["service_name"],
			testOpensearchPattern.Primary.Attributes["cluster_id"],
			testOpensearchPattern.Primary.Attributes["id"],
		), nil
	}
}
