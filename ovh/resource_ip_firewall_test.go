package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccIPFirewallConfig = `
resource "ovh_ip_firewall" "firewall" {
	ip             = "%s"
	ip_on_firewall = "%s"
}
`

const testAccIPFirewallUpdatedConfig = `
resource "ovh_ip_firewall" "firewall" {
	ip             = "%s"
	ip_on_firewall = "%s"
	enabled        = true
}
`

func TestAccIPFirewall_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_TEST")

	config := fmt.Sprintf(testAccIPFirewallConfig, ip, ip)
	updatedConfig := fmt.Sprintf(testAccIPFirewallUpdatedConfig, ip, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "ip_on_firewall", ip),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "enabled", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "ip_on_firewall", ip),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall.firewall", "enabled", "true"),
				),
			},
		},
	})
}
