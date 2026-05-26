package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccResourceVPSSnapshotRevertBasic = `
resource "ovh_vps_snapshot_revert" "test" {
  service_name = "%s"
}
`

func TestAccResourceVPSSnapshotRevert_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME")
	config := fmt.Sprintf(testAccResourceVPSSnapshotRevertBasic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkEnvOrSkip(t, "OVH_VPS_SNAPSHOT_SERVICE_NAME")
			checkVPSServiceOptionSubscribed(t, os.Getenv("OVH_VPS_SNAPSHOT_SERVICE_NAME"), "snapshot")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_snapshot_revert.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_snapshot_revert.test", "task_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_snapshot_revert.test", "task_state"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_snapshot_revert.test", "reverted_at"),
				),
			},
		},
	})
}
