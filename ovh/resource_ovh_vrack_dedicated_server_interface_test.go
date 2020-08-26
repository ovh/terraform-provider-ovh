package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccVrackDedicatedServerInterfaceConfig = fmt.Sprintf(`
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}

resource "ovh_vrack_dedicated_server_interface" "vdsi" {
  vrack_id = "%s"
  interface_id = data.ovh_dedicated_server.server.enabled_vrack_vnis[0]
}
`, os.Getenv("OVH_DEDICATED_SERVER"), os.Getenv("OVH_VRACK"))

func TestAccVrackDedicatedServerInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackDedicatedServerInterfacePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackDedicatedServerInterfaceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_server_interface.vdsi", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttrSet("ovh_vrack_dedicated_server_interface.vdsi", "interface_id"),
				),
			},
		},
	})
}

func testAccCheckVrackDedicatedServerInterfacePreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccPreCheckDedicatedServer(t)
}
