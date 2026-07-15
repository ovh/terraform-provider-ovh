package ovh

import (
	"fmt"
	"os"
	"regexp"
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

// TestAccDataSourceCloudFloatingIP_serviceNameFromEnv mirrors the _basic test
// but omits service_name from the data block: the data source must resolve the
// project id from the OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccDataSourceCloudFloatingIP_serviceNameFromEnv(t *testing.T) {
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

data "ovh_cloud_floating_ip" "test" {
  id = ovh_cloud_floating_ip.test.id
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

// TestAccDataSourceCloudFloatingIP_missingServiceName asserts that the data
// source errors clearly when neither the configuration nor the
// OVH_CLOUD_PROJECT_SERVICE environment variable supplies service_name. The
// error is raised before any API call, so no project/region is needed.
func TestAccDataSourceCloudFloatingIP_missingServiceName(t *testing.T) {
	// Force the env var empty so the inline fallback cannot resolve it.
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "")

	// 203.0.113.1 is a documentation/test IP (RFC 5737); it is never queried
	// because the missing service_name is rejected first.
	config := `
data "ovh_cloud_floating_ip" "test" {
  id = "203.0.113.1"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Missing service_name"),
			},
		},
	})
}
