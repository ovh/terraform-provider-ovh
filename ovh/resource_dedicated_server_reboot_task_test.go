package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServerReboot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerRebootConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "monitoring", "true"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_update.server", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reboot_task.server_reboot", "function", "hardReboot"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reboot_task.server_reboot", "comment", "Reboot asked"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_reboot_task.server_reboot", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerRebootConfig() string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	return fmt.Sprintf(
		testAccDedicatedServerRebootConfig_Basic,
		dedicated_server,
	)

}

const testAccDedicatedServerRebootConfig_Basic = `
data ovh_dedicated_server_boots "rescue" {
  service_name = "%s"
  boot_type    = "rescue"
  kernel       = "rescue64-pro"
}

resource ovh_dedicated_server_update "server" {
  service_name = data.ovh_dedicated_server_boots.rescue.service_name
  boot_id      = data.ovh_dedicated_server_boots.rescue.result[0]
  monitoring   = true
  state        = "ok"
}

resource ovh_dedicated_server_reboot_task "server_reboot" {
  service_name = data.ovh_dedicated_server_boots.rescue.service_name

  keepers = [
     ovh_dedicated_server_update.server.boot_id,
  ]
}
`
