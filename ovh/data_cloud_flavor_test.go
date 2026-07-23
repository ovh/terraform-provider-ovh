package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudFlavor_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)

	config := fmt.Sprintf(`
data "ovh_cloud_flavor" "test" {
  service_name = "%s"
  id           = "%s"
}
`, serviceName, flavorID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_flavor.test", "id", flavorID),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavor.test", "name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavor.test", "vcpus"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavor.test", "ram"),
				),
			},
		},
	})
}
