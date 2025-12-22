package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpLoadbalancingNatIpsDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data ovh_iploadbalancing_nat_ips ips {
						service_name = %q
					}`, serviceName),
				Check: resource.TestCheckResourceAttrSet("data.ovh_iploadbalancing_nat_ips.ips", "nat_ips.#"),
			},
		},
	})
}
