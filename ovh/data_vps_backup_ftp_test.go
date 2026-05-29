package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSBackupFtpDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSBackupFtpDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, "ftpbackup")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_backup_ftp.backup", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_backup_ftp.backup", "ftp_backup_name"),
				),
			},
		},
	})
}

const testAccVPSBackupFtpDatasourceConfig_Basic = `
data "ovh_vps_backup_ftp" "backup" {
  service_name = "%s"
}
`
