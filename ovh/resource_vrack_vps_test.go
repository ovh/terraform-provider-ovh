package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccVrackVpsConfig = fmt.Sprintf(`
resource "ovh_vrack_vps" "vrackvps" {
  service_name     = "%s"
  vps_service_name = "%s"
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VPS"))

func TestAccVrackVps_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackVpsPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackVpsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_vps.vrackvps", "service_name", os.Getenv("OVH_VRACK_SERVICE_TEST")),
					resource.TestCheckResourceAttr("ovh_vrack_vps.vrackvps", "vps_service_name", os.Getenv("OVH_VPS")),
				),
			},
		},
	})
}

func testAccCheckVrackVpsPreCheck(t *testing.T) {
	testAccPreCheckVRack(t)
	testAccCheckVRackExists(t)
	testAccPreCheckVPS(t)
}
