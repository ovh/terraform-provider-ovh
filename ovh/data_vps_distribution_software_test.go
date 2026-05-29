package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDistributionSoftwareDatasourceConfig_Basic = `
data "ovh_vps_distribution_software" "installed" {
  service_name = "%s"
}
`

func TestAccVPSDistributionSoftwareDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSDistributionSoftwareDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/distribution/software", os.Getenv("OVH_VPS")))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_distribution_software.installed", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution_software.installed", "software_ids.#"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution_software.installed", "software.#"),
				),
			},
		},
	})
}
