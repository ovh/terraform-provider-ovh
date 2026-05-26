package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSSecondaryDNSDomainDatasourceConfig = `
data "ovh_vps_secondary_dns_domain" "dom" {
  service_name = "%s"
  domain       = "%s"
}
`

func TestAccVPSSecondaryDNSDomainDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	domain := os.Getenv("OVH_VPS_SECONDARY_DNS_DOMAIN")
	config := fmt.Sprintf(testAccVPSSecondaryDNSDomainDatasourceConfig, vps, domain)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			if domain == "" {
				t.Skip("OVH_VPS_SECONDARY_DNS_DOMAIN must be set")
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vps_secondary_dns_domain.dom", "domain", domain),
					resource.TestCheckResourceAttrSet("data.ovh_vps_secondary_dns_domain.dom", "dns"),
				),
			},
		},
	})
}
