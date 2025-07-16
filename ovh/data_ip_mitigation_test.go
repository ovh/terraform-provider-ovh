package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccIPMitigationDataSourceConfig = `
resource "ovh_ip_mitigation" "mitigation" {
	ip               = "%s"
	ip_on_mitigation = "%s"
}

data "ovh_ip_mitigation" "mitigation_data" {
	ip               = ovh_ip_mitigation.mitigation.ip
	ip_on_mitigation = ovh_ip_mitigation.mitigation.ip_on_mitigation
}
`

func TestAccIPMitigationData_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_TEST")

	config := fmt.Sprintf(testAccIPMitigationDataSourceConfig, ip, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_ip_mitigation.mitigation_data", "ip_on_mitigation", ip),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_mitigation.mitigation_data", "state", "ok"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_mitigation.mitigation_data", "auto", "false"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_mitigation.mitigation_data", "permanent", "true"),
					resource.TestCheckResourceAttrSet("data.ovh_ip_mitigation.mitigation_data", "id"),
				),
			},
		},
	})
}
