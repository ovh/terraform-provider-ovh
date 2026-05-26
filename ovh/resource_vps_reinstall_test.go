package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSReinstall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			checkEnvOrSkip(t, "OVH_VPS")
			checkEnvOrSkip(t, "OVH_VPS_TEMPLATE_ID")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSReinstallConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_reinstall.reinstall", "task_state", "done"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_reinstall.reinstall", "task_id"),
				),
			},
		},
	})
}

func testAccVPSReinstallConfig() string {
	serviceName := os.Getenv("OVH_VPS")
	templateID := os.Getenv("OVH_VPS_TEMPLATE_ID")
	return fmt.Sprintf(`
resource "ovh_vps_reinstall" "reinstall" {
  service_name = "%s"
  template_id  = %s
  language     = "en"
}
`, serviceName, templateID)
}
