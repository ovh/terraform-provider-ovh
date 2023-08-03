package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccresourceDedicatedServerNetworking(t *testing.T) {
	confTemplates := []string{
		testAccDedicatedServerNetworkingPublicVrack,
		testAccDedicatedServerNetworkingSingleVrack,
		testAccDedicatedServerNetworkingDualVrack,
	}

	for _, confTemplate := range confTemplates {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheckCredentials(t)
				testAccPreCheckDedicatedServer(t)
				testAccPreCheckDedicatedServerNetworking(t)
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDedicatedServerNetworkingConfig(confTemplate),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"ovh_dedicated_server_networking.server", "status", "active"),
					),
				},
			},
		})
	}
}

func testAccDedicatedServerNetworkingConfig(confTemplate string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	renderedConfig := fmt.Sprintf(
		confTemplate,
		dedicated_server,
	)
	return renderedConfig
}

const testAccDedicatedServerNetworkingPublicVrack = `
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
`

const testAccDedicatedServerNetworkingSingleVrack = `
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}

resource "ovh_dedicated_server_networking" "server" {
  service_name = data.ovh_dedicated_server.server.service_name
  interfaces {
    macs = flatten(data.ovh_dedicated_server.server.vnis.*.nics)
    type = "vrack"
  }
}
`

const testAccDedicatedServerNetworkingDualVrack = `
data "ovh_dedicated_server" "server" {
  service_name = "%s"
}
resource "ovh_dedicated_server_networking" "server" {
  service_name = data.ovh_dedicated_server.server.service_name
  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 0, 2)
    type = "vrack"
  }
  interfaces {
    macs = slice(sort(flatten(data.ovh_dedicated_server.server.vnis.*.nics)), 2, 4)
    type = "vrack"
  }
}
`
