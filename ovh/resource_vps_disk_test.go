package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPreCheckVPSDisk(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VPS")
	checkEnvOrSkip(t, "OVH_VPS_DISK_ID")
}

const testAccVPSDiskResourceConfig = `
resource "ovh_vps_disk" "disk" {
  service_name             = "%s"
  disk_id                  = %s
  monitoring               = true
  low_free_space_threshold = 1024
}
`

func TestAccResourceVPSDisk_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	diskID := os.Getenv("OVH_VPS_DISK_ID")
	config := fmt.Sprintf(testAccVPSDiskResourceConfig, vps, diskID)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPSDisk(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps_disk.disk", "service_name", vps),
					resource.TestCheckResourceAttr("ovh_vps_disk.disk", "monitoring", "true"),
					resource.TestCheckResourceAttr("ovh_vps_disk.disk", "low_free_space_threshold", "1024"),
					resource.TestCheckResourceAttrSet("ovh_vps_disk.disk", "type"),
					resource.TestCheckResourceAttrSet("ovh_vps_disk.disk", "state"),
				),
			},
		},
	})
}
