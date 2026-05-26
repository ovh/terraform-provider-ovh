package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSIpReverseConfig = `
resource "ovh_vps_ip_reverse" "rev" {
  service_name = "%s"
  ip_address   = "%s"
  reverse      = "%s"
}
`

func TestAccVPSIpReverse_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	ip := os.Getenv("OVH_VPS_IP")
	rev := os.Getenv("OVH_VPS_IP_REVERSE")
	if vps == "" || ip == "" || rev == "" {
		t.Skip("OVH_VPS, OVH_VPS_IP and OVH_VPS_IP_REVERSE must be set for this acceptance test")
	}
	config := fmt.Sprintf(testAccVPSIpReverseConfig, vps, ip, rev)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_ip_reverse.rev", "reverse", rev),
					resource.TestCheckResourceAttr(
						"ovh_vps_ip_reverse.rev", "ip_address", ip),
				),
			},
		},
	})
}
