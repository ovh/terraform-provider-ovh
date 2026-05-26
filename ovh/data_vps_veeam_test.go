package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVeeamDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSVeeamDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, "veeam")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_veeam.server", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_veeam.server", "backup"),
				),
			},
		},
	})
}

const testAccVPSVeeamDatasourceConfig_Basic = `
data "ovh_vps_veeam" "server" {
  service_name = "%s"
}
`
