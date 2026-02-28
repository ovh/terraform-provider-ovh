package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedDatabasePostgresqlConnectionPoolConfig = `
resource "ovh_cloud_managed_database" "db" {
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

resource "ovh_cloud_managed_database_database" "database" {
  service_name  = ovh_cloud_project_database.db.service_name
  engine        = ovh_cloud_project_database.db.engine
  cluster_id    = ovh_cloud_project_database.db.id
  name          = "mydatabase"
}

resource "ovh_cloud_managed_database_postgresql_user" "user" {
	service_name = ovh_cloud_project_database.db.service_name
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
	roles 		 = ["%s"]
}

resource "ovh_cloud_managed_database_postgresql_connection_pool" "connection_pool" {
  service_name = ovh_cloud_project_database.db.service_name
  cluster_id   = ovh_cloud_project_database.db.id
  database_id = ovh_cloud_project_database_database.database.id
  name = "%s"
  user_id = ovh_cloud_project_database_postgresql_user.user.id
  mode = "%s"
  size = %d
}
`

func TestAccCloudManagedDatabasePostgresqlConnectionPool_basic(t *testing.T) {
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
		testAccCloudManagedDatabasePostgresqlConnectionPoolConfig,
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "cluster_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "database_id"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "name", connectionPoolName),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "mode", mode),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "size", strconv.Itoa(size)),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "port"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "ssl_mode"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "uri"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database_postgresql_connection_pool.connection_pool", "user_id"),
				),
			},
		},
	})
}
