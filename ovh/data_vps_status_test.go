package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSStatusConfig = `
data "ovh_vps_status" "status" {
  service_name = "%s"
}
`

func TestAccVPSStatusDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	if vps == "" {
		t.Skip("OVH_VPS must be set for this acceptance test")
	}
	config := fmt.Sprintf(testAccVPSStatusConfig, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/status", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_status.status", "probes.#"),
				),
			},
		},
	})
}
