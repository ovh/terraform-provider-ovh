package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccResourceVPSAutomatedBackupRescheduleConfig = `
resource "ovh_vps_automated_backup_reschedule" "schedule" {
  service_name = "%s"
  schedule     = "%s"
}
`

func TestAccResourceVPSAutomatedBackupReschedule_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS")
	schedule := os.Getenv("OVH_VPS_AUTOMATED_BACKUP_SCHEDULE")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkEnvOrSkip(t, "OVH_VPS_AUTOMATED_BACKUP_SCHEDULE")
			checkVPSOptionSubscribed(t, "automatedBackup")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceVPSAutomatedBackupRescheduleConfig, serviceName, schedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_automated_backup_reschedule.schedule", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"ovh_vps_automated_backup_reschedule.schedule", "schedule", schedule),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_automated_backup_reschedule.schedule", "state"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_automated_backup_reschedule.schedule", "rotation"),
				),
			},
		},
	})
}
