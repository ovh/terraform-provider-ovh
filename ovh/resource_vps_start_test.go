package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsStartConfig = `
resource "ovh_vps_start" "start" {
  service_name = "%s"

  triggers = {
    nonce = "1"
  }
}
`

func TestAccVpsStart_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpsStartConfig, os.Getenv("OVH_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_start.start", "task_state", "done"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_start.start", "task_id"),
				),
			},
		},
	})
}
