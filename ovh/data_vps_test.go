package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			checkEnvOrSkip(t, "OVH_VPS")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVPSDataSourceConfig, vps),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vps.vps", "id", vps),
					resource.TestCheckResourceAttrSet("data.ovh_vps.vps", "display_name"),
					resource.TestCheckResourceAttrSet("data.ovh_vps.vps", "current_state.power_state"),
					resource.TestCheckResourceAttrSet("data.ovh_vps.vps", "current_state.location.region"),
					resource.TestCheckResourceAttrSet("data.ovh_vps.vps", "current_state.ip_addresses.#"),
				),
			},
		},
	})
}

const testAccVPSDataSourceConfig = `
data "ovh_vps" "vps" {
  id = "%s"
}
`
