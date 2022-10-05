package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseIntegrationConfig = `
resource "ovh_cloud_project_database" "db1" {
	service_name = "%s"
	description  = "%s"
	engine       = "postgresql"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database" "db2" {
	service_name = "%s"
	description  = "%s"
	engine       = "opensearch"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_integration" "integration" {
	service_name = ovh_cloud_project_database.db1.service_name
	engine		 = ovh_cloud_project_database.db1.engine
	cluster_id   = ovh_cloud_project_database.db1.id
	source_service_id	= ovh_cloud_project_database.db1.id
	destination_service_id	= ovh_cloud_project_database.db2.id
	type = "opensearchLogs"
}
`

func TestAccCloudProjectDatabaseIntegration_basic(t *testing.T) {
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
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_integration.integration", "source_service_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_integration.integration", "destination_service_id"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_integration.integration", "type", "opensearchLogs"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_integration.integration", "status", "READY"),
				),
			},
		},
	})
}
