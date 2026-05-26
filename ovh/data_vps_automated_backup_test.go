package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSAutomatedBackupDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSAutomatedBackupDatasourceConfig_Basic, vps)

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
					resource.TestCheckResourceAttr(
						"data.ovh_vps_automated_backup.ab", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_automated_backup.ab", "state"),
				),
			},
		},
	})
}

const testAccVPSAutomatedBackupDatasourceConfig_Basic = `
data "ovh_vps_automated_backup" "ab" {
  service_name = "%s"
}
`
