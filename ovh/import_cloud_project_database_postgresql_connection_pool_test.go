package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectDatabasePostgresqlConnectionPool_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_POSTGRESQL_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "johndoe"
	replication := "replication"
	connectionPoolName := "test_connection_pool"
	mode := "session"
	size := 13

	config := fmt.Sprintf(
		testAccCloudProjectDatabasePostgresqlConnectionPoolConfig,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		replication,
		connectionPoolName,
		mode,
		size,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            "ovh_cloud_project_database_postgresql_connection_pool.connection_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudProjectDatabasePostgresqlConnectionPoolImportId("ovh_cloud_project_database_postgresql_connection_pool.connection_pool"),
				ImportStateVerifyIgnore: []string{"uri"},
			},
		},
	})
}

func testAccCloudProjectDatabasePostgresqlConnectionPoolImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testPostgresqlConnectionPool, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_database_postgresql_connection_pool not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			testPostgresqlConnectionPool.Primary.Attributes["service_name"],
			testPostgresqlConnectionPool.Primary.Attributes["cluster_id"],
			testPostgresqlConnectionPool.Primary.Attributes["id"],
		), nil
	}
}
