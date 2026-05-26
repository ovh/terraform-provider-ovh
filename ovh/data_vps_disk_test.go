package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDiskDatasourceConfig = `
data "ovh_vps_disk" "disk" {
  service_name = "%s"
  disk_id      = %s
}
`

func TestAccVPSDiskDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	diskID := os.Getenv("OVH_VPS_DISK_ID")
	config := fmt.Sprintf(testAccVPSDiskDatasourceConfig, vps, diskID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPSDisk(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vps_disk.disk", "service_name", vps),
					resource.TestCheckResourceAttr("data.ovh_vps_disk.disk", "disk_id", diskID),
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk.disk", "type"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk.disk", "state"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_disk.disk", "size"),
				),
			},
		},
	})
}
