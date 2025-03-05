package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVrackDataSource_basic(t *testing.T) {
	vrack := os.Getenv("OVH_VRACK_SERVICE_TEST")
	config := fmt.Sprintf(testAccVrackDatasourceConfig_Basic, vrack)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckVrack(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_vrack.vrack", "id", vrack),
					resource.TestCheckResourceAttr(
						"data.ovh_vrack.vrack", "service_name", vrack),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vrack.vrack", "iam"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vrack.vrack", "name"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_vrack.vrack", "description"),
				),
			},
		},
	})
}

const testAccVrackDatasourceConfig_Basic = `
data "ovh_vrack" "vrack" {
  service_name  = "%s"
}
`
