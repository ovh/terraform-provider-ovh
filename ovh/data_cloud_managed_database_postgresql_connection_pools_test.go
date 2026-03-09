package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudManagedDatabasePostgresqlConnectionPoolsDatasourceConfig = testAccCloudManagedDatabasePostgresqlConnectionPoolConfig + `

data "ovh_cloud_managed_database_postgresql_connection_pools" "connection_pools" {
  service_name = ovh_cloud_project_database_postgresql_connection_pool.connection_pool.service_name
  cluster_id   = ovh_cloud_project_database_postgresql_connection_pool.connection_pool.cluster_id
}
`

func TestAccCloudManagedDatabasePostgresqlConnectionPoolsDataSource_basic(t *testing.T) {
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
		testAccCloudManagedDatabasePostgresqlConnectionPoolsDatasourceConfig,
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
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_database_postgresql_connection_pools.connection_pools",
						"connection_pool_ids.#",
					),
				),
			},
		},
	})
}
