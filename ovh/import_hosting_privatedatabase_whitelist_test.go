package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccHostingPrivateDatabaseWhitelistImportBasic = `
resource "ovh_hosting_privatedatabase_whitelist" "description" {
    service_name = "%s"
    ip           = "%s"
    name         = "%s"
    sftp         = "%s"
    service      = "%s"
}
`

func TestAccHostingPrivateDatabaseWhitelist_importBasic(t *testing.T) {
	resourceName := "ovh_hosting_privatedatabase_whitelist.description"
	serviceName := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	ip := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST")
	name := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_NAME_TEST")
	service := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SERVICE_TEST")
	sftp := os.Getenv("OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SFTP_TEST")

	config := fmt.Sprintf(testAccHostingPrivateDatabaseWhitelistImportBasic, serviceName, ip, name, service, sftp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabaseWhitelist(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccHostingPrivateDatabaseWhitelistImportId(resourceName),
			},
		},
	})
}

func testAccHostingPrivateDatabaseWhitelistImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ds, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_hosting_privatedatabase_whitelist not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			ds.Primary.Attributes["service_name"],
			ds.Primary.Attributes["ip"],
		), nil
	}
}
