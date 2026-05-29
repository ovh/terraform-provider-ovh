package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSSecondaryDNSDomainsDatasourceConfig = `
data "ovh_vps_secondary_dns_domains" "doms" {
  service_name = "%s"
}
`

func TestAccVPSSecondaryDNSDomainsDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSSecondaryDNSDomainsDatasourceConfig, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_vps_secondary_dns_domains.doms",
					"result.#",
				),
			},
		},
	})
}
