package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testAccIpLoadbalancingTcpFarmConfig = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}
resource "ovh_iploadbalancing_tcp_farm" "testfarm" {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port             = "%d"
  zone             = "%s"
  balance 		   = "roundrobin"
  probe {
        interval = 30
        type = "oco"
  }
}
`
)

func TestAccIpLoadbalancingTcpFarmBasicCreate(t *testing.T) {
	displayName1 := acctest.RandomWithPrefix(test_prefix)
	displayName2 := acctest.RandomWithPrefix(test_prefix)
	config1 := fmt.Sprintf(
		testAccIpLoadbalancingTcpFarmConfig,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName1,
		12345,
		"all",
	)
	config2 := fmt.Sprintf(
		testAccIpLoadbalancingTcpFarmConfig,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName2,
		12346,
		"all",
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "display_name", displayName1),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "zone", "all"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "port", "12345"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "probe.0.interval", "30"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "display_name", displayName2),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "zone", "all"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "port", "12346"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "probe.0.interval", "30"),
				),
			},
		},
	})
}
