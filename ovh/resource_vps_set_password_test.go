package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsSetPasswordConfig = `
resource "ovh_vps_set_password" "pwd" {
  service_name = "%s"

  triggers = {
    nonce = "1"
  }
}
`

func TestAccVpsSetPassword_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpsSetPasswordConfig, os.Getenv("OVH_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_set_password.pwd", "task_state", "done"),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_set_password.pwd", "task_id"),
					resource.TestCheckResourceAttr(
						"ovh_vps_set_password.pwd", "password_sent_via_email", "true"),
				),
			},
		},
	})
}
