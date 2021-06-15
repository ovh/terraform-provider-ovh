package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPSDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps.server", "name", vps),
					resource.TestCheckResourceAttr(
						"data.ovh_vps.server", "service_name", vps),
				),
			},
		},
	})
}

const testAccVPSDatasourceConfig_Basic = `
data "ovh_vps" "server" {
  service_name  = "%s"
}
`
