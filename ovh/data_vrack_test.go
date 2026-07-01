package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVrackDatasourceConfig_Basic = `
data "ovh_vrack" "vrack" {
  service_name  = "%s"
}
`

func TestAccDatasourceVrack_basic(t *testing.T) {
	vrackName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	config := fmt.Sprintf(testAccVrackDatasourceConfig_Basic, vrackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVrackData(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vrack.vrack", "service_name", vrackName),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vrack.vrack", "iam.urn"),
				),
			},
		},
	})
}
