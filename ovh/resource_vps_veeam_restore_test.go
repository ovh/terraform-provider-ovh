package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVeeamRestoreResource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	id := os.Getenv("OVH_VPS_VEEAM_RESTORE_POINT_ID")
	if id == "" {
		id = "1"
	}
	config := fmt.Sprintf(testAccVPSVeeamRestoreResourceConfig_Basic, vps, id)

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
						"ovh_vps_veeam_restore.r", "service_name", vps),
					resource.TestCheckResourceAttr(
						"ovh_vps_veeam_restore.r", "export", "nfs"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_veeam_restore.r", "state"),
				),
			},
		},
	})
}

const testAccVPSVeeamRestoreResourceConfig_Basic = `
resource "ovh_vps_veeam_restore" "r" {
  service_name     = "%s"
  restore_point_id = %s
  full             = false
  export           = "nfs"
}
`
