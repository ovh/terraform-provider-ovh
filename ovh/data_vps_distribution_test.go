package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDistributionDatasourceConfig_Basic = `
data "ovh_vps_distribution" "current" {
  service_name = "%s"
}
`

func TestAccVPSDistributionDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSDistributionDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/distribution", os.Getenv("OVH_VPS")))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_distribution.current", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution.current", "name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution.current", "distribution"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution.current", "bit_format"),
				),
			},
		},
	})
}
