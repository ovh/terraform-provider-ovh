package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVrackIpLoadbalancing_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCheckVrackIpLoadbalancingPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVrackIpLoadbalancingConfig,
			},
			{
				ResourceName:      "ovh_vrack_iploadbalancing.viplb",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVrackIpLoadbalancingImportId("ovh_vrack_iploadbalancing.viplb"),
			},
		},
	})
}

func testAccVrackIpLoadbalancingImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		subnet, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("vrack ip_loadbalancing not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s/%s",
			subnet.Primary.Attributes["service_name"],
			subnet.Primary.Attributes["ip_loadbalancing"],
		), nil
	}
}
