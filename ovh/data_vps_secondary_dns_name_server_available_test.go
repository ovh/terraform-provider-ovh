package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSSecondaryDNSNameServerAvailableDatasourceConfig = `
data "ovh_vps_secondary_dns_name_server_available" "ns" {
  service_name = "%s"
}
`

func TestAccVPSSecondaryDNSNameServerAvailableDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSSecondaryDNSNameServerAvailableDatasourceConfig, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_vps_secondary_dns_name_server_available.ns", "hostname"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_secondary_dns_name_server_available.ns", "ip"),
				),
			},
		},
	})
}
