package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDedicatedServerVirtualMac_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
			checkEnvOrSkip(t, "OVH_DEDICATED_SERVER_VMAC_IP_TEST")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerVirtualMacConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_virtual_mac.vmac", "type", "ovh"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_virtual_mac.vmac", "virtual_machine_name", "terraform-acc-vmac"),
					resource.TestCheckResourceAttrSet(
						"ovh_dedicated_server_virtual_mac.vmac", "mac_address"),
				),
			},
			{
				ResourceName:      "ovh_dedicated_server_virtual_mac.vmac",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDedicatedServerVirtualMacConfig() string {
	return fmt.Sprintf(
		testAccDedicatedServerVirtualMacConfig_Basic,
		os.Getenv("OVH_DEDICATED_SERVER"),
		os.Getenv("OVH_DEDICATED_SERVER_VMAC_IP_TEST"),
	)
}

const testAccDedicatedServerVirtualMacConfig_Basic = `
resource "ovh_dedicated_server_virtual_mac" "vmac" {
  service_name         = "%s"
  ip_address           = "%s"
  type                 = "ovh"
  virtual_machine_name = "terraform-acc-vmac"
}
`
