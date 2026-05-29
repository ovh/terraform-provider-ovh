package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVeeamRestorePointsDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSVeeamRestorePointsDatasourceConfig_Basic, vps)

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
						"data.ovh_vps_veeam_restore_points.points", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_veeam_restore_points.points", "result.#"),
				),
			},
		},
	})
}

const testAccVPSVeeamRestorePointsDatasourceConfig_Basic = `
data "ovh_vps_veeam_restore_points" "points" {
  service_name = "%s"
}
`
