package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudNetworkPrivateSubnet_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccCheckcCloudNetworkPrivateSubnetPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkPrivateSubnetConfig(),
			},
			{
				ResourceName:      "ovh_cloud_network_private_subnet.subnet",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateSubnetImportId("ovh_cloud_network_private_subnet.subnet"),
			},
		},
	})
}

func testAccCloudNetworkPrivateSubnetImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		subnet, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("subnet not found: %s", resourceName)
		}

		return fmt.Sprintf(
			"%s/%s/%s",
			subnet.Primary.Attributes["project_id"],
			subnet.Primary.Attributes["network_id"],
			subnet.Primary.ID,
		), nil
	}
}
