package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDomainZoneDataSource_basic(t *testing.T) {
	zoneName := os.Getenv("OVH_ZONE")
	config := fmt.Sprintf(testAccDomainZoneDatasourceConfig_Basic, zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDomain(t); testAccCheckDomainZoneExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainZoneHasNameServers("data.ovh_domain_zone.rootzone", t),
					resource.TestCheckResourceAttr(
						"data.ovh_domain_zone.rootzone", "id", zoneName),
				),
			},
		},
	})
}

func testAccCheckDomainZoneHasNameServers(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		nsAttributeCount, err := strconv.Atoi(rs.Primary.Attributes["name_servers.#"])
		if err != nil || nsAttributeCount < 1 {
			return fmt.Errorf("No Name servers are set")
		}

		return nil
	}
}

const testAccDomainZoneDatasourceConfig_Basic = `
data "ovh_domain_zone" "rootzone" {
  name = "%s"
}
`
