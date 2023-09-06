package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceIpLoadbalancingFailoverIpAttach = `
data "ovh_iploadbalancing_failover_ip_attach" "myfailoverip" {
 service_name = "%s"
 ip = "%s"
}
`

func TestAccDataSourceIpLoadbalancingFailoverIpAttach(t *testing.T) {
	iplbService := os.Getenv("OVH_IPLB_SERVICE_TEST")
	ipAddress := os.Getenv("OVH_IP_BLOCK_TEST")
	config := fmt.Sprintf(
		testAccDataSourceIpLoadbalancingFailoverIpAttach,
		iplbService,
		ipAddress,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckIpLoadbalancingFailoverIpAttach(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing_failover_ip_attach.myfailoverip",
						"service_name",
						iplbService,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing_failover_ip_attach.myfailoverip",
						"ip",
						ipAddress,
					),
				),
			},
		},
	})
}
