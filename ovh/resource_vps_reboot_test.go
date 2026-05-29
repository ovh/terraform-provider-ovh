package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsRebootConfig = `
resource "ovh_vps_reboot" "reboot" {
  service_name = "%s"

  triggers = {
    nonce = "1"
  }
}
`

func TestAccVpsReboot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpsRebootConfig, os.Getenv("OVH_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_reboot.reboot", "task_state", "done"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_reboot.reboot", "task_id"),
				),
			},
		},
	})
}
