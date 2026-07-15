package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudPublicIPs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}

data "ovh_cloud_public_ips" "test" {
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
					resource.TestCheckResourceAttr("data.ovh_cloud_public_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrWith("data.ovh_cloud_public_ips.test", "public_ips.#", testAccCheckCloudPublicIPListNotEmpty),
				),
			},
		},
	})
}

// TestAccDataSourceCloudPublicIPs_serviceNameFromEnv mirrors the _basic test but
// omits service_name from the data block: the data source must resolve the
// project id from the OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccDataSourceCloudPublicIPs_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	description := acctest.RandomWithPrefix(testAccResourceCloudFloatingIPDescriptionPrefix)

	// The resource keeps an explicit service_name so setup stays deterministic;
	// only the data block relies on the environment fallback.
	config := fmt.Sprintf(`
resource "ovh_cloud_floating_ip" "test" {
  service_name = "%s"
  region       = "%s"
  description  = "%s"
}

data "ovh_cloud_public_ips" "test" {
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
					resource.TestCheckResourceAttr("data.ovh_cloud_public_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrWith("data.ovh_cloud_public_ips.test", "public_ips.#", testAccCheckCloudPublicIPListNotEmpty),
				),
			},
		},
	})
}
