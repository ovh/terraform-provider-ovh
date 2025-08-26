package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_ip_firewall", &resource.Sweeper{
		Name: "ovh_ip_firewall",
		F:    testSweepIPFirewall,
	})
}

func testSweepIPFirewall(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	ip := os.Getenv("OVH_IP_FIREWALL_TEST")
	endpoint := fmt.Sprintf("/ip/%s/firewall", url.PathEscape(ip))

	var ips []string
	if err := client.Get(endpoint, &ips); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	for _, ipOnFirewall := range ips {
		if err := client.Delete(fmt.Sprintf("/ip/%s/firewall/%s", url.PathEscape(ip), url.PathEscape(ipOnFirewall)), nil); err != nil {
			return fmt.Errorf("Error deleting ip firewall %s on ip %s:\n\t %q", ipOnFirewall, ip, err)
		}
	}

	return nil
}

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
	ip := os.Getenv("OVH_IP_FIREWALL_TEST")

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
