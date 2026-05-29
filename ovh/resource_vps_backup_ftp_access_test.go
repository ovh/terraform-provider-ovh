package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSBackupFtpAccessResource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	ipBlock := os.Getenv("OVH_VPS_BACKUP_FTP_IP_BLOCK")
	if ipBlock == "" {
		ipBlock = "192.0.2.0/24"
	}
	config := fmt.Sprintf(testAccVPSBackupFtpAccessResourceConfig_Basic, vps, ipBlock)

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
						"ovh_vps_backup_ftp_access.acl", "service_name", vps),
					resource.TestCheckResourceAttr(
						"ovh_vps_backup_ftp_access.acl", "ip_block", ipBlock),
					resource.TestCheckResourceAttr(
						"ovh_vps_backup_ftp_access.acl", "cifs", "true"),
					resource.TestCheckResourceAttr(
						"ovh_vps_backup_ftp_access.acl", "nfs", "false"),
				),
			},
			{
				ResourceName:      "ovh_vps_backup_ftp_access.acl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s|%s", vps, ipBlock),
			},
		},
	})
}

const testAccVPSBackupFtpAccessResourceConfig_Basic = `
resource "ovh_vps_backup_ftp_access" "acl" {
  service_name = "%s"
  ip_block     = "%s"
  cifs         = true
  nfs          = false
  ftp          = false
}
`
