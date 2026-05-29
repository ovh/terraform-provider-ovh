package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSDatacentersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			skipIfEndpointMissing(t, "/vps/datacenter")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPSDatacentersDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_datacenters.all", "datacenters.#"),
				),
			},
		},
	})
}

const testAccVPSDatacentersDatasourceConfig = `
data "ovh_vps_datacenters" "all" {}
`
