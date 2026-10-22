package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPFirewallRule_basic(t *testing.T) {
	ip := os.Getenv("OVH_IP_FIREWALL_TEST")
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

func TestParsePortRange(t *testing.T) {
	tests := []struct {
		desc     string
		from, to int64
		wantErr  bool
	}{
		{desc: "range 443 445", from: 443, to: 445},
		{desc: "range 1 65535", from: 1, to: 65535},
		{desc: "eq 80", wantErr: true},
		{desc: "", wantErr: true},
		{desc: "range 443", wantErr: true},
		{desc: "range a b", wantErr: true},
	}

	for _, tt := range tests {
		got, err := parsePortRange(tt.desc)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parsePortRange(%q): expected an error, got %v", tt.desc, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parsePortRange(%q): unexpected error: %s", tt.desc, err)
			continue
		}
		if got.From.ValueInt64() != tt.from || got.To.ValueInt64() != tt.to {
			t.Errorf("parsePortRange(%q) = {%d, %d}, expected {%d, %d}",
				tt.desc, got.From.ValueInt64(), got.To.ValueInt64(), tt.from, tt.to)
		}
	}
}

func TestParsePortDesc(t *testing.T) {
	tests := []struct {
		desc      string
		port      *int64
		portRange *[2]int64
		wantErr   bool
	}{
		{desc: ""},
		{desc: "eq 22", port: func(v int64) *int64 { return &v }(22)},
		{desc: "range 8000 8010", portRange: &[2]int64{8000, 8010}},
		{desc: "eq", wantErr: true},
		{desc: "eq abc", wantErr: true},
		{desc: "gt 1024", wantErr: true},
		{desc: "range 8000", wantErr: true},
	}

	for _, tt := range tests {
		port, portRange, err := parsePortDesc(tt.desc)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parsePortDesc(%q): expected an error", tt.desc)
			}
			continue
		}
		if err != nil {
			t.Errorf("parsePortDesc(%q): unexpected error: %s", tt.desc, err)
			continue
		}
		if tt.port == nil && !port.IsNull() {
			t.Errorf("parsePortDesc(%q): port = %d, expected null", tt.desc, port.ValueInt64())
		}
		if tt.port != nil && (port.IsNull() || port.ValueInt64() != *tt.port) {
			t.Errorf("parsePortDesc(%q): port = %v, expected %d", tt.desc, port, *tt.port)
		}
		if tt.portRange == nil && portRange != nil {
			t.Errorf("parsePortDesc(%q): range = %v, expected none", tt.desc, portRange)
		}
		if tt.portRange != nil && (portRange == nil || portRange.From.ValueInt64() != tt.portRange[0] || portRange.To.ValueInt64() != tt.portRange[1]) {
			t.Errorf("parsePortDesc(%q): range = %v, expected %v", tt.desc, portRange, *tt.portRange)
		}
	}
}

func testAccIPFirewallRulePortRangeConfig(ip, rule string) string {
	return fmt.Sprintf(`
		resource "ovh_ip_firewall" "firewall" {
			ip             = "%s"
			ip_on_firewall = "%s"
		}

		resource "ovh_ip_firewall_rule" "rule" {
			ip             = ovh_ip_firewall.firewall.ip
			ip_on_firewall = ovh_ip_firewall.firewall.ip_on_firewall
			sequence       = 0
			action         = "permit"
			%s
		}
	`, ip, ip, rule)
}

func TestAccIPFirewallRule_portRange(t *testing.T) {
	ip := os.Getenv("OVH_IP_FIREWALL_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckIp(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPFirewallRulePortRangeConfig(ip, `
					protocol               = "tcp"
					destination_port_range = { from = 443, to = 443 }
				`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("to must be greater than from"),
			},
			{
				Config: testAccIPFirewallRulePortRangeConfig(ip, `
					protocol               = "icmp"
					destination_port_range = { from = 1000, to = 2000 }
				`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid protocol for port range"),
			},
			{
				Config: testAccIPFirewallRulePortRangeConfig(ip, `
					protocol   = "udp"
					tcp_option = "established"
				`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid protocol for tcp_option"),
			},
			{
				Config: testAccIPFirewallRulePortRangeConfig(ip, `
					protocol  = "icmp"
					fragments = false
				`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid protocol for fragments"),
			},
			{
				Config: testAccIPFirewallRulePortRangeConfig(ip, `
					protocol               = "tcp"
					source_port            = 44
					destination_port_range = { from = 1000, to = 2000 }
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "state", "ok"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port_range.from", "1000"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port_range.to", "2000"),
					resource.TestCheckResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port_desc", "range 1000 2000"),
					resource.TestCheckNoResourceAttr(
						"ovh_ip_firewall_rule.rule", "destination_port"),
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
