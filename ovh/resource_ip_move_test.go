package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

var testAccIpMoveConfig = `
resource "ovh_ip_move" "move" {
    ip = "%s"
    routed_to {
		service_name = "%s"
	}
}
`

func init() {}

func TestAccIpMove_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_MOVE_TEST")
	routedToServiceName := os.Getenv("OVH_IP_MOVE_SERVICE_NAME_TEST")

	moveConfig := fmt.Sprintf(testAccIpMoveConfig, ip, routedToServiceName)
	parkConfig := fmt.Sprintf(testAccIpMoveConfig, ip, "")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpMove(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: moveConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "ip", ip),
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", routedToServiceName),
				),
			},
			{
				Config: parkConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_ip_move.move", "ip", ip),
					resource.TestCheckResourceAttr("ovh_ip_move.move", "routed_to.0.service_name", ""),
				),
			},
		},
	})
}
