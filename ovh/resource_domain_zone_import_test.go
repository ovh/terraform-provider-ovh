package ovh

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainZoneImport_Basic(t *testing.T) {
	zone := os.Getenv("OVH_ZONE_TEST")
	zoneFileContent := `$TTL 3600\n@\tIN SOA dns11.ovh.net. tech.ovh.net. (2024090445 86400 3600 3600000 60)\n        IN NS     ns11.ovh.net.\nwww        IN TXT     \"3|hey\"\n`
	zoneFileContentUpdated := `$TTL 3600\n@\tIN SOA dns11.ovh.net. tech.ovh.net. (2024090445 86400 3600 3600000 60)\n        IN NS     ns11.ovh.net.\nwww        IN TXT     \"3|hello\"\n`

	config := `
	resource "ovh_domain_zone_import" "import" {
		zone_name = "%s"
		zone_file = "%s"
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(config, zone, zoneFileContent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_import.import", "zone_name", zone),
					resource.TestCheckResourceAttrWith("ovh_domain_zone_import.import", "exported_content", func(value string) error {
						if !strings.Contains(value, "hey") {
							return errors.New("unexpected content in zone export")
						}
						return nil
					}),
				),
			},
			{
				Config: fmt.Sprintf(config, zone, zoneFileContentUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_import.import", "zone_name", zone),
					resource.TestCheckResourceAttrWith("ovh_domain_zone_import.import", "exported_content", func(value string) error {
						if !strings.Contains(value, "hello") {
							return errors.New("unexpected content in zone export")
						}
						return nil
					}),
				),
			},
		},
	})
}
