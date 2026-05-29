package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSChangeContactResource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	contact := os.Getenv("OVH_NICHANDLE_TEST")
	if contact == "" {
		t.Skip("OVH_NICHANDLE_TEST is not set, skipping")
	}
	config := fmt.Sprintf(testAccVPSChangeContactResourceConfig_Basic, vps, contact)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVPS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_vps_change_contact.action", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"ovh_vps_change_contact.action", "id"),
				),
			},
		},
	})
}

const testAccVPSChangeContactResourceConfig_Basic = `
resource "ovh_vps_change_contact" "action" {
  service_name  = "%s"
  contact_admin = "%s"
}
`
