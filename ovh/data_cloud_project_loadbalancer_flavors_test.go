package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDataSourceCloudProjectLoadbalancerFlavors_serviceNameFromEnv omits
// service_name from the data block: the data source must resolve the project id
// from the OVH_CLOUD_PROJECT_SERVICE environment variable and store it in state.
func TestAccDataSourceCloudProjectLoadbalancerFlavors_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	config := fmt.Sprintf(`
data "ovh_cloud_project_loadbalancer_flavors" "test" {
  region_name = "%s"
}
`, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudRegion(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_loadbalancer_flavors.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_loadbalancer_flavors.test", "region_name", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_loadbalancer_flavors.test", "flavors.0.id"),
				),
			},
		},
	})
}
