package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDiskMonitoringDatasourceConfig = `
data "ovh_vps_disk_monitoring" "mon" {
  service_name = "%s"
  disk_id      = %s
  period       = "lastday"
  type         = "cpu:used"
}
`

func TestAccVPSDiskMonitoringDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	diskID := os.Getenv("OVH_VPS_DISK_ID")
	config := fmt.Sprintf(testAccVPSDiskMonitoringDatasourceConfig, vps, diskID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPSDisk(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk_monitoring.mon", "unit"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk_monitoring.mon", "values.#"),
				),
			},
		},
	})
}
