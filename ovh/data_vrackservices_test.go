package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVrackservicesDatasourceConfig_Basic = `
data "ovh_vrackservices" "vrackservices" {
  vrack_services_id  = "%s"
}
`

func TestAccVrackservicesDataSource_basic(t *testing.T) {
	id := os.Getenv("OVH_TESTACC_VRACK_SERVICES_ID_TEST")
	config := fmt.Sprintf(testAccVrackservicesDatasourceConfig_Basic, id)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVrackServicesData(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vrackservices.vrackservices", "id", id),
					resource.TestCheckResourceAttr("data.ovh_vrackservices.vrackservices", "vrack_services_id", id),
					resource.TestCheckResourceAttr("data.ovh_vrackservices.vrackservices", "resource_status", "READY"),
				),
			},
		},
	})
}
