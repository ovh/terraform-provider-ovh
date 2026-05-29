package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVeeamRestorePointDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	id := os.Getenv("OVH_VPS_VEEAM_RESTORE_POINT_ID")
	if id == "" {
		id = "1"
	}
	config := fmt.Sprintf(testAccVPSVeeamRestorePointDatasourceConfig_Basic, vps, id)

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
						"data.ovh_vps_veeam_restore_point.point", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_veeam_restore_point.point", "creation_time"),
				),
			},
		},
	})
}

const testAccVPSVeeamRestorePointDatasourceConfig_Basic = `
data "ovh_vps_veeam_restore_point" "point" {
  service_name     = "%s"
  restore_point_id = %s
}
`
