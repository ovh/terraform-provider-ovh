package ovh

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_managed_database", &resource.Sweeper{
		Name: "ovh_cloud_managed_database",
		F:    testSweepCloudManagedDatabase,
	})
}

func testSweepCloudManagedDatabase(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No database to sweep")
		return nil
	}
	databases := []string{"cassandra", "grafana", "kafka", "kafkaConnect", "kafkaMirrorMaker", "m3aggregator", "m3db", "mongodb", "mysql", "opensearch", "postgresql", "redis", "valkey"}

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
			endpoint = fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, database, id)
			res := &CloudProjectDatabaseResponse{}
			log.Printf("[DEBUG] read database %s from project: %s", id, serviceName)
			if err := client.Get(endpoint, res); err != nil {
				return err
			}
			if !strings.HasPrefix(res.Description, test_prefix) {
				continue
			}

			err = retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
				if err := client.Delete(fmt.Sprintf("/cloud/project/%s/database/%s/%s", serviceName, database, id), nil); err != nil {
					return retry.RetryableError(err)
				}
				// Successful delete
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil

}

var testAccCloudManagedDatabaseConfig = `
resource "ovh_cloud_managed_database" "db" {
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
	maintenance_time = "01:00:00"
}
`

func TestAccCloudManagedDatabase_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	engine := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
	ip := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST")
	description := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(
		testAccCloudManagedDatabaseConfig,
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
						"ovh_cloud_managed_database.db", "backup_regions.#"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "backup_time"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "created_at"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "description", description),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "endpoints.#"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "endpoints.0.component"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "endpoints.0.domain"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "endpoints.0.ssl"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "engine", engine),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "description", description),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "flavor", flavor),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "maintenance_time"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "network_type"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "nodes.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "nodes.0.region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "ip_restrictions.#"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "ip_restrictions.0.description", description),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "ip_restrictions.0.ip", ip),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "plan", "essential"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "version", version),
					resource.TestCheckResourceAttr(
						"ovh_cloud_managed_database.db", "maintenance_time", "01:00:00"),
				),
			},
		},
	})
}

func TestAccCloudManagedDatabase_invalidBackupTime(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_KAFKA_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")

	config := fmt.Sprintf(`
		resource "ovh_cloud_managed_database" "db" {
			service_name = "%s"
			description  = "Invalid kafka"
			engine       = "kafka"
			version      = "%s"
			plan         = "business"
			backup_time  = "00:11:00"
			flavor       = "%s"
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
		}
	`, serviceName, version, flavor, region, region, region)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`backup_time is not customizable for engine \\\"kafka\\\"`),
			},
		},
	})
}

func TestAccCloudManagedDatabase_invalidBackupTimeUpdate(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_OPENSEARCH_VERSION_TEST")
	if version == "" {
		version = os.Getenv("OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	}
	region := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	flavor := os.Getenv("OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")

	config := fmt.Sprintf(`
		resource "ovh_cloud_managed_database" "db" {
			service_name = "%s"
			description  = "Invalid opensearch"
			engine       = "opensearch"
			version      = "%s"
			plan         = "business"
			flavor       = "%s"
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
		}
	`, serviceName, version, flavor, region, region, region)

	configUpdated := fmt.Sprintf(`
		resource "ovh_cloud_managed_database" "db" {
			service_name = "%s"
			description  = "Invalid opensearch"
			engine       = "opensearch"
			version      = "%s"
			plan         = "business"
			backup_time  = "00:11:00"
			flavor       = "%s"
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
			nodes {
				region     = "%s"
			}
		}
	`, serviceName, version, flavor, region, region, region)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDatabaseNoEngine(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_managed_database.db", "backup_time"),
				),
			},
			{
				Config:      configUpdated,
				ExpectError: regexp.MustCompile(`backup_time is not customizable for engine \\\"opensearch\\\"`),
			},
		},
	})
}
