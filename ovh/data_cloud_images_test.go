package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudImages_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	config := fmt.Sprintf(`
data "ovh_cloud_images" "all" {
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
					resource.TestCheckResourceAttrSet("data.ovh_cloud_images.all", "images.#"),
					// The region filter must yield at least one populated image.
					resource.TestCheckResourceAttrSet("data.ovh_cloud_images.all", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_images.all", "images.0.visibility"),
				),
			},
		},
	})
}
