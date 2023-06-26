package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccIpMoveConfig = `
resource "ovh_ip_move" "move" {
    ip = "%s"
    to = "%s"
}
`

func TestAccIpMove_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_TEST")
	service := os.Getenv("OVH_IPLB_SERVICE_TEST")

	config := fmt.Sprintf(testAccIpMoveConfig, ip, ip, service)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpMove(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.ip", "ip", ip),
					resource.TestCheckResourceAttr("ovh_ip_move.to", "service", service),
				),
			},
		},
	})
}
