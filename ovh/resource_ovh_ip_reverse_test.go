package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccIpReverseConfig = fmt.Sprintf(`
resource "ovh_ip_reverse" "reverse" {
    ip = "%s"
    ipreverse = "%s"
    reverse = "%s"
}
`, os.Getenv("OVH_IP_BLOCK"), os.Getenv("OVH_IP"), os.Getenv("OVH_IP_REVERSE"))

func TestAccIpReverse_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckIp(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpReverseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpReverseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpReverseExists("ovh_ip_reverse.reverse", t),
				),
			},
		},
	})
}

func testAccCheckIpReverseExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip"] == "" {
			return fmt.Errorf("No IP block is set")
		}

		if rs.Primary.Attributes["ipreverse"] == "" {
			return fmt.Errorf("No IP is set")
		}

		return resourceOvhIpReverseExists(rs.Primary.Attributes["ip"], rs.Primary.Attributes["ipreverse"], config.OVHClient)
	}
}

func testAccCheckIpReverseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_ip_reverse" {
			continue
		}

		err := resourceOvhIpReverseExists(rs.Primary.Attributes["ip"], rs.Primary.Attributes["ipreverse"], config.OVHClient)
		if err == nil {
			return fmt.Errorf("IP Reverse still exists")
		}
	}
	return nil
}
