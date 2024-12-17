package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabaseMongodbPrometheus_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseMongodbPrometheusConfig,
		serviceName,
		description,
		version,
		region,
		region,
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
				ResourceName:            "ovh_cloud_project_database_mongodb_prometheus.prometheus",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudProjectDatabaseMongodbPrometheusImportId("ovh_cloud_project_database_mongodb_prometheus.prometheus"),
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCloudProjectDatabaseMongodbPrometheusImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testUser, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_mongodb_prometheus not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testUser.Primary.Attributes["service_name"],
			testUser.Primary.Attributes["cluster_id"],
		), nil
	}
}
