package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsStopConfig = `
resource "ovh_vps_stop" "stop" {
  service_name = "%s"

  triggers = {
    nonce = "1"
  }
}
`

func TestAccVpsStop_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpsStopConfig, os.Getenv("OVH_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_stop.stop", "task_state", "done"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_stop.stop", "task_id"),
				),
			},
		},
	})
}
