package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPFirewall_data(t *testing.T) {
	ip := os.Getenv("OVH_IP_FIREWALL_TEST")
	testAccIPFirewallConfig := fmt.Sprintf(`
		resource "ovh_ip_firewall" "firewall" {
			ip             = "%s"
			ip_on_firewall = "%s"
		}

		data "ovh_ip_firewall" "firewall_data" {
			ip             = ovh_ip_firewall.firewall.ip
			ip_on_firewall = ovh_ip_firewall.firewall.ip_on_firewall
		}
	`, ip, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPFirewallConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall.firewall_data", "ip_on_firewall", ip),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall.firewall_data", "state", "ok"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall.firewall_data", "enabled", "false"),
				),
			},
		},
	})
}
