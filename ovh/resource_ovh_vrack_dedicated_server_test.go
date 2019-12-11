package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAccVrackDedicatedServerConfig = fmt.Sprintf(`
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}

resource "ovh_vrack_dedicated_server" "vds" {
  vrack_id = "%s"
  server_id = data.ovh_dedicated_server.server.id
}
`, os.Getenv("OVH_DEDICATED_SERVER"), os.Getenv("OVH_VRACK"))

func TestAccVrackDedicatedServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackDedicatedServerPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackDedicatedServerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_server.vds", "vrack_id", os.Getenv("OVH_VRACK")),
					resource.TestCheckResourceAttr("ovh_vrack_dedicated_server.vds", "server_id", os.Getenv("OVH_DEDICATED_SERVER")),
				),
			},
		},
	})
}

func testAccCheckVrackDedicatedServerPreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccPreCheckDedicatedServer(t)
}
