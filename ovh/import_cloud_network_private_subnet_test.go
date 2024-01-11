package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudProjectNetworkPrivateSubnet_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccCheckcCloudProjectNetworkPrivateSubnetPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectNetworkPrivateSubnetConfig(testAccCloudProjectNetworkPrivateSubnetConfig_basic),
			},
			{
				ResourceName:      "ovh_cloud_project_network_private_subnet.subnet",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudProjectNetworkPrivateSubnetImportId("ovh_cloud_project_network_private_subnet.subnet"),
			},
		},
	})
}

func testAccCloudProjectNetworkPrivateSubnetImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		subnet, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("subnet not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s/%s/%s",
			subnet.Primary.Attributes["service_name"],
			subnet.Primary.Attributes["network_id"],
			subnet.Primary.ID,
		), nil
	}
}
