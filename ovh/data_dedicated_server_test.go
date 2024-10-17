package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerDataSource_basic(t *testing.T) {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	config := fmt.Sprintf(`
	data "ovh_dedicated_server" "server" {
		service_name  = "%s"
	}`, dedicated_server)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedServer(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "name", dedicated_server),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "service_name", dedicated_server),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dedicated_server.server", "vnis.#"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "vnis.0.server_name", dedicated_server),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "boot_script", ""),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dedicated_server.server", "urn"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dedicated_server.server", "display_name"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "region", "ca-east-bhs"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "availability_zone", "ca-east-bhs-a"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "new_upgrade_system", "true"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "no_intervention", "false"),
					resource.TestCheckResourceAttr(
						"data.ovh_dedicated_server.server", "power_state", "poweron"),
				),
			},
		},
	})
}
