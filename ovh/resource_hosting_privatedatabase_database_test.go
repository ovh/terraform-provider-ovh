package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccHostingPrivateDatabaseDatabaseConfig = `
resource "ovh_hosting_privatedatabase_database" "database" {
    service_name  = "%s"
    database_name = "%s"
}
`

func init() {
	resource.AddTestSweepers("ovh_hosting_privatedatabase_database", &resource.Sweeper{
		Name: "ovh_hosting_privatedatabase_database",
		F:    testSweepHostingPrivateDatabaseDatabase,
	})
}

func testSweepHostingPrivateDatabaseDatabase(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ds := HostingPrivateDatabaseDatabase{}
	testServiceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	testDatabaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/database/%s", url.PathEscape(testServiceName), url.PathEscape(testDatabaseName))

	if err := client.Get(endpoint, &ds); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// no service_name, nothing to sweep
			return nil
		}

		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Deleting database %v", ds)
		if err := client.Delete(endpoint, &ds); err != nil {
			return resource.RetryableError(err)
		}
		// Successful delete
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func TestAccHostingPrivateDatabaseDatabase_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")

	config := fmt.Sprintf(
		testAccHostingPrivateDatabaseDatabaseConfig,
		serviceName,
		databaseName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_database.database",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_database.database",
						"database_name",
						databaseName,
					),
				),
			},
		},
	})
}
