package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIpLoadbalancingDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	config := fmt.Sprintf(testAccIpLoadbalancingDatasourceConfig_Basic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "id", serviceName),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingDataSource_statevrack(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE")
	config := fmt.Sprintf(testAccIpLoadbalancingDatasourceConfig_StateAndVrack, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "id", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "state", "ok"),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
				),
			},
		},
	})
}

const testAccIpLoadbalancingDatasourceConfig_Basic = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}
`
const testAccIpLoadbalancingDatasourceConfig_StateAndVrack = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
  state = "ok"
  vrack_eligibility = true
}
`
