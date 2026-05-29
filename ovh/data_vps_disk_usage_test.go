package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDiskUsageDatasourceConfig = `
data "ovh_vps_disk_usage" "usage" {
  service_name = "%s"
  disk_id      = %s
  type         = "used"
}
`

func TestAccVPSDiskUsageDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	diskID := os.Getenv("OVH_VPS_DISK_ID")
	config := fmt.Sprintf(testAccVPSDiskUsageDatasourceConfig, vps, diskID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPSDisk(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk_usage.usage", "unit"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk_usage.usage", "value"),
				),
			},
		},
	})
}
