package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccHostingPrivateDatabaseUserGrantBasic = `
resource "ovh_hosting_privatedatabase_user_grant" "description" {
    service_name  = "%s"
    user_name     = "%s"
    database_name = "%s"
    grant         = "%s"
}
`

func init() {
	resource.AddTestSweepers("ovh_hosting_privatedatabase_user_grant", &resource.Sweeper{
		Name: "ovh_hosting_privatedatabase_user_grant",
		F:    testSweepHostingPrivateDatabaseUserGrant,
	})
}

func testSweepHostingPrivateDatabaseUserGrant(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ds := HostingPrivateDatabaseUserGrant{}
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s/grant", url.PathEscape(serviceName), url.PathEscape(userName))

	if err := client.Get(endpoint, &ds); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// no service_name, nothing to sweep
			return nil
		}

		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := client.Delete(endpoint, nil); err != nil {
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

func TestAccHostingPrivateDatabaseUserGrant_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	databaseName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	grantName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_GRANT_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseUserGrantBasic, serviceName, userName, databaseName, grantName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUserGrant(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user_grant.description",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user_grant.description",
						"user_name",
						userName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user_grant.description",
						"database_name",
						databaseName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user_grant.description",
						"grant",
						grantName,
					),
				),
			},
		},
	})
}
