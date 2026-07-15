package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudAdditionalIP_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	// Additional IPs cannot be created from Terraform: the test requires an
	// existing additional IP in the test project.
	additionalIP := os.Getenv("OVH_CLOUD_PUBLIC_IP_ADDITIONAL_TEST")
	if additionalIP == "" {
		t.Skip("OVH_CLOUD_PUBLIC_IP_ADDITIONAL_TEST not set")
	}

	config := fmt.Sprintf(`
data "ovh_cloud_additional_ip" "test" {
  service_name = "%s"
  id           = "%s"
}
`, serviceName, additionalIP)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ip.test", "id", additionalIP),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "current_state.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "current_state.ip"),
				),
			},
		},
	})
}

// TestAccDataSourceCloudAdditionalIP_serviceNameFromEnv mirrors the _basic test
// but omits service_name from the data block: the data source must resolve the
// project id from the OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccDataSourceCloudAdditionalIP_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	// Additional IPs cannot be created from Terraform: the test requires an
	// existing additional IP in the test project.
	additionalIP := os.Getenv("OVH_CLOUD_PUBLIC_IP_ADDITIONAL_TEST")
	if additionalIP == "" {
		t.Skip("OVH_CLOUD_PUBLIC_IP_ADDITIONAL_TEST not set")
	}

	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	config := fmt.Sprintf(`
data "ovh_cloud_additional_ip" "test" {
  id = "%s"
}
`, additionalIP)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ip.test", "id", additionalIP),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "current_state.id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ip.test", "current_state.ip"),
				),
			},
		},
	})
}
