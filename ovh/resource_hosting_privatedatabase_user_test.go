package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccHostingPrivateDatabaseUserBasic = `
resource "ovh_hosting_privatedatabase_user" "description" {
    service_name  = "%s"
    password      = "%s"
    user_name     = "%s"
}
`

func init() {
	resource.AddTestSweepers("ovh_hosting_privatedatabase_user", &resource.Sweeper{
		Name: "ovh_hosting_privatedatabase_user",
		F:    testSweepHostingPrivateDatabaseUser,
	})
}

func testSweepHostingPrivateDatabaseUser(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ds := HostingPrivateDatabaseUser{}
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	if serviceName == "" || userName == "" {
		return nil
	}
	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/user/%s", url.PathEscape(serviceName), url.PathEscape(userName))

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

func TestAccHostingPrivateDatabaseUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	password := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
	userName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_USER_TEST")

	config := fmt.Sprintf(
		testAccHostingPrivateDatabaseUserBasic,
		serviceName,
		password,
		userName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseUser(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user.description",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user.description",
						"user_name",
						userName,
					),
					resource.TestCheckResourceAttr(
						"ovh_hosting_privatedatabase_user.description",
						"password",
						password,
					),
				),
			},
		},
	})
}
