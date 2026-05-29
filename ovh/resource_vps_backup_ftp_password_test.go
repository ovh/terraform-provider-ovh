package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVPSBackupFtpPassword_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccResourceVPSBackupFtpPasswordConfig_Basic, vps)

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
						"ovh_vps_backup_ftp_password.rotate", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_backup_ftp_password.rotate", "task_id"),
				),
			},
		},
	})
}

const testAccResourceVPSBackupFtpPasswordConfig_Basic = `
resource "ovh_vps_backup_ftp_password" "rotate" {
  service_name = "%s"

  triggers = {
    rotation = "1"
  }
}
`
