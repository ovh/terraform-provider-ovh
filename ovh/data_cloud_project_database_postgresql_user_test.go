package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectDatabasePostgresqlUserDatasourceConfig_Basic = `
resource "ovh_cloud_project_database" "db" {
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

resource "ovh_cloud_project_database_postgresql_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
	roles = ["%s"]
}

data "ovh_cloud_project_database_postgresql_user" "user" {
  service_name = ovh_cloud_project_database_postgresql_user.user.service_name
  cluster_id   = ovh_cloud_project_database_postgresql_user.user.cluster_id
  name         = ovh_cloud_project_database_postgresql_user.user.name
}
`

func TestAccCloudProjectDatabasePostgresqlUserDataSource_basic(t *testing.T) {
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

	config := fmt.Sprintf(
		testAccCloudProjectDatabasePostgresqlUserDatasourceConfig_Basic,
		serviceName,
		description,
		version,
		region,
		flavor,
		name,
		replication,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_user.user", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_user.user", "roles.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_postgresql_user.user", "roles.0", replication),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_user.user", "status"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_postgresql_user.user", "name", name,
					),
				),
			},
		},
	})
}
