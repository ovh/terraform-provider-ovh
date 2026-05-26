package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVrackVps_basic(t *testing.T) {
	config := fmt.Sprintf(`
resource "ovh_vrack_vps" "vrackvps" {
  service_name     = "%s"
  vps_service_name = "%s"
}

data "ovh_vrack_vps" "vrackvps" {
  service_name     = ovh_vrack_vps.vrackvps.service_name
  vps_service_name = ovh_vrack_vps.vrackvps.vps_service_name
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VPS"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackVpsPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vrack_vps.vrackvps", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr("data.ovh_vrack_vps.vrackvps", "vps_service_name", os.Getenv("OVH_VPS")),
				),
			},
		},
	})
}
