package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectDatabaseDatabaseConfig = `
resource "ovh_cloud_project_database" "db" {
	service_name = "%s"
	description  = "%s"
	engine       = "%s"
	version      = "%s"
	plan         = "essential"
	nodes {
		region     = "%s"
	}
	flavor = "%s"
}

resource "ovh_cloud_project_database_database" "database" {
	service_name = ovh_cloud_project_database.db.service_name
	engine		 = ovh_cloud_project_database.db.engine
	cluster_id   = ovh_cloud_project_database.db.id
	name		 = "%s"
}
`

func init() {
	resource.AddTestSweepers("ovh_cloud_project_database_database", &resource.Sweeper{
		Name: "ovh_cloud_project_database_database",
		F:    testSweepCloudDatabaseDatabase,
	})
}

func testSweepCloudDatabaseDatabase(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	databases := []string{"cassandra", "grafana", "kafka", "kafkaConnect", "kafkaMirrorMaker", "m3aggregator", "m3db", "mongodb", "mysql", "opensearch", "postgresql", "redis"}

	for _, database := range databases {

		idsToSweep := []string{}
		endpoint := fmt.Sprintf("/cloud/project/%s/database/%s", serviceName, database)
		if err := client.Get(endpoint, &idsToSweep); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if len(idsToSweep) == 0 {
			log.Printf("[INFO] No %s database  to sweep", database)
		}

		for _, id := range idsToSweep {
			log.Printf("[INFO] sweeping %s database with id %s", database, id)
			endpoint = fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, database, id)
			if err := client.Delete(endpoint, nil); err != nil {
				return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
			}
		}
	}
	return nil

}
func TestAccCloudProjectDatabaseDatabase_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	description := acctest.RandomWithPrefix(test_prefix)
	name := "mydatabase"

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseDatabaseConfig,
		serviceName,
		description,
		engine,
		version,
		region,
		flavor,
		name,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database_database.database", "default"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database_database.database", "name", name,
					),
				),
			},
		},
	})
}
