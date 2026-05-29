package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSCurrentImageDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSCurrentImageDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/images/current", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_current_image.cur", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_current_image.cur", "id"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_current_image.cur", "name"),
				),
			},
		},
	})
}

const testAccVPSCurrentImageDatasourceConfig_Basic = `
data "ovh_vps_current_image" "cur" {
  service_name = "%s"
}
`
