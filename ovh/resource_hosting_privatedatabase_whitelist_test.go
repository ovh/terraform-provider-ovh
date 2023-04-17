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

const testAccHostingPrivateDatabaseWhitelistBasic = `
resource "ovh_hosting_privatedatabase_whitelist" "description" {
    service_name = "%s"
    ip           = "%s"
    name         = "%s"
    sftp         = "%s"
    service      = "%s"
}
`

func init() {
	resource.AddTestSweepers("ovh_hosting_privatedatabase_whitelist", &resource.Sweeper{
		Name: "ovh_hosting_privatedatabase_whitelist",
		F:    testSweepHostingPrivateDatabaseWhitelist,
	})
}

func testSweepHostingPrivateDatabaseWhitelist(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ds := HostingPrivateDatabaseWhitelist{}
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	ip := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST")

	endpoint := fmt.Sprintf("/hosting/privateDatabase/%s/whitelist/%s", url.PathEscape(serviceName), url.PathEscape(ip))

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

func TestAccHostingPrivateDatabaseWhitelist_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	ip := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST")
	name := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_NAME_TEST")
	service := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SERVICE_TEST")
	sftp := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SFTP_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseWhitelistBasic, serviceName, ip, name, service, sftp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseWhitelist(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_hosting_privatedatabase_whitelist.description",
						"name",
					),
					resource.TestCheckResourceAttrSet(
						"ovh_hosting_privatedatabase_whitelist.description",
						"service",
					),
				),
			},
		},
	})
}
