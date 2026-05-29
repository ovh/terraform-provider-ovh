package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSTemplateDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	tplID := os.Getenv("OVH_VPS_TEMPLATE_ID")
	if vps == "" || tplID == "" {
		t.Skip("OVH_VPS or OVH_VPS_TEMPLATE_ID not set")
	}
	config := fmt.Sprintf(testAccVPSTemplateDatasourceConfig_Basic, vps, tplID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/templates/1", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_template.tpl", "name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_template.tpl", "distribution"),
				),
			},
		},
	})
}

const testAccVPSTemplateDatasourceConfig_Basic = `
data "ovh_vps_template" "tpl" {
  service_name = "%s"
  template_id  = %s
}
`
