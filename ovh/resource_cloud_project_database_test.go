package ovh

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

		err = retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, engineName, databaseId), nil); err != nil {
				return retry.RetryableError(err)
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
		region = "%s"
	}
	ip_restrictions {
		description = "%s"
		ip = "%s"
	}
	flavor = "%s"
}
`

func TestAccCloudProjectDatabase_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	ip := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudProjectDatabaseConfig,
		serviceName,
		description,
		engine,
		version,
		region,
		description,
		ip,
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
						"ovh_cloud_project_database.db", "backup_regions.#"),
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
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_database.db", "ip_restrictions.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "ip_restrictions.0.description", description),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "ip_restrictions.0.ip", ip),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "plan", "essential"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_database.db", "version", version),
				),
			},
		},
	})
}
