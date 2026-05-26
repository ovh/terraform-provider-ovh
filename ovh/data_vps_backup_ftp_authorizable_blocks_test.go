package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSBackupFtpAuthorizableBlocksDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSBackupFtpAuthorizableBlocksDatasourceConfig_Basic, vps)

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
						"data.ovh_vps_backup_ftp_authorizable_blocks.blocks", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_backup_ftp_authorizable_blocks.blocks", "blocks.#"),
				),
			},
		},
	})
}

const testAccVPSBackupFtpAuthorizableBlocksDatasourceConfig_Basic = `
data "ovh_vps_backup_ftp_authorizable_blocks" "blocks" {
  service_name = "%s"
}
`
