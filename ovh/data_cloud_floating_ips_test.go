package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudFloatingIPs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}

data "ovh_cloud_floating_ips" "test" {
  service_name = ovh_cloud_floating_ip.test.service_name

  depends_on = [ovh_cloud_floating_ip.test]
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
					resource.TestCheckResourceAttr("data.ovh_cloud_floating_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrWith("data.ovh_cloud_floating_ips.test", "floating_ips.#", testAccCheckCloudPublicIPListNotEmpty),
				),
			},
		},
	})
}
