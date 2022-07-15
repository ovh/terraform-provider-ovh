package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine 		 = "%s"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

data "ovh_cloud_project_database" "db" {
  service_name = ovh_cloud_project_database.db.service_name
  engine 	   = ovh_cloud_project_database.db.engine
  cluster_id   = ovh_cloud_project_database.db.id
}
`

func TestAccCloudProjectDatabaseDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseDatasourceConfig_Basic,
		serviceName,
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "backup_time"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "created_at"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "description", description),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "endpoints.#"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "endpoints.0.component"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "endpoints.0.domain"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "endpoints.0.ssl"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "endpoints.0.ssl_mode"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "engine", engine),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "flavor", flavor),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "maintenance_time"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "network_type"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "nodes.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "nodes.0.region", region),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "plan", "essential"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database.db", "status"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database.db", "version", version),
				),
			},
		},
	})
}
