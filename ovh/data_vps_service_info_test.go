package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSServiceInfoDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSServiceInfoDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_service_info.info", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_service_info.info", "status"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_service_info.info", "renewal_type"),
				),
			},
		},
	})
}

const testAccVPSServiceInfoDatasourceConfig_Basic = `
data "ovh_vps_service_info" "info" {
  service_name = "%s"
}
`
