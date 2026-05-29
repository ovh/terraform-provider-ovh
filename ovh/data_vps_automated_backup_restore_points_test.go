package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAutomatedBackupRestorePointsDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSAutomatedBackupRestorePointsDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, "automatedBackup")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_vps_automated_backup_restore_points.rp", "restore_points.#"),
			},
		},
	})
}

const testAccVPSAutomatedBackupRestorePointsDatasourceConfig_Basic = `
data "ovh_vps_automated_backup_restore_points" "rp" {
  service_name = "%s"
  state        = "available"
}
`
