package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudFlavors_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_INSTANCE_REGION_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_flavors" "all" {
  service_name = "%s"
  region       = "%s"
}
`, serviceName, region)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudInstance(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavors.all", "flavors.#"),
					// The region filter must yield at least one populated flavor.
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavors.all", "flavors.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_flavors.all", "flavors.0.vcpus"),
				),
			},
		},
	})
}
