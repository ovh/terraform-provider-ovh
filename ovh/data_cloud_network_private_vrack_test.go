package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudNetworkPrivateVrack_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	networkName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateVrackNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "test" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}

data "ovh_cloud_network_private_vrack" "test" {
  service_name = ovh_cloud_network_private_vrack.test.service_name
  id           = ovh_cloud_network_private_vrack.test.id
}
`, serviceName, networkName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateVrack(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack.test", "name", networkName),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack.test", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack.test", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_network_private_vrack.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_network_private_vrack.test", "current_state.location.region"),
				),
			},
		},
	})
}
