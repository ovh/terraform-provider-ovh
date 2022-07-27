package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_project_database", &resource.Sweeper{
		Name: "ovh_cloud_project_database",
		F:    testSweepCloudProjectDatabase,
	})
}

func testSweepCloudProjectDatabase(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No database to sweep")
		return nil
	}

	engineName := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	if engineName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST is not set. No database to sweep")
		return nil
	}

	databaseIds := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/cloud/project/%s/database/%s", serviceName, engineName), &databaseIds); err != nil {
		return fmt.Errorf("Error calling GET /cloud/project/%s/database/%s:\n\t %q", serviceName, engineName, err)
	}
	for _, databaseId := range databaseIds {
		endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, engineName, databaseId)
		res := &CloudProjectDatabaseResponse{}
		log.Printf("[DEBUG] read database %s from project: %s", databaseId, serviceName)
		if err := client.Get(endpoint, res); err != nil {
			return err
		}
		if !strings.HasPrefix(res.Description, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, engineName, databaseId), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}

var testAccCloudProjectDatabaseConfig = `
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
`

func TestAccCloudProjectDatabase_basic(t *testing.T) {
	description := acctest.RandomWithPrefix(test_prefix)
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectDatabaseConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
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
						"ovh_cloud_project_database.db", "backup_time"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "created_at"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "description", description),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "endpoints.#"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "endpoints.0.component"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "endpoints.0.domain"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "endpoints.0.ssl"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "endpoints.0.ssl_mode"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "engine", engine),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "description", description),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "flavor", flavor),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "maintenance_time"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "network_type"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "nodes.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "nodes.0.region", region),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "plan", "essential"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "version", version),
				),
			},
		},
	})
}
