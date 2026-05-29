package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSVeeamRestoredBackupDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSVeeamRestoredBackupDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			checkVPSOptionSubscribed(t, "veeam")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_veeam_restored_backup.current", "state"),
				),
			},
		},
	})
}

const testAccVPSVeeamRestoredBackupDatasourceConfig_Basic = `
data "ovh_vps_veeam_restored_backup" "current" {
  service_name = "%s"
}
`
