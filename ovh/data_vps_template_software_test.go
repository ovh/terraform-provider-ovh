package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSTemplateSoftwareDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	tplID := os.Getenv("OVH_VPS_TEMPLATE_ID")
	swID := os.Getenv("OVH_VPS_SOFTWARE_ID")
	if vps == "" || tplID == "" || swID == "" {
		t.Skip("OVH_VPS, OVH_VPS_TEMPLATE_ID, or OVH_VPS_SOFTWARE_ID not set")
	}
	config := fmt.Sprintf(testAccVPSTemplateSoftwareDatasourceConfig_Basic, vps, tplID, swID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/templates/1/software", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_template_software.sw", "name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_template_software.sw", "type"),
				),
			},
		},
	})
}

const testAccVPSTemplateSoftwareDatasourceConfig_Basic = `
data "ovh_vps_template_software" "sw" {
  service_name = "%s"
  template_id  = %s
  software_id  = %s
}
`
