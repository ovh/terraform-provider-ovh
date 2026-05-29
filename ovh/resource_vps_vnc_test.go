package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVnc_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSVncConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_vnc.session", "service_name", vps),
					resource.TestCheckResourceAttr(
						"ovh_vps_vnc.session", "protocol", "VNC"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_vnc.session", "host"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_vnc.session", "port"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_vnc.session", "password"),
				),
			},
		},
	})
}

const testAccVPSVncConfig_Basic = `
resource "ovh_vps_vnc" "session" {
  service_name = "%s"
  protocol     = "VNC"
}
`
