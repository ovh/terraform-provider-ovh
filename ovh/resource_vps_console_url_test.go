package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSConsoleURL_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSConsoleURLConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_console_url.console", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_console_url.console", "url"),
				),
			},
		},
	})
}

const testAccVPSConsoleURLConfig_Basic = `
resource "ovh_vps_console_url" "console" {
  service_name = "%s"
}
`
