package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSOptionsDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSOptionsDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_options.opts", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_options.opts", "options.#"),
				),
			},
		},
	})
}

const testAccVPSOptionsDatasourceConfig_Basic = `
data "ovh_vps_options" "opts" {
  service_name = "%s"
}
`
