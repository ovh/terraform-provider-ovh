package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSOptionDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	option := os.Getenv("OVH_VPS_OPTION")
	if option == "" {
		option = "snapshot"
	}
	config := fmt.Sprintf(testAccVPSOptionDatasourceConfig_Basic, vps, option)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, option)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_option.opt", "service_name", vps),
					resource.TestCheckResourceAttr(
						"data.ovh_vps_option.opt", "option", option),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_option.opt", "state"),
				),
			},
		},
	})
}

const testAccVPSOptionDatasourceConfig_Basic = `
data "ovh_vps_option" "opt" {
  service_name = "%s"
  option       = "%s"
}
`
