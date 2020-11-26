package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIpLoadbalancingVrackNetwork_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackIpLoadbalancingPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpLoadbalancingVrackNetworkConfig_basic,
			},
			{
				ResourceName:      "ovh_iploadbalancing_vrack_network.network",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIpLoadbalancingVrackNetworkImportId("ovh_iploadbalancing_vrack_network.network"),
			},
		},
	})
}

func testAccIpLoadbalancingVrackNetworkImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		subnet, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_ip_loadbalancing_vrack_network not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s/%s",
			subnet.Primary.Attributes["service_name"],
			subnet.Primary.Attributes["vrack_network_id"],
		), nil
	}
}
