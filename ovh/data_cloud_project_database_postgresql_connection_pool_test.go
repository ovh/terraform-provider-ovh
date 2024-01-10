package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabasePostgresqlConnectionPoolDatasourceConfig_Basic = testAccCloudProjectDatabasePostgresqlConnectionPoolConfig + `

data "ovh_cloud_project_database_postgresql_connection_pool" "connection_pool" {
  service_name = ovh_cloud_project_database_postgresql_connection_pool.connection_pool.service_name
  cluster_id   = ovh_cloud_project_database_postgresql_connection_pool.connection_pool.cluster_id
  name         = ovh_cloud_project_database_postgresql_connection_pool.connection_pool.name
}
`

func TestAccCloudProjectDatabasePostgresqlConnectionPoolDataSource_basic(t *testing.T) {
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
		testAccCloudProjectDatabasePostgresqlConnectionPoolDatasourceConfig_Basic,
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
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "cluster_id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "database_id"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "name", connectionPoolName),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "mode", mode),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "size", strconv.Itoa(size)),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "port"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "ssl_mode"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "uri"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pool.connection_pool", "user_id"),
				),
			},
		},
	})
}
