package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPFirewallRule_data(t *testing.T) {
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

		data "ovh_ip_firewall_rule" "rule_data" {
			ip = ovh_ip_firewall_rule.rule.ip
			ip_on_firewall = ovh_ip_firewall_rule.rule.ip_on_firewall
			sequence = 0
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
						"data.ovh_ip_firewall_rule.rule_data", "ip_on_firewall", ip),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "state", "ok"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "action", "permit"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "source", "any"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "tcp_option", "established"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "destination_port", "eq 22"),
					resource.TestCheckResourceAttr(
						"data.ovh_ip_firewall_rule.rule_data", "source_port", "eq 44"),
				),
			},
		},
	})
}
