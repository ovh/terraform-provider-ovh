package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceIpMove = `
data "ovh_ip_move" "mymoveip" {
 ip = "%s"
}
`

func TestAccDataSourceIpMove(t *testing.T) {
	ipAddress := os.Getenv("OVH_IP_TEST")
	config := fmt.Sprintf(
		testAccDataSourceIpMove,
		ipAddress,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckIpMove(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_ip_move.mymoveip",
						"ip",
						ipAddress,
					),
				),
			},
		},
	})
}
