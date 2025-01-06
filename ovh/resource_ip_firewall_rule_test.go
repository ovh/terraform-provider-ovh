package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPFirewallRule_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_TEST")
	testAccIPFirewallRuleConfig := fmt.Sprintf(`
		resource "ovh_ip_firewall" "firewall" {
			ip             = "%s"
			ip_on_firewall = "%s"
		}

		data "ovh_ip_firewall" "firewall_data" {
			ip             = ovh_ip_firewall.firewall.ip
			ip_on_firewall = ovh_ip_firewall.firewall.ip_on_firewall
		}

		resource "ovh_ip_firewall_rule" "rule" {
			ip = data.ovh_ip_firewall.firewall_data.ip
			ip_on_firewall = data.ovh_ip_firewall.firewall_data.ip_on_firewall

			action = "permit"
			protocol = "tcp"
			sequence = 0
			tcp_option = "established"
			destination_port = 22
			source_port = 44
		}
	`, ip, ip)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPFirewallRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "ip_on_firewall", ip),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "action", "permit"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "source", "any"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "tcp_option", "established"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port", "22"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port_desc", "eq 22"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "source_port", "44"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "source_port_desc", "eq 44"),
				),
			},
			{
				ImportStateId:                        fmt.Sprintf("%s|%s|0", ip, ip),
				ResourceName:                         "ovh_ip_firewall_rule.rule",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "ip_on_firewall",
			},
		},
	})
}
