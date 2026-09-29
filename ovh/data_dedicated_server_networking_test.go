package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDedicatedServerNetworking_basic(t *testing.T) {
	dedicatedServer := os.Getenv("OVH_DEDICATED_SERVER")
	config := fmt.Sprintf(testAccDataSourceDedicatedServerNetworking, dedicatedServer, dedicatedServer)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
			testAccPreCheckDedicatedServerNetworking(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server_networking.server", "status", "active"),
					// interfaces are sorted by type, so the public aggregation is index 0
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server_networking.server", "interfaces.0.type", "public"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dedicated_server_networking.server", "interfaces.0.aggregation_fallback"),
				),
			},
		},
	})
}

const testAccDataSourceDedicatedServerNetworking = `
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}

resource "ovh_dedicated_server_networking" "server" {
  service_name = data.ovh_dedicated_server.server.service_name
  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 0, 2)
    type = "public"
  }
  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 2, 4)
    type = "vrack"
  }
}

data "ovh_dedicated_server_networking" "server" {
  service_name = "%s"
  depends_on   = [ovh_dedicated_server_networking.server]
}
`
