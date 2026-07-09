package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudFloatingIP_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}

data "ovh_cloud_floating_ip" "test" {
  service_name = ovh_cloud_floating_ip.test.service_name
  id           = ovh_cloud_floating_ip.test.id
}
`, serviceName, region, description)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_floating_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_floating_ip.test", "description", description),
					resource.TestCheckResourceAttr("data.ovh_cloud_floating_ip.test", "location.region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_floating_ip.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_floating_ip.test", "checksum"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_floating_ip.test", "current_state.ip"),
				),
			},
		},
	})
}
