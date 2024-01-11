package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabaseM3dbNamespaceDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "m3db"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_m3db_namespace" "namespace" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
	resolution 	 = "%s"
	retention_period_duration = "%s"
}

data "ovh_cloud_project_database_m3db_namespace" "namespace" {
  service_name = ovh_cloud_project_database_m3db_namespace.namespace.service_name
  cluster_id   = ovh_cloud_project_database_m3db_namespace.namespace.cluster_id
  name     = ovh_cloud_project_database_m3db_namespace.namespace.name
}
`

func TestAccCloudProjectDatabaseM3dbNamespaceDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_M3DB_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "mynamespace"
	resolution := "P2D"
	periodDuration := "P2D"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseM3dbNamespaceDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		resolution,
		periodDuration,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_m3db_namespace.namespace", "name", name),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_m3db_namespace.namespace", "resolution", resolution),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_m3db_namespace.namespace", "retention_period_duration", periodDuration),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_m3db_namespace.namespace", "type", "aggregated"),
				),
			},
		},
	})
}
