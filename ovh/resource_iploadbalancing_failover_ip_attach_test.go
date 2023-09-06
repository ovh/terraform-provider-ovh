package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceIpLoadbalancingFailoverIpAttach = `
resource "ovh_iploadbalancing_failover_ip_attach" "myfailoverip" {
 service_name = "%s"
 ip = "%s"
 to = "%s"
}
`

func TestAccResourceIpLoadbalancingFailoverIpAttach(t *testing.T) {
	iplbService := os.Getenv("OVH_IPLB_SERVICE_TEST")
	ipAddress := os.Getenv("OVH_IP_BLOCK_TEST")
	config := fmt.Sprintf(
		testAccResourceIpLoadbalancingFailoverIpAttach,
		iplbService,
		ipAddress,
		iplbService,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckIpLoadbalancingFailoverIpAttach(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_failover_ip_attach.myfailoverip",
						"service_name",
						iplbService,
					),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_failover_ip_attach.myfailoverip",
						"ip",
						ipAddress,
					),
					resource.TestCheckResourceAttr(
						"ovh_iploadbalancing_failover_ip_attach.myfailoverip",
						"to",
						iplbService,
					),
				),
			},
		},
	})
}
