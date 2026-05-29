package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPSModelsDataSource_basic(t *testing.T) {
	vps := os.Getenv("OVH_VPS")
	config := fmt.Sprintf(testAccVPSModelsDatasourceConfig, vps)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPS(t)
			skipIfEndpointMissing(t, fmt.Sprintf("/vps/%s/models", os.Getenv("OVH_VPS")))
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_vps_models.test", "models.#"),
				),
			},
		},
	})
}

const testAccVPSModelsDatasourceConfig = `
data "ovh_vps_models" "test" {
  service_name = "%s"
}
`
