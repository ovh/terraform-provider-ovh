package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccIPMitigationConfig = `
resource "ovh_ip_mitigation" "mitigation" {
	ip               = "%s"
	ip_on_mitigation = "%s"
}
`

const testAccIPMitigationUpdatedConfig = `
resource "ovh_ip_mitigation" "mitigation" {
	ip               = "%s"
	ip_on_mitigation = "%s"
	permanent        = false
}
`

func TestAccIPMitigation_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_TEST")

	config := fmt.Sprintf(testAccIPMitigationConfig, ip, ip)
	updatedConfig := fmt.Sprintf(testAccIPMitigationUpdatedConfig, ip, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "ip_on_mitigation", ip),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "auto", "false"),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "permanent", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "ip_on_mitigation", ip),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "auto", "false"),
					resource.TestCheckResourceAttr(
						"ovh_ip_mitigation.mitigation", "permanent", "false"),
				),
			},
		},
	})
}
