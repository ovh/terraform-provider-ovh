package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAutomatedBackupAttachedDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSAutomatedBackupAttachedDatasourceConfig_Basic, vps)

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
					"data.ovh_vps_automated_backup_attached.att", "attached_backups.#"),
			},
		},
	})
}

const testAccVPSAutomatedBackupAttachedDatasourceConfig_Basic = `
data "ovh_vps_automated_backup_attached" "att" {
  service_name = "%s"
}
`
