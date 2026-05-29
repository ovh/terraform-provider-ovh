package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSIpDatasourceConfig_Basic = `
data "ovh_vps_ip" "ip" {
  service_name = "%s"
  ip_address   = "%s"
}
`

func TestAccVPSIpDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	ip := os.Getenv("OVH_VPS_IP")
	if vps == "" || ip == "" {
		t.Skip("OVH_VPS and OVH_VPS_IP must be set for this acceptance test")
	}
	config := fmt.Sprintf(testAccVPSIpDatasourceConfig_Basic, vps, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_ip.ip", "ip_address", ip),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_ip.ip", "version"),
				),
			},
		},
	})
}
