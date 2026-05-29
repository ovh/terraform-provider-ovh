package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAutomatedBackupRestoreResource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	restorePoint := os.Getenv("OVH_VPS_RESTORE_POINT")
	if restorePoint == "" {
		t.Skip("OVH_VPS_RESTORE_POINT not set")
	}
	config := fmt.Sprintf(testAccVPSAutomatedBackupRestoreResourceConfig_Basic, vps, restorePoint)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, "automatedBackup")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_vps_automated_backup_restore.r", "task_id"),
				),
			},
		},
	})
}

const testAccVPSAutomatedBackupRestoreResourceConfig_Basic = `
resource "ovh_vps_automated_backup_restore" "r" {
  service_name  = "%s"
  restore_point = "%s"
  type          = "file"
}
`
