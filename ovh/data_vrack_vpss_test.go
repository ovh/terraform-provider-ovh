package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVrackVpss_basic(t *testing.T) {
	config := fmt.Sprintf(`
data "ovh_vrack_vpss" "vrackvpss" {
  service_name = "%s"
}
`, os.Getenv("OVH_VRACK_SERVICE_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVRack(t); testAccCheckVRackExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttrSet(
					"data.ovh_vrack_vpss.vrackvpss",
					"result.#",
				),
			},
		},
	})
}
