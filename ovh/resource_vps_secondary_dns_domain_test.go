package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSSecondaryDNSDomainConfig = `
resource "ovh_vps_secondary_dns_domain" "dom" {
  service_name = "%s"
  domain       = "%s"
  ip           = "%s"
}
`

func TestAccVPSSecondaryDNSDomain_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	domain := os.Getenv("OVH_VPS_SECONDARY_DNS_DOMAIN")
	ip := os.Getenv("OVH_VPS_SECONDARY_DNS_IP")
	config := fmt.Sprintf(testAccVPSSecondaryDNSDomainConfig, vps, domain, ip)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			if domain == "" || ip == "" {
				t.Skip("OVH_VPS_SECONDARY_DNS_DOMAIN and OVH_VPS_SECONDARY_DNS_IP must be set")
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps_secondary_dns_domain.dom", "domain", domain),
					resource.TestCheckResourceAttr("ovh_vps_secondary_dns_domain.dom", "ip", ip),
					resource.TestCheckResourceAttrSet("ovh_vps_secondary_dns_domain.dom", "dns"),
					resource.TestCheckResourceAttrSet("ovh_vps_secondary_dns_domain.dom", "creation_date"),
				),
			},
		},
	})
}
