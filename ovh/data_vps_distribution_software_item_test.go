package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVPSDistributionSoftwareItemDatasourceConfig_Basic = `
data "ovh_vps_distribution_software" "installed" {
  service_name = "%s"
}

data "ovh_vps_distribution_software_item" "first" {
  service_name = "%s"
  software_id  = data.ovh_vps_distribution_software.installed.software_ids[0]
}
`

func TestAccVPSDistributionSoftwareItemDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSDistributionSoftwareItemDatasourceConfig_Basic, vps, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/distribution/software/1", os.Getenv("OVH_VPS")))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vps_distribution_software_item.first", "service_name", vps),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution_software_item.first", "name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution_software_item.first", "type"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_distribution_software_item.first", "status"),
				),
			},
		},
	})
}
