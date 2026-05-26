package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSTemplatesDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	if vps == "" {
		t.Skip("OVH_VPS not set")
	}
	config := fmt.Sprintf(testAccVPSTemplatesDatasourceConfig_Basic, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/templates", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_templates.tpls", "template_ids.#"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_templates.tpls", "templates.#"),
				),
			},
		},
	})
}

const testAccVPSTemplatesDatasourceConfig_Basic = `
data "ovh_vps_templates" "tpls" {
  service_name = "%s"
}
`
