package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVPSBackupFtpAccess_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	ipBlock := os.Getenv("OVH_VPS_BACKUP_FTP_IP_BLOCK")
	if ipBlock == "" {
		ipBlock = "203.0.113.0/24"
	}
	config := fmt.Sprintf(testAccDataSourceVPSBackupFtpAccessConfig_Basic, vps, ipBlock)

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
						"data.ovh_vps_backup_ftp_access.entry", "service_name", vps),
					resource.TestCheckResourceAttr(
						"data.ovh_vps_backup_ftp_access.entry", "ip_block", ipBlock),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_backup_ftp_access.entry", "last_update"),
				),
			},
		},
	})
}

const testAccDataSourceVPSBackupFtpAccessConfig_Basic = `
data "ovh_vps_backup_ftp_access" "entry" {
  service_name = "%s"
  ip_block     = "%s"
}
`
