package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudPrivateEndpointServiceAvailabilities_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_private_endpoint_service_availabilities" "availabilities" {
  service_name = "%s"
  region       = "%s"
}
`, serviceName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_REGION_TEST")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_private_endpoint_service_availabilities.availabilities", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_private_endpoint_service_availabilities.availabilities", "region", region),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_private_endpoint_service_availabilities.availabilities", "service_availabilities.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_private_endpoint_service_availabilities.availabilities", "service_availabilities.0.tiers.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_private_endpoint_service_availabilities.availabilities", "service_availabilities.0.tiers.0.region"),
				),
			},
		},
	})
}
