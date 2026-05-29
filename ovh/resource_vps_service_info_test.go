package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSServiceInfoResource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSServiceInfoResourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_service_info.info", "service_name", vps),
					resource.TestCheckResourceAttr(
						"ovh_vps_service_info.info", "renew_automatic", "false"),
				),
			},
		},
	})
}

const testAccVPSServiceInfoResourceConfig_Basic = `
resource "ovh_vps_service_info" "info" {
  service_name    = "%s"
  renew_automatic = false
}
`
