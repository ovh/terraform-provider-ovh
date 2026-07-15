package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudAdditionalIPs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	// Additional IPs cannot be created from Terraform: this test only checks
	// that the listing works, even when the project has no additional IP.
	config := fmt.Sprintf(`
data "ovh_cloud_additional_ips" "test" {
  service_name = "%s"
}
`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ips.test", "additional_ips.#"),
				),
			},
		},
	})
}

// TestAccDataSourceCloudAdditionalIPs_serviceNameFromEnv mirrors the _basic test
// but omits service_name from the data block: the data source must resolve the
// project id from the OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccDataSourceCloudAdditionalIPs_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	// Additional IPs cannot be created from Terraform: this test only checks
	// that the listing works, even when the project has no additional IP.
	config := `
data "ovh_cloud_additional_ips" "test" {
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudPublicIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_additional_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_additional_ips.test", "additional_ips.#"),
				),
			},
		},
	})
}
